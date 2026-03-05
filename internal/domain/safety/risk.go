package safety

// RiskFactors holds the inputs for risk score calculation.
type RiskFactors struct {
	IsNightTime       bool
	ProviderTrustTier string
	DistanceKm        float64
	PriorSOSCount     int
	IsHighRiskZone    bool
}

// ComputeRiskScore calculates a risk score (0-100) from the given factors.
//
// Weights:
//   - Night time (20:00-05:00 PHT): +25
//   - Basic trust tier provider: +15
//   - Distance > 30km: +25, Distance > 20km: +15
//   - Prior SOS history: +20
//   - High risk zone: +15
//
// The score is clamped to 100.
func ComputeRiskScore(factors RiskFactors) int {
	score := 0
	if factors.IsNightTime {
		score += 25
	}
	if factors.ProviderTrustTier == "basic" {
		score += 15
	}
	if factors.DistanceKm > 30 {
		score += 25
	} else if factors.DistanceKm > 20 {
		score += 15
	}
	if factors.PriorSOSCount > 0 {
		score += 20
	}
	if factors.IsHighRiskZone {
		score += 15
	}
	if score > 100 {
		score = 100
	}
	return score
}
