package cache

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"

	"github.com/davidlramirez95/towcommand/internal/usecase/port"
)

const (
	wsConnectionKeyPrefix = "ws:connection:"
	wsReverseKeyPrefix    = "ws:reverse:"
)

// RedisSessionCache implements port.SessionCache using Redis string keys with TTL.
type RedisSessionCache struct {
	client *redis.Client
}

// NewRedisSessionCache creates a new RedisSessionCache.
func NewRedisSessionCache(client *redis.Client) *RedisSessionCache {
	return &RedisSessionCache{client: client}
}

// Compile-time interface check.
var _ port.SessionCache = (*RedisSessionCache)(nil)

// MapConnection stores a WebSocket connection ID for a user with the given TTL.
func (s *RedisSessionCache) MapConnection(ctx context.Context, userID, connectionID string, ttl time.Duration) error {
	key := wsConnectionKeyPrefix + userID
	return s.client.Set(ctx, key, connectionID, ttl).Err()
}

// GetConnection returns the connection ID for a user, or empty string if not found.
func (s *RedisSessionCache) GetConnection(ctx context.Context, userID string) (string, error) {
	key := wsConnectionKeyPrefix + userID
	val, err := s.client.Get(ctx, key).Result()
	if err == redis.Nil {
		return "", nil
	}
	if err != nil {
		return "", fmt.Errorf("getting connection for user %s: %w", userID, err)
	}
	return val, nil
}

// RemoveConnection deletes the connection mapping for a user.
func (s *RedisSessionCache) RemoveConnection(ctx context.Context, userID string) error {
	key := wsConnectionKeyPrefix + userID
	return s.client.Del(ctx, key).Err()
}

// MapReverseConnection stores a reverse mapping from connectionID to userID with the given TTL.
func (s *RedisSessionCache) MapReverseConnection(ctx context.Context, connectionID, userID string, ttl time.Duration) error {
	key := wsReverseKeyPrefix + connectionID
	return s.client.Set(ctx, key, userID, ttl).Err()
}

// GetUserByConnection returns the user ID for a connection, or empty string if not found.
func (s *RedisSessionCache) GetUserByConnection(ctx context.Context, connectionID string) (string, error) {
	key := wsReverseKeyPrefix + connectionID
	val, err := s.client.Get(ctx, key).Result()
	if err == redis.Nil {
		return "", nil
	}
	if err != nil {
		return "", fmt.Errorf("getting user for connection %s: %w", connectionID, err)
	}
	return val, nil
}

// RemoveReverseConnection deletes the reverse connection mapping.
func (s *RedisSessionCache) RemoveReverseConnection(ctx context.Context, connectionID string) error {
	key := wsReverseKeyPrefix + connectionID
	return s.client.Del(ctx, key).Err()
}
