package paymentuc

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/davidlramirez95/towcommand/internal/domain/booking"
	"github.com/davidlramirez95/towcommand/internal/domain/payment"
	"github.com/davidlramirez95/towcommand/internal/domain/provider"
	"github.com/davidlramirez95/towcommand/internal/domain/user"
)

// ---------------------------------------------------------------------------
// Mock implementations for earnings use case
// ---------------------------------------------------------------------------

type mockEarningsProviderFinder struct {
	FindByIDFunc func(ctx context.Context, providerID string) (*provider.Provider, error)
}

func (m *mockEarningsProviderFinder) FindByID(ctx context.Context, providerID string) (*provider.Provider, error) {
	if m.FindByIDFunc != nil {
		return m.FindByIDFunc(ctx, providerID)
	}
	return nil, nil
}

type mockEarningsBookingLister struct {
	FindByProviderFunc func(ctx context.Context, providerID string) ([]booking.Booking, error)
}

func (m *mockEarningsBookingLister) FindByProvider(ctx context.Context, providerID string) ([]booking.Booking, error) {
	if m.FindByProviderFunc != nil {
		return m.FindByProviderFunc(ctx, providerID)
	}
	return nil, nil
}

type mockEarningsPaymentLister struct {
	FindByProviderBookingsFunc func(ctx context.Context, bookingIDs []string) ([]payment.Payment, error)
}

func (m *mockEarningsPaymentLister) FindByProviderBookings(ctx context.Context, bookingIDs []string) ([]payment.Payment, error) {
	if m.FindByProviderBookingsFunc != nil {
		return m.FindByProviderBookingsFunc(ctx, bookingIDs)
	}
	return nil, nil
}

// ---------------------------------------------------------------------------
// Tests
// ---------------------------------------------------------------------------

func TestGetProviderEarningsUseCase(t *testing.T) {
	// Fixed reference time: Wednesday 2026-03-04 10:00:00 UTC
	refTime := time.Date(2026, 3, 4, 10, 0, 0, 0, time.UTC)

	testProvider := &provider.Provider{
		ProviderID: "prov-1",
		TrustTier:  user.TrustTierBasic, // 25% commission
	}

	testBookings := []booking.Booking{
		{BookingID: "bk-1", ProviderID: "prov-1", Status: booking.BookingStatusCompleted},
		{BookingID: "bk-2", ProviderID: "prov-1", Status: booking.BookingStatusCompleted},
		{BookingID: "bk-3", ProviderID: "prov-1", Status: booking.BookingStatusCompleted},
	}

	capturedToday := time.Date(2026, 3, 4, 8, 0, 0, 0, time.UTC)
	capturedThisWeek := time.Date(2026, 3, 2, 14, 0, 0, 0, time.UTC) // Monday
	capturedLastMonth := time.Date(2026, 2, 15, 12, 0, 0, 0, time.UTC)

	tests := []struct {
		name       string
		providerID string
		setup      func(pf *mockEarningsProviderFinder, bl *mockEarningsBookingLister, pl *mockEarningsPaymentLister)
		want       *EarningsOutput
		wantErr    bool
		errMsg     string
	}{
		{
			name:       "provider not found",
			providerID: "prov-missing",
			setup: func(pf *mockEarningsProviderFinder, bl *mockEarningsBookingLister, pl *mockEarningsPaymentLister) {
				pf.FindByIDFunc = func(_ context.Context, _ string) (*provider.Provider, error) {
					return nil, nil
				}
			},
			wantErr: true,
			errMsg:  "not found",
		},
		{
			name:       "provider finder error",
			providerID: "prov-1",
			setup: func(pf *mockEarningsProviderFinder, bl *mockEarningsBookingLister, pl *mockEarningsPaymentLister) {
				pf.FindByIDFunc = func(_ context.Context, _ string) (*provider.Provider, error) {
					return nil, errors.New("dynamo timeout")
				}
			},
			wantErr: true,
			errMsg:  "finding provider",
		},
		{
			name:       "zero bookings returns zero earnings",
			providerID: "prov-1",
			setup: func(pf *mockEarningsProviderFinder, bl *mockEarningsBookingLister, pl *mockEarningsPaymentLister) {
				pf.FindByIDFunc = func(_ context.Context, _ string) (*provider.Provider, error) {
					return testProvider, nil
				}
				bl.FindByProviderFunc = func(_ context.Context, _ string) ([]booking.Booking, error) {
					return nil, nil
				}
			},
			want: &EarningsOutput{ProviderID: "prov-1"},
		},
		{
			name:       "booking lister error",
			providerID: "prov-1",
			setup: func(pf *mockEarningsProviderFinder, bl *mockEarningsBookingLister, pl *mockEarningsPaymentLister) {
				pf.FindByIDFunc = func(_ context.Context, _ string) (*provider.Provider, error) {
					return testProvider, nil
				}
				bl.FindByProviderFunc = func(_ context.Context, _ string) ([]booking.Booking, error) {
					return nil, errors.New("scan failure")
				}
			},
			wantErr: true,
			errMsg:  "listing bookings",
		},
		{
			name:       "payment lister error",
			providerID: "prov-1",
			setup: func(pf *mockEarningsProviderFinder, bl *mockEarningsBookingLister, pl *mockEarningsPaymentLister) {
				pf.FindByIDFunc = func(_ context.Context, _ string) (*provider.Provider, error) {
					return testProvider, nil
				}
				bl.FindByProviderFunc = func(_ context.Context, _ string) ([]booking.Booking, error) {
					return testBookings, nil
				}
				pl.FindByProviderBookingsFunc = func(_ context.Context, _ []string) ([]payment.Payment, error) {
					return nil, errors.New("batch query failure")
				}
			},
			wantErr: true,
			errMsg:  "listing payments",
		},
		{
			name:       "only captured payments count",
			providerID: "prov-1",
			setup: func(pf *mockEarningsProviderFinder, bl *mockEarningsBookingLister, pl *mockEarningsPaymentLister) {
				pf.FindByIDFunc = func(_ context.Context, _ string) (*provider.Provider, error) {
					return testProvider, nil
				}
				bl.FindByProviderFunc = func(_ context.Context, _ string) ([]booking.Booking, error) {
					return testBookings[:1], nil
				}
				pl.FindByProviderBookingsFunc = func(_ context.Context, _ []string) ([]payment.Payment, error) {
					return []payment.Payment{
						{PaymentID: "pay-1", BookingID: "bk-1", Amount: 100000, Status: payment.PaymentStatusPending, CreatedAt: capturedToday},
						{PaymentID: "pay-2", BookingID: "bk-1", Amount: 200000, Status: payment.PaymentStatusRefunded, CreatedAt: capturedToday},
						{PaymentID: "pay-3", BookingID: "bk-1", Amount: 300000, Status: payment.PaymentStatusFailed, CreatedAt: capturedToday},
					}, nil
				}
			},
			want: &EarningsOutput{ProviderID: "prov-1"},
		},
		{
			name:       "aggregation across multiple periods",
			providerID: "prov-1",
			setup: func(pf *mockEarningsProviderFinder, bl *mockEarningsBookingLister, pl *mockEarningsPaymentLister) {
				pf.FindByIDFunc = func(_ context.Context, _ string) (*provider.Provider, error) {
					return testProvider, nil
				}
				bl.FindByProviderFunc = func(_ context.Context, _ string) ([]booking.Booking, error) {
					return testBookings, nil
				}
				pl.FindByProviderBookingsFunc = func(_ context.Context, ids []string) ([]payment.Payment, error) {
					assert.Len(t, ids, 3)
					return []payment.Payment{
						// bk-1: captured today (100000 centavos = 1000 PHP)
						{PaymentID: "pay-1", BookingID: "bk-1", Amount: 100000, Status: payment.PaymentStatusCaptured, CapturedAt: &capturedToday, CreatedAt: capturedToday},
						// bk-2: captured this week (200000 centavos)
						{PaymentID: "pay-2", BookingID: "bk-2", Amount: 200000, Status: payment.PaymentStatusCaptured, CapturedAt: &capturedThisWeek, CreatedAt: capturedThisWeek},
						// bk-3: captured last month (300000 centavos)
						{PaymentID: "pay-3", BookingID: "bk-3", Amount: 300000, Status: payment.PaymentStatusCaptured, CapturedAt: &capturedLastMonth, CreatedAt: capturedLastMonth},
					}, nil
				}
			},
			want: &EarningsOutput{
				ProviderID: "prov-1",
				// Today: only pay-1 (100000, 25% commission = 25000, net = 75000)
				Today: EarningsPeriod{
					GrossAmount:  100000,
					Commission:   25000,
					NetAmount:    75000,
					BookingCount: 1,
				},
				// This week: pay-1 + pay-2 (100000 + 200000 = 300000)
				ThisWeek: EarningsPeriod{
					GrossAmount:  300000,
					Commission:   75000,
					NetAmount:    225000,
					BookingCount: 2,
				},
				// This month: pay-1 + pay-2 (same as week since month is March)
				ThisMonth: EarningsPeriod{
					GrossAmount:  300000,
					Commission:   75000,
					NetAmount:    225000,
					BookingCount: 2,
				},
				// All-time: all three payments (600000)
				AllTime: EarningsPeriod{
					GrossAmount:  600000,
					Commission:   150000,
					NetAmount:    450000,
					BookingCount: 3,
				},
			},
		},
		{
			name:       "uses CapturedAt when available otherwise CreatedAt",
			providerID: "prov-1",
			setup: func(pf *mockEarningsProviderFinder, bl *mockEarningsBookingLister, pl *mockEarningsPaymentLister) {
				pf.FindByIDFunc = func(_ context.Context, _ string) (*provider.Provider, error) {
					return testProvider, nil
				}
				bl.FindByProviderFunc = func(_ context.Context, _ string) ([]booking.Booking, error) {
					return testBookings[:1], nil
				}
				pl.FindByProviderBookingsFunc = func(_ context.Context, _ []string) ([]payment.Payment, error) {
					// CapturedAt is nil; CreatedAt is today, so it counts for today.
					return []payment.Payment{
						{PaymentID: "pay-1", BookingID: "bk-1", Amount: 50000, Status: payment.PaymentStatusCaptured, CreatedAt: capturedToday},
					}, nil
				}
			},
			want: &EarningsOutput{
				ProviderID: "prov-1",
				Today:      EarningsPeriod{GrossAmount: 50000, Commission: 12500, NetAmount: 37500, BookingCount: 1},
				ThisWeek:   EarningsPeriod{GrossAmount: 50000, Commission: 12500, NetAmount: 37500, BookingCount: 1},
				ThisMonth:  EarningsPeriod{GrossAmount: 50000, Commission: 12500, NetAmount: 37500, BookingCount: 1},
				AllTime:    EarningsPeriod{GrossAmount: 50000, Commission: 12500, NetAmount: 37500, BookingCount: 1},
			},
		},
		{
			name:       "verified tier lower commission rate",
			providerID: "prov-1",
			setup: func(pf *mockEarningsProviderFinder, bl *mockEarningsBookingLister, pl *mockEarningsPaymentLister) {
				verifiedProvider := &provider.Provider{
					ProviderID: "prov-1",
					TrustTier:  user.TrustTierSukiGold, // 18% commission
				}
				pf.FindByIDFunc = func(_ context.Context, _ string) (*provider.Provider, error) {
					return verifiedProvider, nil
				}
				bl.FindByProviderFunc = func(_ context.Context, _ string) ([]booking.Booking, error) {
					return testBookings[:1], nil
				}
				pl.FindByProviderBookingsFunc = func(_ context.Context, _ []string) ([]payment.Payment, error) {
					return []payment.Payment{
						{PaymentID: "pay-1", BookingID: "bk-1", Amount: 100000, Status: payment.PaymentStatusCaptured, CapturedAt: &capturedToday, CreatedAt: capturedToday},
					}, nil
				}
			},
			want: &EarningsOutput{
				ProviderID: "prov-1",
				// 18% of 100000 = 18000, net = 82000
				Today:     EarningsPeriod{GrossAmount: 100000, Commission: 18000, NetAmount: 82000, BookingCount: 1},
				ThisWeek:  EarningsPeriod{GrossAmount: 100000, Commission: 18000, NetAmount: 82000, BookingCount: 1},
				ThisMonth: EarningsPeriod{GrossAmount: 100000, Commission: 18000, NetAmount: 82000, BookingCount: 1},
				AllTime:   EarningsPeriod{GrossAmount: 100000, Commission: 18000, NetAmount: 82000, BookingCount: 1},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pf := &mockEarningsProviderFinder{}
			bl := &mockEarningsBookingLister{}
			pl := &mockEarningsPaymentLister{}

			if tt.setup != nil {
				tt.setup(pf, bl, pl)
			}

			uc := NewGetProviderEarningsUseCase(pf, bl, pl)
			uc.now = func() time.Time { return refTime }

			result, err := uc.Execute(context.Background(), tt.providerID)

			if tt.wantErr {
				require.Error(t, err)
				if tt.errMsg != "" {
					assert.Contains(t, err.Error(), tt.errMsg)
				}
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tt.want, result)
		})
	}
}

func TestStartOfWeek(t *testing.T) {
	tests := []struct {
		name string
		in   time.Time
		want time.Time
	}{
		{
			name: "wednesday returns monday",
			in:   time.Date(2026, 3, 4, 15, 30, 0, 0, time.UTC),
			want: time.Date(2026, 3, 2, 0, 0, 0, 0, time.UTC),
		},
		{
			name: "monday returns same monday",
			in:   time.Date(2026, 3, 2, 10, 0, 0, 0, time.UTC),
			want: time.Date(2026, 3, 2, 0, 0, 0, 0, time.UTC),
		},
		{
			name: "sunday returns previous monday",
			in:   time.Date(2026, 3, 8, 23, 59, 0, 0, time.UTC),
			want: time.Date(2026, 3, 2, 0, 0, 0, 0, time.UTC),
		},
		{
			name: "saturday returns monday of same week",
			in:   time.Date(2026, 3, 7, 12, 0, 0, 0, time.UTC),
			want: time.Date(2026, 3, 2, 0, 0, 0, 0, time.UTC),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := startOfWeek(tt.in)
			assert.Equal(t, tt.want, got)
		})
	}
}
