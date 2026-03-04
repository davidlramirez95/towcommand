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

type mockCancelBookingRepo struct{ mock.Mock }

func (m *mockCancelBookingRepo) FindByID(ctx context.Context, bookingID string) (*booking.Booking, error) {
	args := m.Called(ctx, bookingID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*booking.Booking), args.Error(1)
}

func (m *mockCancelBookingRepo) UpdateStatus(ctx context.Context, bookingID string, status booking.BookingStatus, metadata map[string]any) error {
	args := m.Called(ctx, bookingID, status, metadata)
	return args.Error(0)
}

// --- Tests ---

func TestCancelBookingUseCase_Execute_PendingSuccess(t *testing.T) {
	repoMock := new(mockCancelBookingRepo)
	eventsMock := new(mockEventPublisher)
	uc := NewCancelBookingUseCase(repoMock, eventsMock)

	b := &booking.Booking{BookingID: "BK-1", CustomerID: "user-1", Status: booking.BookingStatusPending}
	repoMock.On("FindByID", mock.Anything, "BK-1").Return(b, nil)
	repoMock.On("UpdateStatus", mock.Anything, "BK-1", booking.BookingStatusCancelled, mock.Anything).Return(nil)
	eventsMock.On("Publish", mock.Anything, eventSourceBooking, eventBookingCancelled, mock.Anything, mock.Anything).Return(nil)

	result, err := uc.Execute(context.Background(), CancelBookingInput{
		BookingID: "BK-1", CallerID: "user-1", Reason: "changed my mind",
	})

	require.NoError(t, err)
	assert.Equal(t, "BK-1", result.BookingID)
	assert.Equal(t, string(booking.BookingStatusCancelled), result.Status)
	assert.Equal(t, int64(0), result.CancellationFee) // PENDING = no fee
	repoMock.AssertExpectations(t)
	eventsMock.AssertExpectations(t)
}

func TestCancelBookingUseCase_Execute_MatchedWithFee(t *testing.T) {
	repoMock := new(mockCancelBookingRepo)
	eventsMock := new(mockEventPublisher)
	uc := NewCancelBookingUseCase(repoMock, eventsMock)

	b := &booking.Booking{BookingID: "BK-2", CustomerID: "user-1", ProviderID: "prov-1", Status: booking.BookingStatusMatched}
	repoMock.On("FindByID", mock.Anything, "BK-2").Return(b, nil)
	repoMock.On("UpdateStatus", mock.Anything, "BK-2", booking.BookingStatusCancelled, mock.Anything).Return(nil)
	eventsMock.On("Publish", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)

	result, err := uc.Execute(context.Background(), CancelBookingInput{
		BookingID: "BK-2", CallerID: "user-1",
	})

	require.NoError(t, err)
	assert.Equal(t, int64(10_000), result.CancellationFee) // MATCHED = ₱100
}

func TestCancelBookingUseCase_Execute_EnRouteWithHigherFee(t *testing.T) {
	repoMock := new(mockCancelBookingRepo)
	eventsMock := new(mockEventPublisher)
	uc := NewCancelBookingUseCase(repoMock, eventsMock)

	b := &booking.Booking{BookingID: "BK-3", CustomerID: "user-1", Status: booking.BookingStatusEnRoute}
	repoMock.On("FindByID", mock.Anything, "BK-3").Return(b, nil)
	repoMock.On("UpdateStatus", mock.Anything, "BK-3", booking.BookingStatusCancelled, mock.Anything).Return(nil)
	eventsMock.On("Publish", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)

	result, err := uc.Execute(context.Background(), CancelBookingInput{
		BookingID: "BK-3", CallerID: "user-1",
	})

	require.NoError(t, err)
	assert.Equal(t, int64(25_000), result.CancellationFee) // EN_ROUTE = ₱250
}

func TestCancelBookingUseCase_Execute_NotFound(t *testing.T) {
	repoMock := new(mockCancelBookingRepo)
	eventsMock := new(mockEventPublisher)
	uc := NewCancelBookingUseCase(repoMock, eventsMock)

	repoMock.On("FindByID", mock.Anything, "BK-GONE").Return(nil, nil)

	_, err := uc.Execute(context.Background(), CancelBookingInput{
		BookingID: "BK-GONE", CallerID: "user-1",
	})

	require.Error(t, err)
	var appErr *domainerrors.AppError
	require.True(t, errors.As(err, &appErr))
	assert.Equal(t, domainerrors.CodeNotFound, appErr.Code)
}

func TestCancelBookingUseCase_Execute_NotOwner(t *testing.T) {
	repoMock := new(mockCancelBookingRepo)
	eventsMock := new(mockEventPublisher)
	uc := NewCancelBookingUseCase(repoMock, eventsMock)

	b := &booking.Booking{BookingID: "BK-1", CustomerID: "user-1", Status: booking.BookingStatusPending}
	repoMock.On("FindByID", mock.Anything, "BK-1").Return(b, nil)

	_, err := uc.Execute(context.Background(), CancelBookingInput{
		BookingID: "BK-1", CallerID: "user-HACKER",
	})

	require.Error(t, err)
	var appErr *domainerrors.AppError
	require.True(t, errors.As(err, &appErr))
	assert.Equal(t, domainerrors.CodeForbidden, appErr.Code)
}

func TestCancelBookingUseCase_Execute_InvalidTransition(t *testing.T) {
	nonCancellableStatuses := []booking.BookingStatus{
		booking.BookingStatusArrived,
		booking.BookingStatusConditionReport,
		booking.BookingStatusOTPVerified,
		booking.BookingStatusLoading,
		booking.BookingStatusInTransit,
		booking.BookingStatusArrivedDropoff,
		booking.BookingStatusOTPDropoff,
		booking.BookingStatusCompleted,
		booking.BookingStatusCancelled,
	}

	for _, status := range nonCancellableStatuses {
		t.Run(string(status), func(t *testing.T) {
			repoMock := new(mockCancelBookingRepo)
			eventsMock := new(mockEventPublisher)
			uc := NewCancelBookingUseCase(repoMock, eventsMock)

			b := &booking.Booking{BookingID: "BK-X", CustomerID: "user-1", Status: status}
			repoMock.On("FindByID", mock.Anything, "BK-X").Return(b, nil)

			_, err := uc.Execute(context.Background(), CancelBookingInput{
				BookingID: "BK-X", CallerID: "user-1",
			})

			require.Error(t, err)
			var appErr *domainerrors.AppError
			require.True(t, errors.As(err, &appErr))
			assert.Equal(t, domainerrors.CodeBookingNotCancellable, appErr.Code)
		})
	}
}

func TestCancelBookingUseCase_Execute_UpdateStatusError(t *testing.T) {
	repoMock := new(mockCancelBookingRepo)
	eventsMock := new(mockEventPublisher)
	uc := NewCancelBookingUseCase(repoMock, eventsMock)

	b := &booking.Booking{BookingID: "BK-1", CustomerID: "user-1", Status: booking.BookingStatusPending}
	repoMock.On("FindByID", mock.Anything, "BK-1").Return(b, nil)
	repoMock.On("UpdateStatus", mock.Anything, "BK-1", booking.BookingStatusCancelled, mock.Anything).
		Return(domainerrors.NewInternalError("db error"))

	_, err := uc.Execute(context.Background(), CancelBookingInput{
		BookingID: "BK-1", CallerID: "user-1",
	})

	assert.Error(t, err)
}

func TestCalculateCancellationFee(t *testing.T) {
	tests := []struct {
		status booking.BookingStatus
		want   int64
	}{
		{booking.BookingStatusPending, 0},
		{booking.BookingStatusMatched, 10_000},
		{booking.BookingStatusEnRoute, 25_000},
		{booking.BookingStatusArrived, 0},
	}
	for _, tt := range tests {
		t.Run(string(tt.status), func(t *testing.T) {
			assert.Equal(t, tt.want, calculateCancellationFee(tt.status))
		})
	}
}
