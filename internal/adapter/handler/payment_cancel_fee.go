package handler

import (
	"context"
	"net/http"

	"github.com/aws/aws-lambda-go/events"

	domainerrors "github.com/davidlramirez95/towcommand/internal/domain/errors"
	"github.com/davidlramirez95/towcommand/internal/domain/payment"
	paymentuc "github.com/davidlramirez95/towcommand/internal/usecase/payment"
)

// CancelFeeHandler handles POST /bookings/{id}/cancel-fee requests.
// It initiates a cash payment for the cancellation fee of a cancelled booking.
// Since cancel fees are always collected in cash, no payment method is needed
// from the request body.
type CancelFeeHandler struct {
	uc *paymentuc.InitiatePaymentUseCase
}

// NewCancelFeeHandler constructs a CancelFeeHandler.
func NewCancelFeeHandler(uc *paymentuc.InitiatePaymentUseCase) *CancelFeeHandler {
	return &CancelFeeHandler{uc: uc}
}

// Handle processes a cancel-fee API Gateway event.
func (h *CancelFeeHandler) Handle(ctx context.Context, event *events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	userID := ExtractUserID(event)
	if userID == "" {
		return ErrorResponse(domainerrors.NewUnauthorizedError()), nil
	}

	bookingID := ParsePathParam(event, "id")
	if bookingID == "" {
		return ErrorResponse(domainerrors.NewValidationError("booking ID is required")), nil
	}

	result, err := h.uc.Execute(ctx, &paymentuc.InitiatePaymentInput{
		BookingID: bookingID,
		Method:    payment.PaymentMethodCash,
		UserID:    userID,
	})
	if err != nil {
		return ErrorResponse(err), nil
	}

	return SuccessResponse(http.StatusCreated, result), nil
}
