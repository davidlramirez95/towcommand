package gateway

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/eventbridge"
	ebtypes "github.com/aws/aws-sdk-go-v2/service/eventbridge/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/davidlramirez95/towcommand/internal/domain/event"
	"github.com/davidlramirez95/towcommand/internal/platform/logger"
	"github.com/davidlramirez95/towcommand/internal/usecase/port"
)

// mockEventBridgeClient implements EventBridgeAPI for testing.
type mockEventBridgeClient struct {
	putEventsFunc func(ctx context.Context, params *eventbridge.PutEventsInput, optFns ...func(*eventbridge.Options)) (*eventbridge.PutEventsOutput, error)
}

func (m *mockEventBridgeClient) PutEvents(ctx context.Context, params *eventbridge.PutEventsInput, optFns ...func(*eventbridge.Options)) (*eventbridge.PutEventsOutput, error) {
	return m.putEventsFunc(ctx, params, optFns...)
}

func newTestLogger() *slog.Logger {
	return slog.New(slog.NewJSONHandler(io.Discard, nil))
}

func TestEventBridgePublisher_Publish(t *testing.T) {
	t.Parallel()

	type bookingDetail struct {
		BookingID   string `json:"bookingId"`
		ServiceType string `json:"serviceType"`
		AmountCents int64  `json:"amountCents"`
	}

	tests := []struct {
		name           string
		source         string
		detailType     string
		detail         any
		actor          *port.Actor
		correlationID  string
		mockOutput     *eventbridge.PutEventsOutput
		mockErr        error
		wantErr        bool
		wantErrContain string
		validate       func(t *testing.T, input *eventbridge.PutEventsInput)
	}{
		{
			name:       "successful publish with actor and correlation ID",
			source:     event.SourceBooking,
			detailType: event.BookingCreated,
			detail: bookingDetail{
				BookingID:   "BK-001",
				ServiceType: "FLATBED_TOW",
				AmountCents: 250000,
			},
			actor: &port.Actor{
				UserID:   "USR-123",
				UserType: "customer",
			},
			correlationID: "corr-abc-123",
			mockOutput: &eventbridge.PutEventsOutput{
				FailedEntryCount: 0,
				Entries: []ebtypes.PutEventsResultEntry{
					{EventId: strPtr("evt-1")},
				},
			},
			validate: func(t *testing.T, input *eventbridge.PutEventsInput) {
				t.Helper()
				require.Len(t, input.Entries, 1)
				entry := input.Entries[0]
				assert.Equal(t, "towcommand-test", *entry.EventBusName)
				assert.Equal(t, event.SourceBooking, *entry.Source)
				assert.Equal(t, event.BookingCreated, *entry.DetailType)
				assert.NotNil(t, entry.Time)

				var envelope eventEnvelope
				require.NoError(t, json.Unmarshal([]byte(*entry.Detail), &envelope))
				assert.Equal(t, event.SourceBooking, envelope.Source)
				assert.Equal(t, event.BookingCreated, envelope.DetailType)
				assert.Equal(t, "corr-abc-123", envelope.Metadata.CorrelationID)
				assert.Equal(t, "1.0", envelope.Metadata.Version)
				assert.NotEmpty(t, envelope.Metadata.EventID)
				assert.NotEmpty(t, envelope.Metadata.Timestamp)
				require.NotNil(t, envelope.Metadata.Actor)
				assert.Equal(t, "USR-123", envelope.Metadata.Actor.UserID)
				assert.Equal(t, "customer", envelope.Metadata.Actor.UserType)

				// Verify the detail is properly serialized.
				detailBytes, err := json.Marshal(envelope.Detail)
				require.NoError(t, err)
				var bd bookingDetail
				require.NoError(t, json.Unmarshal(detailBytes, &bd))
				assert.Equal(t, "BK-001", bd.BookingID)
				assert.Equal(t, "FLATBED_TOW", bd.ServiceType)
				assert.Equal(t, int64(250000), bd.AmountCents)
			},
		},
		{
			name:       "successful publish without actor uses event ID as correlation ID",
			source:     event.SourceTracking,
			detailType: event.LocationUpdated,
			detail:     map[string]any{"lat": 14.5995, "lng": 120.9842},
			actor:      nil,
			mockOutput: &eventbridge.PutEventsOutput{
				FailedEntryCount: 0,
				Entries: []ebtypes.PutEventsResultEntry{
					{EventId: strPtr("evt-2")},
				},
			},
			validate: func(t *testing.T, input *eventbridge.PutEventsInput) {
				t.Helper()
				var envelope eventEnvelope
				require.NoError(t, json.Unmarshal([]byte(*input.Entries[0].Detail), &envelope))
				assert.Nil(t, envelope.Metadata.Actor)
				// When no correlation ID in context, eventID is used.
				assert.Equal(t, envelope.Metadata.EventID, envelope.Metadata.CorrelationID)
			},
		},
		{
			name:           "PutEvents API error returns external service error",
			source:         event.SourcePayment,
			detailType:     event.PaymentInitiated,
			detail:         map[string]string{"paymentId": "PAY-001"},
			mockErr:        fmt.Errorf("connection refused"),
			wantErr:        true,
			wantErrContain: "EventBridge",
		},
		{
			name:       "FailedEntryCount > 0 returns error",
			source:     event.SourceSOS,
			detailType: event.SOSTriggered,
			detail:     map[string]string{"alertId": "SOS-001"},
			mockOutput: &eventbridge.PutEventsOutput{
				FailedEntryCount: 1,
				Entries: []ebtypes.PutEventsResultEntry{
					{
						ErrorCode:    strPtr("InternalFailure"),
						ErrorMessage: strPtr("EventBridge internal error"),
					},
				},
			},
			wantErr:        true,
			wantErrContain: "EventBridge internal error",
		},
		{
			name:       "FailedEntryCount > 0 with no error message",
			source:     event.SourceMatching,
			detailType: event.MatchingStarted,
			detail:     map[string]string{"matchId": "M-001"},
			mockOutput: &eventbridge.PutEventsOutput{
				FailedEntryCount: 1,
				Entries:          []ebtypes.PutEventsResultEntry{{}},
			},
			wantErr:        true,
			wantErrContain: "unknown",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			var capturedInput *eventbridge.PutEventsInput
			mock := &mockEventBridgeClient{
				putEventsFunc: func(_ context.Context, params *eventbridge.PutEventsInput, _ ...func(*eventbridge.Options)) (*eventbridge.PutEventsOutput, error) {
					capturedInput = params
					if tt.mockErr != nil {
						return nil, tt.mockErr
					}
					return tt.mockOutput, nil
				},
			}

			pub := NewEventBridgePublisher(mock, "towcommand-test", newTestLogger())

			ctx := context.Background()
			if tt.correlationID != "" {
				ctx = logger.SetCorrelationID(ctx, tt.correlationID)
			}

			err := pub.Publish(ctx, tt.source, tt.detailType, tt.detail, tt.actor)

			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.wantErrContain)
				return
			}

			require.NoError(t, err)
			if tt.validate != nil {
				tt.validate(t, capturedInput)
			}
		})
	}
}

func TestEventBridgePublisher_ImplementsPort(t *testing.T) {
	t.Parallel()
	var _ port.EventPublisher = (*EventBridgePublisher)(nil)
}

func TestGenerateEventID(t *testing.T) {
	t.Parallel()
	id1 := generateEventID()
	id2 := generateEventID()
	assert.NotEmpty(t, id1)
	assert.NotEmpty(t, id2)
	assert.NotEqual(t, id1, id2)
	// UUID v4 format: 8-4-4-4-12 hex chars.
	assert.Regexp(t, `^[0-9a-f]{8}-[0-9a-f]{4}-4[0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}$`, id1)
}

func strPtr(s string) *string { return &s }
