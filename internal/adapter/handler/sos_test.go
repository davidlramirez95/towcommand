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

	"github.com/davidlramirez95/towcommand/internal/adapter/handler"
	"github.com/davidlramirez95/towcommand/internal/domain/booking"
	"github.com/davidlramirez95/towcommand/internal/domain/provider"
	"github.com/davidlramirez95/towcommand/internal/domain/safety"
	"github.com/davidlramirez95/towcommand/internal/usecase/port"
	safetyuc "github.com/davidlramirez95/towcommand/internal/usecase/safety"
)

// ---------------------------------------------------------------------------
// Mock implementations for safety use case dependencies
// ---------------------------------------------------------------------------

type mockSOSSaver struct {
	SaveFunc func(ctx context.Context, alert *safety.SOSAlert) error
}

func (m *mockSOSSaver) Save(ctx context.Context, alert *safety.SOSAlert) error {
	if m.SaveFunc != nil {
		return m.SaveFunc(ctx, alert)
	}
	return nil
}

type mockSOSFinder struct {
	FindByIDFunc func(ctx context.Context, alertID string) (*safety.SOSAlert, error)
}

func (m *mockSOSFinder) FindByID(ctx context.Context, alertID string) (*safety.SOSAlert, error) {
	if m.FindByIDFunc != nil {
		return m.FindByIDFunc(ctx, alertID)
	}
	return nil, nil
}

type mockSOSResolver struct {
	ResolveFunc func(ctx context.Context, alertID, resolvedBy string, resolvedAt time.Time) error
}

func (m *mockSOSResolver) Resolve(ctx context.Context, alertID, resolvedBy string, resolvedAt time.Time) error {
	if m.ResolveFunc != nil {
		return m.ResolveFunc(ctx, alertID, resolvedBy, resolvedAt)
	}
	return nil
}

type mockSOSActiveLister struct {
	FindActiveFunc func(ctx context.Context, limit int32) ([]safety.SOSAlert, error)
}

func (m *mockSOSActiveLister) FindActive(ctx context.Context, limit int32) ([]safety.SOSAlert, error) {
	if m.FindActiveFunc != nil {
		return m.FindActiveFunc(ctx, limit)
	}
	return nil, nil
}

type mockSafetyBookingFinder struct {
	FindByIDFunc func(ctx context.Context, bookingID string) (*booking.Booking, error)
}

func (m *mockSafetyBookingFinder) FindByID(ctx context.Context, bookingID string) (*booking.Booking, error) {
	if m.FindByIDFunc != nil {
		return m.FindByIDFunc(ctx, bookingID)
	}
	return nil, nil
}

type mockSafetyProviderFinder struct {
	FindByIDFunc func(ctx context.Context, providerID string) (*provider.Provider, error)
}

func (m *mockSafetyProviderFinder) FindByID(ctx context.Context, providerID string) (*provider.Provider, error) {
	if m.FindByIDFunc != nil {
		return m.FindByIDFunc(ctx, providerID)
	}
	return nil, nil
}

type mockSafetyEventPublisher struct {
	PublishFunc func(ctx context.Context, source, detailType string, detail any, actor *port.Actor) error
}

func (m *mockSafetyEventPublisher) Publish(ctx context.Context, source, detailType string, detail any, actor *port.Actor) error {
	if m.PublishFunc != nil {
		return m.PublishFunc(ctx, source, detailType, detail, actor)
	}
	return nil
}

// ---------------------------------------------------------------------------
// Helper: build API Gateway event with Cognito auth and user type
// ---------------------------------------------------------------------------

func apiEventWithAuthAndRole(userID, userType string) *events.APIGatewayProxyRequest {
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
// TriggerSOSHandler tests
// ---------------------------------------------------------------------------

func TestTriggerSOSHandler(t *testing.T) {
	testBooking := &booking.Booking{
		BookingID:  "bk-1",
		CustomerID: "user-1",
		ProviderID: "prov-1",
		DropoffLocation: booking.GeoLocation{
			Lat: 14.5995,
			Lng: 120.9842,
		},
	}

	testProvider := &provider.Provider{
		ProviderID: "prov-1",
		TrustTier:  "verified",
	}

	tests := []struct {
		name        string
		event       *events.APIGatewayProxyRequest
		setupMocks  func(ss *mockSOSSaver, bf *mockSafetyBookingFinder, pf *mockSafetyProviderFinder, ep *mockSafetyEventPublisher)
		wantStatus  int
		wantErrCode string
		checkBody   func(t *testing.T, body string)
	}{
		{
			name: "success",
			event: func() *events.APIGatewayProxyRequest {
				e := apiEventWithAuth("user-1")
				e.PathParameters = map[string]string{"id": "bk-1"}
				e.Body = `{"triggerType":"TRIPLE_TAP","lat":14.5995,"lng":120.9842}`
				return e
			}(),
			setupMocks: func(ss *mockSOSSaver, bf *mockSafetyBookingFinder, pf *mockSafetyProviderFinder, ep *mockSafetyEventPublisher) {
				bf.FindByIDFunc = func(_ context.Context, _ string) (*booking.Booking, error) {
					return testBooking, nil
				}
				pf.FindByIDFunc = func(_ context.Context, _ string) (*provider.Provider, error) {
					return testProvider, nil
				}
			},
			wantStatus: http.StatusCreated,
			checkBody: func(t *testing.T, body string) {
				t.Helper()
				var alert safety.SOSAlert
				require.NoError(t, json.Unmarshal([]byte(body), &alert))
				assert.Equal(t, "bk-1", alert.BookingID)
				assert.Equal(t, "user-1", alert.TriggeredBy)
				assert.Equal(t, safety.TriggerTypeTripleTap, alert.TriggerType)
				assert.NotEmpty(t, alert.AlertID)
			},
		},
		{
			name: "unauthorized - no user ID",
			event: &events.APIGatewayProxyRequest{
				PathParameters: map[string]string{"id": "bk-1"},
				Body:           `{"triggerType":"TRIPLE_TAP","lat":14.5995,"lng":120.9842}`,
			},
			wantStatus:  http.StatusUnauthorized,
			wantErrCode: "UNAUTHORIZED",
		},
		{
			name: "missing booking ID",
			event: func() *events.APIGatewayProxyRequest {
				e := apiEventWithAuth("user-1")
				e.Body = `{"triggerType":"TRIPLE_TAP","lat":14.5995,"lng":120.9842}`
				return e
			}(),
			wantStatus:  http.StatusBadRequest,
			wantErrCode: "VALIDATION_ERROR",
		},
		{
			name: "invalid body - bad JSON",
			event: func() *events.APIGatewayProxyRequest {
				e := apiEventWithAuth("user-1")
				e.PathParameters = map[string]string{"id": "bk-1"}
				e.Body = `{not json}`
				return e
			}(),
			wantStatus:  http.StatusBadRequest,
			wantErrCode: "VALIDATION_ERROR",
		},
		{
			name: "invalid body - missing trigger type",
			event: func() *events.APIGatewayProxyRequest {
				e := apiEventWithAuth("user-1")
				e.PathParameters = map[string]string{"id": "bk-1"}
				e.Body = `{"lat":14.5995,"lng":120.9842}`
				return e
			}(),
			wantStatus:  http.StatusBadRequest,
			wantErrCode: "VALIDATION_ERROR",
		},
		{
			name: "use case error - booking not found",
			event: func() *events.APIGatewayProxyRequest {
				e := apiEventWithAuth("user-1")
				e.PathParameters = map[string]string{"id": "bk-missing"}
				e.Body = `{"triggerType":"BUTTON","lat":14.5995,"lng":120.9842}`
				return e
			}(),
			setupMocks: func(ss *mockSOSSaver, bf *mockSafetyBookingFinder, pf *mockSafetyProviderFinder, ep *mockSafetyEventPublisher) {
				bf.FindByIDFunc = func(_ context.Context, _ string) (*booking.Booking, error) {
					return nil, nil
				}
			},
			wantStatus:  http.StatusNotFound,
			wantErrCode: "NOT_FOUND",
		},
		{
			name: "use case error - internal error from repo",
			event: func() *events.APIGatewayProxyRequest {
				e := apiEventWithAuth("user-1")
				e.PathParameters = map[string]string{"id": "bk-1"}
				e.Body = `{"triggerType":"SHAKE","lat":14.5995,"lng":120.9842}`
				return e
			}(),
			setupMocks: func(ss *mockSOSSaver, bf *mockSafetyBookingFinder, pf *mockSafetyProviderFinder, ep *mockSafetyEventPublisher) {
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
			ss := &mockSOSSaver{}
			bf := &mockSafetyBookingFinder{}
			pf := &mockSafetyProviderFinder{}
			ep := &mockSafetyEventPublisher{}

			if tt.setupMocks != nil {
				tt.setupMocks(ss, bf, pf, ep)
			}

			uc := safetyuc.NewTriggerSOSUseCase(bf, pf, ss, ep)
			h := handler.NewTriggerSOSHandler(uc)

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

// ---------------------------------------------------------------------------
// ResolveSOSHandler tests
// ---------------------------------------------------------------------------

func TestResolveSOSHandler(t *testing.T) {
	now := time.Now().UTC()

	activeAlert := &safety.SOSAlert{
		AlertID:     "SOS-2026-abc123",
		BookingID:   "bk-1",
		TriggeredBy: "user-1",
		TriggerType: safety.TriggerTypeButton,
		Lat:         14.5995,
		Lng:         120.9842,
		Resolved:    false,
		Timestamp:   now,
	}

	tests := []struct {
		name        string
		event       *events.APIGatewayProxyRequest
		setupMocks  func(sf *mockSOSFinder, sr *mockSOSResolver, ep *mockSafetyEventPublisher)
		wantStatus  int
		wantErrCode string
		checkBody   func(t *testing.T, body string)
	}{
		{
			name: "success",
			event: func() *events.APIGatewayProxyRequest {
				e := apiEventWithAuth("admin-1")
				e.PathParameters = map[string]string{"id": "SOS-2026-abc123"}
				return e
			}(),
			setupMocks: func(sf *mockSOSFinder, sr *mockSOSResolver, ep *mockSafetyEventPublisher) {
				sf.FindByIDFunc = func(_ context.Context, _ string) (*safety.SOSAlert, error) {
					return activeAlert, nil
				}
			},
			wantStatus: http.StatusOK,
			checkBody: func(t *testing.T, body string) {
				t.Helper()
				var alert safety.SOSAlert
				require.NoError(t, json.Unmarshal([]byte(body), &alert))
				assert.Equal(t, "SOS-2026-abc123", alert.AlertID)
				assert.True(t, alert.Resolved)
				assert.Equal(t, "admin-1", alert.ResolvedBy)
			},
		},
		{
			name: "unauthorized - no user ID",
			event: &events.APIGatewayProxyRequest{
				PathParameters: map[string]string{"id": "SOS-2026-abc123"},
			},
			wantStatus:  http.StatusUnauthorized,
			wantErrCode: "UNAUTHORIZED",
		},
		{
			name: "missing alert ID",
			event: func() *events.APIGatewayProxyRequest {
				return apiEventWithAuth("admin-1")
			}(),
			wantStatus:  http.StatusBadRequest,
			wantErrCode: "VALIDATION_ERROR",
		},
		{
			name: "use case error - alert not found",
			event: func() *events.APIGatewayProxyRequest {
				e := apiEventWithAuth("admin-1")
				e.PathParameters = map[string]string{"id": "SOS-missing"}
				return e
			}(),
			setupMocks: func(sf *mockSOSFinder, sr *mockSOSResolver, ep *mockSafetyEventPublisher) {
				sf.FindByIDFunc = func(_ context.Context, _ string) (*safety.SOSAlert, error) {
					return nil, nil
				}
			},
			wantStatus:  http.StatusNotFound,
			wantErrCode: "NOT_FOUND",
		},
		{
			name: "use case error - already resolved",
			event: func() *events.APIGatewayProxyRequest {
				e := apiEventWithAuth("admin-1")
				e.PathParameters = map[string]string{"id": "SOS-2026-abc123"}
				return e
			}(),
			setupMocks: func(sf *mockSOSFinder, sr *mockSOSResolver, ep *mockSafetyEventPublisher) {
				sf.FindByIDFunc = func(_ context.Context, _ string) (*safety.SOSAlert, error) {
					resolved := *activeAlert
					resolved.Resolved = true
					return &resolved, nil
				}
			},
			wantStatus:  http.StatusConflict,
			wantErrCode: "CONFLICT",
		},
		{
			name: "use case error - internal error from repo",
			event: func() *events.APIGatewayProxyRequest {
				e := apiEventWithAuth("admin-1")
				e.PathParameters = map[string]string{"id": "SOS-2026-abc123"}
				return e
			}(),
			setupMocks: func(sf *mockSOSFinder, sr *mockSOSResolver, ep *mockSafetyEventPublisher) {
				sf.FindByIDFunc = func(_ context.Context, _ string) (*safety.SOSAlert, error) {
					return nil, errors.New("dynamo timeout")
				}
			},
			wantStatus:  http.StatusInternalServerError,
			wantErrCode: "INTERNAL_ERROR",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sf := &mockSOSFinder{}
			sr := &mockSOSResolver{}
			ep := &mockSafetyEventPublisher{}

			if tt.setupMocks != nil {
				tt.setupMocks(sf, sr, ep)
			}

			uc := safetyuc.NewResolveSOSUseCase(sf, sr, ep)
			h := handler.NewResolveSOSHandler(uc)

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

// ---------------------------------------------------------------------------
// AdminActiveSOSHandler tests
// ---------------------------------------------------------------------------

func TestAdminActiveSOSHandler(t *testing.T) {
	now := time.Now().UTC()

	activeAlerts := []safety.SOSAlert{
		{
			AlertID:     "SOS-2026-111",
			BookingID:   "bk-1",
			TriggeredBy: "user-1",
			TriggerType: safety.TriggerTypeTripleTap,
			Lat:         14.5995,
			Lng:         120.9842,
			Resolved:    false,
			Timestamp:   now,
		},
		{
			AlertID:     "SOS-2026-222",
			BookingID:   "bk-2",
			TriggeredBy: "user-2",
			TriggerType: safety.TriggerTypeButton,
			Lat:         14.6000,
			Lng:         120.9850,
			Resolved:    false,
			Timestamp:   now.Add(-5 * time.Minute),
		},
	}

	tests := []struct {
		name        string
		event       *events.APIGatewayProxyRequest
		setupMocks  func(lister *mockSOSActiveLister)
		wantStatus  int
		wantErrCode string
		checkBody   func(t *testing.T, body string)
	}{
		{
			name: "success - default limit",
			event: func() *events.APIGatewayProxyRequest {
				return apiEventWithAuth("admin-1")
			}(),
			setupMocks: func(lister *mockSOSActiveLister) {
				lister.FindActiveFunc = func(_ context.Context, limit int32) ([]safety.SOSAlert, error) {
					assert.Equal(t, int32(50), limit)
					return activeAlerts, nil
				}
			},
			wantStatus: http.StatusOK,
			checkBody: func(t *testing.T, body string) {
				t.Helper()
				var alerts []safety.SOSAlert
				require.NoError(t, json.Unmarshal([]byte(body), &alerts))
				assert.Len(t, alerts, 2)
				assert.Equal(t, "SOS-2026-111", alerts[0].AlertID)
			},
		},
		{
			name: "success - custom limit",
			event: func() *events.APIGatewayProxyRequest {
				e := apiEventWithAuth("admin-1")
				e.QueryStringParameters = map[string]string{"limit": "10"}
				return e
			}(),
			setupMocks: func(lister *mockSOSActiveLister) {
				lister.FindActiveFunc = func(_ context.Context, limit int32) ([]safety.SOSAlert, error) {
					assert.Equal(t, int32(10), limit)
					return activeAlerts[:1], nil
				}
			},
			wantStatus: http.StatusOK,
			checkBody: func(t *testing.T, body string) {
				t.Helper()
				var alerts []safety.SOSAlert
				require.NoError(t, json.Unmarshal([]byte(body), &alerts))
				assert.Len(t, alerts, 1)
			},
		},
		{
			name:        "unauthorized - no user ID",
			event:       &events.APIGatewayProxyRequest{},
			wantStatus:  http.StatusUnauthorized,
			wantErrCode: "UNAUTHORIZED",
		},
		{
			name: "invalid limit - not a number",
			event: func() *events.APIGatewayProxyRequest {
				e := apiEventWithAuth("admin-1")
				e.QueryStringParameters = map[string]string{"limit": "abc"}
				return e
			}(),
			wantStatus:  http.StatusBadRequest,
			wantErrCode: "VALIDATION_ERROR",
		},
		{
			name: "invalid limit - zero",
			event: func() *events.APIGatewayProxyRequest {
				e := apiEventWithAuth("admin-1")
				e.QueryStringParameters = map[string]string{"limit": "0"}
				return e
			}(),
			wantStatus:  http.StatusBadRequest,
			wantErrCode: "VALIDATION_ERROR",
		},
		{
			name: "use case error - internal error from repo",
			event: func() *events.APIGatewayProxyRequest {
				return apiEventWithAuth("admin-1")
			}(),
			setupMocks: func(lister *mockSOSActiveLister) {
				lister.FindActiveFunc = func(_ context.Context, _ int32) ([]safety.SOSAlert, error) {
					return nil, errors.New("dynamo timeout")
				}
			},
			wantStatus:  http.StatusInternalServerError,
			wantErrCode: "INTERNAL_ERROR",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lister := &mockSOSActiveLister{}

			if tt.setupMocks != nil {
				tt.setupMocks(lister)
			}

			uc := safetyuc.NewListActiveSOSUseCase(lister)
			h := handler.NewAdminActiveSOSHandler(uc)

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
