package websocket

import (
	"context"
	"log/slog"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	domainerrors "github.com/davidlramirez95/towcommand/internal/domain/errors"
)

// Mocks are shared from chat_message_test.go (mockConnectionLookup, mockConnectionPoster).

func TestBookingStatusUseCase_Execute_Success(t *testing.T) {
	sessions := new(mockConnectionLookup)
	poster := new(mockConnectionPoster)
	logger := slog.Default()
	uc := NewBookingStatusUseCase(sessions, poster, logger)

	input := BookingStatusInput{
		BookingID: "booking-1",
		Status:    "DRIVER_ASSIGNED",
		UserID:    "user-1",
	}

	sessions.On("GetConnection", mock.Anything, "user-1").Return("conn-abc", nil)
	poster.On("PostToConnection", mock.Anything, "conn-abc", mock.MatchedBy(func(data any) bool {
		m, ok := data.(map[string]any)
		if !ok {
			return false
		}
		return m["action"] == "bookingStatus" &&
			m["bookingId"] == "booking-1" &&
			m["status"] == "DRIVER_ASSIGNED"
	})).Return(nil)

	err := uc.Execute(context.Background(), input)

	require.NoError(t, err)
	sessions.AssertExpectations(t)
	poster.AssertExpectations(t)
}

func TestBookingStatusUseCase_Execute_UserNotConnected(t *testing.T) {
	sessions := new(mockConnectionLookup)
	poster := new(mockConnectionPoster)
	logger := slog.Default()
	uc := NewBookingStatusUseCase(sessions, poster, logger)

	input := BookingStatusInput{
		BookingID: "booking-1",
		Status:    "EN_ROUTE",
		UserID:    "user-offline",
	}

	sessions.On("GetConnection", mock.Anything, "user-offline").Return("", nil)

	err := uc.Execute(context.Background(), input)

	require.NoError(t, err)
	poster.AssertNotCalled(t, "PostToConnection")
}

func TestBookingStatusUseCase_Execute_SessionLookupError(t *testing.T) {
	sessions := new(mockConnectionLookup)
	poster := new(mockConnectionPoster)
	logger := slog.Default()
	uc := NewBookingStatusUseCase(sessions, poster, logger)

	input := BookingStatusInput{
		BookingID: "booking-1",
		Status:    "COMPLETED",
		UserID:    "user-1",
	}

	sessions.On("GetConnection", mock.Anything, "user-1").
		Return("", domainerrors.NewExternalServiceError("Redis", nil))

	err := uc.Execute(context.Background(), input)

	assert.Error(t, err)
	poster.AssertNotCalled(t, "PostToConnection")
}

func TestBookingStatusUseCase_Execute_PostError_NonFatal(t *testing.T) {
	sessions := new(mockConnectionLookup)
	poster := new(mockConnectionPoster)
	logger := slog.Default()
	uc := NewBookingStatusUseCase(sessions, poster, logger)

	input := BookingStatusInput{
		BookingID: "booking-1",
		Status:    "ARRIVED",
		UserID:    "user-1",
	}

	sessions.On("GetConnection", mock.Anything, "user-1").Return("conn-abc", nil)
	poster.On("PostToConnection", mock.Anything, "conn-abc", mock.Anything).
		Return(domainerrors.NewExternalServiceError("apigateway", nil))

	// Post error is logged but non-fatal.
	err := uc.Execute(context.Background(), input)

	require.NoError(t, err)
	sessions.AssertExpectations(t)
	poster.AssertExpectations(t)
}
