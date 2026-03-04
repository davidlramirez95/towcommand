// Package port defines interfaces that adapters must implement.
// Use-case services depend on these ports, not on concrete adapter types.
package port

import (
	"context"
	"time"
)

// ProviderDistance represents a provider with their distance from a query point.
type ProviderDistance struct {
	ProviderID string
	DistanceKm float64
}

// GeoCache manages geospatial indexing of provider locations.
type GeoCache interface {
	AddProviderLocation(ctx context.Context, providerID string, lat, lng float64) error
	FindNearbyProviders(ctx context.Context, lat, lng, radiusKm float64) ([]ProviderDistance, error)
	RemoveProvider(ctx context.Context, providerID string) error
}

// SessionCache manages WebSocket connection mappings.
// It supports both forward (userID -> connectionID) and reverse
// (connectionID -> userID) lookups to enable clean disconnect handling.
type SessionCache interface {
	MapConnection(ctx context.Context, userID, connectionID string, ttl time.Duration) error
	GetConnection(ctx context.Context, userID string) (string, error)
	RemoveConnection(ctx context.Context, userID string) error
	MapReverseConnection(ctx context.Context, connectionID, userID string, ttl time.Duration) error
	GetUserByConnection(ctx context.Context, connectionID string) (string, error)
	RemoveReverseConnection(ctx context.Context, connectionID string) error
}

// OTPCache stores hashed OTP codes with automatic expiration.
type OTPCache interface {
	StoreOTP(ctx context.Context, bookingID, otpType, hashedOTP string, ttl time.Duration) error
	GetOTP(ctx context.Context, bookingID, otpType string) (string, error)
	DeleteOTP(ctx context.Context, bookingID, otpType string) error
}

// RateLimiter performs sliding-window rate-limit checks.
type RateLimiter interface {
	CheckRateLimit(ctx context.Context, key string, maxRequests, windowSec int) (allowed bool, remaining int, err error)
}

// SurgeCache tracks area demand counters for surge pricing.
type SurgeCache interface {
	IncrementAreaDemand(ctx context.Context, areaID string, ttl time.Duration) error
	GetAreaDemand(ctx context.Context, areaID string) (int, error)
}
