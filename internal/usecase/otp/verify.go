package otpuc

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/davidlramirez95/towcommand/internal/domain/booking"
	domainerrors "github.com/davidlramirez95/towcommand/internal/domain/errors"
	"github.com/davidlramirez95/towcommand/internal/domain/otp"
)

const eventOTPVerified = "OTPVerified"

// VerifyOTPInput holds the data needed to verify an OTP.
type VerifyOTPInput struct {
	BookingID string
	OTPType   otp.OTPType
	Code      string
	Lat       float64
	Lng       float64
	CallerID  string
}

// VerifyOTPOutput is the response for OTP verification.
type VerifyOTPOutput struct {
	Success   bool   `json:"success"`
	BookingID string `json:"bookingId"`
	OTPType   string `json:"otpType"`
	NewStatus string `json:"newStatus"`
}

// VerifyOTPUseCase orchestrates OTP verification with Redis + DynamoDB fallback.
type VerifyOTPUseCase struct {
	bookings BookingFinder
	statuses BookingStatusUpdater
	otpCache OTPCache
	otpRepo  OTPRepo
	events   EventPublisher
	logger   *slog.Logger
}

// NewVerifyOTPUseCase constructs a VerifyOTPUseCase with its dependencies.
func NewVerifyOTPUseCase(
	bookings BookingFinder,
	statuses BookingStatusUpdater,
	otpCache OTPCache,
	otpRepo OTPRepo,
	events EventPublisher,
	logger *slog.Logger,
) *VerifyOTPUseCase {
	return &VerifyOTPUseCase{
		bookings: bookings,
		statuses: statuses,
		otpCache: otpCache,
		otpRepo:  otpRepo,
		events:   events,
		logger:   logger,
	}
}

// Execute verifies an OTP code against the stored hash, transitions the booking
// status, and cleans up the OTP from both stores.
func (uc *VerifyOTPUseCase) Execute(ctx context.Context, input *VerifyOTPInput) (*VerifyOTPOutput, error) {
	// 1. Look up the booking.
	b, err := uc.bookings.FindByID(ctx, input.BookingID)
	if err != nil {
		return nil, err
	}
	if b == nil {
		return nil, domainerrors.NewNotFoundError("Booking", input.BookingID)
	}

	// 2. Validate booking status for verification.
	// Verification expects the booking in CONDITION_REPORT (pickup) or ARRIVED_DROPOFF (dropoff).
	if err := validateBookingStatusForVerify(b, input.OTPType); err != nil {
		return nil, err
	}

	// 3. Get stored hash from Redis (primary), fallback to DynamoDB.
	otpRecord, err := uc.resolveOTPRecord(ctx, input.BookingID, input.OTPType)
	if err != nil {
		return nil, err
	}

	// 4. Validate the OTP code and location.
	if err := otpRecord.Validate(input.Code, input.Lat, input.Lng); err != nil {
		return nil, mapOTPDomainError(err)
	}

	// 5. Determine the target booking status after verification.
	newStatus, err := targetStatusAfterVerification(input.OTPType)
	if err != nil {
		return nil, err
	}

	// 6. Transition booking status.
	if err := uc.statuses.UpdateStatus(ctx, input.BookingID, newStatus, map[string]any{
		"verifiedBy": input.CallerID,
		"otpType":    string(input.OTPType),
	}); err != nil {
		return nil, err
	}

	// 7. Clean up: delete from Redis, mark verified in DynamoDB.
	if err := uc.otpCache.DeleteOTP(ctx, input.BookingID, string(input.OTPType)); err != nil {
		uc.logger.ErrorContext(ctx, "failed to delete OTP from cache", "error", err, "booking_id", input.BookingID)
	}
	if err := uc.otpRepo.MarkVerified(ctx, input.BookingID, string(input.OTPType)); err != nil {
		uc.logger.ErrorContext(ctx, "failed to mark OTP verified in DynamoDB", "error", err, "booking_id", input.BookingID)
	}

	// 8. Publish OTPVerified event.
	_ = uc.events.Publish(ctx, eventSourceOTP, eventOTPVerified, map[string]any{
		"bookingId": input.BookingID,
		"otpType":   string(input.OTPType),
		"newStatus": string(newStatus),
	}, &Actor{UserID: input.CallerID, UserType: "system"})

	return &VerifyOTPOutput{
		Success:   true,
		BookingID: input.BookingID,
		OTPType:   string(input.OTPType),
		NewStatus: string(newStatus),
	}, nil
}

// resolveOTPRecord attempts to get the OTP from Redis first, falling back to DynamoDB.
func (uc *VerifyOTPUseCase) resolveOTPRecord(ctx context.Context, bookingID string, otpType otp.OTPType) (*otp.OTP, error) {
	// Try Redis first.
	hash, err := uc.otpCache.GetOTP(ctx, bookingID, string(otpType))
	if err != nil {
		uc.logger.ErrorContext(ctx, "Redis OTP cache error, falling back to DynamoDB",
			"error", err, "booking_id", bookingID)
	}

	if hash != "" {
		// Redis has the hash; reconstruct a minimal OTP record for validation.
		// We still need the DynamoDB record for full metadata (lat, lng, attempts, expiresAt).
		record, err := uc.otpRepo.FindByBookingAndType(ctx, bookingID, string(otpType))
		if err != nil {
			return nil, domainerrors.NewInternalError("failed to retrieve OTP backup").WithCause(err)
		}
		if record != nil {
			return record, nil
		}
		// DynamoDB backup missing but Redis has hash — this is an edge case.
		// We cannot validate without full record metadata. Return not found.
		uc.logger.WarnContext(ctx, "OTP hash found in Redis but not in DynamoDB backup",
			"booking_id", bookingID, "otp_type", string(otpType))
		return nil, domainerrors.NewNotFoundError("OTP", fmt.Sprintf("%s:%s", bookingID, string(otpType)))
	}

	// Fallback to DynamoDB.
	record, err := uc.otpRepo.FindByBookingAndType(ctx, bookingID, string(otpType))
	if err != nil {
		return nil, domainerrors.NewInternalError("failed to retrieve OTP").WithCause(err)
	}
	if record == nil {
		return nil, domainerrors.NewNotFoundError("OTP", fmt.Sprintf("%s:%s", bookingID, string(otpType)))
	}

	uc.logger.InfoContext(ctx, "OTP resolved from DynamoDB fallback",
		"booking_id", bookingID, "otp_type", string(otpType))
	return record, nil
}

// validateBookingStatusForVerify checks that the booking is in the correct status
// for OTP verification.
// PICKUP verification requires CONDITION_REPORT status (CONDITION_REPORT -> OTP_VERIFIED).
// DROPOFF verification requires ARRIVED_DROPOFF status (ARRIVED_DROPOFF -> OTP_DROPOFF).
func validateBookingStatusForVerify(b *booking.Booking, otpType otp.OTPType) error {
	switch otpType {
	case otp.OTPTypePickup:
		if b.Status != booking.BookingStatusConditionReport {
			return domainerrors.NewConflictError(
				fmt.Sprintf("booking must be in CONDITION_REPORT status for pickup OTP verification, current: %s", b.Status))
		}
	case otp.OTPTypeDropoff:
		if b.Status != booking.BookingStatusArrivedDropoff {
			return domainerrors.NewConflictError(
				fmt.Sprintf("booking must be in ARRIVED_DROPOFF status for dropoff OTP verification, current: %s", b.Status))
		}
	default:
		return domainerrors.NewValidationError(fmt.Sprintf("invalid OTP type: %s", otpType))
	}
	return nil
}

// targetStatusAfterVerification determines the booking status to transition to
// after successful OTP verification.
func targetStatusAfterVerification(otpType otp.OTPType) (booking.BookingStatus, error) {
	switch otpType {
	case otp.OTPTypePickup:
		// ARRIVED → CONDITION_REPORT (via status update), then CONDITION_REPORT → OTP_VERIFIED
		return booking.BookingStatusOTPVerified, nil
	case otp.OTPTypeDropoff:
		// ARRIVED_DROPOFF → OTP_DROPOFF
		return booking.BookingStatusOTPDropoff, nil
	default:
		return "", domainerrors.NewValidationError(fmt.Sprintf("invalid OTP type: %s", otpType))
	}
}

// mapOTPDomainError maps domain OTP errors to AppError types.
func mapOTPDomainError(err error) error {
	switch {
	case errors.Is(err, otp.ErrExpired):
		return domainerrors.NewOTPExpiredError()
	case errors.Is(err, otp.ErrInvalidCode):
		return domainerrors.NewOTPInvalidError()
	case errors.Is(err, otp.ErrMaxAttempts):
		return domainerrors.NewOTPInvalidError().WithDetails(map[string]any{"reason": "maximum attempts exceeded"})
	case errors.Is(err, otp.ErrAlreadyVerified):
		return domainerrors.NewConflictError("OTP already verified")
	case errors.Is(err, otp.ErrNotInProximity):
		return domainerrors.NewValidationError("not within required proximity for verification")
	default:
		return domainerrors.NewInternalError("OTP validation failed").WithCause(err)
	}
}
