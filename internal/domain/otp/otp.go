package otp

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"math"
	"math/big"
	"time"
)

const (
	otpLength       = 6
	maxAttempts     = 3
	expiryDuration  = 5 * time.Minute
	proximityMeters = 500.0
	earthRadiusKm   = 6371.0
)

// OTPType distinguishes pickup from dropoff verification.
type OTPType string

const (
	OTPTypePickup  OTPType = "PICKUP"
	OTPTypeDropoff OTPType = "DROPOFF"
)

// OTP represents a one-time password for verifying physical presence.
type OTP struct {
	OTPID     string    `json:"otpId" validate:"required"`
	BookingID string    `json:"bookingId" validate:"required"`
	Type      OTPType   `json:"type" validate:"required"`
	CodeHash  string    `json:"codeHash" validate:"required"`
	Lat       float64   `json:"lat" validate:"required,latitude"`
	Lng       float64   `json:"lng" validate:"required,longitude"`
	Attempts  int       `json:"attempts"`
	ExpiresAt time.Time `json:"expiresAt"`
	Verified  bool      `json:"verified"`
	CreatedAt time.Time `json:"createdAt"`
}

var (
	ErrExpired         = errors.New("otp: code has expired")
	ErrMaxAttempts     = errors.New("otp: maximum attempts exceeded")
	ErrAlreadyVerified = errors.New("otp: already verified")
	ErrInvalidCode     = errors.New("otp: invalid code")
	ErrNotInProximity  = errors.New("otp: not within required proximity")
)

// Generate creates a new OTP, returning the OTP record and the plaintext code.
// The plaintext code should be sent to the customer; only the hash is stored.
func Generate(otpID, bookingID string, otpType OTPType, lat, lng float64) (*OTP, string, error) {
	code, err := generateCode()
	if err != nil {
		return nil, "", fmt.Errorf("otp: generating code: %w", err)
	}

	now := time.Now()
	o := &OTP{
		OTPID:     otpID,
		BookingID: bookingID,
		Type:      otpType,
		CodeHash:  hashCode(code),
		Lat:       lat,
		Lng:       lng,
		Attempts:  0,
		ExpiresAt: now.Add(expiryDuration),
		Verified:  false,
		CreatedAt: now,
	}
	return o, code, nil
}

// Validate checks the provided code and location against the OTP.
// It increments the attempt counter on each call regardless of outcome.
func (o *OTP) Validate(code string, lat, lng float64) error {
	if o.Verified {
		return ErrAlreadyVerified
	}
	if o.Attempts >= maxAttempts {
		return ErrMaxAttempts
	}
	if time.Now().After(o.ExpiresAt) {
		return ErrExpired
	}

	o.Attempts++

	if !o.inProximity(lat, lng) {
		return ErrNotInProximity
	}
	if hashCode(code) != o.CodeHash {
		return ErrInvalidCode
	}

	o.Verified = true
	return nil
}

// IsExpired reports whether the OTP has passed its expiry time.
func (o *OTP) IsExpired() bool {
	return time.Now().After(o.ExpiresAt)
}

// generateCode produces a cryptographically random 6-digit numeric code.
func generateCode() (string, error) {
	maxVal := new(big.Int).Exp(big.NewInt(10), big.NewInt(otpLength), nil)
	n, err := rand.Int(rand.Reader, maxVal)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%0*d", otpLength, n), nil
}

// hashCode returns the SHA-256 hex digest of the given code.
func hashCode(code string) string {
	h := sha256.Sum256([]byte(code))
	return hex.EncodeToString(h[:])
}

// inProximity checks whether (lat, lng) is within the required distance of the OTP location.
func (o *OTP) inProximity(lat, lng float64) bool {
	return haversineMeters(o.Lat, o.Lng, lat, lng) <= proximityMeters
}

// haversineMeters computes the great-circle distance in meters between two coordinates.
func haversineMeters(lat1, lng1, lat2, lng2 float64) float64 {
	dLat := degreesToRadians(lat2 - lat1)
	dLng := degreesToRadians(lng2 - lng1)
	a := math.Sin(dLat/2)*math.Sin(dLat/2) +
		math.Cos(degreesToRadians(lat1))*math.Cos(degreesToRadians(lat2))*
			math.Sin(dLng/2)*math.Sin(dLng/2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
	return earthRadiusKm * c * 1000
}

func degreesToRadians(deg float64) float64 {
	return deg * math.Pi / 180
}
