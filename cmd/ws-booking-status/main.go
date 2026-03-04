// Package main is the composition root for the WebSocket bookingStatus Lambda.
package main

import (
	"context"
	"log/slog"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"

	"github.com/davidlramirez95/towcommand/internal/adapter/cache"
	"github.com/davidlramirez95/towcommand/internal/adapter/handler"
	"github.com/davidlramirez95/towcommand/internal/platform/awsclient"
	"github.com/davidlramirez95/towcommand/internal/platform/config"
	"github.com/davidlramirez95/towcommand/internal/platform/logger"
	ws "github.com/davidlramirez95/towcommand/internal/usecase/websocket"
)

// wsPoster adapts handler.SendToConnection to the websocket.ConnectionPoster interface.
type wsPoster struct {
	client handler.ConnectionPoster
}

func (p *wsPoster) PostToConnection(ctx context.Context, connectionID string, data any) error {
	return handler.SendToConnection(ctx, p.client, connectionID, data)
}

func main() {
	cfg := config.Load()
	log := logger.New(cfg.Stage, cfg.FunctionName, cfg.FunctionVersion, cfg.LogLevel)
	slog.SetDefault(log)

	apigwClient := awsclient.APIGatewayManagementClient(cfg, "")
	redisClient := cache.NewRedisClient(cache.Options{
		Host:     cfg.RedisHost,
		Port:     cfg.RedisPort,
		Password: cfg.RedisPassword,
	})

	sessions := cache.NewRedisSessionCache(redisClient)
	poster := &wsPoster{client: apigwClient}
	uc := ws.NewBookingStatusUseCase(sessions, poster, log)

	h := func(ctx context.Context, event *events.APIGatewayWebsocketProxyRequest) (events.APIGatewayProxyResponse, error) {
		input, err := handler.ParseWSBody[ws.BookingStatusInput](event)
		if err != nil {
			slog.ErrorContext(ctx, "invalid booking status payload", slog.String("error", err.Error()))
			return events.APIGatewayProxyResponse{StatusCode: 400}, nil
		}

		if err := uc.Execute(ctx, input); err != nil {
			slog.ErrorContext(ctx, "booking status push failed", slog.String("error", err.Error()))
			return events.APIGatewayProxyResponse{StatusCode: 500}, nil
		}

		return events.APIGatewayProxyResponse{StatusCode: 200}, nil
	}

	lambda.Start(h)
}
