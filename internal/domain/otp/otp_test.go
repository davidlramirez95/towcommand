package otp

import (
	"testing"
	"time"
)

func TestGenerate(t *testing.T) {
	o, code, err := Generate("otp-1", "booking-1", OTPTypePickup, 14.5995, 120.9842)
	if err != nil {
		t.Fatalf("Generate() error = %v", err)
	}

	if len(code) != otpLength {
		t.Errorf("code length = %d, want %d", len(code), otpLength)
	}

	if o.OTPID != "otp-1" {
		t.Errorf("OTPID = %q, want %q", o.OTPID, "otp-1")
	}
	if o.BookingID != "booking-1" {
		t.Errorf("BookingID = %q, want %q", o.BookingID, "booking-1")
	}
	if o.Type != OTPTypePickup {
		t.Errorf("Type = %q, want %q", o.Type, OTPTypePickup)
	}
	if o.Verified {
		t.Error("new OTP should not be verified")
	}
	if o.Attempts != 0 {
		t.Errorf("Attempts = %d, want 0", o.Attempts)
	}
	if o.CodeHash == "" {
		t.Error("CodeHash should not be empty")
	}
	if o.CodeHash == code {
		t.Error("CodeHash should not equal plaintext code")
	}
}

func TestGenerate_UniqueCodesAndHashes(t *testing.T) {
	codes := make(map[string]bool)
	hashes := make(map[string]bool)
	for i := 0; i < 100; i++ {
		o, code, err := Generate("otp", "booking", OTPTypePickup, 14.5, 120.9)
		if err != nil {
			t.Fatalf("Generate() error on iteration %d: %v", i, err)
		}
		codes[code] = true
		hashes[o.CodeHash] = true
	}
	// With 6-digit codes (1M possibilities) and 100 samples, collisions are extremely unlikely.
	if len(codes) < 95 {
		t.Errorf("expected mostly unique codes, got %d unique out of 100", len(codes))
	}
	if len(hashes) < 95 {
		t.Errorf("expected mostly unique hashes, got %d unique out of 100", len(hashes))
	}
}

func TestValidate_Success(t *testing.T) {
	o, code, err := Generate("otp-1", "booking-1", OTPTypePickup, 14.5995, 120.9842)
	if err != nil {
		t.Fatalf("Generate() error = %v", err)
	}

	// Validate with correct code and nearby location (~0m away)
	err = o.Validate(code, 14.5995, 120.9842)
	if err != nil {
		t.Errorf("Validate() error = %v, want nil", err)
	}
	if !o.Verified {
		t.Error("OTP should be verified after successful validation")
	}
	if o.Attempts != 1 {
		t.Errorf("Attempts = %d, want 1", o.Attempts)
	}
}

func TestValidate_WrongCode(t *testing.T) {
	o, _, err := Generate("otp-1", "booking-1", OTPTypePickup, 14.5995, 120.9842)
	if err != nil {
		t.Fatalf("Generate() error = %v", err)
	}

	err = o.Validate("000000", 14.5995, 120.9842)
	if err != ErrInvalidCode {
		t.Errorf("Validate() error = %v, want %v", err, ErrInvalidCode)
	}
	if o.Verified {
		t.Error("OTP should not be verified after wrong code")
	}
	if o.Attempts != 1 {
		t.Errorf("Attempts = %d, want 1", o.Attempts)
	}
}

func TestValidate_MaxAttempts(t *testing.T) {
	o, _, err := Generate("otp-1", "booking-1", OTPTypePickup, 14.5995, 120.9842)
	if err != nil {
		t.Fatalf("Generate() error = %v", err)
	}

	// Exhaust all attempts
	for i := 0; i < maxAttempts; i++ {
		_ = o.Validate("000000", 14.5995, 120.9842)
	}

	if o.Attempts != maxAttempts {
		t.Errorf("Attempts = %d, want %d", o.Attempts, maxAttempts)
	}

	// Next attempt should fail with max attempts error
	err = o.Validate("123456", 14.5995, 120.9842)
	if err != ErrMaxAttempts {
		t.Errorf("Validate() error = %v, want %v", err, ErrMaxAttempts)
	}
}

func TestValidate_Expired(t *testing.T) {
	o, code, err := Generate("otp-1", "booking-1", OTPTypePickup, 14.5995, 120.9842)
	if err != nil {
		t.Fatalf("Generate() error = %v", err)
	}

	// Force expiry
	o.ExpiresAt = time.Now().Add(-1 * time.Second)

	err = o.Validate(code, 14.5995, 120.9842)
	if err != ErrExpired {
		t.Errorf("Validate() error = %v, want %v", err, ErrExpired)
	}
}

func TestValidate_AlreadyVerified(t *testing.T) {
	o, code, err := Generate("otp-1", "booking-1", OTPTypePickup, 14.5995, 120.9842)
	if err != nil {
		t.Fatalf("Generate() error = %v", err)
	}

	_ = o.Validate(code, 14.5995, 120.9842)

	err = o.Validate(code, 14.5995, 120.9842)
	if err != ErrAlreadyVerified {
		t.Errorf("Validate() error = %v, want %v", err, ErrAlreadyVerified)
	}
}

func TestValidate_NotInProximity(t *testing.T) {
	o, code, err := Generate("otp-1", "booking-1", OTPTypePickup, 14.5995, 120.9842)
	if err != nil {
		t.Fatalf("Generate() error = %v", err)
	}

	// ~10km away
	err = o.Validate(code, 14.6900, 120.9842)
	if err != ErrNotInProximity {
		t.Errorf("Validate() error = %v, want %v", err, ErrNotInProximity)
	}
}

func TestValidate_WithinProximity(t *testing.T) {
	o, code, err := Generate("otp-1", "booking-1", OTPTypePickup, 14.5995, 120.9842)
	if err != nil {
		t.Fatalf("Generate() error = %v", err)
	}

	// ~100m offset (approx 0.001 degrees in Manila)
	err = o.Validate(code, 14.6003, 120.9842)
	if err != nil {
		t.Errorf("Validate() with nearby location error = %v, want nil", err)
	}
}

func TestIsExpired(t *testing.T) {
	o, _, err := Generate("otp-1", "booking-1", OTPTypePickup, 14.5995, 120.9842)
	if err != nil {
		t.Fatalf("Generate() error = %v", err)
	}

	if o.IsExpired() {
		t.Error("newly generated OTP should not be expired")
	}

	o.ExpiresAt = time.Now().Add(-1 * time.Second)
	if !o.IsExpired() {
		t.Error("OTP with past expiry should be expired")
	}
}

func TestHaversineMeters(t *testing.T) {
	tests := []struct {
		name    string
		lat1    float64
		lng1    float64
		lat2    float64
		lng2    float64
		wantMin float64
		wantMax float64
	}{
		{"same point", 14.5995, 120.9842, 14.5995, 120.9842, 0, 0.1},
		{"Manila to Makati ~5km", 14.5995, 120.9842, 14.5547, 121.0244, 4000, 7000},
		{"within 500m", 14.5995, 120.9842, 14.6020, 120.9842, 200, 400},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := haversineMeters(tt.lat1, tt.lng1, tt.lat2, tt.lng2)
			if d < tt.wantMin || d > tt.wantMax {
				t.Errorf("haversineMeters() = %f, want between %f and %f", d, tt.wantMin, tt.wantMax)
			}
		})
	}
}

func TestHashCode_Deterministic(t *testing.T) {
	h1 := hashCode("123456")
	h2 := hashCode("123456")
	if h1 != h2 {
		t.Error("hashCode should be deterministic")
	}

	h3 := hashCode("654321")
	if h1 == h3 {
		t.Error("different codes should produce different hashes")
	}
}
