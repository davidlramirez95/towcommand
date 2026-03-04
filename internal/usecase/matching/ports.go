// Package matching implements the provider matching use case following CLEAN architecture.
// Each use case declares only the port interfaces it needs (ISP).
package matching

import (
	"context"

	"github.com/davidlramirez95/towcommand/internal/domain/booking"
	"github.com/davidlramirez95/towcommand/internal/domain/provider"
	"github.com/davidlramirez95/towcommand/internal/usecase/port"
)

// BookingFinder retrieves a booking by its ID.
type BookingFinder interface {
	FindByID(ctx context.Context, bookingID string) (*booking.Booking, error)
}

// ProviderFinder retrieves a provider by their ID.
type ProviderFinder interface {
	FindByID(ctx context.Context, providerID string) (*provider.Provider, error)
}

// GeoCache performs geospatial queries for nearby providers.
type GeoCache interface {
	FindNearbyProviders(ctx context.Context, lat, lng, radiusKm float64) ([]port.ProviderDistance, error)
}

// SurgeCache reads area demand counters for surge mode detection.
type SurgeCache interface {
	GetAreaDemand(ctx context.Context, areaID string) (int, error)
}
