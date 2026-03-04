package port

import (
	"context"

	"github.com/davidlramirez95/towcommand/internal/domain/evidence"
)

// EvidenceSaver persists a new condition report.
type EvidenceSaver interface {
	Save(ctx context.Context, r *evidence.ConditionReport) error
}

// EvidenceByBookingLister lists condition reports for a given booking.
type EvidenceByBookingLister interface {
	FindByBooking(ctx context.Context, bookingID string) ([]evidence.ConditionReport, error)
}

// MediaItemAdder adds a media item to a booking's evidence collection.
type MediaItemAdder interface {
	AddMediaItem(ctx context.Context, bookingID string, item *evidence.MediaItem) error
}
