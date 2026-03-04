package provider_test

import (
	"context"
	"errors"
	"log/slog"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	domainerrors "github.com/davidlramirez95/towcommand/internal/domain/errors"
	provider "github.com/davidlramirez95/towcommand/internal/usecase/provider"
)

func TestUpdateLocationUseCase_Execute(t *testing.T) {
	validInput := provider.UpdateLocationInput{
		ProviderID: "PROV-123",
		Lat:        14.5995,
		Lng:        120.9842,
		Heading:    90,
		Speed:      40,
	}

	tests := []struct {
		name     string
		input    provider.UpdateLocationInput
		geoErr   error
		repoErr  error
		pubErr   error
		wantErr  bool
		errCode  domainerrors.ErrorCode
		checkOut func(t *testing.T, out *provider.UpdateLocationOutput)
	}{
		{
			name:  "successful location update",
			input: validInput,
			checkOut: func(t *testing.T, out *provider.UpdateLocationOutput) {
				t.Helper()
				assert.Equal(t, "PROV-123", out.ProviderID)
				assert.Equal(t, 14.5995, out.Lat)
				assert.Equal(t, 120.9842, out.Lng)
			},
		},
		{
			name: "invalid coordinates outside Philippines",
			input: provider.UpdateLocationInput{
				ProviderID: "PROV-123",
				Lat:        48.8566, // Paris
				Lng:        2.3522,
			},
			wantErr: true,
			errCode: domainerrors.CodeValidationError,
		},
		{
			name:    "geo cache fails returns external service error",
			input:   validInput,
			geoErr:  errors.New("redis down"),
			wantErr: true,
			errCode: domainerrors.CodeExternalService,
		},
		{
			name:    "repo update fails returns internal error",
			input:   validInput,
			repoErr: errors.New("dynamo timeout"),
			wantErr: true,
			errCode: domainerrors.CodeInternalError,
		},
		{
			name:   "publish fails does not fail update",
			input:  validInput,
			pubErr: errors.New("eventbridge error"),
			checkOut: func(t *testing.T, out *provider.UpdateLocationOutput) {
				t.Helper()
				assert.Equal(t, "PROV-123", out.ProviderID)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &mockProviderLocationUpdater{err: tt.repoErr}
			geo := &mockGeoCache{addErr: tt.geoErr}
			pub := &mockEventPublisher{err: tt.pubErr}
			log := slog.Default()

			uc := provider.NewUpdateLocationUseCase(repo, geo, pub, log)
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
