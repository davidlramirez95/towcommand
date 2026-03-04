//go:build integration

package repository

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/davidlramirez95/towcommand/internal/domain/rating"
)

func TestRatingRepository(t *testing.T) {
	client := newTestClient(t)
	repo := NewRatingRepository(client, testTableName)
	ctx := context.Background()

	now := time.Now().UTC().Truncate(time.Second)

	r := &rating.Rating{
		BookingID:  "book-rate-001",
		CustomerID: "user-rate-001",
		ProviderID: "prov-rate-001",
		Score:      5,
		Comment:    "Excellent service, very professional!",
		Tags:       []string{"fast", "professional"},
		CreatedAt:  now,
	}

	t.Run("Save and FindByBooking", func(t *testing.T) {
		err := repo.Save(ctx, r)
		require.NoError(t, err)

		got, err := repo.FindByBooking(ctx, "book-rate-001")
		require.NoError(t, err)
		require.NotNil(t, got)
		assert.Equal(t, r.Score, got.Score)
		assert.Equal(t, r.Comment, got.Comment)
		assert.Equal(t, r.Tags, got.Tags)
	})

	t.Run("FindByBooking not found", func(t *testing.T) {
		got, err := repo.FindByBooking(ctx, "nonexistent")
		require.NoError(t, err)
		assert.Nil(t, got)
	})

	t.Run("FindByProvider", func(t *testing.T) {
		results, err := repo.FindByProvider(ctx, "prov-rate-001", 10)
		require.NoError(t, err)
		assert.Len(t, results, 1)
		assert.Equal(t, "book-rate-001", results[0].BookingID)
		assert.Equal(t, 5, results[0].Score)
	})

	t.Run("FindByProvider empty", func(t *testing.T) {
		results, err := repo.FindByProvider(ctx, "nonexistent", 10)
		require.NoError(t, err)
		assert.Empty(t, results)
	})

	t.Run("multiple ratings by provider", func(t *testing.T) {
		r2 := &rating.Rating{
			BookingID:  "book-rate-002",
			CustomerID: "user-rate-002",
			ProviderID: "prov-rate-001",
			Score:      4,
			Comment:    "Good service",
			CreatedAt:  now.Add(time.Hour),
		}
		err := repo.Save(ctx, r2)
		require.NoError(t, err)

		results, err := repo.FindByProvider(ctx, "prov-rate-001", 10)
		require.NoError(t, err)
		assert.Len(t, results, 2)
		// Most recent first (descending order).
		assert.Equal(t, "book-rate-002", results[0].BookingID)
		assert.Equal(t, "book-rate-001", results[1].BookingID)
	})
}
