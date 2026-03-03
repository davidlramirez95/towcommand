import type { APIGatewayProxyEvent, APIGatewayProxyResult } from 'aws-lambda';
import { updateLocationSchema, AppError, ErrorCode } from '@towcommand/core';
import { isValidPhilippineCoordinate } from '@towcommand/core';
import { GeoCache } from '@towcommand/cache';
import { ProviderRepository } from '@towcommand/db';
import { publishEvent, EVENT_CATALOG } from '@towcommand/events';
import { handleError, successResponse } from '../../middleware/error-handler';

const geoCache = new GeoCache();
const providerRepo = new ProviderRepository();

export async function handler(event: APIGatewayProxyEvent): Promise<APIGatewayProxyResult> {
  try {
    const providerId = event.requestContext.authorizer?.providerId as string;
    const body = updateLocationSchema.parse(JSON.parse(event.body ?? '{}'));

    if (!isValidPhilippineCoordinate(body.lat, body.lng)) {
      throw AppError.badRequest(ErrorCode.INVALID_LOCATION, 'Coordinates must be within the Philippines');
    }

    // Update geo cache for real-time matching
    await geoCache.updateProviderLocation(providerId, body.lat, body.lng);

    // Persist last known location to DynamoDB
    await providerRepo.update(providerId, {
      currentLat: body.lat,
      currentLng: body.lng,
      lastLocationUpdate: new Date().toISOString(),
    });

    await publishEvent(
      EVENT_CATALOG.tracking.source,
      EVENT_CATALOG.tracking.events.LocationUpdated,
      {
        providerId,
        lat: body.lat,
        lng: body.lng,
        heading: body.heading,
        speed: body.speed,
        timestamp: new Date().toISOString(),
      },
      { userId: providerId, userType: 'provider' },
    );

    return successResponse({ providerId, lat: body.lat, lng: body.lng });
  } catch (error) {
    return handleError(error);
  }
}
