import { BaseRepository } from './base.repo';
import { bookingKeys, toBookingItem, bookingHistoryKeys } from '../entities/booking';
import type { Booking, BookingStatus } from '@towcommand/core';

export class BookingRepository extends BaseRepository {
  async getById(bookingId: string): Promise<Booking | null> {
    const { PK, SK } = bookingKeys(bookingId);
    return this.getItem<Booking>(PK, SK);
  }

  async create(booking: Booking): Promise<void> {
    await this.putItem(toBookingItem(booking));
  }

  async updateStatus(bookingId: string, status: BookingStatus, metadata?: Record<string, unknown>): Promise<void> {
    const { PK, SK } = bookingKeys(bookingId);
    const now = new Date().toISOString();
    await this.updateItem(PK, SK, { status, updatedAt: now, GSI2PK: `STATUS#${status}`, GSI2SK: now });

    await this.putItem({
      ...bookingHistoryKeys(bookingId, now),
      entityType: 'BookingHistory',
      status,
      changedAt: now,
      metadata,
    });
  }

  async getByUser(userId: string, limit = 25): Promise<Booking[]> {
    return this.query<Booking>({
      IndexName: 'GSI1-UserJobs',
      KeyConditionExpression: 'GSI1PK = :pk AND begins_with(GSI1SK, :sk)',
      ExpressionAttributeValues: { ':pk': `USER#${userId}`, ':sk': 'JOB#' },
      ScanIndexForward: false,
      Limit: limit,
    });
  }

  async getByStatus(status: BookingStatus, limit = 50): Promise<Booking[]> {
    return this.query<Booking>({
      IndexName: 'GSI2-StatusJobs',
      KeyConditionExpression: 'GSI2PK = :pk',
      ExpressionAttributeValues: { ':pk': `STATUS#${status}` },
      ScanIndexForward: false,
      Limit: limit,
    });
  }
}
