package port

import (
	"context"

	"github.com/davidlramirez95/towcommand/internal/domain/payment"
)

// PaymentSaver persists a new payment.
type PaymentSaver interface {
	Save(ctx context.Context, p *payment.Payment) error
}

// PaymentFinder retrieves a payment by its ID.
type PaymentFinder interface {
	FindByID(ctx context.Context, paymentID string) (*payment.Payment, error)
}

// PaymentByBookingLister lists payments for a given booking via GSI1.
type PaymentByBookingLister interface {
	FindByBooking(ctx context.Context, bookingID string) ([]payment.Payment, error)
}

// PaymentStatusUpdater changes a payment's status.
type PaymentStatusUpdater interface {
	UpdateStatus(ctx context.Context, paymentID string, status payment.PaymentStatus) error
}
