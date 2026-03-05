package handler_test

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"testing"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/davidlramirez95/towcommand/internal/adapter/gateway"
	"github.com/davidlramirez95/towcommand/internal/adapter/handler"
	"github.com/davidlramirez95/towcommand/internal/domain/booking"
	"github.com/davidlramirez95/towcommand/internal/domain/payment"
	"github.com/davidlramirez95/towcommand/internal/domain/provider"
	"github.com/davidlramirez95/towcommand/internal/domain/user"
	paymentuc "github.com/davidlramirez95/towcommand/internal/usecase/payment"
	"github.com/davidlramirez95/towcommand/internal/usecase/port"
)

// ---------------------------------------------------------------------------
// Mock implementations for payment use case dependencies
// ---------------------------------------------------------------------------

type mockPaymentSaver struct {
	SaveFunc func(ctx context.Context, p *payment.Payment) error
}

func (m *mockPaymentSaver) Save(ctx context.Context, p *payment.Payment) error {
	if m.SaveFunc != nil {
		return m.SaveFunc(ctx, p)
	}
	return nil
}

type mockPaymentFinder struct {
	FindByIDFunc func(ctx context.Context, paymentID string) (*payment.Payment, error)
}

func (m *mockPaymentFinder) FindByID(ctx context.Context, paymentID string) (*payment.Payment, error) {
	if m.FindByIDFunc != nil {
		return m.FindByIDFunc(ctx, paymentID)
	}
	return nil, nil
}

type mockPaymentByBookingLister struct {
	FindByBookingFunc func(ctx context.Context, bookingID string) ([]payment.Payment, error)
}

func (m *mockPaymentByBookingLister) FindByBooking(ctx context.Context, bookingID string) ([]payment.Payment, error) {
	if m.FindByBookingFunc != nil {
		return m.FindByBookingFunc(ctx, bookingID)
	}
	return nil, nil
}

type mockPaymentStatusUpdater struct {
	UpdateStatusFunc func(ctx context.Context, paymentID string, status payment.PaymentStatus) error
}

func (m *mockPaymentStatusUpdater) UpdateStatus(ctx context.Context, paymentID string, status payment.PaymentStatus) error {
	if m.UpdateStatusFunc != nil {
		return m.UpdateStatusFunc(ctx, paymentID, status)
	}
	return nil
}

type mockPaymentBookingFinder struct {
	FindByIDFunc func(ctx context.Context, bookingID string) (*booking.Booking, error)
}

func (m *mockPaymentBookingFinder) FindByID(ctx context.Context, bookingID string) (*booking.Booking, error) {
	if m.FindByIDFunc != nil {
		return m.FindByIDFunc(ctx, bookingID)
	}
	return nil, nil
}

type mockPaymentProviderFinder struct {
	FindByIDFunc func(ctx context.Context, providerID string) (*provider.Provider, error)
}

func (m *mockPaymentProviderFinder) FindByID(ctx context.Context, providerID string) (*provider.Provider, error) {
	if m.FindByIDFunc != nil {
		return m.FindByIDFunc(ctx, providerID)
	}
	return nil, nil
}

type mockPaymentEventPublisher struct {
	PublishFunc func(ctx context.Context, source, detailType string, detail any, actor *port.Actor) error
}

func (m *mockPaymentEventPublisher) Publish(ctx context.Context, source, detailType string, detail any, actor *port.Actor) error {
	if m.PublishFunc != nil {
		return m.PublishFunc(ctx, source, detailType, detail, actor)
	}
	return nil
}

// ---------------------------------------------------------------------------
// Helper: build API Gateway event with Cognito auth (reuse pattern)
// ---------------------------------------------------------------------------

func paymentEventWithAuth(userID string) *events.APIGatewayProxyRequest {
	return &events.APIGatewayProxyRequest{
		RequestContext: events.APIGatewayProxyRequestContext{
			Authorizer: map[string]interface{}{
				"claims": map[string]interface{}{
					"sub": userID,
				},
			},
		},
	}
}

func paymentParseErrorBody(t *testing.T, body string) errorBody {
	t.Helper()
	var eb errorBody
	require.NoError(t, json.Unmarshal([]byte(body), &eb))
	return eb
}

// ---------------------------------------------------------------------------
// InitiatePaymentHandler tests
// ---------------------------------------------------------------------------

func TestInitiatePaymentHandler(t *testing.T) {
	completedBooking := &booking.Booking{
		BookingID:  "bk-1",
		CustomerID: "user-1",
		ProviderID: "prov-1",
		Status:     booking.BookingStatusCompleted,
		Price: booking.PriceBreakdown{
			Total:    250000,
			Currency: "PHP",
		},
	}

	tests := []struct {
		name        string
		event       *events.APIGatewayProxyRequest
		setupMocks  func(bf *mockPaymentBookingFinder, pbl *mockPaymentByBookingLister, ps *mockPaymentSaver, ep *mockPaymentEventPublisher)
		wantStatus  int
		wantErrCode string
		checkBody   func(t *testing.T, body string)
	}{
		{
			name: "success - cash payment",
			event: func() *events.APIGatewayProxyRequest {
				e := paymentEventWithAuth("user-1")
				e.PathParameters = map[string]string{"id": "bk-1"}
				e.Body = `{"method":"cash"}`
				return e
			}(),
			setupMocks: func(bf *mockPaymentBookingFinder, pbl *mockPaymentByBookingLister, ps *mockPaymentSaver, ep *mockPaymentEventPublisher) {
				bf.FindByIDFunc = func(_ context.Context, _ string) (*booking.Booking, error) {
					return completedBooking, nil
				}
				pbl.FindByBookingFunc = func(_ context.Context, _ string) ([]payment.Payment, error) {
					return nil, nil
				}
			},
			wantStatus: http.StatusCreated,
			checkBody: func(t *testing.T, body string) {
				t.Helper()
				var p payment.Payment
				require.NoError(t, json.Unmarshal([]byte(body), &p))
				assert.Equal(t, "bk-1", p.BookingID)
				assert.Equal(t, "user-1", p.UserID)
				assert.Equal(t, payment.PaymentMethodCash, p.Method)
				assert.Equal(t, payment.PaymentStatusCaptured, p.Status)
				assert.Equal(t, int64(250000), p.Amount)
				assert.Equal(t, "PHP", p.Currency)
			},
		},
		{
			name: "success - gcash payment",
			event: func() *events.APIGatewayProxyRequest {
				e := paymentEventWithAuth("user-1")
				e.PathParameters = map[string]string{"id": "bk-1"}
				e.Body = `{"method":"gcash"}`
				return e
			}(),
			setupMocks: func(bf *mockPaymentBookingFinder, pbl *mockPaymentByBookingLister, ps *mockPaymentSaver, ep *mockPaymentEventPublisher) {
				bf.FindByIDFunc = func(_ context.Context, _ string) (*booking.Booking, error) {
					return completedBooking, nil
				}
				pbl.FindByBookingFunc = func(_ context.Context, _ string) ([]payment.Payment, error) {
					return nil, nil
				}
			},
			wantStatus: http.StatusCreated,
			checkBody: func(t *testing.T, body string) {
				t.Helper()
				var p payment.Payment
				require.NoError(t, json.Unmarshal([]byte(body), &p))
				assert.Equal(t, payment.PaymentMethodGCash, p.Method)
				assert.Equal(t, payment.PaymentStatusPending, p.Status)
				assert.Contains(t, p.GatewayRef, "mock-")
			},
		},
		{
			name: "unauthorized - no user ID",
			event: &events.APIGatewayProxyRequest{
				PathParameters: map[string]string{"id": "bk-1"},
				Body:           `{"method":"cash"}`,
			},
			wantStatus:  http.StatusUnauthorized,
			wantErrCode: "UNAUTHORIZED",
		},
		{
			name: "missing booking ID",
			event: func() *events.APIGatewayProxyRequest {
				e := paymentEventWithAuth("user-1")
				e.Body = `{"method":"cash"}`
				return e
			}(),
			wantStatus:  http.StatusBadRequest,
			wantErrCode: "VALIDATION_ERROR",
		},
		{
			name: "invalid body - bad JSON",
			event: func() *events.APIGatewayProxyRequest {
				e := paymentEventWithAuth("user-1")
				e.PathParameters = map[string]string{"id": "bk-1"}
				e.Body = `{not json}`
				return e
			}(),
			wantStatus:  http.StatusBadRequest,
			wantErrCode: "VALIDATION_ERROR",
		},
		{
			name: "invalid body - invalid method",
			event: func() *events.APIGatewayProxyRequest {
				e := paymentEventWithAuth("user-1")
				e.PathParameters = map[string]string{"id": "bk-1"}
				e.Body = `{"method":"bitcoin"}`
				return e
			}(),
			wantStatus:  http.StatusBadRequest,
			wantErrCode: "VALIDATION_ERROR",
		},
		{
			name: "use case error - booking not found",
			event: func() *events.APIGatewayProxyRequest {
				e := paymentEventWithAuth("user-1")
				e.PathParameters = map[string]string{"id": "bk-missing"}
				e.Body = `{"method":"cash"}`
				return e
			}(),
			setupMocks: func(bf *mockPaymentBookingFinder, pbl *mockPaymentByBookingLister, ps *mockPaymentSaver, ep *mockPaymentEventPublisher) {
				bf.FindByIDFunc = func(_ context.Context, _ string) (*booking.Booking, error) {
					return nil, nil
				}
			},
			wantStatus:  http.StatusNotFound,
			wantErrCode: "NOT_FOUND",
		},
		{
			name: "use case error - booking not completed",
			event: func() *events.APIGatewayProxyRequest {
				e := paymentEventWithAuth("user-1")
				e.PathParameters = map[string]string{"id": "bk-1"}
				e.Body = `{"method":"cash"}`
				return e
			}(),
			setupMocks: func(bf *mockPaymentBookingFinder, pbl *mockPaymentByBookingLister, ps *mockPaymentSaver, ep *mockPaymentEventPublisher) {
				bf.FindByIDFunc = func(_ context.Context, _ string) (*booking.Booking, error) {
					return &booking.Booking{
						BookingID:  "bk-1",
						CustomerID: "user-1",
						Status:     booking.BookingStatusPending,
					}, nil
				}
			},
			wantStatus:  http.StatusBadRequest,
			wantErrCode: "VALIDATION_ERROR",
		},
		{
			name: "use case error - duplicate payment",
			event: func() *events.APIGatewayProxyRequest {
				e := paymentEventWithAuth("user-1")
				e.PathParameters = map[string]string{"id": "bk-1"}
				e.Body = `{"method":"cash"}`
				return e
			}(),
			setupMocks: func(bf *mockPaymentBookingFinder, pbl *mockPaymentByBookingLister, ps *mockPaymentSaver, ep *mockPaymentEventPublisher) {
				bf.FindByIDFunc = func(_ context.Context, _ string) (*booking.Booking, error) {
					return completedBooking, nil
				}
				pbl.FindByBookingFunc = func(_ context.Context, _ string) ([]payment.Payment, error) {
					return []payment.Payment{
						{PaymentID: "pay-1", Status: payment.PaymentStatusCaptured},
					}, nil
				}
			},
			wantStatus:  http.StatusConflict,
			wantErrCode: "CONFLICT",
		},
		{
			name: "use case error - internal error from repo",
			event: func() *events.APIGatewayProxyRequest {
				e := paymentEventWithAuth("user-1")
				e.PathParameters = map[string]string{"id": "bk-1"}
				e.Body = `{"method":"cash"}`
				return e
			}(),
			setupMocks: func(bf *mockPaymentBookingFinder, pbl *mockPaymentByBookingLister, ps *mockPaymentSaver, ep *mockPaymentEventPublisher) {
				bf.FindByIDFunc = func(_ context.Context, _ string) (*booking.Booking, error) {
					return nil, errors.New("dynamo timeout")
				}
			},
			wantStatus:  http.StatusInternalServerError,
			wantErrCode: "INTERNAL_ERROR",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bf := &mockPaymentBookingFinder{}
			pbl := &mockPaymentByBookingLister{}
			ps := &mockPaymentSaver{}
			ep := &mockPaymentEventPublisher{}

			if tt.setupMocks != nil {
				tt.setupMocks(bf, pbl, ps, ep)
			}

			gw := gateway.NewMockPaymentGateway("test-secret")
			uc := paymentuc.NewInitiatePaymentUseCase(bf, pbl, ps, gw, ep)
			h := handler.NewInitiatePaymentHandler(uc)

			resp, err := h.Handle(context.Background(), tt.event)

			require.NoError(t, err)
			assert.Equal(t, tt.wantStatus, resp.StatusCode)

			if tt.wantErrCode != "" {
				eb := paymentParseErrorBody(t, resp.Body)
				assert.Equal(t, tt.wantErrCode, eb.Error.Code)
			}

			if tt.checkBody != nil {
				tt.checkBody(t, resp.Body)
			}
		})
	}
}

// ---------------------------------------------------------------------------
// CapturePaymentHandler tests
// ---------------------------------------------------------------------------

func TestCapturePaymentHandler(t *testing.T) {
	now := time.Now().UTC()

	pendingPayment := &payment.Payment{
		PaymentID:  "pay-1",
		BookingID:  "bk-1",
		UserID:     "user-1",
		Amount:     250000,
		Currency:   "PHP",
		Method:     payment.PaymentMethodGCash,
		Status:     payment.PaymentStatusPending,
		GatewayRef: "mock-ref-1",
		CreatedAt:  now,
		UpdatedAt:  now,
	}

	testBooking := &booking.Booking{
		BookingID:  "bk-1",
		CustomerID: "user-1",
		ProviderID: "prov-1",
		Status:     booking.BookingStatusCompleted,
	}

	testProvider := &provider.Provider{
		ProviderID: "prov-1",
		TrustTier:  user.TrustTierBasic,
	}

	tests := []struct {
		name        string
		event       *events.APIGatewayProxyRequest
		setupMocks  func(pf *mockPaymentFinder, bf *mockPaymentBookingFinder, provf *mockPaymentProviderFinder, pu *mockPaymentStatusUpdater, ep *mockPaymentEventPublisher)
		wantStatus  int
		wantErrCode string
		checkBody   func(t *testing.T, body string)
	}{
		{
			name: "success",
			event: func() *events.APIGatewayProxyRequest {
				e := paymentEventWithAuth("admin-1")
				e.PathParameters = map[string]string{"id": "pay-1"}
				return e
			}(),
			setupMocks: func(pf *mockPaymentFinder, bf *mockPaymentBookingFinder, provf *mockPaymentProviderFinder, pu *mockPaymentStatusUpdater, ep *mockPaymentEventPublisher) {
				pf.FindByIDFunc = func(_ context.Context, _ string) (*payment.Payment, error) {
					// Return a fresh copy to avoid mutation issues across tests.
					p := *pendingPayment
					return &p, nil
				}
				bf.FindByIDFunc = func(_ context.Context, _ string) (*booking.Booking, error) {
					return testBooking, nil
				}
				provf.FindByIDFunc = func(_ context.Context, _ string) (*provider.Provider, error) {
					return testProvider, nil
				}
			},
			wantStatus: http.StatusOK,
			checkBody: func(t *testing.T, body string) {
				t.Helper()
				var result paymentuc.CaptureResult
				require.NoError(t, json.Unmarshal([]byte(body), &result))
				assert.Equal(t, payment.PaymentStatusCaptured, result.Payment.Status)
				assert.Greater(t, result.Commission, int64(0))
				assert.Greater(t, result.NetAmount, int64(0))
			},
		},
		{
			name: "unauthorized - no user ID",
			event: &events.APIGatewayProxyRequest{
				PathParameters: map[string]string{"id": "pay-1"},
			},
			wantStatus:  http.StatusUnauthorized,
			wantErrCode: "UNAUTHORIZED",
		},
		{
			name: "missing payment ID",
			event: func() *events.APIGatewayProxyRequest {
				return paymentEventWithAuth("admin-1")
			}(),
			wantStatus:  http.StatusBadRequest,
			wantErrCode: "VALIDATION_ERROR",
		},
		{
			name: "use case error - payment not found",
			event: func() *events.APIGatewayProxyRequest {
				e := paymentEventWithAuth("admin-1")
				e.PathParameters = map[string]string{"id": "pay-missing"}
				return e
			}(),
			setupMocks: func(pf *mockPaymentFinder, bf *mockPaymentBookingFinder, provf *mockPaymentProviderFinder, pu *mockPaymentStatusUpdater, ep *mockPaymentEventPublisher) {
				pf.FindByIDFunc = func(_ context.Context, _ string) (*payment.Payment, error) {
					return nil, nil
				}
			},
			wantStatus:  http.StatusNotFound,
			wantErrCode: "NOT_FOUND",
		},
		{
			name: "use case error - payment not pending",
			event: func() *events.APIGatewayProxyRequest {
				e := paymentEventWithAuth("admin-1")
				e.PathParameters = map[string]string{"id": "pay-1"}
				return e
			}(),
			setupMocks: func(pf *mockPaymentFinder, bf *mockPaymentBookingFinder, provf *mockPaymentProviderFinder, pu *mockPaymentStatusUpdater, ep *mockPaymentEventPublisher) {
				pf.FindByIDFunc = func(_ context.Context, _ string) (*payment.Payment, error) {
					return &payment.Payment{
						PaymentID: "pay-1",
						Status:    payment.PaymentStatusCaptured,
					}, nil
				}
			},
			wantStatus:  http.StatusConflict,
			wantErrCode: "CONFLICT",
		},
		{
			name: "use case error - internal error from repo",
			event: func() *events.APIGatewayProxyRequest {
				e := paymentEventWithAuth("admin-1")
				e.PathParameters = map[string]string{"id": "pay-1"}
				return e
			}(),
			setupMocks: func(pf *mockPaymentFinder, bf *mockPaymentBookingFinder, provf *mockPaymentProviderFinder, pu *mockPaymentStatusUpdater, ep *mockPaymentEventPublisher) {
				pf.FindByIDFunc = func(_ context.Context, _ string) (*payment.Payment, error) {
					return nil, errors.New("dynamo timeout")
				}
			},
			wantStatus:  http.StatusInternalServerError,
			wantErrCode: "INTERNAL_ERROR",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pf := &mockPaymentFinder{}
			bf := &mockPaymentBookingFinder{}
			provf := &mockPaymentProviderFinder{}
			pu := &mockPaymentStatusUpdater{}
			ep := &mockPaymentEventPublisher{}

			if tt.setupMocks != nil {
				tt.setupMocks(pf, bf, provf, pu, ep)
			}

			uc := paymentuc.NewCapturePaymentUseCase(pf, bf, provf, pu, ep)
			h := handler.NewCapturePaymentHandler(uc)

			resp, err := h.Handle(context.Background(), tt.event)

			require.NoError(t, err)
			assert.Equal(t, tt.wantStatus, resp.StatusCode)

			if tt.wantErrCode != "" {
				eb := paymentParseErrorBody(t, resp.Body)
				assert.Equal(t, tt.wantErrCode, eb.Error.Code)
			}

			if tt.checkBody != nil {
				tt.checkBody(t, resp.Body)
			}
		})
	}
}

// ---------------------------------------------------------------------------
// RefundPaymentHandler tests
// ---------------------------------------------------------------------------

func TestRefundPaymentHandler(t *testing.T) {
	now := time.Now().UTC()

	capturedPayment := &payment.Payment{
		PaymentID:  "pay-1",
		BookingID:  "bk-1",
		UserID:     "user-1",
		Amount:     250000,
		Currency:   "PHP",
		Method:     payment.PaymentMethodGCash,
		Status:     payment.PaymentStatusCaptured,
		GatewayRef: "mock-ref-1",
		CapturedAt: &now,
		CreatedAt:  now,
		UpdatedAt:  now,
	}

	tests := []struct {
		name        string
		event       *events.APIGatewayProxyRequest
		setupMocks  func(pf *mockPaymentFinder, pu *mockPaymentStatusUpdater, ep *mockPaymentEventPublisher)
		wantStatus  int
		wantErrCode string
		checkBody   func(t *testing.T, body string)
	}{
		{
			name: "success",
			event: func() *events.APIGatewayProxyRequest {
				e := paymentEventWithAuth("admin-1")
				e.PathParameters = map[string]string{"id": "pay-1"}
				e.Body = `{"reason":"Customer requested refund"}`
				return e
			}(),
			setupMocks: func(pf *mockPaymentFinder, pu *mockPaymentStatusUpdater, ep *mockPaymentEventPublisher) {
				pf.FindByIDFunc = func(_ context.Context, _ string) (*payment.Payment, error) {
					p := *capturedPayment
					return &p, nil
				}
			},
			wantStatus: http.StatusOK,
			checkBody: func(t *testing.T, body string) {
				t.Helper()
				var p payment.Payment
				require.NoError(t, json.Unmarshal([]byte(body), &p))
				assert.Equal(t, payment.PaymentStatusRefunded, p.Status)
				assert.Equal(t, "Customer requested refund", p.RefundReason)
			},
		},
		{
			name: "unauthorized - no user ID",
			event: &events.APIGatewayProxyRequest{
				PathParameters: map[string]string{"id": "pay-1"},
				Body:           `{"reason":"test"}`,
			},
			wantStatus:  http.StatusUnauthorized,
			wantErrCode: "UNAUTHORIZED",
		},
		{
			name: "missing payment ID",
			event: func() *events.APIGatewayProxyRequest {
				e := paymentEventWithAuth("admin-1")
				e.Body = `{"reason":"test"}`
				return e
			}(),
			wantStatus:  http.StatusBadRequest,
			wantErrCode: "VALIDATION_ERROR",
		},
		{
			name: "invalid body - bad JSON",
			event: func() *events.APIGatewayProxyRequest {
				e := paymentEventWithAuth("admin-1")
				e.PathParameters = map[string]string{"id": "pay-1"}
				e.Body = `{not json}`
				return e
			}(),
			wantStatus:  http.StatusBadRequest,
			wantErrCode: "VALIDATION_ERROR",
		},
		{
			name: "invalid body - missing reason",
			event: func() *events.APIGatewayProxyRequest {
				e := paymentEventWithAuth("admin-1")
				e.PathParameters = map[string]string{"id": "pay-1"}
				e.Body = `{}`
				return e
			}(),
			wantStatus:  http.StatusBadRequest,
			wantErrCode: "VALIDATION_ERROR",
		},
		{
			name: "use case error - payment not found",
			event: func() *events.APIGatewayProxyRequest {
				e := paymentEventWithAuth("admin-1")
				e.PathParameters = map[string]string{"id": "pay-missing"}
				e.Body = `{"reason":"test"}`
				return e
			}(),
			setupMocks: func(pf *mockPaymentFinder, pu *mockPaymentStatusUpdater, ep *mockPaymentEventPublisher) {
				pf.FindByIDFunc = func(_ context.Context, _ string) (*payment.Payment, error) {
					return nil, nil
				}
			},
			wantStatus:  http.StatusNotFound,
			wantErrCode: "NOT_FOUND",
		},
		{
			name: "use case error - payment not captured",
			event: func() *events.APIGatewayProxyRequest {
				e := paymentEventWithAuth("admin-1")
				e.PathParameters = map[string]string{"id": "pay-1"}
				e.Body = `{"reason":"test"}`
				return e
			}(),
			setupMocks: func(pf *mockPaymentFinder, pu *mockPaymentStatusUpdater, ep *mockPaymentEventPublisher) {
				pf.FindByIDFunc = func(_ context.Context, _ string) (*payment.Payment, error) {
					return &payment.Payment{
						PaymentID: "pay-1",
						Status:    payment.PaymentStatusPending,
					}, nil
				}
			},
			wantStatus:  http.StatusConflict,
			wantErrCode: "CONFLICT",
		},
		{
			name: "use case error - internal error from repo",
			event: func() *events.APIGatewayProxyRequest {
				e := paymentEventWithAuth("admin-1")
				e.PathParameters = map[string]string{"id": "pay-1"}
				e.Body = `{"reason":"test"}`
				return e
			}(),
			setupMocks: func(pf *mockPaymentFinder, pu *mockPaymentStatusUpdater, ep *mockPaymentEventPublisher) {
				pf.FindByIDFunc = func(_ context.Context, _ string) (*payment.Payment, error) {
					return nil, errors.New("dynamo timeout")
				}
			},
			wantStatus:  http.StatusInternalServerError,
			wantErrCode: "INTERNAL_ERROR",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pf := &mockPaymentFinder{}
			pu := &mockPaymentStatusUpdater{}
			ep := &mockPaymentEventPublisher{}

			if tt.setupMocks != nil {
				tt.setupMocks(pf, pu, ep)
			}

			gw := gateway.NewMockPaymentGateway("test-secret")
			uc := paymentuc.NewRefundPaymentUseCase(pf, pu, gw, ep)
			h := handler.NewRefundPaymentHandler(uc)

			resp, err := h.Handle(context.Background(), tt.event)

			require.NoError(t, err)
			assert.Equal(t, tt.wantStatus, resp.StatusCode)

			if tt.wantErrCode != "" {
				eb := paymentParseErrorBody(t, resp.Body)
				assert.Equal(t, tt.wantErrCode, eb.Error.Code)
			}

			if tt.checkBody != nil {
				tt.checkBody(t, resp.Body)
			}
		})
	}
}

// ---------------------------------------------------------------------------
// PaymentWebhookHandler tests
// ---------------------------------------------------------------------------

func TestPaymentWebhookHandler(t *testing.T) {
	webhookSecret := "test-webhook-secret"
	mockGW := gateway.NewMockPaymentGateway(webhookSecret)

	now := time.Now().UTC()

	pendingPayment := &payment.Payment{
		PaymentID:  "pay-1",
		BookingID:  "bk-1",
		UserID:     "user-1",
		Amount:     250000,
		Currency:   "PHP",
		Method:     payment.PaymentMethodGCash,
		Status:     payment.PaymentStatusPending,
		GatewayRef: "mock-ref-1",
		CreatedAt:  now,
		UpdatedAt:  now,
	}

	tests := []struct {
		name        string
		event       *events.APIGatewayProxyRequest
		setupMocks  func(pf *mockPaymentFinder, pu *mockPaymentStatusUpdater)
		wantStatus  int
		wantErrCode string
		checkBody   func(t *testing.T, body string)
	}{
		{
			name: "success - payment captured",
			event: func() *events.APIGatewayProxyRequest {
				payload := `{"paymentId":"pay-1","event":"payment.captured"}`
				return &events.APIGatewayProxyRequest{
					Headers: map[string]string{
						"x-webhook-signature": mockGW.SignPayload([]byte(payload)),
					},
					Body: payload,
				}
			}(),
			setupMocks: func(pf *mockPaymentFinder, pu *mockPaymentStatusUpdater) {
				pf.FindByIDFunc = func(_ context.Context, _ string) (*payment.Payment, error) {
					p := *pendingPayment
					return &p, nil
				}
			},
			wantStatus: http.StatusOK,
			checkBody: func(t *testing.T, body string) {
				t.Helper()
				var p payment.Payment
				require.NoError(t, json.Unmarshal([]byte(body), &p))
				assert.Equal(t, payment.PaymentStatusCaptured, p.Status)
				assert.NotNil(t, p.CapturedAt)
			},
		},
		{
			name: "success - payment refunded",
			event: func() *events.APIGatewayProxyRequest {
				payload := `{"paymentId":"pay-1","event":"payment.refunded"}`
				return &events.APIGatewayProxyRequest{
					Headers: map[string]string{
						"x-webhook-signature": mockGW.SignPayload([]byte(payload)),
					},
					Body: payload,
				}
			}(),
			setupMocks: func(pf *mockPaymentFinder, pu *mockPaymentStatusUpdater) {
				pf.FindByIDFunc = func(_ context.Context, _ string) (*payment.Payment, error) {
					p := *pendingPayment
					return &p, nil
				}
			},
			wantStatus: http.StatusOK,
			checkBody: func(t *testing.T, body string) {
				t.Helper()
				var p payment.Payment
				require.NoError(t, json.Unmarshal([]byte(body), &p))
				assert.Equal(t, payment.PaymentStatusRefunded, p.Status)
			},
		},
		{
			name: "missing signature header",
			event: &events.APIGatewayProxyRequest{
				Headers: map[string]string{},
				Body:    `{"paymentId":"pay-1","event":"payment.captured"}`,
			},
			wantStatus:  http.StatusUnauthorized,
			wantErrCode: "UNAUTHORIZED",
		},
		{
			name: "invalid signature",
			event: &events.APIGatewayProxyRequest{
				Headers: map[string]string{
					"x-webhook-signature": "invalid-signature",
				},
				Body: `{"paymentId":"pay-1","event":"payment.captured"}`,
			},
			wantStatus:  http.StatusUnauthorized,
			wantErrCode: "UNAUTHORIZED",
		},
		{
			name: "invalid payload JSON",
			event: func() *events.APIGatewayProxyRequest {
				payload := `{not valid json}`
				return &events.APIGatewayProxyRequest{
					Headers: map[string]string{
						"x-webhook-signature": mockGW.SignPayload([]byte(payload)),
					},
					Body: payload,
				}
			}(),
			wantStatus:  http.StatusBadRequest,
			wantErrCode: "VALIDATION_ERROR",
		},
		{
			name: "payment not found",
			event: func() *events.APIGatewayProxyRequest {
				payload := `{"paymentId":"pay-missing","event":"payment.captured"}`
				return &events.APIGatewayProxyRequest{
					Headers: map[string]string{
						"x-webhook-signature": mockGW.SignPayload([]byte(payload)),
					},
					Body: payload,
				}
			}(),
			setupMocks: func(pf *mockPaymentFinder, pu *mockPaymentStatusUpdater) {
				pf.FindByIDFunc = func(_ context.Context, _ string) (*payment.Payment, error) {
					return nil, nil
				}
			},
			wantStatus:  http.StatusNotFound,
			wantErrCode: "NOT_FOUND",
		},
		{
			name: "idempotent - already in target status",
			event: func() *events.APIGatewayProxyRequest {
				payload := `{"paymentId":"pay-1","event":"payment.captured"}`
				return &events.APIGatewayProxyRequest{
					Headers: map[string]string{
						"x-webhook-signature": mockGW.SignPayload([]byte(payload)),
					},
					Body: payload,
				}
			}(),
			setupMocks: func(pf *mockPaymentFinder, pu *mockPaymentStatusUpdater) {
				pf.FindByIDFunc = func(_ context.Context, _ string) (*payment.Payment, error) {
					return &payment.Payment{
						PaymentID: "pay-1",
						Status:    payment.PaymentStatusCaptured,
					}, nil
				}
			},
			wantStatus: http.StatusOK,
			checkBody: func(t *testing.T, body string) {
				t.Helper()
				var p payment.Payment
				require.NoError(t, json.Unmarshal([]byte(body), &p))
				assert.Equal(t, payment.PaymentStatusCaptured, p.Status)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pf := &mockPaymentFinder{}
			pu := &mockPaymentStatusUpdater{}

			if tt.setupMocks != nil {
				tt.setupMocks(pf, pu)
			}

			uc := paymentuc.NewProcessWebhookUseCase(mockGW, pf, pu)
			h := handler.NewPaymentWebhookHandler(uc)

			resp, err := h.Handle(context.Background(), tt.event)

			require.NoError(t, err)
			assert.Equal(t, tt.wantStatus, resp.StatusCode)

			if tt.wantErrCode != "" {
				eb := paymentParseErrorBody(t, resp.Body)
				assert.Equal(t, tt.wantErrCode, eb.Error.Code)
			}

			if tt.checkBody != nil {
				tt.checkBody(t, resp.Body)
			}
		})
	}
}

// ---------------------------------------------------------------------------
// CancelFeeHandler tests
// ---------------------------------------------------------------------------

func TestCancelFeeHandler(t *testing.T) {
	completedBooking := &booking.Booking{
		BookingID:  "bk-1",
		CustomerID: "user-1",
		ProviderID: "prov-1",
		Status:     booking.BookingStatusCompleted,
		Price: booking.PriceBreakdown{
			Total:    100000,
			Currency: "PHP",
		},
	}

	tests := []struct {
		name        string
		event       *events.APIGatewayProxyRequest
		setupMocks  func(bf *mockPaymentBookingFinder, pbl *mockPaymentByBookingLister, ps *mockPaymentSaver, ep *mockPaymentEventPublisher)
		wantStatus  int
		wantErrCode string
		checkBody   func(t *testing.T, body string)
	}{
		{
			name: "success",
			event: func() *events.APIGatewayProxyRequest {
				e := paymentEventWithAuth("user-1")
				e.PathParameters = map[string]string{"id": "bk-1"}
				return e
			}(),
			setupMocks: func(bf *mockPaymentBookingFinder, pbl *mockPaymentByBookingLister, ps *mockPaymentSaver, ep *mockPaymentEventPublisher) {
				bf.FindByIDFunc = func(_ context.Context, _ string) (*booking.Booking, error) {
					return completedBooking, nil
				}
				pbl.FindByBookingFunc = func(_ context.Context, _ string) ([]payment.Payment, error) {
					return nil, nil
				}
			},
			wantStatus: http.StatusCreated,
			checkBody: func(t *testing.T, body string) {
				t.Helper()
				var p payment.Payment
				require.NoError(t, json.Unmarshal([]byte(body), &p))
				assert.Equal(t, payment.PaymentMethodCash, p.Method)
				assert.Equal(t, payment.PaymentStatusCaptured, p.Status)
			},
		},
		{
			name: "unauthorized - no user ID",
			event: &events.APIGatewayProxyRequest{
				PathParameters: map[string]string{"id": "bk-1"},
			},
			wantStatus:  http.StatusUnauthorized,
			wantErrCode: "UNAUTHORIZED",
		},
		{
			name: "missing booking ID",
			event: func() *events.APIGatewayProxyRequest {
				return paymentEventWithAuth("user-1")
			}(),
			wantStatus:  http.StatusBadRequest,
			wantErrCode: "VALIDATION_ERROR",
		},
		{
			name: "use case error - booking not found",
			event: func() *events.APIGatewayProxyRequest {
				e := paymentEventWithAuth("user-1")
				e.PathParameters = map[string]string{"id": "bk-missing"}
				return e
			}(),
			setupMocks: func(bf *mockPaymentBookingFinder, pbl *mockPaymentByBookingLister, ps *mockPaymentSaver, ep *mockPaymentEventPublisher) {
				bf.FindByIDFunc = func(_ context.Context, _ string) (*booking.Booking, error) {
					return nil, nil
				}
			},
			wantStatus:  http.StatusNotFound,
			wantErrCode: "NOT_FOUND",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bf := &mockPaymentBookingFinder{}
			pbl := &mockPaymentByBookingLister{}
			ps := &mockPaymentSaver{}
			ep := &mockPaymentEventPublisher{}

			if tt.setupMocks != nil {
				tt.setupMocks(bf, pbl, ps, ep)
			}

			gw := gateway.NewMockPaymentGateway("test-secret")
			uc := paymentuc.NewInitiatePaymentUseCase(bf, pbl, ps, gw, ep)
			h := handler.NewCancelFeeHandler(uc)

			resp, err := h.Handle(context.Background(), tt.event)

			require.NoError(t, err)
			assert.Equal(t, tt.wantStatus, resp.StatusCode)

			if tt.wantErrCode != "" {
				eb := paymentParseErrorBody(t, resp.Body)
				assert.Equal(t, tt.wantErrCode, eb.Error.Code)
			}

			if tt.checkBody != nil {
				tt.checkBody(t, resp.Body)
			}
		})
	}
}
