package booking

import (
	"time"

	"github.com/davidlramirez95/towcommand/internal/domain/user"
)

// ServiceType represents the kind of roadside service requested.
type ServiceType string

const (
	ServiceTypeFlatbedTow       ServiceType = "FLATBED_TOW"
	ServiceTypeWheelLift        ServiceType = "WHEEL_LIFT"
	ServiceTypeJumpstart        ServiceType = "JUMPSTART"
	ServiceTypeTireChange       ServiceType = "TIRE_CHANGE"
	ServiceTypeFuelDelivery     ServiceType = "FUEL_DELIVERY"
	ServiceTypeLockout          ServiceType = "LOCKOUT"
	ServiceTypeAccidentRecovery ServiceType = "ACCIDENT_RECOVERY"
)

// GeoLocation represents a geographic coordinate with an optional address.
type GeoLocation struct {
	Lat     float64 `json:"lat" validate:"required,latitude"`
	Lng     float64 `json:"lng" validate:"required,longitude"`
	Address string  `json:"address,omitempty"`
}

// PriceBreakdown itemizes the cost of a booking in centavos.
type PriceBreakdown struct {
	Base            int64   `json:"base"`
	Distance        int64   `json:"distance"`
	Weight          int64   `json:"weight"`
	TimeSurcharge   int64   `json:"timeSurcharge"`
	SurgePricing    int64   `json:"surgePricing"`
	Total           int64   `json:"total"`
	Currency        string  `json:"currency" validate:"required,eq=PHP"`
	SurgeMultiplier float64 `json:"surgeMultiplier,omitempty"`
}

// Booking represents a tow/roadside assistance service request.
type Booking struct {
	BookingID          string           `json:"bookingId" validate:"required"`
	CustomerID         string           `json:"customerId" validate:"required"`
	ProviderID         string           `json:"providerId,omitempty"`
	VehicleID          string           `json:"vehicleId" validate:"required"`
	ServiceType        ServiceType      `json:"serviceType" validate:"required"`
	Status             BookingStatus    `json:"status" validate:"required"`
	PickupLocation     GeoLocation      `json:"pickupLocation" validate:"required"`
	DropoffLocation    GeoLocation      `json:"dropoffLocation" validate:"required"`
	WeightClass        user.WeightClass `json:"weightClass" validate:"required"`
	Price              PriceBreakdown   `json:"price"`
	EstimateID         string           `json:"estimateId" validate:"required"`
	Notes              string           `json:"notes,omitempty"`
	CancellationReason string           `json:"cancellationReason,omitempty"`
	CancellationFee    int64            `json:"cancellationFee,omitempty"`
	MatchedAt          *time.Time       `json:"matchedAt,omitempty"`
	CompletedAt        *time.Time       `json:"completedAt,omitempty"`
	CreatedAt          time.Time        `json:"createdAt"`
	UpdatedAt          time.Time        `json:"updatedAt"`
}

// BookingEstimate represents a price quote for a potential booking.
type BookingEstimate struct {
	EstimateID          string           `json:"estimateId" validate:"required"`
	PickupLocation      GeoLocation      `json:"pickupLocation" validate:"required"`
	DropoffLocation     GeoLocation      `json:"dropoffLocation" validate:"required"`
	ServiceType         ServiceType      `json:"serviceType" validate:"required"`
	WeightClass         user.WeightClass `json:"weightClass" validate:"required"`
	DistanceKm          float64          `json:"distanceKm" validate:"required,gt=0"`
	Price               PriceBreakdown   `json:"price"`
	AvailableProviders  int              `json:"availableProviders"`
	EstimatedETAMinutes int              `json:"estimatedEtaMinutes"`
	ExpiresAt           time.Time        `json:"expiresAt"`
	CreatedAt           time.Time        `json:"createdAt"`
}
