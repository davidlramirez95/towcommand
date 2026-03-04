package main

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/service/s3"

	"github.com/davidlramirez95/towcommand/internal/adapter/gateway"
	"github.com/davidlramirez95/towcommand/internal/adapter/handler"
	"github.com/davidlramirez95/towcommand/internal/adapter/repository"
	domainerrors "github.com/davidlramirez95/towcommand/internal/domain/errors"
	"github.com/davidlramirez95/towcommand/internal/platform/awsclient"
	"github.com/davidlramirez95/towcommand/internal/platform/config"
	"github.com/davidlramirez95/towcommand/internal/platform/logger"
	evidenceuc "github.com/davidlramirez95/towcommand/internal/usecase/evidence"
)

func main() {
	cfg := config.Load()
	log := logger.New(cfg.Stage, cfg.FunctionName, cfg.FunctionVersion, cfg.LogLevel)
	slog.SetDefault(log)

	s3Client := awsclient.S3Client(cfg)
	presignClient := s3.NewPresignClient(s3Client)
	ddb := awsclient.DynamoDBClient(cfg)

	presigner := gateway.NewS3EvidenceAdapter(presignClient, cfg.S3Bucket)
	bookingRepo := repository.NewBookingRepository(ddb, cfg.DynamoDBTable)
	uc := evidenceuc.NewGenerateUploadURLUseCase(bookingRepo, presigner)

	h := func(ctx context.Context, event *events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
		userID := handler.ExtractUserID(event)
		if userID == "" {
			return handler.ErrorResponse(domainerrors.NewUnauthorizedError()), nil
		}

		body, err := handler.ParseBody[evidenceuc.GenerateUploadURLInput](event)
		if err != nil {
			return handler.ErrorResponse(err), nil
		}

		result, err := uc.Execute(ctx, &body)
		if err != nil {
			return handler.ErrorResponse(err), nil
		}

		return handler.SuccessResponse(http.StatusOK, result), nil
	}

	lambda.Start(handler.WithRecover(handler.WithLogging(handler.WithCorrelationID(h))))
}
