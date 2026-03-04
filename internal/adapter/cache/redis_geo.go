package cache

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"

	"github.com/davidlramirez95/towcommand/internal/usecase/port"
)

const geoKey = "provider:geo"

// RedisGeoCache implements port.GeoCache using Redis geospatial commands.
type RedisGeoCache struct {
	client *redis.Client
}

// NewRedisGeoCache creates a new RedisGeoCache.
func NewRedisGeoCache(client *redis.Client) *RedisGeoCache {
	return &RedisGeoCache{client: client}
}

// Compile-time interface check.
var _ port.GeoCache = (*RedisGeoCache)(nil)

// AddProviderLocation adds or updates a provider's geospatial position.
func (g *RedisGeoCache) AddProviderLocation(ctx context.Context, providerID string, lat, lng float64) error {
	return g.client.GeoAdd(ctx, geoKey, &redis.GeoLocation{
		Name:      providerID,
		Longitude: lng,
		Latitude:  lat,
	}).Err()
}

// FindNearbyProviders returns providers within radiusKm, sorted by distance ascending.
// NOTE: Uses GeoRadius because miniredis v2.37.0 does not support GEOSEARCH yet.
// TODO: Migrate to GeoSearchLocation once miniredis adds GEOSEARCH support.
func (g *RedisGeoCache) FindNearbyProviders(ctx context.Context, lat, lng, radiusKm float64) ([]port.ProviderDistance, error) {
	//nolint:staticcheck // GeoRadius is deprecated but miniredis lacks GEOSEARCH support.
	results, err := g.client.GeoRadius(ctx, geoKey, lng, lat, &redis.GeoRadiusQuery{
		Radius:   radiusKm,
		Unit:     "km",
		WithDist: true,
		Sort:     "ASC",
	}).Result()
	if err != nil {
		return nil, fmt.Errorf("finding nearby providers: %w", err)
	}

	providers := make([]port.ProviderDistance, len(results))
	for i, r := range results {
		providers[i] = port.ProviderDistance{
			ProviderID: r.Name,
			DistanceKm: r.Dist,
		}
	}
	return providers, nil
}

// RemoveProvider removes a provider from the geo index.
func (g *RedisGeoCache) RemoveProvider(ctx context.Context, providerID string) error {
	return g.client.ZRem(ctx, geoKey, providerID).Err()
}
