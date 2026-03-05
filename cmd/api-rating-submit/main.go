package main

import (
	"log/slog"

	"github.com/aws/aws-lambda-go/lambda"

	"github.com/davidlramirez95/towcommand/internal/adapter/gateway"
	"github.com/davidlramirez95/towcommand/internal/adapter/handler"
	"github.com/davidlramirez95/towcommand/internal/adapter/repository"
	"github.com/davidlramirez95/towcommand/internal/platform/awsclient"
	"github.com/davidlramirez95/towcommand/internal/platform/config"
	"github.com/davidlramirez95/towcommand/internal/platform/logger"
	ratinguc "github.com/davidlramirez95/towcommand/internal/usecase/rating"
)

func main() {
	cfg := config.Load()
	log := logger.New(cfg.Stage, cfg.FunctionName, cfg.FunctionVersion, cfg.LogLevel)
	slog.SetDefault(log)

	ddb := awsclient.DynamoDBClient(cfg)
	eb := awsclient.EventBridgeClient(cfg)

	ratingRepo := repository.NewRatingRepository(ddb, cfg.DynamoDBTable)
	providerRepo := repository.NewProviderRepository(ddb, cfg.DynamoDBTable)
	bookingRepo := repository.NewBookingRepository(ddb, cfg.DynamoDBTable)
	events := gateway.NewEventBridgePublisher(eb, cfg.EventBusName, log)

	uc := ratinguc.NewSubmitRatingUseCase(ratingRepo, ratingRepo, ratingRepo, bookingRepo, providerRepo, providerRepo, events)
	h := handler.NewSubmitRatingHandler(uc)

	lambda.Start(handler.WithRecover(handler.WithLogging(handler.WithCorrelationID(h.Handle))))
}
