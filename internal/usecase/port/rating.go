package port

import (
	"context"

	"github.com/davidlramirez95/towcommand/internal/domain/rating"
)

// RatingSaver persists a new rating.
type RatingSaver interface {
	Save(ctx context.Context, r *rating.Rating) error
}

// RatingByBookingFinder retrieves the rating for a given booking.
type RatingByBookingFinder interface {
	FindByBooking(ctx context.Context, bookingID string) (*rating.Rating, error)
}

// RatingByProviderLister lists ratings for a given provider via GSI1, ordered by date descending.
type RatingByProviderLister interface {
	FindByProvider(ctx context.Context, providerID string, limit int32) ([]rating.Rating, error)
}
