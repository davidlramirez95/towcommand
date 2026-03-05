package paymentuc

import (
	"context"
	"crypto/rand"
	"fmt"
	"time"

	"github.com/davidlramirez95/towcommand/internal/domain/booking"
	domainerrors "github.com/davidlramirez95/towcommand/internal/domain/errors"
	"github.com/davidlramirez95/towcommand/internal/domain/event"
	"github.com/davidlramirez95/towcommand/internal/domain/payment"
)

// InitiatePaymentInput holds the data needed to initiate a payment.
type InitiatePaymentInput struct {
	BookingID string
	Method    payment.PaymentMethod
	UserID    string
}

// InitiatePaymentUseCase orchestrates the creation of a new payment for a booking.
type InitiatePaymentUseCase struct {
	bookings BookingFinder
	payments PaymentByBookingLister
	saver    PaymentSaver
	gateway  PaymentGateway
	events   EventPublisher
	idGen    func() string
	now      func() time.Time
}

// NewInitiatePaymentUseCase constructs an InitiatePaymentUseCase with its dependencies.
func NewInitiatePaymentUseCase(
	bookings BookingFinder,
	payments PaymentByBookingLister,
	saver PaymentSaver,
	gateway PaymentGateway,
	events EventPublisher,
) *InitiatePaymentUseCase {
	return &InitiatePaymentUseCase{
		bookings: bookings,
		payments: payments,
		saver:    saver,
		gateway:  gateway,
		events:   events,
		idGen:    generatePaymentID,
		now:      func() time.Time { return time.Now().UTC() },
	}
}

// Execute creates a payment for the given booking.
//
// For cash payments the payment is immediately marked as captured. For digital
// methods (gcash, maya, card) the gateway is called to create a pending charge.
func (uc *InitiatePaymentUseCase) Execute(ctx context.Context, input *InitiatePaymentInput) (*payment.Payment, error) {
	b, err := uc.bookings.FindByID(ctx, input.BookingID)
	if err != nil {
		return nil, fmt.Errorf("finding booking %s: %w", input.BookingID, err)
	}
	if b == nil {
		return nil, domainerrors.NewNotFoundError("booking", input.BookingID)
	}

	if b.Status != booking.BookingStatusCompleted {
		return nil, domainerrors.NewValidationError(
			fmt.Sprintf("booking %s is in status %s, expected COMPLETED", input.BookingID, b.Status),
		)
	}

	// Check for existing payments to prevent duplicates.
	existing, err := uc.payments.FindByBooking(ctx, input.BookingID)
	if err != nil {
		return nil, fmt.Errorf("listing payments for booking %s: %w", input.BookingID, err)
	}
	for _, ep := range existing {
		if ep.Status == payment.PaymentStatusCaptured || ep.Status == payment.PaymentStatusPending {
			return nil, domainerrors.NewConflictError(
				fmt.Sprintf("booking %s already has a %s payment", input.BookingID, ep.Status),
			)
		}
	}

	now := uc.now()
	paymentID := uc.idGen()

	p := &payment.Payment{
		PaymentID: paymentID,
		BookingID: input.BookingID,
		UserID:    input.UserID,
		Amount:    b.Price.Total,
		Currency:  b.Price.Currency,
		Method:    input.Method,
		CreatedAt: now,
		UpdatedAt: now,
	}

	if input.Method == payment.PaymentMethodCash {
		p.Status = payment.PaymentStatusCaptured
		p.CapturedAt = &now
	} else {
		chargeResult, chargeErr := uc.gateway.Charge(ctx, paymentID, b.Price.Total, b.Price.Currency, string(input.Method))
		if chargeErr != nil {
			return nil, domainerrors.NewPaymentFailedError(chargeErr.Error())
		}
		p.Status = payment.PaymentStatusPending
		p.GatewayRef = chargeResult.GatewayRef
	}

	if err := uc.saver.Save(ctx, p); err != nil {
		return nil, fmt.Errorf("saving payment %s: %w", paymentID, err)
	}

	detailType := event.PaymentInitiated
	if input.Method == payment.PaymentMethodCash {
		detailType = event.PaymentCaptured
	}

	_ = uc.events.Publish(ctx, event.SourcePayment, detailType, map[string]any{
		"paymentId": paymentID,
		"bookingId": input.BookingID,
		"amount":    p.Amount,
		"currency":  p.Currency,
		"method":    string(input.Method),
		"status":    string(p.Status),
	}, &Actor{UserID: input.UserID, UserType: "customer"})

	return p, nil
}

// generatePaymentID produces a payment ID in the format PAY-<year>-<hex>.
func generatePaymentID() string {
	b := make([]byte, 10)
	_, _ = rand.Read(b)
	return fmt.Sprintf("PAY-%d-%X", time.Now().Year(), b)
}
