package otpuc

import (
	"context"
	"crypto/rand"
	"fmt"
	"log/slog"
	"time"

	"github.com/davidlramirez95/towcommand/internal/domain/booking"
	domainerrors "github.com/davidlramirez95/towcommand/internal/domain/errors"
	"github.com/davidlramirez95/towcommand/internal/domain/otp"
)

const (
	eventSourceOTP    = "tc.otp"
	eventOTPGenerated = "OTPGenerated"

	otpTTL               = 5 * time.Minute
	rateLimitMaxGenerate = 3
	rateLimitWindowSec   = 300
)

// GenerateOTPInput holds the data needed to generate an OTP.
type GenerateOTPInput struct {
	BookingID string
	OTPType   otp.OTPType
	Lat       float64
	Lng       float64
	CallerID  string
}

// GenerateOTPOutput is the response for OTP generation.
// The plaintext code is NEVER included in the API response.
type GenerateOTPOutput struct {
	Success   bool   `json:"success"`
	BookingID string `json:"bookingId"`
	OTPType   string `json:"otpType"`
}

// GenerateOTPUseCase orchestrates OTP generation with rate limiting and dual-store persistence.
type GenerateOTPUseCase struct {
	bookings    BookingFinder
	users       UserFinder
	otpCache    OTPCache
	otpRepo     OTPRepo
	rateLimiter RateLimiter
	sms         SMSSender
	events      EventPublisher
	logger      *slog.Logger
	idGen       func() string
}

// NewGenerateOTPUseCase constructs a GenerateOTPUseCase with its dependencies.
func NewGenerateOTPUseCase(
	bookings BookingFinder,
	users UserFinder,
	otpCache OTPCache,
	otpRepo OTPRepo,
	rateLimiter RateLimiter,
	sms SMSSender,
	events EventPublisher,
	logger *slog.Logger,
) *GenerateOTPUseCase {
	return &GenerateOTPUseCase{
		bookings:    bookings,
		users:       users,
		otpCache:    otpCache,
		otpRepo:     otpRepo,
		rateLimiter: rateLimiter,
		sms:         sms,
		events:      events,
		logger:      logger,
		idGen:       generateOTPID,
	}
}

// Execute generates an OTP for a booking, sends it via SMS, and persists it in Redis + DynamoDB.
func (uc *GenerateOTPUseCase) Execute(ctx context.Context, input *GenerateOTPInput) (*GenerateOTPOutput, error) {
	// 1. Look up the booking.
	b, err := uc.bookings.FindByID(ctx, input.BookingID)
	if err != nil {
		return nil, err
	}
	if b == nil {
		return nil, domainerrors.NewNotFoundError("Booking", input.BookingID)
	}

	// 2. Validate booking is in the correct status for the requested OTP type.
	if err := validateBookingStatusForOTP(b, input.OTPType); err != nil {
		return nil, err
	}

	// 3. Rate-limit OTP generation per booking.
	rateLimitKey := fmt.Sprintf("otp:generate:%s", input.BookingID)
	allowed, _, err := uc.rateLimiter.CheckRateLimit(ctx, rateLimitKey, rateLimitMaxGenerate, rateLimitWindowSec)
	if err != nil {
		uc.logger.ErrorContext(ctx, "rate limiter error", "error", err, "booking_id", input.BookingID)
		return nil, domainerrors.NewInternalError("rate limiter unavailable").WithCause(err)
	}
	if !allowed {
		return nil, domainerrors.NewRateLimitedError(rateLimitWindowSec)
	}

	// 4. Generate OTP (domain logic).
	otpID := uc.idGen()
	otpRecord, plainCode, err := otp.Generate(otpID, input.BookingID, input.OTPType, input.Lat, input.Lng)
	if err != nil {
		return nil, domainerrors.NewInternalError("OTP generation failed").WithCause(err)
	}

	// 5. Store hash in Redis with TTL.
	if err := uc.otpCache.StoreOTP(ctx, input.BookingID, string(input.OTPType), otpRecord.CodeHash, otpTTL); err != nil {
		uc.logger.ErrorContext(ctx, "failed to store OTP in cache", "error", err, "booking_id", input.BookingID)
		return nil, domainerrors.NewInternalError("OTP storage failed").WithCause(err)
	}

	// 6. Backup in DynamoDB.
	if err := uc.otpRepo.Save(ctx, otpRecord); err != nil {
		uc.logger.ErrorContext(ctx, "failed to backup OTP in DynamoDB", "error", err, "booking_id", input.BookingID)
		// Non-fatal: Redis is the primary store. Log and continue.
	}

	// 7. Look up customer phone and send SMS.
	u, err := uc.users.FindByID(ctx, b.CustomerID)
	if err != nil {
		uc.logger.ErrorContext(ctx, "failed to find customer for SMS", "error", err, "customer_id", b.CustomerID)
	}
	if u != nil && u.Phone != "" {
		smsMsg := fmt.Sprintf("Your TowCommand verification code: %s. Valid for 5 minutes. Do not share this code.", plainCode)
		if err := uc.sms.SendSMS(ctx, u.Phone, smsMsg); err != nil {
			uc.logger.ErrorContext(ctx, "failed to send OTP SMS", "error", err, "booking_id", input.BookingID)
			// Non-fatal: OTP is stored, customer can request resend.
		}
	}

	// 8. Publish OTPGenerated event.
	_ = uc.events.Publish(ctx, eventSourceOTP, eventOTPGenerated, map[string]any{
		"bookingId": input.BookingID,
		"otpType":   string(input.OTPType),
		"otpId":     otpID,
	}, &Actor{UserID: input.CallerID, UserType: "system"})

	return &GenerateOTPOutput{
		Success:   true,
		BookingID: input.BookingID,
		OTPType:   string(input.OTPType),
	}, nil
}

// validateBookingStatusForOTP checks that the booking is in the correct status
// for the requested OTP type.
// PICKUP OTP requires ARRIVED status.
// DROPOFF OTP requires ARRIVED_DROPOFF status.
func validateBookingStatusForOTP(b *booking.Booking, otpType otp.OTPType) error {
	switch otpType {
	case otp.OTPTypePickup:
		if b.Status != booking.BookingStatusArrived {
			return domainerrors.NewConflictError(
				fmt.Sprintf("booking must be in ARRIVED status for pickup OTP, current: %s", b.Status))
		}
	case otp.OTPTypeDropoff:
		if b.Status != booking.BookingStatusArrivedDropoff {
			return domainerrors.NewConflictError(
				fmt.Sprintf("booking must be in ARRIVED_DROPOFF status for dropoff OTP, current: %s", b.Status))
		}
	default:
		return domainerrors.NewValidationError(fmt.Sprintf("invalid OTP type: %s", otpType))
	}
	return nil
}

// generateOTPID produces a random OTP ID.
func generateOTPID() string {
	b := make([]byte, 8)
	_, _ = rand.Read(b)
	return fmt.Sprintf("OTP-%X", b)
}
