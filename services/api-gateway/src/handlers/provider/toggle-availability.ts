import type { APIGatewayProxyEvent, APIGatewayProxyResult } from 'aws-lambda';
import { AppError, ErrorCode, ProviderStatus } from '@towcommand/core';
import { GeoCache } from '@towcommand/cache';
import { ProviderRepository } from '@towcommand/db';
import { publishEvent, EVENT_CATALOG } from '@towcommand/events';
import { handleError, successResponse } from '../../middleware/error-handler';

const geoCache = new GeoCache();
const providerRepo = new ProviderRepository();

export async function handler(event: APIGatewayProxyEvent): Promise<APIGatewayProxyResult> {
  try {
    const providerId = event.requestContext.authorizer?.providerId as string;
    const body = JSON.parse(event.body ?? '{}');
    const goOnline = body.online === true;

    const provider = await providerRepo.getById(providerId);

    if (!provider) {
      throw AppError.notFound('Provider');
    }

    if (provider.status !== ProviderStatus.ACTIVE) {
      throw AppError.badRequest(
        ErrorCode.PROVIDER_NOT_VERIFIED,
        'Only verified providers can go online',
      );
    }

    await providerRepo.update(providerId, { isOnline: goOnline });

    if (goOnline) {
      // Add to geo index so they appear in nearby searches
      if (provider.currentLat && provider.currentLng) {
        await geoCache.updateProviderLocation(
          providerId, provider.currentLat, provider.currentLng,
          provider.serviceAreas[0] ?? 'NCR',
        );
      }

      await publishEvent(
        EVENT_CATALOG.provider.source,
        EVENT_CATALOG.provider.events.ProviderOnline,
        {
          providerId,
          lat: provider.currentLat,
          lng: provider.currentLng,
          serviceArea: provider.serviceAreas[0] ?? 'NCR',
        },
        { userId: providerId, userType: 'provider' },
      );
    } else {
      // Remove from geo index
      await geoCache.removeProvider(providerId, provider.serviceAreas[0] ?? 'NCR');

      await publishEvent(
        EVENT_CATALOG.provider.source,
        EVENT_CATALOG.provider.events.ProviderOffline,
        { providerId },
        { userId: providerId, userType: 'provider' },
      );
    }

    return successResponse({ providerId, online: goOnline });
  } catch (error) {
    return handleError(error);
  }
}
