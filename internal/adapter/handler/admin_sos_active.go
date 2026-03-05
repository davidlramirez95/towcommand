package handler

import (
	"context"
	"net/http"
	"strconv"

	"github.com/aws/aws-lambda-go/events"

	domainerrors "github.com/davidlramirez95/towcommand/internal/domain/errors"
	safetyuc "github.com/davidlramirez95/towcommand/internal/usecase/safety"
)

// AdminActiveSOSHandler handles GET /admin/sos/active requests.
type AdminActiveSOSHandler struct {
	uc *safetyuc.ListActiveSOSUseCase
}

// NewAdminActiveSOSHandler constructs an AdminActiveSOSHandler.
func NewAdminActiveSOSHandler(uc *safetyuc.ListActiveSOSUseCase) *AdminActiveSOSHandler {
	return &AdminActiveSOSHandler{uc: uc}
}

// Handle processes an admin list-active-SOS API Gateway event.
func (h *AdminActiveSOSHandler) Handle(ctx context.Context, event *events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	userID := ExtractUserID(event)
	if userID == "" {
		return ErrorResponse(domainerrors.NewUnauthorizedError()), nil
	}

	var limit int32 = 50
	if limitStr := ParseQueryParam(event, "limit"); limitStr != "" {
		n, err := strconv.ParseInt(limitStr, 10, 32)
		if err != nil || n <= 0 {
			return ErrorResponse(domainerrors.NewValidationError("limit must be a positive integer")), nil
		}
		limit = int32(n)
	}

	alerts, err := h.uc.Execute(ctx, limit)
	if err != nil {
		return ErrorResponse(err), nil
	}

	return SuccessResponse(http.StatusOK, alerts), nil
}
