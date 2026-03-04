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

func newOTPTestClient(t *testing.T) (*cache.RedisOTPCache, *miniredis.Miniredis) {
	t.Helper()
	mr := miniredis.RunT(t)
	client := cache.NewRedisClient(cache.Options{
		Host: mr.Host(),
		Port: mr.Server().Addr().Port,
	})
	t.Cleanup(func() { _ = client.Close() })
	return cache.NewRedisOTPCache(client), mr
}

func TestRedisOTPCache_StoreAndGetOTP(t *testing.T) {
	otpCache, _ := newOTPTestClient(t)
	ctx := context.Background()

	tests := []struct {
		name      string
		bookingID string
		otpType   string
		hashedOTP string
		ttl       time.Duration
	}{
		{
			name:      "pickup OTP",
			bookingID: "booking-1",
			otpType:   "PICKUP",
			hashedOTP: "sha256-hash-pickup",
			ttl:       5 * time.Minute,
		},
		{
			name:      "dropoff OTP",
			bookingID: "booking-1",
			otpType:   "DROPOFF",
			hashedOTP: "sha256-hash-dropoff",
			ttl:       5 * time.Minute,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := otpCache.StoreOTP(ctx, tt.bookingID, tt.otpType, tt.hashedOTP, tt.ttl)
			require.NoError(t, err)

			got, err := otpCache.GetOTP(ctx, tt.bookingID, tt.otpType)
			require.NoError(t, err)
			assert.Equal(t, tt.hashedOTP, got)
		})
	}
}

func TestRedisOTPCache_GetOTP_NotFound(t *testing.T) {
	otpCache, _ := newOTPTestClient(t)
	ctx := context.Background()

	got, err := otpCache.GetOTP(ctx, "nonexistent", "PICKUP")
	require.NoError(t, err)
	assert.Empty(t, got)
}

func TestRedisOTPCache_DeleteOTP(t *testing.T) {
	otpCache, _ := newOTPTestClient(t)
	ctx := context.Background()

	require.NoError(t, otpCache.StoreOTP(ctx, "booking-del", "PICKUP", "hash123", 5*time.Minute))

	err := otpCache.DeleteOTP(ctx, "booking-del", "PICKUP")
	require.NoError(t, err)

	got, err := otpCache.GetOTP(ctx, "booking-del", "PICKUP")
	require.NoError(t, err)
	assert.Empty(t, got)
}

func TestRedisOTPCache_AutoExpiry(t *testing.T) {
	otpCache, mr := newOTPTestClient(t)
	ctx := context.Background()

	require.NoError(t, otpCache.StoreOTP(ctx, "booking-exp", "PICKUP", "hash-exp", 1*time.Second))

	// Verify exists before expiry.
	got, err := otpCache.GetOTP(ctx, "booking-exp", "PICKUP")
	require.NoError(t, err)
	assert.Equal(t, "hash-exp", got)

	// Fast-forward past TTL.
	mr.FastForward(2 * time.Second)

	got, err = otpCache.GetOTP(ctx, "booking-exp", "PICKUP")
	require.NoError(t, err)
	assert.Empty(t, got)
}

func TestRedisOTPCache_SeparateKeysByType(t *testing.T) {
	otpCache, _ := newOTPTestClient(t)
	ctx := context.Background()

	require.NoError(t, otpCache.StoreOTP(ctx, "booking-sep", "PICKUP", "pickup-hash", 5*time.Minute))
	require.NoError(t, otpCache.StoreOTP(ctx, "booking-sep", "DROPOFF", "dropoff-hash", 5*time.Minute))

	pickupOTP, err := otpCache.GetOTP(ctx, "booking-sep", "PICKUP")
	require.NoError(t, err)
	assert.Equal(t, "pickup-hash", pickupOTP)

	dropoffOTP, err := otpCache.GetOTP(ctx, "booking-sep", "DROPOFF")
	require.NoError(t, err)
	assert.Equal(t, "dropoff-hash", dropoffOTP)

	// Deleting one type should not affect the other.
	require.NoError(t, otpCache.DeleteOTP(ctx, "booking-sep", "PICKUP"))

	pickupOTP, err = otpCache.GetOTP(ctx, "booking-sep", "PICKUP")
	require.NoError(t, err)
	assert.Empty(t, pickupOTP)

	dropoffOTP, err = otpCache.GetOTP(ctx, "booking-sep", "DROPOFF")
	require.NoError(t, err)
	assert.Equal(t, "dropoff-hash", dropoffOTP)
}
