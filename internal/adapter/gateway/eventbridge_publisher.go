// Package gateway provides adapter implementations for external service gateways.
package gateway

import (
	"context"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"log/slog"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/eventbridge"
	ebtypes "github.com/aws/aws-sdk-go-v2/service/eventbridge/types"

	"github.com/davidlramirez95/towcommand/internal/domain/errors"
	"github.com/davidlramirez95/towcommand/internal/platform/logger"
	"github.com/davidlramirez95/towcommand/internal/usecase/port"
)

// EventBridgeAPI is the subset of the EventBridge client needed by the publisher.
type EventBridgeAPI interface {
	PutEvents(ctx context.Context, params *eventbridge.PutEventsInput, optFns ...func(*eventbridge.Options)) (*eventbridge.PutEventsOutput, error)
}

// eventEnvelope wraps the event detail with metadata, matching the TypeScript TowCommandEvent shape.
type eventEnvelope struct {
	Source     string        `json:"source"`
	DetailType string        `json:"detailType"`
	Detail     any           `json:"detail"`
	Metadata   eventMetadata `json:"metadata"`
}

type eventMetadata struct {
	EventID       string     `json:"eventId"`
	CorrelationID string     `json:"correlationId"`
	Timestamp     string     `json:"timestamp"`
	Version       string     `json:"version"`
	Actor         *actorInfo `json:"actor,omitempty"`
}

type actorInfo struct {
	UserID   string `json:"userId"`
	UserType string `json:"userType"`
}

// EventBridgePublisher publishes domain events to AWS EventBridge.
// It implements the port.EventPublisher interface.
type EventBridgePublisher struct {
	client       EventBridgeAPI
	eventBusName string
	logger       *slog.Logger
}

// NewEventBridgePublisher creates a new EventBridge publisher adapter.
func NewEventBridgePublisher(client EventBridgeAPI, eventBusName string, log *slog.Logger) *EventBridgePublisher {
	return &EventBridgePublisher{
		client:       client,
		eventBusName: eventBusName,
		logger:       log,
	}
}

// Publish sends a domain event to EventBridge.
// The detail is JSON-marshalled into the EventBridge Detail field as a string.
// correlation_id is extracted from ctx; if absent, the generated event ID is used.
func (p *EventBridgePublisher) Publish(ctx context.Context, source, detailType string, detail any, actor *port.Actor) error {
	correlationID, _ := ctx.Value(logger.CorrelationIDKey).(string)
	eventID := generateEventID()
	if correlationID == "" {
		correlationID = eventID
	}

	metadata := eventMetadata{
		EventID:       eventID,
		CorrelationID: correlationID,
		Timestamp:     time.Now().UTC().Format(time.RFC3339),
		Version:       "1.0",
	}
	if actor != nil {
		metadata.Actor = &actorInfo{
			UserID:   actor.UserID,
			UserType: actor.UserType,
		}
	}

	envelope := eventEnvelope{
		Source:     source,
		DetailType: detailType,
		Detail:     detail,
		Metadata:   metadata,
	}

	detailJSON, err := json.Marshal(envelope)
	if err != nil {
		return errors.NewInternalError("failed to marshal event detail").WithCause(err)
	}

	now := time.Now()
	detailStr := string(detailJSON)

	output, err := p.client.PutEvents(ctx, &eventbridge.PutEventsInput{
		Entries: []ebtypes.PutEventsRequestEntry{
			{
				EventBusName: &p.eventBusName,
				Source:       &source,
				DetailType:   &detailType,
				Detail:       &detailStr,
				Time:         &now,
			},
		},
	})
	if err != nil {
		return errors.NewExternalServiceError("EventBridge", err)
	}

	if output.FailedEntryCount > 0 {
		errMsg := "unknown"
		if len(output.Entries) > 0 && output.Entries[0].ErrorMessage != nil {
			errMsg = *output.Entries[0].ErrorMessage
		}
		p.logger.ErrorContext(ctx, "EventBridge publish failed",
			"failed_count", output.FailedEntryCount,
			"source", source,
			"detail_type", detailType,
			"error_message", errMsg,
		)
		return errors.NewExternalServiceError("EventBridge",
			fmt.Errorf("failed to publish event: %s", errMsg))
	}

	p.logger.InfoContext(ctx, "event published",
		"event_id", eventID,
		"source", source,
		"detail_type", detailType,
	)

	return nil
}

// generateEventID produces a random UUID v4 string for event identification.
func generateEventID() string {
	b := make([]byte, 16)
	_, _ = rand.Read(b)
	// Set version 4 and variant bits per RFC 4122.
	b[6] = (b[6] & 0x0f) | 0x40
	b[8] = (b[8] & 0x3f) | 0x80
	return fmt.Sprintf("%08x-%04x-%04x-%04x-%012x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
}
