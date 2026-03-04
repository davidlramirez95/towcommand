package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/davidlramirez95/towcommand/internal/domain/otp"
)

// otpRecord is the DynamoDB storage representation of an OTP entity.
type otpRecord struct {
	PK         string    `json:"PK"`
	SK         string    `json:"SK"`
	EntityType string    `json:"entityType"`
	OTPID      string    `json:"otpId"`
	BookingID  string    `json:"bookingId"`
	Type       string    `json:"type"`
	CodeHash   string    `json:"codeHash"`
	Lat        float64   `json:"lat"`
	Lng        float64   `json:"lng"`
	Attempts   int       `json:"attempts"`
	ExpiresAt  time.Time `json:"expiresAt"`
	Verified   bool      `json:"verified"`
	CreatedAt  time.Time `json:"createdAt"`
}

// DynamoOTPRepository implements OTP backup persistence against DynamoDB.
// Key schema: PK=JOB#{bookingId}, SK=OTP#{type}
type DynamoOTPRepository struct {
	baseRepository
}

// NewOTPRepository creates a new DynamoDB-backed OTP repository.
func NewOTPRepository(client DynamoDBAPI, tableName string) *DynamoOTPRepository {
	return &DynamoOTPRepository{
		baseRepository: baseRepository{client: client, tableName: tableName},
	}
}

// Save persists an OTP record as a backup for Redis.
func (r *DynamoOTPRepository) Save(ctx context.Context, o *otp.OTP) error {
	rec := otpRecord{
		PK:         PrefixJob + o.BookingID,
		SK:         PrefixOTP + string(o.Type),
		EntityType: "OTP",
		OTPID:      o.OTPID,
		BookingID:  o.BookingID,
		Type:       string(o.Type),
		CodeHash:   o.CodeHash,
		Lat:        o.Lat,
		Lng:        o.Lng,
		Attempts:   o.Attempts,
		ExpiresAt:  o.ExpiresAt,
		Verified:   o.Verified,
		CreatedAt:  o.CreatedAt,
	}

	item, err := marshalItem(rec)
	if err != nil {
		return fmt.Errorf("marshal OTP record: %w", err)
	}
	return r.putItem(ctx, item)
}

// FindByBookingAndType retrieves an OTP record by booking ID and OTP type.
// Returns nil if not found.
func (r *DynamoOTPRepository) FindByBookingAndType(ctx context.Context, bookingID, otpType string) (*otp.OTP, error) {
	var rec otpRecord
	found, err := r.getItem(ctx, PrefixJob+bookingID, PrefixOTP+otpType, &rec)
	if err != nil {
		return nil, err
	}
	if !found {
		return nil, nil
	}
	return &otp.OTP{
		OTPID:     rec.OTPID,
		BookingID: rec.BookingID,
		Type:      otp.OTPType(rec.Type),
		CodeHash:  rec.CodeHash,
		Lat:       rec.Lat,
		Lng:       rec.Lng,
		Attempts:  rec.Attempts,
		ExpiresAt: rec.ExpiresAt,
		Verified:  rec.Verified,
		CreatedAt: rec.CreatedAt,
	}, nil
}

// MarkVerified marks an OTP record as verified in DynamoDB.
func (r *DynamoOTPRepository) MarkVerified(ctx context.Context, bookingID, otpType string) error {
	return r.updateItem(ctx, PrefixJob+bookingID, PrefixOTP+otpType, map[string]any{
		"verified": true,
	})
}
