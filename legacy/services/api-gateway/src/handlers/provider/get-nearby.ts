import type { APIGatewayProxyEvent, APIGatewayProxyResult } from 'aws-lambda';
import { AppError, ErrorCode } from '@towcommand/core';
import { isValidPhilippineCoordinate, estimateEtaMinutes } from '@towcommand/core';
import { GeoCache } from '@towcommand/cache';
import { ProviderRepository } from '@towcommand/db';
import { handleError, successResponse } from '../../middleware/error-handler';

const geoCache = new GeoCache();
const providerRepo = new ProviderRepository();

export async function handler(event: APIGatewayProxyEvent): Promise<APIGatewayProxyResult> {
  try {
    const queryParams = event.queryStringParameters ?? {};
    const lat = parseFloat(queryParams.lat ?? '0');
    const lng = parseFloat(queryParams.lng ?? '0');
    const radiusKm = Math.min(parseInt(queryParams.radius ?? '10', 10), 50);
    const limit = Math.min(parseInt(queryParams.limit ?? '10', 10), 20);
    const city = queryParams.city ?? 'NCR';

    if (!isValidPhilippineCoordinate(lat, lng)) {
      throw AppError.badRequest(ErrorCode.INVALID_LOCATION, 'Coordinates must be within the Philippines');
    }

    const nearby = await geoCache.getNearbyProviders(lat, lng, radiusKm, city, limit);

    // Enrich with provider details
    const providers = await Promise.all(
      nearby.map(async ({ providerId, distance }) => {
        const provider = await providerRepo.getById(providerId);
        if (!provider || !provider.isOnline) return null;
        return {
          providerId,
          name: provider.name,
          truckType: provider.truckType,
          rating: provider.rating,
          totalJobsCompleted: provider.totalJobsCompleted,
          plateNumber: provider.plateNumber,
          distanceKm: Math.round(distance * 10) / 10,
          etaMinutes: estimateEtaMinutes(distance),
        };
      }),
    );

    const filtered = providers.filter(Boolean);

    return successResponse({ providers: filtered, count: filtered.length });
  } catch (error) {
    return handleError(error);
  }
}
