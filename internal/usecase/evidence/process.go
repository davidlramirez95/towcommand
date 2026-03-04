package evidenceuc

import (
	"context"
	"crypto/rand"
	"fmt"
	"time"

	domainerrors "github.com/davidlramirez95/towcommand/internal/domain/errors"
	"github.com/davidlramirez95/towcommand/internal/domain/evidence"
)

const (
	eventSourceEvidence    = "towcommand.evidence"
	eventEvidenceValidated = "EvidenceValidated"
)

// ProcessPhotoInput holds the data needed to process an uploaded evidence photo.
type ProcessPhotoInput struct {
	BookingID string                 `json:"bookingId" validate:"required"`
	S3Key     string                 `json:"s3Key" validate:"required"`
	S3Bucket  string                 `json:"s3Bucket" validate:"required"`
	Position  evidence.PhotoPosition `json:"position" validate:"required"`
	MimeType  string                 `json:"mimeType" validate:"required"`
	FileHash  string                 `json:"fileHash" validate:"required"`
}

// ProcessPhotoUseCase orchestrates photo validation, persistence, and event publishing.
type ProcessPhotoUseCase struct {
	validator ImageValidator
	media     MediaItemAdder
	events    EventPublisher
	idGen     func() string
	now       func() time.Time
}

// NewProcessPhotoUseCase constructs a ProcessPhotoUseCase with its dependencies.
func NewProcessPhotoUseCase(validator ImageValidator, media MediaItemAdder, events EventPublisher) *ProcessPhotoUseCase {
	return &ProcessPhotoUseCase{
		validator: validator,
		media:     media,
		events:    events,
		idGen:     generateMediaID,
		now:       func() time.Time { return time.Now().UTC() },
	}
}

// Execute validates the photo via Rekognition, creates a MediaItem, saves it, and publishes an event.
func (uc *ProcessPhotoUseCase) Execute(ctx context.Context, input *ProcessPhotoInput) (*evidence.MediaItem, error) {
	if err := validateProcessInput(input); err != nil {
		return nil, err
	}

	result, err := uc.validator.ValidateVehiclePhoto(ctx, input.S3Bucket, input.S3Key)
	if err != nil {
		return nil, domainerrors.NewExternalServiceError("Rekognition", err)
	}
	if !result.IsValid {
		return nil, domainerrors.NewEvidenceValidationFailedError(result.Reason)
	}

	mediaItem := &evidence.MediaItem{
		MediaID:  uc.idGen(),
		S3Key:    input.S3Key,
		Position: input.Position,
		MimeType: input.MimeType,
		Integrity: evidence.HashIntegrity{
			Algorithm: "SHA-256",
			Hash:      input.FileHash,
		},
		CapturedAt: uc.now(),
	}

	if err := uc.media.AddMediaItem(ctx, input.BookingID, mediaItem); err != nil {
		return nil, err
	}

	_ = uc.events.Publish(ctx, eventSourceEvidence, eventEvidenceValidated, map[string]any{
		"bookingId": input.BookingID,
		"mediaId":   mediaItem.MediaID,
		"s3Key":     input.S3Key,
		"position":  input.Position,
		"labels":    result.Labels,
	}, nil)

	return mediaItem, nil
}

// validateProcessInput performs basic validation on the process photo input.
func validateProcessInput(input *ProcessPhotoInput) error {
	if input.BookingID == "" {
		return domainerrors.NewValidationError("bookingId is required")
	}
	if input.S3Key == "" {
		return domainerrors.NewValidationError("s3Key is required")
	}
	if input.S3Bucket == "" {
		return domainerrors.NewValidationError("s3Bucket is required")
	}
	if input.Position == "" {
		return domainerrors.NewValidationError("position is required")
	}
	if input.MimeType == "" {
		return domainerrors.NewValidationError("mimeType is required")
	}
	if input.FileHash == "" {
		return domainerrors.NewValidationError("fileHash is required")
	}
	return nil
}

// generateMediaID produces a random media ID.
func generateMediaID() string {
	b := make([]byte, 16)
	_, _ = rand.Read(b)
	return fmt.Sprintf("media-%x", b)
}
