package handler

import (
	"context"
	"net/http"

	"github.com/aws/aws-lambda-go/events"

	domainerrors "github.com/davidlramirez95/towcommand/internal/domain/errors"
	paymentuc "github.com/davidlramirez95/towcommand/internal/usecase/payment"
)

// ProviderEarningsHandler handles GET /providers/{id}/earnings requests.
type ProviderEarningsHandler struct {
	uc *paymentuc.GetProviderEarningsUseCase
}

// NewProviderEarningsHandler constructs a ProviderEarningsHandler.
func NewProviderEarningsHandler(uc *paymentuc.GetProviderEarningsUseCase) *ProviderEarningsHandler {
	return &ProviderEarningsHandler{uc: uc}
}

// Handle processes a provider-earnings API Gateway event.
// The caller must be authenticated and may only request their own earnings
// unless they have the "admin" user type.
func (h *ProviderEarningsHandler) Handle(ctx context.Context, event *events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	userID := ExtractUserID(event)
	if userID == "" {
		return ErrorResponse(domainerrors.NewUnauthorizedError()), nil
	}

	providerID := ParsePathParam(event, "id")
	if providerID == "" {
		return ErrorResponse(domainerrors.NewValidationError("provider ID is required")), nil
	}

	// Providers may only view their own earnings unless they are an admin.
	userType := ExtractUserType(event)
	if userID != providerID && userType != "admin" {
		return ErrorResponse(domainerrors.NewForbiddenError("you may only view your own earnings")), nil
	}

	result, err := h.uc.Execute(ctx, providerID)
	if err != nil {
		return ErrorResponse(err), nil
	}

	return SuccessResponse(http.StatusOK, result), nil
}
