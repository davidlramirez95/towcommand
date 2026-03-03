import { BaseRepository } from './base.repo';
import { paymentKeys, toPaymentItem } from '../entities/payment';
import { KEY_PREFIXES } from '../table-design';
import type { Payment } from '@towcommand/core';

export class PaymentRepository extends BaseRepository {
  async getById(paymentId: string): Promise<Payment | null> {
    const { PK, SK } = paymentKeys(paymentId);
    return this.getItem<Payment>(PK, SK);
  }

  async create(payment: Payment): Promise<void> {
    await this.putItem(toPaymentItem(payment));
  }

  async update(paymentId: string, updates: Partial<Payment>): Promise<void> {
    const { PK, SK } = paymentKeys(paymentId);
    await this.updateItem(PK, SK, { ...updates, updatedAt: new Date().toISOString() });
  }

  async getByBooking(bookingId: string): Promise<Payment[]> {
    return this.query<Payment>({
      IndexName: 'GSI1-UserJobs',
      KeyConditionExpression: 'GSI1PK = :pk AND begins_with(GSI1SK, :sk)',
      ExpressionAttributeValues: {
        ':pk': `${KEY_PREFIXES.JOB}${bookingId}`,
        ':sk': `${KEY_PREFIXES.PAYMENT}`,
      },
      ScanIndexForward: false,
    });
  }
}
