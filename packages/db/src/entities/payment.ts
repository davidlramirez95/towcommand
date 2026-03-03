import { KEY_PREFIXES, buildKey } from '../table-design';
import type { Payment } from '@towcommand/core';

export function paymentKeys(transactionId: string) {
  return {
    PK: buildKey(KEY_PREFIXES.TRANSACTION, transactionId),
    SK: KEY_PREFIXES.DETAILS,
  };
}

export function toPaymentItem(payment: Payment) {
  return {
    ...paymentKeys(payment.paymentId),
    GSI1PK: buildKey(KEY_PREFIXES.JOB, payment.bookingId),
    GSI1SK: buildKey(KEY_PREFIXES.PAYMENT, payment.createdAt),
    entityType: 'Payment',
    ...payment,
  };
}
