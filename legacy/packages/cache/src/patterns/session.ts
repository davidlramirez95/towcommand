import { getRedisClient } from '../client';
import { CACHE_KEYS, CACHE_TTL } from '../keys';

export class SessionCache {
  private redis = getRedisClient();

  async setUserClaims(cognitoSub: string, claims: Record<string, unknown>): Promise<void> {
    await this.redis.setex(CACHE_KEYS.userClaims(cognitoSub), CACHE_TTL.userClaims, JSON.stringify(claims));
  }

  async getUserClaims(cognitoSub: string): Promise<Record<string, unknown> | null> {
    const data = await this.redis.get(CACHE_KEYS.userClaims(cognitoSub));
    return data ? JSON.parse(data) : null;
  }

  async setWebSocketConnection(userId: string, connectionId: string): Promise<void> {
    await this.redis.set(CACHE_KEYS.wsConnection(userId), connectionId);
  }

  async getWebSocketConnection(userId: string): Promise<string | null> {
    return this.redis.get(CACHE_KEYS.wsConnection(userId));
  }

  async removeWebSocketConnection(userId: string): Promise<void> {
    await this.redis.del(CACHE_KEYS.wsConnection(userId));
  }
}
