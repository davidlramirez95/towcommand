package cache

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"

	"github.com/davidlramirez95/towcommand/internal/usecase/port"
)

const otpKeyPrefix = "otp:"

// RedisOTPCache implements port.OTPCache using Redis with TTL-based auto-expiration.
type RedisOTPCache struct {
	client *redis.Client
}

// NewRedisOTPCache creates a new RedisOTPCache.
func NewRedisOTPCache(client *redis.Client) *RedisOTPCache {
	return &RedisOTPCache{client: client}
}

// Compile-time interface check.
var _ port.OTPCache = (*RedisOTPCache)(nil)

// StoreOTP stores a hashed OTP that auto-expires after the given TTL.
func (o *RedisOTPCache) StoreOTP(ctx context.Context, bookingID, otpType, hashedOTP string, ttl time.Duration) error {
	key := otpKeyPrefix + bookingID + ":" + otpType
	return o.client.Set(ctx, key, hashedOTP, ttl).Err()
}

// GetOTP retrieves a hashed OTP, returning empty string if not found or expired.
func (o *RedisOTPCache) GetOTP(ctx context.Context, bookingID, otpType string) (string, error) {
	key := otpKeyPrefix + bookingID + ":" + otpType
	val, err := o.client.Get(ctx, key).Result()
	if err == redis.Nil {
		return "", nil
	}
	if err != nil {
		return "", fmt.Errorf("getting OTP for booking %s: %w", bookingID, err)
	}
	return val, nil
}

// DeleteOTP removes an OTP entry immediately.
func (o *RedisOTPCache) DeleteOTP(ctx context.Context, bookingID, otpType string) error {
	key := otpKeyPrefix + bookingID + ":" + otpType
	return o.client.Del(ctx, key).Err()
}
