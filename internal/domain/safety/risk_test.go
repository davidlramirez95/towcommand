package safety

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestComputeRiskScore_AllCombinations(t *testing.T) {
	// Test all 32 combinations of the 5 boolean-like factors.
	// DistanceKm is treated as 3 states but for the 2^5 matrix we test
	// distance <= 20 (off) vs distance > 30 (on). The >20 band is tested
	// separately below.
	boolVals := []bool{false, true}

	for _, nightTime := range boolVals {
		for _, basicTier := range boolVals {
			for _, longDistance := range boolVals {
				for _, priorSOS := range boolVals {
					for _, highRisk := range boolVals {
						name := fmt.Sprintf("night=%t/basic=%t/far=%t/sos=%t/zone=%t",
							nightTime, basicTier, longDistance, priorSOS, highRisk)

						t.Run(name, func(t *testing.T) {
							tier := "verified"
							if basicTier {
								tier = "basic"
							}

							dist := 10.0
							if longDistance {
								dist = 35.0
							}

							sosCount := 0
							if priorSOS {
								sosCount = 2
							}

							factors := RiskFactors{
								IsNightTime:       nightTime,
								ProviderTrustTier: tier,
								DistanceKm:        dist,
								PriorSOSCount:     sosCount,
								IsHighRiskZone:    highRisk,
							}

							got := ComputeRiskScore(factors)

							want := 0
							if nightTime {
								want += 25
							}
							if basicTier {
								want += 15
							}
							if longDistance {
								want += 25
							}
							if priorSOS {
								want += 20
							}
							if highRisk {
								want += 15
							}
							if want > 100 {
								want = 100
							}

							assert.Equal(t, want, got)
						})
					}
				}
			}
		}
	}
}

func TestComputeRiskScore_MediumDistanceBand(t *testing.T) {
	// Test the >20 but <=30 distance band.
	factors := RiskFactors{
		IsNightTime:       false,
		ProviderTrustTier: "verified",
		DistanceKm:        25.0,
		PriorSOSCount:     0,
		IsHighRiskZone:    false,
	}
	got := ComputeRiskScore(factors)
	assert.Equal(t, 15, got)
}

func TestComputeRiskScore_DistanceBoundary20(t *testing.T) {
	// Exactly 20 km should not add any distance score.
	factors := RiskFactors{
		DistanceKm: 20.0,
	}
	assert.Equal(t, 0, ComputeRiskScore(factors))
}

func TestComputeRiskScore_DistanceBoundary30(t *testing.T) {
	// Exactly 30 km falls in the >20 band (score +15).
	factors := RiskFactors{
		DistanceKm: 30.0,
	}
	assert.Equal(t, 15, ComputeRiskScore(factors))
}

func TestComputeRiskScore_ClampAt100(t *testing.T) {
	// All factors active: 25+15+25+20+15 = 100, exactly at cap.
	factors := RiskFactors{
		IsNightTime:       true,
		ProviderTrustTier: "basic",
		DistanceKm:        50.0,
		PriorSOSCount:     5,
		IsHighRiskZone:    true,
	}
	got := ComputeRiskScore(factors)
	assert.Equal(t, 100, got)
	assert.LessOrEqual(t, got, 100, "score must never exceed 100")
}

func TestComputeRiskScore_ZeroFactors(t *testing.T) {
	factors := RiskFactors{
		IsNightTime:       false,
		ProviderTrustTier: "suki_gold",
		DistanceKm:        5.0,
		PriorSOSCount:     0,
		IsHighRiskZone:    false,
	}
	assert.Equal(t, 0, ComputeRiskScore(factors))
}
