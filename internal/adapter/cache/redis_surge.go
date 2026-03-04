package cache

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"

	"github.com/davidlramirez95/towcommand/internal/usecase/port"
)

const surgeKeyPrefix = "surge:demand:"

// RedisSurgeCache implements port.SurgeCache using Redis INCR counters with TTL.
type RedisSurgeCache struct {
	client *redis.Client
}

// NewRedisSurgeCache creates a new RedisSurgeCache.
func NewRedisSurgeCache(client *redis.Client) *RedisSurgeCache {
	return &RedisSurgeCache{client: client}
}

// Compile-time interface check.
var _ port.SurgeCache = (*RedisSurgeCache)(nil)

// IncrementAreaDemand atomically increments the demand counter for an area.
// The TTL is set on the first increment and refreshed on subsequent calls.
func (s *RedisSurgeCache) IncrementAreaDemand(ctx context.Context, areaID string, ttl time.Duration) error {
	key := surgeKeyPrefix + areaID
	pipe := s.client.Pipeline()
	pipe.Incr(ctx, key)
	pipe.Expire(ctx, key, ttl)
	_, err := pipe.Exec(ctx)
	if err != nil {
		return fmt.Errorf("incrementing area demand for %s: %w", areaID, err)
	}
	return nil
}

// GetAreaDemand returns the current demand count for an area, or 0 if not set.
func (s *RedisSurgeCache) GetAreaDemand(ctx context.Context, areaID string) (int, error) {
	key := surgeKeyPrefix + areaID
	val, err := s.client.Get(ctx, key).Int()
	if err == redis.Nil {
		return 0, nil
	}
	if err != nil {
		return 0, fmt.Errorf("getting area demand for %s: %w", areaID, err)
	}
	return val, nil
}
