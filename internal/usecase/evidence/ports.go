// Package evidenceuc implements evidence use cases following CLEAN architecture.
// Each use case declares only the port interfaces it needs (ISP).
package evidenceuc

import (
	"context"
	"time"

	"github.com/davidlramirez95/towcommand/internal/domain/booking"
	"github.com/davidlramirez95/towcommand/internal/domain/evidence"
	"github.com/davidlramirez95/towcommand/internal/usecase/port"
)

// BookingFinder retrieves a booking by its ID. Returns nil if not found.
type BookingFinder interface {
	FindByID(ctx context.Context, bookingID string) (*booking.Booking, error)
}

// PresignedURLGenerator generates presigned URLs for S3 uploads.
type PresignedURLGenerator interface {
	GenerateUploadURL(ctx context.Context, key, contentType string, expiry time.Duration) (string, error)
}

// ImageValidator validates images against domain-specific criteria.
type ImageValidator interface {
	ValidateVehiclePhoto(ctx context.Context, s3Bucket, s3Key string) (*port.ImageValidationResult, error)
}

// MediaItemAdder adds a media item to a booking's evidence collection.
type MediaItemAdder interface {
	AddMediaItem(ctx context.Context, bookingID string, item *evidence.MediaItem) error
}

// EvidenceSaver persists a new condition report.
type EvidenceSaver interface {
	Save(ctx context.Context, r *evidence.ConditionReport) error
}

// EvidenceByBookingLister lists condition reports for a given booking.
type EvidenceByBookingLister interface {
	FindByBooking(ctx context.Context, bookingID string) ([]evidence.ConditionReport, error)
}

// EventPublisher publishes domain events to an event bus.
type EventPublisher interface {
	Publish(ctx context.Context, source, detailType string, detail any, actor *port.Actor) error
}

// Actor is a type alias for port.Actor to avoid import stuttering in use case code.
type Actor = port.Actor
