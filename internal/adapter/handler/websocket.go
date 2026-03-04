package handler

import (
	"context"
	"encoding/json"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/apigatewaymanagementapi"

	domainerrors "github.com/davidlramirez95/towcommand/internal/domain/errors"
)

// ConnectionPoster abstracts the PostToConnection API for testability (ISP).
type ConnectionPoster interface {
	PostToConnection(ctx context.Context, params *apigatewaymanagementapi.PostToConnectionInput, optFns ...func(*apigatewaymanagementapi.Options)) (*apigatewaymanagementapi.PostToConnectionOutput, error)
}

// ParseWSBody unmarshals and validates the JSON body from a WebSocket event.
func ParseWSBody[T any](event *events.APIGatewayWebsocketProxyRequest) (T, error) {
	var body T
	if err := json.Unmarshal([]byte(event.Body), &body); err != nil {
		return body, domainerrors.NewValidationError("invalid JSON body").WithCause(err)
	}
	if err := validate.Struct(body); err != nil {
		return body, domainerrors.NewValidationError(err.Error()).WithCause(err)
	}
	return body, nil
}

// ExtractConnectionID returns the WebSocket connection ID from the event context.
func ExtractConnectionID(event *events.APIGatewayWebsocketProxyRequest) string {
	return event.RequestContext.ConnectionID
}

// SendToConnection sends JSON-encoded data to a WebSocket connection
// via the API Gateway Management API.
func SendToConnection(ctx context.Context, client ConnectionPoster, connectionID string, data any) error {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return domainerrors.NewInternalError("failed to marshal WebSocket message").WithCause(err)
	}
	_, err = client.PostToConnection(ctx, &apigatewaymanagementapi.PostToConnectionInput{
		ConnectionId: aws.String(connectionID),
		Data:         jsonData,
	})
	if err != nil {
		return domainerrors.NewExternalServiceError("apigateway-management", err)
	}
	return nil
}
