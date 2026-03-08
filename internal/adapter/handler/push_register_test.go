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
	"github.com/davidlramirez95/towcommand/internal/usecase/port"
)

// ---------------------------------------------------------------------------
// Mock implementations
// ---------------------------------------------------------------------------

type mockPushTokenRegistrar struct {
	RegisterFunc   func(ctx context.Context, token *port.PushToken) error
	FindByUserFunc func(ctx context.Context, userID string) ([]port.PushToken, error)
	DeleteFunc     func(ctx context.Context, userID, deviceID string) error
}

func (m *mockPushTokenRegistrar) Register(ctx context.Context, token *port.PushToken) error {
	if m.RegisterFunc != nil {
		return m.RegisterFunc(ctx, token)
	}
	return nil
}

func (m *mockPushTokenRegistrar) FindByUserID(ctx context.Context, userID string) ([]port.PushToken, error) {
	if m.FindByUserFunc != nil {
		return m.FindByUserFunc(ctx, userID)
	}
	return nil, nil
}

func (m *mockPushTokenRegistrar) Delete(ctx context.Context, userID, deviceID string) error {
	if m.DeleteFunc != nil {
		return m.DeleteFunc(ctx, userID, deviceID)
	}
	return nil
}

type mockPushEndpointCreator struct {
	CreateEndpointFunc func(ctx context.Context, platform port.PushPlatform, token string) (string, error)
}

func (m *mockPushEndpointCreator) CreateEndpoint(ctx context.Context, platform port.PushPlatform, token string) (string, error) {
	if m.CreateEndpointFunc != nil {
		return m.CreateEndpointFunc(ctx, platform, token)
	}
	return "arn:aws:sns:ap-southeast-1:123:endpoint/GCM/App/default", nil
}

// ---------------------------------------------------------------------------
// Tests
// ---------------------------------------------------------------------------

func TestRegisterPushTokenHandler(t *testing.T) {
	tests := []struct {
		name        string
		event       *events.APIGatewayProxyRequest
		setupMocks  func(tr *mockPushTokenRegistrar, ec *mockPushEndpointCreator)
		wantStatus  int
		wantErrCode string
		checkBody   func(t *testing.T, body string)
	}{
		{
			name: "success - FCM token",
			event: func() *events.APIGatewayProxyRequest {
				e := apiEventWithAuth("user-1")
				e.PathParameters = map[string]string{"id": "user-1"}
				e.Body = `{"token":"fcm-token-abc","platform":"FCM","deviceId":"device-001"}`
				return e
			}(),
			setupMocks: func(tr *mockPushTokenRegistrar, ec *mockPushEndpointCreator) {
				ec.CreateEndpointFunc = func(_ context.Context, p port.PushPlatform, _ string) (string, error) {
					return "arn:aws:sns:ap-southeast-1:123:endpoint/GCM/App/new-ep", nil
				}
			},
			wantStatus: http.StatusCreated,
			checkBody: func(t *testing.T, body string) {
				t.Helper()
				var token port.PushToken
				require.NoError(t, json.Unmarshal([]byte(body), &token))
				assert.Equal(t, "user-1", token.UserID)
				assert.Equal(t, "fcm-token-abc", token.Token)
				assert.Equal(t, port.PushPlatformFCM, token.Platform)
				assert.Equal(t, "device-001", token.DeviceID)
				assert.Equal(t, "arn:aws:sns:ap-southeast-1:123:endpoint/GCM/App/new-ep", token.EndpointArn)
			},
		},
		{
			name: "success - APNS token",
			event: func() *events.APIGatewayProxyRequest {
				e := apiEventWithAuth("user-2")
				e.PathParameters = map[string]string{"id": "user-2"}
				e.Body = `{"token":"apns-token-xyz","platform":"APNS","deviceId":"iphone-001"}`
				return e
			}(),
			setupMocks: func(tr *mockPushTokenRegistrar, ec *mockPushEndpointCreator) {
				ec.CreateEndpointFunc = func(_ context.Context, p port.PushPlatform, _ string) (string, error) {
					return "arn:aws:sns:ap-southeast-1:123:endpoint/APNS/App/new-ep", nil
				}
			},
			wantStatus: http.StatusCreated,
			checkBody: func(t *testing.T, body string) {
				t.Helper()
				var token port.PushToken
				require.NoError(t, json.Unmarshal([]byte(body), &token))
				assert.Equal(t, port.PushPlatformAPNS, token.Platform)
				assert.Equal(t, "iphone-001", token.DeviceID)
			},
		},
		{
			name: "unauthorized - no user ID",
			event: &events.APIGatewayProxyRequest{
				PathParameters: map[string]string{"id": "user-1"},
				Body:           `{"token":"abc","platform":"FCM","deviceId":"dev-1"}`,
			},
			wantStatus:  http.StatusUnauthorized,
			wantErrCode: "UNAUTHORIZED",
		},
		{
			name: "missing path parameter - user ID",
			event: func() *events.APIGatewayProxyRequest {
				e := apiEventWithAuth("user-1")
				e.Body = `{"token":"abc","platform":"FCM","deviceId":"dev-1"}`
				return e
			}(),
			wantStatus:  http.StatusBadRequest,
			wantErrCode: "VALIDATION_ERROR",
		},
		{
			name: "forbidden - different user",
			event: func() *events.APIGatewayProxyRequest {
				e := apiEventWithAuth("user-1")
				e.PathParameters = map[string]string{"id": "user-2"}
				e.Body = `{"token":"abc","platform":"FCM","deviceId":"dev-1"}`
				return e
			}(),
			wantStatus:  http.StatusForbidden,
			wantErrCode: "FORBIDDEN",
		},
		{
			name: "admin can register for another user",
			event: func() *events.APIGatewayProxyRequest {
				e := &events.APIGatewayProxyRequest{
					PathParameters: map[string]string{"id": "user-2"},
					Body:           `{"token":"abc","platform":"FCM","deviceId":"dev-1"}`,
					RequestContext: events.APIGatewayProxyRequestContext{
						Authorizer: map[string]interface{}{
							"claims": map[string]interface{}{
								"sub":             "admin-user",
								"custom:userType": "admin",
							},
						},
					},
				}
				return e
			}(),
			wantStatus: http.StatusCreated,
		},
		{
			name: "invalid body - bad JSON",
			event: func() *events.APIGatewayProxyRequest {
				e := apiEventWithAuth("user-1")
				e.PathParameters = map[string]string{"id": "user-1"}
				e.Body = `{not json}`
				return e
			}(),
			wantStatus:  http.StatusBadRequest,
			wantErrCode: "VALIDATION_ERROR",
		},
		{
			name: "invalid body - missing token",
			event: func() *events.APIGatewayProxyRequest {
				e := apiEventWithAuth("user-1")
				e.PathParameters = map[string]string{"id": "user-1"}
				e.Body = `{"platform":"FCM","deviceId":"dev-1"}`
				return e
			}(),
			wantStatus:  http.StatusBadRequest,
			wantErrCode: "VALIDATION_ERROR",
		},
		{
			name: "invalid body - missing platform",
			event: func() *events.APIGatewayProxyRequest {
				e := apiEventWithAuth("user-1")
				e.PathParameters = map[string]string{"id": "user-1"}
				e.Body = `{"token":"abc","deviceId":"dev-1"}`
				return e
			}(),
			wantStatus:  http.StatusBadRequest,
			wantErrCode: "VALIDATION_ERROR",
		},
		{
			name: "invalid body - invalid platform value",
			event: func() *events.APIGatewayProxyRequest {
				e := apiEventWithAuth("user-1")
				e.PathParameters = map[string]string{"id": "user-1"}
				e.Body = `{"token":"abc","platform":"WEB_PUSH","deviceId":"dev-1"}`
				return e
			}(),
			wantStatus:  http.StatusBadRequest,
			wantErrCode: "VALIDATION_ERROR",
		},
		{
			name: "invalid body - missing deviceId",
			event: func() *events.APIGatewayProxyRequest {
				e := apiEventWithAuth("user-1")
				e.PathParameters = map[string]string{"id": "user-1"}
				e.Body = `{"token":"abc","platform":"FCM"}`
				return e
			}(),
			wantStatus:  http.StatusBadRequest,
			wantErrCode: "VALIDATION_ERROR",
		},
		{
			name: "endpoint creation error",
			event: func() *events.APIGatewayProxyRequest {
				e := apiEventWithAuth("user-1")
				e.PathParameters = map[string]string{"id": "user-1"}
				e.Body = `{"token":"bad-token","platform":"FCM","deviceId":"dev-1"}`
				return e
			}(),
			setupMocks: func(tr *mockPushTokenRegistrar, ec *mockPushEndpointCreator) {
				ec.CreateEndpointFunc = func(_ context.Context, _ port.PushPlatform, _ string) (string, error) {
					return "", errors.New("SNS: invalid token")
				}
			},
			wantStatus: http.StatusInternalServerError,
		},
		{
			name: "token registration error",
			event: func() *events.APIGatewayProxyRequest {
				e := apiEventWithAuth("user-1")
				e.PathParameters = map[string]string{"id": "user-1"}
				e.Body = `{"token":"abc","platform":"FCM","deviceId":"dev-1"}`
				return e
			}(),
			setupMocks: func(tr *mockPushTokenRegistrar, ec *mockPushEndpointCreator) {
				tr.RegisterFunc = func(_ context.Context, _ *port.PushToken) error {
					return errors.New("dynamo timeout")
				}
			},
			wantStatus:  http.StatusInternalServerError,
			wantErrCode: "INTERNAL_ERROR",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tr := &mockPushTokenRegistrar{}
			ec := &mockPushEndpointCreator{}

			if tt.setupMocks != nil {
				tt.setupMocks(tr, ec)
			}

			h := handler.NewRegisterPushTokenHandler(tr, ec)
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
