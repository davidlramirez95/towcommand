package bookinguc

import (
	"context"

	"github.com/davidlramirez95/towcommand/internal/domain/booking"
	domainerrors "github.com/davidlramirez95/towcommand/internal/domain/errors"
)

const (
	eventBookingStatusChanged = "BookingStatusChanged"
	eventBookingCompleted     = "BookingCompleted"
)

// providerStatuses are statuses that only the assigned provider or admin can set.
var providerStatuses = map[booking.BookingStatus]bool{
	booking.BookingStatusEnRoute:         true,
	booking.BookingStatusArrived:         true,
	booking.BookingStatusConditionReport: true,
	booking.BookingStatusLoading:         true,
	booking.BookingStatusInTransit:       true,
	booking.BookingStatusArrivedDropoff:  true,
	booking.BookingStatusCompleted:       true,
}

// UpdateBookingStatusInput holds the data needed to update a booking's status.
type UpdateBookingStatusInput struct {
	BookingID  string
	CallerID   string
	CallerType string
	NewStatus  booking.BookingStatus
	Metadata   map[string]any
}

// UpdateBookingStatusOutput is the response for a status update.
type UpdateBookingStatusOutput struct {
	BookingID      string `json:"bookingId"`
	PreviousStatus string `json:"previousStatus"`
	Status         string `json:"status"`
}

// UpdateStatusRepo combines the interfaces needed by the update-status use case.
type UpdateStatusRepo interface {
	BookingFinder
	BookingStatusUpdater
}

// UpdateBookingStatusUseCase orchestrates booking status transitions.
type UpdateBookingStatusUseCase struct {
	repo   UpdateStatusRepo
	events EventPublisher
}

// NewUpdateBookingStatusUseCase constructs an UpdateBookingStatusUseCase.
func NewUpdateBookingStatusUseCase(repo UpdateStatusRepo, events EventPublisher) *UpdateBookingStatusUseCase {
	return &UpdateBookingStatusUseCase{repo: repo, events: events}
}

// Execute updates a booking's status, enforcing authorization and valid transitions.
func (uc *UpdateBookingStatusUseCase) Execute(ctx context.Context, input UpdateBookingStatusInput) (*UpdateBookingStatusOutput, error) {
	b, err := uc.repo.FindByID(ctx, input.BookingID)
	if err != nil {
		return nil, err
	}
	if b == nil {
		return nil, domainerrors.NewNotFoundError("Booking", input.BookingID)
	}

	// Authorization: provider statuses require assigned provider or admin.
	if providerStatuses[input.NewStatus] {
		isAssigned := b.ProviderID == input.CallerID
		if !isAssigned && !isAdmin(input.CallerType) {
			return nil, domainerrors.NewForbiddenError("Only the assigned provider can update this status")
		}
	} else if !isAdmin(input.CallerType) {
		return nil, domainerrors.NewForbiddenError("Only admins can set this status")
	}

	if !booking.CanTransition(b.Status, input.NewStatus) {
		return nil, domainerrors.NewInvalidStatusTransitionError(string(b.Status), string(input.NewStatus))
	}

	metadata := map[string]any{"changedBy": input.CallerID}
	for k, v := range input.Metadata {
		metadata[k] = v
	}
	if err := uc.repo.UpdateStatus(ctx, input.BookingID, input.NewStatus, metadata); err != nil {
		return nil, err
	}

	actor := &Actor{UserID: input.CallerID, UserType: input.CallerType}
	_ = uc.events.Publish(ctx, eventSourceBooking, eventBookingStatusChanged, map[string]any{
		"bookingId":      input.BookingID,
		"previousStatus": b.Status,
		"newStatus":      input.NewStatus,
		"changedBy":      input.CallerID,
		"metadata":       input.Metadata,
	}, actor)

	if input.NewStatus == booking.BookingStatusCompleted && b.ProviderID != "" {
		_ = uc.events.Publish(ctx, eventSourceBooking, eventBookingCompleted, map[string]any{
			"bookingId":  input.BookingID,
			"customerId": b.CustomerID,
			"providerId": b.ProviderID,
			"price":      b.Price,
		}, actor)
	}

	return &UpdateBookingStatusOutput{
		BookingID:      input.BookingID,
		PreviousStatus: string(b.Status),
		Status:         string(input.NewStatus),
	}, nil
}
