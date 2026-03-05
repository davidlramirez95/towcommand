package paymentuc

import (
	"context"
	"fmt"
	"time"

	domainerrors "github.com/davidlramirez95/towcommand/internal/domain/errors"
	"github.com/davidlramirez95/towcommand/internal/domain/event"
	"github.com/davidlramirez95/towcommand/internal/domain/payment"
)

// CaptureResult holds the outcome of a successful payment capture including
// the commission breakdown.
type CaptureResult struct {
	Payment    *payment.Payment
	Commission int64
	NetAmount  int64
}

// CapturePaymentUseCase orchestrates capturing a pending payment after gateway
// confirmation.
type CapturePaymentUseCase struct {
	payments PaymentFinder
	bookings BookingFinder
	provider ProviderFinder
	updater  PaymentStatusUpdater
	events   EventPublisher
	now      func() time.Time
}

// NewCapturePaymentUseCase constructs a CapturePaymentUseCase with its dependencies.
func NewCapturePaymentUseCase(
	payments PaymentFinder,
	bookings BookingFinder,
	provider ProviderFinder,
	updater PaymentStatusUpdater,
	events EventPublisher,
) *CapturePaymentUseCase {
	return &CapturePaymentUseCase{
		payments: payments,
		bookings: bookings,
		provider: provider,
		updater:  updater,
		events:   events,
		now:      func() time.Time { return time.Now().UTC() },
	}
}

// Execute captures a pending payment, calculates the commission split, and
// publishes a PaymentCaptured event.
func (uc *CapturePaymentUseCase) Execute(ctx context.Context, paymentID string) (*CaptureResult, error) {
	p, err := uc.payments.FindByID(ctx, paymentID)
	if err != nil {
		return nil, fmt.Errorf("finding payment %s: %w", paymentID, err)
	}
	if p == nil {
		return nil, domainerrors.NewNotFoundError("payment", paymentID)
	}

	if p.Status != payment.PaymentStatusPending {
		return nil, domainerrors.NewConflictError(
			fmt.Sprintf("payment %s is in status %s, expected pending", paymentID, p.Status),
		)
	}

	b, err := uc.bookings.FindByID(ctx, p.BookingID)
	if err != nil {
		return nil, fmt.Errorf("finding booking %s for payment %s: %w", p.BookingID, paymentID, err)
	}
	if b == nil {
		return nil, domainerrors.NewNotFoundError("booking", p.BookingID)
	}

	prov, err := uc.provider.FindByID(ctx, b.ProviderID)
	if err != nil {
		return nil, fmt.Errorf("finding provider %s for payment %s: %w", b.ProviderID, paymentID, err)
	}
	if prov == nil {
		return nil, domainerrors.NewNotFoundError("provider", b.ProviderID)
	}

	commission, commissionRate := CalculateCommission(p.Amount, prov.TrustTier)
	netAmount := p.Amount - commission

	if err := uc.updater.UpdateStatus(ctx, paymentID, payment.PaymentStatusCaptured); err != nil {
		return nil, fmt.Errorf("updating payment %s status to captured: %w", paymentID, err)
	}

	now := uc.now()
	p.Status = payment.PaymentStatusCaptured
	p.CapturedAt = &now

	_ = uc.events.Publish(ctx, event.SourcePayment, event.PaymentCaptured, map[string]any{
		"paymentId":      paymentID,
		"bookingId":      p.BookingID,
		"amount":         p.Amount,
		"commission":     commission,
		"commissionRate": commissionRate,
		"netAmount":      netAmount,
		"providerId":     b.ProviderID,
		"providerTier":   string(prov.TrustTier),
	}, &Actor{UserID: p.UserID, UserType: "system"})

	return &CaptureResult{
		Payment:    p,
		Commission: commission,
		NetAmount:  netAmount,
	}, nil
}
