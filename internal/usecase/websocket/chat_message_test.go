package websocket

import (
	"context"
	"log/slog"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/davidlramirez95/towcommand/internal/domain/booking"
	domainerrors "github.com/davidlramirez95/towcommand/internal/domain/errors"
)

// --- Mocks ---

type mockChatSaver struct{ mock.Mock }

func (m *mockChatSaver) Save(ctx context.Context, msg *booking.ChatMessage) error {
	args := m.Called(ctx, msg)
	return args.Error(0)
}

type mockConnectionLookup struct{ mock.Mock }

func (m *mockConnectionLookup) GetConnection(ctx context.Context, userID string) (string, error) {
	args := m.Called(ctx, userID)
	return args.String(0), args.Error(1)
}

type mockConnectionPoster struct{ mock.Mock }

func (m *mockConnectionPoster) PostToConnection(ctx context.Context, connectionID string, data any) error {
	args := m.Called(ctx, connectionID, data)
	return args.Error(0)
}

// --- Tests ---

func TestChatMessageUseCase_Execute_Success(t *testing.T) {
	chat := new(mockChatSaver)
	sessions := new(mockConnectionLookup)
	poster := new(mockConnectionPoster)
	logger := slog.Default()

	uc := NewChatMessageUseCase(chat, sessions, poster, logger)
	uc.idGen = func() string { return "MSG-TEST123" }
	fixedTime := time.Date(2026, 3, 4, 10, 0, 0, 0, time.UTC)
	uc.now = func() time.Time { return fixedTime }

	input := ChatMessageInput{
		BookingID: "booking-1",
		Message:   "On my way!",
		SenderID:  "provider-1",
	}

	chat.On("Save", mock.Anything, mock.MatchedBy(func(msg *booking.ChatMessage) bool {
		return msg.MessageID == "MSG-TEST123" &&
			msg.BookingID == "booking-1" &&
			msg.SenderID == "provider-1" &&
			msg.Message == "On my way!"
	})).Return(nil)
	sessions.On("GetConnection", mock.Anything, "customer-1").Return("conn-abc", nil)
	poster.On("PostToConnection", mock.Anything, "conn-abc", mock.Anything).Return(nil)

	err := uc.Execute(context.Background(), input, "customer-1")

	require.NoError(t, err)
	chat.AssertExpectations(t)
	sessions.AssertExpectations(t)
	poster.AssertExpectations(t)
}

func TestChatMessageUseCase_Execute_RecipientOffline(t *testing.T) {
	chat := new(mockChatSaver)
	sessions := new(mockConnectionLookup)
	poster := new(mockConnectionPoster)
	logger := slog.Default()

	uc := NewChatMessageUseCase(chat, sessions, poster, logger)
	uc.idGen = func() string { return "MSG-OFFLINE" }
	uc.now = func() time.Time { return time.Now().UTC() }

	input := ChatMessageInput{
		BookingID: "booking-1",
		Message:   "Hello?",
		SenderID:  "user-1",
	}

	chat.On("Save", mock.Anything, mock.Anything).Return(nil)
	sessions.On("GetConnection", mock.Anything, "user-2").Return("", nil) // Not connected.

	err := uc.Execute(context.Background(), input, "user-2")

	require.NoError(t, err)
	chat.AssertExpectations(t)
	sessions.AssertExpectations(t)
	poster.AssertNotCalled(t, "PostToConnection")
}

func TestChatMessageUseCase_Execute_SaveError(t *testing.T) {
	chat := new(mockChatSaver)
	sessions := new(mockConnectionLookup)
	poster := new(mockConnectionPoster)
	logger := slog.Default()

	uc := NewChatMessageUseCase(chat, sessions, poster, logger)
	uc.idGen = func() string { return "MSG-ERR" }
	uc.now = func() time.Time { return time.Now().UTC() }

	input := ChatMessageInput{
		BookingID: "booking-1",
		Message:   "fail",
		SenderID:  "user-1",
	}

	chat.On("Save", mock.Anything, mock.Anything).
		Return(domainerrors.NewInternalError("db error"))

	err := uc.Execute(context.Background(), input, "user-2")

	assert.Error(t, err)
	sessions.AssertNotCalled(t, "GetConnection")
	poster.AssertNotCalled(t, "PostToConnection")
}

func TestChatMessageUseCase_Execute_PostError_NonFatal(t *testing.T) {
	chat := new(mockChatSaver)
	sessions := new(mockConnectionLookup)
	poster := new(mockConnectionPoster)
	logger := slog.Default()

	uc := NewChatMessageUseCase(chat, sessions, poster, logger)
	uc.idGen = func() string { return "MSG-POSTERR" }
	uc.now = func() time.Time { return time.Now().UTC() }

	input := ChatMessageInput{
		BookingID: "booking-1",
		Message:   "hi",
		SenderID:  "user-1",
	}

	chat.On("Save", mock.Anything, mock.Anything).Return(nil)
	sessions.On("GetConnection", mock.Anything, "user-2").Return("conn-xyz", nil)
	poster.On("PostToConnection", mock.Anything, "conn-xyz", mock.Anything).
		Return(domainerrors.NewExternalServiceError("apigateway", nil))

	// Post error is non-fatal; message was saved.
	err := uc.Execute(context.Background(), input, "user-2")

	require.NoError(t, err)
	chat.AssertExpectations(t)
	sessions.AssertExpectations(t)
	poster.AssertExpectations(t)
}
