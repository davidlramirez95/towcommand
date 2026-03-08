//go:build integration

package repository

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/davidlramirez95/towcommand/internal/usecase/port"
)

func TestPushRepository(t *testing.T) {
	client := newTestClient(t)
	repo := NewPushRepository(client, testTableName)
	ctx := context.Background()

	now := time.Now().UTC().Truncate(time.Second)

	token := &port.PushToken{
		UserID:      "user-push-001",
		Token:       "fcm-token-abc123",
		Platform:    port.PushPlatformFCM,
		DeviceID:    "device-001",
		EndpointArn: "arn:aws:sns:ap-southeast-1:123456789:endpoint/GCM/TowCommand/abc",
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	t.Run("Register and FindByUserID", func(t *testing.T) {
		err := repo.Register(ctx, token)
		require.NoError(t, err)

		tokens, err := repo.FindByUserID(ctx, "user-push-001")
		require.NoError(t, err)
		require.Len(t, tokens, 1)
		assert.Equal(t, "user-push-001", tokens[0].UserID)
		assert.Equal(t, "fcm-token-abc123", tokens[0].Token)
		assert.Equal(t, port.PushPlatformFCM, tokens[0].Platform)
		assert.Equal(t, "device-001", tokens[0].DeviceID)
	})

	t.Run("FindByUserID not found", func(t *testing.T) {
		tokens, err := repo.FindByUserID(ctx, "nonexistent")
		require.NoError(t, err)
		assert.Empty(t, tokens)
	})

	t.Run("multiple tokens per user", func(t *testing.T) {
		token2 := &port.PushToken{
			UserID:      "user-push-001",
			Token:       "apns-token-xyz789",
			Platform:    port.PushPlatformAPNS,
			DeviceID:    "device-002",
			EndpointArn: "arn:aws:sns:ap-southeast-1:123456789:endpoint/APNS/TowCommand/xyz",
			CreatedAt:   now.Add(time.Hour),
			UpdatedAt:   now.Add(time.Hour),
		}
		err := repo.Register(ctx, token2)
		require.NoError(t, err)

		tokens, err := repo.FindByUserID(ctx, "user-push-001")
		require.NoError(t, err)
		assert.Len(t, tokens, 2)
	})

	t.Run("Delete", func(t *testing.T) {
		err := repo.Delete(ctx, "user-push-001", "device-001")
		require.NoError(t, err)

		tokens, err := repo.FindByUserID(ctx, "user-push-001")
		require.NoError(t, err)
		assert.Len(t, tokens, 1)
		assert.Equal(t, "device-002", tokens[0].DeviceID)
	})

	t.Run("Delete nonexistent is no-op", func(t *testing.T) {
		err := repo.Delete(ctx, "user-push-001", "nonexistent-device")
		require.NoError(t, err)
	})
}
