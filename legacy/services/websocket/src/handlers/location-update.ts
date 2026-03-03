import type { APIGatewayProxyEvent, APIGatewayProxyResult } from 'aws-lambda';
import { GeoCache, SessionCache } from '@towcommand/cache';
import { BookingRepository } from '@towcommand/db';
import { estimateEtaMinutes } from '@towcommand/core';
import { ConnectionManager } from '../lib/connection-manager';

const geoCache = new GeoCache();
const sessionCache = new SessionCache();
const bookingRepo = new BookingRepository();
const connectionManager = new ConnectionManager();

export async function handler(event: APIGatewayProxyEvent): Promise<APIGatewayProxyResult> {
  try {
    const body = JSON.parse(event.body ?? '{}');
    const { action, data } = body;
    const { providerId, lat, lng, heading, speed, bookingId } = data ?? {};

    if (!providerId || !lat || !lng) {
      return { statusCode: 400, body: 'Missing required fields' };
    }

    // Update provider location in geo cache
    await geoCache.updateProviderLocation(providerId, lat, lng);

    // If there's an active booking, broadcast location to the customer
    if (bookingId) {
      const booking = await bookingRepo.getById(bookingId);
      if (booking?.customerId) {
        const customerConnectionId = await sessionCache.getWebSocketConnection(booking.customerId);
        if (customerConnectionId) {
          const distanceToPickup = 0; // Would be calculated from booking pickup
          await connectionManager.sendMessage(customerConnectionId, {
            type: 'location.broadcast',
            data: {
              providerId,
              lat,
              lng,
              heading,
              speed,
              eta: estimateEtaMinutes(distanceToPickup),
              bookingId,
              timestamp: new Date().toISOString(),
            },
          });
        }
      }
    }

    return { statusCode: 200, body: 'Location updated' };
  } catch (error) {
    console.error('Location update error:', error);
    return { statusCode: 500, body: 'Update failed' };
  }
}
