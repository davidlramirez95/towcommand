package paymentuc

import "github.com/davidlramirez95/towcommand/internal/domain/user"

// commissionRates maps each trust tier to its commission rate.
// Higher tiers earn a lower platform commission as a loyalty incentive.
var commissionRates = map[user.TrustTier]float64{
	user.TrustTierBasic:      0.25,
	user.TrustTierVerified:   0.22,
	user.TrustTierSukiSilver: 0.20,
	user.TrustTierSukiGold:   0.18,
	user.TrustTierSukiElite:  0.15,
}

// CalculateCommission computes the platform commission for a given amount and
// provider trust tier. Unknown tiers fall back to the basic rate. All amounts
// are in centavos.
func CalculateCommission(amountCentavos int64, tier user.TrustTier) (commission int64, rate float64) {
	rate, ok := commissionRates[tier]
	if !ok {
		rate = commissionRates[user.TrustTierBasic]
	}
	commission = int64(float64(amountCentavos) * rate)
	return commission, rate
}
