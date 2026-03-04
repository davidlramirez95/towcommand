//go:build integration

package repository

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/davidlramirez95/towcommand/internal/domain/booking"
	"github.com/davidlramirez95/towcommand/internal/domain/user"
)

func TestBookingRepository(t *testing.T) {
	client := newTestClient(t)
	repo := NewBookingRepository(client, testTableName)
	ctx := context.Background()

	now := time.Now().UTC().Truncate(time.Second)

	b := &booking.Booking{
		BookingID:  "book-001",
		CustomerID: "user-001",
		ProviderID: "prov-001",
		VehicleID:  "veh-001",
		ServiceType: booking.ServiceTypeFlatbedTow,
		Status:      booking.BookingStatusPending,
		PickupLocation: booking.GeoLocation{
			Lat: 14.5995, Lng: 120.9842, Address: "Manila",
		},
		DropoffLocation: booking.GeoLocation{
			Lat: 14.6042, Lng: 120.9822, Address: "Quezon City",
		},
		WeightClass: user.WeightClassMedium,
		Price: booking.PriceBreakdown{
			Base: 50000, Distance: 15000, Total: 65000, Currency: "PHP",
		},
		EstimateID: "est-001",
		CreatedAt:  now,
		UpdatedAt:  now,
	}

	t.Run("Save and FindByID", func(t *testing.T) {
		err := repo.Save(ctx, b)
		require.NoError(t, err)

		got, err := repo.FindByID(ctx, "book-001")
		require.NoError(t, err)
		require.NotNil(t, got)
		assert.Equal(t, b.BookingID, got.BookingID)
		assert.Equal(t, b.CustomerID, got.CustomerID)
		assert.Equal(t, b.Status, got.Status)
		assert.Equal(t, b.Price.Total, got.Price.Total)
	})

	t.Run("FindByID not found", func(t *testing.T) {
		got, err := repo.FindByID(ctx, "nonexistent")
		require.NoError(t, err)
		assert.Nil(t, got)
	})

	t.Run("FindByUser", func(t *testing.T) {
		results, err := repo.FindByUser(ctx, "user-001", 10)
		require.NoError(t, err)
		assert.Len(t, results, 1)
		assert.Equal(t, "book-001", results[0].BookingID)
	})

	t.Run("FindByStatus", func(t *testing.T) {
		results, err := repo.FindByStatus(ctx, booking.BookingStatusPending, 10)
		require.NoError(t, err)
		assert.Len(t, results, 1)
		assert.Equal(t, "book-001", results[0].BookingID)
	})

	t.Run("UpdateStatus", func(t *testing.T) {
		err := repo.UpdateStatus(ctx, "book-001", booking.BookingStatusMatched, map[string]any{
			"reason": "provider accepted",
		})
		require.NoError(t, err)

		got, err := repo.FindByID(ctx, "book-001")
		require.NoError(t, err)
		require.NotNil(t, got)
		assert.Equal(t, booking.BookingStatusMatched, got.Status)
	})

	t.Run("FindByStatus after update", func(t *testing.T) {
		results, err := repo.FindByStatus(ctx, booking.BookingStatusPending, 10)
		require.NoError(t, err)
		assert.Empty(t, results)

		results, err = repo.FindByStatus(ctx, booking.BookingStatusMatched, 10)
		require.NoError(t, err)
		assert.Len(t, results, 1)
	})
}
