//go:build integration

package repository

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/davidlramirez95/towcommand/internal/domain/provider"
	"github.com/davidlramirez95/towcommand/internal/domain/user"
)

func TestProviderRepository(t *testing.T) {
	client := newTestClient(t)
	repo := NewProviderRepository(client, testTableName)
	ctx := context.Background()

	now := time.Now().UTC().Truncate(time.Second)

	p := &provider.Provider{
		ProviderID:          "prov-100",
		CognitoSub:          "sub-prov-100",
		Name:                "Pedro's Towing",
		Phone:               "+639181234567",
		Email:               "pedro@towing.ph",
		Status:              provider.ProviderStatusActive,
		TrustTier:           user.TrustTierSukiGold,
		TruckType:           provider.TruckTypeFlatbed,
		MaxWeightCapacityKg: 3500,
		PlateNumber:         "XYZ 5678",
		LTORegistration:     "LTO-12345",
		NBIClearanceStatus:  provider.ClearanceStatusApproved,
		DrugTestStatus:      provider.ClearanceStatusApproved,
		MMADAccredited:      true,
		Rating:              4.8,
		TotalJobsCompleted:  150,
		AcceptanceRate:      0.92,
		IsOnline:            true,
		ServiceAreas:        []string{"NCR"},
		CreatedAt:           now,
		UpdatedAt:           now,
	}

	t.Run("Save and FindByID", func(t *testing.T) {
		err := repo.Save(ctx, p)
		require.NoError(t, err)

		got, err := repo.FindByID(ctx, "prov-100")
		require.NoError(t, err)
		require.NotNil(t, got)
		assert.Equal(t, p.Name, got.Name)
		assert.Equal(t, p.TrustTier, got.TrustTier)
		assert.InDelta(t, p.Rating, got.Rating, 0.01)
	})

	t.Run("FindByID not found", func(t *testing.T) {
		got, err := repo.FindByID(ctx, "nonexistent")
		require.NoError(t, err)
		assert.Nil(t, got)
	})

	t.Run("FindByTierAndCity", func(t *testing.T) {
		results, err := repo.FindByTierAndCity(ctx, user.TrustTierSukiGold, "NCR", 10)
		require.NoError(t, err)
		assert.Len(t, results, 1)
		assert.Equal(t, "prov-100", results[0].ProviderID)
	})

	t.Run("FindByTierAndCity empty", func(t *testing.T) {
		results, err := repo.FindByTierAndCity(ctx, user.TrustTierSukiElite, "NCR", 10)
		require.NoError(t, err)
		assert.Empty(t, results)
	})

	t.Run("UpdateLocation", func(t *testing.T) {
		err := repo.UpdateLocation(ctx, "prov-100", 14.6500, 121.0500)
		require.NoError(t, err)

		got, err := repo.FindByID(ctx, "prov-100")
		require.NoError(t, err)
		require.NotNil(t, got)
		require.NotNil(t, got.CurrentLat)
		require.NotNil(t, got.CurrentLng)
		assert.InDelta(t, 14.65, *got.CurrentLat, 0.001)
		assert.InDelta(t, 121.05, *got.CurrentLng, 0.001)
	})

	t.Run("UpdateAvailability", func(t *testing.T) {
		err := repo.UpdateAvailability(ctx, "prov-100", false)
		require.NoError(t, err)

		got, err := repo.FindByID(ctx, "prov-100")
		require.NoError(t, err)
		require.NotNil(t, got)
		assert.False(t, got.IsOnline)
	})

	t.Run("UploadDoc and GetDocs", func(t *testing.T) {
		doc := &provider.ProviderDoc{
			ProviderID: "prov-100",
			DocType:    provider.DocTypeNBIClearance,
			S3Key:      "docs/prov-100/nbi.pdf",
			Status:     provider.DocStatusApproved,
			UploadedAt: now,
		}
		err := repo.UploadDoc(ctx, doc)
		require.NoError(t, err)

		docs, err := repo.GetDocs(ctx, "prov-100")
		require.NoError(t, err)
		assert.Len(t, docs, 1)
		assert.Equal(t, provider.DocTypeNBIClearance, docs[0].DocType)
	})
}
