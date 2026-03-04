package evidenceuc

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/davidlramirez95/towcommand/internal/domain/booking"
	domainerrors "github.com/davidlramirez95/towcommand/internal/domain/errors"
	"github.com/davidlramirez95/towcommand/internal/domain/evidence"
)

// --- Mocks ---

type mockBookingFinder struct{ mock.Mock }

func (m *mockBookingFinder) FindByID(ctx context.Context, bookingID string) (*booking.Booking, error) {
	args := m.Called(ctx, bookingID)
	if v := args.Get(0); v != nil {
		return v.(*booking.Booking), args.Error(1)
	}
	return nil, args.Error(1)
}

type mockPresigner struct{ mock.Mock }

func (m *mockPresigner) GenerateUploadURL(ctx context.Context, key, contentType string, expiry time.Duration) (string, error) {
	args := m.Called(ctx, key, contentType, expiry)
	return args.String(0), args.Error(1)
}

// --- Tests ---

func TestGenerateUploadURLUseCase_Execute_Success(t *testing.T) {
	bookingMock := new(mockBookingFinder)
	presignerMock := new(mockPresigner)

	uc := NewGenerateUploadURLUseCase(bookingMock, presignerMock)
	fixedTime := time.Date(2026, 3, 4, 10, 0, 0, 0, time.UTC)
	uc.now = func() time.Time { return fixedTime }

	bookingMock.On("FindByID", mock.Anything, "booking-123").
		Return(&booking.Booking{BookingID: "booking-123"}, nil)

	expectedKey := "evidence/booking-123/pickup/FRONT_1772618400.jpg"
	presignerMock.On("GenerateUploadURL", mock.Anything, expectedKey, "image/jpeg", 900*time.Second).
		Return("https://s3.amazonaws.com/presigned-url", nil)

	result, err := uc.Execute(context.Background(), &GenerateUploadURLInput{
		BookingID:   "booking-123",
		Phase:       "pickup",
		Position:    evidence.PhotoPositionFront,
		ContentType: "image/jpeg",
	})

	require.NoError(t, err)
	assert.Equal(t, "https://s3.amazonaws.com/presigned-url", result.UploadURL)
	assert.Equal(t, expectedKey, result.S3Key)
	assert.Equal(t, 900, result.ExpiresIn)
	bookingMock.AssertExpectations(t)
	presignerMock.AssertExpectations(t)
}

func TestGenerateUploadURLUseCase_Execute_BookingNotFound(t *testing.T) {
	bookingMock := new(mockBookingFinder)
	presignerMock := new(mockPresigner)

	uc := NewGenerateUploadURLUseCase(bookingMock, presignerMock)

	bookingMock.On("FindByID", mock.Anything, "booking-missing").Return(nil, nil)

	_, err := uc.Execute(context.Background(), &GenerateUploadURLInput{
		BookingID:   "booking-missing",
		Phase:       "pickup",
		Position:    evidence.PhotoPositionFront,
		ContentType: "image/jpeg",
	})

	require.Error(t, err)
	var appErr *domainerrors.AppError
	require.ErrorAs(t, err, &appErr)
	assert.Equal(t, domainerrors.CodeNotFound, appErr.Code)
	presignerMock.AssertNotCalled(t, "GenerateUploadURL")
}

func TestGenerateUploadURLUseCase_Execute_ValidationErrors(t *testing.T) {
	tests := []struct {
		name  string
		input GenerateUploadURLInput
		want  string
	}{
		{
			name:  "missing bookingId",
			input: GenerateUploadURLInput{Phase: "pickup", Position: evidence.PhotoPositionFront, ContentType: "image/jpeg"},
			want:  "bookingId is required",
		},
		{
			name:  "invalid phase",
			input: GenerateUploadURLInput{BookingID: "b-1", Phase: "invalid", Position: evidence.PhotoPositionFront, ContentType: "image/jpeg"},
			want:  "phase must be pickup or dropoff",
		},
		{
			name:  "missing position",
			input: GenerateUploadURLInput{BookingID: "b-1", Phase: "pickup", ContentType: "image/jpeg"},
			want:  "position is required",
		},
		{
			name:  "missing contentType",
			input: GenerateUploadURLInput{BookingID: "b-1", Phase: "dropoff", Position: evidence.PhotoPositionRear},
			want:  "contentType is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			uc := NewGenerateUploadURLUseCase(new(mockBookingFinder), new(mockPresigner))

			_, err := uc.Execute(context.Background(), &tt.input)

			require.Error(t, err)
			var appErr *domainerrors.AppError
			require.ErrorAs(t, err, &appErr)
			assert.Equal(t, domainerrors.CodeValidationError, appErr.Code)
			assert.Contains(t, appErr.Message, tt.want)
		})
	}
}

func TestGenerateUploadURLUseCase_Execute_PresignError(t *testing.T) {
	bookingMock := new(mockBookingFinder)
	presignerMock := new(mockPresigner)

	uc := NewGenerateUploadURLUseCase(bookingMock, presignerMock)
	uc.now = func() time.Time { return time.Date(2026, 3, 4, 10, 0, 0, 0, time.UTC) }

	bookingMock.On("FindByID", mock.Anything, "booking-123").
		Return(&booking.Booking{BookingID: "booking-123"}, nil)
	presignerMock.On("GenerateUploadURL", mock.Anything, mock.Anything, mock.Anything, mock.Anything).
		Return("", assert.AnError)

	_, err := uc.Execute(context.Background(), &GenerateUploadURLInput{
		BookingID:   "booking-123",
		Phase:       "pickup",
		Position:    evidence.PhotoPositionFront,
		ContentType: "image/jpeg",
	})

	require.Error(t, err)
	var appErr *domainerrors.AppError
	require.ErrorAs(t, err, &appErr)
	assert.Equal(t, domainerrors.CodeExternalService, appErr.Code)
}

func TestGenerateUploadURLUseCase_Execute_BookingRepoError(t *testing.T) {
	bookingMock := new(mockBookingFinder)
	presignerMock := new(mockPresigner)

	uc := NewGenerateUploadURLUseCase(bookingMock, presignerMock)

	bookingMock.On("FindByID", mock.Anything, "booking-123").
		Return(nil, domainerrors.NewInternalError("db error"))

	_, err := uc.Execute(context.Background(), &GenerateUploadURLInput{
		BookingID:   "booking-123",
		Phase:       "pickup",
		Position:    evidence.PhotoPositionFront,
		ContentType: "image/jpeg",
	})

	require.Error(t, err)
	presignerMock.AssertNotCalled(t, "GenerateUploadURL")
}
