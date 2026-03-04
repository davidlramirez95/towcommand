package matching

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/davidlramirez95/towcommand/internal/domain/booking"
	domainerrors "github.com/davidlramirez95/towcommand/internal/domain/errors"
	"github.com/davidlramirez95/towcommand/internal/domain/provider"
	"github.com/davidlramirez95/towcommand/internal/domain/user"
	"github.com/davidlramirez95/towcommand/internal/usecase/port"
)

// --- Mocks ---

type mockBookingFinder struct {
	booking *booking.Booking
	err     error
}

func (m *mockBookingFinder) FindByID(_ context.Context, _ string) (*booking.Booking, error) {
	return m.booking, m.err
}

type mockProviderFinderMap struct {
	providers map[string]*provider.Provider
	err       error
}

func (m *mockProviderFinderMap) FindByID(_ context.Context, providerID string) (*provider.Provider, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.providers[providerID], nil
}

type mockGeoCache struct {
	// resultsByRadius maps radius -> results, allowing cascade simulation.
	resultsByRadius map[float64][]port.ProviderDistance
	err             error
}

func (m *mockGeoCache) FindNearbyProviders(_ context.Context, _, _, radiusKm float64) ([]port.ProviderDistance, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.resultsByRadius[radiusKm], nil
}

type mockSurgeCache struct {
	demand int
	err    error
}

func (m *mockSurgeCache) GetAreaDemand(_ context.Context, _ string) (int, error) {
	return m.demand, m.err
}

// --- Test fixtures ---

func testBooking() *booking.Booking {
	return &booking.Booking{
		BookingID:   "BK-001",
		CustomerID:  "USR-001",
		ServiceType: booking.ServiceTypeFlatbedTow,
		WeightClass: user.WeightClassLight,
		PickupLocation: booking.GeoLocation{
			Lat: 14.5995, Lng: 120.9842, Address: "Manila",
		},
	}
}

func testProviders() map[string]*provider.Provider {
	return map[string]*provider.Provider{
		"PROV-1": {
			ProviderID:          "PROV-1",
			TrustTier:           user.TrustTierSukiGold,
			AcceptanceRate:      0.95,
			TruckType:           provider.TruckTypeFlatbed,
			MaxWeightCapacityKg: 5000,
			IsOnline:            true,
		},
		"PROV-2": {
			ProviderID:          "PROV-2",
			TrustTier:           user.TrustTierVerified,
			AcceptanceRate:      0.80,
			TruckType:           provider.TruckTypeFlatbed,
			MaxWeightCapacityKg: 3000,
			IsOnline:            true,
		},
		"PROV-OFFLINE": {
			ProviderID:          "PROV-OFFLINE",
			TrustTier:           user.TrustTierSukiElite,
			AcceptanceRate:      0.99,
			TruckType:           provider.TruckTypeFlatbed,
			MaxWeightCapacityKg: 8000,
			IsOnline:            false,
		},
	}
}

// --- Tests ---

func TestMatchBookingUseCase_Execute(t *testing.T) {
	tests := []struct {
		name        string
		booking     *booking.Booking
		bookingErr  error
		providers   map[string]*provider.Provider
		providerErr error
		geoResults  map[float64][]port.ProviderDistance
		geoErr      error
		surgeDemand int
		surgeErr    error
		wantErr     bool
		wantErrCode domainerrors.ErrorCode
		checkResult func(t *testing.T, result *MatchResult)
	}{
		{
			name:      "happy path: finds and ranks providers at first radius",
			booking:   testBooking(),
			providers: testProviders(),
			geoResults: map[float64][]port.ProviderDistance{
				5: {
					{ProviderID: "PROV-1", DistanceKm: 3.0},
					{ProviderID: "PROV-2", DistanceKm: 4.5},
				},
			},
			surgeDemand: 2,
			checkResult: func(t *testing.T, result *MatchResult) {
				t.Helper()
				assert.False(t, result.SurgeMode)
				require.GreaterOrEqual(t, len(result.Scores), 1)
				// PROV-1 should rank higher: gold tier, higher acceptance, closer.
				assert.Equal(t, "PROV-1", result.Scores[0].ProviderID)
				// Scores should be descending.
				for i := 1; i < len(result.Scores); i++ {
					assert.GreaterOrEqual(t, result.Scores[i-1].TotalScore, result.Scores[i].TotalScore)
				}
			},
		},
		{
			name:      "cascade expansion: no providers at 5km, found at 10km",
			booking:   testBooking(),
			providers: testProviders(),
			geoResults: map[float64][]port.ProviderDistance{
				5:  {},
				10: {{ProviderID: "PROV-2", DistanceKm: 8.0}},
			},
			surgeDemand: 0,
			checkResult: func(t *testing.T, result *MatchResult) {
				t.Helper()
				assert.False(t, result.SurgeMode)
				require.Len(t, result.Scores, 1)
				assert.Equal(t, "PROV-2", result.Scores[0].ProviderID)
			},
		},
		{
			name:      "surge mode activation when demand >= 10",
			booking:   testBooking(),
			providers: testProviders(),
			geoResults: map[float64][]port.ProviderDistance{
				5: {{ProviderID: "PROV-1", DistanceKm: 2.0}},
			},
			surgeDemand: 15,
			checkResult: func(t *testing.T, result *MatchResult) {
				t.Helper()
				assert.True(t, result.SurgeMode)
				require.GreaterOrEqual(t, len(result.Scores), 1)
			},
		},
		{
			name:      "no providers found at any radius",
			booking:   testBooking(),
			providers: testProviders(),
			geoResults: map[float64][]port.ProviderDistance{
				5: {}, 10: {}, 20: {}, 30: {},
			},
			wantErr:     true,
			wantErrCode: domainerrors.CodeProviderUnavailable,
		},
		{
			name:        "booking not found",
			booking:     nil,
			wantErr:     true,
			wantErrCode: domainerrors.CodeNotFound,
		},
		{
			name:       "booking finder returns error",
			bookingErr: errors.New("db connection failed"),
			wantErr:    true,
		},
		{
			name:        "geo cache error returns external service error",
			booking:     testBooking(),
			geoErr:      errors.New("redis timeout"),
			wantErr:     true,
			wantErrCode: domainerrors.CodeExternalService,
		},
		{
			name:      "surge cache error is non-fatal and defaults to normal mode",
			booking:   testBooking(),
			providers: testProviders(),
			geoResults: map[float64][]port.ProviderDistance{
				5: {{ProviderID: "PROV-1", DistanceKm: 3.0}},
			},
			surgeErr: errors.New("redis down"),
			checkResult: func(t *testing.T, result *MatchResult) {
				t.Helper()
				assert.False(t, result.SurgeMode)
				require.GreaterOrEqual(t, len(result.Scores), 1)
			},
		},
		{
			name:      "offline providers filtered by RankProviders",
			booking:   testBooking(),
			providers: testProviders(),
			geoResults: map[float64][]port.ProviderDistance{
				5: {
					{ProviderID: "PROV-OFFLINE", DistanceKm: 1.0},
					{ProviderID: "PROV-1", DistanceKm: 3.0},
				},
			},
			checkResult: func(t *testing.T, result *MatchResult) {
				t.Helper()
				for _, s := range result.Scores {
					assert.NotEqual(t, "PROV-OFFLINE", s.ProviderID)
				}
			},
		},
		{
			name:      "surge mode at exact threshold boundary (demand == 10)",
			booking:   testBooking(),
			providers: testProviders(),
			geoResults: map[float64][]port.ProviderDistance{
				5: {{ProviderID: "PROV-1", DistanceKm: 2.0}},
			},
			surgeDemand: 10,
			checkResult: func(t *testing.T, result *MatchResult) {
				t.Helper()
				assert.True(t, result.SurgeMode)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bookingFinder := &mockBookingFinder{booking: tt.booking, err: tt.bookingErr}
			providerFinder := &mockProviderFinderMap{providers: tt.providers, err: tt.providerErr}
			geo := &mockGeoCache{resultsByRadius: tt.geoResults, err: tt.geoErr}
			surge := &mockSurgeCache{demand: tt.surgeDemand, err: tt.surgeErr}

			uc := NewMatchBookingUseCase(bookingFinder, providerFinder, geo, surge)
			result, err := uc.Execute(context.Background(), "BK-001")

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
			require.NotNil(t, result)
			if tt.checkResult != nil {
				tt.checkResult(t, result)
			}
		})
	}
}

func TestCascadeSearch_SkipsMissingProviders(t *testing.T) {
	// Provider in geo cache but not in provider finder should be skipped.
	geo := &mockGeoCache{
		resultsByRadius: map[float64][]port.ProviderDistance{
			5: {
				{ProviderID: "PROV-EXISTS", DistanceKm: 2.0},
				{ProviderID: "PROV-GHOST", DistanceKm: 3.0},
			},
		},
	}
	finder := &mockProviderFinderMap{
		providers: map[string]*provider.Provider{
			"PROV-EXISTS": {
				ProviderID:          "PROV-EXISTS",
				TrustTier:           user.TrustTierBasic,
				AcceptanceRate:      0.8,
				TruckType:           provider.TruckTypeFlatbed,
				MaxWeightCapacityKg: 5000,
				IsOnline:            true,
			},
		},
	}

	candidates, err := cascadeSearch(context.Background(), geo, finder, 14.5, 121.0)

	require.NoError(t, err)
	require.Len(t, candidates, 1)
	assert.Equal(t, "PROV-EXISTS", candidates[0].ProviderID)
	assert.Equal(t, 2.0, candidates[0].DistanceKm)
}

func TestCascadeSearch_ReturnsNilWhenAllRadiiEmpty(t *testing.T) {
	geo := &mockGeoCache{
		resultsByRadius: map[float64][]port.ProviderDistance{
			5: {}, 10: {}, 20: {}, 30: {},
		},
	}
	finder := &mockProviderFinderMap{providers: map[string]*provider.Provider{}}

	candidates, err := cascadeSearch(context.Background(), geo, finder, 14.5, 121.0)

	require.NoError(t, err)
	assert.Nil(t, candidates)
}

func TestCascadeSearch_StopsAtFirstRadiusWithResults(t *testing.T) {
	geo := &mockGeoCache{
		resultsByRadius: map[float64][]port.ProviderDistance{
			5:  {},
			10: {{ProviderID: "PROV-1", DistanceKm: 8.0}},
			20: {{ProviderID: "PROV-2", DistanceKm: 15.0}},
		},
	}
	finder := &mockProviderFinderMap{
		providers: map[string]*provider.Provider{
			"PROV-1": {
				ProviderID: "PROV-1", TrustTier: user.TrustTierBasic,
				AcceptanceRate: 0.7, TruckType: provider.TruckTypeFlatbed,
				MaxWeightCapacityKg: 5000, IsOnline: true,
			},
			"PROV-2": {
				ProviderID: "PROV-2", TrustTier: user.TrustTierBasic,
				AcceptanceRate: 0.7, TruckType: provider.TruckTypeFlatbed,
				MaxWeightCapacityKg: 5000, IsOnline: true,
			},
		},
	}

	candidates, err := cascadeSearch(context.Background(), geo, finder, 14.5, 121.0)

	require.NoError(t, err)
	// Should only return the 10km result, not the 20km one.
	require.Len(t, candidates, 1)
	assert.Equal(t, "PROV-1", candidates[0].ProviderID)
}
