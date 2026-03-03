// Package errors defines domain error types for the TowCommand platform.
//
// It provides an AppError type that implements the standard error interface
// with support for error codes, HTTP status mapping, and Go error wrapping
// via errors.Is and errors.As.
package errors

import (
	"fmt"
	"net/http"
)

// ErrorCode categorises domain errors.
type ErrorCode string

const (
	// CodeValidationError is the error code for request validation failures.
	CodeValidationError ErrorCode = "VALIDATION_ERROR"
	// CodeNotFound is the error code for missing resources.
	CodeNotFound ErrorCode = "NOT_FOUND"
	// CodeUnauthorized is the error code for unauthenticated requests.
	CodeUnauthorized ErrorCode = "UNAUTHORIZED"
	// CodeForbidden is the error code for insufficient permissions.
	CodeForbidden ErrorCode = "FORBIDDEN"
	// CodeConflict is the error code for state conflicts.
	CodeConflict ErrorCode = "CONFLICT"
	// CodeRateLimited is the error code for rate-limit violations.
	CodeRateLimited ErrorCode = "RATE_LIMITED"
	// CodeInternalError is the error code for unexpected internal failures.
	CodeInternalError ErrorCode = "INTERNAL_ERROR"
	// CodeExternalService is the error code for upstream service failures.
	CodeExternalService ErrorCode = "EXTERNAL_SERVICE_ERROR"

	// CodeInvalidStatusTransition is the error code for illegal booking state transitions.
	CodeInvalidStatusTransition ErrorCode = "INVALID_STATUS_TRANSITION"
	// CodeBookingNotCancellable is the error code for bookings that cannot be cancelled.
	CodeBookingNotCancellable ErrorCode = "BOOKING_NOT_CANCELLABLE"

	// CodeProviderUnavailable is the error code for no available tow providers.
	CodeProviderUnavailable ErrorCode = "PROVIDER_UNAVAILABLE"

	// CodePaymentFailed is the error code for payment processing failures.
	CodePaymentFailed ErrorCode = "PAYMENT_FAILED"

	// CodeOTPExpired is the error code for expired one-time passwords.
	CodeOTPExpired ErrorCode = "OTP_EXPIRED"
	// CodeOTPInvalid is the error code for incorrect one-time passwords.
	CodeOTPInvalid ErrorCode = "OTP_INVALID"

	// CodeEvidenceValidationFailed is the error code for evidence that fails validation.
	CodeEvidenceValidationFailed ErrorCode = "EVIDENCE_VALIDATION_FAILED"

	// CodeSOSActive is the error code for actions blocked by an active SOS session.
	CodeSOSActive ErrorCode = "SOS_ACTIVE"
)

// httpStatusMap maps each ErrorCode to an HTTP status code.
var httpStatusMap = map[ErrorCode]int{
	CodeValidationError:          http.StatusBadRequest,          // 400
	CodeNotFound:                 http.StatusNotFound,            // 404
	CodeUnauthorized:             http.StatusUnauthorized,        // 401
	CodeForbidden:                http.StatusForbidden,           // 403
	CodeConflict:                 http.StatusConflict,            // 409
	CodeRateLimited:              http.StatusTooManyRequests,     // 429
	CodeInternalError:            http.StatusInternalServerError, // 500
	CodeExternalService:          http.StatusBadGateway,          // 502
	CodeInvalidStatusTransition:  http.StatusConflict,            // 409
	CodeBookingNotCancellable:    http.StatusConflict,            // 409
	CodeProviderUnavailable:      http.StatusServiceUnavailable,  // 503
	CodePaymentFailed:            http.StatusBadGateway,          // 502
	CodeOTPExpired:               http.StatusBadRequest,          // 400
	CodeOTPInvalid:               http.StatusBadRequest,          // 400
	CodeEvidenceValidationFailed: http.StatusBadRequest,          // 400
	CodeSOSActive:                http.StatusConflict,            // 409
}

// AppError is the domain error type used throughout the application.
type AppError struct {
	// Code identifies the category of the error.
	Code ErrorCode
	// Message is a human-readable description.
	Message string
	// IsOperational indicates whether the error is expected (true) or a
	// programming/infrastructure bug (false).
	IsOperational bool
	// Details carries optional structured metadata about the error.
	Details map[string]any
	// cause is the wrapped underlying error, if any.
	cause error
}

// Error implements the error interface.
func (e *AppError) Error() string {
	if e.cause != nil {
		return fmt.Sprintf("%s: %s: %v", e.Code, e.Message, e.cause)
	}
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

// HTTPStatusCode returns the HTTP status code that corresponds to this error's Code.
// Unknown codes default to 500.
func (e *AppError) HTTPStatusCode() int {
	if status, ok := httpStatusMap[e.Code]; ok {
		return status
	}
	return http.StatusInternalServerError
}

// Unwrap returns the underlying cause, enabling errors.Unwrap.
func (e *AppError) Unwrap() error {
	return e.cause
}

// Is reports whether target matches this error. It enables errors.Is support.
//
// When target is an *AppError with a zero-value Code, any AppError matches.
// When target carries a specific Code, only errors with the same Code match.
func (e *AppError) Is(target error) bool {
	t, ok := target.(*AppError)
	if !ok {
		return false
	}
	if t.Code == "" {
		return true
	}
	return e.Code == t.Code
}

// WithCause returns a copy of the error with the given cause wrapped inside.
func (e *AppError) WithCause(err error) *AppError {
	clone := *e
	clone.cause = err
	return &clone
}

// WithDetails returns a copy of the error with the given details merged in.
func (e *AppError) WithDetails(details map[string]any) *AppError {
	clone := *e
	merged := make(map[string]any, len(e.Details)+len(details))
	for k, v := range e.Details {
		merged[k] = v
	}
	for k, v := range details {
		merged[k] = v
	}
	clone.Details = merged
	return &clone
}

// ---------------------------------------------------------------------------
// Constructor functions
// ---------------------------------------------------------------------------

// NewValidationError creates a 400 validation error.
func NewValidationError(msg string) *AppError {
	return &AppError{
		Code:          CodeValidationError,
		Message:       msg,
		IsOperational: true,
	}
}

// NewNotFoundError creates a 404 error for a missing resource.
func NewNotFoundError(resource, id string) *AppError {
	return &AppError{
		Code:          CodeNotFound,
		Message:       fmt.Sprintf("%s not found", resource),
		IsOperational: true,
		Details:       map[string]any{"resource": resource, "id": id},
	}
}

// NewUnauthorizedError creates a 401 error.
func NewUnauthorizedError() *AppError {
	return &AppError{
		Code:          CodeUnauthorized,
		Message:       "Unauthorized",
		IsOperational: true,
	}
}

// NewForbiddenError creates a 403 error.
func NewForbiddenError(msg string) *AppError {
	return &AppError{
		Code:          CodeForbidden,
		Message:       msg,
		IsOperational: true,
	}
}

// NewConflictError creates a 409 conflict error.
func NewConflictError(msg string) *AppError {
	return &AppError{
		Code:          CodeConflict,
		Message:       msg,
		IsOperational: true,
	}
}

// NewRateLimitedError creates a 429 rate-limit error.
func NewRateLimitedError(retryAfterSecs int) *AppError {
	return &AppError{
		Code:          CodeRateLimited,
		Message:       "Too many requests",
		IsOperational: true,
		Details:       map[string]any{"retryAfter": retryAfterSecs},
	}
}

// NewInternalError creates a 500 internal error. These are marked as
// non-operational because they represent unexpected failures.
func NewInternalError(msg string) *AppError {
	return &AppError{
		Code:          CodeInternalError,
		Message:       msg,
		IsOperational: false,
	}
}

// NewExternalServiceError creates a 502 error for failures in upstream
// services. The original error is wrapped as the cause.
func NewExternalServiceError(service string, err error) *AppError {
	return &AppError{
		Code:          CodeExternalService,
		Message:       fmt.Sprintf("external service failure: %s", service),
		IsOperational: true,
		Details:       map[string]any{"service": service},
		cause:         err,
	}
}

// NewInvalidStatusTransitionError creates a 409 error for illegal booking
// state transitions.
func NewInvalidStatusTransitionError(from, to string) *AppError {
	return &AppError{
		Code:          CodeInvalidStatusTransition,
		Message:       fmt.Sprintf("cannot transition from %s to %s", from, to),
		IsOperational: true,
		Details:       map[string]any{"from": from, "to": to},
	}
}

// NewBookingNotCancellableError creates a 409 error when a booking cannot
// be cancelled in its current state.
func NewBookingNotCancellableError(bookingID, status string) *AppError {
	return &AppError{
		Code:          CodeBookingNotCancellable,
		Message:       fmt.Sprintf("booking %s cannot be cancelled in status %s", bookingID, status),
		IsOperational: true,
		Details:       map[string]any{"bookingID": bookingID, "status": status},
	}
}

// NewProviderUnavailableError creates a 503 error when no providers are
// available to fulfil a request.
func NewProviderUnavailableError() *AppError {
	return &AppError{
		Code:          CodeProviderUnavailable,
		Message:       "No providers available",
		IsOperational: true,
	}
}

// NewPaymentFailedError creates a 502 error for payment processing failures.
func NewPaymentFailedError(reason string) *AppError {
	return &AppError{
		Code:          CodePaymentFailed,
		Message:       fmt.Sprintf("payment failed: %s", reason),
		IsOperational: true,
		Details:       map[string]any{"reason": reason},
	}
}

// NewOTPExpiredError creates a 400 error for an expired OTP.
func NewOTPExpiredError() *AppError {
	return &AppError{
		Code:          CodeOTPExpired,
		Message:       "OTP has expired",
		IsOperational: true,
	}
}

// NewOTPInvalidError creates a 400 error for an incorrect OTP.
func NewOTPInvalidError() *AppError {
	return &AppError{
		Code:          CodeOTPInvalid,
		Message:       "Invalid OTP",
		IsOperational: true,
	}
}

// NewEvidenceValidationFailedError creates a 400 error when submitted
// evidence does not pass validation.
func NewEvidenceValidationFailedError(reason string) *AppError {
	return &AppError{
		Code:          CodeEvidenceValidationFailed,
		Message:       fmt.Sprintf("evidence validation failed: %s", reason),
		IsOperational: true,
		Details:       map[string]any{"reason": reason},
	}
}

// NewSOSActiveError creates a 409 error when an action is blocked by an
// active SOS session.
func NewSOSActiveError() *AppError {
	return &AppError{
		Code:          CodeSOSActive,
		Message:       "An SOS session is already active",
		IsOperational: true,
	}
}
