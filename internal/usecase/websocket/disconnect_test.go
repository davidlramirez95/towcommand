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

// --- Mocks ---

type mockConnectionRemover struct{ mock.Mock }

func (m *mockConnectionRemover) GetUserByConnection(ctx context.Context, connectionID string) (string, error) {
	args := m.Called(ctx, connectionID)
	return args.String(0), args.Error(1)
}

func (m *mockConnectionRemover) RemoveConnection(ctx context.Context, userID string) error {
	args := m.Called(ctx, userID)
	return args.Error(0)
}

func (m *mockConnectionRemover) RemoveReverseConnection(ctx context.Context, connectionID string) error {
	args := m.Called(ctx, connectionID)
	return args.Error(0)
}

// --- Tests ---

func TestDisconnectUseCase_Execute_Success(t *testing.T) {
	remover := new(mockConnectionRemover)
	logger := slog.Default()
	uc := NewDisconnectUseCase(remover, logger)

	remover.On("GetUserByConnection", mock.Anything, "conn-abc").Return("user-1", nil)
	remover.On("RemoveConnection", mock.Anything, "user-1").Return(nil)
	remover.On("RemoveReverseConnection", mock.Anything, "conn-abc").Return(nil)

	err := uc.Execute(context.Background(), "conn-abc")

	require.NoError(t, err)
	remover.AssertExpectations(t)
}

func TestDisconnectUseCase_Execute_UserNotFound(t *testing.T) {
	remover := new(mockConnectionRemover)
	logger := slog.Default()
	uc := NewDisconnectUseCase(remover, logger)

	// Reverse mapping already expired.
	remover.On("GetUserByConnection", mock.Anything, "conn-gone").Return("", nil)
	remover.On("RemoveReverseConnection", mock.Anything, "conn-gone").Return(nil)

	err := uc.Execute(context.Background(), "conn-gone")

	require.NoError(t, err)
	remover.AssertNotCalled(t, "RemoveConnection")
	remover.AssertExpectations(t)
}

func TestDisconnectUseCase_Execute_ReverseLookupError(t *testing.T) {
	remover := new(mockConnectionRemover)
	logger := slog.Default()
	uc := NewDisconnectUseCase(remover, logger)

	remover.On("GetUserByConnection", mock.Anything, "conn-err").
		Return("", domainerrors.NewExternalServiceError("Redis", nil))

	err := uc.Execute(context.Background(), "conn-err")

	assert.Error(t, err)
	remover.AssertNotCalled(t, "RemoveConnection")
	remover.AssertNotCalled(t, "RemoveReverseConnection")
}

func TestDisconnectUseCase_Execute_RemoveConnectionError_NonFatal(t *testing.T) {
	remover := new(mockConnectionRemover)
	logger := slog.Default()
	uc := NewDisconnectUseCase(remover, logger)

	remover.On("GetUserByConnection", mock.Anything, "conn-abc").Return("user-1", nil)
	remover.On("RemoveConnection", mock.Anything, "user-1").
		Return(domainerrors.NewExternalServiceError("Redis", nil))
	remover.On("RemoveReverseConnection", mock.Anything, "conn-abc").Return(nil)

	// RemoveConnection error is logged but does not fail the use case.
	err := uc.Execute(context.Background(), "conn-abc")

	require.NoError(t, err)
	remover.AssertExpectations(t)
}
