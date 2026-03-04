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

func TestRegisterUseCase_Execute(t *testing.T) {
	validInput := provider.RegisterInput{
		CognitoSub:          "cognito-sub-123",
		Name:                "Juan Dela Cruz",
		Phone:               "+639171234567",
		Email:               "juan@example.com",
		TruckType:           "flatbed",
		MaxWeightCapacityKg: 5000,
		PlateNumber:         "ABC-1234",
		LTORegistration:     "LTO-REG-001",
		ServiceAreas:        []string{"NCR"},
	}

	tests := []struct {
		name      string
		input     provider.RegisterInput
		saveErr   error
		pubErr    error
		wantErr   bool
		errCode   domainerrors.ErrorCode
		checkResp func(t *testing.T, out *provider.RegisterOutput)
	}{
		{
			name:  "successful registration",
			input: validInput,
			checkResp: func(t *testing.T, out *provider.RegisterOutput) {
				t.Helper()
				assert.Contains(t, out.ProviderID, "PROV-")
				assert.Equal(t, provdomain.ProviderStatusPendingVerification, out.Status)
			},
		},
		{
			name:    "save fails returns internal error",
			input:   validInput,
			saveErr: errors.New("dynamo timeout"),
			wantErr: true,
			errCode: domainerrors.CodeInternalError,
		},
		{
			name:   "publish fails does not fail registration",
			input:  validInput,
			pubErr: errors.New("eventbridge error"),
			checkResp: func(t *testing.T, out *provider.RegisterOutput) {
				t.Helper()
				assert.Contains(t, out.ProviderID, "PROV-")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			saver := &mockProviderSaver{err: tt.saveErr}
			pub := &mockEventPublisher{err: tt.pubErr}
			log := slog.Default()

			uc := provider.NewRegisterUseCase(saver, pub, log)
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

			// Verify saved provider has correct defaults
			if saver.saved != nil {
				assert.Equal(t, provdomain.ProviderStatusPendingVerification, saver.saved.Status)
				assert.Equal(t, user.TrustTierBasic, saver.saved.TrustTier)
				assert.Equal(t, provdomain.ClearanceStatusPending, saver.saved.NBIClearanceStatus)
				assert.Equal(t, provdomain.ClearanceStatusPending, saver.saved.DrugTestStatus)
				assert.False(t, saver.saved.MMADAccredited)
				assert.False(t, saver.saved.IsOnline)
				assert.Equal(t, float64(0), saver.saved.Rating)
				assert.Equal(t, 1.0, saver.saved.AcceptanceRate)
			}

			if tt.checkResp != nil {
				tt.checkResp(t, result)
			}
		})
	}
}
