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

type mockBookingFinder struct{ mock.Mock }

func (m *mockBookingFinder) FindByID(ctx context.Context, bookingID string) (*booking.Booking, error) {
	args := m.Called(ctx, bookingID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*booking.Booking), args.Error(1)
}

// --- Tests ---

func TestGetBookingUseCase_Execute_OwnerAccess(t *testing.T) {
	repoMock := new(mockBookingFinder)
	uc := NewGetBookingUseCase(repoMock)

	b := &booking.Booking{BookingID: "BK-1", CustomerID: "user-1", Status: booking.BookingStatusPending}
	repoMock.On("FindByID", mock.Anything, "BK-1").Return(b, nil)

	result, err := uc.Execute(context.Background(), GetBookingInput{
		BookingID: "BK-1", CallerID: "user-1", CallerType: "customer",
	})

	require.NoError(t, err)
	assert.Equal(t, "BK-1", result.BookingID)
	repoMock.AssertExpectations(t)
}

func TestGetBookingUseCase_Execute_ProviderAccess(t *testing.T) {
	repoMock := new(mockBookingFinder)
	uc := NewGetBookingUseCase(repoMock)

	b := &booking.Booking{BookingID: "BK-1", CustomerID: "user-1", ProviderID: "prov-1", Status: booking.BookingStatusMatched}
	repoMock.On("FindByID", mock.Anything, "BK-1").Return(b, nil)

	result, err := uc.Execute(context.Background(), GetBookingInput{
		BookingID: "BK-1", CallerID: "prov-1", CallerType: "provider",
	})

	require.NoError(t, err)
	assert.Equal(t, "BK-1", result.BookingID)
}

func TestGetBookingUseCase_Execute_AdminAccess(t *testing.T) {
	repoMock := new(mockBookingFinder)
	uc := NewGetBookingUseCase(repoMock)

	b := &booking.Booking{BookingID: "BK-1", CustomerID: "user-1", Status: booking.BookingStatusPending}
	repoMock.On("FindByID", mock.Anything, "BK-1").Return(b, nil)

	for _, callerType := range []string{"admin", "ops_agent"} {
		result, err := uc.Execute(context.Background(), GetBookingInput{
			BookingID: "BK-1", CallerID: "admin-1", CallerType: callerType,
		})
		require.NoError(t, err)
		assert.Equal(t, "BK-1", result.BookingID)
	}
}

func TestGetBookingUseCase_Execute_NotFound(t *testing.T) {
	repoMock := new(mockBookingFinder)
	uc := NewGetBookingUseCase(repoMock)

	repoMock.On("FindByID", mock.Anything, "BK-GONE").Return(nil, nil)

	_, err := uc.Execute(context.Background(), GetBookingInput{
		BookingID: "BK-GONE", CallerID: "user-1", CallerType: "customer",
	})

	require.Error(t, err)
	var appErr *domainerrors.AppError
	require.True(t, errors.As(err, &appErr))
	assert.Equal(t, domainerrors.CodeNotFound, appErr.Code)
}

func TestGetBookingUseCase_Execute_Forbidden(t *testing.T) {
	repoMock := new(mockBookingFinder)
	uc := NewGetBookingUseCase(repoMock)

	b := &booking.Booking{BookingID: "BK-1", CustomerID: "user-1", ProviderID: "prov-1", Status: booking.BookingStatusMatched}
	repoMock.On("FindByID", mock.Anything, "BK-1").Return(b, nil)

	_, err := uc.Execute(context.Background(), GetBookingInput{
		BookingID: "BK-1", CallerID: "user-OTHER", CallerType: "customer",
	})

	require.Error(t, err)
	var appErr *domainerrors.AppError
	require.True(t, errors.As(err, &appErr))
	assert.Equal(t, domainerrors.CodeForbidden, appErr.Code)
}

func TestGetBookingUseCase_Execute_RepoError(t *testing.T) {
	repoMock := new(mockBookingFinder)
	uc := NewGetBookingUseCase(repoMock)

	repoMock.On("FindByID", mock.Anything, "BK-ERR").Return(nil, domainerrors.NewInternalError("db error"))

	_, err := uc.Execute(context.Background(), GetBookingInput{
		BookingID: "BK-ERR", CallerID: "user-1", CallerType: "customer",
	})

	assert.Error(t, err)
}
