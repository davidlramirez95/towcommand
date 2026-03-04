package handler

import (
	"encoding/json"

	"github.com/aws/aws-lambda-go/events"

	domainerrors "github.com/davidlramirez95/towcommand/internal/domain/errors"
)

// ParseEventDetail unmarshals and validates the Detail field of a CloudWatch Event.
func ParseEventDetail[T any](event *events.CloudWatchEvent) (T, error) {
	var detail T
	if err := json.Unmarshal(event.Detail, &detail); err != nil {
		return detail, domainerrors.NewValidationError("invalid event detail JSON").WithCause(err)
	}
	if err := validate.Struct(detail); err != nil {
		return detail, domainerrors.NewValidationError(err.Error()).WithCause(err)
	}
	return detail, nil
}

// ExtractCorrelationID returns the unique event ID from a CloudWatch Event,
// suitable for use as a correlation ID for tracing.
func ExtractCorrelationID(event *events.CloudWatchEvent) string {
	return event.ID
}
