package provider

import (
	"context"
	"math"

	"github.com/davidlramirez95/towcommand/internal/domain/errors"
	"github.com/davidlramirez95/towcommand/internal/usecase/port"
)

const (
	maxRadiusKm = 50.0
	maxLimit    = 20
	// avgSpeedKmH is the assumed average tow-truck speed for ETA estimation.
	avgSpeedKmH = 30.0
)

// GetNearbyInput carries the validated fields for a nearby providers query.
type GetNearbyInput struct {
	Lat      float64 `json:"lat"`
	Lng      float64 `json:"lng"`
	RadiusKm float64 `json:"radiusKm"`
	Limit    int     `json:"limit"`
}

// NearbyProvider is a single provider result with distance and ETA.
type NearbyProvider struct {
	ProviderID         string  `json:"providerId"`
	Name               string  `json:"name"`
	TruckType          string  `json:"truckType"`
	Rating             float64 `json:"rating"`
	TotalJobsCompleted int     `json:"totalJobsCompleted"`
	PlateNumber        string  `json:"plateNumber"`
	DistanceKm         float64 `json:"distanceKm"`
	ETAMinutes         int     `json:"etaMinutes"`
}

// GetNearbyOutput is the response for a nearby providers query.
type GetNearbyOutput struct {
	Providers []NearbyProvider `json:"providers"`
	Count     int              `json:"count"`
}

// GetNearbyUseCase orchestrates geospatial provider searches.
type GetNearbyUseCase struct {
	geo    port.GeoCache
	finder port.ProviderFinder
}

// NewGetNearbyUseCase creates a new GetNearbyUseCase.
func NewGetNearbyUseCase(geo port.GeoCache, finder port.ProviderFinder) *GetNearbyUseCase {
	return &GetNearbyUseCase{geo: geo, finder: finder}
}

// Execute finds nearby online providers, enriches them with profile data, and sorts by distance.
func (uc *GetNearbyUseCase) Execute(ctx context.Context, input GetNearbyInput) (*GetNearbyOutput, error) {
	if !isValidPhilippineCoordinate(input.Lat, input.Lng) {
		return nil, errors.NewValidationError("coordinates must be within the Philippines")
	}

	radiusKm := math.Min(input.RadiusKm, maxRadiusKm)
	if radiusKm <= 0 {
		radiusKm = 10
	}

	nearby, err := uc.geo.FindNearbyProviders(ctx, input.Lat, input.Lng, radiusKm)
	if err != nil {
		return nil, errors.NewExternalServiceError("Redis", err)
	}

	limit := input.Limit
	if limit <= 0 || limit > maxLimit {
		limit = maxLimit
	}
	if len(nearby) > limit {
		nearby = nearby[:limit]
	}

	providers := make([]NearbyProvider, 0, len(nearby))
	for _, pd := range nearby {
		p, err := uc.finder.FindByID(ctx, pd.ProviderID)
		if err != nil || p == nil || !p.IsOnline {
			continue
		}
		providers = append(providers, NearbyProvider{
			ProviderID:         pd.ProviderID,
			Name:               p.Name,
			TruckType:          string(p.TruckType),
			Rating:             p.Rating,
			TotalJobsCompleted: p.TotalJobsCompleted,
			PlateNumber:        p.PlateNumber,
			DistanceKm:         math.Round(pd.DistanceKm*10) / 10,
			ETAMinutes:         estimateETAMinutes(pd.DistanceKm),
		})
	}

	return &GetNearbyOutput{
		Providers: providers,
		Count:     len(providers),
	}, nil
}

// estimateETAMinutes estimates travel time based on distance and average speed.
func estimateETAMinutes(distanceKm float64) int {
	minutes := (distanceKm / avgSpeedKmH) * 60
	eta := int(math.Ceil(minutes))
	if eta < 1 {
		return 1
	}
	return eta
}
