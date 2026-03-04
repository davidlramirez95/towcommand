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
	eventConditionReportCreated = "ConditionReportCreated"
)

// CreateConditionReportInput holds the data needed to create a condition report.
type CreateConditionReportInput struct {
	BookingID  string `json:"bookingId" validate:"required"`
	ProviderID string `json:"providerId" validate:"required"`
	Phase      string `json:"phase" validate:"required,oneof=pickup dropoff"`
	Notes      string `json:"notes"`
}

// CreateConditionReportUseCase orchestrates condition report creation.
type CreateConditionReportUseCase struct {
	bookings BookingFinder
	evidence EvidenceSaver
	events   EventPublisher
	idGen    func() string
	now      func() time.Time
}

// NewCreateConditionReportUseCase constructs a CreateConditionReportUseCase with its dependencies.
func NewCreateConditionReportUseCase(bookings BookingFinder, evidence EvidenceSaver, events EventPublisher) *CreateConditionReportUseCase {
	return &CreateConditionReportUseCase{
		bookings: bookings,
		evidence: evidence,
		events:   events,
		idGen:    generateReportID,
		now:      func() time.Time { return time.Now().UTC() },
	}
}

// Execute validates the booking exists, creates a condition report, saves it, and publishes an event.
func (uc *CreateConditionReportUseCase) Execute(ctx context.Context, input *CreateConditionReportInput) (*evidence.ConditionReport, error) {
	if err := validateReportInput(input); err != nil {
		return nil, err
	}

	b, err := uc.bookings.FindByID(ctx, input.BookingID)
	if err != nil {
		return nil, err
	}
	if b == nil {
		return nil, domainerrors.NewNotFoundError("Booking", input.BookingID)
	}

	report := &evidence.ConditionReport{
		ReportID:   uc.idGen(),
		BookingID:  input.BookingID,
		ProviderID: input.ProviderID,
		Phase:      input.Phase,
		Media:      []evidence.MediaItem{},
		Notes:      input.Notes,
		CreatedAt:  uc.now(),
	}

	if err := uc.evidence.Save(ctx, report); err != nil {
		return nil, err
	}

	_ = uc.events.Publish(ctx, eventSourceEvidence, eventConditionReportCreated, map[string]any{
		"reportId":   report.ReportID,
		"bookingId":  input.BookingID,
		"providerId": input.ProviderID,
		"phase":      input.Phase,
	}, &Actor{UserID: input.ProviderID, UserType: "provider"})

	return report, nil
}

// validateReportInput performs basic validation on the condition report input.
func validateReportInput(input *CreateConditionReportInput) error {
	if input.BookingID == "" {
		return domainerrors.NewValidationError("bookingId is required")
	}
	if input.ProviderID == "" {
		return domainerrors.NewValidationError("providerId is required")
	}
	if input.Phase != "pickup" && input.Phase != "dropoff" {
		return domainerrors.NewValidationError("phase must be pickup or dropoff")
	}
	return nil
}

// generateReportID produces a random report ID.
func generateReportID() string {
	b := make([]byte, 16)
	_, _ = rand.Read(b)
	return fmt.Sprintf("report-%x", b)
}

// --- CheckCompletenessUseCase ---

// CompletenessOutput holds the result of a completeness check on evidence for a booking.
type CompletenessOutput struct {
	IsComplete       bool                     `json:"isComplete"`
	TotalPhotos      int                      `json:"totalPhotos"`
	RequiredPhotos   int                      `json:"requiredPhotos"`
	MissingPositions []evidence.PhotoPosition `json:"missingPositions"`
}

// CheckCompletenessUseCase checks whether all 8 required photo positions are covered.
type CheckCompletenessUseCase struct {
	evidence EvidenceByBookingLister
}

// NewCheckCompletenessUseCase constructs a CheckCompletenessUseCase with its dependencies.
func NewCheckCompletenessUseCase(evidence EvidenceByBookingLister) *CheckCompletenessUseCase {
	return &CheckCompletenessUseCase{evidence: evidence}
}

// Execute retrieves all condition reports for a booking and determines coverage.
func (uc *CheckCompletenessUseCase) Execute(ctx context.Context, bookingID string) (*CompletenessOutput, error) {
	if bookingID == "" {
		return nil, domainerrors.NewValidationError("bookingId is required")
	}

	reports, err := uc.evidence.FindByBooking(ctx, bookingID)
	if err != nil {
		return nil, err
	}

	covered := make(map[evidence.PhotoPosition]bool, len(evidence.AllPhotoPositions))
	totalPhotos := 0

	for i := range reports {
		for _, m := range reports[i].Media {
			covered[m.Position] = true
			totalPhotos++
		}
	}

	var missing []evidence.PhotoPosition
	for _, pos := range evidence.AllPhotoPositions {
		if !covered[pos] {
			missing = append(missing, pos)
		}
	}

	return &CompletenessOutput{
		IsComplete:       len(missing) == 0,
		TotalPhotos:      totalPhotos,
		RequiredPhotos:   len(evidence.AllPhotoPositions),
		MissingPositions: missing,
	}, nil
}
