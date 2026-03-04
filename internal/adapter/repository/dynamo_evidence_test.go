//go:build integration

package repository

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/davidlramirez95/towcommand/internal/domain/evidence"
)

func TestEvidenceRepository(t *testing.T) {
	client := newTestClient(t)
	repo := NewEvidenceRepository(client, testTableName)
	ctx := context.Background()

	now := time.Now().UTC().Truncate(time.Second)

	cr := &evidence.ConditionReport{
		ReportID:   "rpt-001",
		BookingID:  "book-ev-001",
		ProviderID: "prov-ev-001",
		Phase:      "pickup",
		Media: []evidence.MediaItem{
			{
				MediaID:  "media-001",
				S3Key:    "evidence/book-ev-001/front.jpg",
				Position: evidence.PhotoPositionFront,
				MimeType: "image/jpeg",
				Integrity: evidence.HashIntegrity{
					Algorithm: "SHA-256",
					Hash:      "abc123",
				},
				CapturedAt: now,
			},
		},
		Notes:     "Minor scratch on front bumper",
		CreatedAt: now,
	}

	t.Run("Save and FindByBooking", func(t *testing.T) {
		err := repo.Save(ctx, cr)
		require.NoError(t, err)

		reports, err := repo.FindByBooking(ctx, "book-ev-001")
		require.NoError(t, err)
		assert.Len(t, reports, 1)
		assert.Equal(t, "rpt-001", reports[0].ReportID)
		assert.Equal(t, "pickup", reports[0].Phase)
		assert.Equal(t, "Minor scratch on front bumper", reports[0].Notes)
		assert.Len(t, reports[0].Media, 1)
	})

	t.Run("FindByBooking empty", func(t *testing.T) {
		reports, err := repo.FindByBooking(ctx, "nonexistent")
		require.NoError(t, err)
		assert.Empty(t, reports)
	})

	t.Run("AddMediaItem", func(t *testing.T) {
		mi := &evidence.MediaItem{
			MediaID:  "media-002",
			S3Key:    "evidence/book-ev-001/rear.jpg",
			Position: evidence.PhotoPositionRear,
			MimeType: "image/jpeg",
			Integrity: evidence.HashIntegrity{
				Algorithm: "SHA-256",
				Hash:      "def456",
			},
			CapturedAt: now,
		}
		err := repo.AddMediaItem(ctx, "book-ev-001", mi)
		require.NoError(t, err)

		// The media item is stored separately under MEDIA# SK, not under the report.
		// Verify via a direct GetItem.
		var got evidence.MediaItem
		found, err := repo.getItem(ctx, PrefixJob+"book-ev-001", PrefixMedia+"media-002", &got)
		require.NoError(t, err)
		assert.True(t, found)
		assert.Equal(t, "media-002", got.MediaID)
		assert.Equal(t, evidence.PhotoPositionRear, got.Position)
	})

	t.Run("multiple reports for same booking", func(t *testing.T) {
		cr2 := &evidence.ConditionReport{
			ReportID:   "rpt-002",
			BookingID:  "book-ev-001",
			ProviderID: "prov-ev-001",
			Phase:      "dropoff",
			CreatedAt:  now.Add(time.Hour),
		}
		err := repo.Save(ctx, cr2)
		require.NoError(t, err)

		reports, err := repo.FindByBooking(ctx, "book-ev-001")
		require.NoError(t, err)
		assert.Len(t, reports, 2)
	})
}
