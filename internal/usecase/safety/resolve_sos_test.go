package safetyuc

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	domainerrors "github.com/davidlramirez95/towcommand/internal/domain/errors"
	"github.com/davidlramirez95/towcommand/internal/domain/safety"
)

func TestResolveSOSUseCase_Execute(t *testing.T) {
	fixedTime := time.Date(2026, 3, 5, 14, 0, 0, 0, time.UTC)

	tests := []struct {
		name        string
		input       *ResolveSOSInput
		setupMocks  func(sf *mockSOSFinder, sr *mockSOSResolver, ep *mockEventPublisher)
		wantErr     bool
		wantErrCode domainerrors.ErrorCode
		checkResult func(t *testing.T, alert *safety.SOSAlert)
	}{
		{
			name: "success",
			input: &ResolveSOSInput{
				AlertID:    "SOS-2026-abc",
				ResolvedBy: "ADMIN-001",
			},
			setupMocks: func(sf *mockSOSFinder, sr *mockSOSResolver, ep *mockEventPublisher) {
				sf.On("FindByID", mock.Anything, "SOS-2026-abc").Return(&safety.SOSAlert{
					AlertID:     "SOS-2026-abc",
					BookingID:   "BK-001",
					TriggeredBy: "USR-001",
					TriggerType: safety.TriggerTypeButton,
					Resolved:    false,
				}, nil)
				sr.On("Resolve", mock.Anything, "SOS-2026-abc", "ADMIN-001", fixedTime).Return(nil)
				ep.On("Publish", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)
			},
			checkResult: func(t *testing.T, alert *safety.SOSAlert) {
				t.Helper()
				assert.True(t, alert.Resolved)
				assert.Equal(t, "ADMIN-001", alert.ResolvedBy)
				require.NotNil(t, alert.ResolvedAt)
				assert.Equal(t, fixedTime, *alert.ResolvedAt)
			},
		},
		{
			name: "not found",
			input: &ResolveSOSInput{
				AlertID:    "SOS-MISSING",
				ResolvedBy: "ADMIN-001",
			},
			setupMocks: func(sf *mockSOSFinder, sr *mockSOSResolver, ep *mockEventPublisher) {
				sf.On("FindByID", mock.Anything, "SOS-MISSING").Return(nil, nil)
			},
			wantErr:     true,
			wantErrCode: domainerrors.CodeNotFound,
		},
		{
			name: "already resolved",
			input: &ResolveSOSInput{
				AlertID:    "SOS-RESOLVED",
				ResolvedBy: "ADMIN-001",
			},
			setupMocks: func(sf *mockSOSFinder, sr *mockSOSResolver, ep *mockEventPublisher) {
				resolved := time.Date(2026, 3, 5, 10, 0, 0, 0, time.UTC)
				sf.On("FindByID", mock.Anything, "SOS-RESOLVED").Return(&safety.SOSAlert{
					AlertID:    "SOS-RESOLVED",
					Resolved:   true,
					ResolvedBy: "ADMIN-002",
					ResolvedAt: &resolved,
				}, nil)
			},
			wantErr:     true,
			wantErrCode: domainerrors.CodeConflict,
		},
		{
			name: "finder error",
			input: &ResolveSOSInput{
				AlertID:    "SOS-ERR",
				ResolvedBy: "ADMIN-001",
			},
			setupMocks: func(sf *mockSOSFinder, sr *mockSOSResolver, ep *mockEventPublisher) {
				sf.On("FindByID", mock.Anything, "SOS-ERR").Return(nil, errors.New("db error"))
			},
			wantErr: true,
		},
		{
			name: "resolver error",
			input: &ResolveSOSInput{
				AlertID:    "SOS-2026-abc",
				ResolvedBy: "ADMIN-001",
			},
			setupMocks: func(sf *mockSOSFinder, sr *mockSOSResolver, ep *mockEventPublisher) {
				sf.On("FindByID", mock.Anything, "SOS-2026-abc").Return(&safety.SOSAlert{
					AlertID:  "SOS-2026-abc",
					Resolved: false,
				}, nil)
				sr.On("Resolve", mock.Anything, "SOS-2026-abc", "ADMIN-001", fixedTime).Return(errors.New("ddb error"))
			},
			wantErr: true,
		},
		{
			name: "validation error: missing alertId",
			input: &ResolveSOSInput{
				ResolvedBy: "ADMIN-001",
			},
			setupMocks:  func(sf *mockSOSFinder, sr *mockSOSResolver, ep *mockEventPublisher) {},
			wantErr:     true,
			wantErrCode: domainerrors.CodeValidationError,
		},
		{
			name: "validation error: missing resolvedBy",
			input: &ResolveSOSInput{
				AlertID: "SOS-2026-abc",
			},
			setupMocks:  func(sf *mockSOSFinder, sr *mockSOSResolver, ep *mockEventPublisher) {},
			wantErr:     true,
			wantErrCode: domainerrors.CodeValidationError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sf := new(mockSOSFinder)
			sr := new(mockSOSResolver)
			ep := new(mockEventPublisher)

			tt.setupMocks(sf, sr, ep)

			uc := NewResolveSOSUseCase(sf, sr, ep)
			uc.now = func() time.Time { return fixedTime }

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
