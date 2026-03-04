package bookinguc

import (
	"context"
	"crypto/rand"
	"fmt"
	"math"
	"time"

	"github.com/davidlramirez95/towcommand/internal/domain/booking"
	"github.com/davidlramirez95/towcommand/internal/domain/user"
)

const (
	eventSourceBooking  = "towcommand.booking"
	eventBookingCreated = "BookingCreated"
)

// CreateBookingInput holds the data needed to create a booking.
type CreateBookingInput struct {
	CustomerID      string
	VehicleID       string
	ServiceType     booking.ServiceType
	PickupLocation  booking.GeoLocation
	DropoffLocation booking.GeoLocation
	EstimateID      string
	Notes           string
}

// CreateBookingUseCase orchestrates booking creation.
type CreateBookingUseCase struct {
	repo   BookingSaver
	events EventPublisher
	idGen  func() string
	now    func() time.Time
}

// NewCreateBookingUseCase constructs a CreateBookingUseCase with its dependencies.
func NewCreateBookingUseCase(repo BookingSaver, events EventPublisher) *CreateBookingUseCase {
	return &CreateBookingUseCase{
		repo:   repo,
		events: events,
		idGen:  generateBookingID,
		now:    func() time.Time { return time.Now().UTC() },
	}
}

// Execute creates a new booking, saves it, and publishes a BookingCreated event.
func (uc *CreateBookingUseCase) Execute(ctx context.Context, input *CreateBookingInput) (*booking.Booking, error) {
	now := uc.now()
	bookingID := uc.idGen()

	distanceKm := haversineDistance(
		input.PickupLocation.Lat, input.PickupLocation.Lng,
		input.DropoffLocation.Lat, input.DropoffLocation.Lng,
	)

	weightClass := user.WeightClassLight
	price := calculatePrice(input.ServiceType, weightClass, distanceKm, isNightTimePHT(now))

	b := &booking.Booking{
		BookingID:       bookingID,
		CustomerID:      input.CustomerID,
		VehicleID:       input.VehicleID,
		ServiceType:     input.ServiceType,
		Status:          booking.BookingStatusPending,
		PickupLocation:  input.PickupLocation,
		DropoffLocation: input.DropoffLocation,
		WeightClass:     weightClass,
		Price:           price,
		EstimateID:      input.EstimateID,
		Notes:           input.Notes,
		CreatedAt:       now,
		UpdatedAt:       now,
	}

	if err := uc.repo.Save(ctx, b); err != nil {
		return nil, err
	}

	_ = uc.events.Publish(ctx, eventSourceBooking, eventBookingCreated, map[string]any{
		"bookingId":       bookingID,
		"customerId":      input.CustomerID,
		"serviceType":     input.ServiceType,
		"pickupLocation":  input.PickupLocation,
		"dropoffLocation": input.DropoffLocation,
		"price":           price,
	}, &Actor{UserID: input.CustomerID, UserType: string(user.UserTypeCustomer)})

	return b, nil
}

// generateBookingID produces a booking ID in the format TC-<year>-<random>.
func generateBookingID() string {
	b := make([]byte, 10)
	_, _ = rand.Read(b)
	return fmt.Sprintf("TC-%d-%X", time.Now().Year(), b)
}

// haversineDistance calculates the great-circle distance in km between two points.
func haversineDistance(lat1, lng1, lat2, lng2 float64) float64 {
	const earthRadiusKm = 6371.0
	dLat := (lat2 - lat1) * math.Pi / 180
	dLng := (lng2 - lng1) * math.Pi / 180
	lat1Rad := lat1 * math.Pi / 180
	lat2Rad := lat2 * math.Pi / 180

	a := math.Sin(dLat/2)*math.Sin(dLat/2) +
		math.Cos(lat1Rad)*math.Cos(lat2Rad)*math.Sin(dLng/2)*math.Sin(dLng/2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
	return earthRadiusKm * c
}

// isNightTimePHT returns true if the given UTC time falls within the
// night surcharge window (22:00–06:00 PHT, UTC+8).
func isNightTimePHT(t time.Time) bool {
	hour := t.Add(8 * time.Hour).Hour()
	return hour >= 22 || hour < 6
}

// baseRateCentavos returns the base rate in centavos for a service type.
func baseRateCentavos(st booking.ServiceType) int64 {
	switch st {
	case booking.ServiceTypeFlatbedTow:
		return 150_000
	case booking.ServiceTypeWheelLift:
		return 120_000
	case booking.ServiceTypeJumpstart:
		return 50_000
	case booking.ServiceTypeTireChange:
		return 60_000
	case booking.ServiceTypeFuelDelivery:
		return 40_000
	case booking.ServiceTypeLockout:
		return 55_000
	case booking.ServiceTypeAccidentRecovery:
		return 250_000
	default:
		return 100_000
	}
}

// weightMultiplier returns the weight-based price multiplier.
func weightMultiplier(wc user.WeightClass) float64 {
	switch wc {
	case user.WeightClassMotorcycle:
		return 0.5
	case user.WeightClassLight:
		return 1.0
	case user.WeightClassMedium:
		return 1.3
	case user.WeightClassHeavy:
		return 1.8
	case user.WeightClassSuperHeavy:
		return 2.5
	default:
		return 1.0
	}
}

// calculatePrice builds a PriceBreakdown for the given parameters.
// Rates follow MMDA Regulation 24-004 guidelines.
func calculatePrice(st booking.ServiceType, wc user.WeightClass, distanceKm float64, nightTime bool) booking.PriceBreakdown {
	base := baseRateCentavos(st)
	distanceFee := int64(distanceKm * 5_000) // ₱50/km
	weightFee := int64(float64(base) * (weightMultiplier(wc) - 1))

	subtotal := base + distanceFee + weightFee
	var timeSurcharge int64
	if nightTime {
		timeSurcharge = int64(float64(subtotal) * 0.30)
	}
	total := subtotal + timeSurcharge

	return booking.PriceBreakdown{
		Base:          base,
		Distance:      distanceFee,
		Weight:        weightFee,
		TimeSurcharge: timeSurcharge,
		SurgePricing:  0,
		Total:         total,
		Currency:      "PHP",
	}
}
