package paymentuc

import (
	"context"
	"fmt"
	"time"

	domainerrors "github.com/davidlramirez95/towcommand/internal/domain/errors"
	"github.com/davidlramirez95/towcommand/internal/domain/event"
	"github.com/davidlramirez95/towcommand/internal/domain/payment"
)

// RefundPaymentInput holds the data needed to refund a payment.
type RefundPaymentInput struct {
	PaymentID string
	Reason    string
}

// RefundPaymentUseCase orchestrates refunding a captured payment via the
// payment gateway.
type RefundPaymentUseCase struct {
	payments PaymentFinder
	updater  PaymentStatusUpdater
	gateway  PaymentGateway
	events   EventPublisher
	now      func() time.Time
}

// NewRefundPaymentUseCase constructs a RefundPaymentUseCase with its dependencies.
func NewRefundPaymentUseCase(
	payments PaymentFinder,
	updater PaymentStatusUpdater,
	gateway PaymentGateway,
	events EventPublisher,
) *RefundPaymentUseCase {
	return &RefundPaymentUseCase{
		payments: payments,
		updater:  updater,
		gateway:  gateway,
		events:   events,
		now:      func() time.Time { return time.Now().UTC() },
	}
}

// Execute refunds a captured payment. It calls the gateway to process the
// refund, updates the payment status, and publishes a PaymentRefunded event.
func (uc *RefundPaymentUseCase) Execute(ctx context.Context, input *RefundPaymentInput) (*payment.Payment, error) {
	p, err := uc.payments.FindByID(ctx, input.PaymentID)
	if err != nil {
		return nil, fmt.Errorf("finding payment %s: %w", input.PaymentID, err)
	}
	if p == nil {
		return nil, domainerrors.NewNotFoundError("payment", input.PaymentID)
	}

	if p.Status != payment.PaymentStatusCaptured {
		return nil, domainerrors.NewConflictError(
			fmt.Sprintf("payment %s is in status %s, expected captured", input.PaymentID, p.Status),
		)
	}

	_, err = uc.gateway.Refund(ctx, p.GatewayRef, p.Amount)
	if err != nil {
		return nil, domainerrors.NewPaymentFailedError(fmt.Sprintf("refund gateway error: %v", err))
	}

	if err := uc.updater.UpdateStatus(ctx, input.PaymentID, payment.PaymentStatusRefunded); err != nil {
		return nil, fmt.Errorf("updating payment %s status to refunded: %w", input.PaymentID, err)
	}

	now := uc.now()
	p.Status = payment.PaymentStatusRefunded
	p.RefundedAt = &now
	p.RefundReason = input.Reason

	_ = uc.events.Publish(ctx, event.SourcePayment, event.PaymentRefunded, map[string]any{
		"paymentId": input.PaymentID,
		"bookingId": p.BookingID,
		"amount":    p.Amount,
		"reason":    input.Reason,
	}, &Actor{UserID: p.UserID, UserType: "system"})

	return p, nil
}
