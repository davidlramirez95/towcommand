// Package safetyuc implements safety use cases following CLEAN architecture.
// Each use case declares only the port interfaces it needs (ISP).
package safetyuc

import (
	"context"
	"time"

	"github.com/davidlramirez95/towcommand/internal/domain/booking"
	"github.com/davidlramirez95/towcommand/internal/domain/provider"
	"github.com/davidlramirez95/towcommand/internal/domain/safety"
	"github.com/davidlramirez95/towcommand/internal/usecase/port"
)

// SOSSaver persists a new SOS alert.
type SOSSaver interface {
	Save(ctx context.Context, alert *safety.SOSAlert) error
}

// SOSFinder retrieves an SOS alert by its ID.
type SOSFinder interface {
	FindByID(ctx context.Context, alertID string) (*safety.SOSAlert, error)
}

// SOSResolver marks an SOS alert as resolved.
type SOSResolver interface {
	Resolve(ctx context.Context, alertID string, resolvedBy string, resolvedAt time.Time) error
}

// SOSActiveLister queries for active (unresolved) SOS alerts.
type SOSActiveLister interface {
	FindActive(ctx context.Context, limit int32) ([]safety.SOSAlert, error)
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

// Actor is a type alias for port.Actor to avoid import stuttering in use case code.
type Actor = port.Actor
