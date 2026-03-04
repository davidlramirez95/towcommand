package provider

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/davidlramirez95/towcommand/internal/domain/booking"
	"github.com/davidlramirez95/towcommand/internal/domain/user"
)

func TestScoreDistance(t *testing.T) {
	tests := []struct {
		name       string
		distanceKm float64
		want       float64
	}{
		{"zero distance", 0, 100},
		{"15km", 15, 50},
		{"30km", 30, 0},
		{"beyond 30km", 40, 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.InDelta(t, tt.want, scoreDistance(tt.distanceKm), 0.01)
		})
	}
}

func TestScoreTrustTier(t *testing.T) {
	tests := []struct {
		tier user.TrustTier
		want float64
	}{
		{user.TrustTierBasic, 20},
		{user.TrustTierVerified, 40},
		{user.TrustTierSukiSilver, 60},
		{user.TrustTierSukiGold, 80},
		{user.TrustTierSukiElite, 100},
	}
	for _, tt := range tests {
		t.Run(string(tt.tier), func(t *testing.T) {
			assert.Equal(t, tt.want, scoreTrustTier(tt.tier))
		})
	}
}

func TestScoreTrustTier_UnknownDefaultsToBasic(t *testing.T) {
	assert.Equal(t, 20.0, scoreTrustTier("unknown_tier"))
}

func TestScoreAcceptanceRate(t *testing.T) {
	tests := []struct {
		name string
		rate float64
		want float64
	}{
		{"zero rate", 0, 0},
		{"50%", 0.5, 50},
		{"perfect", 1.0, 100},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.InDelta(t, tt.want, scoreAcceptanceRate(tt.rate), 0.01)
		})
	}
}

func TestScoreVehicleCompat(t *testing.T) {
	tests := []struct {
		name     string
		truck    TruckType
		maxKg    int
		svcType  booking.ServiceType
		weightKg int
		want     float64
	}{
		{"flatbed for flatbed tow with enough capacity", TruckTypeFlatbed, 5000, booking.ServiceTypeFlatbedTow, 1500, 100},
		{"wheel_lift for flatbed tow is incompatible", TruckTypeWheelLift, 5000, booking.ServiceTypeFlatbedTow, 1500, 0},
		{"flatbed but insufficient capacity", TruckTypeFlatbed, 1000, booking.ServiceTypeFlatbedTow, 1500, 0},
		{"any truck for jumpstart", TruckTypeMotorcycleCarrier, 200, booking.ServiceTypeJumpstart, 0, 100},
		{"boom for accident recovery", TruckTypeBoom, 8000, booking.ServiceTypeAccidentRecovery, 5000, 100},
		{"wheel_lift for wheel lift service", TruckTypeWheelLift, 3000, booking.ServiceTypeWheelLift, 2000, 100},
		{"unknown service type returns zero", TruckTypeFlatbed, 5000, "UNKNOWN_SERVICE", 1500, 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, scoreVehicleCompat(tt.truck, tt.maxKg, tt.svcType, tt.weightKg))
		})
	}
}

func TestScoreCurrentLoad(t *testing.T) {
	assert.Equal(t, 100.0, scoreCurrentLoad(0))
	assert.Equal(t, 50.0, scoreCurrentLoad(1))
	assert.Equal(t, 0.0, scoreCurrentLoad(2))
	assert.Equal(t, 0.0, scoreCurrentLoad(5))
}

func TestScoreProvider_DefaultWeights(t *testing.T) {
	candidate := MatchCandidate{
		ProviderID:          "prov-1",
		TrustTier:           user.TrustTierSukiGold,
		AcceptanceRate:      0.95,
		TruckType:           TruckTypeFlatbed,
		MaxWeightCapacityKg: 5000,
		ActiveJobCount:      0,
		DistanceKm:          5,
		IsOnline:            true,
	}

	score := ScoreProvider(&candidate, booking.ServiceTypeFlatbedTow, 1500, false)

	assert.Equal(t, "prov-1", score.ProviderID)
	assert.Greater(t, score.TotalScore, 0.70)
	assert.InDelta(t, 83.33, score.Factors.Distance, 0.5)
	assert.Equal(t, 80.0, score.Factors.TrustTier)
	assert.InDelta(t, 95.0, score.Factors.AcceptanceRate, 0.01)
	assert.Equal(t, 100.0, score.Factors.VehicleCompat)
	assert.Equal(t, 100.0, score.Factors.CurrentLoad)
}

func TestScoreProvider_SurgeMode(t *testing.T) {
	candidate := MatchCandidate{
		ProviderID:          "prov-1",
		TrustTier:           user.TrustTierBasic,
		AcceptanceRate:      0.5,
		TruckType:           TruckTypeFlatbed,
		MaxWeightCapacityKg: 5000,
		ActiveJobCount:      0,
		DistanceKm:          2,
		IsOnline:            true,
	}

	normal := ScoreProvider(&candidate, booking.ServiceTypeFlatbedTow, 1500, false)
	surge := ScoreProvider(&candidate, booking.ServiceTypeFlatbedTow, 1500, true)

	// In surge, distance weight increases from 0.40 to 0.50 -- nearby providers score higher.
	assert.Greater(t, surge.TotalScore, normal.TotalScore)
}

func TestScoreProvider_IncompatibleTruckGetsLowScore(t *testing.T) {
	candidate := MatchCandidate{
		ProviderID:          "prov-bad",
		TrustTier:           user.TrustTierSukiElite,
		AcceptanceRate:      1.0,
		TruckType:           TruckTypeMotorcycleCarrier,
		MaxWeightCapacityKg: 300,
		ActiveJobCount:      0,
		DistanceKm:          1,
		IsOnline:            true,
	}

	score := ScoreProvider(&candidate, booking.ServiceTypeFlatbedTow, 1500, false)

	// Vehicle compat is 0, so total should be reduced significantly.
	assert.Equal(t, 0.0, score.Factors.VehicleCompat)
	assert.Less(t, score.TotalScore, 0.90)
}

func TestRankProviders(t *testing.T) {
	candidates := []MatchCandidate{
		{ProviderID: "far-basic", TrustTier: user.TrustTierBasic, AcceptanceRate: 0.3, TruckType: TruckTypeFlatbed, MaxWeightCapacityKg: 5000, ActiveJobCount: 2, DistanceKm: 25, IsOnline: true},
		{ProviderID: "close-elite", TrustTier: user.TrustTierSukiElite, AcceptanceRate: 0.99, TruckType: TruckTypeFlatbed, MaxWeightCapacityKg: 5000, ActiveJobCount: 0, DistanceKm: 3, IsOnline: true},
		{ProviderID: "mid-gold", TrustTier: user.TrustTierSukiGold, AcceptanceRate: 0.8, TruckType: TruckTypeFlatbed, MaxWeightCapacityKg: 5000, ActiveJobCount: 0, DistanceKm: 10, IsOnline: true},
		{ProviderID: "offline", TrustTier: user.TrustTierSukiElite, AcceptanceRate: 0.99, TruckType: TruckTypeFlatbed, MaxWeightCapacityKg: 5000, ActiveJobCount: 0, DistanceKm: 1, IsOnline: false},
	}

	scores := RankProviders(candidates, booking.ServiceTypeFlatbedTow, 1500, false)

	assert.GreaterOrEqual(t, len(scores), 2) // at least close-elite and mid-gold pass threshold
	assert.Equal(t, "close-elite", scores[0].ProviderID)
	// Offline provider should not appear.
	for _, s := range scores {
		assert.NotEqual(t, "offline", s.ProviderID)
	}
	// Scores should be descending.
	for i := 1; i < len(scores); i++ {
		assert.GreaterOrEqual(t, scores[i-1].TotalScore, scores[i].TotalScore)
	}
}

func TestRankProviders_IncompatibleTruck(t *testing.T) {
	candidates := []MatchCandidate{
		{ProviderID: "motorcycle-carrier", TrustTier: user.TrustTierSukiElite, AcceptanceRate: 0.99, TruckType: TruckTypeMotorcycleCarrier, MaxWeightCapacityKg: 300, ActiveJobCount: 0, DistanceKm: 1, IsOnline: true},
	}

	scores := RankProviders(candidates, booking.ServiceTypeFlatbedTow, 1500, false)

	// Motorcycle carrier can't do flatbed tow — vehicle compat is 0,
	// but other high factors keep total above threshold.
	// Verify that vehicle compat factor is zero.
	if len(scores) > 0 {
		assert.Equal(t, 0.0, scores[0].Factors.VehicleCompat)
	}
}

func TestRankProviders_EmptyInput(t *testing.T) {
	scores := RankProviders(nil, booking.ServiceTypeFlatbedTow, 1500, false)
	assert.Empty(t, scores)
}

func TestRankProviders_AllOffline(t *testing.T) {
	candidates := []MatchCandidate{
		{ProviderID: "off-1", IsOnline: false, TruckType: TruckTypeFlatbed, MaxWeightCapacityKg: 5000, TrustTier: user.TrustTierSukiElite, AcceptanceRate: 0.99, DistanceKm: 1},
		{ProviderID: "off-2", IsOnline: false, TruckType: TruckTypeFlatbed, MaxWeightCapacityKg: 5000, TrustTier: user.TrustTierSukiGold, AcceptanceRate: 0.9, DistanceKm: 2},
	}

	scores := RankProviders(candidates, booking.ServiceTypeFlatbedTow, 1500, false)
	assert.Empty(t, scores)
}

func TestWeightClassToKg(t *testing.T) {
	assert.Equal(t, 200, WeightClassToKg(user.WeightClassMotorcycle))
	assert.Equal(t, 1500, WeightClassToKg(user.WeightClassLight))
	assert.Equal(t, 2500, WeightClassToKg(user.WeightClassMedium))
	assert.Equal(t, 5000, WeightClassToKg(user.WeightClassHeavy))
	assert.Equal(t, 10000, WeightClassToKg(user.WeightClassSuperHeavy))
	assert.Equal(t, 1500, WeightClassToKg("unknown")) // default
}
