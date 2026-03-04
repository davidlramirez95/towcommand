//go:build integration

package repository

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/davidlramirez95/towcommand/internal/domain/user"
)

func TestUserRepository(t *testing.T) {
	client := newTestClient(t)
	repo := NewUserRepository(client, testTableName)
	ctx := context.Background()

	now := time.Now().UTC().Truncate(time.Second)

	u := &user.User{
		UserID:     "user-100",
		CognitoSub: "sub-100",
		Email:      "juan@example.com",
		Phone:      "+639171234567",
		Name:       "Juan Dela Cruz",
		UserType:   user.UserTypeCustomer,
		TrustTier:  user.TrustTierBasic,
		Language:   user.LanguageFilipino,
		Status:     user.UserStatusActive,
		CreatedAt:  now,
		UpdatedAt:  now,
	}

	t.Run("Save and FindByID", func(t *testing.T) {
		err := repo.Save(ctx, u)
		require.NoError(t, err)

		got, err := repo.FindByID(ctx, "user-100")
		require.NoError(t, err)
		require.NotNil(t, got)
		assert.Equal(t, u.Email, got.Email)
		assert.Equal(t, u.Phone, got.Phone)
		assert.Equal(t, u.Name, got.Name)
	})

	t.Run("FindByID not found", func(t *testing.T) {
		got, err := repo.FindByID(ctx, "nonexistent")
		require.NoError(t, err)
		assert.Nil(t, got)
	})

	t.Run("FindByEmail", func(t *testing.T) {
		got, err := repo.FindByEmail(ctx, "juan@example.com")
		require.NoError(t, err)
		require.NotNil(t, got)
		assert.Equal(t, "user-100", got.UserID)
	})

	t.Run("FindByEmail not found", func(t *testing.T) {
		got, err := repo.FindByEmail(ctx, "nobody@example.com")
		require.NoError(t, err)
		assert.Nil(t, got)
	})

	t.Run("FindByPhone", func(t *testing.T) {
		got, err := repo.FindByPhone(ctx, "+639171234567")
		require.NoError(t, err)
		require.NotNil(t, got)
		assert.Equal(t, "user-100", got.UserID)
	})

	t.Run("FindByPhone not found", func(t *testing.T) {
		got, err := repo.FindByPhone(ctx, "+639170000000")
		require.NoError(t, err)
		assert.Nil(t, got)
	})

	t.Run("AddVehicle and GetVehicles", func(t *testing.T) {
		v := &user.UserVehicle{
			VehicleID:   "veh-100",
			UserID:      "user-100",
			Make:        "Toyota",
			Model:       "Vios",
			Year:        2022,
			PlateNumber: "ABC 1234",
			WeightClass: user.WeightClassLight,
			Color:       "White",
			IsDefault:   true,
			CreatedAt:   now,
		}
		err := repo.AddVehicle(ctx, v)
		require.NoError(t, err)

		vehicles, err := repo.GetVehicles(ctx, "user-100")
		require.NoError(t, err)
		assert.Len(t, vehicles, 1)
		assert.Equal(t, "veh-100", vehicles[0].VehicleID)
		assert.Equal(t, "Toyota", vehicles[0].Make)
	})

	t.Run("GetVehicles empty", func(t *testing.T) {
		vehicles, err := repo.GetVehicles(ctx, "user-no-vehicles")
		require.NoError(t, err)
		assert.Empty(t, vehicles)
	})
}
