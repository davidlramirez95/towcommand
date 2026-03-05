package handler

import (
	"context"
	"net/http"

	"github.com/aws/aws-lambda-go/events"

	domainerrors "github.com/davidlramirez95/towcommand/internal/domain/errors"
	"github.com/davidlramirez95/towcommand/internal/domain/payment"
	paymentuc "github.com/davidlramirez95/towcommand/internal/usecase/payment"
)

// InitiatePaymentRequest is the expected JSON body for POST /bookings/{id}/payments.
type InitiatePaymentRequest struct {
	Method string `json:"method" validate:"required,oneof=gcash maya card cash corporate"`
}

// InitiatePaymentHandler handles POST /bookings/{id}/payments requests.
type InitiatePaymentHandler struct {
	uc *paymentuc.InitiatePaymentUseCase
}

// NewInitiatePaymentHandler constructs an InitiatePaymentHandler.
func NewInitiatePaymentHandler(uc *paymentuc.InitiatePaymentUseCase) *InitiatePaymentHandler {
	return &InitiatePaymentHandler{uc: uc}
}

// Handle processes an initiate-payment API Gateway event.
func (h *InitiatePaymentHandler) Handle(ctx context.Context, event *events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	userID := ExtractUserID(event)
	if userID == "" {
		return ErrorResponse(domainerrors.NewUnauthorizedError()), nil
	}

	bookingID := ParsePathParam(event, "id")
	if bookingID == "" {
		return ErrorResponse(domainerrors.NewValidationError("booking ID is required")), nil
	}

	body, err := ParseBody[InitiatePaymentRequest](event)
	if err != nil {
		return ErrorResponse(err), nil
	}

	result, err := h.uc.Execute(ctx, &paymentuc.InitiatePaymentInput{
		BookingID: bookingID,
		Method:    payment.PaymentMethod(body.Method),
		UserID:    userID,
	})
	if err != nil {
		return ErrorResponse(err), nil
	}

	return SuccessResponse(http.StatusCreated, result), nil
}
