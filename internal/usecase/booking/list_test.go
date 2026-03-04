package bookinguc

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/davidlramirez95/towcommand/internal/domain/booking"
	domainerrors "github.com/davidlramirez95/towcommand/internal/domain/errors"
)

// --- Mocks ---

type mockBookingLister struct{ mock.Mock }

func (m *mockBookingLister) FindByUser(ctx context.Context, userID string, limit int32) ([]booking.Booking, error) {
	args := m.Called(ctx, userID, limit)
	return args.Get(0).([]booking.Booking), args.Error(1)
}

func (m *mockBookingLister) FindByStatus(ctx context.Context, status booking.BookingStatus, limit int32) ([]booking.Booking, error) {
	args := m.Called(ctx, status, limit)
	return args.Get(0).([]booking.Booking), args.Error(1)
}

// --- Tests ---

func TestListBookingsUseCase_Execute_CustomerListsOwnBookings(t *testing.T) {
	repoMock := new(mockBookingLister)
	uc := NewListBookingsUseCase(repoMock)

	bookings := []booking.Booking{
		{BookingID: "BK-1", CustomerID: "user-1", Status: booking.BookingStatusPending},
		{BookingID: "BK-2", CustomerID: "user-1", Status: booking.BookingStatusCompleted},
	}
	repoMock.On("FindByUser", mock.Anything, "user-1", int32(25)).Return(bookings, nil)

	result, err := uc.Execute(context.Background(), ListBookingsInput{
		CallerID: "user-1", CallerType: "customer", Limit: 25,
	})

	require.NoError(t, err)
	assert.Len(t, result.Items, 2)
	assert.Equal(t, 2, result.Count)
	repoMock.AssertExpectations(t)
}

func TestListBookingsUseCase_Execute_AdminFiltersByStatus(t *testing.T) {
	repoMock := new(mockBookingLister)
	uc := NewListBookingsUseCase(repoMock)

	bookings := []booking.Booking{
		{BookingID: "BK-1", Status: booking.BookingStatusPending},
	}
	repoMock.On("FindByStatus", mock.Anything, booking.BookingStatusPending, int32(25)).Return(bookings, nil)

	result, err := uc.Execute(context.Background(), ListBookingsInput{
		CallerID: "admin-1", CallerType: "admin", Limit: 25, StatusFilter: "PENDING",
	})

	require.NoError(t, err)
	assert.Len(t, result.Items, 1)
	repoMock.AssertExpectations(t)
}

func TestListBookingsUseCase_Execute_OpsAgentFiltersByStatus(t *testing.T) {
	repoMock := new(mockBookingLister)
	uc := NewListBookingsUseCase(repoMock)

	bookings := []booking.Booking{
		{BookingID: "BK-1", Status: booking.BookingStatusMatched},
	}
	repoMock.On("FindByStatus", mock.Anything, booking.BookingStatusMatched, int32(25)).Return(bookings, nil)

	result, err := uc.Execute(context.Background(), ListBookingsInput{
		CallerID: "ops-1", CallerType: "ops_agent", Limit: 25, StatusFilter: "MATCHED",
	})

	require.NoError(t, err)
	assert.Len(t, result.Items, 1)
}

func TestListBookingsUseCase_Execute_CustomerStatusFilterClientSide(t *testing.T) {
	repoMock := new(mockBookingLister)
	uc := NewListBookingsUseCase(repoMock)

	bookings := []booking.Booking{
		{BookingID: "BK-1", CustomerID: "user-1", Status: booking.BookingStatusPending},
		{BookingID: "BK-2", CustomerID: "user-1", Status: booking.BookingStatusCompleted},
		{BookingID: "BK-3", CustomerID: "user-1", Status: booking.BookingStatusPending},
	}
	repoMock.On("FindByUser", mock.Anything, "user-1", int32(25)).Return(bookings, nil)

	result, err := uc.Execute(context.Background(), ListBookingsInput{
		CallerID: "user-1", CallerType: "customer", Limit: 25, StatusFilter: "PENDING",
	})

	require.NoError(t, err)
	assert.Len(t, result.Items, 2)
	assert.Equal(t, 2, result.Count)
	for _, b := range result.Items {
		assert.Equal(t, booking.BookingStatusPending, b.Status)
	}
}

func TestListBookingsUseCase_Execute_DefaultAndMaxLimit(t *testing.T) {
	tests := []struct {
		name      string
		inputLim  int32
		expectLim int32
	}{
		{"zero becomes 25", 0, 25},
		{"negative becomes 25", -1, 25},
		{"over 100 capped to 100", 200, 100},
		{"normal passes through", 50, 50},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repoMock := new(mockBookingLister)
			uc := NewListBookingsUseCase(repoMock)

			repoMock.On("FindByUser", mock.Anything, "u1", tt.expectLim).Return([]booking.Booking{}, nil)

			_, err := uc.Execute(context.Background(), ListBookingsInput{
				CallerID: "u1", CallerType: "customer", Limit: tt.inputLim,
			})
			require.NoError(t, err)
			repoMock.AssertExpectations(t)
		})
	}
}

func TestListBookingsUseCase_Execute_RepoError(t *testing.T) {
	repoMock := new(mockBookingLister)
	uc := NewListBookingsUseCase(repoMock)

	repoMock.On("FindByUser", mock.Anything, "user-1", int32(25)).
		Return([]booking.Booking{}, domainerrors.NewInternalError("db error"))

	_, err := uc.Execute(context.Background(), ListBookingsInput{
		CallerID: "user-1", CallerType: "customer", Limit: 25,
	})

	assert.Error(t, err)
}
