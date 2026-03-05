// Package paymentuc implements payment use cases following CLEAN architecture.
// Each use case declares only the port interfaces it needs (ISP).
package paymentuc

import (
	"context"

	"github.com/davidlramirez95/towcommand/internal/domain/booking"
	"github.com/davidlramirez95/towcommand/internal/domain/payment"
	"github.com/davidlramirez95/towcommand/internal/domain/provider"
	"github.com/davidlramirez95/towcommand/internal/usecase/port"
)

// PaymentSaver persists a new payment.
type PaymentSaver interface {
	Save(ctx context.Context, p *payment.Payment) error
}

// PaymentFinder retrieves a payment by its ID.
type PaymentFinder interface {
	FindByID(ctx context.Context, paymentID string) (*payment.Payment, error)
}

// PaymentByBookingLister lists payments for a given booking.
type PaymentByBookingLister interface {
	FindByBooking(ctx context.Context, bookingID string) ([]payment.Payment, error)
}

// PaymentStatusUpdater changes a payment's status.
type PaymentStatusUpdater interface {
	UpdateStatus(ctx context.Context, paymentID string, status payment.PaymentStatus) error
}

// BookingFinder retrieves a booking by its ID.
type BookingFinder interface {
	FindByID(ctx context.Context, bookingID string) (*booking.Booking, error)
}

// ProviderFinder retrieves a provider by their ID.
type ProviderFinder interface {
	FindByID(ctx context.Context, providerID string) (*provider.Provider, error)
}

// EventPublisher publishes domain events to an event bus.
type EventPublisher interface {
	Publish(ctx context.Context, source, detailType string, detail any, actor *port.Actor) error
}

// PaymentGateway defines the interface for external payment providers.
type PaymentGateway interface {
	Charge(ctx context.Context, paymentID string, amountCentavos int64, currency, method string) (*port.ChargeResult, error)
	Refund(ctx context.Context, gatewayRef string, amountCentavos int64) (*port.RefundResult, error)
	VerifyWebhookSignature(payload []byte, signature string) error
}

// Actor is a type alias for port.Actor to avoid import stuttering in use case code.
type Actor = port.Actor
