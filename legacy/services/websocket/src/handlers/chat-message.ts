import type { APIGatewayProxyEvent, APIGatewayProxyResult } from 'aws-lambda';
import { BookingRepository } from '@towcommand/db';
import { SessionCache } from '@towcommand/cache';
import { ConnectionManager } from '../lib/connection-manager';
import { ulid } from 'ulid';

const bookingRepo = new BookingRepository();
const sessionCache = new SessionCache();
const connectionManager = new ConnectionManager();

export async function handler(event: APIGatewayProxyEvent): Promise<APIGatewayProxyResult> {
  try {
    const connectionId = event.requestContext.connectionId as string;
    const body = JSON.parse(event.body ?? '{}');
    const { bookingId, senderId, message } = body.data ?? {};

    if (!bookingId || !senderId || !message) {
      return { statusCode: 400, body: 'Missing required fields' };
    }

    const booking = await bookingRepo.getById(bookingId);
    if (!booking) {
      return { statusCode: 404, body: 'Booking not found' };
    }

    // Determine the recipient (the other party in the booking)
    const recipientId = senderId === booking.customerId
      ? booking.providerId
      : booking.customerId;

    if (!recipientId) {
      return { statusCode: 400, body: 'No recipient for this booking' };
    }

    const chatPayload = {
      type: 'chat.receive',
      data: {
        messageId: ulid(),
        bookingId,
        senderId,
        message: message.substring(0, 1000), // Limit message length
        timestamp: new Date().toISOString(),
      },
    };

    // Send to recipient
    const recipientConnId = await sessionCache.getWebSocketConnection(recipientId);
    if (recipientConnId) {
      await connectionManager.sendMessage(recipientConnId, chatPayload);
    }

    return { statusCode: 200, body: 'Message sent' };
  } catch (error) {
    console.error('Chat message error:', error);
    return { statusCode: 500, body: 'Message failed' };
  }
}
