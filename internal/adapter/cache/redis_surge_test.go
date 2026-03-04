package cache_test

import (
	"context"
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/davidlramirez95/towcommand/internal/adapter/cache"
)

func newSurgeTestClient(t *testing.T) (*cache.RedisSurgeCache, *miniredis.Miniredis) {
	t.Helper()
	mr := miniredis.RunT(t)
	client := cache.NewRedisClient(cache.Options{
		Host: mr.Host(),
		Port: mr.Server().Addr().Port,
	})
	t.Cleanup(func() { client.Close() })
	return cache.NewRedisSurgeCache(client), mr
}

func TestRedisSurgeCache_IncrementAndGetDemand(t *testing.T) {
	surge, _ := newSurgeTestClient(t)
	ctx := context.Background()

	tests := []struct {
		name       string
		areaID     string
		increments int
		wantCount  int
	}{
		{
			name:       "single increment",
			areaID:     "area-1",
			increments: 1,
			wantCount:  1,
		},
		{
			name:       "multiple increments",
			areaID:     "area-2",
			increments: 5,
			wantCount:  5,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for i := 0; i < tt.increments; i++ {
				err := surge.IncrementAreaDemand(ctx, tt.areaID, 5*time.Minute)
				require.NoError(t, err)
			}

			got, err := surge.GetAreaDemand(ctx, tt.areaID)
			require.NoError(t, err)
			assert.Equal(t, tt.wantCount, got)
		})
	}
}

func TestRedisSurgeCache_GetAreaDemand_NotFound(t *testing.T) {
	surge, _ := newSurgeTestClient(t)
	ctx := context.Background()

	got, err := surge.GetAreaDemand(ctx, "nonexistent-area")
	require.NoError(t, err)
	assert.Equal(t, 0, got)
}

func TestRedisSurgeCache_TTLExpiry(t *testing.T) {
	surge, mr := newSurgeTestClient(t)
	ctx := context.Background()

	require.NoError(t, surge.IncrementAreaDemand(ctx, "area-exp", 1*time.Second))

	got, err := surge.GetAreaDemand(ctx, "area-exp")
	require.NoError(t, err)
	assert.Equal(t, 1, got)

	mr.FastForward(2 * time.Second)

	got, err = surge.GetAreaDemand(ctx, "area-exp")
	require.NoError(t, err)
	assert.Equal(t, 0, got)
}

func TestRedisSurgeCache_SeparateAreas(t *testing.T) {
	surge, _ := newSurgeTestClient(t)
	ctx := context.Background()

	for i := 0; i < 3; i++ {
		require.NoError(t, surge.IncrementAreaDemand(ctx, "manila", 5*time.Minute))
	}
	require.NoError(t, surge.IncrementAreaDemand(ctx, "cebu", 5*time.Minute))

	manila, err := surge.GetAreaDemand(ctx, "manila")
	require.NoError(t, err)
	assert.Equal(t, 3, manila)

	cebu, err := surge.GetAreaDemand(ctx, "cebu")
	require.NoError(t, err)
	assert.Equal(t, 1, cebu)
}
