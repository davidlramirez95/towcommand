package otpuc

import (
	"context"
	"errors"
	"io"
	"log/slog"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/davidlramirez95/towcommand/internal/domain/booking"
	domainerrors "github.com/davidlramirez95/towcommand/internal/domain/errors"
	"github.com/davidlramirez95/towcommand/internal/domain/otp"
	"github.com/davidlramirez95/towcommand/internal/domain/user"
	"github.com/davidlramirez95/towcommand/internal/usecase/port"
)

// --- Mocks ---

type mockBookingFinder struct{ mock.Mock }

func (m *mockBookingFinder) FindByID(ctx context.Context, bookingID string) (*booking.Booking, error) {
	args := m.Called(ctx, bookingID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*booking.Booking), args.Error(1)
}

type mockUserFinder struct{ mock.Mock }

func (m *mockUserFinder) FindByID(ctx context.Context, userID string) (*user.User, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*user.User), args.Error(1)
}

type mockOTPCache struct{ mock.Mock }

func (m *mockOTPCache) StoreOTP(ctx context.Context, bookingID, otpType, hashedOTP string, ttl time.Duration) error {
	args := m.Called(ctx, bookingID, otpType, hashedOTP, ttl)
	return args.Error(0)
}

func (m *mockOTPCache) GetOTP(ctx context.Context, bookingID, otpType string) (string, error) {
	args := m.Called(ctx, bookingID, otpType)
	return args.String(0), args.Error(1)
}

func (m *mockOTPCache) DeleteOTP(ctx context.Context, bookingID, otpType string) error {
	args := m.Called(ctx, bookingID, otpType)
	return args.Error(0)
}

type mockOTPRepo struct{ mock.Mock }

func (m *mockOTPRepo) Save(ctx context.Context, o *otp.OTP) error {
	args := m.Called(ctx, o)
	return args.Error(0)
}

func (m *mockOTPRepo) FindByBookingAndType(ctx context.Context, bookingID, otpType string) (*otp.OTP, error) {
	args := m.Called(ctx, bookingID, otpType)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*otp.OTP), args.Error(1)
}

func (m *mockOTPRepo) MarkVerified(ctx context.Context, bookingID, otpType string) error {
	args := m.Called(ctx, bookingID, otpType)
	return args.Error(0)
}

type mockRateLimiter struct{ mock.Mock }

func (m *mockRateLimiter) CheckRateLimit(ctx context.Context, key string, maxRequests, windowSec int) (allowed bool, remaining int, err error) {
	args := m.Called(ctx, key, maxRequests, windowSec)
	return args.Bool(0), args.Int(1), args.Error(2)
}

type mockSMSSender struct{ mock.Mock }

func (m *mockSMSSender) SendSMS(ctx context.Context, phoneNumber, message string) error {
	args := m.Called(ctx, phoneNumber, message)
	return args.Error(0)
}

type mockEventPublisher struct{ mock.Mock }

func (m *mockEventPublisher) Publish(ctx context.Context, source, detailType string, detail any, actor *port.Actor) error {
	args := m.Called(ctx, source, detailType, detail, actor)
	return args.Error(0)
}

func testLogger() *slog.Logger {
	return slog.New(slog.NewJSONHandler(io.Discard, nil))
}

// --- Tests ---

func TestGenerateOTPUseCase_Execute_Success(t *testing.T) {
	bookingsMock := new(mockBookingFinder)
	usersMock := new(mockUserFinder)
	cacheMock := new(mockOTPCache)
	repoMock := new(mockOTPRepo)
	rateMock := new(mockRateLimiter)
	smsMock := new(mockSMSSender)
	eventsMock := new(mockEventPublisher)

	uc := NewGenerateOTPUseCase(bookingsMock, usersMock, cacheMock, repoMock, rateMock, smsMock, eventsMock, testLogger())
	uc.idGen = func() string { return "OTP-TEST-001" }

	b := &booking.Booking{
		BookingID:  "BK-001",
		CustomerID: "user-123",
		Status:     booking.BookingStatusArrived,
	}
	u := &user.User{
		UserID: "user-123",
		Phone:  "+639171234567",
	}

	bookingsMock.On("FindByID", mock.Anything, "BK-001").Return(b, nil)
	rateMock.On("CheckRateLimit", mock.Anything, "otp:generate:BK-001", 3, 300).Return(true, 2, nil)
	cacheMock.On("StoreOTP", mock.Anything, "BK-001", "PICKUP", mock.AnythingOfType("string"), 5*time.Minute).Return(nil)
	repoMock.On("Save", mock.Anything, mock.AnythingOfType("*otp.OTP")).Return(nil)
	usersMock.On("FindByID", mock.Anything, "user-123").Return(u, nil)
	smsMock.On("SendSMS", mock.Anything, "+639171234567", mock.MatchedBy(func(msg string) bool {
		return msg != "" && assert.Contains(t, msg, "TowCommand verification code")
	})).Return(nil)
	eventsMock.On("Publish", mock.Anything, eventSourceOTP, eventOTPGenerated, mock.Anything, mock.Anything).Return(nil)

	result, err := uc.Execute(context.Background(), &GenerateOTPInput{
		BookingID: "BK-001",
		OTPType:   otp.OTPTypePickup,
		Lat:       14.5995,
		Lng:       120.9842,
		CallerID:  "prov-001",
	})

	require.NoError(t, err)
	assert.True(t, result.Success)
	assert.Equal(t, "BK-001", result.BookingID)
	assert.Equal(t, "PICKUP", result.OTPType)

	bookingsMock.AssertExpectations(t)
	rateMock.AssertExpectations(t)
	cacheMock.AssertExpectations(t)
	repoMock.AssertExpectations(t)
	usersMock.AssertExpectations(t)
	smsMock.AssertExpectations(t)
	eventsMock.AssertExpectations(t)
}

func TestGenerateOTPUseCase_Execute_DropoffSuccess(t *testing.T) {
	bookingsMock := new(mockBookingFinder)
	usersMock := new(mockUserFinder)
	cacheMock := new(mockOTPCache)
	repoMock := new(mockOTPRepo)
	rateMock := new(mockRateLimiter)
	smsMock := new(mockSMSSender)
	eventsMock := new(mockEventPublisher)

	uc := NewGenerateOTPUseCase(bookingsMock, usersMock, cacheMock, repoMock, rateMock, smsMock, eventsMock, testLogger())
	uc.idGen = func() string { return "OTP-TEST-002" }

	b := &booking.Booking{
		BookingID:  "BK-002",
		CustomerID: "user-456",
		Status:     booking.BookingStatusArrivedDropoff,
	}
	u := &user.User{UserID: "user-456", Phone: "+639179876543"}

	bookingsMock.On("FindByID", mock.Anything, "BK-002").Return(b, nil)
	rateMock.On("CheckRateLimit", mock.Anything, mock.Anything, 3, 300).Return(true, 2, nil)
	cacheMock.On("StoreOTP", mock.Anything, "BK-002", "DROPOFF", mock.Anything, 5*time.Minute).Return(nil)
	repoMock.On("Save", mock.Anything, mock.Anything).Return(nil)
	usersMock.On("FindByID", mock.Anything, "user-456").Return(u, nil)
	smsMock.On("SendSMS", mock.Anything, "+639179876543", mock.Anything).Return(nil)
	eventsMock.On("Publish", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)

	result, err := uc.Execute(context.Background(), &GenerateOTPInput{
		BookingID: "BK-002",
		OTPType:   otp.OTPTypeDropoff,
		Lat:       14.5547,
		Lng:       121.0244,
		CallerID:  "prov-002",
	})

	require.NoError(t, err)
	assert.True(t, result.Success)
	assert.Equal(t, "DROPOFF", result.OTPType)
}

func TestGenerateOTPUseCase_Execute_RateLimited(t *testing.T) {
	bookingsMock := new(mockBookingFinder)
	usersMock := new(mockUserFinder)
	cacheMock := new(mockOTPCache)
	repoMock := new(mockOTPRepo)
	rateMock := new(mockRateLimiter)
	smsMock := new(mockSMSSender)
	eventsMock := new(mockEventPublisher)

	uc := NewGenerateOTPUseCase(bookingsMock, usersMock, cacheMock, repoMock, rateMock, smsMock, eventsMock, testLogger())

	b := &booking.Booking{BookingID: "BK-003", CustomerID: "user-789", Status: booking.BookingStatusArrived}
	bookingsMock.On("FindByID", mock.Anything, "BK-003").Return(b, nil)
	rateMock.On("CheckRateLimit", mock.Anything, "otp:generate:BK-003", 3, 300).Return(false, 0, nil)

	_, err := uc.Execute(context.Background(), &GenerateOTPInput{
		BookingID: "BK-003",
		OTPType:   otp.OTPTypePickup,
		Lat:       14.5995,
		Lng:       120.9842,
		CallerID:  "prov-003",
	})

	require.Error(t, err)
	var appErr *domainerrors.AppError
	require.True(t, errors.As(err, &appErr))
	assert.Equal(t, domainerrors.CodeRateLimited, appErr.Code)
}

func TestGenerateOTPUseCase_Execute_WrongBookingStatus(t *testing.T) {
	bookingsMock := new(mockBookingFinder)
	usersMock := new(mockUserFinder)
	cacheMock := new(mockOTPCache)
	repoMock := new(mockOTPRepo)
	rateMock := new(mockRateLimiter)
	smsMock := new(mockSMSSender)
	eventsMock := new(mockEventPublisher)

	uc := NewGenerateOTPUseCase(bookingsMock, usersMock, cacheMock, repoMock, rateMock, smsMock, eventsMock, testLogger())

	// Booking is PENDING, not ARRIVED — should fail.
	b := &booking.Booking{BookingID: "BK-004", CustomerID: "user-100", Status: booking.BookingStatusPending}
	bookingsMock.On("FindByID", mock.Anything, "BK-004").Return(b, nil)

	_, err := uc.Execute(context.Background(), &GenerateOTPInput{
		BookingID: "BK-004",
		OTPType:   otp.OTPTypePickup,
		Lat:       14.5995,
		Lng:       120.9842,
		CallerID:  "prov-004",
	})

	require.Error(t, err)
	var appErr *domainerrors.AppError
	require.True(t, errors.As(err, &appErr))
	assert.Equal(t, domainerrors.CodeConflict, appErr.Code)
}

func TestGenerateOTPUseCase_Execute_BookingNotFound(t *testing.T) {
	bookingsMock := new(mockBookingFinder)
	usersMock := new(mockUserFinder)
	cacheMock := new(mockOTPCache)
	repoMock := new(mockOTPRepo)
	rateMock := new(mockRateLimiter)
	smsMock := new(mockSMSSender)
	eventsMock := new(mockEventPublisher)

	uc := NewGenerateOTPUseCase(bookingsMock, usersMock, cacheMock, repoMock, rateMock, smsMock, eventsMock, testLogger())

	bookingsMock.On("FindByID", mock.Anything, "BK-GONE").Return(nil, nil)

	_, err := uc.Execute(context.Background(), &GenerateOTPInput{
		BookingID: "BK-GONE",
		OTPType:   otp.OTPTypePickup,
		Lat:       14.5995,
		Lng:       120.9842,
		CallerID:  "prov-005",
	})

	require.Error(t, err)
	var appErr *domainerrors.AppError
	require.True(t, errors.As(err, &appErr))
	assert.Equal(t, domainerrors.CodeNotFound, appErr.Code)
}

func TestGenerateOTPUseCase_Execute_SMSFailureNonFatal(t *testing.T) {
	bookingsMock := new(mockBookingFinder)
	usersMock := new(mockUserFinder)
	cacheMock := new(mockOTPCache)
	repoMock := new(mockOTPRepo)
	rateMock := new(mockRateLimiter)
	smsMock := new(mockSMSSender)
	eventsMock := new(mockEventPublisher)

	uc := NewGenerateOTPUseCase(bookingsMock, usersMock, cacheMock, repoMock, rateMock, smsMock, eventsMock, testLogger())
	uc.idGen = func() string { return "OTP-TEST-SMS" }

	b := &booking.Booking{BookingID: "BK-005", CustomerID: "user-sms", Status: booking.BookingStatusArrived}
	u := &user.User{UserID: "user-sms", Phone: "+639170000000"}

	bookingsMock.On("FindByID", mock.Anything, "BK-005").Return(b, nil)
	rateMock.On("CheckRateLimit", mock.Anything, mock.Anything, 3, 300).Return(true, 1, nil)
	cacheMock.On("StoreOTP", mock.Anything, "BK-005", "PICKUP", mock.Anything, 5*time.Minute).Return(nil)
	repoMock.On("Save", mock.Anything, mock.Anything).Return(nil)
	usersMock.On("FindByID", mock.Anything, "user-sms").Return(u, nil)
	smsMock.On("SendSMS", mock.Anything, "+639170000000", mock.Anything).Return(domainerrors.NewExternalServiceError("SNS", nil))
	eventsMock.On("Publish", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)

	result, err := uc.Execute(context.Background(), &GenerateOTPInput{
		BookingID: "BK-005",
		OTPType:   otp.OTPTypePickup,
		Lat:       14.5995,
		Lng:       120.9842,
		CallerID:  "prov-sms",
	})

	// SMS failure should not fail the use case.
	require.NoError(t, err)
	assert.True(t, result.Success)
}

func TestGenerateOTPUseCase_Execute_DDBBackupFailureNonFatal(t *testing.T) {
	bookingsMock := new(mockBookingFinder)
	usersMock := new(mockUserFinder)
	cacheMock := new(mockOTPCache)
	repoMock := new(mockOTPRepo)
	rateMock := new(mockRateLimiter)
	smsMock := new(mockSMSSender)
	eventsMock := new(mockEventPublisher)

	uc := NewGenerateOTPUseCase(bookingsMock, usersMock, cacheMock, repoMock, rateMock, smsMock, eventsMock, testLogger())
	uc.idGen = func() string { return "OTP-TEST-DDB" }

	b := &booking.Booking{BookingID: "BK-006", CustomerID: "user-ddb", Status: booking.BookingStatusArrived}
	u := &user.User{UserID: "user-ddb", Phone: "+639171111111"}

	bookingsMock.On("FindByID", mock.Anything, "BK-006").Return(b, nil)
	rateMock.On("CheckRateLimit", mock.Anything, mock.Anything, 3, 300).Return(true, 1, nil)
	cacheMock.On("StoreOTP", mock.Anything, "BK-006", "PICKUP", mock.Anything, 5*time.Minute).Return(nil)
	repoMock.On("Save", mock.Anything, mock.Anything).Return(domainerrors.NewInternalError("DDB write failed"))
	usersMock.On("FindByID", mock.Anything, "user-ddb").Return(u, nil)
	smsMock.On("SendSMS", mock.Anything, "+639171111111", mock.Anything).Return(nil)
	eventsMock.On("Publish", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)

	result, err := uc.Execute(context.Background(), &GenerateOTPInput{
		BookingID: "BK-006",
		OTPType:   otp.OTPTypePickup,
		Lat:       14.5995,
		Lng:       120.9842,
		CallerID:  "prov-ddb",
	})

	// DynamoDB backup failure should not fail the use case.
	require.NoError(t, err)
	assert.True(t, result.Success)
}
