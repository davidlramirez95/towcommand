package handler_test

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"testing"

	"github.com/aws/aws-lambda-go/events"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/davidlramirez95/towcommand/internal/adapter/handler"
	"github.com/davidlramirez95/towcommand/internal/domain/booking"
	"github.com/davidlramirez95/towcommand/internal/domain/payment"
	"github.com/davidlramirez95/towcommand/internal/domain/provider"
	"github.com/davidlramirez95/towcommand/internal/domain/user"
	paymentuc "github.com/davidlramirez95/towcommand/internal/usecase/payment"
)

// ---------------------------------------------------------------------------
// Mocks for provider earnings handler tests
// ---------------------------------------------------------------------------

type mockEarningsHandlerProviderFinder struct {
	FindByIDFunc func(ctx context.Context, providerID string) (*provider.Provider, error)
}

func (m *mockEarningsHandlerProviderFinder) FindByID(ctx context.Context, providerID string) (*provider.Provider, error) {
	if m.FindByIDFunc != nil {
		return m.FindByIDFunc(ctx, providerID)
	}
	return nil, nil
}

type mockEarningsHandlerBookingLister struct {
	FindByProviderFunc func(ctx context.Context, providerID string) ([]booking.Booking, error)
}

func (m *mockEarningsHandlerBookingLister) FindByProvider(ctx context.Context, providerID string) ([]booking.Booking, error) {
	if m.FindByProviderFunc != nil {
		return m.FindByProviderFunc(ctx, providerID)
	}
	return nil, nil
}

type mockEarningsHandlerPaymentLister struct {
	FindByProviderBookingsFunc func(ctx context.Context, bookingIDs []string) ([]payment.Payment, error)
}

func (m *mockEarningsHandlerPaymentLister) FindByProviderBookings(ctx context.Context, bookingIDs []string) ([]payment.Payment, error) {
	if m.FindByProviderBookingsFunc != nil {
		return m.FindByProviderBookingsFunc(ctx, bookingIDs)
	}
	return nil, nil
}

// earningsEventWithAuth builds an API Gateway event with Cognito auth claims and optional user type.
func earningsEventWithAuth(userID, userType string) *events.APIGatewayProxyRequest {
	claims := map[string]interface{}{
		"sub": userID,
	}
	if userType != "" {
		claims["custom:userType"] = userType
	}
	return &events.APIGatewayProxyRequest{
		RequestContext: events.APIGatewayProxyRequestContext{
			Authorizer: map[string]interface{}{
				"claims": claims,
			},
		},
	}
}

// ---------------------------------------------------------------------------
// Tests
// ---------------------------------------------------------------------------

func TestProviderEarningsHandler(t *testing.T) {
	testProvider := &provider.Provider{
		ProviderID: "prov-1",
		TrustTier:  user.TrustTierBasic,
	}

	tests := []struct {
		name        string
		event       *events.APIGatewayProxyRequest
		setupMocks  func(pf *mockEarningsHandlerProviderFinder, bl *mockEarningsHandlerBookingLister, pl *mockEarningsHandlerPaymentLister)
		wantStatus  int
		wantErrCode string
		checkBody   func(t *testing.T, body string)
	}{
		{
			name: "success - own earnings",
			event: func() *events.APIGatewayProxyRequest {
				e := earningsEventWithAuth("prov-1", "provider")
				e.PathParameters = map[string]string{"id": "prov-1"}
				return e
			}(),
			setupMocks: func(pf *mockEarningsHandlerProviderFinder, bl *mockEarningsHandlerBookingLister, pl *mockEarningsHandlerPaymentLister) {
				pf.FindByIDFunc = func(_ context.Context, _ string) (*provider.Provider, error) {
					return testProvider, nil
				}
				bl.FindByProviderFunc = func(_ context.Context, _ string) ([]booking.Booking, error) {
					return nil, nil
				}
			},
			wantStatus: http.StatusOK,
			checkBody: func(t *testing.T, body string) {
				t.Helper()
				var out paymentuc.EarningsOutput
				require.NoError(t, json.Unmarshal([]byte(body), &out))
				assert.Equal(t, "prov-1", out.ProviderID)
				assert.Equal(t, int64(0), out.AllTime.GrossAmount)
			},
		},
		{
			name: "success - admin views another provider earnings",
			event: func() *events.APIGatewayProxyRequest {
				e := earningsEventWithAuth("admin-1", "admin")
				e.PathParameters = map[string]string{"id": "prov-1"}
				return e
			}(),
			setupMocks: func(pf *mockEarningsHandlerProviderFinder, bl *mockEarningsHandlerBookingLister, pl *mockEarningsHandlerPaymentLister) {
				pf.FindByIDFunc = func(_ context.Context, _ string) (*provider.Provider, error) {
					return testProvider, nil
				}
				bl.FindByProviderFunc = func(_ context.Context, _ string) ([]booking.Booking, error) {
					return nil, nil
				}
			},
			wantStatus: http.StatusOK,
			checkBody: func(t *testing.T, body string) {
				t.Helper()
				var out paymentuc.EarningsOutput
				require.NoError(t, json.Unmarshal([]byte(body), &out))
				assert.Equal(t, "prov-1", out.ProviderID)
			},
		},
		{
			name: "unauthorized - no user ID",
			event: &events.APIGatewayProxyRequest{
				PathParameters: map[string]string{"id": "prov-1"},
			},
			wantStatus:  http.StatusUnauthorized,
			wantErrCode: "UNAUTHORIZED",
		},
		{
			name: "missing provider ID",
			event: func() *events.APIGatewayProxyRequest {
				return earningsEventWithAuth("prov-1", "provider")
			}(),
			wantStatus:  http.StatusBadRequest,
			wantErrCode: "VALIDATION_ERROR",
		},
		{
			name: "forbidden - different provider viewing another",
			event: func() *events.APIGatewayProxyRequest {
				e := earningsEventWithAuth("prov-2", "provider")
				e.PathParameters = map[string]string{"id": "prov-1"}
				return e
			}(),
			wantStatus:  http.StatusForbidden,
			wantErrCode: "FORBIDDEN",
		},
		{
			name: "use case error - provider not found",
			event: func() *events.APIGatewayProxyRequest {
				e := earningsEventWithAuth("prov-missing", "provider")
				e.PathParameters = map[string]string{"id": "prov-missing"}
				return e
			}(),
			setupMocks: func(pf *mockEarningsHandlerProviderFinder, bl *mockEarningsHandlerBookingLister, pl *mockEarningsHandlerPaymentLister) {
				pf.FindByIDFunc = func(_ context.Context, _ string) (*provider.Provider, error) {
					return nil, nil
				}
			},
			wantStatus:  http.StatusNotFound,
			wantErrCode: "NOT_FOUND",
		},
		{
			name: "use case error - internal error from repo",
			event: func() *events.APIGatewayProxyRequest {
				e := earningsEventWithAuth("prov-1", "provider")
				e.PathParameters = map[string]string{"id": "prov-1"}
				return e
			}(),
			setupMocks: func(pf *mockEarningsHandlerProviderFinder, bl *mockEarningsHandlerBookingLister, pl *mockEarningsHandlerPaymentLister) {
				pf.FindByIDFunc = func(_ context.Context, _ string) (*provider.Provider, error) {
					return nil, errors.New("dynamo timeout")
				}
			},
			wantStatus:  http.StatusInternalServerError,
			wantErrCode: "INTERNAL_ERROR",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pf := &mockEarningsHandlerProviderFinder{}
			bl := &mockEarningsHandlerBookingLister{}
			pl := &mockEarningsHandlerPaymentLister{}

			if tt.setupMocks != nil {
				tt.setupMocks(pf, bl, pl)
			}

			uc := paymentuc.NewGetProviderEarningsUseCase(pf, bl, pl)
			h := handler.NewProviderEarningsHandler(uc)

			resp, err := h.Handle(context.Background(), tt.event)

			require.NoError(t, err)
			assert.Equal(t, tt.wantStatus, resp.StatusCode)

			if tt.wantErrCode != "" {
				eb := parseErrorBody(t, resp.Body)
				assert.Equal(t, tt.wantErrCode, eb.Error.Code)
			}

			if tt.checkBody != nil {
				tt.checkBody(t, resp.Body)
			}
		})
	}
}
