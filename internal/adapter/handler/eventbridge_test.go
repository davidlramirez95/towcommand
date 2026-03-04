package handler_test

import (
	"encoding/json"
	"errors"
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

type bookingCreatedDetail struct {
	BookingID   string `json:"booking_id" validate:"required"`
	CustomerID  string `json:"customer_id" validate:"required"`
	ServiceType string `json:"service_type"`
}

// ---------------------------------------------------------------------------
// ParseEventDetail
// ---------------------------------------------------------------------------

func TestParseEventDetail(t *testing.T) {
	tests := []struct {
		name    string
		detail  json.RawMessage
		wantErr bool
		errCode domainerrors.ErrorCode
	}{
		{
			name:    "valid detail",
			detail:  json.RawMessage(`{"booking_id":"bk-1","customer_id":"cust-1","service_type":"FLATBED_TOW"}`),
			wantErr: false,
		},
		{
			name:    "invalid JSON",
			detail:  json.RawMessage(`{broken}`),
			wantErr: true,
			errCode: domainerrors.CodeValidationError,
		},
		{
			name:    "missing required field",
			detail:  json.RawMessage(`{"booking_id":"bk-1"}`),
			wantErr: true,
			errCode: domainerrors.CodeValidationError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			event := &events.CloudWatchEvent{Detail: tt.detail}
			result, err := handler.ParseEventDetail[bookingCreatedDetail](event)

			if tt.wantErr {
				require.Error(t, err)
				var appErr *domainerrors.AppError
				require.True(t, errors.As(err, &appErr))
				assert.Equal(t, tt.errCode, appErr.Code)
			} else {
				require.NoError(t, err)
				assert.Equal(t, "bk-1", result.BookingID)
				assert.Equal(t, "cust-1", result.CustomerID)
				assert.Equal(t, "FLATBED_TOW", result.ServiceType)
			}
		})
	}
}

// ---------------------------------------------------------------------------
// ExtractCorrelationID
// ---------------------------------------------------------------------------

func TestExtractCorrelationID(t *testing.T) {
	tests := []struct {
		name string
		id   string
		want string
	}{
		{
			name: "has event ID",
			id:   "evt-abc-123-def",
			want: "evt-abc-123-def",
		},
		{
			name: "empty event ID",
			id:   "",
			want: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			event := &events.CloudWatchEvent{ID: tt.id}
			assert.Equal(t, tt.want, handler.ExtractCorrelationID(event))
		})
	}
}
