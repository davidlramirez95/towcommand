package gateway

import (
	"context"
	"encoding/json"
	"errors"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/bedrockruntime"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	diagnosisuc "github.com/davidlramirez95/towcommand/internal/usecase/diagnosis"
)

// --- Mock Bedrock Runtime API ---

type mockBedrockRuntimeAPI struct{ mock.Mock }

func (m *mockBedrockRuntimeAPI) InvokeModel(ctx context.Context, params *bedrockruntime.InvokeModelInput, optFns ...func(*bedrockruntime.Options)) (*bedrockruntime.InvokeModelOutput, error) {
	args := m.Called(ctx, params, optFns)
	if v := args.Get(0); v != nil {
		return v.(*bedrockruntime.InvokeModelOutput), args.Error(1)
	}
	return nil, args.Error(1)
}

// --- Helper: build mock Bedrock response body ---

func buildBedrockResponseBody(t *testing.T, result *diagnosisuc.DiagnosisResult) []byte {
	t.Helper()
	resultJSON, err := json.Marshal(result)
	require.NoError(t, err)

	resp := bedrockResponse{
		Content: []bedrockContentBlock{
			{Type: "text", Text: string(resultJSON)},
		},
	}
	body, err := json.Marshal(resp)
	require.NoError(t, err)
	return body
}

// --- Tests ---

func TestBedrockDiagnosisEngine_Success(t *testing.T) {
	brMock := new(mockBedrockRuntimeAPI)
	engine := NewBedrockDiagnosisEngine(brMock)

	expected := &diagnosisuc.DiagnosisResult{
		RecommendedService: "JUMPSTART",
		UrgencyLevel:       "HIGH",
		EstimatedCostMin:   100000,
		EstimatedCostMax:   200000,
		Description:        "Dead battery detected. Jumpstart service recommended.",
		SafetyWarnings:     []string{"Turn off electrical accessories"},
	}

	brMock.On("InvokeModel", mock.Anything, mock.MatchedBy(func(input *bedrockruntime.InvokeModelInput) bool {
		assert.Equal(t, bedrockModelID, *input.ModelId)
		assert.Equal(t, "application/json", *input.ContentType)

		var req bedrockRequest
		require.NoError(t, json.Unmarshal(input.Body, &req))
		assert.Equal(t, "bedrock-2023-05-31", req.AnthropicVersion)
		assert.Equal(t, 1024, req.MaxTokens)
		assert.Contains(t, req.System, "TowCommand PH")
		assert.Len(t, req.Messages, 1)
		assert.Equal(t, "user", req.Messages[0].Role)
		assert.Contains(t, req.Messages[0].Content, "battery is dead")
		return true
	}), mock.Anything).Return(&bedrockruntime.InvokeModelOutput{
		Body: buildBedrockResponseBody(t, expected),
	}, nil)

	result, err := engine.Diagnose(context.Background(), &diagnosisuc.DiagnosisInput{
		Description: "My car battery is dead and the lights are dim",
		VehicleType: "sedan",
		Location:    &diagnosisuc.LatLng{Lat: 14.5995, Lng: 120.9842},
		PhotoURLs:   []string{"https://s3.example.com/photo1.jpg"},
	})

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, "JUMPSTART", result.RecommendedService)
	assert.Equal(t, "HIGH", result.UrgencyLevel)
	assert.Equal(t, int64(100000), result.EstimatedCostMin)
	assert.Equal(t, int64(200000), result.EstimatedCostMax)
	assert.Equal(t, "Dead battery detected. Jumpstart service recommended.", result.Description)
	assert.Equal(t, []string{"Turn off electrical accessories"}, result.SafetyWarnings)
	brMock.AssertExpectations(t)
}

func TestBedrockDiagnosisEngine_MinimalInput(t *testing.T) {
	brMock := new(mockBedrockRuntimeAPI)
	engine := NewBedrockDiagnosisEngine(brMock)

	expected := &diagnosisuc.DiagnosisResult{
		RecommendedService: "TIRE_CHANGE",
		UrgencyLevel:       "MEDIUM",
		EstimatedCostMin:   80000,
		EstimatedCostMax:   150000,
		Description:        "Flat tire detected.",
		SafetyWarnings:     []string{},
	}

	brMock.On("InvokeModel", mock.Anything, mock.MatchedBy(func(input *bedrockruntime.InvokeModelInput) bool {
		var req bedrockRequest
		require.NoError(t, json.Unmarshal(input.Body, &req))
		// Minimal input should not contain vehicle type or location
		assert.NotContains(t, req.Messages[0].Content, "Vehicle type")
		assert.NotContains(t, req.Messages[0].Content, "Location")
		assert.NotContains(t, req.Messages[0].Content, "photo")
		return true
	}), mock.Anything).Return(&bedrockruntime.InvokeModelOutput{
		Body: buildBedrockResponseBody(t, expected),
	}, nil)

	result, err := engine.Diagnose(context.Background(), &diagnosisuc.DiagnosisInput{
		Description: "I have a flat tire on my car",
	})

	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, "TIRE_CHANGE", result.RecommendedService)
	brMock.AssertExpectations(t)
}

func TestBedrockDiagnosisEngine_InvokeModelError(t *testing.T) {
	brMock := new(mockBedrockRuntimeAPI)
	engine := NewBedrockDiagnosisEngine(brMock)

	brMock.On("InvokeModel", mock.Anything, mock.Anything, mock.Anything).
		Return(nil, errors.New("bedrock throttled"))

	result, err := engine.Diagnose(context.Background(), &diagnosisuc.DiagnosisInput{
		Description: "My car won't start and I need help",
	})

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "invoking bedrock model")
	brMock.AssertExpectations(t)
}

func TestBedrockDiagnosisEngine_InvalidResponseJSON(t *testing.T) {
	brMock := new(mockBedrockRuntimeAPI)
	engine := NewBedrockDiagnosisEngine(brMock)

	brMock.On("InvokeModel", mock.Anything, mock.Anything, mock.Anything).
		Return(&bedrockruntime.InvokeModelOutput{
			Body: []byte(`not valid json`),
		}, nil)

	result, err := engine.Diagnose(context.Background(), &diagnosisuc.DiagnosisInput{
		Description: "My car is making strange noises from the engine",
	})

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "unmarshalling bedrock response")
	brMock.AssertExpectations(t)
}

func TestBedrockDiagnosisEngine_EmptyContent(t *testing.T) {
	brMock := new(mockBedrockRuntimeAPI)
	engine := NewBedrockDiagnosisEngine(brMock)

	resp := bedrockResponse{Content: []bedrockContentBlock{}}
	body, _ := json.Marshal(resp)

	brMock.On("InvokeModel", mock.Anything, mock.Anything, mock.Anything).
		Return(&bedrockruntime.InvokeModelOutput{
			Body: body,
		}, nil)

	result, err := engine.Diagnose(context.Background(), &diagnosisuc.DiagnosisInput{
		Description: "My car is overheating and steam is coming out",
	})

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "empty content")
	brMock.AssertExpectations(t)
}

func TestBedrockDiagnosisEngine_MalformedDiagnosisJSON(t *testing.T) {
	brMock := new(mockBedrockRuntimeAPI)
	engine := NewBedrockDiagnosisEngine(brMock)

	resp := bedrockResponse{
		Content: []bedrockContentBlock{
			{Type: "text", Text: `{"recommendedService": "JUMPSTART", "urgencyLevel": INVALID}`},
		},
	}
	body, _ := json.Marshal(resp)

	brMock.On("InvokeModel", mock.Anything, mock.Anything, mock.Anything).
		Return(&bedrockruntime.InvokeModelOutput{
			Body: body,
		}, nil)

	result, err := engine.Diagnose(context.Background(), &diagnosisuc.DiagnosisInput{
		Description: "My car has a strange noise from the transmission",
	})

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "parsing diagnosis JSON")
	brMock.AssertExpectations(t)
}

func TestBuildUserPrompt(t *testing.T) {
	tests := []struct {
		name     string
		input    *diagnosisuc.DiagnosisInput
		contains []string
		excludes []string
	}{
		{
			name: "full input",
			input: &diagnosisuc.DiagnosisInput{
				Description: "Dead battery",
				VehicleType: "sedan",
				Location:    &diagnosisuc.LatLng{Lat: 14.5995, Lng: 120.9842},
				PhotoURLs:   []string{"url1", "url2"},
			},
			contains: []string{
				"Dead battery",
				"Vehicle type: sedan",
				"lat=14.599500",
				"lng=120.984200",
				"2 photo(s)",
			},
		},
		{
			name: "minimal input",
			input: &diagnosisuc.DiagnosisInput{
				Description: "Flat tire",
			},
			contains: []string{"Flat tire"},
			excludes: []string{"Vehicle type", "Location", "photo"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			prompt := buildUserPrompt(tt.input)
			for _, c := range tt.contains {
				assert.Contains(t, prompt, c)
			}
			for _, e := range tt.excludes {
				assert.NotContains(t, prompt, e)
			}
		})
	}
}
