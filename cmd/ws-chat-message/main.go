// Package main is the composition root for the WebSocket chatMessage Lambda.
package main

import (
	"context"
	"log/slog"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"

	"github.com/davidlramirez95/towcommand/internal/adapter/cache"
	"github.com/davidlramirez95/towcommand/internal/adapter/handler"
	"github.com/davidlramirez95/towcommand/internal/adapter/repository"
	"github.com/davidlramirez95/towcommand/internal/platform/awsclient"
	"github.com/davidlramirez95/towcommand/internal/platform/config"
	"github.com/davidlramirez95/towcommand/internal/platform/logger"
	ws "github.com/davidlramirez95/towcommand/internal/usecase/websocket"
)

// chatMessageBody is the JSON body expected from the WebSocket message.
type chatMessageBody struct {
	BookingID   string `json:"bookingId" validate:"required"`
	Message     string `json:"message" validate:"required,max=1000"`
	RecipientID string `json:"recipientId" validate:"required"`
}

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

	ddb := awsclient.DynamoDBClient(cfg)
	apigwClient := awsclient.APIGatewayManagementClient(cfg, "")
	redisClient := cache.NewRedisClient(cache.Options{
		Host:     cfg.RedisHost,
		Port:     cfg.RedisPort,
		Password: cfg.RedisPassword,
	})

	chatRepo := repository.NewChatRepository(ddb, cfg.DynamoDBTable)
	sessions := cache.NewRedisSessionCache(redisClient)
	poster := &wsPoster{client: apigwClient}
	uc := ws.NewChatMessageUseCase(chatRepo, sessions, poster, log)

	h := func(ctx context.Context, event *events.APIGatewayWebsocketProxyRequest) (events.APIGatewayProxyResponse, error) {
		body, err := handler.ParseWSBody[chatMessageBody](event)
		if err != nil {
			slog.ErrorContext(ctx, "invalid chat message payload", slog.String("error", err.Error()))
			return events.APIGatewayProxyResponse{StatusCode: 400}, nil
		}

		// Extract sender from the authorizer context or query params.
		var senderID string
		if auth, ok := event.RequestContext.Authorizer.(map[string]interface{}); ok {
			senderID, _ = auth["principalId"].(string)
		}
		if senderID == "" {
			senderID = event.QueryStringParameters["userId"]
		}

		input := ws.ChatMessageInput{
			BookingID: body.BookingID,
			Message:   body.Message,
			SenderID:  senderID,
		}

		if err := uc.Execute(ctx, input, body.RecipientID); err != nil {
			slog.ErrorContext(ctx, "chat message failed", slog.String("error", err.Error()))
			return events.APIGatewayProxyResponse{StatusCode: 500}, nil
		}

		return events.APIGatewayProxyResponse{StatusCode: 200}, nil
	}

	lambda.Start(h)
}
