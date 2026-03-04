package rating

import "time"

// Rating represents a customer's review of a completed booking.
type Rating struct {
	BookingID  string    `json:"bookingId" validate:"required"`
	CustomerID string    `json:"customerId" validate:"required"`
	ProviderID string    `json:"providerId" validate:"required"`
	Score      int       `json:"rating" validate:"required,min=1,max=5"`
	Comment    string    `json:"comment,omitempty"`
	Tags       []string  `json:"tags,omitempty"`
	CreatedAt  time.Time `json:"createdAt"`
}
