// Package ratinguc implements rating use cases following CLEAN architecture.
// Each use case declares only the port interfaces it needs (ISP).
package ratinguc

import (
	"context"

	"github.com/davidlramirez95/towcommand/internal/domain/booking"
	"github.com/davidlramirez95/towcommand/internal/domain/provider"
	"github.com/davidlramirez95/towcommand/internal/domain/rating"
	"github.com/davidlramirez95/towcommand/internal/usecase/port"
)

// RatingSaver persists a new rating.
type RatingSaver interface {
	Save(ctx context.Context, r *rating.Rating) error
}

// RatingByBookingFinder retrieves the rating for a given booking.
type RatingByBookingFinder interface {
	FindByBooking(ctx context.Context, bookingID string) (*rating.Rating, error)
}

// RatingByProviderLister lists ratings for a given provider, ordered by date descending.
type RatingByProviderLister interface {
	FindByProvider(ctx context.Context, providerID string, limit int32) ([]rating.Rating, error)
}

// BookingFinder retrieves a booking by its ID.
type BookingFinder interface {
	FindByID(ctx context.Context, bookingID string) (*booking.Booking, error)
}

// ProviderFinder retrieves a provider by its ID.
type ProviderFinder interface {
	FindByID(ctx context.Context, providerID string) (*provider.Provider, error)
}

// ProviderSaver persists a provider entity.
type ProviderSaver interface {
	Save(ctx context.Context, p *provider.Provider) error
}

// EventPublisher publishes domain events to an event bus.
type EventPublisher interface {
	Publish(ctx context.Context, source, detailType string, detail any, actor *port.Actor) error
}

// Actor is a type alias for port.Actor to avoid import stuttering in use case code.
type Actor = port.Actor
