package websocket

import (
	"context"
	"log/slog"
)

// DisconnectUseCase handles the $disconnect WebSocket route.
// It uses the reverse mapping (connectionId -> userId) to clean up both
// forward and reverse entries, fixing the TS bug where disconnect could
// not determine which user to remove.
type DisconnectUseCase struct {
	sessions ConnectionRemover
	logger   *slog.Logger
}

// NewDisconnectUseCase creates a new DisconnectUseCase.
func NewDisconnectUseCase(sessions ConnectionRemover, logger *slog.Logger) *DisconnectUseCase {
	return &DisconnectUseCase{sessions: sessions, logger: logger}
}

// Execute performs the reverse lookup and cleans up both mappings.
func (uc *DisconnectUseCase) Execute(ctx context.Context, connectionID string) error {
	userID, err := uc.sessions.GetUserByConnection(ctx, connectionID)
	if err != nil {
		return err
	}

	// If the reverse mapping is already gone (TTL expired or duplicate disconnect),
	// we can only clean up the reverse key.
	if userID != "" {
		if err := uc.sessions.RemoveConnection(ctx, userID); err != nil {
			uc.logger.WarnContext(ctx, "failed to remove forward connection mapping",
				slog.String("userId", userID),
				slog.String("error", err.Error()),
			)
		}
	}

	if err := uc.sessions.RemoveReverseConnection(ctx, connectionID); err != nil {
		uc.logger.WarnContext(ctx, "failed to remove reverse connection mapping",
			slog.String("connectionId", connectionID),
			slog.String("error", err.Error()),
		)
	}

	uc.logger.InfoContext(ctx, "WebSocket disconnected",
		slog.String("userId", userID),
		slog.String("connectionId", connectionID),
	)
	return nil
}
