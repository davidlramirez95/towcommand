package paymentuc

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	domainerrors "github.com/davidlramirez95/towcommand/internal/domain/errors"
	"github.com/davidlramirez95/towcommand/internal/domain/payment"
	"github.com/davidlramirez95/towcommand/internal/usecase/port"
)

// --- Test helpers ---

func capturedPayment() *payment.Payment {
	capturedAt := fixedTime.Add(-30 * time.Minute)
	return &payment.Payment{
		PaymentID:  "PAY-001",
		BookingID:  "BK-001",
		UserID:     "user-123",
		Amount:     200_000,
		Currency:   "PHP",
		Method:     payment.PaymentMethodGCash,
		Status:     payment.PaymentStatusCaptured,
		GatewayRef: "gw-ref-123",
		CapturedAt: &capturedAt,
		CreatedAt:  fixedTime.Add(-1 * time.Hour),
		UpdatedAt:  fixedTime.Add(-30 * time.Minute),
	}
}

// --- Tests ---

func TestRefundPayment_Success(t *testing.T) {
	paymentMock := new(mockPaymentFinder)
	updaterMock := new(mockPaymentStatusUpdater)
	gwMock := new(mockPaymentGateway)
	eventsMock := new(mockEventPublisher)

	uc := NewRefundPaymentUseCase(paymentMock, updaterMock, gwMock, eventsMock)
	uc.now = func() time.Time { return fixedTime }

	paymentMock.On("FindByID", mock.Anything, "PAY-001").Return(capturedPayment(), nil)
	gwMock.On("Refund", mock.Anything, "gw-ref-123", int64(200_000)).
		Return(&port.RefundResult{GatewayRef: "mock-refund-xyz"}, nil)
	updaterMock.On("UpdateStatus", mock.Anything, "PAY-001", payment.PaymentStatusRefunded).Return(nil)
	eventsMock.On("Publish", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)

	result, err := uc.Execute(context.Background(), &RefundPaymentInput{
		PaymentID: "PAY-001",
		Reason:    "customer request",
	})

	require.NoError(t, err)
	assert.Equal(t, payment.PaymentStatusRefunded, result.Status)
	assert.NotNil(t, result.RefundedAt)
	assert.Equal(t, "customer request", result.RefundReason)

	gwMock.AssertExpectations(t)
	updaterMock.AssertExpectations(t)
	eventsMock.AssertExpectations(t)
}

func TestRefundPayment_NotFound(t *testing.T) {
	paymentMock := new(mockPaymentFinder)
	updaterMock := new(mockPaymentStatusUpdater)
	gwMock := new(mockPaymentGateway)
	eventsMock := new(mockEventPublisher)

	uc := NewRefundPaymentUseCase(paymentMock, updaterMock, gwMock, eventsMock)

	paymentMock.On("FindByID", mock.Anything, "PAY-MISSING").Return(nil, nil)

	_, err := uc.Execute(context.Background(), &RefundPaymentInput{
		PaymentID: "PAY-MISSING",
		Reason:    "test",
	})

	require.Error(t, err)
	var appErr *domainerrors.AppError
	require.ErrorAs(t, err, &appErr)
	assert.Equal(t, domainerrors.CodeNotFound, appErr.Code)
}

func TestRefundPayment_InvalidStatus(t *testing.T) {
	paymentMock := new(mockPaymentFinder)
	updaterMock := new(mockPaymentStatusUpdater)
	gwMock := new(mockPaymentGateway)
	eventsMock := new(mockEventPublisher)

	uc := NewRefundPaymentUseCase(paymentMock, updaterMock, gwMock, eventsMock)

	pending := capturedPayment()
	pending.Status = payment.PaymentStatusPending
	paymentMock.On("FindByID", mock.Anything, "PAY-001").Return(pending, nil)

	_, err := uc.Execute(context.Background(), &RefundPaymentInput{
		PaymentID: "PAY-001",
		Reason:    "test",
	})

	require.Error(t, err)
	var appErr *domainerrors.AppError
	require.ErrorAs(t, err, &appErr)
	assert.Equal(t, domainerrors.CodeConflict, appErr.Code)
}

func TestRefundPayment_GatewayFailure(t *testing.T) {
	paymentMock := new(mockPaymentFinder)
	updaterMock := new(mockPaymentStatusUpdater)
	gwMock := new(mockPaymentGateway)
	eventsMock := new(mockEventPublisher)

	uc := NewRefundPaymentUseCase(paymentMock, updaterMock, gwMock, eventsMock)

	paymentMock.On("FindByID", mock.Anything, "PAY-001").Return(capturedPayment(), nil)
	gwMock.On("Refund", mock.Anything, "gw-ref-123", int64(200_000)).
		Return(nil, assert.AnError)

	_, err := uc.Execute(context.Background(), &RefundPaymentInput{
		PaymentID: "PAY-001",
		Reason:    "test",
	})

	require.Error(t, err)
	var appErr *domainerrors.AppError
	require.ErrorAs(t, err, &appErr)
	assert.Equal(t, domainerrors.CodePaymentFailed, appErr.Code)

	// Status should NOT have been updated.
	updaterMock.AssertNotCalled(t, "UpdateStatus")
}

func TestRefundPayment_AlreadyRefunded(t *testing.T) {
	paymentMock := new(mockPaymentFinder)
	updaterMock := new(mockPaymentStatusUpdater)
	gwMock := new(mockPaymentGateway)
	eventsMock := new(mockEventPublisher)

	uc := NewRefundPaymentUseCase(paymentMock, updaterMock, gwMock, eventsMock)

	refunded := capturedPayment()
	refunded.Status = payment.PaymentStatusRefunded
	paymentMock.On("FindByID", mock.Anything, "PAY-001").Return(refunded, nil)

	_, err := uc.Execute(context.Background(), &RefundPaymentInput{
		PaymentID: "PAY-001",
		Reason:    "double refund attempt",
	})

	require.Error(t, err)
	var appErr *domainerrors.AppError
	require.ErrorAs(t, err, &appErr)
	assert.Equal(t, domainerrors.CodeConflict, appErr.Code)
}
