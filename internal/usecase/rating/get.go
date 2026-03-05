package ratinguc

import (
	"context"

	domainerrors "github.com/davidlramirez95/towcommand/internal/domain/errors"
	"github.com/davidlramirez95/towcommand/internal/domain/rating"
)

// GetBookingRatingUseCase retrieves the rating for a given booking.
type GetBookingRatingUseCase struct {
	ratings RatingByBookingFinder
}

// NewGetBookingRatingUseCase constructs a GetBookingRatingUseCase with its dependencies.
func NewGetBookingRatingUseCase(ratings RatingByBookingFinder) *GetBookingRatingUseCase {
	return &GetBookingRatingUseCase{ratings: ratings}
}

// Execute retrieves the rating for the specified booking.
// Returns a NotFoundError if no rating exists for the booking.
func (uc *GetBookingRatingUseCase) Execute(ctx context.Context, bookingID string) (*rating.Rating, error) {
	r, err := uc.ratings.FindByBooking(ctx, bookingID)
	if err != nil {
		return nil, err
	}
	if r == nil {
		return nil, domainerrors.NewNotFoundError("rating", bookingID)
	}
	return r, nil
}
