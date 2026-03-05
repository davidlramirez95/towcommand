package handler

import (
	"context"
	"net/http"

	"github.com/aws/aws-lambda-go/events"

	domainerrors "github.com/davidlramirez95/towcommand/internal/domain/errors"
	paymentuc "github.com/davidlramirez95/towcommand/internal/usecase/payment"
)

// PaymentWebhookHandler handles incoming payment gateway webhook callbacks.
// It does NOT require Cognito authentication since webhooks originate from
// external payment providers, not end users.
type PaymentWebhookHandler struct {
	uc *paymentuc.ProcessWebhookUseCase
}

// NewPaymentWebhookHandler constructs a PaymentWebhookHandler.
func NewPaymentWebhookHandler(uc *paymentuc.ProcessWebhookUseCase) *PaymentWebhookHandler {
	return &PaymentWebhookHandler{uc: uc}
}

// Handle processes a payment webhook API Gateway event.
func (h *PaymentWebhookHandler) Handle(ctx context.Context, event *events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	signature := event.Headers["x-webhook-signature"]
	if signature == "" {
		signature = event.Headers["X-Webhook-Signature"]
	}
	if signature == "" {
		return ErrorResponse(domainerrors.NewUnauthorizedError()), nil
	}

	result, err := h.uc.Execute(ctx, &paymentuc.ProcessWebhookInput{
		Payload:   []byte(event.Body),
		Signature: signature,
	})
	if err != nil {
		return ErrorResponse(err), nil
	}

	return SuccessResponse(http.StatusOK, result), nil
}
