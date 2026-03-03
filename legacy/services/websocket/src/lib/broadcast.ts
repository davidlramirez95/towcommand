import { SessionCache } from '@towcommand/cache';
import { BookingRepository } from '@towcommand/db';
import { ConnectionManager } from './connection-manager';

const sessionCache = new SessionCache();
const bookingRepo = new BookingRepository();
const connectionManager = new ConnectionManager();

export class BroadcastService {
  async broadcastToUser(userId: string, data: unknown): Promise<void> {
    const connectionId = await sessionCache.getWebSocketConnection(userId);
    if (connectionId) {
      await connectionManager.sendMessage(connectionId, data);
    }
  }

  async broadcastToBooking(bookingId: string, data: unknown): Promise<void> {
    const booking = await bookingRepo.getById(bookingId);
    if (!booking) return;

    const connectionIds: string[] = [];

    const customerConn = await sessionCache.getWebSocketConnection(booking.customerId);
    if (customerConn) connectionIds.push(customerConn);

    if (booking.providerId) {
      const providerConn = await sessionCache.getWebSocketConnection(booking.providerId);
      if (providerConn) connectionIds.push(providerConn);
    }

    await connectionManager.broadcast(connectionIds, data);
  }

  async broadcastToProviders(region: string, data: unknown): Promise<void> {
    // In production, maintain a Redis set of online provider connectionIds per region
    // For MVP, this would be populated by the connect handler
    console.log(`Broadcasting to providers in ${region}:`, data);
  }
}
