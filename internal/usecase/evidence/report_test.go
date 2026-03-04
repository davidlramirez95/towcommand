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

type mockEvidenceSaver struct{ mock.Mock }

func (m *mockEvidenceSaver) Save(ctx context.Context, r *evidence.ConditionReport) error {
	args := m.Called(ctx, r)
	return args.Error(0)
}

type mockEvidenceByBookingLister struct{ mock.Mock }

func (m *mockEvidenceByBookingLister) FindByBooking(ctx context.Context, bookingID string) ([]evidence.ConditionReport, error) {
	args := m.Called(ctx, bookingID)
	if v := args.Get(0); v != nil {
		return v.([]evidence.ConditionReport), args.Error(1)
	}
	return nil, args.Error(1)
}

// --- CreateConditionReportUseCase Tests ---

func TestCreateConditionReportUseCase_Execute_Success(t *testing.T) {
	bookingMock := new(mockBookingFinder)
	evidenceMock := new(mockEvidenceSaver)
	eventsMock := new(mockEventPublisher)

	uc := NewCreateConditionReportUseCase(bookingMock, evidenceMock, eventsMock)
	uc.idGen = func() string { return "report-test-id" }
	fixedTime := time.Date(2026, 3, 4, 10, 0, 0, 0, time.UTC)
	uc.now = func() time.Time { return fixedTime }

	bookingMock.On("FindByID", mock.Anything, "booking-123").
		Return(&booking.Booking{BookingID: "booking-123"}, nil)
	evidenceMock.On("Save", mock.Anything, mock.MatchedBy(func(r *evidence.ConditionReport) bool {
		return r.ReportID == "report-test-id" &&
			r.BookingID == "booking-123" &&
			r.ProviderID == "provider-456" &&
			r.Phase == "pickup" &&
			r.Notes == "minor scratch on hood"
	})).Return(nil)
	eventsMock.On("Publish", mock.Anything, eventSourceEvidence, eventConditionReportCreated, mock.Anything, mock.Anything).
		Return(nil)

	result, err := uc.Execute(context.Background(), &CreateConditionReportInput{
		BookingID:  "booking-123",
		ProviderID: "provider-456",
		Phase:      "pickup",
		Notes:      "minor scratch on hood",
	})

	require.NoError(t, err)
	assert.Equal(t, "report-test-id", result.ReportID)
	assert.Equal(t, "booking-123", result.BookingID)
	assert.Equal(t, "provider-456", result.ProviderID)
	assert.Equal(t, "pickup", result.Phase)
	assert.Equal(t, fixedTime, result.CreatedAt)
	assert.Empty(t, result.Media)
	bookingMock.AssertExpectations(t)
	evidenceMock.AssertExpectations(t)
	eventsMock.AssertExpectations(t)
}

func TestCreateConditionReportUseCase_Execute_BookingNotFound(t *testing.T) {
	bookingMock := new(mockBookingFinder)
	evidenceMock := new(mockEvidenceSaver)
	eventsMock := new(mockEventPublisher)

	uc := NewCreateConditionReportUseCase(bookingMock, evidenceMock, eventsMock)

	bookingMock.On("FindByID", mock.Anything, "booking-missing").Return(nil, nil)

	_, err := uc.Execute(context.Background(), &CreateConditionReportInput{
		BookingID:  "booking-missing",
		ProviderID: "provider-456",
		Phase:      "pickup",
	})

	require.Error(t, err)
	var appErr *domainerrors.AppError
	require.ErrorAs(t, err, &appErr)
	assert.Equal(t, domainerrors.CodeNotFound, appErr.Code)
	evidenceMock.AssertNotCalled(t, "Save")
}

func TestCreateConditionReportUseCase_Execute_SaveError(t *testing.T) {
	bookingMock := new(mockBookingFinder)
	evidenceMock := new(mockEvidenceSaver)
	eventsMock := new(mockEventPublisher)

	uc := NewCreateConditionReportUseCase(bookingMock, evidenceMock, eventsMock)
	uc.idGen = func() string { return "report-test-id" }
	uc.now = func() time.Time { return time.Now().UTC() }

	bookingMock.On("FindByID", mock.Anything, "booking-123").
		Return(&booking.Booking{BookingID: "booking-123"}, nil)
	evidenceMock.On("Save", mock.Anything, mock.Anything).
		Return(domainerrors.NewInternalError("db error"))

	_, err := uc.Execute(context.Background(), &CreateConditionReportInput{
		BookingID:  "booking-123",
		ProviderID: "provider-456",
		Phase:      "pickup",
	})

	require.Error(t, err)
	eventsMock.AssertNotCalled(t, "Publish")
}

func TestCreateConditionReportUseCase_Execute_ValidationErrors(t *testing.T) {
	tests := []struct {
		name  string
		input CreateConditionReportInput
		want  string
	}{
		{
			name:  "missing bookingId",
			input: CreateConditionReportInput{ProviderID: "p-1", Phase: "pickup"},
			want:  "bookingId is required",
		},
		{
			name:  "missing providerId",
			input: CreateConditionReportInput{BookingID: "b-1", Phase: "pickup"},
			want:  "providerId is required",
		},
		{
			name:  "invalid phase",
			input: CreateConditionReportInput{BookingID: "b-1", ProviderID: "p-1", Phase: "invalid"},
			want:  "phase must be pickup or dropoff",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			uc := NewCreateConditionReportUseCase(new(mockBookingFinder), new(mockEvidenceSaver), new(mockEventPublisher))

			_, err := uc.Execute(context.Background(), &tt.input)

			require.Error(t, err)
			var appErr *domainerrors.AppError
			require.ErrorAs(t, err, &appErr)
			assert.Equal(t, domainerrors.CodeValidationError, appErr.Code)
			assert.Contains(t, appErr.Message, tt.want)
		})
	}
}

func TestCreateConditionReportUseCase_Execute_EventErrorDoesNotFail(t *testing.T) {
	bookingMock := new(mockBookingFinder)
	evidenceMock := new(mockEvidenceSaver)
	eventsMock := new(mockEventPublisher)

	uc := NewCreateConditionReportUseCase(bookingMock, evidenceMock, eventsMock)
	uc.idGen = func() string { return "report-test-id" }
	uc.now = func() time.Time { return time.Now().UTC() }

	bookingMock.On("FindByID", mock.Anything, "booking-123").
		Return(&booking.Booking{BookingID: "booking-123"}, nil)
	evidenceMock.On("Save", mock.Anything, mock.Anything).Return(nil)
	eventsMock.On("Publish", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
		Return(domainerrors.NewExternalServiceError("EventBridge", nil))

	result, err := uc.Execute(context.Background(), &CreateConditionReportInput{
		BookingID:  "booking-123",
		ProviderID: "provider-456",
		Phase:      "dropoff",
	})

	require.NoError(t, err)
	assert.NotEmpty(t, result.ReportID)
}

// --- CheckCompletenessUseCase Tests ---

func TestCheckCompletenessUseCase_Execute_Complete(t *testing.T) {
	evidenceMock := new(mockEvidenceByBookingLister)

	uc := NewCheckCompletenessUseCase(evidenceMock)

	// Build a report with all 8 positions covered
	media := make([]evidence.MediaItem, 0, len(evidence.AllPhotoPositions))
	for _, pos := range evidence.AllPhotoPositions {
		media = append(media, evidence.MediaItem{
			MediaID:  "m-" + string(pos),
			S3Key:    "evidence/photo.jpg",
			Position: pos,
			MimeType: "image/jpeg",
		})
	}

	evidenceMock.On("FindByBooking", mock.Anything, "booking-123").
		Return([]evidence.ConditionReport{
			{ReportID: "r-1", BookingID: "booking-123", Phase: "pickup", Media: media},
		}, nil)

	result, err := uc.Execute(context.Background(), "booking-123")

	require.NoError(t, err)
	assert.True(t, result.IsComplete)
	assert.Equal(t, 8, result.TotalPhotos)
	assert.Equal(t, 8, result.RequiredPhotos)
	assert.Empty(t, result.MissingPositions)
	evidenceMock.AssertExpectations(t)
}

func TestCheckCompletenessUseCase_Execute_Incomplete(t *testing.T) {
	evidenceMock := new(mockEvidenceByBookingLister)

	uc := NewCheckCompletenessUseCase(evidenceMock)

	// Only 3 positions covered
	media := []evidence.MediaItem{
		{MediaID: "m-1", Position: evidence.PhotoPositionFront},
		{MediaID: "m-2", Position: evidence.PhotoPositionRear},
		{MediaID: "m-3", Position: evidence.PhotoPositionLeft},
	}

	evidenceMock.On("FindByBooking", mock.Anything, "booking-123").
		Return([]evidence.ConditionReport{
			{ReportID: "r-1", BookingID: "booking-123", Phase: "pickup", Media: media},
		}, nil)

	result, err := uc.Execute(context.Background(), "booking-123")

	require.NoError(t, err)
	assert.False(t, result.IsComplete)
	assert.Equal(t, 3, result.TotalPhotos)
	assert.Equal(t, 8, result.RequiredPhotos)
	assert.Len(t, result.MissingPositions, 5)
	assert.Contains(t, result.MissingPositions, evidence.PhotoPositionRight)
	assert.Contains(t, result.MissingPositions, evidence.PhotoPositionFrontLeft)
	assert.Contains(t, result.MissingPositions, evidence.PhotoPositionFrontRight)
	assert.Contains(t, result.MissingPositions, evidence.PhotoPositionRearLeft)
	assert.Contains(t, result.MissingPositions, evidence.PhotoPositionRearRight)
}

func TestCheckCompletenessUseCase_Execute_NoReports(t *testing.T) {
	evidenceMock := new(mockEvidenceByBookingLister)

	uc := NewCheckCompletenessUseCase(evidenceMock)

	evidenceMock.On("FindByBooking", mock.Anything, "booking-123").
		Return([]evidence.ConditionReport{}, nil)

	result, err := uc.Execute(context.Background(), "booking-123")

	require.NoError(t, err)
	assert.False(t, result.IsComplete)
	assert.Equal(t, 0, result.TotalPhotos)
	assert.Len(t, result.MissingPositions, 8)
}

func TestCheckCompletenessUseCase_Execute_MultipleReports(t *testing.T) {
	evidenceMock := new(mockEvidenceByBookingLister)

	uc := NewCheckCompletenessUseCase(evidenceMock)

	// Spread positions across two reports
	media1 := []evidence.MediaItem{
		{MediaID: "m-1", Position: evidence.PhotoPositionFront},
		{MediaID: "m-2", Position: evidence.PhotoPositionRear},
		{MediaID: "m-3", Position: evidence.PhotoPositionLeft},
		{MediaID: "m-4", Position: evidence.PhotoPositionRight},
	}
	media2 := []evidence.MediaItem{
		{MediaID: "m-5", Position: evidence.PhotoPositionFrontLeft},
		{MediaID: "m-6", Position: evidence.PhotoPositionFrontRight},
		{MediaID: "m-7", Position: evidence.PhotoPositionRearLeft},
		{MediaID: "m-8", Position: evidence.PhotoPositionRearRight},
	}

	evidenceMock.On("FindByBooking", mock.Anything, "booking-123").
		Return([]evidence.ConditionReport{
			{ReportID: "r-1", Phase: "pickup", Media: media1},
			{ReportID: "r-2", Phase: "dropoff", Media: media2},
		}, nil)

	result, err := uc.Execute(context.Background(), "booking-123")

	require.NoError(t, err)
	assert.True(t, result.IsComplete)
	assert.Equal(t, 8, result.TotalPhotos)
	assert.Empty(t, result.MissingPositions)
}

func TestCheckCompletenessUseCase_Execute_MissingBookingID(t *testing.T) {
	uc := NewCheckCompletenessUseCase(new(mockEvidenceByBookingLister))

	_, err := uc.Execute(context.Background(), "")

	require.Error(t, err)
	var appErr *domainerrors.AppError
	require.ErrorAs(t, err, &appErr)
	assert.Equal(t, domainerrors.CodeValidationError, appErr.Code)
}

func TestCheckCompletenessUseCase_Execute_RepoError(t *testing.T) {
	evidenceMock := new(mockEvidenceByBookingLister)

	uc := NewCheckCompletenessUseCase(evidenceMock)

	evidenceMock.On("FindByBooking", mock.Anything, "booking-123").
		Return(nil, domainerrors.NewInternalError("db error"))

	_, err := uc.Execute(context.Background(), "booking-123")

	require.Error(t, err)
}
