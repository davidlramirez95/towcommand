package handler_test

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go-v2/service/apigatewaymanagementapi"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	domainerrors "github.com/davidlramirez95/towcommand/internal/domain/errors"

	"github.com/davidlramirez95/towcommand/internal/adapter/handler"
)

// ---------------------------------------------------------------------------
// Test DTOs
// ---------------------------------------------------------------------------

type wsMessage struct {
	Action string `json:"action" validate:"required"`
	RoomID string `json:"room_id" validate:"required"`
}

// ---------------------------------------------------------------------------
// Mock ConnectionPoster
// ---------------------------------------------------------------------------

type mockConnectionPoster struct {
	postedData   []byte
	postedConnID string
	err          error
}

func (m *mockConnectionPoster) PostToConnection(_ context.Context, params *apigatewaymanagementapi.PostToConnectionInput, _ ...func(*apigatewaymanagementapi.Options)) (*apigatewaymanagementapi.PostToConnectionOutput, error) {
	m.postedData = params.Data
	if params.ConnectionId != nil {
		m.postedConnID = *params.ConnectionId
	}
	return &apigatewaymanagementapi.PostToConnectionOutput{}, m.err
}

// ---------------------------------------------------------------------------
// ParseWSBody
// ---------------------------------------------------------------------------

func TestParseWSBody(t *testing.T) {
	tests := []struct {
		name    string
		body    string
		wantErr bool
		errCode domainerrors.ErrorCode
	}{
		{
			name:    "valid body",
			body:    `{"action":"join","room_id":"room-1"}`,
			wantErr: false,
		},
		{
			name:    "invalid JSON",
			body:    `not-json`,
			wantErr: true,
			errCode: domainerrors.CodeValidationError,
		},
		{
			name:    "missing required field",
			body:    `{"action":"join"}`,
			wantErr: true,
			errCode: domainerrors.CodeValidationError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			event := &events.APIGatewayWebsocketProxyRequest{Body: tt.body}
			result, err := handler.ParseWSBody[wsMessage](event)

			if tt.wantErr {
				require.Error(t, err)
				var appErr *domainerrors.AppError
				require.True(t, errors.As(err, &appErr))
				assert.Equal(t, tt.errCode, appErr.Code)
			} else {
				require.NoError(t, err)
				assert.Equal(t, "join", result.Action)
				assert.Equal(t, "room-1", result.RoomID)
			}
		})
	}
}

// ---------------------------------------------------------------------------
// ExtractConnectionID
// ---------------------------------------------------------------------------

func TestExtractConnectionID(t *testing.T) {
	tests := []struct {
		name string
		id   string
		want string
	}{
		{
			name: "has connection ID",
			id:   "conn-abc-123",
			want: "conn-abc-123",
		},
		{
			name: "empty connection ID",
			id:   "",
			want: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			event := &events.APIGatewayWebsocketProxyRequest{
				RequestContext: events.APIGatewayWebsocketProxyRequestContext{
					ConnectionID: tt.id,
				},
			}
			assert.Equal(t, tt.want, handler.ExtractConnectionID(event))
		})
	}
}

// ---------------------------------------------------------------------------
// SendToConnection
// ---------------------------------------------------------------------------

func TestSendToConnection(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mock := &mockConnectionPoster{}
		data := map[string]string{"message": "hello"}

		err := handler.SendToConnection(context.Background(), mock, "conn-1", data)

		require.NoError(t, err)
		assert.Equal(t, "conn-1", mock.postedConnID)
		assert.Contains(t, string(mock.postedData), `"message":"hello"`)
	})

	t.Run("API error", func(t *testing.T) {
		mock := &mockConnectionPoster{err: fmt.Errorf("gone")}

		err := handler.SendToConnection(context.Background(), mock, "conn-1", "data")

		require.Error(t, err)
		var appErr *domainerrors.AppError
		require.True(t, errors.As(err, &appErr))
		assert.Equal(t, domainerrors.CodeExternalService, appErr.Code)
	})

	t.Run("unmarshalable data", func(t *testing.T) {
		mock := &mockConnectionPoster{}

		err := handler.SendToConnection(context.Background(), mock, "conn-1", make(chan int))

		require.Error(t, err)
		var appErr *domainerrors.AppError
		require.True(t, errors.As(err, &appErr))
		assert.Equal(t, domainerrors.CodeInternalError, appErr.Code)
	})
}
