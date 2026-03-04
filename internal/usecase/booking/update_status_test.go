package bookinguc

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/davidlramirez95/towcommand/internal/domain/booking"
	domainerrors "github.com/davidlramirez95/towcommand/internal/domain/errors"
)

// --- Mocks ---

type mockUpdateStatusRepo struct{ mock.Mock }

func (m *mockUpdateStatusRepo) FindByID(ctx context.Context, bookingID string) (*booking.Booking, error) {
	args := m.Called(ctx, bookingID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*booking.Booking), args.Error(1)
}

func (m *mockUpdateStatusRepo) UpdateStatus(ctx context.Context, bookingID string, status booking.BookingStatus, metadata map[string]any) error {
	args := m.Called(ctx, bookingID, status, metadata)
	return args.Error(0)
}

// --- Tests ---

func TestUpdateBookingStatusUseCase_Execute_ProviderTransition(t *testing.T) {
	repoMock := new(mockUpdateStatusRepo)
	eventsMock := new(mockEventPublisher)
	uc := NewUpdateBookingStatusUseCase(repoMock, eventsMock)

	b := &booking.Booking{BookingID: "BK-1", CustomerID: "user-1", ProviderID: "prov-1", Status: booking.BookingStatusMatched}
	repoMock.On("FindByID", mock.Anything, "BK-1").Return(b, nil)
	repoMock.On("UpdateStatus", mock.Anything, "BK-1", booking.BookingStatusEnRoute, mock.Anything).Return(nil)
	eventsMock.On("Publish", mock.Anything, eventSourceBooking, eventBookingStatusChanged, mock.Anything, mock.Anything).Return(nil)

	result, err := uc.Execute(context.Background(), UpdateBookingStatusInput{
		BookingID: "BK-1", CallerID: "prov-1", CallerType: "provider", NewStatus: booking.BookingStatusEnRoute,
	})

	require.NoError(t, err)
	assert.Equal(t, "BK-1", result.BookingID)
	assert.Equal(t, "MATCHED", result.PreviousStatus)
	assert.Equal(t, "EN_ROUTE", result.Status)
	repoMock.AssertExpectations(t)
	eventsMock.AssertExpectations(t)
}

func TestUpdateBookingStatusUseCase_Execute_AdminTransition(t *testing.T) {
	repoMock := new(mockUpdateStatusRepo)
	eventsMock := new(mockEventPublisher)
	uc := NewUpdateBookingStatusUseCase(repoMock, eventsMock)

	b := &booking.Booking{BookingID: "BK-1", CustomerID: "user-1", ProviderID: "prov-1", Status: booking.BookingStatusPending}
	repoMock.On("FindByID", mock.Anything, "BK-1").Return(b, nil)
	repoMock.On("UpdateStatus", mock.Anything, "BK-1", booking.BookingStatusMatched, mock.Anything).Return(nil)
	eventsMock.On("Publish", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)

	result, err := uc.Execute(context.Background(), UpdateBookingStatusInput{
		BookingID: "BK-1", CallerID: "admin-1", CallerType: "admin", NewStatus: booking.BookingStatusMatched,
	})

	require.NoError(t, err)
	assert.Equal(t, "MATCHED", result.Status)
}

func TestUpdateBookingStatusUseCase_Execute_CompletionPublishesTwoEvents(t *testing.T) {
	repoMock := new(mockUpdateStatusRepo)
	eventsMock := new(mockEventPublisher)
	uc := NewUpdateBookingStatusUseCase(repoMock, eventsMock)

	b := &booking.Booking{
		BookingID:  "BK-1",
		CustomerID: "user-1",
		ProviderID: "prov-1",
		Status:     booking.BookingStatusOTPDropoff,
		Price:      booking.PriceBreakdown{Total: 200_000, Currency: "PHP"},
	}
	repoMock.On("FindByID", mock.Anything, "BK-1").Return(b, nil)
	repoMock.On("UpdateStatus", mock.Anything, "BK-1", booking.BookingStatusCompleted, mock.Anything).Return(nil)
	eventsMock.On("Publish", mock.Anything, eventSourceBooking, eventBookingStatusChanged, mock.Anything, mock.Anything).Return(nil)
	eventsMock.On("Publish", mock.Anything, eventSourceBooking, eventBookingCompleted, mock.Anything, mock.Anything).Return(nil)

	result, err := uc.Execute(context.Background(), UpdateBookingStatusInput{
		BookingID: "BK-1", CallerID: "prov-1", CallerType: "provider", NewStatus: booking.BookingStatusCompleted,
	})

	require.NoError(t, err)
	assert.Equal(t, "COMPLETED", result.Status)
	eventsMock.AssertNumberOfCalls(t, "Publish", 2)
}

func TestUpdateBookingStatusUseCase_Execute_NotFound(t *testing.T) {
	repoMock := new(mockUpdateStatusRepo)
	eventsMock := new(mockEventPublisher)
	uc := NewUpdateBookingStatusUseCase(repoMock, eventsMock)

	repoMock.On("FindByID", mock.Anything, "BK-GONE").Return(nil, nil)

	_, err := uc.Execute(context.Background(), UpdateBookingStatusInput{
		BookingID: "BK-GONE", CallerID: "prov-1", CallerType: "provider", NewStatus: booking.BookingStatusEnRoute,
	})

	require.Error(t, err)
	var appErr *domainerrors.AppError
	require.True(t, errors.As(err, &appErr))
	assert.Equal(t, domainerrors.CodeNotFound, appErr.Code)
}

func TestUpdateBookingStatusUseCase_Execute_UnauthorizedProvider(t *testing.T) {
	repoMock := new(mockUpdateStatusRepo)
	eventsMock := new(mockEventPublisher)
	uc := NewUpdateBookingStatusUseCase(repoMock, eventsMock)

	b := &booking.Booking{BookingID: "BK-1", CustomerID: "user-1", ProviderID: "prov-1", Status: booking.BookingStatusMatched}
	repoMock.On("FindByID", mock.Anything, "BK-1").Return(b, nil)

	_, err := uc.Execute(context.Background(), UpdateBookingStatusInput{
		BookingID: "BK-1", CallerID: "prov-OTHER", CallerType: "provider", NewStatus: booking.BookingStatusEnRoute,
	})

	require.Error(t, err)
	var appErr *domainerrors.AppError
	require.True(t, errors.As(err, &appErr))
	assert.Equal(t, domainerrors.CodeForbidden, appErr.Code)
}

func TestUpdateBookingStatusUseCase_Execute_NonAdminCannotSetNonProviderStatus(t *testing.T) {
	repoMock := new(mockUpdateStatusRepo)
	eventsMock := new(mockEventPublisher)
	uc := NewUpdateBookingStatusUseCase(repoMock, eventsMock)

	b := &booking.Booking{BookingID: "BK-1", CustomerID: "user-1", ProviderID: "prov-1", Status: booking.BookingStatusConditionReport}
	repoMock.On("FindByID", mock.Anything, "BK-1").Return(b, nil)

	// OTP_VERIFIED is not a provider status, so a non-admin customer cannot set it
	_, err := uc.Execute(context.Background(), UpdateBookingStatusInput{
		BookingID: "BK-1", CallerID: "user-1", CallerType: "customer", NewStatus: booking.BookingStatusOTPVerified,
	})

	require.Error(t, err)
	var appErr *domainerrors.AppError
	require.True(t, errors.As(err, &appErr))
	assert.Equal(t, domainerrors.CodeForbidden, appErr.Code)
}

func TestUpdateBookingStatusUseCase_Execute_InvalidTransition(t *testing.T) {
	repoMock := new(mockUpdateStatusRepo)
	eventsMock := new(mockEventPublisher)
	uc := NewUpdateBookingStatusUseCase(repoMock, eventsMock)

	b := &booking.Booking{BookingID: "BK-1", CustomerID: "user-1", ProviderID: "prov-1", Status: booking.BookingStatusPending}
	repoMock.On("FindByID", mock.Anything, "BK-1").Return(b, nil)

	// Cannot skip MATCHED and go to EN_ROUTE
	_, err := uc.Execute(context.Background(), UpdateBookingStatusInput{
		BookingID: "BK-1", CallerID: "prov-1", CallerType: "admin", NewStatus: booking.BookingStatusEnRoute,
	})

	require.Error(t, err)
	var appErr *domainerrors.AppError
	require.True(t, errors.As(err, &appErr))
	assert.Equal(t, domainerrors.CodeInvalidStatusTransition, appErr.Code)
}

func TestUpdateBookingStatusUseCase_Execute_AllLinearTransitions(t *testing.T) {
	transitions := []struct {
		from booking.BookingStatus
		to   booking.BookingStatus
	}{
		{booking.BookingStatusPending, booking.BookingStatusMatched},
		{booking.BookingStatusMatched, booking.BookingStatusEnRoute},
		{booking.BookingStatusEnRoute, booking.BookingStatusArrived},
		{booking.BookingStatusArrived, booking.BookingStatusConditionReport},
		{booking.BookingStatusConditionReport, booking.BookingStatusOTPVerified},
		{booking.BookingStatusOTPVerified, booking.BookingStatusLoading},
		{booking.BookingStatusLoading, booking.BookingStatusInTransit},
		{booking.BookingStatusInTransit, booking.BookingStatusArrivedDropoff},
		{booking.BookingStatusArrivedDropoff, booking.BookingStatusOTPDropoff},
		{booking.BookingStatusOTPDropoff, booking.BookingStatusCompleted},
	}

	for _, tt := range transitions {
		t.Run(string(tt.from)+"_to_"+string(tt.to), func(t *testing.T) {
			repoMock := new(mockUpdateStatusRepo)
			eventsMock := new(mockEventPublisher)
			uc := NewUpdateBookingStatusUseCase(repoMock, eventsMock)

			b := &booking.Booking{BookingID: "BK-1", CustomerID: "user-1", ProviderID: "prov-1", Status: tt.from}
			repoMock.On("FindByID", mock.Anything, "BK-1").Return(b, nil)
			repoMock.On("UpdateStatus", mock.Anything, "BK-1", tt.to, mock.Anything).Return(nil)
			eventsMock.On("Publish", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)

			result, err := uc.Execute(context.Background(), UpdateBookingStatusInput{
				BookingID: "BK-1", CallerID: "admin-1", CallerType: "admin", NewStatus: tt.to,
			})

			require.NoError(t, err)
			assert.Equal(t, string(tt.to), result.Status)
			assert.Equal(t, string(tt.from), result.PreviousStatus)
		})
	}
}

func TestUpdateBookingStatusUseCase_Execute_UpdateStatusError(t *testing.T) {
	repoMock := new(mockUpdateStatusRepo)
	eventsMock := new(mockEventPublisher)
	uc := NewUpdateBookingStatusUseCase(repoMock, eventsMock)

	b := &booking.Booking{BookingID: "BK-1", CustomerID: "user-1", ProviderID: "prov-1", Status: booking.BookingStatusMatched}
	repoMock.On("FindByID", mock.Anything, "BK-1").Return(b, nil)
	repoMock.On("UpdateStatus", mock.Anything, "BK-1", booking.BookingStatusEnRoute, mock.Anything).
		Return(domainerrors.NewInternalError("db error"))

	_, err := uc.Execute(context.Background(), UpdateBookingStatusInput{
		BookingID: "BK-1", CallerID: "prov-1", CallerType: "provider", NewStatus: booking.BookingStatusEnRoute,
	})

	assert.Error(t, err)
}

func TestUpdateBookingStatusUseCase_Execute_WithMetadata(t *testing.T) {
	repoMock := new(mockUpdateStatusRepo)
	eventsMock := new(mockEventPublisher)
	uc := NewUpdateBookingStatusUseCase(repoMock, eventsMock)

	b := &booking.Booking{BookingID: "BK-1", CustomerID: "user-1", ProviderID: "prov-1", Status: booking.BookingStatusEnRoute}
	repoMock.On("FindByID", mock.Anything, "BK-1").Return(b, nil)
	repoMock.On("UpdateStatus", mock.Anything, "BK-1", booking.BookingStatusArrived, mock.MatchedBy(func(m map[string]any) bool {
		_, hasChangedBy := m["changedBy"]
		_, hasNote := m["note"]
		return hasChangedBy && hasNote
	})).Return(nil)
	eventsMock.On("Publish", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)

	result, err := uc.Execute(context.Background(), UpdateBookingStatusInput{
		BookingID:  "BK-1",
		CallerID:   "prov-1",
		CallerType: "provider",
		NewStatus:  booking.BookingStatusArrived,
		Metadata:   map[string]any{"note": "arrived at pickup"},
	})

	require.NoError(t, err)
	assert.Equal(t, "ARRIVED", result.Status)
	repoMock.AssertExpectations(t)
}

func TestUpdateBookingStatusUseCase_Execute_CompletionNoProviderSkipsCompletedEvent(t *testing.T) {
	repoMock := new(mockUpdateStatusRepo)
	eventsMock := new(mockEventPublisher)
	uc := NewUpdateBookingStatusUseCase(repoMock, eventsMock)

	// Booking without provider (edge case)
	b := &booking.Booking{BookingID: "BK-1", CustomerID: "user-1", ProviderID: "", Status: booking.BookingStatusOTPDropoff}
	repoMock.On("FindByID", mock.Anything, "BK-1").Return(b, nil)
	repoMock.On("UpdateStatus", mock.Anything, "BK-1", booking.BookingStatusCompleted, mock.Anything).Return(nil)
	eventsMock.On("Publish", mock.Anything, eventSourceBooking, eventBookingStatusChanged, mock.Anything, mock.Anything).Return(nil)

	_, err := uc.Execute(context.Background(), UpdateBookingStatusInput{
		BookingID: "BK-1", CallerID: "admin-1", CallerType: "admin", NewStatus: booking.BookingStatusCompleted,
	})

	require.NoError(t, err)
	// Only BookingStatusChanged, not BookingCompleted
	eventsMock.AssertNumberOfCalls(t, "Publish", 1)
}
