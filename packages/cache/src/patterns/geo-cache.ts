import { getRedisClient } from '../client';
import { CACHE_KEYS, CACHE_TTL } from '../keys';

export class GeoCache {
  private redis = getRedisClient();

  async updateProviderLocation(providerId: string, lat: number, lng: number, city = 'NCR'): Promise<void> {
    const pipeline = this.redis.pipeline();
    pipeline.geoadd(CACHE_KEYS.providerGeo(city), lng, lat, providerId);
    pipeline.setex(CACHE_KEYS.providerLocation(providerId), CACHE_TTL.providerLocation, JSON.stringify({ lat, lng, updatedAt: Date.now() }));
    await pipeline.exec();
  }

  async getNearbyProviders(lat: number, lng: number, radiusKm: number, city = 'NCR', limit = 20): Promise<Array<{ providerId: string; distance: number }>> {
    const results = await this.redis.georadius(CACHE_KEYS.providerGeo(city), lng, lat, radiusKm, 'km', 'WITHDIST', 'ASC', 'COUNT', limit);
    return (results as Array<[string, string]>).map(([providerId, dist]) => ({
      providerId,
      distance: parseFloat(dist),
    }));
  }

  async removeProvider(providerId: string, city = 'NCR'): Promise<void> {
    const pipeline = this.redis.pipeline();
    pipeline.zrem(CACHE_KEYS.providerGeo(city), providerId);
    pipeline.del(CACHE_KEYS.providerLocation(providerId));
    await pipeline.exec();
  }

  async getProviderLocation(providerId: string): Promise<{ lat: number; lng: number } | null> {
    const data = await this.redis.get(CACHE_KEYS.providerLocation(providerId));
    return data ? JSON.parse(data) : null;
  }
}
