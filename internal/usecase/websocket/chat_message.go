package websocket

import (
	"context"
	"crypto/rand"
	"fmt"
	"log/slog"
	"time"

	"github.com/davidlramirez95/towcommand/internal/domain/booking"
)

// ChatMessageInput carries the validated fields from a WebSocket chat message.
type ChatMessageInput struct {
	BookingID string `json:"bookingId" validate:"required"`
	Message   string `json:"message" validate:"required,max=1000"`
	SenderID  string `json:"-"` // Extracted from auth context, not from body.
}

// ChatMessageUseCase handles chat messages sent during an active booking.
type ChatMessageUseCase struct {
	chat     ChatSaver
	sessions ConnectionLookup
	poster   ConnectionPoster
	logger   *slog.Logger
	idGen    func() string
	now      func() time.Time
}

// NewChatMessageUseCase creates a new ChatMessageUseCase.
func NewChatMessageUseCase(
	chat ChatSaver,
	sessions ConnectionLookup,
	poster ConnectionPoster,
	logger *slog.Logger,
) *ChatMessageUseCase {
	return &ChatMessageUseCase{
		chat:     chat,
		sessions: sessions,
		poster:   poster,
		logger:   logger,
		idGen:    generateMessageID,
		now:      func() time.Time { return time.Now().UTC() },
	}
}

// Execute saves a chat message and delivers it to the recipient via WebSocket.
func (uc *ChatMessageUseCase) Execute(ctx context.Context, input ChatMessageInput, recipientID string) error {
	msg := &booking.ChatMessage{
		MessageID: uc.idGen(),
		BookingID: input.BookingID,
		SenderID:  input.SenderID,
		Message:   input.Message,
		CreatedAt: uc.now(),
	}

	if err := uc.chat.Save(ctx, msg); err != nil {
		return fmt.Errorf("saving chat message: %w", err)
	}

	// Look up the recipient's WebSocket connection.
	connID, err := uc.sessions.GetConnection(ctx, recipientID)
	if err != nil {
		uc.logger.WarnContext(ctx, "failed to look up recipient connection",
			slog.String("recipientId", recipientID),
			slog.String("error", err.Error()),
		)
		return nil // Message saved; delivery is best-effort.
	}

	if connID == "" {
		uc.logger.DebugContext(ctx, "recipient not connected, message saved for later",
			slog.String("recipientId", recipientID),
			slog.String("bookingId", input.BookingID),
		)
		return nil
	}

	// Deliver the message via WebSocket.
	if err := uc.poster.PostToConnection(ctx, connID, map[string]any{
		"action":    "chatMessage",
		"messageId": msg.MessageID,
		"bookingId": msg.BookingID,
		"senderId":  msg.SenderID,
		"message":   msg.Message,
		"createdAt": msg.CreatedAt.Format(time.RFC3339),
	}); err != nil {
		uc.logger.WarnContext(ctx, "failed to deliver chat message via WebSocket",
			slog.String("recipientId", recipientID),
			slog.String("connectionId", connID),
			slog.String("error", err.Error()),
		)
	}

	return nil
}

// generateMessageID produces a random message ID.
func generateMessageID() string {
	b := make([]byte, 12)
	_, _ = rand.Read(b)
	return fmt.Sprintf("MSG-%X", b)
}
