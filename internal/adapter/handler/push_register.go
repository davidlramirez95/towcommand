package handler

import (
	"context"
	"net/http"
	"time"

	"github.com/aws/aws-lambda-go/events"

	domainerrors "github.com/davidlramirez95/towcommand/internal/domain/errors"
	"github.com/davidlramirez95/towcommand/internal/usecase/port"
)

// RegisterPushTokenRequest is the expected JSON body for POST /users/{id}/push-token.
type RegisterPushTokenRequest struct {
	Token    string `json:"token" validate:"required"`
	Platform string `json:"platform" validate:"required,oneof=FCM APNS"`
	DeviceID string `json:"deviceId" validate:"required"`
}

// RegisterPushTokenHandler handles POST /users/{id}/push-token requests.
type RegisterPushTokenHandler struct {
	tokens   port.PushTokenRegistrar
	endpoint port.PushEndpointCreator
}

// NewRegisterPushTokenHandler constructs a RegisterPushTokenHandler.
func NewRegisterPushTokenHandler(tokens port.PushTokenRegistrar, endpoint port.PushEndpointCreator) *RegisterPushTokenHandler {
	return &RegisterPushTokenHandler{
		tokens:   tokens,
		endpoint: endpoint,
	}
}

// Handle processes a push token registration API Gateway event.
func (h *RegisterPushTokenHandler) Handle(ctx context.Context, event *events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	userID := ExtractUserID(event)
	if userID == "" {
		return ErrorResponse(domainerrors.NewUnauthorizedError()), nil
	}

	pathUserID := ParsePathParam(event, "id")
	if pathUserID == "" {
		return ErrorResponse(domainerrors.NewValidationError("user ID is required")), nil
	}

	// The authenticated user must match the path user (or be an admin).
	userType := ExtractUserType(event)
	if pathUserID != userID && userType != "admin" {
		return ErrorResponse(domainerrors.NewForbiddenError("cannot register push token for another user")), nil
	}

	body, err := ParseBody[RegisterPushTokenRequest](event)
	if err != nil {
		return ErrorResponse(err), nil
	}

	platform := port.PushPlatform(body.Platform)

	// Create SNS platform endpoint for the device token.
	endpointArn, err := h.endpoint.CreateEndpoint(ctx, platform, body.Token)
	if err != nil {
		return ErrorResponse(err), nil
	}

	now := time.Now().UTC()
	token := &port.PushToken{
		UserID:      pathUserID,
		Token:       body.Token,
		Platform:    platform,
		DeviceID:    body.DeviceID,
		EndpointArn: endpointArn,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	if err := h.tokens.Register(ctx, token); err != nil {
		return ErrorResponse(domainerrors.NewInternalError("failed to save push token").WithCause(err)), nil
	}

	return SuccessResponse(http.StatusCreated, token), nil
}
