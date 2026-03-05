package safetyuc

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/davidlramirez95/towcommand/internal/domain/booking"
	domainerrors "github.com/davidlramirez95/towcommand/internal/domain/errors"
	"github.com/davidlramirez95/towcommand/internal/domain/provider"
	"github.com/davidlramirez95/towcommand/internal/domain/safety"
	"github.com/davidlramirez95/towcommand/internal/domain/user"
	"github.com/davidlramirez95/towcommand/internal/usecase/port"
)

// --- Mocks ---

type mockBookingFinder struct{ mock.Mock }

func (m *mockBookingFinder) FindByID(ctx context.Context, bookingID string) (*booking.Booking, error) {
	args := m.Called(ctx, bookingID)
	if v := args.Get(0); v != nil {
		return v.(*booking.Booking), args.Error(1)
	}
	return nil, args.Error(1)
}

type mockProviderFinder struct{ mock.Mock }

func (m *mockProviderFinder) FindByID(ctx context.Context, providerID string) (*provider.Provider, error) {
	args := m.Called(ctx, providerID)
	if v := args.Get(0); v != nil {
		return v.(*provider.Provider), args.Error(1)
	}
	return nil, args.Error(1)
}

type mockSOSSaver struct{ mock.Mock }

func (m *mockSOSSaver) Save(ctx context.Context, alert *safety.SOSAlert) error {
	args := m.Called(ctx, alert)
	return args.Error(0)
}

type mockSOSFinder struct{ mock.Mock }

func (m *mockSOSFinder) FindByID(ctx context.Context, alertID string) (*safety.SOSAlert, error) {
	args := m.Called(ctx, alertID)
	if v := args.Get(0); v != nil {
		return v.(*safety.SOSAlert), args.Error(1)
	}
	return nil, args.Error(1)
}

type mockSOSResolver struct{ mock.Mock }

func (m *mockSOSResolver) Resolve(ctx context.Context, alertID string, resolvedBy string, resolvedAt time.Time) error {
	args := m.Called(ctx, alertID, resolvedBy, resolvedAt)
	return args.Error(0)
}

type mockSOSActiveLister struct{ mock.Mock }

func (m *mockSOSActiveLister) FindActive(ctx context.Context, limit int32) ([]safety.SOSAlert, error) {
	args := m.Called(ctx, limit)
	return args.Get(0).([]safety.SOSAlert), args.Error(1)
}

type mockEventPublisher struct{ mock.Mock }

func (m *mockEventPublisher) Publish(ctx context.Context, source, detailType string, detail any, actor *port.Actor) error {
	args := m.Called(ctx, source, detailType, detail, actor)
	return args.Error(0)
}

// --- Test fixtures ---

func testBookingWithProvider() *booking.Booking {
	return &booking.Booking{
		BookingID:  "BK-001",
		CustomerID: "USR-001",
		ProviderID: "PROV-001",
		PickupLocation: booking.GeoLocation{
			Lat: 14.5995, Lng: 120.9842,
		},
		DropoffLocation: booking.GeoLocation{
			Lat: 14.6500, Lng: 121.0500,
		},
	}
}

func testProvider() *provider.Provider {
	return &provider.Provider{
		ProviderID: "PROV-001",
		TrustTier:  user.TrustTierBasic,
	}
}

// daytimePHT returns a UTC time that is daytime in PHT (e.g. 14:00 PHT = 06:00 UTC).
func daytimePHT() time.Time {
	return time.Date(2026, 3, 5, 6, 0, 0, 0, time.UTC) // 14:00 PHT
}

// nighttimePHT returns a UTC time that is nighttime in PHT (e.g. 22:00 PHT = 14:00 UTC).
func nighttimePHT() time.Time {
	return time.Date(2026, 3, 5, 14, 0, 0, 0, time.UTC) // 22:00 PHT
}

// --- TriggerSOSUseCase Tests ---

func TestTriggerSOSUseCase_Execute(t *testing.T) {
	tests := []struct {
		name        string
		input       *TriggerSOSInput
		setupMocks  func(bf *mockBookingFinder, pf *mockProviderFinder, ss *mockSOSSaver, ep *mockEventPublisher)
		nowFunc     func() time.Time
		wantErr     bool
		wantErrCode domainerrors.ErrorCode
		checkResult func(t *testing.T, alert *safety.SOSAlert)
	}{
		{
			name: "success: daytime, basic provider, short distance",
			input: &TriggerSOSInput{
				BookingID:   "BK-001",
				TriggeredBy: "USR-001",
				TriggerType: safety.TriggerTypeButton,
				Lat:         14.5995,
				Lng:         120.9842,
			},
			nowFunc: daytimePHT,
			setupMocks: func(bf *mockBookingFinder, pf *mockProviderFinder, ss *mockSOSSaver, ep *mockEventPublisher) {
				bf.On("FindByID", mock.Anything, "BK-001").Return(testBookingWithProvider(), nil)
				pf.On("FindByID", mock.Anything, "PROV-001").Return(testProvider(), nil)
				ss.On("Save", mock.Anything, mock.AnythingOfType("*safety.SOSAlert")).Return(nil)
				ep.On("Publish", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)
			},
			checkResult: func(t *testing.T, alert *safety.SOSAlert) {
				t.Helper()
				assert.NotEmpty(t, alert.AlertID)
				assert.Equal(t, "BK-001", alert.BookingID)
				assert.Equal(t, "USR-001", alert.TriggeredBy)
				assert.Equal(t, safety.TriggerTypeButton, alert.TriggerType)
				assert.False(t, alert.Resolved)
				// basic tier = +15, daytime = 0, short distance = 0
				assert.Equal(t, 15, alert.Risk.Score)
				assert.Equal(t, "low", alert.Risk.Level)
				assert.Contains(t, alert.Risk.Factors, "basic_tier_provider")
			},
		},
		{
			name: "success: nighttime elevates risk",
			input: &TriggerSOSInput{
				BookingID:   "BK-001",
				TriggeredBy: "USR-001",
				TriggerType: safety.TriggerTypeTripleTap,
				Lat:         14.5995,
				Lng:         120.9842,
			},
			nowFunc: nighttimePHT,
			setupMocks: func(bf *mockBookingFinder, pf *mockProviderFinder, ss *mockSOSSaver, ep *mockEventPublisher) {
				bf.On("FindByID", mock.Anything, "BK-001").Return(testBookingWithProvider(), nil)
				pf.On("FindByID", mock.Anything, "PROV-001").Return(testProvider(), nil)
				ss.On("Save", mock.Anything, mock.AnythingOfType("*safety.SOSAlert")).Return(nil)
				ep.On("Publish", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)
			},
			checkResult: func(t *testing.T, alert *safety.SOSAlert) {
				t.Helper()
				// night = +25, basic = +15 => 40
				assert.Equal(t, 40, alert.Risk.Score)
				assert.Equal(t, "medium", alert.Risk.Level)
				assert.Contains(t, alert.Risk.Factors, "night_time")
				assert.Contains(t, alert.Risk.Factors, "basic_tier_provider")
			},
		},
		{
			name: "success: gold tier provider has lower risk",
			input: &TriggerSOSInput{
				BookingID:   "BK-001",
				TriggeredBy: "USR-001",
				TriggerType: safety.TriggerTypeButton,
				Lat:         14.5995,
				Lng:         120.9842,
			},
			nowFunc: daytimePHT,
			setupMocks: func(bf *mockBookingFinder, pf *mockProviderFinder, ss *mockSOSSaver, ep *mockEventPublisher) {
				bf.On("FindByID", mock.Anything, "BK-001").Return(testBookingWithProvider(), nil)
				goldProvider := &provider.Provider{
					ProviderID: "PROV-001",
					TrustTier:  user.TrustTierSukiGold,
				}
				pf.On("FindByID", mock.Anything, "PROV-001").Return(goldProvider, nil)
				ss.On("Save", mock.Anything, mock.AnythingOfType("*safety.SOSAlert")).Return(nil)
				ep.On("Publish", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)
			},
			checkResult: func(t *testing.T, alert *safety.SOSAlert) {
				t.Helper()
				// daytime, gold tier, short distance => score 0
				assert.Equal(t, 0, alert.Risk.Score)
				assert.Equal(t, "low", alert.Risk.Level)
			},
		},
		{
			name: "booking not found",
			input: &TriggerSOSInput{
				BookingID:   "BK-MISSING",
				TriggeredBy: "USR-001",
				TriggerType: safety.TriggerTypeButton,
				Lat:         14.5,
				Lng:         121.0,
			},
			nowFunc: daytimePHT,
			setupMocks: func(bf *mockBookingFinder, pf *mockProviderFinder, ss *mockSOSSaver, ep *mockEventPublisher) {
				bf.On("FindByID", mock.Anything, "BK-MISSING").Return(nil, nil)
			},
			wantErr:     true,
			wantErrCode: domainerrors.CodeNotFound,
		},
		{
			name: "booking finder error",
			input: &TriggerSOSInput{
				BookingID:   "BK-001",
				TriggeredBy: "USR-001",
				TriggerType: safety.TriggerTypeButton,
				Lat:         14.5,
				Lng:         121.0,
			},
			nowFunc: daytimePHT,
			setupMocks: func(bf *mockBookingFinder, pf *mockProviderFinder, ss *mockSOSSaver, ep *mockEventPublisher) {
				bf.On("FindByID", mock.Anything, "BK-001").Return(nil, errors.New("db error"))
			},
			wantErr: true,
		},
		{
			name: "provider finder error",
			input: &TriggerSOSInput{
				BookingID:   "BK-001",
				TriggeredBy: "USR-001",
				TriggerType: safety.TriggerTypeButton,
				Lat:         14.5,
				Lng:         121.0,
			},
			nowFunc: daytimePHT,
			setupMocks: func(bf *mockBookingFinder, pf *mockProviderFinder, ss *mockSOSSaver, ep *mockEventPublisher) {
				bf.On("FindByID", mock.Anything, "BK-001").Return(testBookingWithProvider(), nil)
				pf.On("FindByID", mock.Anything, "PROV-001").Return(nil, errors.New("db error"))
			},
			wantErr: true,
		},
		{
			name: "save error",
			input: &TriggerSOSInput{
				BookingID:   "BK-001",
				TriggeredBy: "USR-001",
				TriggerType: safety.TriggerTypeButton,
				Lat:         14.5,
				Lng:         121.0,
			},
			nowFunc: daytimePHT,
			setupMocks: func(bf *mockBookingFinder, pf *mockProviderFinder, ss *mockSOSSaver, ep *mockEventPublisher) {
				bf.On("FindByID", mock.Anything, "BK-001").Return(testBookingWithProvider(), nil)
				pf.On("FindByID", mock.Anything, "PROV-001").Return(testProvider(), nil)
				ss.On("Save", mock.Anything, mock.Anything).Return(errors.New("ddb error"))
			},
			wantErr: true,
		},
		{
			name: "validation error: missing bookingId",
			input: &TriggerSOSInput{
				TriggeredBy: "USR-001",
				TriggerType: safety.TriggerTypeButton,
			},
			nowFunc:     daytimePHT,
			setupMocks:  func(bf *mockBookingFinder, pf *mockProviderFinder, ss *mockSOSSaver, ep *mockEventPublisher) {},
			wantErr:     true,
			wantErrCode: domainerrors.CodeValidationError,
		},
		{
			name: "validation error: missing triggeredBy",
			input: &TriggerSOSInput{
				BookingID:   "BK-001",
				TriggerType: safety.TriggerTypeButton,
			},
			nowFunc:     daytimePHT,
			setupMocks:  func(bf *mockBookingFinder, pf *mockProviderFinder, ss *mockSOSSaver, ep *mockEventPublisher) {},
			wantErr:     true,
			wantErrCode: domainerrors.CodeValidationError,
		},
		{
			name: "validation error: missing triggerType",
			input: &TriggerSOSInput{
				BookingID:   "BK-001",
				TriggeredBy: "USR-001",
			},
			nowFunc:     daytimePHT,
			setupMocks:  func(bf *mockBookingFinder, pf *mockProviderFinder, ss *mockSOSSaver, ep *mockEventPublisher) {},
			wantErr:     true,
			wantErrCode: domainerrors.CodeValidationError,
		},
		{
			name: "booking without provider uses basic tier default",
			input: &TriggerSOSInput{
				BookingID:   "BK-NOPROV",
				TriggeredBy: "USR-001",
				TriggerType: safety.TriggerTypeShake,
				Lat:         14.5,
				Lng:         121.0,
			},
			nowFunc: daytimePHT,
			setupMocks: func(bf *mockBookingFinder, pf *mockProviderFinder, ss *mockSOSSaver, ep *mockEventPublisher) {
				bf.On("FindByID", mock.Anything, "BK-NOPROV").Return(&booking.Booking{
					BookingID: "BK-NOPROV",
					// No ProviderID
					DropoffLocation: booking.GeoLocation{Lat: 14.55, Lng: 121.05},
				}, nil)
				ss.On("Save", mock.Anything, mock.AnythingOfType("*safety.SOSAlert")).Return(nil)
				ep.On("Publish", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)
			},
			checkResult: func(t *testing.T, alert *safety.SOSAlert) {
				t.Helper()
				// basic tier default = +15
				assert.Equal(t, 15, alert.Risk.Score)
				assert.Contains(t, alert.Risk.Factors, "basic_tier_provider")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bf := new(mockBookingFinder)
			pf := new(mockProviderFinder)
			ss := new(mockSOSSaver)
			ep := new(mockEventPublisher)

			tt.setupMocks(bf, pf, ss, ep)

			uc := NewTriggerSOSUseCase(bf, pf, ss, ep)
			uc.now = tt.nowFunc
			uc.idGen = func() string { return "SOS-2026-test123" }

			alert, err := uc.Execute(context.Background(), tt.input)

			if tt.wantErr {
				require.Error(t, err)
				if tt.wantErrCode != "" {
					var appErr *domainerrors.AppError
					require.True(t, errors.As(err, &appErr), "expected AppError, got %T: %v", err, err)
					assert.Equal(t, tt.wantErrCode, appErr.Code)
				}
				return
			}

			require.NoError(t, err)
			require.NotNil(t, alert)
			if tt.checkResult != nil {
				tt.checkResult(t, alert)
			}
		})
	}
}

func TestIsNightTime(t *testing.T) {
	tests := []struct {
		name string
		utc  time.Time
		want bool
	}{
		{
			name: "22:00 PHT is night",
			utc:  time.Date(2026, 3, 5, 14, 0, 0, 0, time.UTC), // 22:00 PHT
			want: true,
		},
		{
			name: "04:00 PHT is night",
			utc:  time.Date(2026, 3, 5, 20, 0, 0, 0, time.UTC), // 04:00 PHT next day
			want: true,
		},
		{
			name: "05:00 PHT is daytime",
			utc:  time.Date(2026, 3, 5, 21, 0, 0, 0, time.UTC), // 05:00 PHT next day
			want: false,
		},
		{
			name: "14:00 PHT is daytime",
			utc:  time.Date(2026, 3, 5, 6, 0, 0, 0, time.UTC), // 14:00 PHT
			want: false,
		},
		{
			name: "20:00 PHT boundary is night",
			utc:  time.Date(2026, 3, 5, 12, 0, 0, 0, time.UTC), // 20:00 PHT
			want: true,
		},
		{
			name: "19:59 PHT is daytime",
			utc:  time.Date(2026, 3, 5, 11, 59, 0, 0, time.UTC), // 19:59 PHT
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, isNightTime(tt.utc))
		})
	}
}

func TestHaversine(t *testing.T) {
	// Manila to Makati: roughly 8-10 km
	dist := haversine(14.5995, 120.9842, 14.5547, 121.0244)
	assert.InDelta(t, 6.0, dist, 2.0, "Manila-Makati distance should be roughly 5-8 km")

	// Same point
	assert.Equal(t, 0.0, haversine(14.5, 121.0, 14.5, 121.0))
}
