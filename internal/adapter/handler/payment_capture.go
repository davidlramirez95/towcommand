package handler

import (
	"context"
	"net/http"

	"github.com/aws/aws-lambda-go/events"

	domainerrors "github.com/davidlramirez95/towcommand/internal/domain/errors"
	paymentuc "github.com/davidlramirez95/towcommand/internal/usecase/payment"
)

// CapturePaymentHandler handles POST /payments/{id}/capture requests (admin).
type CapturePaymentHandler struct {
	uc *paymentuc.CapturePaymentUseCase
}

// NewCapturePaymentHandler constructs a CapturePaymentHandler.
func NewCapturePaymentHandler(uc *paymentuc.CapturePaymentUseCase) *CapturePaymentHandler {
	return &CapturePaymentHandler{uc: uc}
}

// Handle processes a capture-payment API Gateway event.
func (h *CapturePaymentHandler) Handle(ctx context.Context, event *events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	userID := ExtractUserID(event)
	if userID == "" {
		return ErrorResponse(domainerrors.NewUnauthorizedError()), nil
	}

	paymentID := ParsePathParam(event, "id")
	if paymentID == "" {
		return ErrorResponse(domainerrors.NewValidationError("payment ID is required")), nil
	}

	result, err := h.uc.Execute(ctx, paymentID)
	if err != nil {
		return ErrorResponse(err), nil
	}

	return SuccessResponse(http.StatusOK, result), nil
}
