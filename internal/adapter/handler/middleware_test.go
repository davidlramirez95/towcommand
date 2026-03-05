package handler_test

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/aws/aws-lambda-go/events"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/davidlramirez95/towcommand/internal/adapter/handler"
	"github.com/davidlramirez95/towcommand/internal/platform/logger"
)

// ---------------------------------------------------------------------------
// WithCorrelationID
// ---------------------------------------------------------------------------

func TestWithCorrelationID_FromHeader(t *testing.T) {
	var capturedCtx context.Context
	base := handler.APIGatewayHandler(func(ctx context.Context, _ *events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
		capturedCtx = ctx
		return events.APIGatewayProxyResponse{StatusCode: http.StatusOK}, nil
	})

	wrapped := handler.WithCorrelationID(base)
	event := &events.APIGatewayProxyRequest{
		Headers: map[string]string{"X-Correlation-ID": "req-abc-123"},
	}

	_, err := wrapped(context.Background(), event)
	require.NoError(t, err)

	id, ok := capturedCtx.Value(logger.CorrelationIDKey).(string)
	require.True(t, ok)
	assert.Equal(t, "req-abc-123", id)
}

func TestWithCorrelationID_LowercaseHeader(t *testing.T) {
	var capturedCtx context.Context
	base := handler.APIGatewayHandler(func(ctx context.Context, _ *events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
		capturedCtx = ctx
		return events.APIGatewayProxyResponse{StatusCode: http.StatusOK}, nil
	})

	wrapped := handler.WithCorrelationID(base)
	event := &events.APIGatewayProxyRequest{
		Headers: map[string]string{"x-correlation-id": "req-lower-456"},
	}

	_, err := wrapped(context.Background(), event)
	require.NoError(t, err)

	id, ok := capturedCtx.Value(logger.CorrelationIDKey).(string)
	require.True(t, ok)
	assert.Equal(t, "req-lower-456", id)
}

func TestWithCorrelationID_GeneratesWhenMissing(t *testing.T) {
	var capturedCtx context.Context
	base := handler.APIGatewayHandler(func(ctx context.Context, _ *events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
		capturedCtx = ctx
		return events.APIGatewayProxyResponse{StatusCode: http.StatusOK}, nil
	})

	wrapped := handler.WithCorrelationID(base)
	event := &events.APIGatewayProxyRequest{}

	_, err := wrapped(context.Background(), event)
	require.NoError(t, err)

	id, ok := capturedCtx.Value(logger.CorrelationIDKey).(string)
	require.True(t, ok)
	assert.NotEmpty(t, id)
}

// ---------------------------------------------------------------------------
// WithLogging
// ---------------------------------------------------------------------------

func TestWithLogging_PassesThrough(t *testing.T) {
	called := false
	base := handler.APIGatewayHandler(func(_ context.Context, _ *events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
		called = true
		return events.APIGatewayProxyResponse{StatusCode: http.StatusOK}, nil
	})

	wrapped := handler.WithLogging(base)
	event := &events.APIGatewayProxyRequest{
		HTTPMethod: "GET",
		Path:       "/bookings",
	}

	resp, err := wrapped(context.Background(), event)
	require.NoError(t, err)
	assert.True(t, called)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestWithLogging_PreservesResponse(t *testing.T) {
	base := handler.APIGatewayHandler(func(_ context.Context, _ *events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusCreated,
			Body:       `{"id":"test"}`,
		}, nil
	})

	wrapped := handler.WithLogging(base)
	resp, err := wrapped(context.Background(), &events.APIGatewayProxyRequest{})
	require.NoError(t, err)
	assert.Equal(t, http.StatusCreated, resp.StatusCode)
	assert.Equal(t, `{"id":"test"}`, resp.Body)
}

// ---------------------------------------------------------------------------
// WithRecover
// ---------------------------------------------------------------------------

func TestWithRecover_NoPanic(t *testing.T) {
	base := handler.APIGatewayHandler(func(_ context.Context, _ *events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
		return events.APIGatewayProxyResponse{StatusCode: http.StatusOK}, nil
	})

	wrapped := handler.WithRecover(base)
	resp, err := wrapped(context.Background(), &events.APIGatewayProxyRequest{})
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestWithRecover_CatchesPanic(t *testing.T) {
	base := handler.APIGatewayHandler(func(_ context.Context, _ *events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
		panic("something went wrong")
	})

	wrapped := handler.WithRecover(base)
	resp, err := wrapped(context.Background(), &events.APIGatewayProxyRequest{})

	require.NoError(t, err, "WithRecover should not return an error")
	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)

	var body struct {
		Error struct {
			Code    string `json:"code"`
			Message string `json:"message"`
		} `json:"error"`
	}
	require.NoError(t, json.Unmarshal([]byte(resp.Body), &body))
	assert.Equal(t, "INTERNAL_ERROR", body.Error.Code)
	assert.Equal(t, "internal server error", body.Error.Message)
}

func TestWithRecover_CatchesPanicWithError(t *testing.T) {
	base := handler.APIGatewayHandler(func(_ context.Context, _ *events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
		panic("nil pointer dereference")
	})

	wrapped := handler.WithRecover(base)
	resp, err := wrapped(context.Background(), &events.APIGatewayProxyRequest{})

	require.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
}

// ---------------------------------------------------------------------------
// Middleware composition
// ---------------------------------------------------------------------------

func TestMiddlewareComposition(t *testing.T) {
	var capturedCtx context.Context
	base := handler.APIGatewayHandler(func(ctx context.Context, _ *events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
		capturedCtx = ctx
		return events.APIGatewayProxyResponse{StatusCode: http.StatusOK}, nil
	})

	composed := handler.WithLogging(
		handler.WithCorrelationID(
			handler.WithRecover(base),
		),
	)

	event := &events.APIGatewayProxyRequest{
		HTTPMethod: "POST",
		Path:       "/bookings",
		Headers:    map[string]string{"X-Correlation-ID": "test-corr-id"},
	}

	resp, err := composed(context.Background(), event)

	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	id, ok := capturedCtx.Value(logger.CorrelationIDKey).(string)
	require.True(t, ok)
	assert.Equal(t, "test-corr-id", id)
}

func TestMiddlewareComposition_PanicRecovery(t *testing.T) {
	base := handler.APIGatewayHandler(func(_ context.Context, _ *events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
		panic("boom")
	})

	composed := handler.WithLogging(
		handler.WithCorrelationID(
			handler.WithRecover(base),
		),
	)

	resp, err := composed(context.Background(), &events.APIGatewayProxyRequest{})

	require.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
}

// ---------------------------------------------------------------------------
// RequireRole
// ---------------------------------------------------------------------------

func TestRequireRole_AllowedRole(t *testing.T) {
	called := false
	base := handler.APIGatewayHandler(func(_ context.Context, _ *events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
		called = true
		return events.APIGatewayProxyResponse{StatusCode: http.StatusOK}, nil
	})

	wrapped := handler.RequireRole("admin", "ops_agent")(base)
	event := &events.APIGatewayProxyRequest{
		RequestContext: events.APIGatewayProxyRequestContext{
			Authorizer: map[string]interface{}{
				"claims": map[string]interface{}{
					"sub":             "user-123",
					"custom:userType": "admin",
				},
			},
		},
	}

	resp, err := wrapped(context.Background(), event)
	require.NoError(t, err)
	assert.True(t, called)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestRequireRole_DeniedRole(t *testing.T) {
	called := false
	base := handler.APIGatewayHandler(func(_ context.Context, _ *events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
		called = true
		return events.APIGatewayProxyResponse{StatusCode: http.StatusOK}, nil
	})

	wrapped := handler.RequireRole("admin", "ops_agent")(base)
	event := &events.APIGatewayProxyRequest{
		RequestContext: events.APIGatewayProxyRequestContext{
			Authorizer: map[string]interface{}{
				"claims": map[string]interface{}{
					"sub":             "user-456",
					"custom:userType": "customer",
				},
			},
		},
	}

	resp, err := wrapped(context.Background(), event)
	require.NoError(t, err)
	assert.False(t, called)
	assert.Equal(t, http.StatusForbidden, resp.StatusCode)

	var body struct {
		Error struct {
			Code    string `json:"code"`
			Message string `json:"message"`
		} `json:"error"`
	}
	require.NoError(t, json.Unmarshal([]byte(resp.Body), &body))
	assert.Equal(t, "FORBIDDEN", body.Error.Code)
	assert.Equal(t, "insufficient permissions", body.Error.Message)
}

func TestRequireRole_MissingUserType(t *testing.T) {
	called := false
	base := handler.APIGatewayHandler(func(_ context.Context, _ *events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
		called = true
		return events.APIGatewayProxyResponse{StatusCode: http.StatusOK}, nil
	})

	wrapped := handler.RequireRole("admin", "ops_agent")(base)
	event := &events.APIGatewayProxyRequest{
		RequestContext: events.APIGatewayProxyRequestContext{
			Authorizer: map[string]interface{}{
				"claims": map[string]interface{}{
					"sub": "user-789",
				},
			},
		},
	}

	resp, err := wrapped(context.Background(), event)
	require.NoError(t, err)
	assert.False(t, called)
	assert.Equal(t, http.StatusForbidden, resp.StatusCode)

	var body struct {
		Error struct {
			Code    string `json:"code"`
			Message string `json:"message"`
		} `json:"error"`
	}
	require.NoError(t, json.Unmarshal([]byte(resp.Body), &body))
	assert.Equal(t, "FORBIDDEN", body.Error.Code)
	assert.Equal(t, "missing user type", body.Error.Message)
}

func TestRequireRole_MultipleAllowedRoles(t *testing.T) {
	called := false
	base := handler.APIGatewayHandler(func(_ context.Context, _ *events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
		called = true
		return events.APIGatewayProxyResponse{StatusCode: http.StatusOK}, nil
	})

	wrapped := handler.RequireRole("admin", "ops_agent")(base)
	event := &events.APIGatewayProxyRequest{
		RequestContext: events.APIGatewayProxyRequestContext{
			Authorizer: map[string]interface{}{
				"claims": map[string]interface{}{
					"sub":             "user-ops",
					"custom:userType": "ops_agent",
				},
			},
		},
	}

	resp, err := wrapped(context.Background(), event)
	require.NoError(t, err)
	assert.True(t, called)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}
