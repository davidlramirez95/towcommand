import { KEY_PREFIXES, buildKey } from '../table-design';

export interface SukiTier {
  userId: string;
  points: number;
  tier: 'silver' | 'gold' | 'elite';
  totalSpent: number;
  totalBookings: number;
  lastEarnedAt: string;
  updatedAt: string;
}

export function sukiKeys(userId: string) {
  return {
    PK: buildKey(KEY_PREFIXES.USER, userId),
    SK: KEY_PREFIXES.SUKI,
  };
}
