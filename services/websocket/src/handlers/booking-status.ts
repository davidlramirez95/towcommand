import type { APIGatewayProxyEvent, APIGatewayProxyResult } from 'aws-lambda';
import { BookingRepository } from '@towcommand/db';
import { SessionCache } from '@towcommand/cache';
import { ConnectionManager } from '../lib/connection-manager';

const bookingRepo = new BookingRepository();
const sessionCache = new SessionCache();
const connectionManager = new ConnectionManager();

export async function handler(event: APIGatewayProxyEvent): Promise<APIGatewayProxyResult> {
  try {
    const body = JSON.parse(event.body ?? '{}');
    const { bookingId, status, metadata } = body.data ?? {};

    if (!bookingId || !status) {
      return { statusCode: 400, body: 'Missing bookingId or status' };
    }

    const booking = await bookingRepo.getById(bookingId);
    if (!booking) {
      return { statusCode: 404, body: 'Booking not found' };
    }

    const statusPayload = {
      type: 'booking.status',
      data: {
        bookingId,
        status,
        previousStatus: booking.status,
        metadata,
        timestamp: new Date().toISOString(),
      },
    };

    // Notify the customer
    const customerConnId = await sessionCache.getWebSocketConnection(booking.customerId);
    if (customerConnId) {
      await connectionManager.sendMessage(customerConnId, statusPayload);
    }

    // Notify the provider if assigned
    if (booking.providerId) {
      const providerConnId = await sessionCache.getWebSocketConnection(booking.providerId);
      if (providerConnId) {
        await connectionManager.sendMessage(providerConnId, statusPayload);
      }
    }

    return { statusCode: 200, body: 'Status broadcast sent' };
  } catch (error) {
    console.error('Booking status broadcast error:', error);
    return { statusCode: 500, body: 'Broadcast failed' };
  }
}
