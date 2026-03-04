package matching

import (
	"context"
	"fmt"

	"github.com/davidlramirez95/towcommand/internal/domain/provider"
)

// cascadeRadii defines the expanding search radii in km.
// The search stops at the first radius that yields ranked candidates.
var cascadeRadii = []float64{5, 10, 20, 30}

// cascadeSearch expands through progressively larger radii to find nearby
// providers. For each radius, it queries the GeoCache, hydrates full provider
// records from the ProviderFinder, and converts them to MatchCandidates.
// It returns candidates from the first radius that yields at least one result.
func cascadeSearch(ctx context.Context, geo GeoCache, finder ProviderFinder, lat, lng float64) ([]provider.MatchCandidate, error) {
	for _, radius := range cascadeRadii {
		nearby, err := geo.FindNearbyProviders(ctx, lat, lng, radius)
		if err != nil {
			return nil, fmt.Errorf("finding nearby providers at radius %.0fkm: %w", radius, err)
		}
		if len(nearby) == 0 {
			continue
		}

		candidates := make([]provider.MatchCandidate, 0, len(nearby))
		for _, pd := range nearby {
			p, err := finder.FindByID(ctx, pd.ProviderID)
			if err != nil || p == nil {
				continue
			}
			candidates = append(candidates, provider.MatchCandidate{
				ProviderID:          p.ProviderID,
				TrustTier:           p.TrustTier,
				AcceptanceRate:      p.AcceptanceRate,
				TruckType:           p.TruckType,
				MaxWeightCapacityKg: p.MaxWeightCapacityKg,
				ActiveJobCount:      0, // TODO(#22): integrate active job counter
				DistanceKm:          pd.DistanceKm,
				IsOnline:            p.IsOnline,
			})
		}
		if len(candidates) > 0 {
			return candidates, nil
		}
	}
	return nil, nil
}
