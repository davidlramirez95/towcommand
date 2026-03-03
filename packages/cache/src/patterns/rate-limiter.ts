import { getRedisClient } from '../client';
import { CACHE_KEYS, CACHE_TTL } from '../keys';

export class RateLimiter {
  private redis = getRedisClient();

  async checkLimit(userId: string, maxRequests = 100): Promise<{ allowed: boolean; remaining: number; resetIn: number }> {
    const key = CACHE_KEYS.rateLimit(userId);
    const current = await this.redis.incr(key);

    if (current === 1) {
      await this.redis.expire(key, CACHE_TTL.rateLimit);
    }

    const ttl = await this.redis.ttl(key);

    return {
      allowed: current <= maxRequests,
      remaining: Math.max(0, maxRequests - current),
      resetIn: ttl > 0 ? ttl : CACHE_TTL.rateLimit,
    };
  }

  async acquireJobLock(jobId: string): Promise<boolean> {
    const key = CACHE_KEYS.jobLock(jobId);
    const result = await this.redis.set(key, '1', 'EX', CACHE_TTL.jobLock, 'NX');
    return result === 'OK';
  }

  async releaseJobLock(jobId: string): Promise<void> {
    await this.redis.del(CACHE_KEYS.jobLock(jobId));
  }
}
