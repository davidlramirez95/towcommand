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
	"github.com/davidlramirez95/towcommand/internal/domain/provider"
	"github.com/davidlramirez95/towcommand/internal/domain/user"
)

// --- Mocks for capture ---

type mockPaymentFinder struct{ mock.Mock }

func (m *mockPaymentFinder) FindByID(ctx context.Context, paymentID string) (*payment.Payment, error) {
	args := m.Called(ctx, paymentID)
	if v := args.Get(0); v != nil {
		return v.(*payment.Payment), args.Error(1)
	}
	return nil, args.Error(1)
}

type mockPaymentStatusUpdater struct{ mock.Mock }

func (m *mockPaymentStatusUpdater) UpdateStatus(ctx context.Context, paymentID string, status payment.PaymentStatus) error {
	args := m.Called(ctx, paymentID, status)
	return args.Error(0)
}

type mockProviderFinder struct{ mock.Mock }

func (m *mockProviderFinder) FindByID(ctx context.Context, providerID string) (*provider.Provider, error) {
	args := m.Called(ctx, providerID)
	if v := args.Get(0); v != nil {
		return v.(*provider.Provider), args.Error(1)
	}
	return nil, args.Error(1)
}

// --- Test helpers ---

func pendingPayment() *payment.Payment {
	return &payment.Payment{
		PaymentID:  "PAY-001",
		BookingID:  "BK-001",
		UserID:     "user-123",
		Amount:     200_000,
		Currency:   "PHP",
		Method:     payment.PaymentMethodGCash,
		Status:     payment.PaymentStatusPending,
		GatewayRef: "gw-ref-123",
		CreatedAt:  fixedTime.Add(-1 * time.Hour),
		UpdatedAt:  fixedTime.Add(-1 * time.Hour),
	}
}

func testProvider() *provider.Provider {
	return &provider.Provider{
		ProviderID: "prov-456",
		TrustTier:  user.TrustTierSukiSilver,
	}
}

// --- Tests ---

func TestCapturePayment_Success(t *testing.T) {
	paymentMock := new(mockPaymentFinder)
	bookingMock := new(mockBookingFinder)
	providerMock := new(mockProviderFinder)
	updaterMock := new(mockPaymentStatusUpdater)
	eventsMock := new(mockEventPublisher)

	uc := NewCapturePaymentUseCase(paymentMock, bookingMock, providerMock, updaterMock, eventsMock)
	uc.now = func() time.Time { return fixedTime }

	paymentMock.On("FindByID", mock.Anything, "PAY-001").Return(pendingPayment(), nil)
	bookingMock.On("FindByID", mock.Anything, "BK-001").Return(completedBooking(), nil)
	providerMock.On("FindByID", mock.Anything, "prov-456").Return(testProvider(), nil)
	updaterMock.On("UpdateStatus", mock.Anything, "PAY-001", payment.PaymentStatusCaptured).Return(nil)
	eventsMock.On("Publish", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)

	result, err := uc.Execute(context.Background(), "PAY-001")

	require.NoError(t, err)
	assert.Equal(t, payment.PaymentStatusCaptured, result.Payment.Status)
	assert.NotNil(t, result.Payment.CapturedAt)

	// Suki Silver = 20% commission.
	assert.Equal(t, int64(40_000), result.Commission)
	assert.Equal(t, int64(160_000), result.NetAmount)

	updaterMock.AssertExpectations(t)
	eventsMock.AssertExpectations(t)
}

func TestCapturePayment_NotFound(t *testing.T) {
	paymentMock := new(mockPaymentFinder)
	bookingMock := new(mockBookingFinder)
	providerMock := new(mockProviderFinder)
	updaterMock := new(mockPaymentStatusUpdater)
	eventsMock := new(mockEventPublisher)

	uc := NewCapturePaymentUseCase(paymentMock, bookingMock, providerMock, updaterMock, eventsMock)

	paymentMock.On("FindByID", mock.Anything, "PAY-MISSING").Return(nil, nil)

	_, err := uc.Execute(context.Background(), "PAY-MISSING")

	require.Error(t, err)
	var appErr *domainerrors.AppError
	require.ErrorAs(t, err, &appErr)
	assert.Equal(t, domainerrors.CodeNotFound, appErr.Code)
}

func TestCapturePayment_InvalidStatus(t *testing.T) {
	paymentMock := new(mockPaymentFinder)
	bookingMock := new(mockBookingFinder)
	providerMock := new(mockProviderFinder)
	updaterMock := new(mockPaymentStatusUpdater)
	eventsMock := new(mockEventPublisher)

	uc := NewCapturePaymentUseCase(paymentMock, bookingMock, providerMock, updaterMock, eventsMock)

	captured := pendingPayment()
	captured.Status = payment.PaymentStatusCaptured
	paymentMock.On("FindByID", mock.Anything, "PAY-001").Return(captured, nil)

	_, err := uc.Execute(context.Background(), "PAY-001")

	require.Error(t, err)
	var appErr *domainerrors.AppError
	require.ErrorAs(t, err, &appErr)
	assert.Equal(t, domainerrors.CodeConflict, appErr.Code)
}

func TestCapturePayment_ProviderNotFound(t *testing.T) {
	paymentMock := new(mockPaymentFinder)
	bookingMock := new(mockBookingFinder)
	providerMock := new(mockProviderFinder)
	updaterMock := new(mockPaymentStatusUpdater)
	eventsMock := new(mockEventPublisher)

	uc := NewCapturePaymentUseCase(paymentMock, bookingMock, providerMock, updaterMock, eventsMock)

	paymentMock.On("FindByID", mock.Anything, "PAY-001").Return(pendingPayment(), nil)
	bookingMock.On("FindByID", mock.Anything, "BK-001").Return(completedBooking(), nil)
	providerMock.On("FindByID", mock.Anything, "prov-456").Return(nil, nil)

	_, err := uc.Execute(context.Background(), "PAY-001")

	require.Error(t, err)
	var appErr *domainerrors.AppError
	require.ErrorAs(t, err, &appErr)
	assert.Equal(t, domainerrors.CodeNotFound, appErr.Code)
}

func TestCapturePayment_CommissionByTier(t *testing.T) {
	tests := []struct {
		name           string
		tier           user.TrustTier
		wantCommission int64
		wantNet        int64
	}{
		{"basic 25%", user.TrustTierBasic, 50_000, 150_000},
		{"verified 22%", user.TrustTierVerified, 44_000, 156_000},
		{"suki_silver 20%", user.TrustTierSukiSilver, 40_000, 160_000},
		{"suki_gold 18%", user.TrustTierSukiGold, 36_000, 164_000},
		{"suki_elite 15%", user.TrustTierSukiElite, 30_000, 170_000},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			paymentMock := new(mockPaymentFinder)
			bookingMock := new(mockBookingFinder)
			providerMock := new(mockProviderFinder)
			updaterMock := new(mockPaymentStatusUpdater)
			eventsMock := new(mockEventPublisher)

			uc := NewCapturePaymentUseCase(paymentMock, bookingMock, providerMock, updaterMock, eventsMock)
			uc.now = func() time.Time { return fixedTime }

			paymentMock.On("FindByID", mock.Anything, "PAY-001").Return(pendingPayment(), nil)
			bookingMock.On("FindByID", mock.Anything, "BK-001").Return(completedBooking(), nil)
			providerMock.On("FindByID", mock.Anything, "prov-456").Return(&provider.Provider{
				ProviderID: "prov-456",
				TrustTier:  tt.tier,
			}, nil)
			updaterMock.On("UpdateStatus", mock.Anything, mock.Anything, mock.Anything).Return(nil)
			eventsMock.On("Publish", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)

			result, err := uc.Execute(context.Background(), "PAY-001")

			require.NoError(t, err)
			assert.Equal(t, tt.wantCommission, result.Commission)
			assert.Equal(t, tt.wantNet, result.NetAmount)
		})
	}
}
