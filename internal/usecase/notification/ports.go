// Package notification implements notification routing and delivery for domain events.
// Each event type is dispatched to the appropriate channel (SMS, email) with
// Filipino/English message templates tailored for the Philippine market.
package notification

import (
	"context"

	"github.com/davidlramirez95/towcommand/internal/domain/booking"
	"github.com/davidlramirez95/towcommand/internal/domain/user"
)

// SMSSender sends SMS messages to phone numbers.
type SMSSender interface {
	SendSMS(ctx context.Context, phoneNumber, message string) error
}

// EmailSender sends email messages.
type EmailSender interface {
	SendEmail(ctx context.Context, to, subject, htmlBody string) error
}

// UserFinder retrieves a user by their ID.
type UserFinder interface {
	FindByID(ctx context.Context, userID string) (*user.User, error)
}

// BookingFinder retrieves a booking by its ID.
type BookingFinder interface {
	FindByID(ctx context.Context, bookingID string) (*booking.Booking, error)
}
