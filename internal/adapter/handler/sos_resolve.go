package handler

import (
	"context"
	"net/http"

	"github.com/aws/aws-lambda-go/events"

	domainerrors "github.com/davidlramirez95/towcommand/internal/domain/errors"
	safetyuc "github.com/davidlramirez95/towcommand/internal/usecase/safety"
)

// ResolveSOSHandler handles POST /sos/{id}/resolve requests.
type ResolveSOSHandler struct {
	uc *safetyuc.ResolveSOSUseCase
}

// NewResolveSOSHandler constructs a ResolveSOSHandler.
func NewResolveSOSHandler(uc *safetyuc.ResolveSOSUseCase) *ResolveSOSHandler {
	return &ResolveSOSHandler{uc: uc}
}

// Handle processes a resolve-SOS API Gateway event.
func (h *ResolveSOSHandler) Handle(ctx context.Context, event *events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	userID := ExtractUserID(event)
	if userID == "" {
		return ErrorResponse(domainerrors.NewUnauthorizedError()), nil
	}

	alertID := ParsePathParam(event, "id")
	if alertID == "" {
		return ErrorResponse(domainerrors.NewValidationError("alert ID is required")), nil
	}

	result, err := h.uc.Execute(ctx, &safetyuc.ResolveSOSInput{
		AlertID:    alertID,
		ResolvedBy: userID,
	})
	if err != nil {
		return ErrorResponse(err), nil
	}

	return SuccessResponse(http.StatusOK, result), nil
}
