package paymentuc

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/davidlramirez95/towcommand/internal/domain/user"
)

func TestCalculateCommission(t *testing.T) {
	tests := []struct {
		name           string
		amount         int64
		tier           user.TrustTier
		wantCommission int64
		wantRate       float64
	}{
		{
			name:           "basic tier 25%",
			amount:         100_000,
			tier:           user.TrustTierBasic,
			wantCommission: 25_000,
			wantRate:       0.25,
		},
		{
			name:           "verified tier 22%",
			amount:         100_000,
			tier:           user.TrustTierVerified,
			wantCommission: 22_000,
			wantRate:       0.22,
		},
		{
			name:           "suki silver tier 20%",
			amount:         100_000,
			tier:           user.TrustTierSukiSilver,
			wantCommission: 20_000,
			wantRate:       0.20,
		},
		{
			name:           "suki gold tier 18%",
			amount:         100_000,
			tier:           user.TrustTierSukiGold,
			wantCommission: 18_000,
			wantRate:       0.18,
		},
		{
			name:           "suki elite tier 15%",
			amount:         200_000,
			tier:           user.TrustTierSukiElite,
			wantCommission: 30_000,
			wantRate:       0.15,
		},
		{
			name:           "unknown tier falls back to basic 25%",
			amount:         100_000,
			tier:           user.TrustTier("unknown"),
			wantCommission: 25_000,
			wantRate:       0.25,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			commission, rate := CalculateCommission(tt.amount, tt.tier)
			assert.Equal(t, tt.wantCommission, commission)
			assert.InDelta(t, tt.wantRate, rate, 0.001)
		})
	}
}
