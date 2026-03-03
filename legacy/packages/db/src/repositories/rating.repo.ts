import { BaseRepository } from './base.repo';
import { ratingKeys, toRatingItem, type Rating } from '../entities/rating';
import { KEY_PREFIXES } from '../table-design';

export class RatingRepository extends BaseRepository {
  async getByBooking(bookingId: string): Promise<Rating | null> {
    const { PK, SK } = ratingKeys(bookingId);
    return this.getItem<Rating>(PK, SK);
  }

  async create(rating: Rating): Promise<void> {
    await this.putItem(toRatingItem(rating));
  }

  async getByProvider(providerId: string, limit = 25): Promise<Rating[]> {
    return this.query<Rating>({
      IndexName: 'GSI1-UserJobs',
      KeyConditionExpression: 'GSI1PK = :pk AND begins_with(GSI1SK, :sk)',
      ExpressionAttributeValues: {
        ':pk': `${KEY_PREFIXES.PROVIDER}${providerId}`,
        ':sk': 'RATE#',
      },
      ScanIndexForward: false,
      Limit: limit,
    });
  }
}
