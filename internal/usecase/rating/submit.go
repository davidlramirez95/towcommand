package ratinguc

import (
	"context"
	"time"

	"github.com/davidlramirez95/towcommand/internal/domain/booking"
	domainerrors "github.com/davidlramirez95/towcommand/internal/domain/errors"
	"github.com/davidlramirez95/towcommand/internal/domain/event"
	"github.com/davidlramirez95/towcommand/internal/domain/rating"
)

// SubmitRatingInput holds the data needed to submit a rating for a completed booking.
type SubmitRatingInput struct {
	BookingID  string
	CustomerID string
	Score      int
	Comment    string
	Tags       []string
}

// SubmitRatingUseCase orchestrates rating submission, provider average recalculation,
// and event publishing.
type SubmitRatingUseCase struct {
	ratings         RatingSaver
	ratingsFinder   RatingByBookingFinder
	ratingsLister   RatingByProviderLister
	bookings        BookingFinder
	providers       ProviderFinder
	providerSaver   ProviderSaver
	events          EventPublisher
	now             func() time.Time
}

// NewSubmitRatingUseCase constructs a SubmitRatingUseCase with its dependencies.
func NewSubmitRatingUseCase(
	ratings RatingSaver,
	ratingsFinder RatingByBookingFinder,
	ratingsLister RatingByProviderLister,
	bookings BookingFinder,
	providers ProviderFinder,
	providerSaver ProviderSaver,
	events EventPublisher,
) *SubmitRatingUseCase {
	return &SubmitRatingUseCase{
		ratings:       ratings,
		ratingsFinder: ratingsFinder,
		ratingsLister: ratingsLister,
		bookings:      bookings,
		providers:     providers,
		providerSaver: providerSaver,
		events:        events,
		now:           func() time.Time { return time.Now().UTC() },
	}
}

// Execute validates the input, creates a rating, recalculates the provider's average
// rating from all existing ratings (self-healing), and publishes a RatingSubmitted event.
func (uc *SubmitRatingUseCase) Execute(ctx context.Context, input *SubmitRatingInput) (*rating.Rating, error) {
	// 1. Find booking
	b, err := uc.bookings.FindByID(ctx, input.BookingID)
	if err != nil {
		return nil, err
	}
	if b == nil {
		return nil, domainerrors.NewNotFoundError("booking", input.BookingID)
	}

	// 2. Verify booking is completed
	if b.Status != booking.BookingStatusCompleted {
		return nil, domainerrors.NewConflictError("booking is not completed")
	}

	// 3. Verify caller is the booking's customer
	if input.CustomerID != b.CustomerID {
		return nil, domainerrors.NewForbiddenError("only the booking customer can submit a rating")
	}

	// 4. Check for existing rating (duplicate guard)
	existing, err := uc.ratingsFinder.FindByBooking(ctx, input.BookingID)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, domainerrors.NewConflictError("rating already submitted for this booking")
	}

	// 5. Validate score range
	if input.Score < 1 || input.Score > 5 {
		return nil, domainerrors.NewValidationError("score must be between 1 and 5")
	}

	// 6. Find provider
	p, err := uc.providers.FindByID(ctx, b.ProviderID)
	if err != nil {
		return nil, err
	}
	if p == nil {
		return nil, domainerrors.NewNotFoundError("provider", b.ProviderID)
	}

	// 7. Create rating entity
	r := &rating.Rating{
		BookingID:  input.BookingID,
		CustomerID: input.CustomerID,
		ProviderID: b.ProviderID,
		Score:      input.Score,
		Comment:    input.Comment,
		Tags:       input.Tags,
		CreatedAt:  uc.now(),
	}

	// 8. Save rating
	if err := uc.ratings.Save(ctx, r); err != nil {
		return nil, err
	}

	// 9. Recalculate provider average from all ratings (self-healing)
	allRatings, err := uc.ratingsLister.FindByProvider(ctx, b.ProviderID, 1000)
	if err == nil && len(allRatings) > 0 {
		var total int
		for i := range allRatings {
			total += allRatings[i].Score
		}
		p.Rating = float64(total) / float64(len(allRatings))
		_ = uc.providerSaver.Save(ctx, p)
	}

	// 10. Publish RatingSubmitted event (best-effort)
	_ = uc.events.Publish(ctx, event.SourceRating, event.RatingSubmitted, map[string]any{
		"bookingId":  r.BookingID,
		"providerId": r.ProviderID,
		"score":      r.Score,
	}, &Actor{UserID: input.CustomerID, UserType: "customer"})

	return r, nil
}
