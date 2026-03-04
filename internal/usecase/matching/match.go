package matching

import (
	"context"
	"fmt"

	"github.com/davidlramirez95/towcommand/internal/domain/errors"
	"github.com/davidlramirez95/towcommand/internal/domain/provider"
)

// surgeDemandThreshold is the minimum area demand count that activates surge mode.
// In surge mode, distance weight increases from 40% to 50% to favour closer providers.
const surgeDemandThreshold = 10

// MatchResult holds the output of a matching operation.
type MatchResult struct {
	Scores    []provider.MatchScore `json:"scores"`
	SurgeMode bool                  `json:"surgeMode"`
}

// MatchBookingUseCase orchestrates the matching of a booking to available providers.
type MatchBookingUseCase struct {
	bookings  BookingFinder
	providers ProviderFinder
	geo       GeoCache
	surge     SurgeCache
}

// NewMatchBookingUseCase creates a MatchBookingUseCase with its dependencies.
func NewMatchBookingUseCase(bookings BookingFinder, providers ProviderFinder, geo GeoCache, surge SurgeCache) *MatchBookingUseCase {
	return &MatchBookingUseCase{
		bookings:  bookings,
		providers: providers,
		geo:       geo,
		surge:     surge,
	}
}

// Execute matches a booking to the best available providers.
//
// It looks up the booking by ID, detects surge mode via area demand, performs
// a cascade radius search to find nearby provider candidates, ranks them using
// the weighted scoring algorithm, and returns the ranked results.
func (uc *MatchBookingUseCase) Execute(ctx context.Context, bookingID string) (*MatchResult, error) {
	b, err := uc.bookings.FindByID(ctx, bookingID)
	if err != nil {
		return nil, fmt.Errorf("finding booking %s: %w", bookingID, err)
	}
	if b == nil {
		return nil, errors.NewNotFoundError("booking", bookingID)
	}

	// Derive area ID from pickup coordinates for surge lookup.
	areaID := fmt.Sprintf("%.2f,%.2f", b.PickupLocation.Lat, b.PickupLocation.Lng)
	demand, err := uc.surge.GetAreaDemand(ctx, areaID)
	if err != nil {
		// Surge cache failure is non-fatal; default to normal mode.
		demand = 0
	}
	surgeMode := demand >= surgeDemandThreshold

	candidates, err := cascadeSearch(ctx, uc.geo, uc.providers, b.PickupLocation.Lat, b.PickupLocation.Lng)
	if err != nil {
		return nil, errors.NewExternalServiceError("GeoCache", err)
	}
	if len(candidates) == 0 {
		return nil, errors.NewProviderUnavailableError()
	}

	weightKg := provider.WeightClassToKg(b.WeightClass)
	scores := provider.RankProviders(candidates, b.ServiceType, weightKg, surgeMode)
	if len(scores) == 0 {
		return nil, errors.NewProviderUnavailableError()
	}

	return &MatchResult{
		Scores:    scores,
		SurgeMode: surgeMode,
	}, nil
}
