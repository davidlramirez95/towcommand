package provider

import (
	"math"
	"sort"

	"github.com/davidlramirez95/towcommand/internal/domain/booking"
	"github.com/davidlramirez95/towcommand/internal/domain/user"
)

// MatchCandidate holds provider data needed for scoring.
type MatchCandidate struct {
	ProviderID          string
	TrustTier           user.TrustTier
	AcceptanceRate      float64
	TruckType           TruckType
	MaxWeightCapacityKg int
	ActiveJobCount      int
	DistanceKm          float64
	IsOnline            bool
}

// MatchFactors holds the individual scoring components.
type MatchFactors struct {
	Distance       float64
	TrustTier      float64
	AcceptanceRate float64
	VehicleCompat  float64
	CurrentLoad    float64
}

// MatchScore holds the weighted total score and breakdown for a candidate.
type MatchScore struct {
	ProviderID string
	TotalScore float64
	Factors    MatchFactors
	DistanceKm float64
}

const minMatchThreshold = 0.30

// Weight configuration per TRS spec.
var (
	defaultWeights = weights{Distance: 0.40, TrustTier: 0.25, AcceptanceRate: 0.15, VehicleCompat: 0.10, CurrentLoad: 0.10}
	surgeWeights   = weights{Distance: 0.50, TrustTier: 0.20, AcceptanceRate: 0.12, VehicleCompat: 0.08, CurrentLoad: 0.10}
)

type weights struct {
	Distance       float64
	TrustTier      float64
	AcceptanceRate float64
	VehicleCompat  float64
	CurrentLoad    float64
}

// trustTierScores maps trust tier to a 0-100 score.
var trustTierScores = map[user.TrustTier]float64{
	user.TrustTierBasic:      20,
	user.TrustTierVerified:   40,
	user.TrustTierSukiSilver: 60,
	user.TrustTierSukiGold:   80,
	user.TrustTierSukiElite:  100,
}

// serviceTypeCapabilities maps service type to capable truck types.
var serviceTypeCapabilities = map[booking.ServiceType][]TruckType{
	booking.ServiceTypeFlatbedTow:       {TruckTypeFlatbed},
	booking.ServiceTypeWheelLift:        {TruckTypeWheelLift, TruckTypeFlatbed},
	booking.ServiceTypeAccidentRecovery: {TruckTypeBoom, TruckTypeFlatbed},
	booking.ServiceTypeJumpstart:        {TruckTypeFlatbed, TruckTypeWheelLift, TruckTypeBoom, TruckTypeMotorcycleCarrier},
	booking.ServiceTypeTireChange:       {TruckTypeFlatbed, TruckTypeWheelLift, TruckTypeBoom, TruckTypeMotorcycleCarrier},
	booking.ServiceTypeFuelDelivery:     {TruckTypeFlatbed, TruckTypeWheelLift, TruckTypeBoom, TruckTypeMotorcycleCarrier},
	booking.ServiceTypeLockout:          {TruckTypeFlatbed, TruckTypeWheelLift, TruckTypeBoom, TruckTypeMotorcycleCarrier},
}

// weightClassKg maps weight classes to approximate kg values for capacity comparison.
var weightClassKg = map[user.WeightClass]int{
	user.WeightClassMotorcycle: 200,
	user.WeightClassLight:      1500,
	user.WeightClassMedium:     2500,
	user.WeightClassHeavy:      5000,
	user.WeightClassSuperHeavy: 10000,
}

// WeightClassToKg returns the approximate kg value for a weight class.
func WeightClassToKg(wc user.WeightClass) int {
	if kg, ok := weightClassKg[wc]; ok {
		return kg
	}
	return 1500 // default to light
}

// ScoreProvider calculates a match score for a single candidate.
func ScoreProvider(candidate *MatchCandidate, serviceType booking.ServiceType, weightKg int, surgeMode bool) MatchScore {
	w := defaultWeights
	if surgeMode {
		w = surgeWeights
	}

	factors := MatchFactors{
		Distance:       scoreDistance(candidate.DistanceKm),
		TrustTier:      scoreTrustTier(candidate.TrustTier),
		AcceptanceRate: scoreAcceptanceRate(candidate.AcceptanceRate),
		VehicleCompat:  scoreVehicleCompat(candidate.TruckType, candidate.MaxWeightCapacityKg, serviceType, weightKg),
		CurrentLoad:    scoreCurrentLoad(candidate.ActiveJobCount),
	}

	total := factors.Distance*w.Distance +
		factors.TrustTier*w.TrustTier +
		factors.AcceptanceRate*w.AcceptanceRate +
		factors.VehicleCompat*w.VehicleCompat +
		factors.CurrentLoad*w.CurrentLoad

	// Normalize to 0-1 range (each factor is 0-100, weights sum to 1.0).
	total /= 100.0

	return MatchScore{
		ProviderID: candidate.ProviderID,
		TotalScore: total,
		Factors:    factors,
		DistanceKm: candidate.DistanceKm,
	}
}

// RankProviders scores and ranks candidates, filtering below minimum threshold.
func RankProviders(candidates []MatchCandidate, serviceType booking.ServiceType, weightKg int, surgeMode bool) []MatchScore {
	var scores []MatchScore
	for _, c := range candidates {
		if !c.IsOnline {
			continue
		}
		score := ScoreProvider(&c, serviceType, weightKg, surgeMode)
		if score.TotalScore >= minMatchThreshold {
			scores = append(scores, score)
		}
	}
	sort.Slice(scores, func(i, j int) bool {
		return scores[i].TotalScore > scores[j].TotalScore
	})
	return scores
}

func scoreDistance(distanceKm float64) float64 {
	return math.Max(0, (1-distanceKm/30)*100)
}

func scoreTrustTier(tier user.TrustTier) float64 {
	if score, ok := trustTierScores[tier]; ok {
		return score
	}
	return 20 // default to basic
}

func scoreAcceptanceRate(rate float64) float64 {
	return rate * 100
}

func scoreVehicleCompat(truckType TruckType, maxWeightKg int, serviceType booking.ServiceType, weightKg int) float64 {
	capableTrucks, ok := serviceTypeCapabilities[serviceType]
	if !ok {
		return 0
	}
	capable := false
	for _, t := range capableTrucks {
		if t == truckType {
			capable = true
			break
		}
	}
	if !capable {
		return 0
	}
	if maxWeightKg < weightKg {
		return 0
	}
	return 100
}

func scoreCurrentLoad(activeJobs int) float64 {
	switch {
	case activeJobs == 0:
		return 100
	case activeJobs == 1:
		return 50
	default:
		return 0
	}
}
