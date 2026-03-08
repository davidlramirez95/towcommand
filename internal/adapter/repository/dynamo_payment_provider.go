package repository

import (
	"context"
	"fmt"

	"github.com/davidlramirez95/towcommand/internal/domain/payment"
)

// FindByProviderBookings retrieves payments for multiple bookings by querying
// each booking's payments via GSI1. This supports the provider earnings use case
// where we first look up a provider's bookings, then batch-fetch their payments.
func (r *DynamoPaymentRepository) FindByProviderBookings(ctx context.Context, bookingIDs []string) ([]payment.Payment, error) {
	if len(bookingIDs) == 0 {
		return nil, nil
	}

	result := make([]payment.Payment, 0, len(bookingIDs))
	for _, bid := range bookingIDs {
		payments, err := r.FindByBooking(ctx, bid)
		if err != nil {
			return nil, fmt.Errorf("finding payments for booking %s: %w", bid, err)
		}
		result = append(result, payments...)
	}
	return result, nil
}
