package gateway

import (
	"context"
	"errors"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/rekognition"
	rektypes "github.com/aws/aws-sdk-go-v2/service/rekognition/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// --- Mocks ---

type mockRekognitionAPI struct{ mock.Mock }

func (m *mockRekognitionAPI) DetectLabels(ctx context.Context, params *rekognition.DetectLabelsInput, optFns ...func(*rekognition.Options)) (*rekognition.DetectLabelsOutput, error) {
	args := m.Called(ctx, params, optFns)
	if v := args.Get(0); v != nil {
		return v.(*rekognition.DetectLabelsOutput), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *mockRekognitionAPI) DetectModerationLabels(ctx context.Context, params *rekognition.DetectModerationLabelsInput, optFns ...func(*rekognition.Options)) (*rekognition.DetectModerationLabelsOutput, error) {
	args := m.Called(ctx, params, optFns)
	if v := args.Get(0); v != nil {
		return v.(*rekognition.DetectModerationLabelsOutput), args.Error(1)
	}
	return nil, args.Error(1)
}

// --- Tests ---

func TestRekognitionValidator_ValidVehicle(t *testing.T) {
	rekMock := new(mockRekognitionAPI)
	validator := NewRekognitionValidator(rekMock)

	// No moderation labels
	rekMock.On("DetectModerationLabels", mock.Anything, mock.Anything, mock.Anything).
		Return(&rekognition.DetectModerationLabelsOutput{
			ModerationLabels: []rektypes.ModerationLabel{},
		}, nil)

	// Vehicle labels detected
	rekMock.On("DetectLabels", mock.Anything, mock.Anything, mock.Anything).
		Return(&rekognition.DetectLabelsOutput{
			Labels: []rektypes.Label{
				{Name: aws.String("Car"), Confidence: aws.Float32(95.0)},
				{Name: aws.String("Vehicle"), Confidence: aws.Float32(97.0)},
				{Name: aws.String("Wheel"), Confidence: aws.Float32(88.0)},
			},
		}, nil)

	result, err := validator.ValidateVehiclePhoto(context.Background(), "test-bucket", "evidence/photo.jpg")

	require.NoError(t, err)
	assert.True(t, result.IsValid)
	assert.Empty(t, result.Reason)
	assert.Contains(t, result.Labels, "Car")
	assert.Contains(t, result.Labels, "Vehicle")
	rekMock.AssertExpectations(t)
}

func TestRekognitionValidator_NoVehicleDetected(t *testing.T) {
	rekMock := new(mockRekognitionAPI)
	validator := NewRekognitionValidator(rekMock)

	rekMock.On("DetectModerationLabels", mock.Anything, mock.Anything, mock.Anything).
		Return(&rekognition.DetectModerationLabelsOutput{
			ModerationLabels: []rektypes.ModerationLabel{},
		}, nil)

	rekMock.On("DetectLabels", mock.Anything, mock.Anything, mock.Anything).
		Return(&rekognition.DetectLabelsOutput{
			Labels: []rektypes.Label{
				{Name: aws.String("Person"), Confidence: aws.Float32(95.0)},
				{Name: aws.String("Building"), Confidence: aws.Float32(88.0)},
			},
		}, nil)

	result, err := validator.ValidateVehiclePhoto(context.Background(), "test-bucket", "evidence/photo.jpg")

	require.NoError(t, err)
	assert.False(t, result.IsValid)
	assert.Equal(t, "no vehicle detected", result.Reason)
	assert.Contains(t, result.Labels, "Person")
	rekMock.AssertExpectations(t)
}

func TestRekognitionValidator_InappropriateContent(t *testing.T) {
	rekMock := new(mockRekognitionAPI)
	validator := NewRekognitionValidator(rekMock)

	rekMock.On("DetectModerationLabels", mock.Anything, mock.Anything, mock.Anything).
		Return(&rekognition.DetectModerationLabelsOutput{
			ModerationLabels: []rektypes.ModerationLabel{
				{Name: aws.String("Explicit Nudity"), Confidence: aws.Float32(92.0)},
			},
		}, nil)

	result, err := validator.ValidateVehiclePhoto(context.Background(), "test-bucket", "evidence/photo.jpg")

	require.NoError(t, err)
	assert.False(t, result.IsValid)
	assert.Equal(t, "inappropriate content detected", result.Reason)
	assert.Contains(t, result.Labels, "Explicit Nudity")
	// DetectLabels should NOT be called after moderation failure
	rekMock.AssertNotCalled(t, "DetectLabels")
}

func TestRekognitionValidator_ModerationError(t *testing.T) {
	rekMock := new(mockRekognitionAPI)
	validator := NewRekognitionValidator(rekMock)

	rekMock.On("DetectModerationLabels", mock.Anything, mock.Anything, mock.Anything).
		Return(nil, errors.New("rekognition unavailable"))

	result, err := validator.ValidateVehiclePhoto(context.Background(), "test-bucket", "evidence/photo.jpg")

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "detecting moderation labels")
	rekMock.AssertExpectations(t)
}

func TestRekognitionValidator_DetectLabelsError(t *testing.T) {
	rekMock := new(mockRekognitionAPI)
	validator := NewRekognitionValidator(rekMock)

	rekMock.On("DetectModerationLabels", mock.Anything, mock.Anything, mock.Anything).
		Return(&rekognition.DetectModerationLabelsOutput{
			ModerationLabels: []rektypes.ModerationLabel{},
		}, nil)

	rekMock.On("DetectLabels", mock.Anything, mock.Anything, mock.Anything).
		Return(nil, errors.New("label detection failed"))

	result, err := validator.ValidateVehiclePhoto(context.Background(), "test-bucket", "evidence/photo.jpg")

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "detecting labels")
	rekMock.AssertExpectations(t)
}
