package handler_test

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/aws/aws-lambda-go/events"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/davidlramirez95/towcommand/internal/adapter/handler"
	diagnosisuc "github.com/davidlramirez95/towcommand/internal/usecase/diagnosis"
)

// ---------------------------------------------------------------------------
// Mock DiagnosisEngine for handler tests
// ---------------------------------------------------------------------------

type mockDiagnosisEngineForHandler struct {
	DiagnoseFunc func(ctx context.Context, input *diagnosisuc.DiagnosisInput) (*diagnosisuc.DiagnosisResult, error)
}

func (m *mockDiagnosisEngineForHandler) Diagnose(ctx context.Context, input *diagnosisuc.DiagnosisInput) (*diagnosisuc.DiagnosisResult, error) {
	if m.DiagnoseFunc != nil {
		return m.DiagnoseFunc(ctx, input)
	}
	return nil, nil
}

// ---------------------------------------------------------------------------
// DiagnoseHandler tests
// ---------------------------------------------------------------------------

func TestDiagnoseHandler(t *testing.T) {
	tests := []struct {
		name        string
		event       *events.APIGatewayProxyRequest
		setupEngine func(m *mockDiagnosisEngineForHandler)
		wantStatus  int
		wantErrCode string
		checkBody   func(t *testing.T, body string)
	}{
		{
			name: "success - full input with location",
			event: func() *events.APIGatewayProxyRequest {
				e := apiEventWithAuth("user-1")
				e.Body = `{"description":"My car battery is dead and the lights are dim","photoUrls":["https://s3.example.com/photo1.jpg"],"vehicleType":"sedan","lat":14.5995,"lng":120.9842}`
				return e
			}(),
			setupEngine: func(m *mockDiagnosisEngineForHandler) {
				m.DiagnoseFunc = func(_ context.Context, input *diagnosisuc.DiagnosisInput) (*diagnosisuc.DiagnosisResult, error) {
					assert.Equal(t, "sedan", input.VehicleType)
					assert.NotNil(t, input.Location)
					assert.InDelta(t, 14.5995, input.Location.Lat, 0.0001)
					assert.Len(t, input.PhotoURLs, 1)
					return &diagnosisuc.DiagnosisResult{
						RecommendedService: "JUMPSTART",
						UrgencyLevel:       "HIGH",
						EstimatedCostMin:   100000,
						EstimatedCostMax:   200000,
						Description:        "Dead battery detected.",
						SafetyWarnings:     []string{"Turn off accessories"},
					}, nil
				}
			},
			wantStatus: http.StatusOK,
			checkBody: func(t *testing.T, body string) {
				t.Helper()
				var result diagnosisuc.DiagnosisResult
				require.NoError(t, json.Unmarshal([]byte(body), &result))
				assert.Equal(t, "JUMPSTART", result.RecommendedService)
				assert.Equal(t, "HIGH", result.UrgencyLevel)
				assert.Equal(t, int64(100000), result.EstimatedCostMin)
				assert.Equal(t, int64(200000), result.EstimatedCostMax)
				assert.Equal(t, "Dead battery detected.", result.Description)
				assert.Equal(t, []string{"Turn off accessories"}, result.SafetyWarnings)
			},
		},
		{
			name: "success - minimal input without location",
			event: func() *events.APIGatewayProxyRequest {
				e := apiEventWithAuth("user-1")
				e.Body = `{"description":"My motorcycle has a flat tire and I need help"}`
				return e
			}(),
			setupEngine: func(m *mockDiagnosisEngineForHandler) {
				m.DiagnoseFunc = func(_ context.Context, input *diagnosisuc.DiagnosisInput) (*diagnosisuc.DiagnosisResult, error) {
					assert.Nil(t, input.Location)
					assert.Empty(t, input.VehicleType)
					return &diagnosisuc.DiagnosisResult{
						RecommendedService: "TIRE_CHANGE",
						UrgencyLevel:       "MEDIUM",
						EstimatedCostMin:   80000,
						EstimatedCostMax:   150000,
						Description:        "Flat tire detected.",
						SafetyWarnings:     []string{},
					}, nil
				}
			},
			wantStatus: http.StatusOK,
			checkBody: func(t *testing.T, body string) {
				t.Helper()
				var result diagnosisuc.DiagnosisResult
				require.NoError(t, json.Unmarshal([]byte(body), &result))
				assert.Equal(t, "TIRE_CHANGE", result.RecommendedService)
			},
		},
		{
			name: "unauthorized - no user ID",
			event: &events.APIGatewayProxyRequest{
				Body: `{"description":"My car won't start and I need help immediately"}`,
			},
			wantStatus:  http.StatusUnauthorized,
			wantErrCode: "UNAUTHORIZED",
		},
		{
			name: "invalid body - bad JSON",
			event: func() *events.APIGatewayProxyRequest {
				e := apiEventWithAuth("user-1")
				e.Body = `{not json}`
				return e
			}(),
			wantStatus:  http.StatusBadRequest,
			wantErrCode: "VALIDATION_ERROR",
		},
		{
			name: "validation error - description too short",
			event: func() *events.APIGatewayProxyRequest {
				e := apiEventWithAuth("user-1")
				e.Body = `{"description":"help"}`
				return e
			}(),
			wantStatus:  http.StatusBadRequest,
			wantErrCode: "VALIDATION_ERROR",
		},
		{
			name: "validation error - missing description",
			event: func() *events.APIGatewayProxyRequest {
				e := apiEventWithAuth("user-1")
				e.Body = `{"vehicleType":"sedan"}`
				return e
			}(),
			wantStatus:  http.StatusBadRequest,
			wantErrCode: "VALIDATION_ERROR",
		},
		{
			name: "use case error - engine failure",
			event: func() *events.APIGatewayProxyRequest {
				e := apiEventWithAuth("user-1")
				e.Body = `{"description":"My car is stuck in a flood and the engine won't start"}`
				return e
			}(),
			setupEngine: func(m *mockDiagnosisEngineForHandler) {
				m.DiagnoseFunc = func(_ context.Context, _ *diagnosisuc.DiagnosisInput) (*diagnosisuc.DiagnosisResult, error) {
					return nil, assert.AnError
				}
			},
			wantStatus:  http.StatusBadGateway,
			wantErrCode: "EXTERNAL_SERVICE_ERROR",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			engine := &mockDiagnosisEngineForHandler{}
			if tt.setupEngine != nil {
				tt.setupEngine(engine)
			}

			uc := diagnosisuc.NewDiagnoseUseCase(engine)
			h := handler.NewDiagnoseHandler(uc)

			resp, err := h.Handle(context.Background(), tt.event)

			require.NoError(t, err)
			assert.Equal(t, tt.wantStatus, resp.StatusCode)

			if tt.wantErrCode != "" {
				eb := parseErrorBody(t, resp.Body)
				assert.Equal(t, tt.wantErrCode, eb.Error.Code)
			}

			if tt.checkBody != nil {
				tt.checkBody(t, resp.Body)
			}
		})
	}
}
