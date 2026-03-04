package otpuc

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/davidlramirez95/towcommand/internal/domain/booking"
	domainerrors "github.com/davidlramirez95/towcommand/internal/domain/errors"
	"github.com/davidlramirez95/towcommand/internal/domain/otp"
)

// --- Mocks (reusing mockBookingFinder, mockOTPCache, mockOTPRepo, mockEventPublisher from generate_test.go) ---

type mockBookingStatusUpdater struct{ mock.Mock }

func (m *mockBookingStatusUpdater) UpdateStatus(ctx context.Context, bookingID string, status booking.BookingStatus, metadata map[string]any) error {
	args := m.Called(ctx, bookingID, status, metadata)
	return args.Error(0)
}

// hashTestCode produces a SHA-256 hash matching the domain's hashCode function.
func hashTestCode(code string) string {
	h := sha256.Sum256([]byte(code))
	return hex.EncodeToString(h[:])
}

// newTestOTPRecord creates a valid OTP record for testing.
func newTestOTPRecord(bookingID string, otpType otp.OTPType, code string, lat, lng float64) *otp.OTP {
	return &otp.OTP{
		OTPID:     "OTP-TEST",
		BookingID: bookingID,
		Type:      otpType,
		CodeHash:  hashTestCode(code),
		Lat:       lat,
		Lng:       lng,
		Attempts:  0,
		ExpiresAt: time.Now().Add(5 * time.Minute),
		Verified:  false,
		CreatedAt: time.Now(),
	}
}

// --- Tests ---

func TestVerifyOTPUseCase_Execute_PickupSuccess(t *testing.T) {
	bookingsMock := new(mockBookingFinder)
	statusesMock := new(mockBookingStatusUpdater)
	cacheMock := new(mockOTPCache)
	repoMock := new(mockOTPRepo)
	eventsMock := new(mockEventPublisher)

	uc := NewVerifyOTPUseCase(bookingsMock, statusesMock, cacheMock, repoMock, eventsMock, testLogger())

	b := &booking.Booking{
		BookingID:  "BK-001",
		CustomerID: "user-123",
		ProviderID: "prov-001",
		Status:     booking.BookingStatusConditionReport,
	}
	otpRec := newTestOTPRecord("BK-001", otp.OTPTypePickup, "123456", 14.5995, 120.9842)

	bookingsMock.On("FindByID", mock.Anything, "BK-001").Return(b, nil)
	cacheMock.On("GetOTP", mock.Anything, "BK-001", "PICKUP").Return(otpRec.CodeHash, nil)
	repoMock.On("FindByBookingAndType", mock.Anything, "BK-001", "PICKUP").Return(otpRec, nil)
	statusesMock.On("UpdateStatus", mock.Anything, "BK-001", booking.BookingStatusOTPVerified, mock.Anything).Return(nil)
	cacheMock.On("DeleteOTP", mock.Anything, "BK-001", "PICKUP").Return(nil)
	repoMock.On("MarkVerified", mock.Anything, "BK-001", "PICKUP").Return(nil)
	eventsMock.On("Publish", mock.Anything, eventSourceOTP, eventOTPVerified, mock.Anything, mock.Anything).Return(nil)

	result, err := uc.Execute(context.Background(), &VerifyOTPInput{
		BookingID: "BK-001",
		OTPType:   otp.OTPTypePickup,
		Code:      "123456",
		Lat:       14.5995,
		Lng:       120.9842,
		CallerID:  "prov-001",
	})

	require.NoError(t, err)
	assert.True(t, result.Success)
	assert.Equal(t, "BK-001", result.BookingID)
	assert.Equal(t, "PICKUP", result.OTPType)
	assert.Equal(t, "OTP_VERIFIED", result.NewStatus)

	bookingsMock.AssertExpectations(t)
	statusesMock.AssertExpectations(t)
	cacheMock.AssertExpectations(t)
	repoMock.AssertExpectations(t)
	eventsMock.AssertExpectations(t)
}

func TestVerifyOTPUseCase_Execute_DropoffSuccess(t *testing.T) {
	bookingsMock := new(mockBookingFinder)
	statusesMock := new(mockBookingStatusUpdater)
	cacheMock := new(mockOTPCache)
	repoMock := new(mockOTPRepo)
	eventsMock := new(mockEventPublisher)

	uc := NewVerifyOTPUseCase(bookingsMock, statusesMock, cacheMock, repoMock, eventsMock, testLogger())

	b := &booking.Booking{
		BookingID:  "BK-002",
		CustomerID: "user-456",
		ProviderID: "prov-002",
		Status:     booking.BookingStatusArrivedDropoff,
	}
	otpRec := newTestOTPRecord("BK-002", otp.OTPTypeDropoff, "654321", 14.5547, 121.0244)

	bookingsMock.On("FindByID", mock.Anything, "BK-002").Return(b, nil)
	cacheMock.On("GetOTP", mock.Anything, "BK-002", "DROPOFF").Return(otpRec.CodeHash, nil)
	repoMock.On("FindByBookingAndType", mock.Anything, "BK-002", "DROPOFF").Return(otpRec, nil)
	statusesMock.On("UpdateStatus", mock.Anything, "BK-002", booking.BookingStatusOTPDropoff, mock.Anything).Return(nil)
	cacheMock.On("DeleteOTP", mock.Anything, "BK-002", "DROPOFF").Return(nil)
	repoMock.On("MarkVerified", mock.Anything, "BK-002", "DROPOFF").Return(nil)
	eventsMock.On("Publish", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)

	result, err := uc.Execute(context.Background(), &VerifyOTPInput{
		BookingID: "BK-002",
		OTPType:   otp.OTPTypeDropoff,
		Code:      "654321",
		Lat:       14.5547,
		Lng:       121.0244,
		CallerID:  "prov-002",
	})

	require.NoError(t, err)
	assert.True(t, result.Success)
	assert.Equal(t, "OTP_DROPOFF", result.NewStatus)
}

func TestVerifyOTPUseCase_Execute_InvalidCode(t *testing.T) {
	bookingsMock := new(mockBookingFinder)
	statusesMock := new(mockBookingStatusUpdater)
	cacheMock := new(mockOTPCache)
	repoMock := new(mockOTPRepo)
	eventsMock := new(mockEventPublisher)

	uc := NewVerifyOTPUseCase(bookingsMock, statusesMock, cacheMock, repoMock, eventsMock, testLogger())

	b := &booking.Booking{
		BookingID: "BK-003",
		Status:    booking.BookingStatusConditionReport,
	}
	otpRec := newTestOTPRecord("BK-003", otp.OTPTypePickup, "123456", 14.5995, 120.9842)

	bookingsMock.On("FindByID", mock.Anything, "BK-003").Return(b, nil)
	cacheMock.On("GetOTP", mock.Anything, "BK-003", "PICKUP").Return(otpRec.CodeHash, nil)
	repoMock.On("FindByBookingAndType", mock.Anything, "BK-003", "PICKUP").Return(otpRec, nil)

	_, err := uc.Execute(context.Background(), &VerifyOTPInput{
		BookingID: "BK-003",
		OTPType:   otp.OTPTypePickup,
		Code:      "999999", // wrong code
		Lat:       14.5995,
		Lng:       120.9842,
		CallerID:  "prov-003",
	})

	require.Error(t, err)
	var appErr *domainerrors.AppError
	require.True(t, errors.As(err, &appErr))
	assert.Equal(t, domainerrors.CodeOTPInvalid, appErr.Code)
}

func TestVerifyOTPUseCase_Execute_ExpiredOTP(t *testing.T) {
	bookingsMock := new(mockBookingFinder)
	statusesMock := new(mockBookingStatusUpdater)
	cacheMock := new(mockOTPCache)
	repoMock := new(mockOTPRepo)
	eventsMock := new(mockEventPublisher)

	uc := NewVerifyOTPUseCase(bookingsMock, statusesMock, cacheMock, repoMock, eventsMock, testLogger())

	b := &booking.Booking{
		BookingID: "BK-004",
		Status:    booking.BookingStatusConditionReport,
	}
	// Create an expired OTP record.
	expiredOTP := &otp.OTP{
		OTPID:     "OTP-EXPIRED",
		BookingID: "BK-004",
		Type:      otp.OTPTypePickup,
		CodeHash:  hashTestCode("111111"),
		Lat:       14.5995,
		Lng:       120.9842,
		Attempts:  0,
		ExpiresAt: time.Now().Add(-1 * time.Minute), // expired
		Verified:  false,
		CreatedAt: time.Now().Add(-10 * time.Minute),
	}

	bookingsMock.On("FindByID", mock.Anything, "BK-004").Return(b, nil)
	cacheMock.On("GetOTP", mock.Anything, "BK-004", "PICKUP").Return("", nil) // Redis miss
	repoMock.On("FindByBookingAndType", mock.Anything, "BK-004", "PICKUP").Return(expiredOTP, nil)

	_, err := uc.Execute(context.Background(), &VerifyOTPInput{
		BookingID: "BK-004",
		OTPType:   otp.OTPTypePickup,
		Code:      "111111",
		Lat:       14.5995,
		Lng:       120.9842,
		CallerID:  "prov-004",
	})

	require.Error(t, err)
	var appErr *domainerrors.AppError
	require.True(t, errors.As(err, &appErr))
	assert.Equal(t, domainerrors.CodeOTPExpired, appErr.Code)
}

func TestVerifyOTPUseCase_Execute_RedisMissWithDDBFallback(t *testing.T) {
	bookingsMock := new(mockBookingFinder)
	statusesMock := new(mockBookingStatusUpdater)
	cacheMock := new(mockOTPCache)
	repoMock := new(mockOTPRepo)
	eventsMock := new(mockEventPublisher)

	uc := NewVerifyOTPUseCase(bookingsMock, statusesMock, cacheMock, repoMock, eventsMock, testLogger())

	b := &booking.Booking{
		BookingID:  "BK-005",
		CustomerID: "user-500",
		Status:     booking.BookingStatusConditionReport,
	}
	otpRec := newTestOTPRecord("BK-005", otp.OTPTypePickup, "555555", 14.5995, 120.9842)

	bookingsMock.On("FindByID", mock.Anything, "BK-005").Return(b, nil)
	// Redis returns empty (cache miss / expired).
	cacheMock.On("GetOTP", mock.Anything, "BK-005", "PICKUP").Return("", nil)
	// DynamoDB has the record.
	repoMock.On("FindByBookingAndType", mock.Anything, "BK-005", "PICKUP").Return(otpRec, nil)
	statusesMock.On("UpdateStatus", mock.Anything, "BK-005", booking.BookingStatusOTPVerified, mock.Anything).Return(nil)
	cacheMock.On("DeleteOTP", mock.Anything, "BK-005", "PICKUP").Return(nil)
	repoMock.On("MarkVerified", mock.Anything, "BK-005", "PICKUP").Return(nil)
	eventsMock.On("Publish", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)

	result, err := uc.Execute(context.Background(), &VerifyOTPInput{
		BookingID: "BK-005",
		OTPType:   otp.OTPTypePickup,
		Code:      "555555",
		Lat:       14.5995,
		Lng:       120.9842,
		CallerID:  "prov-005",
	})

	require.NoError(t, err)
	assert.True(t, result.Success)
	assert.Equal(t, "OTP_VERIFIED", result.NewStatus)
}

func TestVerifyOTPUseCase_Execute_OTPNotFound(t *testing.T) {
	bookingsMock := new(mockBookingFinder)
	statusesMock := new(mockBookingStatusUpdater)
	cacheMock := new(mockOTPCache)
	repoMock := new(mockOTPRepo)
	eventsMock := new(mockEventPublisher)

	uc := NewVerifyOTPUseCase(bookingsMock, statusesMock, cacheMock, repoMock, eventsMock, testLogger())

	b := &booking.Booking{BookingID: "BK-006", Status: booking.BookingStatusConditionReport}

	bookingsMock.On("FindByID", mock.Anything, "BK-006").Return(b, nil)
	cacheMock.On("GetOTP", mock.Anything, "BK-006", "PICKUP").Return("", nil)
	repoMock.On("FindByBookingAndType", mock.Anything, "BK-006", "PICKUP").Return(nil, nil)

	_, err := uc.Execute(context.Background(), &VerifyOTPInput{
		BookingID: "BK-006",
		OTPType:   otp.OTPTypePickup,
		Code:      "000000",
		Lat:       14.5995,
		Lng:       120.9842,
		CallerID:  "prov-006",
	})

	require.Error(t, err)
	var appErr *domainerrors.AppError
	require.True(t, errors.As(err, &appErr))
	assert.Equal(t, domainerrors.CodeNotFound, appErr.Code)
}

func TestVerifyOTPUseCase_Execute_WrongBookingStatus(t *testing.T) {
	bookingsMock := new(mockBookingFinder)
	statusesMock := new(mockBookingStatusUpdater)
	cacheMock := new(mockOTPCache)
	repoMock := new(mockOTPRepo)
	eventsMock := new(mockEventPublisher)

	uc := NewVerifyOTPUseCase(bookingsMock, statusesMock, cacheMock, repoMock, eventsMock, testLogger())

	// Booking is in ARRIVED, but verification needs CONDITION_REPORT.
	b := &booking.Booking{BookingID: "BK-007", Status: booking.BookingStatusArrived}
	bookingsMock.On("FindByID", mock.Anything, "BK-007").Return(b, nil)

	_, err := uc.Execute(context.Background(), &VerifyOTPInput{
		BookingID: "BK-007",
		OTPType:   otp.OTPTypePickup,
		Code:      "123456",
		Lat:       14.5995,
		Lng:       120.9842,
		CallerID:  "prov-007",
	})

	require.Error(t, err)
	var appErr *domainerrors.AppError
	require.True(t, errors.As(err, &appErr))
	assert.Equal(t, domainerrors.CodeConflict, appErr.Code)
}

func TestVerifyOTPUseCase_Execute_BookingNotFound(t *testing.T) {
	bookingsMock := new(mockBookingFinder)
	statusesMock := new(mockBookingStatusUpdater)
	cacheMock := new(mockOTPCache)
	repoMock := new(mockOTPRepo)
	eventsMock := new(mockEventPublisher)

	uc := NewVerifyOTPUseCase(bookingsMock, statusesMock, cacheMock, repoMock, eventsMock, testLogger())

	bookingsMock.On("FindByID", mock.Anything, "BK-MISSING").Return(nil, nil)

	_, err := uc.Execute(context.Background(), &VerifyOTPInput{
		BookingID: "BK-MISSING",
		OTPType:   otp.OTPTypePickup,
		Code:      "123456",
		Lat:       14.5995,
		Lng:       120.9842,
		CallerID:  "prov-008",
	})

	require.Error(t, err)
	var appErr *domainerrors.AppError
	require.True(t, errors.As(err, &appErr))
	assert.Equal(t, domainerrors.CodeNotFound, appErr.Code)
}

func TestMapOTPDomainError(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		wantCode domainerrors.ErrorCode
	}{
		{"expired", otp.ErrExpired, domainerrors.CodeOTPExpired},
		{"invalid code", otp.ErrInvalidCode, domainerrors.CodeOTPInvalid},
		{"max attempts", otp.ErrMaxAttempts, domainerrors.CodeOTPInvalid},
		{"already verified", otp.ErrAlreadyVerified, domainerrors.CodeConflict},
		{"not in proximity", otp.ErrNotInProximity, domainerrors.CodeValidationError},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := mapOTPDomainError(tt.err)
			var appErr *domainerrors.AppError
			require.True(t, errors.As(result, &appErr))
			assert.Equal(t, tt.wantCode, appErr.Code)
		})
	}
}
