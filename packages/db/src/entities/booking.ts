import { KEY_PREFIXES, buildKey } from '../table-design';
import type { Booking, BookingStatus } from '@towcommand/core';

export function bookingKeys(bookingId: string) {
  return {
    PK: buildKey(KEY_PREFIXES.JOB, bookingId),
    SK: KEY_PREFIXES.DETAILS,
  };
}

export function bookingGSI1Keys(userId: string, createdAt: string) {
  return {
    GSI1PK: buildKey(KEY_PREFIXES.USER, userId),
    GSI1SK: buildKey(KEY_PREFIXES.JOB, createdAt),
  };
}

export function bookingGSI2Keys(status: string, createdAt: string) {
  return {
    GSI2PK: buildKey(KEY_PREFIXES.STATUS, status),
    GSI2SK: createdAt,
  };
}

export function bookingHistoryKeys(bookingId: string, timestamp: string) {
  return {
    PK: buildKey(KEY_PREFIXES.JOB, bookingId),
    SK: buildKey(KEY_PREFIXES.STATUS, timestamp),
  };
}

export function otpKeys(bookingId: string, type: 'PICKUP' | 'DROPOFF') {
  return {
    PK: buildKey(KEY_PREFIXES.JOB, bookingId),
    SK: buildKey(KEY_PREFIXES.OTP, type),
  };
}

export function chatMessageKeys(bookingId: string, timestamp: string) {
  return {
    PK: buildKey(KEY_PREFIXES.CHAT, bookingId),
    SK: buildKey(KEY_PREFIXES.MESSAGE, timestamp),
  };
}

export function evidenceKeys(bookingId: string, reportId: string) {
  return {
    PK: buildKey(KEY_PREFIXES.JOB, bookingId),
    SK: buildKey(KEY_PREFIXES.EVIDENCE, reportId),
  };
}

export function mediaKeys(bookingId: string, mediaId: string) {
  return {
    PK: buildKey(KEY_PREFIXES.JOB, bookingId),
    SK: buildKey(KEY_PREFIXES.MEDIA, mediaId),
  };
}

export function toBookingItem(booking: Booking) {
  return {
    ...bookingKeys(booking.bookingId),
    ...bookingGSI1Keys(booking.customerId, booking.createdAt),
    ...bookingGSI2Keys(booking.status, booking.createdAt),
    entityType: 'Booking',
    ...booking,
  };
}
