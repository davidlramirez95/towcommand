package paymentuc

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	domainerrors "github.com/davidlramirez95/towcommand/internal/domain/errors"
	"github.com/davidlramirez95/towcommand/internal/domain/payment"
)

// --- Test helpers ---

func webhookJSON(t *testing.T, paymentID, eventType string) []byte {
	t.Helper()
	data, err := json.Marshal(webhookPayload{PaymentID: paymentID, Event: eventType})
	require.NoError(t, err)
	return data
}

// --- Tests ---

func TestProcessWebhook_CaptureSuccess(t *testing.T) {
	gwMock := new(mockPaymentGateway)
	paymentMock := new(mockPaymentFinder)
	updaterMock := new(mockPaymentStatusUpdater)

	uc := NewProcessWebhookUseCase(gwMock, paymentMock, updaterMock)
	uc.now = func() time.Time { return fixedTime }

	payload := webhookJSON(t, "PAY-001", "payment.captured")

	gwMock.On("VerifyWebhookSignature", payload, "valid-sig").Return(nil)
	paymentMock.On("FindByID", mock.Anything, "PAY-001").Return(pendingPayment(), nil)
	updaterMock.On("UpdateStatus", mock.Anything, "PAY-001", payment.PaymentStatusCaptured).Return(nil)

	result, err := uc.Execute(context.Background(), &ProcessWebhookInput{
		Payload:   payload,
		Signature: "valid-sig",
	})

	require.NoError(t, err)
	assert.Equal(t, payment.PaymentStatusCaptured, result.Status)
	assert.NotNil(t, result.CapturedAt)

	updaterMock.AssertExpectations(t)
}

func TestProcessWebhook_RefundSuccess(t *testing.T) {
	gwMock := new(mockPaymentGateway)
	paymentMock := new(mockPaymentFinder)
	updaterMock := new(mockPaymentStatusUpdater)

	uc := NewProcessWebhookUseCase(gwMock, paymentMock, updaterMock)
	uc.now = func() time.Time { return fixedTime }

	payload := webhookJSON(t, "PAY-001", "payment.refunded")

	gwMock.On("VerifyWebhookSignature", payload, "valid-sig").Return(nil)
	paymentMock.On("FindByID", mock.Anything, "PAY-001").Return(capturedPayment(), nil)
	updaterMock.On("UpdateStatus", mock.Anything, "PAY-001", payment.PaymentStatusRefunded).Return(nil)

	result, err := uc.Execute(context.Background(), &ProcessWebhookInput{
		Payload:   payload,
		Signature: "valid-sig",
	})

	require.NoError(t, err)
	assert.Equal(t, payment.PaymentStatusRefunded, result.Status)
	assert.NotNil(t, result.RefundedAt)
}

func TestProcessWebhook_InvalidSignature(t *testing.T) {
	gwMock := new(mockPaymentGateway)
	paymentMock := new(mockPaymentFinder)
	updaterMock := new(mockPaymentStatusUpdater)

	uc := NewProcessWebhookUseCase(gwMock, paymentMock, updaterMock)

	payload := webhookJSON(t, "PAY-001", "payment.captured")

	gwMock.On("VerifyWebhookSignature", payload, "bad-sig").Return(assert.AnError)

	_, err := uc.Execute(context.Background(), &ProcessWebhookInput{
		Payload:   payload,
		Signature: "bad-sig",
	})

	require.Error(t, err)
	var appErr *domainerrors.AppError
	require.ErrorAs(t, err, &appErr)
	assert.Equal(t, domainerrors.CodeUnauthorized, appErr.Code)
}

func TestProcessWebhook_Idempotent(t *testing.T) {
	gwMock := new(mockPaymentGateway)
	paymentMock := new(mockPaymentFinder)
	updaterMock := new(mockPaymentStatusUpdater)

	uc := NewProcessWebhookUseCase(gwMock, paymentMock, updaterMock)

	payload := webhookJSON(t, "PAY-001", "payment.captured")

	// Payment already captured — should be idempotent.
	alreadyCaptured := capturedPayment()
	alreadyCaptured.Status = payment.PaymentStatusCaptured

	gwMock.On("VerifyWebhookSignature", payload, "valid-sig").Return(nil)
	paymentMock.On("FindByID", mock.Anything, "PAY-001").Return(alreadyCaptured, nil)

	result, err := uc.Execute(context.Background(), &ProcessWebhookInput{
		Payload:   payload,
		Signature: "valid-sig",
	})

	require.NoError(t, err)
	assert.Equal(t, payment.PaymentStatusCaptured, result.Status)

	// UpdateStatus should NOT have been called.
	updaterMock.AssertNotCalled(t, "UpdateStatus")
}

func TestProcessWebhook_UnknownEventType(t *testing.T) {
	gwMock := new(mockPaymentGateway)
	paymentMock := new(mockPaymentFinder)
	updaterMock := new(mockPaymentStatusUpdater)

	uc := NewProcessWebhookUseCase(gwMock, paymentMock, updaterMock)

	payload := webhookJSON(t, "PAY-001", "payment.unknown")

	gwMock.On("VerifyWebhookSignature", payload, "valid-sig").Return(nil)
	paymentMock.On("FindByID", mock.Anything, "PAY-001").Return(pendingPayment(), nil)

	_, err := uc.Execute(context.Background(), &ProcessWebhookInput{
		Payload:   payload,
		Signature: "valid-sig",
	})

	require.Error(t, err)
	var appErr *domainerrors.AppError
	require.ErrorAs(t, err, &appErr)
	assert.Equal(t, domainerrors.CodeValidationError, appErr.Code)
	assert.Contains(t, appErr.Message, "unknown webhook event type")
}

func TestProcessWebhook_PaymentNotFound(t *testing.T) {
	gwMock := new(mockPaymentGateway)
	paymentMock := new(mockPaymentFinder)
	updaterMock := new(mockPaymentStatusUpdater)

	uc := NewProcessWebhookUseCase(gwMock, paymentMock, updaterMock)

	payload := webhookJSON(t, "PAY-MISSING", "payment.captured")

	gwMock.On("VerifyWebhookSignature", payload, "valid-sig").Return(nil)
	paymentMock.On("FindByID", mock.Anything, "PAY-MISSING").Return(nil, nil)

	_, err := uc.Execute(context.Background(), &ProcessWebhookInput{
		Payload:   payload,
		Signature: "valid-sig",
	})

	require.Error(t, err)
	var appErr *domainerrors.AppError
	require.ErrorAs(t, err, &appErr)
	assert.Equal(t, domainerrors.CodeNotFound, appErr.Code)
}

func TestProcessWebhook_FailedEvent(t *testing.T) {
	gwMock := new(mockPaymentGateway)
	paymentMock := new(mockPaymentFinder)
	updaterMock := new(mockPaymentStatusUpdater)

	uc := NewProcessWebhookUseCase(gwMock, paymentMock, updaterMock)
	uc.now = func() time.Time { return fixedTime }

	payload := webhookJSON(t, "PAY-001", "payment.failed")

	gwMock.On("VerifyWebhookSignature", payload, "valid-sig").Return(nil)
	paymentMock.On("FindByID", mock.Anything, "PAY-001").Return(pendingPayment(), nil)
	updaterMock.On("UpdateStatus", mock.Anything, "PAY-001", payment.PaymentStatusFailed).Return(nil)

	result, err := uc.Execute(context.Background(), &ProcessWebhookInput{
		Payload:   payload,
		Signature: "valid-sig",
	})

	require.NoError(t, err)
	assert.Equal(t, payment.PaymentStatusFailed, result.Status)
}

func TestProcessWebhook_InvalidPayload(t *testing.T) {
	gwMock := new(mockPaymentGateway)
	paymentMock := new(mockPaymentFinder)
	updaterMock := new(mockPaymentStatusUpdater)

	uc := NewProcessWebhookUseCase(gwMock, paymentMock, updaterMock)

	payload := []byte(`not-json`)

	gwMock.On("VerifyWebhookSignature", payload, "valid-sig").Return(nil)

	_, err := uc.Execute(context.Background(), &ProcessWebhookInput{
		Payload:   payload,
		Signature: "valid-sig",
	})

	require.Error(t, err)
	var appErr *domainerrors.AppError
	require.ErrorAs(t, err, &appErr)
	assert.Equal(t, domainerrors.CodeValidationError, appErr.Code)
}

func TestProcessWebhook_MissingPaymentID(t *testing.T) {
	gwMock := new(mockPaymentGateway)
	paymentMock := new(mockPaymentFinder)
	updaterMock := new(mockPaymentStatusUpdater)

	uc := NewProcessWebhookUseCase(gwMock, paymentMock, updaterMock)

	payload := []byte(`{"event":"payment.captured"}`)

	gwMock.On("VerifyWebhookSignature", payload, "valid-sig").Return(nil)

	_, err := uc.Execute(context.Background(), &ProcessWebhookInput{
		Payload:   payload,
		Signature: "valid-sig",
	})

	require.Error(t, err)
	var appErr *domainerrors.AppError
	require.ErrorAs(t, err, &appErr)
	assert.Equal(t, domainerrors.CodeValidationError, appErr.Code)
	assert.Contains(t, appErr.Message, "missing paymentId")
}
