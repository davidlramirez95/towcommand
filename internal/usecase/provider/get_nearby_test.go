package provider_test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	domainerrors "github.com/davidlramirez95/towcommand/internal/domain/errors"
	provdomain "github.com/davidlramirez95/towcommand/internal/domain/provider"
	"github.com/davidlramirez95/towcommand/internal/domain/user"
	provider "github.com/davidlramirez95/towcommand/internal/usecase/provider"
	"github.com/davidlramirez95/towcommand/internal/usecase/port"
)

func TestGetNearbyUseCase_Execute(t *testing.T) {
	onlineProvider := &provdomain.Provider{
		ProviderID:         "PROV-1",
		Name:               "Juan Towing",
		TruckType:          provdomain.TruckTypeFlatbed,
		Rating:             4.5,
		TotalJobsCompleted: 100,
		PlateNumber:        "ABC-1234",
		TrustTier:          user.TrustTierBasic,
		IsOnline:           true,
	}

	offlineProvider := &provdomain.Provider{
		ProviderID: "PROV-2",
		Name:       "Pedro Towing",
		TrustTier:  user.TrustTierBasic,
		IsOnline:   false,
	}

	tests := []struct {
		name      string
		input     provider.GetNearbyInput
		nearby    []port.ProviderDistance
		geoErr    error
		providers map[string]*provdomain.Provider
		findErr   error
		wantErr   bool
		errCode   domainerrors.ErrorCode
		checkOut  func(t *testing.T, out *provider.GetNearbyOutput)
	}{
		{
			name:  "successful nearby query",
			input: provider.GetNearbyInput{Lat: 14.5995, Lng: 120.9842, RadiusKm: 10, Limit: 10},
			nearby: []port.ProviderDistance{
				{ProviderID: "PROV-1", DistanceKm: 2.5},
			},
			providers: map[string]*provdomain.Provider{"PROV-1": onlineProvider},
			checkOut: func(t *testing.T, out *provider.GetNearbyOutput) {
				t.Helper()
				assert.Equal(t, 1, out.Count)
				assert.Equal(t, "PROV-1", out.Providers[0].ProviderID)
				assert.Equal(t, "Juan Towing", out.Providers[0].Name)
				assert.Equal(t, 2.5, out.Providers[0].DistanceKm)
				assert.Greater(t, out.Providers[0].ETAMinutes, 0)
			},
		},
		{
			name:  "filters out offline providers",
			input: provider.GetNearbyInput{Lat: 14.5995, Lng: 120.9842, RadiusKm: 10, Limit: 10},
			nearby: []port.ProviderDistance{
				{ProviderID: "PROV-1", DistanceKm: 2.5},
				{ProviderID: "PROV-2", DistanceKm: 3.0},
			},
			providers: map[string]*provdomain.Provider{
				"PROV-1": onlineProvider,
				"PROV-2": offlineProvider,
			},
			checkOut: func(t *testing.T, out *provider.GetNearbyOutput) {
				t.Helper()
				assert.Equal(t, 1, out.Count)
				assert.Equal(t, "PROV-1", out.Providers[0].ProviderID)
			},
		},
		{
			name: "invalid coordinates outside Philippines",
			input: provider.GetNearbyInput{
				Lat: 48.8566, Lng: 2.3522, RadiusKm: 10, Limit: 10,
			},
			wantErr: true,
			errCode: domainerrors.CodeValidationError,
		},
		{
			name:    "geo cache error",
			input:   provider.GetNearbyInput{Lat: 14.5995, Lng: 120.9842, RadiusKm: 10, Limit: 10},
			geoErr:  errors.New("redis down"),
			wantErr: true,
			errCode: domainerrors.CodeExternalService,
		},
		{
			name:   "empty results",
			input:  provider.GetNearbyInput{Lat: 14.5995, Lng: 120.9842, RadiusKm: 10, Limit: 10},
			nearby: []port.ProviderDistance{},
			checkOut: func(t *testing.T, out *provider.GetNearbyOutput) {
				t.Helper()
				assert.Equal(t, 0, out.Count)
				assert.Empty(t, out.Providers)
			},
		},
		{
			name:  "radius clamped to 50km max",
			input: provider.GetNearbyInput{Lat: 14.5995, Lng: 120.9842, RadiusKm: 100, Limit: 10},
			nearby: []port.ProviderDistance{
				{ProviderID: "PROV-1", DistanceKm: 2.5},
			},
			providers: map[string]*provdomain.Provider{"PROV-1": onlineProvider},
			checkOut: func(t *testing.T, out *provider.GetNearbyOutput) {
				t.Helper()
				assert.Equal(t, 1, out.Count)
			},
		},
		{
			name:  "default radius when 0",
			input: provider.GetNearbyInput{Lat: 14.5995, Lng: 120.9842, RadiusKm: 0, Limit: 10},
			nearby: []port.ProviderDistance{
				{ProviderID: "PROV-1", DistanceKm: 2.5},
			},
			providers: map[string]*provdomain.Provider{"PROV-1": onlineProvider},
			checkOut: func(t *testing.T, out *provider.GetNearbyOutput) {
				t.Helper()
				assert.Equal(t, 1, out.Count)
			},
		},
		{
			name:  "limit defaults to 20 when 0",
			input: provider.GetNearbyInput{Lat: 14.5995, Lng: 120.9842, RadiusKm: 10, Limit: 0},
			nearby: []port.ProviderDistance{
				{ProviderID: "PROV-1", DistanceKm: 2.5},
			},
			providers: map[string]*provdomain.Provider{"PROV-1": onlineProvider},
			checkOut: func(t *testing.T, out *provider.GetNearbyOutput) {
				t.Helper()
				assert.Equal(t, 1, out.Count)
			},
		},
		{
			name:  "skip providers not found in DB",
			input: provider.GetNearbyInput{Lat: 14.5995, Lng: 120.9842, RadiusKm: 10, Limit: 10},
			nearby: []port.ProviderDistance{
				{ProviderID: "PROV-1", DistanceKm: 2.5},
				{ProviderID: "PROV-GHOST", DistanceKm: 3.0},
			},
			providers: map[string]*provdomain.Provider{"PROV-1": onlineProvider},
			checkOut: func(t *testing.T, out *provider.GetNearbyOutput) {
				t.Helper()
				assert.Equal(t, 1, out.Count)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			geo := &mockGeoCache{nearby: tt.nearby, findErr: tt.geoErr}
			finder := &mockProviderFinderMap{providers: tt.providers, err: tt.findErr}

			uc := provider.NewGetNearbyUseCase(geo, finder)
			result, err := uc.Execute(context.Background(), tt.input)

			if tt.wantErr {
				require.Error(t, err)
				var appErr *domainerrors.AppError
				require.True(t, errors.As(err, &appErr))
				assert.Equal(t, tt.errCode, appErr.Code)
				return
			}

			require.NoError(t, err)
			require.NotNil(t, result)
			if tt.checkOut != nil {
				tt.checkOut(t, result)
			}
		})
	}
}
