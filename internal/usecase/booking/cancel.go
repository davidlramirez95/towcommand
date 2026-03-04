package bookinguc

import (
	"context"

	"github.com/davidlramirez95/towcommand/internal/domain/booking"
	domainerrors "github.com/davidlramirez95/towcommand/internal/domain/errors"
	"github.com/davidlramirez95/towcommand/internal/domain/user"
)

const (
	eventBookingCancelled = "BookingCancelled"
)

// CancelBookingInput holds the data needed to cancel a booking.
type CancelBookingInput struct {
	BookingID string
	CallerID  string
	Reason    string
}

// CancelBookingOutput is the response for a cancel operation.
type CancelBookingOutput struct {
	BookingID       string `json:"bookingId"`
	Status          string `json:"status"`
	CancellationFee int64  `json:"cancellationFee"`
}

// CancelBookingRepo combines the interfaces needed by the cancel use case.
type CancelBookingRepo interface {
	BookingFinder
	BookingStatusUpdater
}

// CancelBookingUseCase orchestrates booking cancellation.
type CancelBookingUseCase struct {
	repo   CancelBookingRepo
	events EventPublisher
}

// NewCancelBookingUseCase constructs a CancelBookingUseCase with its dependencies.
func NewCancelBookingUseCase(repo CancelBookingRepo, events EventPublisher) *CancelBookingUseCase {
	return &CancelBookingUseCase{repo: repo, events: events}
}

// Execute cancels a booking, enforcing ownership and valid state transitions.
func (uc *CancelBookingUseCase) Execute(ctx context.Context, input CancelBookingInput) (*CancelBookingOutput, error) {
	b, err := uc.repo.FindByID(ctx, input.BookingID)
	if err != nil {
		return nil, err
	}
	if b == nil {
		return nil, domainerrors.NewNotFoundError("Booking", input.BookingID)
	}

	if b.CustomerID != input.CallerID {
		return nil, domainerrors.NewForbiddenError("You can only cancel your own bookings")
	}

	if !booking.CanTransition(b.Status, booking.BookingStatusCancelled) {
		return nil, domainerrors.NewBookingNotCancellableError(input.BookingID, string(b.Status))
	}

	cancellationFee := calculateCancellationFee(b.Status)

	metadata := map[string]any{
		"cancelledBy":     input.CallerID,
		"reason":          input.Reason,
		"cancellationFee": cancellationFee,
	}
	if err := uc.repo.UpdateStatus(ctx, input.BookingID, booking.BookingStatusCancelled, metadata); err != nil {
		return nil, err
	}

	_ = uc.events.Publish(ctx, eventSourceBooking, eventBookingCancelled, map[string]any{
		"bookingId":       input.BookingID,
		"customerId":      b.CustomerID,
		"providerId":      b.ProviderID,
		"previousStatus":  b.Status,
		"reason":          input.Reason,
		"cancellationFee": cancellationFee,
	}, &Actor{UserID: input.CallerID, UserType: string(user.UserTypeCustomer)})

	return &CancelBookingOutput{
		BookingID:       input.BookingID,
		Status:          string(booking.BookingStatusCancelled),
		CancellationFee: cancellationFee,
	}, nil
}

// calculateCancellationFee returns the fee in centavos based on current status.
func calculateCancellationFee(status booking.BookingStatus) int64 {
	switch status {
	case booking.BookingStatusPending:
		return 0
	case booking.BookingStatusMatched:
		return 10_000 // ₱100
	case booking.BookingStatusEnRoute:
		return 25_000 // ₱250
	default:
		return 0
	}
}
