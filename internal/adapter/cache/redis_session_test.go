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

func newSessionTestClient(t *testing.T) (*cache.RedisSessionCache, *miniredis.Miniredis) {
	t.Helper()
	mr := miniredis.RunT(t)
	client := cache.NewRedisClient(cache.Options{
		Host: mr.Host(),
		Port: mr.Server().Addr().Port,
	})
	t.Cleanup(func() { client.Close() })
	return cache.NewRedisSessionCache(client), mr
}

func TestRedisSessionCache_MapAndGetConnection(t *testing.T) {
	session, _ := newSessionTestClient(t)
	ctx := context.Background()

	tests := []struct {
		name         string
		userID       string
		connectionID string
		ttl          time.Duration
	}{
		{
			name:         "map user connection",
			userID:       "user-1",
			connectionID: "conn-abc123",
			ttl:          10 * time.Minute,
		},
		{
			name:         "map another user connection",
			userID:       "user-2",
			connectionID: "conn-def456",
			ttl:          5 * time.Minute,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := session.MapConnection(ctx, tt.userID, tt.connectionID, tt.ttl)
			require.NoError(t, err)

			got, err := session.GetConnection(ctx, tt.userID)
			require.NoError(t, err)
			assert.Equal(t, tt.connectionID, got)
		})
	}
}

func TestRedisSessionCache_GetConnection_NotFound(t *testing.T) {
	session, _ := newSessionTestClient(t)
	ctx := context.Background()

	got, err := session.GetConnection(ctx, "nonexistent-user")
	require.NoError(t, err)
	assert.Empty(t, got)
}

func TestRedisSessionCache_RemoveConnection(t *testing.T) {
	session, _ := newSessionTestClient(t)
	ctx := context.Background()

	require.NoError(t, session.MapConnection(ctx, "user-rm", "conn-xyz", 10*time.Minute))

	err := session.RemoveConnection(ctx, "user-rm")
	require.NoError(t, err)

	got, err := session.GetConnection(ctx, "user-rm")
	require.NoError(t, err)
	assert.Empty(t, got)
}

func TestRedisSessionCache_TTLExpiry(t *testing.T) {
	session, mr := newSessionTestClient(t)
	ctx := context.Background()

	require.NoError(t, session.MapConnection(ctx, "user-ttl", "conn-ttl", 1*time.Second))

	// Verify exists before expiry.
	got, err := session.GetConnection(ctx, "user-ttl")
	require.NoError(t, err)
	assert.Equal(t, "conn-ttl", got)

	// Fast-forward time in miniredis.
	mr.FastForward(2 * time.Second)

	got, err = session.GetConnection(ctx, "user-ttl")
	require.NoError(t, err)
	assert.Empty(t, got)
}
