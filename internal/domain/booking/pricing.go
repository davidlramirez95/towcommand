package booking

import (
	"math"
	"time"

	"github.com/davidlramirez95/towcommand/internal/domain/user"
)

// manilaLoc is the Asia/Manila timezone (PHT = UTC+8).
var manilaLoc = time.FixedZone("Asia/Manila", 8*60*60)

// pricingTier defines the base and per-km rates in centavos for a weight class.
type pricingTier struct {
	baseRate  int64
	perKmRate int64
}

// MMDA Regulation 24-004 towing rates in centavos.
var mmdaRates = map[user.WeightClass]pricingTier{
	user.WeightClassMotorcycle: {baseRate: 80_000, perKmRate: 10_000},
	user.WeightClassLight:      {baseRate: 180_000, perKmRate: 20_000},
	user.WeightClassMedium:     {baseRate: 300_000, perKmRate: 25_000},
	user.WeightClassHeavy:      {baseRate: 500_000, perKmRate: 35_000},
	user.WeightClassSuperHeavy: {baseRate: 800_000, perKmRate: 50_000},
}

// Roadside service flat rates in centavos.
var roadsideRates = map[ServiceType]int64{
	ServiceTypeJumpstart:    50_000,
	ServiceTypeTireChange:   60_000,
	ServiceTypeFuelDelivery: 40_000,
	ServiceTypeLockout:      70_000,
}

const (
	nightSurchargeRate   = 0.30 // +30 % of base between 22:00-06:00 PHT
	weekendSurchargeRate = 0.15 // +15 % of base on Sat/Sun

	accidentSurchargeCentavos int64 = 50_000  // ₱500
	superHeavyWeightCentavos  int64 = 100_000 // ₱1,000
	maxSurge                        = 2.5
	minSurge                        = 1.0

	waitingFreeMinutes        = 30
	waitingBlockMinutes       = 15
	waitingFeePerBlock  int64 = 10_000 // ₱100 per 15-min block

	platformFeeCentavos int64 = 5_000 // ₱50 flat booking fee
)

// IsNightTime reports whether t falls in the night surcharge window
// (22:00–05:59 PHT).
func IsNightTime(t time.Time) bool {
	hour := t.In(manilaLoc).Hour()
	return hour >= 22 || hour < 6
}

// IsWeekend reports whether t falls on Saturday or Sunday in PHT.
func IsWeekend(t time.Time) bool {
	day := t.In(manilaLoc).Weekday()
	return day == time.Saturday || day == time.Sunday
}

// CalculateEstimate computes a PriceBreakdown for a potential booking.
// All monetary amounts are in centavos (int64).
func CalculateEstimate(
	serviceType ServiceType,
	weightClass user.WeightClass,
	distanceKm float64,
	pickupTime time.Time,
	surgeMultiplier float64,
	isAccident bool,
) PriceBreakdown {
	surge := clampSurge(surgeMultiplier)

	var base, distance int64

	if flatRate, ok := roadsideRates[serviceType]; ok {
		base = flatRate
		distance = 0
	} else {
		tier, ok := mmdaRates[weightClass]
		if !ok {
			tier = mmdaRates[user.WeightClassLight]
		}
		base = tier.baseRate
		distance = int64(math.Ceil(distanceKm)) * tier.perKmRate
	}

	var weight int64
	if weightClass == user.WeightClassSuperHeavy {
		weight = superHeavyWeightCentavos
	}

	var nightCharge int64
	if IsNightTime(pickupTime) {
		nightCharge = roundCentavos(float64(base) * nightSurchargeRate)
	}

	var weekendCharge int64
	if IsWeekend(pickupTime) {
		weekendCharge = roundCentavos(float64(base) * weekendSurchargeRate)
	}

	var accidentCharge int64
	if isAccident {
		accidentCharge = accidentSurchargeCentavos
	}

	timeSurcharge := nightCharge + weekendCharge + accidentCharge

	subtotal := base + distance + weight + timeSurcharge
	surgePricing := roundCentavos(float64(subtotal) * (surge - 1))
	total := subtotal + surgePricing

	pb := PriceBreakdown{
		Base:          base,
		Distance:      distance,
		Weight:        weight,
		TimeSurcharge: timeSurcharge,
		SurgePricing:  surgePricing,
		Total:         total,
		Currency:      "PHP",
	}
	if surge > 1 {
		pb.SurgeMultiplier = surge
	}
	return pb
}

// FinalPrice represents the total price after a completed booking,
// including post-service charges. All amounts in centavos.
type FinalPrice struct {
	Estimate    PriceBreakdown `json:"estimate"`
	WaitingFee  int64          `json:"waitingFee"`
	TollFees    int64          `json:"tollFees"`
	PlatformFee int64          `json:"platformFee"`
	Total       int64          `json:"total"`
	Currency    string         `json:"currency" validate:"required,eq=PHP"`
}

// CalculateFinalPrice computes the total price including waiting time,
// toll fees, and platform fee on top of the original estimate.
func CalculateFinalPrice(estimate PriceBreakdown, waitingMinutes int, tollFees int64) FinalPrice {
	waiting := calculateWaitingFee(waitingMinutes)
	total := estimate.Total + waiting + tollFees + platformFeeCentavos

	return FinalPrice{
		Estimate:    estimate,
		WaitingFee:  waiting,
		TollFees:    tollFees,
		PlatformFee: platformFeeCentavos,
		Total:       total,
		Currency:    "PHP",
	}
}

// CalculateCancellationFee returns the cancellation fee in centavos,
// tiered by booking status:
//   - PENDING:   ₱0
//   - MATCHED:   ₱100
//   - EN_ROUTE:  ₱100 + (distance × ₱30), capped at ₱500
//   - ARRIVED:   ₱500 + min(distance × ₱30, ₱300)
func CalculateCancellationFee(status BookingStatus, distanceKmTravelled float64) int64 {
	switch status {
	case BookingStatusPending:
		return 0
	case BookingStatusMatched:
		return 10_000 // ₱100
	case BookingStatusEnRoute:
		fee := 10_000 + int64(distanceKmTravelled*3_000)
		if fee > 50_000 {
			return 50_000
		}
		return fee
	case BookingStatusArrived:
		distFee := int64(distanceKmTravelled * 3_000)
		if distFee > 30_000 {
			distFee = 30_000
		}
		return 50_000 + distFee
	default:
		return 0
	}
}

// GetSurgeMultiplier calculates the surge multiplier from demand/supply
// counts. Returns a value clamped between 1.0 and 2.5.
func GetSurgeMultiplier(demandCount, supplyCount int) float64 {
	if supplyCount <= 0 {
		return maxSurge
	}
	if demandCount <= 0 {
		return minSurge
	}
	ratio := float64(demandCount) / float64(supplyCount)
	// Linear mapping: ratio 1.0 → 1.0x, ratio 5.0 → 2.5x
	surge := 1.0 + (ratio-1.0)*0.375
	return clampSurge(surge)
}

func clampSurge(surge float64) float64 {
	if surge < minSurge {
		return minSurge
	}
	if surge > maxSurge {
		return maxSurge
	}
	return surge
}

func calculateWaitingFee(waitingMinutes int) int64 {
	if waitingMinutes <= waitingFreeMinutes {
		return 0
	}
	billable := waitingMinutes - waitingFreeMinutes
	blocks := int64(math.Ceil(float64(billable) / float64(waitingBlockMinutes)))
	return blocks * waitingFeePerBlock
}

func roundCentavos(f float64) int64 {
	return int64(math.Round(f))
}
