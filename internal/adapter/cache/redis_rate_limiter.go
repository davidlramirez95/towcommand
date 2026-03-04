package cache

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"

	"github.com/davidlramirez95/towcommand/internal/usecase/port"
)

const rateLimitKeyPrefix = "rate:"

// slidingWindowScript atomically performs a sliding-window rate-limit check.
// It removes expired entries, counts current requests, and conditionally adds a new one.
// Returns {1, remaining} if allowed, {0, 0} if denied.
var slidingWindowScript = redis.NewScript(`
local key = KEYS[1]
local now = tonumber(ARGV[1])
local window = tonumber(ARGV[2])
local limit = tonumber(ARGV[3])
local member = ARGV[4]

redis.call('ZREMRANGEBYSCORE', key, '-inf', now - window)
local count = redis.call('ZCARD', key)

if count < limit then
    redis.call('ZADD', key, now, member)
    redis.call('EXPIRE', key, window)
    return {1, limit - count - 1}
end

redis.call('EXPIRE', key, window)
return {0, 0}
`)

// RedisRateLimiter implements port.RateLimiter using a sorted-set sliding window.
type RedisRateLimiter struct {
	client *redis.Client
}

// NewRedisRateLimiter creates a new RedisRateLimiter.
func NewRedisRateLimiter(client *redis.Client) *RedisRateLimiter {
	return &RedisRateLimiter{client: client}
}

// Compile-time interface check.
var _ port.RateLimiter = (*RedisRateLimiter)(nil)

// CheckRateLimit checks whether a request is allowed under the sliding window.
// Each call with an allowed result records the request in the window.
func (r *RedisRateLimiter) CheckRateLimit(ctx context.Context, key string, maxRequests, windowSec int) (allowed bool, remaining int, err error) {
	redisKey := rateLimitKeyPrefix + key
	now := time.Now().UnixMicro()
	windowMicro := int64(windowSec) * 1_000_000
	member := strconv.FormatInt(now, 10)

	result, err := slidingWindowScript.Run(ctx, r.client, []string{redisKey}, now, windowMicro, maxRequests, member).Int64Slice()
	if err != nil {
		return false, 0, fmt.Errorf("checking rate limit for %s: %w", key, err)
	}

	allowed = result[0] == 1
	remaining = int(result[1])
	return allowed, remaining, nil
}
