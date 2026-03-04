package websocket

import (
	"context"
	"log/slog"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	domainerrors "github.com/davidlramirez95/towcommand/internal/domain/errors"
)

// --- Mocks ---

type mockConnectionMapper struct{ mock.Mock }

func (m *mockConnectionMapper) MapConnection(ctx context.Context, userID, connectionID string, ttl time.Duration) error {
	args := m.Called(ctx, userID, connectionID, ttl)
	return args.Error(0)
}

func (m *mockConnectionMapper) MapReverseConnection(ctx context.Context, connectionID, userID string, ttl time.Duration) error {
	args := m.Called(ctx, connectionID, userID, ttl)
	return args.Error(0)
}

// --- Tests ---

func TestConnectUseCase_Execute_Success(t *testing.T) {
	mapper := new(mockConnectionMapper)
	logger := slog.Default()
	uc := NewConnectUseCase(mapper, logger)

	mapper.On("MapConnection", mock.Anything, "user-1", "conn-abc", connectionTTL).Return(nil)
	mapper.On("MapReverseConnection", mock.Anything, "conn-abc", "user-1", connectionTTL).Return(nil)

	err := uc.Execute(context.Background(), "user-1", "conn-abc")

	require.NoError(t, err)
	mapper.AssertExpectations(t)
}

func TestConnectUseCase_Execute_MapConnectionError(t *testing.T) {
	mapper := new(mockConnectionMapper)
	logger := slog.Default()
	uc := NewConnectUseCase(mapper, logger)

	mapper.On("MapConnection", mock.Anything, "user-1", "conn-abc", connectionTTL).
		Return(domainerrors.NewExternalServiceError("Redis", nil))

	err := uc.Execute(context.Background(), "user-1", "conn-abc")

	assert.Error(t, err)
	mapper.AssertNotCalled(t, "MapReverseConnection")
}

func TestConnectUseCase_Execute_MapReverseConnectionError(t *testing.T) {
	mapper := new(mockConnectionMapper)
	logger := slog.Default()
	uc := NewConnectUseCase(mapper, logger)

	mapper.On("MapConnection", mock.Anything, "user-1", "conn-abc", connectionTTL).Return(nil)
	mapper.On("MapReverseConnection", mock.Anything, "conn-abc", "user-1", connectionTTL).
		Return(domainerrors.NewExternalServiceError("Redis", nil))

	err := uc.Execute(context.Background(), "user-1", "conn-abc")

	assert.Error(t, err)
	mapper.AssertExpectations(t)
}
