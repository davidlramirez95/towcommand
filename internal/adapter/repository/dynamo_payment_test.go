//go:build integration

package repository

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/davidlramirez95/towcommand/internal/domain/payment"
)

func TestPaymentRepository(t *testing.T) {
	client := newTestClient(t)
	repo := NewPaymentRepository(client, testTableName)
	ctx := context.Background()

	now := time.Now().UTC().Truncate(time.Second)

	p := &payment.Payment{
		PaymentID: "pay-001",
		BookingID: "book-pay-001",
		UserID:    "user-pay-001",
		Amount:    65000,
		Currency:  "PHP",
		Method:    payment.PaymentMethodGCash,
		Status:    payment.PaymentStatusPending,
		CreatedAt: now,
		UpdatedAt: now,
	}

	t.Run("Save and FindByID", func(t *testing.T) {
		err := repo.Save(ctx, p)
		require.NoError(t, err)

		got, err := repo.FindByID(ctx, "pay-001")
		require.NoError(t, err)
		require.NotNil(t, got)
		assert.Equal(t, p.Amount, got.Amount)
		assert.Equal(t, p.Method, got.Method)
		assert.Equal(t, p.Status, got.Status)
	})

	t.Run("FindByID not found", func(t *testing.T) {
		got, err := repo.FindByID(ctx, "nonexistent")
		require.NoError(t, err)
		assert.Nil(t, got)
	})

	t.Run("FindByBooking", func(t *testing.T) {
		results, err := repo.FindByBooking(ctx, "book-pay-001")
		require.NoError(t, err)
		assert.Len(t, results, 1)
		assert.Equal(t, "pay-001", results[0].PaymentID)
	})

	t.Run("FindByBooking empty", func(t *testing.T) {
		results, err := repo.FindByBooking(ctx, "nonexistent")
		require.NoError(t, err)
		assert.Empty(t, results)
	})

	t.Run("UpdateStatus", func(t *testing.T) {
		err := repo.UpdateStatus(ctx, "pay-001", payment.PaymentStatusCaptured)
		require.NoError(t, err)

		got, err := repo.FindByID(ctx, "pay-001")
		require.NoError(t, err)
		require.NotNil(t, got)
		assert.Equal(t, payment.PaymentStatusCaptured, got.Status)
	})
}
