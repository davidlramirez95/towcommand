package paymentuc

import (
	"context"
	"fmt"
	"time"

	"github.com/davidlramirez95/towcommand/internal/domain/booking"
	domainerrors "github.com/davidlramirez95/towcommand/internal/domain/errors"
	"github.com/davidlramirez95/towcommand/internal/domain/payment"
	"github.com/davidlramirez95/towcommand/internal/domain/provider"
)

// EarningsPeriod aggregates earnings data for a specific time period.
// All monetary amounts are in centavos.
type EarningsPeriod struct {
	GrossAmount  int64 `json:"grossAmount"`
	Commission   int64 `json:"commission"`
	NetAmount    int64 `json:"netAmount"`
	BookingCount int   `json:"bookingCount"`
}

// EarningsOutput holds the complete earnings breakdown for a provider.
type EarningsOutput struct {
	ProviderID string         `json:"providerId"`
	Today      EarningsPeriod `json:"today"`
	ThisWeek   EarningsPeriod `json:"thisWeek"`
	ThisMonth  EarningsPeriod `json:"thisMonth"`
	AllTime    EarningsPeriod `json:"allTime"`
}

// BookingByProviderLister lists bookings assigned to a provider.
type BookingByProviderLister interface {
	FindByProvider(ctx context.Context, providerID string) ([]booking.Booking, error)
}

// PaymentByProviderBookingsLister retrieves payments across multiple bookings.
type PaymentByProviderBookingsLister interface {
	FindByProviderBookings(ctx context.Context, bookingIDs []string) ([]payment.Payment, error)
}

// GetProviderEarningsUseCase aggregates a provider's earnings from captured
// payments across all their completed bookings.
type GetProviderEarningsUseCase struct {
	providers ProviderFinder
	bookings  BookingByProviderLister
	payments  PaymentByProviderBookingsLister
	now       func() time.Time
}

// NewGetProviderEarningsUseCase constructs a GetProviderEarningsUseCase with its dependencies.
func NewGetProviderEarningsUseCase(
	providers ProviderFinder,
	bookings BookingByProviderLister,
	payments PaymentByProviderBookingsLister,
) *GetProviderEarningsUseCase {
	return &GetProviderEarningsUseCase{
		providers: providers,
		bookings:  bookings,
		payments:  payments,
		now:       func() time.Time { return time.Now().UTC() },
	}
}

// Execute retrieves and aggregates earnings for the given provider, broken down
// into today, this week, this month, and all-time periods.
func (uc *GetProviderEarningsUseCase) Execute(ctx context.Context, providerID string) (*EarningsOutput, error) {
	prov, err := uc.providers.FindByID(ctx, providerID)
	if err != nil {
		return nil, fmt.Errorf("finding provider %s: %w", providerID, err)
	}
	if prov == nil {
		return nil, domainerrors.NewNotFoundError("provider", providerID)
	}

	bookings, err := uc.bookings.FindByProvider(ctx, providerID)
	if err != nil {
		return nil, fmt.Errorf("listing bookings for provider %s: %w", providerID, err)
	}

	if len(bookings) == 0 {
		return &EarningsOutput{ProviderID: providerID}, nil
	}

	bookingIDs := make([]string, 0, len(bookings))
	for i := range bookings {
		bookingIDs = append(bookingIDs, bookings[i].BookingID)
	}

	payments, err := uc.payments.FindByProviderBookings(ctx, bookingIDs)
	if err != nil {
		return nil, fmt.Errorf("listing payments for provider %s bookings: %w", providerID, err)
	}

	return uc.aggregate(providerID, prov, payments), nil
}

// aggregate computes earnings per period from captured payments.
func (uc *GetProviderEarningsUseCase) aggregate(providerID string, prov *provider.Provider, payments []payment.Payment) *EarningsOutput {
	now := uc.now()
	todayStart := startOfDay(now)
	weekStart := startOfWeek(now)
	monthStart := startOfMonth(now)

	out := &EarningsOutput{ProviderID: providerID}

	for i := range payments {
		p := &payments[i]

		// Only count captured payments toward earnings.
		if p.Status != payment.PaymentStatusCaptured {
			continue
		}

		commission, _ := CalculateCommission(p.Amount, prov.TrustTier)
		net := p.Amount - commission

		// Determine the capture time; fall back to CreatedAt if CapturedAt is nil.
		capturedAt := p.CreatedAt
		if p.CapturedAt != nil {
			capturedAt = *p.CapturedAt
		}

		// All-time always counts.
		addToEarnings(&out.AllTime, p.Amount, commission, net)

		if !capturedAt.Before(monthStart) {
			addToEarnings(&out.ThisMonth, p.Amount, commission, net)
		}
		if !capturedAt.Before(weekStart) {
			addToEarnings(&out.ThisWeek, p.Amount, commission, net)
		}
		if !capturedAt.Before(todayStart) {
			addToEarnings(&out.Today, p.Amount, commission, net)
		}
	}

	return out
}

// addToEarnings accumulates amounts into an EarningsPeriod.
func addToEarnings(ep *EarningsPeriod, gross, commission, net int64) {
	ep.GrossAmount += gross
	ep.Commission += commission
	ep.NetAmount += net
	ep.BookingCount++
}

// startOfDay returns midnight UTC for the given time.
func startOfDay(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.UTC)
}

// startOfWeek returns the start of the ISO week (Monday 00:00 UTC).
func startOfWeek(t time.Time) time.Time {
	weekday := t.Weekday()
	if weekday == time.Sunday {
		weekday = 7
	}
	daysBack := int(weekday) - int(time.Monday)
	return startOfDay(t.AddDate(0, 0, -daysBack))
}

// startOfMonth returns the first day of the month at midnight UTC.
func startOfMonth(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), 1, 0, 0, 0, 0, time.UTC)
}
