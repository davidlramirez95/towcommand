import { KEY_PREFIXES, buildKey } from '../table-design';

export interface Rating {
  bookingId: string;
  customerId: string;
  providerId: string;
  rating: number;
  comment?: string;
  tags?: string[];
  createdAt: string;
}

export function ratingKeys(bookingId: string) {
  return {
    PK: buildKey(KEY_PREFIXES.JOB, bookingId),
    SK: KEY_PREFIXES.RATING,
  };
}

export function toRatingItem(r: Rating) {
  return {
    ...ratingKeys(r.bookingId),
    GSI1PK: buildKey(KEY_PREFIXES.PROVIDER, r.providerId),
    GSI1SK: `RATE#${r.createdAt}`,
    entityType: 'Rating',
    ...r,
  };
}
