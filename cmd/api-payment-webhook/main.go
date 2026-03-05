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

	paymentRepo := repository.NewPaymentRepository(ddb, cfg.DynamoDBTable)

	webhookSecret := os.Getenv("PAYMENT_WEBHOOK_SECRET")
	if webhookSecret == "" {
		webhookSecret = "dev-secret"
	}
	gw := gateway.NewMockPaymentGateway(webhookSecret)

	uc := paymentuc.NewProcessWebhookUseCase(gw, paymentRepo, paymentRepo)
	h := handler.NewPaymentWebhookHandler(uc)

	// No Cognito auth middleware for webhooks - only recover, logging, and correlation ID.
	lambda.Start(handler.WithRecover(handler.WithLogging(handler.WithCorrelationID(h.Handle))))
}
