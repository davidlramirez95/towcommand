package safetyuc

import (
	"context"
	"crypto/rand"
	"fmt"
	"math"
	"time"

	domainerrors "github.com/davidlramirez95/towcommand/internal/domain/errors"
	"github.com/davidlramirez95/towcommand/internal/domain/event"
	"github.com/davidlramirez95/towcommand/internal/domain/safety"
)

// phtLocation is the Philippine Standard Time timezone (UTC+8).
var phtLocation = time.FixedZone("PHT", 8*60*60)

// TriggerSOSInput holds the data needed to trigger an SOS alert.
type TriggerSOSInput struct {
	BookingID   string             `json:"bookingId" validate:"required"`
	TriggeredBy string             `json:"triggeredBy" validate:"required"`
	TriggerType safety.TriggerType `json:"triggerType" validate:"required"`
	Lat         float64            `json:"lat" validate:"required"`
	Lng         float64            `json:"lng" validate:"required"`
}

// TriggerSOSUseCase orchestrates the creation of an SOS alert.
type TriggerSOSUseCase struct {
	bookings  BookingFinder
	providers ProviderFinder
	sos       SOSSaver
	events    EventPublisher
	now       func() time.Time
	idGen     func() string
}

// NewTriggerSOSUseCase constructs a TriggerSOSUseCase with its dependencies.
func NewTriggerSOSUseCase(
	bookings BookingFinder,
	providers ProviderFinder,
	sos SOSSaver,
	events EventPublisher,
) *TriggerSOSUseCase {
	return &TriggerSOSUseCase{
		bookings:  bookings,
		providers: providers,
		sos:       sos,
		events:    events,
		now:       func() time.Time { return time.Now().UTC() },
		idGen:     generateAlertID,
	}
}

// Execute triggers an SOS alert: looks up the booking and provider, computes
// risk, persists the alert, and publishes an SOSTriggered event.
func (uc *TriggerSOSUseCase) Execute(ctx context.Context, input *TriggerSOSInput) (*safety.SOSAlert, error) {
	if err := validateTriggerInput(input); err != nil {
		return nil, err
	}

	b, err := uc.bookings.FindByID(ctx, input.BookingID)
	if err != nil {
		return nil, fmt.Errorf("finding booking %s: %w", input.BookingID, err)
	}
	if b == nil {
		return nil, domainerrors.NewNotFoundError("booking", input.BookingID)
	}

	// Determine provider trust tier for risk scoring.
	trustTier := "basic"
	if b.ProviderID != "" {
		p, err := uc.providers.FindByID(ctx, b.ProviderID)
		if err != nil {
			return nil, fmt.Errorf("finding provider %s: %w", b.ProviderID, err)
		}
		if p != nil {
			trustTier = string(p.TrustTier)
		}
	}

	now := uc.now()

	// Compute risk factors.
	factors := safety.RiskFactors{
		IsNightTime:       isNightTime(now),
		ProviderTrustTier: trustTier,
		DistanceKm: haversine(
			input.Lat, input.Lng,
			b.DropoffLocation.Lat, b.DropoffLocation.Lng,
		),
		PriorSOSCount:  0,     // Hardcoded for MVP
		IsHighRiskZone: false, // Hardcoded for MVP
	}

	score := safety.ComputeRiskScore(factors)

	riskScore := safety.RiskScore{
		Score:      score,
		AssessedAt: now,
	}
	riskScore.Level = riskScore.ComputeLevel()
	riskScore.Factors = buildFactorList(factors)

	alert := &safety.SOSAlert{
		AlertID:     uc.idGen(),
		BookingID:   input.BookingID,
		TriggeredBy: input.TriggeredBy,
		TriggerType: input.TriggerType,
		Lat:         input.Lat,
		Lng:         input.Lng,
		Risk:        riskScore,
		Resolved:    false,
		Timestamp:   now,
	}

	if err := uc.sos.Save(ctx, alert); err != nil {
		return nil, fmt.Errorf("saving SOS alert: %w", err)
	}

	_ = uc.events.Publish(ctx, event.SourceSOS, event.SOSTriggered, alert, &Actor{
		UserID:   input.TriggeredBy,
		UserType: "customer",
	})

	return alert, nil
}

// validateTriggerInput performs basic validation on the trigger input.
func validateTriggerInput(input *TriggerSOSInput) error {
	if input.BookingID == "" {
		return domainerrors.NewValidationError("bookingId is required")
	}
	if input.TriggeredBy == "" {
		return domainerrors.NewValidationError("triggeredBy is required")
	}
	if input.TriggerType == "" {
		return domainerrors.NewValidationError("triggerType is required")
	}
	return nil
}

// isNightTime returns true if the given UTC time falls in the 20:00-05:00 PHT window.
func isNightTime(t time.Time) bool {
	pht := t.In(phtLocation)
	hour := pht.Hour()
	return hour >= 20 || hour < 5
}

// haversine computes the great-circle distance between two lat/lng pairs in kilometres.
func haversine(lat1, lng1, lat2, lng2 float64) float64 {
	const earthRadiusKm = 6371.0
	dLat := degreesToRadians(lat2 - lat1)
	dLng := degreesToRadians(lng2 - lng1)
	a := math.Sin(dLat/2)*math.Sin(dLat/2) +
		math.Cos(degreesToRadians(lat1))*math.Cos(degreesToRadians(lat2))*
			math.Sin(dLng/2)*math.Sin(dLng/2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
	return earthRadiusKm * c
}

// degreesToRadians converts degrees to radians.
func degreesToRadians(deg float64) float64 {
	return deg * math.Pi / 180.0
}

// buildFactorList returns human-readable descriptions of active risk factors.
func buildFactorList(f safety.RiskFactors) []string {
	var factors []string
	if f.IsNightTime {
		factors = append(factors, "night_time")
	}
	if f.ProviderTrustTier == "basic" {
		factors = append(factors, "basic_tier_provider")
	}
	if f.DistanceKm > 30 {
		factors = append(factors, "long_distance")
	} else if f.DistanceKm > 20 {
		factors = append(factors, "medium_distance")
	}
	if f.PriorSOSCount > 0 {
		factors = append(factors, "prior_sos_history")
	}
	if f.IsHighRiskZone {
		factors = append(factors, "high_risk_zone")
	}
	return factors
}

// generateAlertID creates a unique alert ID in the format SOS-<year>-<random hex>.
func generateAlertID() string {
	b := make([]byte, 8)
	_, _ = rand.Read(b)
	return fmt.Sprintf("SOS-%d-%x", time.Now().Year(), b)
}
