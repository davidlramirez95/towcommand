package analytics

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"time"

	"github.com/davidlramirez95/towcommand/internal/domain/event"
)

// EventRecorder processes domain events and updates analytics counters.
type EventRecorder struct {
	repo AnalyticsRecorder
}

// NewEventRecorder creates an EventRecorder with its dependencies.
func NewEventRecorder(repo AnalyticsRecorder) *EventRecorder {
	return &EventRecorder{repo: repo}
}

// Record dispatches a domain event to the appropriate counter update logic.
// Errors are logged but never propagated -- analytics never fails the Lambda.
func (r *EventRecorder) Record(ctx context.Context, eventType string, detail json.RawMessage, eventTime time.Time) error {
	date := eventTime.UTC().Format("2006-01-02")

	switch eventType {
	case event.BookingCreated:
		return r.handleBookingCreated(ctx, date, detail)
	case event.BookingCompleted:
		return r.handleBookingCompleted(ctx, date, detail)
	case event.BookingCancelled:
		return r.handleBookingCancelled(ctx, date, detail)
	case event.PaymentCaptured:
		return r.handlePaymentCaptured(ctx, date, detail)
	default:
		slog.DebugContext(ctx, "unhandled event type for analytics", "event_type", eventType)
		return nil
	}
}

// bookingCreatedDetail is the expected detail shape for BookingCreated events.
type bookingCreatedDetail struct {
	BookingID string  `json:"bookingId"`
	PickupLat float64 `json:"pickupLat"`
	PickupLng float64 `json:"pickupLng"`
}

func (r *EventRecorder) handleBookingCreated(ctx context.Context, date string, detail json.RawMessage) error {
	var d bookingCreatedDetail
	if err := json.Unmarshal(detail, &d); err != nil {
		return fmt.Errorf("unmarshalling BookingCreated detail: %w", err)
	}

	if err := r.repo.IncrementDailyCounter(ctx, date, "totalBookings", 1); err != nil {
		slog.ErrorContext(ctx, "failed to increment totalBookings", "error", err)
		return err
	}

	// Update heatmap cell.
	geohash := Geohash6(d.PickupLat, d.PickupLng)
	if err := r.repo.IncrementHeatmapCell(ctx, date, geohash, d.PickupLat, d.PickupLng); err != nil {
		slog.ErrorContext(ctx, "failed to increment heatmap cell", "error", err, "geohash", geohash)
		return err
	}

	return nil
}

// bookingCompletedAnalyticsDetail is the expected detail shape for BookingCompleted events.
type bookingCompletedAnalyticsDetail struct {
	BookingID  string `json:"bookingId"`
	ProviderID string `json:"providerId"`
}

func (r *EventRecorder) handleBookingCompleted(ctx context.Context, date string, detail json.RawMessage) error {
	var d bookingCompletedAnalyticsDetail
	if err := json.Unmarshal(detail, &d); err != nil {
		return fmt.Errorf("unmarshalling BookingCompleted detail: %w", err)
	}

	if err := r.repo.IncrementDailyCounter(ctx, date, "completedBookings", 1); err != nil {
		slog.ErrorContext(ctx, "failed to increment completedBookings", "error", err)
		return err
	}

	if d.ProviderID != "" {
		if err := r.repo.IncrementProviderCounter(ctx, d.ProviderID, date, "completedJobs", 1); err != nil {
			slog.ErrorContext(ctx, "failed to increment provider completedJobs", "error", err, "provider_id", d.ProviderID)
			return err
		}
	}

	return nil
}

// bookingCancelledAnalyticsDetail is the expected detail shape for BookingCancelled events.
type bookingCancelledAnalyticsDetail struct {
	BookingID  string `json:"bookingId"`
	ProviderID string `json:"providerId"`
}

func (r *EventRecorder) handleBookingCancelled(ctx context.Context, date string, detail json.RawMessage) error {
	var d bookingCancelledAnalyticsDetail
	if err := json.Unmarshal(detail, &d); err != nil {
		return fmt.Errorf("unmarshalling BookingCancelled detail: %w", err)
	}

	if err := r.repo.IncrementDailyCounter(ctx, date, "cancelledBookings", 1); err != nil {
		slog.ErrorContext(ctx, "failed to increment cancelledBookings", "error", err)
		return err
	}

	if d.ProviderID != "" {
		if err := r.repo.IncrementProviderCounter(ctx, d.ProviderID, date, "cancelledJobs", 1); err != nil {
			slog.ErrorContext(ctx, "failed to increment provider cancelledJobs", "error", err, "provider_id", d.ProviderID)
			return err
		}
	}

	return nil
}

// paymentCapturedAnalyticsDetail is the expected detail shape for PaymentCaptured events.
type paymentCapturedAnalyticsDetail struct {
	BookingID      string `json:"bookingId"`
	ProviderID     string `json:"providerId"`
	AmountCentavos int64  `json:"amountCentavos"`
}

func (r *EventRecorder) handlePaymentCaptured(ctx context.Context, date string, detail json.RawMessage) error {
	var d paymentCapturedAnalyticsDetail
	if err := json.Unmarshal(detail, &d); err != nil {
		return fmt.Errorf("unmarshalling PaymentCaptured detail: %w", err)
	}

	if err := r.repo.IncrementDailyCounter(ctx, date, "totalRevenueCentavos", d.AmountCentavos); err != nil {
		slog.ErrorContext(ctx, "failed to increment totalRevenueCentavos", "error", err)
		return err
	}

	if d.ProviderID != "" {
		if err := r.repo.IncrementProviderCounter(ctx, d.ProviderID, date, "totalRevenueCentavos", d.AmountCentavos); err != nil {
			slog.ErrorContext(ctx, "failed to increment provider totalRevenueCentavos", "error", err, "provider_id", d.ProviderID)
			return err
		}
	}

	return nil
}

// Geohash6 returns a simple geohash by truncating lat/lng to 3 decimal places.
// This provides ~111m precision which is sufficient for heatmap cells.
func Geohash6(lat, lng float64) string {
	return fmt.Sprintf("%.3f,%.3f", lat, lng)
}
