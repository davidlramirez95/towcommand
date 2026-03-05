package paymentuc

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	domainerrors "github.com/davidlramirez95/towcommand/internal/domain/errors"
	"github.com/davidlramirez95/towcommand/internal/domain/payment"
)

// webhookPayload is the expected JSON structure from the payment gateway webhook.
type webhookPayload struct {
	PaymentID string `json:"paymentId"`
	Event     string `json:"event"`
}

// ProcessWebhookInput holds the raw webhook data for verification and processing.
type ProcessWebhookInput struct {
	Payload   []byte
	Signature string
}

// ProcessWebhookUseCase handles incoming payment gateway webhooks. It verifies
// the signature, parses the event, and routes to the appropriate status update.
type ProcessWebhookUseCase struct {
	gateway  PaymentGateway
	payments PaymentFinder
	updater  PaymentStatusUpdater
	now      func() time.Time
}

// NewProcessWebhookUseCase constructs a ProcessWebhookUseCase with its dependencies.
func NewProcessWebhookUseCase(
	gateway PaymentGateway,
	payments PaymentFinder,
	updater PaymentStatusUpdater,
) *ProcessWebhookUseCase {
	return &ProcessWebhookUseCase{
		gateway:  gateway,
		payments: payments,
		updater:  updater,
		now:      func() time.Time { return time.Now().UTC() },
	}
}

// Execute verifies the webhook signature, parses the payload, and updates the
// payment status accordingly. It is idempotent: if the payment is already in
// the target status, it returns the payment without error.
func (uc *ProcessWebhookUseCase) Execute(ctx context.Context, input *ProcessWebhookInput) (*payment.Payment, error) {
	if err := uc.gateway.VerifyWebhookSignature(input.Payload, input.Signature); err != nil {
		return nil, domainerrors.NewUnauthorizedError()
	}

	var wp webhookPayload
	if err := json.Unmarshal(input.Payload, &wp); err != nil {
		return nil, domainerrors.NewValidationError(fmt.Sprintf("invalid webhook payload: %v", err))
	}

	if wp.PaymentID == "" {
		return nil, domainerrors.NewValidationError("webhook payload missing paymentId")
	}

	p, err := uc.payments.FindByID(ctx, wp.PaymentID)
	if err != nil {
		return nil, fmt.Errorf("finding payment %s: %w", wp.PaymentID, err)
	}
	if p == nil {
		return nil, domainerrors.NewNotFoundError("payment", wp.PaymentID)
	}

	targetStatus, err := webhookEventToStatus(wp.Event)
	if err != nil {
		return nil, err
	}

	// Idempotency: if already in target status, return success.
	if p.Status == targetStatus {
		return p, nil
	}

	if err := uc.updater.UpdateStatus(ctx, wp.PaymentID, targetStatus); err != nil {
		return nil, fmt.Errorf("updating payment %s to %s: %w", wp.PaymentID, targetStatus, err)
	}

	now := uc.now()
	p.Status = targetStatus

	switch targetStatus {
	case payment.PaymentStatusCaptured:
		p.CapturedAt = &now
	case payment.PaymentStatusRefunded:
		p.RefundedAt = &now
	}

	return p, nil
}

// webhookEventToStatus maps a gateway webhook event string to the corresponding
// payment status.
func webhookEventToStatus(eventType string) (payment.PaymentStatus, error) {
	switch eventType {
	case "payment.captured":
		return payment.PaymentStatusCaptured, nil
	case "payment.refunded":
		return payment.PaymentStatusRefunded, nil
	case "payment.failed":
		return payment.PaymentStatusFailed, nil
	default:
		return "", domainerrors.NewValidationError(fmt.Sprintf("unknown webhook event type: %s", eventType))
	}
}
