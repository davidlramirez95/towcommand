package cache_test

import (
	"context"
	"testing"

	"github.com/alicebob/miniredis/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/davidlramirez95/towcommand/internal/adapter/cache"
)

func newRateLimiterTestClient(t *testing.T) (*cache.RedisRateLimiter, *miniredis.Miniredis) {
	t.Helper()
	mr := miniredis.RunT(t)
	client := cache.NewRedisClient(cache.Options{
		Host: mr.Host(),
		Port: mr.Server().Addr().Port,
	})
	t.Cleanup(func() { _ = client.Close() })
	return cache.NewRedisRateLimiter(client), mr
}

func TestRedisRateLimiter_AllowsWithinLimit(t *testing.T) {
	limiter, _ := newRateLimiterTestClient(t)
	ctx := context.Background()

	allowed, remaining, err := limiter.CheckRateLimit(ctx, "user-1", 5, 60)
	require.NoError(t, err)
	assert.True(t, allowed)
	assert.Equal(t, 4, remaining)
}

func TestRedisRateLimiter_DeniesOverLimit(t *testing.T) {
	limiter, _ := newRateLimiterTestClient(t)
	ctx := context.Background()

	maxReqs := 3
	for i := 0; i < maxReqs; i++ {
		allowed, _, err := limiter.CheckRateLimit(ctx, "user-flood", maxReqs, 60)
		require.NoError(t, err)
		assert.True(t, allowed)
	}

	// Next request should be denied.
	allowed, remaining, err := limiter.CheckRateLimit(ctx, "user-flood", maxReqs, 60)
	require.NoError(t, err)
	assert.False(t, allowed)
	assert.Equal(t, 0, remaining)
}

func TestRedisRateLimiter_RemainingDecreases(t *testing.T) {
	limiter, _ := newRateLimiterTestClient(t)
	ctx := context.Background()

	maxReqs := 5
	for i := 0; i < maxReqs; i++ {
		allowed, remaining, err := limiter.CheckRateLimit(ctx, "user-dec", maxReqs, 60)
		require.NoError(t, err)
		assert.True(t, allowed)
		assert.Equal(t, maxReqs-i-1, remaining)
	}
}

func TestRedisRateLimiter_SeparateKeys(t *testing.T) {
	limiter, _ := newRateLimiterTestClient(t)
	ctx := context.Background()

	// Exhaust limit for user-a.
	for i := 0; i < 2; i++ {
		_, _, err := limiter.CheckRateLimit(ctx, "user-a", 2, 60)
		require.NoError(t, err)
	}

	allowed, _, err := limiter.CheckRateLimit(ctx, "user-a", 2, 60)
	require.NoError(t, err)
	assert.False(t, allowed)

	// user-b should still be allowed.
	allowed, remaining, err := limiter.CheckRateLimit(ctx, "user-b", 2, 60)
	require.NoError(t, err)
	assert.True(t, allowed)
	assert.Equal(t, 1, remaining)
}
