package main

import (
	"log/slog"

	"github.com/aws/aws-lambda-go/lambda"

	"github.com/davidlramirez95/towcommand/internal/adapter/handler"
	"github.com/davidlramirez95/towcommand/internal/adapter/repository"
	"github.com/davidlramirez95/towcommand/internal/platform/awsclient"
	"github.com/davidlramirez95/towcommand/internal/platform/config"
	"github.com/davidlramirez95/towcommand/internal/platform/logger"
	paymentuc "github.com/davidlramirez95/towcommand/internal/usecase/payment"
)

func main() {
	cfg := config.Load()
	log := logger.New(cfg.Stage, cfg.FunctionName, cfg.FunctionVersion, cfg.LogLevel)
	slog.SetDefault(log)

	ddb := awsclient.DynamoDBClient(cfg)

	providerRepo := repository.NewProviderRepository(ddb, cfg.DynamoDBTable)
	bookingRepo := repository.NewBookingRepository(ddb, cfg.DynamoDBTable)
	paymentRepo := repository.NewPaymentRepository(ddb, cfg.DynamoDBTable)

	uc := paymentuc.NewGetProviderEarningsUseCase(providerRepo, bookingRepo, paymentRepo)
	h := handler.NewProviderEarningsHandler(uc)

	lambda.Start(handler.WithRecover(handler.WithLogging(handler.WithCorrelationID(h.Handle))))
}
