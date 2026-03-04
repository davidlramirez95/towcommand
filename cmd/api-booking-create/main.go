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
	bookinguc "github.com/davidlramirez95/towcommand/internal/usecase/booking"
)

func main() {
	cfg := config.Load()
	log := logger.New(cfg.Stage, cfg.FunctionName, cfg.FunctionVersion, cfg.LogLevel)
	slog.SetDefault(log)

	ddb := awsclient.DynamoDBClient(cfg)
	eb := awsclient.EventBridgeClient(cfg)

	repo := repository.NewBookingRepository(ddb, cfg.DynamoDBTable)
	events := gateway.NewEventBridgePublisher(eb, cfg.EventBusName, log)
	uc := bookinguc.NewCreateBookingUseCase(repo, events)
	h := handler.NewCreateBookingHandler(uc)

	lambda.Start(handler.WithRecover(handler.WithLogging(handler.WithCorrelationID(h.Handle))))
}
