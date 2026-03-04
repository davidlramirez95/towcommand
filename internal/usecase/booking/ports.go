// Package bookinguc implements booking use cases following CLEAN architecture.
// Each use case declares only the port interfaces it needs (ISP).
package bookinguc

import (
	"context"

	"github.com/davidlramirez95/towcommand/internal/domain/booking"
	"github.com/davidlramirez95/towcommand/internal/usecase/port"
)

// BookingSaver persists a new booking.
type BookingSaver interface {
	Save(ctx context.Context, b *booking.Booking) error
}

// BookingFinder retrieves a booking by its ID. Returns nil if not found.
type BookingFinder interface {
	FindByID(ctx context.Context, bookingID string) (*booking.Booking, error)
}

// BookingByUserLister lists bookings for a given user.
type BookingByUserLister interface {
	FindByUser(ctx context.Context, userID string, limit int32) ([]booking.Booking, error)
}

// BookingStatusUpdater changes a booking's status and records the transition.
type BookingStatusUpdater interface {
	UpdateStatus(ctx context.Context, bookingID string, status booking.BookingStatus, metadata map[string]any) error
}

// BookingByStatusLister lists bookings by their current status.
type BookingByStatusLister interface {
	FindByStatus(ctx context.Context, status booking.BookingStatus, limit int32) ([]booking.Booking, error)
}

// EventPublisher publishes domain events to an event bus.
type EventPublisher interface {
	Publish(ctx context.Context, source, detailType string, detail any, actor *port.Actor) error
}

// Actor is a type alias for port.Actor to avoid import stuttering in use case code.
type Actor = port.Actor
