package evidenceuc

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	domainerrors "github.com/davidlramirez95/towcommand/internal/domain/errors"
	"github.com/davidlramirez95/towcommand/internal/domain/evidence"
	"github.com/davidlramirez95/towcommand/internal/usecase/port"
)

// --- Mocks ---

type mockImageValidator struct{ mock.Mock }

func (m *mockImageValidator) ValidateVehiclePhoto(ctx context.Context, s3Bucket, s3Key string) (*port.ImageValidationResult, error) {
	args := m.Called(ctx, s3Bucket, s3Key)
	if v := args.Get(0); v != nil {
		return v.(*port.ImageValidationResult), args.Error(1)
	}
	return nil, args.Error(1)
}

type mockMediaItemAdder struct{ mock.Mock }

func (m *mockMediaItemAdder) AddMediaItem(ctx context.Context, bookingID string, item *evidence.MediaItem) error {
	args := m.Called(ctx, bookingID, item)
	return args.Error(0)
}

type mockEventPublisher struct{ mock.Mock }

func (m *mockEventPublisher) Publish(ctx context.Context, source, detailType string, detail any, actor *Actor) error {
	args := m.Called(ctx, source, detailType, detail, actor)
	return args.Error(0)
}

// --- Tests ---

func TestProcessPhotoUseCase_Execute_Success(t *testing.T) {
	validatorMock := new(mockImageValidator)
	mediaMock := new(mockMediaItemAdder)
	eventsMock := new(mockEventPublisher)

	uc := NewProcessPhotoUseCase(validatorMock, mediaMock, eventsMock)
	uc.idGen = func() string { return "media-test-id" }
	fixedTime := time.Date(2026, 3, 4, 10, 0, 0, 0, time.UTC)
	uc.now = func() time.Time { return fixedTime }

	validatorMock.On("ValidateVehiclePhoto", mock.Anything, "test-bucket", "evidence/booking-1/pickup/FRONT_123.jpg").
		Return(&port.ImageValidationResult{
			IsValid: true,
			Labels:  []string{"Car", "Vehicle"},
		}, nil)

	mediaMock.On("AddMediaItem", mock.Anything, "booking-1", mock.MatchedBy(func(item *evidence.MediaItem) bool {
		return item.MediaID == "media-test-id" &&
			item.S3Key == "evidence/booking-1/pickup/FRONT_123.jpg" &&
			item.Position == evidence.PhotoPositionFront &&
			item.MimeType == "image/jpeg" &&
			item.Integrity.Algorithm == "SHA-256" &&
			item.Integrity.Hash == "abc123hash"
	})).Return(nil)

	eventsMock.On("Publish", mock.Anything, eventSourceEvidence, eventEvidenceValidated, mock.Anything, mock.Anything).
		Return(nil)

	result, err := uc.Execute(context.Background(), &ProcessPhotoInput{
		BookingID: "booking-1",
		S3Key:     "evidence/booking-1/pickup/FRONT_123.jpg",
		S3Bucket:  "test-bucket",
		Position:  evidence.PhotoPositionFront,
		MimeType:  "image/jpeg",
		FileHash:  "abc123hash",
	})

	require.NoError(t, err)
	assert.Equal(t, "media-test-id", result.MediaID)
	assert.Equal(t, evidence.PhotoPositionFront, result.Position)
	assert.Equal(t, "SHA-256", result.Integrity.Algorithm)
	assert.Equal(t, fixedTime, result.CapturedAt)
	validatorMock.AssertExpectations(t)
	mediaMock.AssertExpectations(t)
	eventsMock.AssertExpectations(t)
}

func TestProcessPhotoUseCase_Execute_ValidationFailed(t *testing.T) {
	validatorMock := new(mockImageValidator)
	mediaMock := new(mockMediaItemAdder)
	eventsMock := new(mockEventPublisher)

	uc := NewProcessPhotoUseCase(validatorMock, mediaMock, eventsMock)

	validatorMock.On("ValidateVehiclePhoto", mock.Anything, "test-bucket", "evidence/photo.jpg").
		Return(&port.ImageValidationResult{
			IsValid: false,
			Labels:  []string{"Person"},
			Reason:  "no vehicle detected",
		}, nil)

	_, err := uc.Execute(context.Background(), &ProcessPhotoInput{
		BookingID: "booking-1",
		S3Key:     "evidence/photo.jpg",
		S3Bucket:  "test-bucket",
		Position:  evidence.PhotoPositionFront,
		MimeType:  "image/jpeg",
		FileHash:  "abc123",
	})

	require.Error(t, err)
	var appErr *domainerrors.AppError
	require.ErrorAs(t, err, &appErr)
	assert.Equal(t, domainerrors.CodeEvidenceValidationFailed, appErr.Code)
	assert.Contains(t, appErr.Message, "no vehicle detected")
	mediaMock.AssertNotCalled(t, "AddMediaItem")
}

func TestProcessPhotoUseCase_Execute_RekognitionError(t *testing.T) {
	validatorMock := new(mockImageValidator)
	mediaMock := new(mockMediaItemAdder)
	eventsMock := new(mockEventPublisher)

	uc := NewProcessPhotoUseCase(validatorMock, mediaMock, eventsMock)

	validatorMock.On("ValidateVehiclePhoto", mock.Anything, "test-bucket", "evidence/photo.jpg").
		Return(nil, assert.AnError)

	_, err := uc.Execute(context.Background(), &ProcessPhotoInput{
		BookingID: "booking-1",
		S3Key:     "evidence/photo.jpg",
		S3Bucket:  "test-bucket",
		Position:  evidence.PhotoPositionFront,
		MimeType:  "image/jpeg",
		FileHash:  "abc123",
	})

	require.Error(t, err)
	var appErr *domainerrors.AppError
	require.ErrorAs(t, err, &appErr)
	assert.Equal(t, domainerrors.CodeExternalService, appErr.Code)
}

func TestProcessPhotoUseCase_Execute_SaveError(t *testing.T) {
	validatorMock := new(mockImageValidator)
	mediaMock := new(mockMediaItemAdder)
	eventsMock := new(mockEventPublisher)

	uc := NewProcessPhotoUseCase(validatorMock, mediaMock, eventsMock)
	uc.idGen = func() string { return "media-test-id" }
	uc.now = func() time.Time { return time.Now().UTC() }

	validatorMock.On("ValidateVehiclePhoto", mock.Anything, mock.Anything, mock.Anything).
		Return(&port.ImageValidationResult{IsValid: true, Labels: []string{"Car"}}, nil)

	mediaMock.On("AddMediaItem", mock.Anything, mock.Anything, mock.Anything).
		Return(domainerrors.NewInternalError("db error"))

	_, err := uc.Execute(context.Background(), &ProcessPhotoInput{
		BookingID: "booking-1",
		S3Key:     "evidence/photo.jpg",
		S3Bucket:  "test-bucket",
		Position:  evidence.PhotoPositionFront,
		MimeType:  "image/jpeg",
		FileHash:  "abc123",
	})

	require.Error(t, err)
	eventsMock.AssertNotCalled(t, "Publish")
}

func TestProcessPhotoUseCase_Execute_EventPublishErrorDoesNotFail(t *testing.T) {
	validatorMock := new(mockImageValidator)
	mediaMock := new(mockMediaItemAdder)
	eventsMock := new(mockEventPublisher)

	uc := NewProcessPhotoUseCase(validatorMock, mediaMock, eventsMock)
	uc.idGen = func() string { return "media-test-id" }
	uc.now = func() time.Time { return time.Now().UTC() }

	validatorMock.On("ValidateVehiclePhoto", mock.Anything, mock.Anything, mock.Anything).
		Return(&port.ImageValidationResult{IsValid: true, Labels: []string{"Car"}}, nil)
	mediaMock.On("AddMediaItem", mock.Anything, mock.Anything, mock.Anything).Return(nil)
	eventsMock.On("Publish", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
		Return(domainerrors.NewExternalServiceError("EventBridge", nil))

	result, err := uc.Execute(context.Background(), &ProcessPhotoInput{
		BookingID: "booking-1",
		S3Key:     "evidence/photo.jpg",
		S3Bucket:  "test-bucket",
		Position:  evidence.PhotoPositionFront,
		MimeType:  "image/jpeg",
		FileHash:  "abc123",
	})

	require.NoError(t, err)
	assert.NotEmpty(t, result.MediaID)
}

func TestProcessPhotoUseCase_Execute_InputValidation(t *testing.T) {
	tests := []struct {
		name  string
		input ProcessPhotoInput
		want  string
	}{
		{
			name:  "missing bookingId",
			input: ProcessPhotoInput{S3Key: "k", S3Bucket: "b", Position: "FRONT", MimeType: "image/jpeg", FileHash: "h"},
			want:  "bookingId is required",
		},
		{
			name:  "missing s3Key",
			input: ProcessPhotoInput{BookingID: "b-1", S3Bucket: "b", Position: "FRONT", MimeType: "image/jpeg", FileHash: "h"},
			want:  "s3Key is required",
		},
		{
			name:  "missing s3Bucket",
			input: ProcessPhotoInput{BookingID: "b-1", S3Key: "k", Position: "FRONT", MimeType: "image/jpeg", FileHash: "h"},
			want:  "s3Bucket is required",
		},
		{
			name:  "missing position",
			input: ProcessPhotoInput{BookingID: "b-1", S3Key: "k", S3Bucket: "b", MimeType: "image/jpeg", FileHash: "h"},
			want:  "position is required",
		},
		{
			name:  "missing mimeType",
			input: ProcessPhotoInput{BookingID: "b-1", S3Key: "k", S3Bucket: "b", Position: "FRONT", FileHash: "h"},
			want:  "mimeType is required",
		},
		{
			name:  "missing fileHash",
			input: ProcessPhotoInput{BookingID: "b-1", S3Key: "k", S3Bucket: "b", Position: "FRONT", MimeType: "image/jpeg"},
			want:  "fileHash is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			uc := NewProcessPhotoUseCase(new(mockImageValidator), new(mockMediaItemAdder), new(mockEventPublisher))

			_, err := uc.Execute(context.Background(), &tt.input)

			require.Error(t, err)
			var appErr *domainerrors.AppError
			require.ErrorAs(t, err, &appErr)
			assert.Equal(t, domainerrors.CodeValidationError, appErr.Code)
			assert.Contains(t, appErr.Message, tt.want)
		})
	}
}
