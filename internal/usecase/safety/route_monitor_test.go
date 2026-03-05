package safetyuc

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCheckRouteDeviation(t *testing.T) {
	// Define a route corridor: pickup (Manila) to dropoff (Quezon City), roughly north.
	pickup := RoutePoint{Lat: 14.5995, Lng: 120.9842}
	dropoff := RoutePoint{Lat: 14.6760, Lng: 121.0437}

	tests := []struct {
		name        string
		current     RoutePoint
		pickup      RoutePoint
		dropoff     RoutePoint
		thresholdKm float64
		wantDeviate bool
	}{
		{
			name:        "on route: at pickup",
			current:     RoutePoint{Lat: 14.5995, Lng: 120.9842},
			pickup:      pickup,
			dropoff:     dropoff,
			thresholdKm: DefaultDeviationThresholdKm,
			wantDeviate: false,
		},
		{
			name:        "on route: at dropoff",
			current:     RoutePoint{Lat: 14.6760, Lng: 121.0437},
			pickup:      pickup,
			dropoff:     dropoff,
			thresholdKm: DefaultDeviationThresholdKm,
			wantDeviate: false,
		},
		{
			name:        "on route: midpoint",
			current:     RoutePoint{Lat: 14.6378, Lng: 121.0140},
			pickup:      pickup,
			dropoff:     dropoff,
			thresholdKm: DefaultDeviationThresholdKm,
			wantDeviate: false,
		},
		{
			name:        "slightly off: under threshold",
			current:     RoutePoint{Lat: 14.6400, Lng: 121.0000},
			pickup:      pickup,
			dropoff:     dropoff,
			thresholdKm: DefaultDeviationThresholdKm,
			wantDeviate: false,
		},
		{
			name:        "significantly off: over threshold",
			current:     RoutePoint{Lat: 14.5500, Lng: 121.1000},
			pickup:      pickup,
			dropoff:     dropoff,
			thresholdKm: DefaultDeviationThresholdKm,
			wantDeviate: true,
		},
		{
			name:        "far away: totally off route",
			current:     RoutePoint{Lat: 14.8000, Lng: 121.3000},
			pickup:      pickup,
			dropoff:     dropoff,
			thresholdKm: DefaultDeviationThresholdKm,
			wantDeviate: true,
		},
		{
			name:        "near pickup endpoint but off line",
			current:     RoutePoint{Lat: 14.5900, Lng: 120.9600},
			pickup:      pickup,
			dropoff:     dropoff,
			thresholdKm: DefaultDeviationThresholdKm,
			wantDeviate: true, // > 2km from segment
		},
		{
			name:        "custom threshold: tight corridor",
			current:     RoutePoint{Lat: 14.6400, Lng: 121.0000},
			pickup:      pickup,
			dropoff:     dropoff,
			thresholdKm: 0.5, // Very tight corridor
			wantDeviate: true,
		},
		{
			name:        "custom threshold: wide corridor",
			current:     RoutePoint{Lat: 14.5500, Lng: 121.1000},
			pickup:      pickup,
			dropoff:     dropoff,
			thresholdKm: 20.0, // Very wide corridor
			wantDeviate: false,
		},
		{
			name:        "degenerate: pickup equals dropoff, on point",
			current:     RoutePoint{Lat: 14.5995, Lng: 120.9842},
			pickup:      RoutePoint{Lat: 14.5995, Lng: 120.9842},
			dropoff:     RoutePoint{Lat: 14.5995, Lng: 120.9842},
			thresholdKm: 1.0,
			wantDeviate: false,
		},
		{
			name:        "degenerate: pickup equals dropoff, away from point",
			current:     RoutePoint{Lat: 14.6200, Lng: 121.0000},
			pickup:      RoutePoint{Lat: 14.5995, Lng: 120.9842},
			dropoff:     RoutePoint{Lat: 14.5995, Lng: 120.9842},
			thresholdKm: 1.0,
			wantDeviate: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := CheckRouteDeviation(tt.current, tt.pickup, tt.dropoff, tt.thresholdKm)
			assert.Equal(t, tt.wantDeviate, got)
		})
	}
}

func TestPointToSegmentDistanceKm(t *testing.T) {
	// Test perpendicular distance: point directly east of segment midpoint.
	a := RoutePoint{Lat: 14.5, Lng: 121.0}
	b := RoutePoint{Lat: 14.7, Lng: 121.0}
	p := RoutePoint{Lat: 14.6, Lng: 121.1} // East of the midpoint

	dist := pointToSegmentDistanceKm(p, a, b)
	// At latitude ~14.6, 0.1 degree longitude is roughly 10.7 km.
	assert.InDelta(t, 10.7, dist, 1.0)

	// Point projects before segment start -> distance to start.
	pBefore := RoutePoint{Lat: 14.4, Lng: 121.0}
	distBefore := pointToSegmentDistanceKm(pBefore, a, b)
	// 0.1 degree latitude is roughly 11.1 km.
	assert.InDelta(t, 11.1, distBefore, 1.0)

	// Point projects after segment end -> distance to end.
	pAfter := RoutePoint{Lat: 14.8, Lng: 121.0}
	distAfter := pointToSegmentDistanceKm(pAfter, a, b)
	assert.InDelta(t, 11.1, distAfter, 1.0)
}
