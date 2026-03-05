package paymentuc

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/davidlramirez95/towcommand/internal/domain/booking"
	domainerrors "github.com/davidlramirez95/towcommand/internal/domain/errors"
	"github.com/davidlramirez95/towcommand/internal/domain/payment"
	"github.com/davidlramirez95/towcommand/internal/usecase/port"
)

// --- Mocks ---

type mockBookingFinder struct{ mock.Mock }

func (m *mockBookingFinder) FindByID(ctx context.Context, bookingID string) (*booking.Booking, error) {
	args := m.Called(ctx, bookingID)
	if v := args.Get(0); v != nil {
		return v.(*booking.Booking), args.Error(1)
	}
	return nil, args.Error(1)
}

type mockPaymentByBookingLister struct{ mock.Mock }

func (m *mockPaymentByBookingLister) FindByBooking(ctx context.Context, bookingID string) ([]payment.Payment, error) {
	args := m.Called(ctx, bookingID)
	if v := args.Get(0); v != nil {
		return v.([]payment.Payment), args.Error(1)
	}
	return nil, args.Error(1)
}

type mockPaymentSaver struct{ mock.Mock }

func (m *mockPaymentSaver) Save(ctx context.Context, p *payment.Payment) error {
	args := m.Called(ctx, p)
	return args.Error(0)
}

type mockPaymentGateway struct{ mock.Mock }

func (m *mockPaymentGateway) Charge(ctx context.Context, paymentID string, amountCentavos int64, currency, method string) (*port.ChargeResult, error) {
	args := m.Called(ctx, paymentID, amountCentavos, currency, method)
	if v := args.Get(0); v != nil {
		return v.(*port.ChargeResult), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *mockPaymentGateway) Refund(ctx context.Context, gatewayRef string, amountCentavos int64) (*port.RefundResult, error) {
	args := m.Called(ctx, gatewayRef, amountCentavos)
	if v := args.Get(0); v != nil {
		return v.(*port.RefundResult), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *mockPaymentGateway) VerifyWebhookSignature(payload []byte, signature string) error {
	args := m.Called(payload, signature)
	return args.Error(0)
}

type mockEventPublisher struct{ mock.Mock }

func (m *mockEventPublisher) Publish(ctx context.Context, source, detailType string, detail any, actor *Actor) error {
	args := m.Called(ctx, source, detailType, detail, actor)
	return args.Error(0)
}

// --- Helpers ---

func completedBooking() *booking.Booking {
	return &booking.Booking{
		BookingID:  "BK-001",
		CustomerID: "user-123",
		ProviderID: "prov-456",
		Status:     booking.BookingStatusCompleted,
		Price: booking.PriceBreakdown{
			Total:    200_000,
			Currency: "PHP",
		},
	}
}

var fixedTime = time.Date(2026, 3, 5, 10, 0, 0, 0, time.UTC)

// --- Tests ---

func TestInitiatePayment_CashAutoCapture(t *testing.T) {
	bookingMock := new(mockBookingFinder)
	paymentLister := new(mockPaymentByBookingLister)
	saverMock := new(mockPaymentSaver)
	gwMock := new(mockPaymentGateway)
	eventsMock := new(mockEventPublisher)

	uc := NewInitiatePaymentUseCase(bookingMock, paymentLister, saverMock, gwMock, eventsMock)
	uc.idGen = func() string { return "PAY-2026-CASH" }
	uc.now = func() time.Time { return fixedTime }

	bookingMock.On("FindByID", mock.Anything, "BK-001").Return(completedBooking(), nil)
	paymentLister.On("FindByBooking", mock.Anything, "BK-001").Return([]payment.Payment{}, nil)
	saverMock.On("Save", mock.Anything, mock.AnythingOfType("*payment.Payment")).Return(nil)
	eventsMock.On("Publish", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)

	result, err := uc.Execute(context.Background(), &InitiatePaymentInput{
		BookingID: "BK-001",
		Method:    payment.PaymentMethodCash,
		UserID:    "user-123",
	})

	require.NoError(t, err)
	assert.Equal(t, "PAY-2026-CASH", result.PaymentID)
	assert.Equal(t, payment.PaymentStatusCaptured, result.Status)
	assert.NotNil(t, result.CapturedAt)
	assert.Equal(t, int64(200_000), result.Amount)
	assert.Empty(t, result.GatewayRef)

	// Verify gateway was NOT called for cash.
	gwMock.AssertNotCalled(t, "Charge")
	saverMock.AssertExpectations(t)
	eventsMock.AssertExpectations(t)
}

func TestInitiatePayment_EWalletPending(t *testing.T) {
	bookingMock := new(mockBookingFinder)
	paymentLister := new(mockPaymentByBookingLister)
	saverMock := new(mockPaymentSaver)
	gwMock := new(mockPaymentGateway)
	eventsMock := new(mockEventPublisher)

	uc := NewInitiatePaymentUseCase(bookingMock, paymentLister, saverMock, gwMock, eventsMock)
	uc.idGen = func() string { return "PAY-2026-GCASH" }
	uc.now = func() time.Time { return fixedTime }

	bookingMock.On("FindByID", mock.Anything, "BK-001").Return(completedBooking(), nil)
	paymentLister.On("FindByBooking", mock.Anything, "BK-001").Return([]payment.Payment{}, nil)
	gwMock.On("Charge", mock.Anything, "PAY-2026-GCASH", int64(200_000), "PHP", "gcash").
		Return(&port.ChargeResult{GatewayRef: "gw-ref-123"}, nil)
	saverMock.On("Save", mock.Anything, mock.AnythingOfType("*payment.Payment")).Return(nil)
	eventsMock.On("Publish", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)

	result, err := uc.Execute(context.Background(), &InitiatePaymentInput{
		BookingID: "BK-001",
		Method:    payment.PaymentMethodGCash,
		UserID:    "user-123",
	})

	require.NoError(t, err)
	assert.Equal(t, "PAY-2026-GCASH", result.PaymentID)
	assert.Equal(t, payment.PaymentStatusPending, result.Status)
	assert.Nil(t, result.CapturedAt)
	assert.Equal(t, "gw-ref-123", result.GatewayRef)

	gwMock.AssertExpectations(t)
	saverMock.AssertExpectations(t)
}

func TestInitiatePayment_BookingNotFound(t *testing.T) {
	bookingMock := new(mockBookingFinder)
	paymentLister := new(mockPaymentByBookingLister)
	saverMock := new(mockPaymentSaver)
	gwMock := new(mockPaymentGateway)
	eventsMock := new(mockEventPublisher)

	uc := NewInitiatePaymentUseCase(bookingMock, paymentLister, saverMock, gwMock, eventsMock)

	bookingMock.On("FindByID", mock.Anything, "BK-MISSING").Return(nil, nil)

	_, err := uc.Execute(context.Background(), &InitiatePaymentInput{
		BookingID: "BK-MISSING",
		Method:    payment.PaymentMethodGCash,
		UserID:    "user-123",
	})

	require.Error(t, err)
	var appErr *domainerrors.AppError
	require.ErrorAs(t, err, &appErr)
	assert.Equal(t, domainerrors.CodeNotFound, appErr.Code)
}

func TestInitiatePayment_BookingNotCompleted(t *testing.T) {
	bookingMock := new(mockBookingFinder)
	paymentLister := new(mockPaymentByBookingLister)
	saverMock := new(mockPaymentSaver)
	gwMock := new(mockPaymentGateway)
	eventsMock := new(mockEventPublisher)

	uc := NewInitiatePaymentUseCase(bookingMock, paymentLister, saverMock, gwMock, eventsMock)

	pendingBooking := completedBooking()
	pendingBooking.Status = booking.BookingStatusPending
	bookingMock.On("FindByID", mock.Anything, "BK-001").Return(pendingBooking, nil)

	_, err := uc.Execute(context.Background(), &InitiatePaymentInput{
		BookingID: "BK-001",
		Method:    payment.PaymentMethodGCash,
		UserID:    "user-123",
	})

	require.Error(t, err)
	var appErr *domainerrors.AppError
	require.ErrorAs(t, err, &appErr)
	assert.Equal(t, domainerrors.CodeValidationError, appErr.Code)
}

func TestInitiatePayment_DuplicatePaymentConflict(t *testing.T) {
	bookingMock := new(mockBookingFinder)
	paymentLister := new(mockPaymentByBookingLister)
	saverMock := new(mockPaymentSaver)
	gwMock := new(mockPaymentGateway)
	eventsMock := new(mockEventPublisher)

	uc := NewInitiatePaymentUseCase(bookingMock, paymentLister, saverMock, gwMock, eventsMock)

	bookingMock.On("FindByID", mock.Anything, "BK-001").Return(completedBooking(), nil)
	paymentLister.On("FindByBooking", mock.Anything, "BK-001").Return([]payment.Payment{
		{PaymentID: "PAY-EXISTING", Status: payment.PaymentStatusCaptured},
	}, nil)

	_, err := uc.Execute(context.Background(), &InitiatePaymentInput{
		BookingID: "BK-001",
		Method:    payment.PaymentMethodCash,
		UserID:    "user-123",
	})

	require.Error(t, err)
	var appErr *domainerrors.AppError
	require.ErrorAs(t, err, &appErr)
	assert.Equal(t, domainerrors.CodeConflict, appErr.Code)
}

func TestInitiatePayment_GatewayChargeFailure(t *testing.T) {
	bookingMock := new(mockBookingFinder)
	paymentLister := new(mockPaymentByBookingLister)
	saverMock := new(mockPaymentSaver)
	gwMock := new(mockPaymentGateway)
	eventsMock := new(mockEventPublisher)

	uc := NewInitiatePaymentUseCase(bookingMock, paymentLister, saverMock, gwMock, eventsMock)
	uc.idGen = func() string { return "PAY-2026-FAIL" }

	bookingMock.On("FindByID", mock.Anything, "BK-001").Return(completedBooking(), nil)
	paymentLister.On("FindByBooking", mock.Anything, "BK-001").Return([]payment.Payment{}, nil)
	gwMock.On("Charge", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
		Return(nil, assert.AnError)

	_, err := uc.Execute(context.Background(), &InitiatePaymentInput{
		BookingID: "BK-001",
		Method:    payment.PaymentMethodMaya,
		UserID:    "user-123",
	})

	require.Error(t, err)
	var appErr *domainerrors.AppError
	require.ErrorAs(t, err, &appErr)
	assert.Equal(t, domainerrors.CodePaymentFailed, appErr.Code)
}

func TestInitiatePayment_DuplicatePendingConflict(t *testing.T) {
	bookingMock := new(mockBookingFinder)
	paymentLister := new(mockPaymentByBookingLister)
	saverMock := new(mockPaymentSaver)
	gwMock := new(mockPaymentGateway)
	eventsMock := new(mockEventPublisher)

	uc := NewInitiatePaymentUseCase(bookingMock, paymentLister, saverMock, gwMock, eventsMock)

	bookingMock.On("FindByID", mock.Anything, "BK-001").Return(completedBooking(), nil)
	paymentLister.On("FindByBooking", mock.Anything, "BK-001").Return([]payment.Payment{
		{PaymentID: "PAY-PENDING", Status: payment.PaymentStatusPending},
	}, nil)

	_, err := uc.Execute(context.Background(), &InitiatePaymentInput{
		BookingID: "BK-001",
		Method:    payment.PaymentMethodGCash,
		UserID:    "user-123",
	})

	require.Error(t, err)
	var appErr *domainerrors.AppError
	require.ErrorAs(t, err, &appErr)
	assert.Equal(t, domainerrors.CodeConflict, appErr.Code)
}
