package diagnosisuc_test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	domainerrors "github.com/davidlramirez95/towcommand/internal/domain/errors"
	diagnosisuc "github.com/davidlramirez95/towcommand/internal/usecase/diagnosis"
)

// ---------------------------------------------------------------------------
// Mock DiagnosisEngine
// ---------------------------------------------------------------------------

type mockDiagnosisEngine struct {
	DiagnoseFunc func(ctx context.Context, input *diagnosisuc.DiagnosisInput) (*diagnosisuc.DiagnosisResult, error)
}

func (m *mockDiagnosisEngine) Diagnose(ctx context.Context, input *diagnosisuc.DiagnosisInput) (*diagnosisuc.DiagnosisResult, error) {
	if m.DiagnoseFunc != nil {
		return m.DiagnoseFunc(ctx, input)
	}
	return nil, nil
}

// ---------------------------------------------------------------------------
// Tests
// ---------------------------------------------------------------------------

func TestDiagnoseUseCase_Execute(t *testing.T) {
	validResult := &diagnosisuc.DiagnosisResult{
		RecommendedService: "JUMPSTART",
		UrgencyLevel:       "HIGH",
		EstimatedCostMin:   100000,
		EstimatedCostMax:   200000,
		Description:        "Your battery appears to be dead. A jumpstart service is recommended.",
		SafetyWarnings:     []string{"Turn off all electrical accessories before jumpstart"},
	}

	tests := []struct {
		name        string
		input       *diagnosisuc.DiagnosisInput
		setupEngine func(m *mockDiagnosisEngine)
		wantResult  *diagnosisuc.DiagnosisResult
		wantErrCode domainerrors.ErrorCode
	}{
		{
			name: "success - full input",
			input: &diagnosisuc.DiagnosisInput{
				Description: "My car won't start, the battery seems dead and lights are dim",
				PhotoURLs:   []string{"https://s3.example.com/photo1.jpg"},
				VehicleType: "sedan",
				Location:    &diagnosisuc.LatLng{Lat: 14.5995, Lng: 120.9842},
			},
			setupEngine: func(m *mockDiagnosisEngine) {
				m.DiagnoseFunc = func(_ context.Context, _ *diagnosisuc.DiagnosisInput) (*diagnosisuc.DiagnosisResult, error) {
					return validResult, nil
				}
			},
			wantResult: validResult,
		},
		{
			name: "success - minimal input",
			input: &diagnosisuc.DiagnosisInput{
				Description: "Flat tire on my car on EDSA highway",
			},
			setupEngine: func(m *mockDiagnosisEngine) {
				m.DiagnoseFunc = func(_ context.Context, _ *diagnosisuc.DiagnosisInput) (*diagnosisuc.DiagnosisResult, error) {
					return &diagnosisuc.DiagnosisResult{
						RecommendedService: "TIRE_CHANGE",
						UrgencyLevel:       "MEDIUM",
						EstimatedCostMin:   80000,
						EstimatedCostMax:   150000,
						Description:        "Tire change recommended.",
						SafetyWarnings:     []string{},
					}, nil
				}
			},
			wantResult: &diagnosisuc.DiagnosisResult{
				RecommendedService: "TIRE_CHANGE",
				UrgencyLevel:       "MEDIUM",
				EstimatedCostMin:   80000,
				EstimatedCostMax:   150000,
				Description:        "Tire change recommended.",
				SafetyWarnings:     []string{},
			},
		},
		{
			name: "validation error - description too short",
			input: &diagnosisuc.DiagnosisInput{
				Description: "help",
			},
			wantErrCode: domainerrors.CodeValidationError,
		},
		{
			name: "validation error - description too long",
			input: &diagnosisuc.DiagnosisInput{
				Description: string(make([]byte, 1001)),
			},
			wantErrCode: domainerrors.CodeValidationError,
		},
		{
			name: "engine error - external service failure",
			input: &diagnosisuc.DiagnosisInput{
				Description: "My motorcycle fell over and won't start now",
			},
			setupEngine: func(m *mockDiagnosisEngine) {
				m.DiagnoseFunc = func(_ context.Context, _ *diagnosisuc.DiagnosisInput) (*diagnosisuc.DiagnosisResult, error) {
					return nil, errors.New("bedrock throttled")
				}
			},
			wantErrCode: domainerrors.CodeExternalService,
		},
		{
			name: "sanitise - invalid service type defaults to FLATBED_TOWING",
			input: &diagnosisuc.DiagnosisInput{
				Description: "Something weird happened to my car engine",
			},
			setupEngine: func(m *mockDiagnosisEngine) {
				m.DiagnoseFunc = func(_ context.Context, _ *diagnosisuc.DiagnosisInput) (*diagnosisuc.DiagnosisResult, error) {
					return &diagnosisuc.DiagnosisResult{
						RecommendedService: "UNKNOWN_SERVICE",
						UrgencyLevel:       "HIGH",
						EstimatedCostMin:   100000,
						EstimatedCostMax:   200000,
						Description:        "Engine failure detected.",
						SafetyWarnings:     []string{},
					}, nil
				}
			},
			wantResult: &diagnosisuc.DiagnosisResult{
				RecommendedService: "FLATBED_TOWING",
				UrgencyLevel:       "HIGH",
				EstimatedCostMin:   100000,
				EstimatedCostMax:   200000,
				Description:        "Engine failure detected.",
				SafetyWarnings:     []string{},
			},
		},
		{
			name: "sanitise - invalid urgency defaults to MEDIUM",
			input: &diagnosisuc.DiagnosisInput{
				Description: "My car has smoke coming from the hood area",
			},
			setupEngine: func(m *mockDiagnosisEngine) {
				m.DiagnoseFunc = func(_ context.Context, _ *diagnosisuc.DiagnosisInput) (*diagnosisuc.DiagnosisResult, error) {
					return &diagnosisuc.DiagnosisResult{
						RecommendedService: "FLATBED_TOWING",
						UrgencyLevel:       "SUPER_URGENT",
						EstimatedCostMin:   200000,
						EstimatedCostMax:   400000,
						Description:        "Possible engine overheating.",
						SafetyWarnings:     nil,
					}, nil
				}
			},
			wantResult: &diagnosisuc.DiagnosisResult{
				RecommendedService: "FLATBED_TOWING",
				UrgencyLevel:       "MEDIUM",
				EstimatedCostMin:   200000,
				EstimatedCostMax:   400000,
				Description:        "Possible engine overheating.",
				SafetyWarnings:     []string{},
			},
		},
		{
			name: "sanitise - negative cost clamped to zero",
			input: &diagnosisuc.DiagnosisInput{
				Description: "I need help with a flat tire on my vehicle",
			},
			setupEngine: func(m *mockDiagnosisEngine) {
				m.DiagnoseFunc = func(_ context.Context, _ *diagnosisuc.DiagnosisInput) (*diagnosisuc.DiagnosisResult, error) {
					return &diagnosisuc.DiagnosisResult{
						RecommendedService: "TIRE_CHANGE",
						UrgencyLevel:       "LOW",
						EstimatedCostMin:   -5000,
						EstimatedCostMax:   100000,
						Description:        "Simple tire change.",
						SafetyWarnings:     []string{},
					}, nil
				}
			},
			wantResult: &diagnosisuc.DiagnosisResult{
				RecommendedService: "TIRE_CHANGE",
				UrgencyLevel:       "LOW",
				EstimatedCostMin:   0,
				EstimatedCostMax:   100000,
				Description:        "Simple tire change.",
				SafetyWarnings:     []string{},
			},
		},
		{
			name: "sanitise - max less than min gets corrected",
			input: &diagnosisuc.DiagnosisInput{
				Description: "My car is stuck in a ditch and needs to be pulled out",
			},
			setupEngine: func(m *mockDiagnosisEngine) {
				m.DiagnoseFunc = func(_ context.Context, _ *diagnosisuc.DiagnosisInput) (*diagnosisuc.DiagnosisResult, error) {
					return &diagnosisuc.DiagnosisResult{
						RecommendedService: "WINCH_RECOVERY",
						UrgencyLevel:       "HIGH",
						EstimatedCostMin:   300000,
						EstimatedCostMax:   200000, // max < min
						Description:        "Winch recovery required.",
						SafetyWarnings:     []string{},
					}, nil
				}
			},
			wantResult: &diagnosisuc.DiagnosisResult{
				RecommendedService: "WINCH_RECOVERY",
				UrgencyLevel:       "HIGH",
				EstimatedCostMin:   300000,
				EstimatedCostMax:   300000, // corrected to match min
				Description:        "Winch recovery required.",
				SafetyWarnings:     []string{},
			},
		},
		{
			name: "sanitise - empty description gets default",
			input: &diagnosisuc.DiagnosisInput{
				Description: "I ran out of fuel on Commonwealth Avenue",
			},
			setupEngine: func(m *mockDiagnosisEngine) {
				m.DiagnoseFunc = func(_ context.Context, _ *diagnosisuc.DiagnosisInput) (*diagnosisuc.DiagnosisResult, error) {
					return &diagnosisuc.DiagnosisResult{
						RecommendedService: "FUEL_DELIVERY",
						UrgencyLevel:       "LOW",
						EstimatedCostMin:   50000,
						EstimatedCostMax:   100000,
						Description:        "", // empty
						SafetyWarnings:     []string{},
					}, nil
				}
			},
			wantResult: &diagnosisuc.DiagnosisResult{
				RecommendedService: "FUEL_DELIVERY",
				UrgencyLevel:       "LOW",
				EstimatedCostMin:   50000,
				EstimatedCostMax:   100000,
				Description:        "AI diagnosis for: I ran out of fuel on Commonwealth Avenue",
				SafetyWarnings:     []string{},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			engine := &mockDiagnosisEngine{}
			if tt.setupEngine != nil {
				tt.setupEngine(engine)
			}

			uc := diagnosisuc.NewDiagnoseUseCase(engine)
			result, err := uc.Execute(context.Background(), tt.input)

			if tt.wantErrCode != "" {
				require.Error(t, err)
				var appErr *domainerrors.AppError
				require.True(t, errors.As(err, &appErr), "expected AppError, got %T", err)
				assert.Equal(t, tt.wantErrCode, appErr.Code)
				assert.Nil(t, result)
				return
			}

			require.NoError(t, err)
			require.NotNil(t, result)
			assert.Equal(t, tt.wantResult.RecommendedService, result.RecommendedService)
			assert.Equal(t, tt.wantResult.UrgencyLevel, result.UrgencyLevel)
			assert.Equal(t, tt.wantResult.EstimatedCostMin, result.EstimatedCostMin)
			assert.Equal(t, tt.wantResult.EstimatedCostMax, result.EstimatedCostMax)
			assert.Equal(t, tt.wantResult.Description, result.Description)
			assert.Equal(t, tt.wantResult.SafetyWarnings, result.SafetyWarnings)
		})
	}
}
