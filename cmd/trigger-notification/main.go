// Package main is the composition root for the trigger-notification EventBridge subscriber Lambda.
// It routes domain events to the appropriate notification channels (SMS, email)
// using Filipino/English message templates.
package main

import (
	"context"
	"log/slog"
	"os"
	"runtime/debug"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"

	"github.com/davidlramirez95/towcommand/internal/adapter/gateway"
	"github.com/davidlramirez95/towcommand/internal/adapter/handler"
	"github.com/davidlramirez95/towcommand/internal/adapter/repository"
	"github.com/davidlramirez95/towcommand/internal/platform/awsclient"
	"github.com/davidlramirez95/towcommand/internal/platform/config"
	"github.com/davidlramirez95/towcommand/internal/platform/logger"
	"github.com/davidlramirez95/towcommand/internal/usecase/notification"
)

func main() {
	cfg := config.Load()
	log := logger.New(cfg.Stage, cfg.FunctionName, cfg.FunctionVersion, cfg.LogLevel)
	slog.SetDefault(log)

	ddb := awsclient.DynamoDBClient(cfg)
	snsClient := awsclient.SNSClient(cfg)
	sesClient := awsclient.SESClient(cfg)

	smsSender := gateway.NewSNSNotificationSender(snsClient)
	emailSender := gateway.NewSESNotificationSender(sesClient)
	userRepo := repository.NewUserRepository(ddb, cfg.DynamoDBTable)
	bookingRepo := repository.NewBookingRepository(ddb, cfg.DynamoDBTable)

	opsPhone := os.Getenv("OPS_PHONE_NUMBER")
	if opsPhone == "" {
		opsPhone = "+639170000000" // default ops phone for dev
	}
	safetyEmail := os.Getenv("SAFETY_EMAIL")
	if safetyEmail == "" {
		safetyEmail = "safety@towcommand.ph"
	}

	router := notification.NewNotificationRouter(
		smsSender, emailSender, userRepo, bookingRepo,
		opsPhone, safetyEmail,
	)

	h := func(ctx context.Context, evt events.CloudWatchEvent) error {
		defer func() {
			if r := recover(); r != nil {
				slog.ErrorContext(ctx, "panic recovered in trigger-notification",
					"panic", r,
					"stack", string(debug.Stack()),
				)
			}
		}()

		correlationID := handler.ExtractCorrelationID(&evt)
		ctx = logger.SetCorrelationID(ctx, correlationID)

		slog.InfoContext(ctx, "trigger-notification invoked",
			"detail_type", evt.DetailType,
			"event_id", evt.ID,
		)

		if err := router.Route(ctx, evt.DetailType, evt.Detail); err != nil {
			slog.ErrorContext(ctx, "notification routing error",
				"error", err,
				"detail_type", evt.DetailType,
			)
		}

		// Notifications are best-effort: always return success to EventBridge.
		return nil
	}

	lambda.Start(h)
}
