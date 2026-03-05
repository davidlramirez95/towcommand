// Package main is the composition root for the trigger-analytics EventBridge subscriber Lambda.
// It processes domain events and updates DynamoDB atomic counters and heatmap cells
// for the TowCommand analytics dashboard.
package main

import (
	"context"
	"log/slog"
	"runtime/debug"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"

	"github.com/davidlramirez95/towcommand/internal/adapter/handler"
	"github.com/davidlramirez95/towcommand/internal/adapter/repository"
	"github.com/davidlramirez95/towcommand/internal/platform/awsclient"
	"github.com/davidlramirez95/towcommand/internal/platform/config"
	"github.com/davidlramirez95/towcommand/internal/platform/logger"
	"github.com/davidlramirez95/towcommand/internal/usecase/analytics"
)

func main() {
	cfg := config.Load()
	log := logger.New(cfg.Stage, cfg.FunctionName, cfg.FunctionVersion, cfg.LogLevel)
	slog.SetDefault(log)

	ddb := awsclient.DynamoDBClient(cfg)
	analyticsRepo := repository.NewDynamoAnalyticsRepository(ddb, cfg.DynamoDBTable)
	recorder := analytics.NewEventRecorder(analyticsRepo)

	h := func(ctx context.Context, evt events.CloudWatchEvent) error {
		defer func() {
			if r := recover(); r != nil {
				slog.ErrorContext(ctx, "panic recovered in trigger-analytics",
					"panic", r,
					"stack", string(debug.Stack()),
				)
			}
		}()

		correlationID := handler.ExtractCorrelationID(&evt)
		ctx = logger.SetCorrelationID(ctx, correlationID)

		slog.InfoContext(ctx, "trigger-analytics invoked",
			"detail_type", evt.DetailType,
			"event_id", evt.ID,
		)

		if err := recorder.Record(ctx, evt.DetailType, evt.Detail, evt.Time); err != nil {
			slog.ErrorContext(ctx, "analytics recording error",
				"error", err,
				"detail_type", evt.DetailType,
			)
		}

		// Analytics are best-effort: always return success to EventBridge.
		return nil
	}

	lambda.Start(h)
}
