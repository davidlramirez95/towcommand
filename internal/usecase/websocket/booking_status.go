package websocket

import (
	"context"
	"log/slog"
)

// BookingStatusInput carries the validated fields for a booking status WebSocket push.
type BookingStatusInput struct {
	BookingID string `json:"bookingId" validate:"required"`
	Status    string `json:"status" validate:"required"`
	UserID    string `json:"userId" validate:"required"`
}

// BookingStatusUseCase pushes booking status updates to connected users via WebSocket.
type BookingStatusUseCase struct {
	sessions ConnectionLookup
	poster   ConnectionPoster
	logger   *slog.Logger
}

// NewBookingStatusUseCase creates a new BookingStatusUseCase.
func NewBookingStatusUseCase(
	sessions ConnectionLookup,
	poster ConnectionPoster,
	logger *slog.Logger,
) *BookingStatusUseCase {
	return &BookingStatusUseCase{sessions: sessions, poster: poster, logger: logger}
}

// Execute looks up the user's connection and sends the status update.
func (uc *BookingStatusUseCase) Execute(ctx context.Context, input BookingStatusInput) error {
	connID, err := uc.sessions.GetConnection(ctx, input.UserID)
	if err != nil {
		return err
	}

	if connID == "" {
		uc.logger.DebugContext(ctx, "user not connected, skipping status push",
			slog.String("userId", input.UserID),
			slog.String("bookingId", input.BookingID),
		)
		return nil
	}

	if err := uc.poster.PostToConnection(ctx, connID, map[string]any{
		"action":    "bookingStatus",
		"bookingId": input.BookingID,
		"status":    input.Status,
	}); err != nil {
		uc.logger.WarnContext(ctx, "failed to push booking status via WebSocket",
			slog.String("userId", input.UserID),
			slog.String("connectionId", connID),
			slog.String("error", err.Error()),
		)
	}

	return nil
}
