package handler

import (
	"context"
	"net/http"

	"github.com/aws/aws-lambda-go/events"

	domainerrors "github.com/davidlramirez95/towcommand/internal/domain/errors"
	paymentuc "github.com/davidlramirez95/towcommand/internal/usecase/payment"
)

// RefundPaymentRequest is the expected JSON body for POST /payments/{id}/refund.
type RefundPaymentRequest struct {
	Reason string `json:"reason" validate:"required"`
}

// RefundPaymentHandler handles POST /payments/{id}/refund requests (admin).
type RefundPaymentHandler struct {
	uc *paymentuc.RefundPaymentUseCase
}

// NewRefundPaymentHandler constructs a RefundPaymentHandler.
func NewRefundPaymentHandler(uc *paymentuc.RefundPaymentUseCase) *RefundPaymentHandler {
	return &RefundPaymentHandler{uc: uc}
}

// Handle processes a refund-payment API Gateway event.
func (h *RefundPaymentHandler) Handle(ctx context.Context, event *events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	userID := ExtractUserID(event)
	if userID == "" {
		return ErrorResponse(domainerrors.NewUnauthorizedError()), nil
	}

	paymentID := ParsePathParam(event, "id")
	if paymentID == "" {
		return ErrorResponse(domainerrors.NewValidationError("payment ID is required")), nil
	}

	body, err := ParseBody[RefundPaymentRequest](event)
	if err != nil {
		return ErrorResponse(err), nil
	}

	result, err := h.uc.Execute(ctx, &paymentuc.RefundPaymentInput{
		PaymentID: paymentID,
		Reason:    body.Reason,
	})
	if err != nil {
		return ErrorResponse(err), nil
	}

	return SuccessResponse(http.StatusOK, result), nil
}
