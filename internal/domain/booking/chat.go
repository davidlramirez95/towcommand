package booking

import "time"

// ChatMessage represents a message sent between customer and provider
// during an active booking.
type ChatMessage struct {
	MessageID string    `json:"messageId" validate:"required"`
	BookingID string    `json:"bookingId" validate:"required"`
	SenderID  string    `json:"senderId" validate:"required"`
	Message   string    `json:"message" validate:"required,max=1000"`
	CreatedAt time.Time `json:"createdAt"`
}
