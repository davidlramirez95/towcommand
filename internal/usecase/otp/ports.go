// Package otpuc implements OTP generation and verification use cases.
// Each use case declares only the port interfaces it needs (ISP).
package otpuc

import (
	"context"
	"time"

	"github.com/davidlramirez95/towcommand/internal/domain/booking"
	"github.com/davidlramirez95/towcommand/internal/domain/otp"
	"github.com/davidlramirez95/towcommand/internal/domain/user"
	"github.com/davidlramirez95/towcommand/internal/usecase/port"
)

// BookingFinder retrieves a booking by its ID. Returns nil if not found.
type BookingFinder interface {
	FindByID(ctx context.Context, bookingID string) (*booking.Booking, error)
}

// BookingStatusUpdater changes a booking's status and records the transition.
type BookingStatusUpdater interface {
	UpdateStatus(ctx context.Context, bookingID string, status booking.BookingStatus, metadata map[string]any) error
}

// UserFinder retrieves a user by their ID. Returns nil if not found.
type UserFinder interface {
	FindByID(ctx context.Context, userID string) (*user.User, error)
}

// OTPCache stores hashed OTP codes with automatic expiration.
type OTPCache interface {
	StoreOTP(ctx context.Context, bookingID, otpType, hashedOTP string, ttl time.Duration) error
	GetOTP(ctx context.Context, bookingID, otpType string) (string, error)
	DeleteOTP(ctx context.Context, bookingID, otpType string) error
}

// RateLimiter performs sliding-window rate-limit checks.
type RateLimiter interface {
	CheckRateLimit(ctx context.Context, key string, maxRequests, windowSec int) (allowed bool, remaining int, err error)
}

// OTPRepo persists OTP records as backup for Redis.
type OTPRepo interface {
	Save(ctx context.Context, o *otp.OTP) error
	FindByBookingAndType(ctx context.Context, bookingID, otpType string) (*otp.OTP, error)
	MarkVerified(ctx context.Context, bookingID, otpType string) error
}

// SMSSender sends SMS messages to phone numbers.
type SMSSender interface {
	SendSMS(ctx context.Context, phoneNumber, message string) error
}

// EventPublisher publishes domain events to an event bus.
type EventPublisher interface {
	Publish(ctx context.Context, source, detailType string, detail any, actor *port.Actor) error
}

// Actor is a type alias for port.Actor to avoid import stuttering in use case code.
type Actor = port.Actor
