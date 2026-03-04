package port

import (
	"context"

	"github.com/davidlramirez95/towcommand/internal/domain/booking"
)

// BookingSaver persists a new booking.
type BookingSaver interface {
	Save(ctx context.Context, b *booking.Booking) error
}

// BookingFinder retrieves a booking by its ID.
type BookingFinder interface {
	FindByID(ctx context.Context, bookingID string) (*booking.Booking, error)
}

// BookingByUserLister lists bookings for a given user, ordered by creation date descending.
type BookingByUserLister interface {
	FindByUser(ctx context.Context, userID string, limit int32) ([]booking.Booking, error)
}

// BookingStatusUpdater changes a booking's status and records the transition history.
type BookingStatusUpdater interface {
	UpdateStatus(ctx context.Context, bookingID string, status booking.BookingStatus, metadata map[string]any) error
}

// BookingByStatusLister lists bookings by their current status, ordered by creation date descending.
type BookingByStatusLister interface {
	FindByStatus(ctx context.Context, status booking.BookingStatus, limit int32) ([]booking.Booking, error)
}
