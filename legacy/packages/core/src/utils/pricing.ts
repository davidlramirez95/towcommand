import { ServiceType, WeightClass, PriceBreakdown } from '../types';

interface PricingTier {
  baseRate: number;
  perKmRate: number;
}

const MMDA_RATES: Record<WeightClass, PricingTier> = {
  [WeightClass.MOTORCYCLE]: { baseRate: 800, perKmRate: 100 },
  [WeightClass.LIGHT]: { baseRate: 1800, perKmRate: 200 },
  [WeightClass.MEDIUM]: { baseRate: 3000, perKmRate: 250 },
  [WeightClass.HEAVY]: { baseRate: 5000, perKmRate: 350 },
  [WeightClass.SUPER_HEAVY]: { baseRate: 8000, perKmRate: 500 },
};

const ROADSIDE_RATES: Record<string, number> = {
  [ServiceType.JUMPSTART]: 500,
  [ServiceType.TIRE_CHANGE]: 600,
  [ServiceType.FUEL_DELIVERY]: 400,
  [ServiceType.LOCKOUT]: 700,
};

export function calculatePrice(
  serviceType: ServiceType,
  weightClass: WeightClass,
  distanceKm: number,
  options: {
    surgeMultiplier?: number;
    isNightTime?: boolean;
    isHoliday?: boolean;
    isAccident?: boolean;
  } = {},
): PriceBreakdown {
  const { surgeMultiplier = 1.0, isNightTime = false, isHoliday = false, isAccident = false } = options;
  const clampedSurge = Math.min(surgeMultiplier, 1.5);

  let base: number;
  let distance: number;

  if (serviceType in ROADSIDE_RATES) {
    base = ROADSIDE_RATES[serviceType] ?? 500;
    distance = 0;
  } else {
    const tier = MMDA_RATES[weightClass];
    base = tier.baseRate;
    distance = Math.ceil(distanceKm) * tier.perKmRate;
  }

  const weight = weightClass === WeightClass.SUPER_HEAVY ? 1000 : 0;
  const timeSurcharge = isNightTime ? Math.round(base * 0.2) : 0;
  const holidaySurcharge = isHoliday ? Math.round(base * 0.15) : 0;
  const accidentSurcharge = isAccident ? 500 : 0;

  const subtotal = base + distance + weight + timeSurcharge + holidaySurcharge + accidentSurcharge;
  const surgePricing = Math.round(subtotal * (clampedSurge - 1));
  const total = subtotal + surgePricing;

  return {
    base,
    distance,
    weight,
    timeSurcharge: timeSurcharge + holidaySurcharge + accidentSurcharge,
    surgePricing,
    total,
    currency: 'PHP',
    surgeMultiplier: clampedSurge > 1 ? clampedSurge : undefined,
  };
}

export function isNightTime(): boolean {
  const hour = new Date().getUTCHours() + 8; // PHT = UTC+8
  const phtHour = hour >= 24 ? hour - 24 : hour;
  return phtHour >= 22 || phtHour < 6;
}

export function calculateCancellationFee(
  status: string,
  distanceKmTravelled: number,
): number {
  switch (status) {
    case 'PENDING': return 0;
    case 'MATCHED': return 100;
    case 'EN_ROUTE': return Math.min(100 + distanceKmTravelled * 30, 500);
    case 'ARRIVED': return 500 + Math.min(distanceKmTravelled * 30, 300);
    default: return 0;
  }
}
