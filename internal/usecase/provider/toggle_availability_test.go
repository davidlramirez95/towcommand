package provider_test

import (
	"context"
	"errors"
	"log/slog"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	domainerrors "github.com/davidlramirez95/towcommand/internal/domain/errors"
	provdomain "github.com/davidlramirez95/towcommand/internal/domain/provider"
	"github.com/davidlramirez95/towcommand/internal/domain/user"
	provider "github.com/davidlramirez95/towcommand/internal/usecase/provider"
)

func TestToggleAvailabilityUseCase_Execute(t *testing.T) {
	lat := 14.5995
	lng := 120.9842

	activeProvider := &provdomain.Provider{
		ProviderID: "PROV-123",
		Status:     provdomain.ProviderStatusActive,
		TrustTier:  user.TrustTierBasic,
		CurrentLat: &lat,
		CurrentLng: &lng,
		IsOnline:   false,
	}

	pendingProvider := &provdomain.Provider{
		ProviderID: "PROV-456",
		Status:     provdomain.ProviderStatusPendingVerification,
		TrustTier:  user.TrustTierBasic,
		IsOnline:   false,
	}

	activeNoLocation := &provdomain.Provider{
		ProviderID: "PROV-789",
		Status:     provdomain.ProviderStatusActive,
		TrustTier:  user.TrustTierBasic,
		IsOnline:   false,
	}

	tests := []struct {
		name      string
		input     provider.ToggleAvailabilityInput
		found     *provdomain.Provider
		findErr   error
		updateErr error
		geoAddErr error
		geoRemErr error
		pubErr    error
		wantErr   bool
		errCode   domainerrors.ErrorCode
		checkOut  func(t *testing.T, out *provider.ToggleAvailabilityOutput)
		checkGeo  func(t *testing.T, geo *mockGeoCache)
	}{
		{
			name:  "go online successfully",
			input: provider.ToggleAvailabilityInput{ProviderID: "PROV-123", Online: true},
			found: activeProvider,
			checkOut: func(t *testing.T, out *provider.ToggleAvailabilityOutput) {
				t.Helper()
				assert.Equal(t, "PROV-123", out.ProviderID)
				assert.True(t, out.Online)
			},
			checkGeo: func(t *testing.T, geo *mockGeoCache) {
				t.Helper()
				assert.True(t, geo.addCalled, "expected geo add to be called")
			},
		},
		{
			name:  "go offline successfully",
			input: provider.ToggleAvailabilityInput{ProviderID: "PROV-123", Online: false},
			found: activeProvider,
			checkOut: func(t *testing.T, out *provider.ToggleAvailabilityOutput) {
				t.Helper()
				assert.Equal(t, "PROV-123", out.ProviderID)
				assert.False(t, out.Online)
			},
			checkGeo: func(t *testing.T, geo *mockGeoCache) {
				t.Helper()
				assert.True(t, geo.removeCalled, "expected geo remove to be called")
			},
		},
		{
			name:    "provider not found",
			input:   provider.ToggleAvailabilityInput{ProviderID: "PROV-999", Online: true},
			found:   nil,
			wantErr: true,
			errCode: domainerrors.CodeNotFound,
		},
		{
			name:    "pending provider cannot go online",
			input:   provider.ToggleAvailabilityInput{ProviderID: "PROV-456", Online: true},
			found:   pendingProvider,
			wantErr: true,
			errCode: domainerrors.CodeValidationError,
		},
		{
			name:    "find provider fails",
			input:   provider.ToggleAvailabilityInput{ProviderID: "PROV-123", Online: true},
			findErr: errors.New("dynamo timeout"),
			wantErr: true,
			errCode: domainerrors.CodeInternalError,
		},
		{
			name:      "update availability fails",
			input:     provider.ToggleAvailabilityInput{ProviderID: "PROV-123", Online: true},
			found:     activeProvider,
			updateErr: errors.New("dynamo timeout"),
			wantErr:   true,
			errCode:   domainerrors.CodeInternalError,
		},
		{
			name:  "go online without location skips geo add",
			input: provider.ToggleAvailabilityInput{ProviderID: "PROV-789", Online: true},
			found: activeNoLocation,
			checkGeo: func(t *testing.T, geo *mockGeoCache) {
				t.Helper()
				assert.False(t, geo.addCalled, "expected geo add to NOT be called")
			},
		},
		{
			name:      "geo add error does not fail toggle",
			input:     provider.ToggleAvailabilityInput{ProviderID: "PROV-123", Online: true},
			found:     activeProvider,
			geoAddErr: errors.New("redis down"),
			checkOut: func(t *testing.T, out *provider.ToggleAvailabilityOutput) {
				t.Helper()
				assert.True(t, out.Online)
			},
		},
		{
			name:      "geo remove error does not fail toggle",
			input:     provider.ToggleAvailabilityInput{ProviderID: "PROV-123", Online: false},
			found:     activeProvider,
			geoRemErr: errors.New("redis down"),
			checkOut: func(t *testing.T, out *provider.ToggleAvailabilityOutput) {
				t.Helper()
				assert.False(t, out.Online)
			},
		},
		{
			name:   "publish error does not fail toggle",
			input:  provider.ToggleAvailabilityInput{ProviderID: "PROV-123", Online: true},
			found:  activeProvider,
			pubErr: errors.New("eventbridge down"),
			checkOut: func(t *testing.T, out *provider.ToggleAvailabilityOutput) {
				t.Helper()
				assert.True(t, out.Online)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			finder := &mockProviderFinder{provider: tt.found, err: tt.findErr}
			updater := &mockProviderAvailabilityUpdater{err: tt.updateErr}
			geo := &mockGeoCache{addErr: tt.geoAddErr, removeErr: tt.geoRemErr}
			pub := &mockEventPublisher{err: tt.pubErr}
			log := slog.Default()

			uc := provider.NewToggleAvailabilityUseCase(finder, updater, geo, pub, log)
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
			if tt.checkGeo != nil {
				tt.checkGeo(t, geo)
			}
		})
	}
}
