import type { APIGatewayProxyEvent, APIGatewayProxyResult } from 'aws-lambda';
import { RatingRepository, ProviderRepository } from '@towcommand/db';
import { handleError, successResponse } from '../../middleware/error-handler';

const ratingRepo = new RatingRepository();
const providerRepo = new ProviderRepository();

export async function handler(event: APIGatewayProxyEvent): Promise<APIGatewayProxyResult> {
  try {
    const providerId = event.pathParameters?.id as string;
    const limit = Math.min(
      parseInt(event.queryStringParameters?.limit ?? '25', 10),
      100,
    );

    const provider = await providerRepo.getById(providerId);
    const ratings = await ratingRepo.getByProvider(providerId, limit);

    const avgRating = provider?.rating ?? 0;
    const totalReviews = ratings.length;

    return successResponse({
      providerId,
      averageRating: avgRating,
      totalReviews,
      ratings: ratings.map((r) => ({
        bookingId: r.bookingId,
        rating: r.rating,
        comment: r.comment,
        tags: r.tags,
        createdAt: r.createdAt,
      })),
    });
  } catch (error) {
    return handleError(error);
  }
}
