package websocket

import (
	"context"
	"log/slog"
	"time"
)

const connectionTTL = 24 * time.Hour

// ConnectUseCase handles the $connect WebSocket route.
// It stores both forward (userId -> connectionId) and reverse
// (connectionId -> userId) mappings to fix the TS disconnect bug.
type ConnectUseCase struct {
	sessions ConnectionMapper
	logger   *slog.Logger
}

// NewConnectUseCase creates a new ConnectUseCase.
func NewConnectUseCase(sessions ConnectionMapper, logger *slog.Logger) *ConnectUseCase {
	return &ConnectUseCase{sessions: sessions, logger: logger}
}

// Execute stores the bidirectional connection mapping.
func (uc *ConnectUseCase) Execute(ctx context.Context, userID, connectionID string) error {
	if err := uc.sessions.MapConnection(ctx, userID, connectionID, connectionTTL); err != nil {
		return err
	}
	if err := uc.sessions.MapReverseConnection(ctx, connectionID, userID, connectionTTL); err != nil {
		return err
	}

	uc.logger.InfoContext(ctx, "WebSocket connected",
		slog.String("userId", userID),
		slog.String("connectionId", connectionID),
	)
	return nil
}
