package main

import (
	"log/slog"

	"github.com/aws/aws-lambda-go/lambda"

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

	repo := repository.NewBookingRepository(ddb, cfg.DynamoDBTable)
	uc := bookinguc.NewListBookingsUseCase(repo)
	h := handler.NewListBookingsHandler(uc)

	lambda.Start(handler.WithRecover(handler.WithLogging(handler.WithCorrelationID(h.Handle))))
}
