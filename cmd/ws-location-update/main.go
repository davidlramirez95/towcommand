// Package main is the composition root for the WebSocket locationUpdate Lambda.
package main

import (
	"context"
	"log/slog"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"

	"github.com/davidlramirez95/towcommand/internal/adapter/cache"
	"github.com/davidlramirez95/towcommand/internal/adapter/gateway"
	"github.com/davidlramirez95/towcommand/internal/adapter/handler"
	"github.com/davidlramirez95/towcommand/internal/platform/awsclient"
	"github.com/davidlramirez95/towcommand/internal/platform/config"
	"github.com/davidlramirez95/towcommand/internal/platform/logger"
	ws "github.com/davidlramirez95/towcommand/internal/usecase/websocket"
)

func main() {
	cfg := config.Load()
	log := logger.New(cfg.Stage, cfg.FunctionName, cfg.FunctionVersion, cfg.LogLevel)
	slog.SetDefault(log)

	ebClient := awsclient.EventBridgeClient(cfg)
	redisClient := cache.NewRedisClient(cache.Options{
		Host:     cfg.RedisHost,
		Port:     cfg.RedisPort,
		Password: cfg.RedisPassword,
	})

	geo := cache.NewRedisGeoCache(redisClient)
	pub := gateway.NewEventBridgePublisher(ebClient, cfg.EventBusName, log)
	uc := ws.NewLocationUpdateUseCase(geo, pub, log)

	h := func(ctx context.Context, event *events.APIGatewayWebsocketProxyRequest) (events.APIGatewayProxyResponse, error) {
		input, err := handler.ParseWSBody[ws.LocationUpdateInput](event)
		if err != nil {
			slog.ErrorContext(ctx, "invalid location update payload", slog.String("error", err.Error()))
			return events.APIGatewayProxyResponse{StatusCode: 400}, nil
		}

		if err := uc.Execute(ctx, input); err != nil {
			slog.ErrorContext(ctx, "location update failed", slog.String("error", err.Error()))
			return events.APIGatewayProxyResponse{StatusCode: 500}, nil
		}

		return events.APIGatewayProxyResponse{StatusCode: 200}, nil
	}

	lambda.Start(h)
}
