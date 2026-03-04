// Package main is the composition root for the OTP generation Lambda.
package main

import (
	"context"
	"log/slog"
	"net/http"

	lambdaevents "github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"

	"github.com/davidlramirez95/towcommand/internal/adapter/cache"
	"github.com/davidlramirez95/towcommand/internal/adapter/gateway"
	"github.com/davidlramirez95/towcommand/internal/adapter/handler"
	"github.com/davidlramirez95/towcommand/internal/adapter/repository"
	domainerrors "github.com/davidlramirez95/towcommand/internal/domain/errors"
	"github.com/davidlramirez95/towcommand/internal/domain/otp"
	"github.com/davidlramirez95/towcommand/internal/platform/awsclient"
	"github.com/davidlramirez95/towcommand/internal/platform/config"
	"github.com/davidlramirez95/towcommand/internal/platform/logger"
	otpuc "github.com/davidlramirez95/towcommand/internal/usecase/otp"
)

// generateOTPRequest is the expected JSON body for POST /bookings/{id}/otp/generate.
type generateOTPRequest struct {
	OTPType string  `json:"otpType" validate:"required,oneof=PICKUP DROPOFF"`
	Lat     float64 `json:"lat" validate:"required,latitude"`
	Lng     float64 `json:"lng" validate:"required,longitude"`
}

func main() {
	cfg := config.Load()
	log := logger.New(cfg.Stage, cfg.FunctionName, cfg.FunctionVersion, cfg.LogLevel)
	slog.SetDefault(log)

	ddb := awsclient.DynamoDBClient(cfg)
	eb := awsclient.EventBridgeClient(cfg)
	snsClient := awsclient.SNSClient(cfg)

	redisClient := cache.NewRedisClient(cache.Options{
		Host:     cfg.RedisHost,
		Port:     cfg.RedisPort,
		Password: cfg.RedisPassword,
	})

	bookingRepo := repository.NewBookingRepository(ddb, cfg.DynamoDBTable)
	userRepo := repository.NewUserRepository(ddb, cfg.DynamoDBTable)
	otpRepo := repository.NewOTPRepository(ddb, cfg.DynamoDBTable)
	otpCache := cache.NewRedisOTPCache(redisClient)
	rateLimiter := cache.NewRedisRateLimiter(redisClient)
	smsSender := gateway.NewSNSNotificationSender(snsClient)
	pub := gateway.NewEventBridgePublisher(eb, cfg.EventBusName, log)

	uc := otpuc.NewGenerateOTPUseCase(bookingRepo, userRepo, otpCache, otpRepo, rateLimiter, smsSender, pub, log)

	h := handler.WithRecover(handler.WithCorrelationID(handler.WithLogging(
		func(ctx context.Context, event *lambdaevents.APIGatewayProxyRequest) (lambdaevents.APIGatewayProxyResponse, error) {
			userID := handler.ExtractUserID(event)
			if userID == "" {
				return handler.ErrorResponse(domainerrors.NewUnauthorizedError()), nil
			}

			bookingID := handler.ParsePathParam(event, "id")
			if bookingID == "" {
				return handler.ErrorResponse(domainerrors.NewValidationError("missing path parameter: id")), nil
			}

			body, err := handler.ParseBody[generateOTPRequest](event)
			if err != nil {
				return handler.ErrorResponse(err), nil
			}

			result, err := uc.Execute(ctx, &otpuc.GenerateOTPInput{
				BookingID: bookingID,
				OTPType:   otp.OTPType(body.OTPType),
				Lat:       body.Lat,
				Lng:       body.Lng,
				CallerID:  userID,
			})
			if err != nil {
				return handler.ErrorResponse(err), nil
			}

			return handler.SuccessResponse(http.StatusOK, result), nil
		},
	)))

	lambda.Start(h)
}
