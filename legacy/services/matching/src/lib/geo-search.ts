import { GeoCache } from '@towcommand/cache';
import { ProviderRepository } from '@towcommand/db';
import type { Provider, ProviderStatus } from '@towcommand/core';

const geoCache = new GeoCache();
const providerRepo = new ProviderRepository();

export interface NearbyProvider {
  providerId: string;
  distance: number;
  provider: Provider;
}

export class GeoSearchService {
  async findNearbyProviders(
    lat: number,
    lng: number,
    radiusKm: number,
    city = 'NCR',
  ): Promise<NearbyProvider[]> {
    const nearby = await geoCache.getNearbyProviders(lat, lng, radiusKm, city, 50);

    const results: NearbyProvider[] = [];
    for (const { providerId, distance } of nearby) {
      const provider = await providerRepo.getById(providerId);
      if (provider && provider.isOnline && provider.status === ('active' as ProviderStatus)) {
        results.push({ providerId, distance, provider });
      }
    }

    return results.sort((a, b) => a.distance - b.distance);
  }

  async findProvidersInZone(zoneId: string): Promise<string[]> {
    // Zone-based fallback: query DynamoDB for providers in a service area
    const providers = await providerRepo.getByTierAndCity('basic', zoneId);
    return providers
      .filter((p) => p.isOnline)
      .map((p) => p.providerId);
  }
}
