package errors_test

import (
	"errors"
	"fmt"
	"net/http"
	"testing"

	domainerrors "github.com/davidlramirez95/towcommand/internal/domain/errors"
)

// ---------------------------------------------------------------------------
// Constructor + HTTP status mapping tests (table-driven)
// ---------------------------------------------------------------------------

func TestConstructorsAndHTTPStatus(t *testing.T) {
	tests := []struct {
		name           string
		err            *domainerrors.AppError
		wantCode       domainerrors.ErrorCode
		wantHTTP       int
		wantOp         bool
		wantMsgContain string
	}{
		{
			name:           "NewValidationError",
			err:            domainerrors.NewValidationError("field is required"),
			wantCode:       domainerrors.CodeValidationError,
			wantHTTP:       http.StatusBadRequest,
			wantOp:         true,
			wantMsgContain: "field is required",
		},
		{
			name:           "NewNotFoundError",
			err:            domainerrors.NewNotFoundError("Booking", "abc-123"),
			wantCode:       domainerrors.CodeNotFound,
			wantHTTP:       http.StatusNotFound,
			wantOp:         true,
			wantMsgContain: "Booking not found",
		},
		{
			name:           "NewUnauthorizedError",
			err:            domainerrors.NewUnauthorizedError(),
			wantCode:       domainerrors.CodeUnauthorized,
			wantHTTP:       http.StatusUnauthorized,
			wantOp:         true,
			wantMsgContain: "Unauthorized",
		},
		{
			name:           "NewForbiddenError",
			err:            domainerrors.NewForbiddenError("admin access required"),
			wantCode:       domainerrors.CodeForbidden,
			wantHTTP:       http.StatusForbidden,
			wantOp:         true,
			wantMsgContain: "admin access required",
		},
		{
			name:           "NewConflictError",
			err:            domainerrors.NewConflictError("resource already exists"),
			wantCode:       domainerrors.CodeConflict,
			wantHTTP:       http.StatusConflict,
			wantOp:         true,
			wantMsgContain: "resource already exists",
		},
		{
			name:           "NewRateLimitedError",
			err:            domainerrors.NewRateLimitedError(60),
			wantCode:       domainerrors.CodeRateLimited,
			wantHTTP:       http.StatusTooManyRequests,
			wantOp:         true,
			wantMsgContain: "Too many requests",
		},
		{
			name:           "NewInternalError",
			err:            domainerrors.NewInternalError("unexpected failure"),
			wantCode:       domainerrors.CodeInternalError,
			wantHTTP:       http.StatusInternalServerError,
			wantOp:         false,
			wantMsgContain: "unexpected failure",
		},
		{
			name:           "NewExternalServiceError",
			err:            domainerrors.NewExternalServiceError("PayMongo", fmt.Errorf("timeout")),
			wantCode:       domainerrors.CodeExternalService,
			wantHTTP:       http.StatusBadGateway,
			wantOp:         true,
			wantMsgContain: "PayMongo",
		},
		{
			name:           "NewInvalidStatusTransitionError",
			err:            domainerrors.NewInvalidStatusTransitionError("PENDING", "COMPLETED"),
			wantCode:       domainerrors.CodeInvalidStatusTransition,
			wantHTTP:       http.StatusConflict,
			wantOp:         true,
			wantMsgContain: "PENDING",
		},
		{
			name:           "NewBookingNotCancellableError",
			err:            domainerrors.NewBookingNotCancellableError("bk-1", "EN_ROUTE"),
			wantCode:       domainerrors.CodeBookingNotCancellable,
			wantHTTP:       http.StatusConflict,
			wantOp:         true,
			wantMsgContain: "bk-1",
		},
		{
			name:           "NewProviderUnavailableError",
			err:            domainerrors.NewProviderUnavailableError(),
			wantCode:       domainerrors.CodeProviderUnavailable,
			wantHTTP:       http.StatusServiceUnavailable,
			wantOp:         true,
			wantMsgContain: "No providers available",
		},
		{
			name:           "NewPaymentFailedError",
			err:            domainerrors.NewPaymentFailedError("card declined"),
			wantCode:       domainerrors.CodePaymentFailed,
			wantHTTP:       http.StatusBadGateway,
			wantOp:         true,
			wantMsgContain: "card declined",
		},
		{
			name:           "NewOTPExpiredError",
			err:            domainerrors.NewOTPExpiredError(),
			wantCode:       domainerrors.CodeOTPExpired,
			wantHTTP:       http.StatusBadRequest,
			wantOp:         true,
			wantMsgContain: "expired",
		},
		{
			name:           "NewOTPInvalidError",
			err:            domainerrors.NewOTPInvalidError(),
			wantCode:       domainerrors.CodeOTPInvalid,
			wantHTTP:       http.StatusBadRequest,
			wantOp:         true,
			wantMsgContain: "Invalid OTP",
		},
		{
			name:           "NewEvidenceValidationFailedError",
			err:            domainerrors.NewEvidenceValidationFailedError("image too small"),
			wantCode:       domainerrors.CodeEvidenceValidationFailed,
			wantHTTP:       http.StatusBadRequest,
			wantOp:         true,
			wantMsgContain: "image too small",
		},
		{
			name:           "NewSOSActiveError",
			err:            domainerrors.NewSOSActiveError(),
			wantCode:       domainerrors.CodeSOSActive,
			wantHTTP:       http.StatusConflict,
			wantOp:         true,
			wantMsgContain: "SOS",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.err.Code != tt.wantCode {
				t.Errorf("Code = %q, want %q", tt.err.Code, tt.wantCode)
			}
			if got := tt.err.HTTPStatusCode(); got != tt.wantHTTP {
				t.Errorf("HTTPStatusCode() = %d, want %d", got, tt.wantHTTP)
			}
			if tt.err.IsOperational != tt.wantOp {
				t.Errorf("IsOperational = %v, want %v", tt.err.IsOperational, tt.wantOp)
			}
			if msg := tt.err.Error(); !containsSubstring(msg, tt.wantMsgContain) {
				t.Errorf("Error() = %q, want substring %q", msg, tt.wantMsgContain)
			}
		})
	}
}

// ---------------------------------------------------------------------------
// Error() string formatting
// ---------------------------------------------------------------------------

func TestErrorString(t *testing.T) {
	t.Run("without cause", func(t *testing.T) {
		err := domainerrors.NewUnauthorizedError()
		want := "UNAUTHORIZED: Unauthorized"
		if got := err.Error(); got != want {
			t.Errorf("Error() = %q, want %q", got, want)
		}
	})

	t.Run("with cause", func(t *testing.T) {
		cause := fmt.Errorf("connection refused")
		err := domainerrors.NewExternalServiceError("Redis", cause)
		got := err.Error()
		if !containsSubstring(got, "Redis") {
			t.Errorf("Error() = %q, want to contain %q", got, "Redis")
		}
		if !containsSubstring(got, "connection refused") {
			t.Errorf("Error() = %q, want to contain wrapped cause", got)
		}
	})
}

// ---------------------------------------------------------------------------
// errors.Is support
// ---------------------------------------------------------------------------

func TestErrorsIs(t *testing.T) {
	t.Run("matches any AppError with zero-value code", func(t *testing.T) {
		err := domainerrors.NewNotFoundError("User", "u-1")
		if !errors.Is(err, &domainerrors.AppError{}) {
			t.Error("errors.Is(err, &AppError{}) should be true")
		}
	})

	t.Run("matches same code", func(t *testing.T) {
		err := domainerrors.NewNotFoundError("User", "u-1")
		target := &domainerrors.AppError{Code: domainerrors.CodeNotFound}
		if !errors.Is(err, target) {
			t.Error("errors.Is should match same code")
		}
	})

	t.Run("does not match different code", func(t *testing.T) {
		err := domainerrors.NewNotFoundError("User", "u-1")
		target := &domainerrors.AppError{Code: domainerrors.CodeForbidden}
		if errors.Is(err, target) {
			t.Error("errors.Is should not match different code")
		}
	})

	t.Run("does not match non-AppError", func(t *testing.T) {
		err := domainerrors.NewInternalError("oops")
		if errors.Is(err, fmt.Errorf("oops")) {
			t.Error("errors.Is should not match plain error")
		}
	})

	t.Run("matches through wrapped chain", func(t *testing.T) {
		inner := domainerrors.NewPaymentFailedError("declined")
		outer := fmt.Errorf("handler: %w", inner)
		if !errors.Is(outer, &domainerrors.AppError{Code: domainerrors.CodePaymentFailed}) {
			t.Error("errors.Is should match through wrapping")
		}
	})
}

// ---------------------------------------------------------------------------
// errors.As support
// ---------------------------------------------------------------------------

func TestErrorsAs(t *testing.T) {
	t.Run("extracts AppError from plain error variable", func(t *testing.T) {
		var err error = domainerrors.NewForbiddenError("no access")
		var appErr *domainerrors.AppError
		if !errors.As(err, &appErr) {
			t.Fatal("errors.As should find *AppError")
		}
		if appErr.Code != domainerrors.CodeForbidden {
			t.Errorf("Code = %q, want %q", appErr.Code, domainerrors.CodeForbidden)
		}
	})

	t.Run("extracts AppError through wrapped chain", func(t *testing.T) {
		inner := domainerrors.NewOTPExpiredError()
		outer := fmt.Errorf("auth: %w", inner)
		var appErr *domainerrors.AppError
		if !errors.As(outer, &appErr) {
			t.Fatal("errors.As should unwrap to find *AppError")
		}
		if appErr.Code != domainerrors.CodeOTPExpired {
			t.Errorf("Code = %q, want %q", appErr.Code, domainerrors.CodeOTPExpired)
		}
	})

	t.Run("returns false for non-AppError", func(t *testing.T) {
		err := fmt.Errorf("plain error")
		var appErr *domainerrors.AppError
		if errors.As(err, &appErr) {
			t.Error("errors.As should return false for non-AppError")
		}
	})
}

// ---------------------------------------------------------------------------
// Unwrap
// ---------------------------------------------------------------------------

func TestUnwrap(t *testing.T) {
	t.Run("nil when no cause", func(t *testing.T) {
		err := domainerrors.NewInternalError("boom")
		if err.Unwrap() != nil {
			t.Error("Unwrap() should be nil when no cause")
		}
	})

	t.Run("returns cause from NewExternalServiceError", func(t *testing.T) {
		cause := fmt.Errorf("dial tcp: timeout")
		err := domainerrors.NewExternalServiceError("DynamoDB", cause)
		if err.Unwrap() != cause {
			t.Errorf("Unwrap() = %v, want %v", err.Unwrap(), cause)
		}
	})

	t.Run("returns cause from WithCause", func(t *testing.T) {
		cause := fmt.Errorf("root cause")
		err := domainerrors.NewInternalError("wrapper").WithCause(cause)
		if err.Unwrap() != cause {
			t.Errorf("Unwrap() = %v, want %v", err.Unwrap(), cause)
		}
	})
}

// ---------------------------------------------------------------------------
// WithCause
// ---------------------------------------------------------------------------

func TestWithCause(t *testing.T) {
	original := domainerrors.NewValidationError("bad input")
	cause := fmt.Errorf("parse error")
	wrapped := original.WithCause(cause)

	if wrapped == original {
		t.Error("WithCause should return a new AppError, not mutate the original")
	}
	if original.Unwrap() != nil {
		t.Error("original should remain unchanged")
	}
	if wrapped.Unwrap() != cause {
		t.Error("wrapped error should have the cause")
	}
	if wrapped.Code != original.Code {
		t.Error("wrapped error should preserve the code")
	}
}

// ---------------------------------------------------------------------------
// WithDetails
// ---------------------------------------------------------------------------

func TestWithDetails(t *testing.T) {
	original := domainerrors.NewNotFoundError("Booking", "bk-1")
	extra := map[string]any{"hint": "check region"}
	withExtra := original.WithDetails(extra)

	if withExtra == original {
		t.Error("WithDetails should return a new AppError")
	}
	if _, ok := withExtra.Details["hint"]; !ok {
		t.Error("WithDetails should merge new keys")
	}
	if _, ok := withExtra.Details["resource"]; !ok {
		t.Error("WithDetails should preserve existing keys")
	}
	// Original should not be mutated.
	if _, ok := original.Details["hint"]; ok {
		t.Error("original should not be mutated")
	}
}

// ---------------------------------------------------------------------------
// HTTPStatusCode unknown code fallback
// ---------------------------------------------------------------------------

func TestHTTPStatusCodeUnknownCode(t *testing.T) {
	err := &domainerrors.AppError{Code: "TOTALLY_UNKNOWN", Message: "wat"}
	if got := err.HTTPStatusCode(); got != http.StatusInternalServerError {
		t.Errorf("HTTPStatusCode() = %d, want %d for unknown code", got, http.StatusInternalServerError)
	}
}

// ---------------------------------------------------------------------------
// Details on specific constructors
// ---------------------------------------------------------------------------

func TestConstructorDetails(t *testing.T) {
	t.Run("NewNotFoundError details", func(t *testing.T) {
		err := domainerrors.NewNotFoundError("Vehicle", "v-99")
		assertDetail(t, err, "resource", "Vehicle")
		assertDetail(t, err, "id", "v-99")
	})

	t.Run("NewRateLimitedError retryAfter", func(t *testing.T) {
		err := domainerrors.NewRateLimitedError(30)
		assertDetail(t, err, "retryAfter", 30)
	})

	t.Run("NewExternalServiceError service", func(t *testing.T) {
		err := domainerrors.NewExternalServiceError("Stripe", nil)
		assertDetail(t, err, "service", "Stripe")
	})

	t.Run("NewInvalidStatusTransitionError from/to", func(t *testing.T) {
		err := domainerrors.NewInvalidStatusTransitionError("MATCHED", "PENDING")
		assertDetail(t, err, "from", "MATCHED")
		assertDetail(t, err, "to", "PENDING")
	})

	t.Run("NewBookingNotCancellableError bookingID/status", func(t *testing.T) {
		err := domainerrors.NewBookingNotCancellableError("bk-5", "COMPLETED")
		assertDetail(t, err, "bookingID", "bk-5")
		assertDetail(t, err, "status", "COMPLETED")
	})

	t.Run("NewPaymentFailedError reason", func(t *testing.T) {
		err := domainerrors.NewPaymentFailedError("insufficient funds")
		assertDetail(t, err, "reason", "insufficient funds")
	})

	t.Run("NewEvidenceValidationFailedError reason", func(t *testing.T) {
		err := domainerrors.NewEvidenceValidationFailedError("blurry photo")
		assertDetail(t, err, "reason", "blurry photo")
	})
}

// ---------------------------------------------------------------------------
// Helpers
// ---------------------------------------------------------------------------

func containsSubstring(s, substr string) bool {
	return len(s) >= len(substr) && searchSubstring(s, substr)
}

func searchSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

func assertDetail(t *testing.T, err *domainerrors.AppError, key string, want any) {
	t.Helper()
	got, ok := err.Details[key]
	if !ok {
		t.Errorf("Details[%q] missing", key)
		return
	}
	if fmt.Sprintf("%v", got) != fmt.Sprintf("%v", want) {
		t.Errorf("Details[%q] = %v, want %v", key, got, want)
	}
}
