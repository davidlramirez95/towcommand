package evidenceuc

import (
	"context"
	"fmt"
	"time"

	domainerrors "github.com/davidlramirez95/towcommand/internal/domain/errors"
	"github.com/davidlramirez95/towcommand/internal/domain/evidence"
)

const (
	uploadURLExpiry = 900 * time.Second // 15 minutes
)

// GenerateUploadURLInput holds the data needed to generate a presigned upload URL.
type GenerateUploadURLInput struct {
	BookingID   string                 `json:"bookingId" validate:"required"`
	Phase       string                 `json:"phase" validate:"required,oneof=pickup dropoff"`
	Position    evidence.PhotoPosition `json:"position" validate:"required"`
	ContentType string                 `json:"contentType" validate:"required"`
}

// GenerateUploadURLOutput is the result of generating a presigned upload URL.
type GenerateUploadURLOutput struct {
	UploadURL string `json:"uploadUrl"`
	S3Key     string `json:"s3Key"`
	ExpiresIn int    `json:"expiresIn"`
}

// GenerateUploadURLUseCase orchestrates presigned URL generation for evidence uploads.
type GenerateUploadURLUseCase struct {
	bookings  BookingFinder
	presigner PresignedURLGenerator
	now       func() time.Time
}

// NewGenerateUploadURLUseCase constructs a GenerateUploadURLUseCase with its dependencies.
func NewGenerateUploadURLUseCase(bookings BookingFinder, presigner PresignedURLGenerator) *GenerateUploadURLUseCase {
	return &GenerateUploadURLUseCase{
		bookings:  bookings,
		presigner: presigner,
		now:       func() time.Time { return time.Now().UTC() },
	}
}

// Execute validates the booking exists and returns a presigned S3 upload URL.
func (uc *GenerateUploadURLUseCase) Execute(ctx context.Context, input *GenerateUploadURLInput) (*GenerateUploadURLOutput, error) {
	if err := validateUploadInput(input); err != nil {
		return nil, err
	}

	b, err := uc.bookings.FindByID(ctx, input.BookingID)
	if err != nil {
		return nil, err
	}
	if b == nil {
		return nil, domainerrors.NewNotFoundError("Booking", input.BookingID)
	}

	timestamp := uc.now().Unix()
	s3Key := fmt.Sprintf("evidence/%s/%s/%s_%d.jpg", input.BookingID, input.Phase, input.Position, timestamp)

	url, err := uc.presigner.GenerateUploadURL(ctx, s3Key, input.ContentType, uploadURLExpiry)
	if err != nil {
		return nil, domainerrors.NewExternalServiceError("S3", err)
	}

	return &GenerateUploadURLOutput{
		UploadURL: url,
		S3Key:     s3Key,
		ExpiresIn: int(uploadURLExpiry.Seconds()),
	}, nil
}

// validateUploadInput performs basic validation on the upload input.
func validateUploadInput(input *GenerateUploadURLInput) error {
	if input.BookingID == "" {
		return domainerrors.NewValidationError("bookingId is required")
	}
	if input.Phase != "pickup" && input.Phase != "dropoff" {
		return domainerrors.NewValidationError("phase must be pickup or dropoff")
	}
	if input.Position == "" {
		return domainerrors.NewValidationError("position is required")
	}
	if input.ContentType == "" {
		return domainerrors.NewValidationError("contentType is required")
	}
	return nil
}
