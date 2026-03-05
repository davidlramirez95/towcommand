package main

import (
	"log/slog"
	"os"

	"github.com/aws/aws-lambda-go/lambda"

	"github.com/davidlramirez95/towcommand/internal/adapter/gateway"
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
	eb := awsclient.EventBridgeClient(cfg)

	bookingRepo := repository.NewBookingRepository(ddb, cfg.DynamoDBTable)
	paymentRepo := repository.NewPaymentRepository(ddb, cfg.DynamoDBTable)
	events := gateway.NewEventBridgePublisher(eb, cfg.EventBusName, log)

	webhookSecret := os.Getenv("PAYMENT_WEBHOOK_SECRET")
	if webhookSecret == "" {
		webhookSecret = "dev-secret"
	}
	gw := gateway.NewMockPaymentGateway(webhookSecret)

	uc := paymentuc.NewInitiatePaymentUseCase(bookingRepo, paymentRepo, paymentRepo, gw, events)
	h := handler.NewInitiatePaymentHandler(uc)

	lambda.Start(handler.WithRecover(handler.WithLogging(handler.WithCorrelationID(h.Handle))))
}
