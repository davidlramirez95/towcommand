// Package main is the composition root for the WebSocket $connect Lambda.
package main

import (
	"context"
	"log/slog"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"

	"github.com/davidlramirez95/towcommand/internal/adapter/cache"
	"github.com/davidlramirez95/towcommand/internal/adapter/handler"
	"github.com/davidlramirez95/towcommand/internal/platform/config"
	"github.com/davidlramirez95/towcommand/internal/platform/logger"
	ws "github.com/davidlramirez95/towcommand/internal/usecase/websocket"
)

func main() {
	cfg := config.Load()
	log := logger.New(cfg.Stage, cfg.FunctionName, cfg.FunctionVersion, cfg.LogLevel)
	slog.SetDefault(log)

	redisClient := cache.NewRedisClient(cache.Options{
		Host:     cfg.RedisHost,
		Port:     cfg.RedisPort,
		Password: cfg.RedisPassword,
	})
	sessions := cache.NewRedisSessionCache(redisClient)
	uc := ws.NewConnectUseCase(sessions, log)

	h := func(ctx context.Context, event *events.APIGatewayWebsocketProxyRequest) (events.APIGatewayProxyResponse, error) {
		connectionID := handler.ExtractConnectionID(event)
		userID := event.QueryStringParameters["userId"]
		if userID == "" {
			return events.APIGatewayProxyResponse{StatusCode: 400, Body: "missing userId query parameter"}, nil
		}

		if err := uc.Execute(ctx, userID, connectionID); err != nil {
			slog.ErrorContext(ctx, "connect failed", slog.String("error", err.Error()))
			return events.APIGatewayProxyResponse{StatusCode: 500}, nil
		}

		return events.APIGatewayProxyResponse{StatusCode: 200}, nil
	}

	lambda.Start(h)
}
