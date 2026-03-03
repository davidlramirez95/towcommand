/**
 * Feature flag service for TowCommand
 * Gates features by user tier (Suki loyalty program)
 * Reads FEATURE_FLAGS env var (JSON string baked in at deploy time)
 * Pattern adapted from gutguard-ai feature flags
 */

export interface FeatureFlags {
  maxBookingsPerDay: number;
  priorityMatching: boolean;
  realTimeTracking: boolean;
  chatWithProvider: boolean;
  multipleVehicles: boolean;
  scheduledBookings: boolean;
  corporateAccounts: boolean;
  advancedDiagnostics: boolean;
  loyaltyRewards: boolean;
  emergencySOS: boolean;
  priceGuarantee: boolean;
  cancellationFeeWaiver: boolean;
}

export type UserTier = 'basic' | 'suki-silver' | 'suki-gold' | 'suki-platinum' | 'provider' | 'admin';

type FlagConfig = Record<UserTier, FeatureFlags>;

const DEFAULT_FLAGS: FlagConfig = {
  basic: {
    maxBookingsPerDay: 3,
    priorityMatching: false,
    realTimeTracking: true,
    chatWithProvider: true,
    multipleVehicles: false,
    scheduledBookings: false,
    corporateAccounts: false,
    advancedDiagnostics: false,
    loyaltyRewards: false,
    emergencySOS: true,
    priceGuarantee: false,
    cancellationFeeWaiver: false,
  },
  'suki-silver': {
    maxBookingsPerDay: 5,
    priorityMatching: false,
    realTimeTracking: true,
    chatWithProvider: true,
    multipleVehicles: true,
    scheduledBookings: false,
    corporateAccounts: false,
    advancedDiagnostics: false,
    loyaltyRewards: true,
    emergencySOS: true,
    priceGuarantee: false,
    cancellationFeeWaiver: false,
  },
  'suki-gold': {
    maxBookingsPerDay: 10,
    priorityMatching: true,
    realTimeTracking: true,
    chatWithProvider: true,
    multipleVehicles: true,
    scheduledBookings: true,
    corporateAccounts: false,
    advancedDiagnostics: true,
    loyaltyRewards: true,
    emergencySOS: true,
    priceGuarantee: true,
    cancellationFeeWaiver: false,
  },
  'suki-platinum': {
    maxBookingsPerDay: -1,
    priorityMatching: true,
    realTimeTracking: true,
    chatWithProvider: true,
    multipleVehicles: true,
    scheduledBookings: true,
    corporateAccounts: true,
    advancedDiagnostics: true,
    loyaltyRewards: true,
    emergencySOS: true,
    priceGuarantee: true,
    cancellationFeeWaiver: true,
  },
  provider: {
    maxBookingsPerDay: -1,
    priorityMatching: false,
    realTimeTracking: true,
    chatWithProvider: true,
    multipleVehicles: false,
    scheduledBookings: false,
    corporateAccounts: false,
    advancedDiagnostics: false,
    loyaltyRewards: false,
    emergencySOS: true,
    priceGuarantee: false,
    cancellationFeeWaiver: false,
  },
  admin: {
    maxBookingsPerDay: -1,
    priorityMatching: true,
    realTimeTracking: true,
    chatWithProvider: true,
    multipleVehicles: true,
    scheduledBookings: true,
    corporateAccounts: true,
    advancedDiagnostics: true,
    loyaltyRewards: true,
    emergencySOS: true,
    priceGuarantee: true,
    cancellationFeeWaiver: true,
  },
};

let cachedConfig: FlagConfig | null = null;

function parseConfig(): FlagConfig {
  if (cachedConfig) return cachedConfig;

  const raw = process.env.FEATURE_FLAGS;
  if (!raw) {
    // Use defaults if no env var set
    cachedConfig = DEFAULT_FLAGS;
    return cachedConfig;
  }

  try {
    cachedConfig = JSON.parse(raw) as FlagConfig;
  } catch {
    console.warn('Failed to parse FEATURE_FLAGS env var, using defaults');
    cachedConfig = DEFAULT_FLAGS;
  }
  return cachedConfig;
}

export function getFeatureFlags(tier: UserTier): FeatureFlags {
  const config = parseConfig();
  const flags = config[tier];
  if (!flags) {
    // Fall back to basic tier for unknown tiers
    return config.basic;
  }
  return flags;
}

export function isFeatureEnabled(tier: UserTier, feature: keyof FeatureFlags): boolean {
  const flags = getFeatureFlags(tier);
  const value = flags[feature];
  if (typeof value === 'boolean') return value;
  // For numeric flags (e.g. maxBookingsPerDay), treat -1 (unlimited) or >0 as enabled
  if (typeof value === 'number') return value !== 0;
  return false;
}

/** Reset cached config - useful for testing */
export function resetFeatureFlagCache(): void {
  cachedConfig = null;
}
