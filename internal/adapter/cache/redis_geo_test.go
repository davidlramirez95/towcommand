package cache_test

import (
	"context"
	"testing"

	"github.com/alicebob/miniredis/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/davidlramirez95/towcommand/internal/adapter/cache"
)

// Manila coordinates for testing (NCR area).
const (
	manilaLat = 14.5995
	manilaLng = 120.9842
)

func newGeoTestClient(t *testing.T) (*cache.RedisGeoCache, *miniredis.Miniredis) {
	t.Helper()
	mr := miniredis.RunT(t)
	client := cache.NewRedisClient(cache.Options{
		Host: mr.Host(),
		Port: mr.Server().Addr().Port,
	})
	t.Cleanup(func() { _ = client.Close() })
	return cache.NewRedisGeoCache(client), mr
}

func TestRedisGeoCache_AddProviderLocation(t *testing.T) {
	geo, _ := newGeoTestClient(t)
	ctx := context.Background()

	tests := []struct {
		name       string
		providerID string
		lat        float64
		lng        float64
	}{
		{
			name:       "add provider in Manila",
			providerID: "provider-1",
			lat:        manilaLat,
			lng:        manilaLng,
		},
		{
			name:       "add provider in Quezon City",
			providerID: "provider-2",
			lat:        14.6760,
			lng:        121.0437,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := geo.AddProviderLocation(ctx, tt.providerID, tt.lat, tt.lng)
			require.NoError(t, err)
		})
	}
}

func TestRedisGeoCache_FindNearbyProviders(t *testing.T) {
	geo, _ := newGeoTestClient(t)
	ctx := context.Background()

	// Seed providers at known locations around Manila.
	providers := []struct {
		id  string
		lat float64
		lng float64
	}{
		{"close-provider", 14.6000, 120.9850},    // ~55m away
		{"medium-provider", 14.6100, 120.9900},   // ~1.2km away
		{"far-provider", 14.7000, 121.0500},      // ~13km away
		{"very-far-provider", 15.0000, 121.5000}, // ~70km away
	}
	for _, p := range providers {
		require.NoError(t, geo.AddProviderLocation(ctx, p.id, p.lat, p.lng))
	}

	tests := []struct {
		name      string
		lat       float64
		lng       float64
		radiusKm  float64
		wantCount int
		wantFirst string
	}{
		{
			name:      "small radius finds closest",
			lat:       manilaLat,
			lng:       manilaLng,
			radiusKm:  2.0,
			wantCount: 2,
			wantFirst: "close-provider",
		},
		{
			name:      "large radius finds all nearby",
			lat:       manilaLat,
			lng:       manilaLng,
			radiusKm:  20.0,
			wantCount: 3,
			wantFirst: "close-provider",
		},
		{
			name:      "no providers in radius",
			lat:       10.0,
			lng:       120.0,
			radiusKm:  1.0,
			wantCount: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			results, err := geo.FindNearbyProviders(ctx, tt.lat, tt.lng, tt.radiusKm)
			require.NoError(t, err)
			assert.Len(t, results, tt.wantCount)

			if tt.wantCount > 0 {
				// Verify sorted by distance ascending.
				assert.Equal(t, tt.wantFirst, results[0].ProviderID)
				for i := 1; i < len(results); i++ {
					assert.GreaterOrEqual(t, results[i].DistanceKm, results[i-1].DistanceKm)
				}
			}
		})
	}
}

func TestRedisGeoCache_RemoveProvider(t *testing.T) {
	geo, _ := newGeoTestClient(t)
	ctx := context.Background()

	require.NoError(t, geo.AddProviderLocation(ctx, "provider-rm", manilaLat, manilaLng))

	// Verify provider is findable.
	results, err := geo.FindNearbyProviders(ctx, manilaLat, manilaLng, 1.0)
	require.NoError(t, err)
	assert.Len(t, results, 1)

	// Remove and verify.
	err = geo.RemoveProvider(ctx, "provider-rm")
	require.NoError(t, err)

	results, err = geo.FindNearbyProviders(ctx, manilaLat, manilaLng, 1.0)
	require.NoError(t, err)
	assert.Empty(t, results)
}
