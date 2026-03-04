package handler_test

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"testing"

	"github.com/aws/aws-lambda-go/events"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	domainerrors "github.com/davidlramirez95/towcommand/internal/domain/errors"

	"github.com/davidlramirez95/towcommand/internal/adapter/handler"
)

// ---------------------------------------------------------------------------
// Test DTOs
// ---------------------------------------------------------------------------

type createBookingRequest struct {
	CustomerID  string `json:"customer_id" validate:"required"`
	ServiceType string `json:"service_type" validate:"required"`
	Notes       string `json:"notes"`
}

type noValidationRequest struct {
	Name string `json:"name"`
}

// ---------------------------------------------------------------------------
// ParseBody
// ---------------------------------------------------------------------------

func TestParseBody(t *testing.T) {
	tests := []struct {
		name    string
		body    string
		wantErr bool
		errCode domainerrors.ErrorCode
	}{
		{
			name:    "valid body with all required fields",
			body:    `{"customer_id":"cust-1","service_type":"FLATBED_TOW","notes":"careful"}`,
			wantErr: false,
		},
		{
			name:    "invalid JSON",
			body:    `{not json}`,
			wantErr: true,
			errCode: domainerrors.CodeValidationError,
		},
		{
			name:    "missing required field",
			body:    `{"customer_id":"cust-1"}`,
			wantErr: true,
			errCode: domainerrors.CodeValidationError,
		},
		{
			name:    "empty body",
			body:    "",
			wantErr: true,
			errCode: domainerrors.CodeValidationError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			event := &events.APIGatewayProxyRequest{Body: tt.body}
			result, err := handler.ParseBody[createBookingRequest](event)

			if tt.wantErr {
				require.Error(t, err)
				var appErr *domainerrors.AppError
				require.True(t, errors.As(err, &appErr))
				assert.Equal(t, tt.errCode, appErr.Code)
			} else {
				require.NoError(t, err)
				assert.Equal(t, "cust-1", result.CustomerID)
				assert.Equal(t, "FLATBED_TOW", result.ServiceType)
				assert.Equal(t, "careful", result.Notes)
			}
		})
	}
}

func TestParseBody_NoValidationTags(t *testing.T) {
	event := &events.APIGatewayProxyRequest{Body: `{"name":"test"}`}
	result, err := handler.ParseBody[noValidationRequest](event)
	require.NoError(t, err)
	assert.Equal(t, "test", result.Name)
}

// ---------------------------------------------------------------------------
// ParsePathParam
// ---------------------------------------------------------------------------

func TestParsePathParam(t *testing.T) {
	tests := []struct {
		name   string
		params map[string]string
		key    string
		want   string
	}{
		{
			name:   "existing key",
			params: map[string]string{"id": "booking-123"},
			key:    "id",
			want:   "booking-123",
		},
		{
			name:   "missing key",
			params: map[string]string{"id": "booking-123"},
			key:    "missing",
			want:   "",
		},
		{
			name:   "nil params",
			params: nil,
			key:    "id",
			want:   "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			event := &events.APIGatewayProxyRequest{PathParameters: tt.params}
			assert.Equal(t, tt.want, handler.ParsePathParam(event, tt.key))
		})
	}
}

// ---------------------------------------------------------------------------
// ParseQueryParam
// ---------------------------------------------------------------------------

func TestParseQueryParam(t *testing.T) {
	tests := []struct {
		name   string
		params map[string]string
		key    string
		want   string
	}{
		{
			name:   "existing key",
			params: map[string]string{"status": "PENDING"},
			key:    "status",
			want:   "PENDING",
		},
		{
			name:   "missing key",
			params: map[string]string{"status": "PENDING"},
			key:    "missing",
			want:   "",
		},
		{
			name:   "nil params",
			params: nil,
			key:    "status",
			want:   "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			event := &events.APIGatewayProxyRequest{QueryStringParameters: tt.params}
			assert.Equal(t, tt.want, handler.ParseQueryParam(event, tt.key))
		})
	}
}

// ---------------------------------------------------------------------------
// ExtractUserID
// ---------------------------------------------------------------------------

func TestExtractUserID(t *testing.T) {
	tests := []struct {
		name       string
		authorizer map[string]interface{}
		want       string
	}{
		{
			name: "valid Cognito claims",
			authorizer: map[string]interface{}{
				"claims": map[string]interface{}{
					"sub":   "user-abc-123",
					"email": "test@example.com",
				},
			},
			want: "user-abc-123",
		},
		{
			name:       "no claims key",
			authorizer: map[string]interface{}{},
			want:       "",
		},
		{
			name: "claims is not a map",
			authorizer: map[string]interface{}{
				"claims": "invalid",
			},
			want: "",
		},
		{
			name: "sub is not a string",
			authorizer: map[string]interface{}{
				"claims": map[string]interface{}{
					"sub": 12345,
				},
			},
			want: "",
		},
		{
			name:       "nil authorizer",
			authorizer: nil,
			want:       "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			event := &events.APIGatewayProxyRequest{
				RequestContext: events.APIGatewayProxyRequestContext{
					Authorizer: tt.authorizer,
				},
			}
			assert.Equal(t, tt.want, handler.ExtractUserID(event))
		})
	}
}

// ---------------------------------------------------------------------------
// SuccessResponse
// ---------------------------------------------------------------------------

func TestSuccessResponse(t *testing.T) {
	type responseData struct {
		ID     string `json:"id"`
		Status string `json:"status"`
	}

	resp := handler.SuccessResponse(http.StatusOK, responseData{ID: "bk-1", Status: "PENDING"})

	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, "application/json", resp.Headers["Content-Type"])
	assert.Equal(t, "*", resp.Headers["Access-Control-Allow-Origin"])

	var body responseData
	require.NoError(t, json.Unmarshal([]byte(resp.Body), &body))
	assert.Equal(t, "bk-1", body.ID)
	assert.Equal(t, "PENDING", body.Status)
}

func TestSuccessResponse_CORSHeaders(t *testing.T) {
	resp := handler.SuccessResponse(http.StatusCreated, map[string]string{"ok": "true"})

	assert.Equal(t, "*", resp.Headers["Access-Control-Allow-Origin"])
	assert.Contains(t, resp.Headers["Access-Control-Allow-Methods"], "POST")
	assert.Contains(t, resp.Headers["Access-Control-Allow-Headers"], "Authorization")
}

func TestSuccessResponse_UnmarshalableBody(t *testing.T) {
	// channels cannot be JSON-marshalled
	resp := handler.SuccessResponse(http.StatusOK, make(chan int))
	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
}

// ---------------------------------------------------------------------------
// ErrorResponse
// ---------------------------------------------------------------------------

func TestErrorResponse(t *testing.T) {
	tests := []struct {
		name       string
		err        error
		wantStatus int
		wantCode   string
		wantMsg    string
	}{
		{
			name:       "AppError validation",
			err:        domainerrors.NewValidationError("name is required"),
			wantStatus: http.StatusBadRequest,
			wantCode:   "VALIDATION_ERROR",
			wantMsg:    "name is required",
		},
		{
			name:       "AppError not found",
			err:        domainerrors.NewNotFoundError("Booking", "bk-1"),
			wantStatus: http.StatusNotFound,
			wantCode:   "NOT_FOUND",
			wantMsg:    "Booking not found",
		},
		{
			name:       "AppError unauthorized",
			err:        domainerrors.NewUnauthorizedError(),
			wantStatus: http.StatusUnauthorized,
			wantCode:   "UNAUTHORIZED",
			wantMsg:    "Unauthorized",
		},
		{
			name:       "wrapped AppError",
			err:        fmt.Errorf("handler: %w", domainerrors.NewForbiddenError("admin only")),
			wantStatus: http.StatusForbidden,
			wantCode:   "FORBIDDEN",
			wantMsg:    "admin only",
		},
		{
			name:       "non-AppError",
			err:        fmt.Errorf("something broke"),
			wantStatus: http.StatusInternalServerError,
			wantCode:   "INTERNAL_ERROR",
			wantMsg:    "An unexpected error occurred",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp := handler.ErrorResponse(tt.err)

			assert.Equal(t, tt.wantStatus, resp.StatusCode)
			assert.Equal(t, "*", resp.Headers["Access-Control-Allow-Origin"])

			var body struct {
				Error struct {
					Code    string `json:"code"`
					Message string `json:"message"`
				} `json:"error"`
			}
			require.NoError(t, json.Unmarshal([]byte(resp.Body), &body))
			assert.Equal(t, tt.wantCode, body.Error.Code)
			assert.Equal(t, tt.wantMsg, body.Error.Message)
		})
	}
}

func TestErrorResponse_CORSHeaders(t *testing.T) {
	resp := handler.ErrorResponse(domainerrors.NewInternalError("boom"))
	assert.Equal(t, "*", resp.Headers["Access-Control-Allow-Origin"])
	assert.Contains(t, resp.Headers["Access-Control-Allow-Methods"], "GET")
}
