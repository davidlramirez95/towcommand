import { getRedisClient } from '../client';
import { CACHE_KEYS } from '../keys';

export class SurgePricingCache {
  private redis = getRedisClient();

  async getSurgeMultiplier(region: string): Promise<number> {
    const value = await this.redis.get(CACHE_KEYS.surgeMultiplier(region));
    return value ? parseFloat(value) : 1.0;
  }

  async setSurgeMultiplier(region: string, multiplier: number, ttlSeconds = 1800): Promise<void> {
    const clamped = Math.min(Math.max(multiplier, 1.0), 1.5);
    await this.redis.setex(CACHE_KEYS.surgeMultiplier(region), ttlSeconds, String(clamped));
  }
}
