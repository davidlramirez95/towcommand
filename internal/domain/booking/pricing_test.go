package booking

import (
	"testing"
	"time"

	"github.com/davidlramirez95/towcommand/internal/domain/user"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// pht returns a time.Time in Asia/Manila (UTC+8).
func pht(year, month, day, hour, minute int) time.Time {
	return time.Date(year, time.Month(month), day, hour, minute, 0, 0, manilaLoc)
}

// --- IsNightTime ---

func TestIsNightTime(t *testing.T) {
	tests := []struct {
		name string
		time time.Time
		want bool
	}{
		{"22:00 is night", pht(2025, 1, 6, 22, 0), true},
		{"23:59 is night", pht(2025, 1, 6, 23, 59), true},
		{"00:00 is night", pht(2025, 1, 7, 0, 0), true},
		{"05:59 is night", pht(2025, 1, 7, 5, 59), true},
		{"06:00 is not night", pht(2025, 1, 7, 6, 0), false},
		{"12:00 is not night", pht(2025, 1, 7, 12, 0), false},
		{"21:59 is not night", pht(2025, 1, 6, 21, 59), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, IsNightTime(tt.time))
		})
	}
}

// --- IsWeekend ---

func TestIsWeekend(t *testing.T) {
	tests := []struct {
		name string
		time time.Time
		want bool
	}{
		// 2025-01-06 is a Monday
		{"Monday", pht(2025, 1, 6, 10, 0), false},
		{"Tuesday", pht(2025, 1, 7, 10, 0), false},
		{"Wednesday", pht(2025, 1, 8, 10, 0), false},
		{"Thursday", pht(2025, 1, 9, 10, 0), false},
		{"Friday", pht(2025, 1, 10, 10, 0), false},
		{"Saturday", pht(2025, 1, 11, 10, 0), true},
		{"Sunday", pht(2025, 1, 12, 10, 0), true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, IsWeekend(tt.time))
		})
	}
}

// --- CalculateEstimate ---

func TestCalculateEstimate(t *testing.T) {
	// Weekday daytime: Monday 2025-01-06 10:00 PHT
	weekdayDay := pht(2025, 1, 6, 10, 0)
	// Weekday night: Monday 2025-01-06 23:00 PHT
	weekdayNight := pht(2025, 1, 6, 23, 0)
	// Weekend night: Saturday 2025-01-11 23:00 PHT
	weekendNight := pht(2025, 1, 11, 23, 0)
	// Weekend day: Saturday 2025-01-11 10:00 PHT
	weekendDay := pht(2025, 1, 11, 10, 0)

	tests := []struct {
		name            string
		serviceType     ServiceType
		weightClass     user.WeightClass
		distanceKm      float64
		pickupTime      time.Time
		surgeMultiplier float64
		isAccident      bool
		want            PriceBreakdown
	}{
		{
			name:            "light 10km daytime no surge",
			serviceType:     ServiceTypeFlatbedTow,
			weightClass:     user.WeightClassLight,
			distanceKm:      10,
			pickupTime:      weekdayDay,
			surgeMultiplier: 1.0,
			want: PriceBreakdown{
				Base:     180_000,
				Distance: 200_000,
				Total:    380_000,
				Currency: "PHP",
			},
		},
		{
			name:            "heavy 25km nighttime 1.5x surge",
			serviceType:     ServiceTypeFlatbedTow,
			weightClass:     user.WeightClassHeavy,
			distanceKm:      25,
			pickupTime:      weekdayNight,
			surgeMultiplier: 1.5,
			want: PriceBreakdown{
				Base:            500_000,
				Distance:        875_000,
				TimeSurcharge:   150_000, // 30% of 500_000
				SurgePricing:    762_500, // (500k+875k+150k) * 0.5
				Total:           2_287_500,
				Currency:        "PHP",
				SurgeMultiplier: 1.5,
			},
		},
		{
			name:            "light 5km weekend night combined surcharges",
			serviceType:     ServiceTypeFlatbedTow,
			weightClass:     user.WeightClassLight,
			distanceKm:      5,
			pickupTime:      weekendNight,
			surgeMultiplier: 1.0,
			want: PriceBreakdown{
				Base:          180_000,
				Distance:      100_000,
				TimeSurcharge: 81_000, // night 54_000 + weekend 27_000
				Total:         361_000,
				Currency:      "PHP",
			},
		},
		{
			name:            "motorcycle 15km daytime no surge",
			serviceType:     ServiceTypeFlatbedTow,
			weightClass:     user.WeightClassMotorcycle,
			distanceKm:      15,
			pickupTime:      weekdayDay,
			surgeMultiplier: 1.0,
			want: PriceBreakdown{
				Base:     80_000,
				Distance: 150_000,
				Total:    230_000,
				Currency: "PHP",
			},
		},
		{
			name:            "super heavy 10km daytime no surge",
			serviceType:     ServiceTypeFlatbedTow,
			weightClass:     user.WeightClassSuperHeavy,
			distanceKm:      10,
			pickupTime:      weekdayDay,
			surgeMultiplier: 1.0,
			want: PriceBreakdown{
				Base:     800_000,
				Distance: 500_000,
				Weight:   100_000,
				Total:    1_400_000,
				Currency: "PHP",
			},
		},
		{
			name:            "medium 10km weekend day no surge",
			serviceType:     ServiceTypeWheelLift,
			weightClass:     user.WeightClassMedium,
			distanceKm:      10,
			pickupTime:      weekendDay,
			surgeMultiplier: 1.0,
			want: PriceBreakdown{
				Base:          300_000,
				Distance:      250_000,
				TimeSurcharge: 45_000, // 15% of 300_000
				Total:         595_000,
				Currency:      "PHP",
			},
		},
		{
			name:            "accident recovery heavy 10km daytime",
			serviceType:     ServiceTypeAccidentRecovery,
			weightClass:     user.WeightClassHeavy,
			distanceKm:      10,
			pickupTime:      weekdayDay,
			surgeMultiplier: 1.0,
			isAccident:      true,
			want: PriceBreakdown{
				Base:          500_000,
				Distance:      350_000,
				TimeSurcharge: 50_000, // accident only
				Total:         900_000,
				Currency:      "PHP",
			},
		},
		{
			name:            "zero distance",
			serviceType:     ServiceTypeFlatbedTow,
			weightClass:     user.WeightClassLight,
			distanceKm:      0,
			pickupTime:      weekdayDay,
			surgeMultiplier: 1.0,
			want: PriceBreakdown{
				Base:     180_000,
				Distance: 0,
				Total:    180_000,
				Currency: "PHP",
			},
		},
		{
			name:            "max surge 2.5x",
			serviceType:     ServiceTypeFlatbedTow,
			weightClass:     user.WeightClassLight,
			distanceKm:      10,
			pickupTime:      weekdayDay,
			surgeMultiplier: 3.0, // clamped to 2.5
			want: PriceBreakdown{
				Base:            180_000,
				Distance:        200_000,
				SurgePricing:    570_000, // 380_000 * 1.5
				Total:           950_000,
				Currency:        "PHP",
				SurgeMultiplier: 2.5,
			},
		},
		{
			name:            "fractional distance rounds up",
			serviceType:     ServiceTypeFlatbedTow,
			weightClass:     user.WeightClassLight,
			distanceKm:      10.1,
			pickupTime:      weekdayDay,
			surgeMultiplier: 1.0,
			want: PriceBreakdown{
				Base:     180_000,
				Distance: 220_000, // 10.1 km rounds up to 11 km
				Total:    400_000,
				Currency: "PHP",
			},
		},
		{
			name:            "surge below 1.0 clamped to 1.0",
			serviceType:     ServiceTypeFlatbedTow,
			weightClass:     user.WeightClassLight,
			distanceKm:      10,
			pickupTime:      weekdayDay,
			surgeMultiplier: 0.5, // clamped to 1.0
			want: PriceBreakdown{
				Base:     180_000,
				Distance: 200_000,
				Total:    380_000,
				Currency: "PHP",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := CalculateEstimate(
				tt.serviceType, tt.weightClass, tt.distanceKm,
				tt.pickupTime, tt.surgeMultiplier, tt.isAccident,
			)
			assert.Equal(t, tt.want, got)
		})
	}
}

// --- Roadside services ---

func TestCalculateEstimate_RoadsideServices(t *testing.T) {
	weekdayDay := pht(2025, 1, 6, 10, 0)

	tests := []struct {
		name        string
		serviceType ServiceType
		wantBase    int64
	}{
		{"jumpstart", ServiceTypeJumpstart, 50_000},
		{"tire change", ServiceTypeTireChange, 60_000},
		{"fuel delivery", ServiceTypeFuelDelivery, 40_000},
		{"lockout", ServiceTypeLockout, 70_000},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := CalculateEstimate(
				tt.serviceType, user.WeightClassLight, 15, // distance ignored
				weekdayDay, 1.0, false,
			)
			assert.Equal(t, tt.wantBase, got.Base)
			assert.Equal(t, int64(0), got.Distance, "roadside service should have zero distance")
			assert.Equal(t, tt.wantBase, got.Total)
			assert.Equal(t, "PHP", got.Currency)
		})
	}
}

// --- Boundary times ---

func TestCalculateEstimate_BoundaryTimes(t *testing.T) {
	lightBase := int64(180_000)

	tests := []struct {
		name          string
		pickupTime    time.Time
		wantNightPart int64
	}{
		{"21:59 no night surcharge", pht(2025, 1, 6, 21, 59), 0},
		{"22:00 night surcharge", pht(2025, 1, 6, 22, 0), roundCentavos(float64(lightBase) * nightSurchargeRate)},
		{"05:59 night surcharge", pht(2025, 1, 7, 5, 59), roundCentavos(float64(lightBase) * nightSurchargeRate)},
		{"06:00 no night surcharge", pht(2025, 1, 7, 6, 0), 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := CalculateEstimate(
				ServiceTypeFlatbedTow, user.WeightClassLight, 0,
				tt.pickupTime, 1.0, false,
			)
			assert.Equal(t, tt.wantNightPart, got.TimeSurcharge)
		})
	}
}

// --- CalculateFinalPrice ---

func TestCalculateFinalPrice(t *testing.T) {
	tests := []struct {
		name           string
		estimate       PriceBreakdown
		waitingMinutes int
		tollFees       int64
		wantWaiting    int64
		wantTotal      int64
	}{
		{
			name: "no waiting no tolls",
			estimate: PriceBreakdown{
				Total:    380_000,
				Currency: "PHP",
			},
			waitingMinutes: 0,
			tollFees:       0,
			wantWaiting:    0,
			wantTotal:      385_000, // 380k + 5k platform
		},
		{
			name: "30min waiting (free)",
			estimate: PriceBreakdown{
				Total:    380_000,
				Currency: "PHP",
			},
			waitingMinutes: 30,
			tollFees:       0,
			wantWaiting:    0,
			wantTotal:      385_000,
		},
		{
			name: "31min waiting (1 block)",
			estimate: PriceBreakdown{
				Total:    380_000,
				Currency: "PHP",
			},
			waitingMinutes: 31,
			tollFees:       0,
			wantWaiting:    10_000,
			wantTotal:      395_000,
		},
		{
			name: "60min waiting + tolls",
			estimate: PriceBreakdown{
				Total:    380_000,
				Currency: "PHP",
			},
			waitingMinutes: 60,
			tollFees:       15_000,
			wantWaiting:    20_000, // 2 blocks
			wantTotal:      420_000,
		},
		{
			name: "45min waiting (1 block)",
			estimate: PriceBreakdown{
				Total:    100_000,
				Currency: "PHP",
			},
			waitingMinutes: 45,
			tollFees:       0,
			wantWaiting:    10_000,
			wantTotal:      115_000,
		},
		{
			name: "44min waiting (1 block, ceil)",
			estimate: PriceBreakdown{
				Total:    100_000,
				Currency: "PHP",
			},
			waitingMinutes: 44,
			tollFees:       0,
			wantWaiting:    10_000,
			wantTotal:      115_000,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := CalculateFinalPrice(tt.estimate, tt.waitingMinutes, tt.tollFees)
			assert.Equal(t, tt.wantWaiting, got.WaitingFee)
			assert.Equal(t, tt.tollFees, got.TollFees)
			assert.Equal(t, platformFeeCentavos, got.PlatformFee)
			assert.Equal(t, tt.wantTotal, got.Total)
			assert.Equal(t, "PHP", got.Currency)
			assert.Equal(t, tt.estimate, got.Estimate)
		})
	}
}

// --- CalculateCancellationFee ---

func TestCalculateCancellationFee(t *testing.T) {
	tests := []struct {
		name    string
		status  BookingStatus
		distKm  float64
		wantFee int64
	}{
		{"pending", BookingStatusPending, 0, 0},
		{"matched", BookingStatusMatched, 10, 10_000},
		{"en_route 5km", BookingStatusEnRoute, 5, 25_000},
		{"en_route 20km (capped)", BookingStatusEnRoute, 20, 50_000},
		{"en_route 0km", BookingStatusEnRoute, 0, 10_000},
		{"arrived 5km", BookingStatusArrived, 5, 65_000},
		{"arrived 15km (dist capped)", BookingStatusArrived, 15, 80_000},
		{"arrived 0km", BookingStatusArrived, 0, 50_000},
		{"loading (non-cancellable)", BookingStatusLoading, 0, 0},
		{"completed", BookingStatusCompleted, 0, 0},
		{"cancelled", BookingStatusCancelled, 0, 0},
		{"unknown status", BookingStatus("UNKNOWN"), 0, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := CalculateCancellationFee(tt.status, tt.distKm)
			assert.Equal(t, tt.wantFee, got)
		})
	}
}

// --- GetSurgeMultiplier ---

func TestGetSurgeMultiplier(t *testing.T) {
	tests := []struct {
		name   string
		demand int
		supply int
		want   float64
	}{
		{"no demand", 0, 10, 1.0},
		{"equal demand/supply", 10, 10, 1.0},
		{"2:1 ratio", 20, 10, 1.375},
		{"3:1 ratio", 30, 10, 1.75},
		{"5:1 ratio (max)", 50, 10, 2.5},
		{"10:1 ratio (clamped)", 100, 10, 2.5},
		{"no supply", 10, 0, 2.5},
		{"both zero", 0, 0, 2.5},
		{"negative demand", -5, 10, 1.0},
		{"negative supply", 10, -5, 2.5},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GetSurgeMultiplier(tt.demand, tt.supply)
			assert.InDelta(t, tt.want, got, 0.001)
		})
	}
}

// --- Waiting fee helper ---

func TestCalculateWaitingFee(t *testing.T) {
	tests := []struct {
		name    string
		minutes int
		want    int64
	}{
		{"0 minutes", 0, 0},
		{"29 minutes", 29, 0},
		{"30 minutes (free)", 30, 0},
		{"31 minutes (1 block)", 31, 10_000},
		{"45 minutes (1 block)", 45, 10_000},
		{"46 minutes (2 blocks)", 46, 20_000},
		{"60 minutes (2 blocks)", 60, 20_000},
		{"90 minutes (4 blocks)", 90, 40_000},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := calculateWaitingFee(tt.minutes)
			assert.Equal(t, tt.want, got)
		})
	}
}

// --- MMDA rate coverage ---

func TestCalculateEstimate_AllWeightClasses(t *testing.T) {
	weekdayDay := pht(2025, 1, 6, 10, 0)

	tests := []struct {
		name        string
		weightClass user.WeightClass
		wantBase    int64
		wantPerKm   int64
	}{
		{"motorcycle", user.WeightClassMotorcycle, 80_000, 10_000},
		{"light", user.WeightClassLight, 180_000, 20_000},
		{"medium", user.WeightClassMedium, 300_000, 25_000},
		{"heavy", user.WeightClassHeavy, 500_000, 35_000},
		{"super_heavy", user.WeightClassSuperHeavy, 800_000, 50_000},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := CalculateEstimate(
				ServiceTypeFlatbedTow, tt.weightClass, 1,
				weekdayDay, 1.0, false,
			)
			assert.Equal(t, tt.wantBase, got.Base)
			assert.Equal(t, tt.wantPerKm, got.Distance) // 1 km
		})
	}
}

// --- All amounts in centavos ---

func TestAllAmountsInCentavos(t *testing.T) {
	weekdayDay := pht(2025, 1, 6, 10, 0)
	pb := CalculateEstimate(
		ServiceTypeFlatbedTow, user.WeightClassLight, 10,
		weekdayDay, 1.0, false,
	)

	require.Equal(t, "PHP", pb.Currency)
	assert.Greater(t, pb.Base, int64(0), "base should be positive centavos")
	assert.Greater(t, pb.Distance, int64(0), "distance should be positive centavos")
	assert.Equal(t, pb.Base+pb.Distance+pb.Weight+pb.TimeSurcharge+pb.SurgePricing, pb.Total,
		"total should equal sum of components")
}

// --- Surge multiplier not set when 1.0 ---

func TestSurgeMultiplierOmittedAtOne(t *testing.T) {
	weekdayDay := pht(2025, 1, 6, 10, 0)
	pb := CalculateEstimate(
		ServiceTypeFlatbedTow, user.WeightClassLight, 10,
		weekdayDay, 1.0, false,
	)
	assert.Equal(t, float64(0), pb.SurgeMultiplier, "surge multiplier should be zero-value when 1.0")
}

func TestSurgeMultiplierSetAboveOne(t *testing.T) {
	weekdayDay := pht(2025, 1, 6, 10, 0)
	pb := CalculateEstimate(
		ServiceTypeFlatbedTow, user.WeightClassLight, 10,
		weekdayDay, 1.5, false,
	)
	assert.Equal(t, 1.5, pb.SurgeMultiplier)
}
