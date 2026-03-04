// Package cache provides Redis-backed cache adapters for the towcommand platform.
// Each adapter implements a port interface from internal/usecase/port.
package cache

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
)

// Options configures the Redis client connection.
type Options struct {
	Host     string
	Port     int
	Password string
	DB       int
	PoolSize int
}

// NewRedisClient creates a go-redis v9 client with connection pooling.
func NewRedisClient(opts Options) *redis.Client {
	poolSize := opts.PoolSize
	if poolSize == 0 {
		poolSize = 10
	}
	return redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", opts.Host, opts.Port),
		Password: opts.Password,
		DB:       opts.DB,
		PoolSize: poolSize,
	})
}

// HealthCheck pings Redis and returns an error if unreachable.
func HealthCheck(ctx context.Context, client *redis.Client) error {
	return client.Ping(ctx).Err()
}
