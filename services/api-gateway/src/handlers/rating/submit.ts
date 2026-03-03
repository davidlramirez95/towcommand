import type { APIGatewayProxyEvent, APIGatewayProxyResult } from 'aws-lambda';
import { submitRatingSchema, AppError, BookingStatus } from '@towcommand/core';
import { BookingRepository, RatingRepository, ProviderRepository } from '@towcommand/db';
import { handleError, successResponse } from '../../middleware/error-handler';

const bookingRepo = new BookingRepository();
const ratingRepo = new RatingRepository();
const providerRepo = new ProviderRepository();

export async function handler(event: APIGatewayProxyEvent): Promise<APIGatewayProxyResult> {
  try {
    const userId = event.requestContext.authorizer?.userId as string;
    const body = submitRatingSchema.parse(JSON.parse(event.body ?? '{}'));

    const booking = await bookingRepo.getById(body.bookingId);

    if (!booking) {
      throw AppError.notFound('Booking');
    }

    if (booking.customerId !== userId) {
      throw AppError.forbidden('You can only rate your own bookings');
    }

    if (booking.status !== BookingStatus.COMPLETED) {
      throw AppError.badRequest('VALIDATION_ERROR' as any, 'Can only rate completed bookings');
    }

    if (!booking.providerId) {
      throw AppError.badRequest('VALIDATION_ERROR' as any, 'Booking has no assigned provider');
    }

    // Check for duplicate rating
    const existing = await ratingRepo.getByBooking(body.bookingId);
    if (existing) {
      throw AppError.badRequest('VALIDATION_ERROR' as any, 'You have already rated this booking');
    }

    const now = new Date().toISOString();

    await ratingRepo.create({
      bookingId: body.bookingId,
      customerId: userId,
      providerId: booking.providerId,
      rating: body.rating,
      comment: body.comment,
      tags: body.tags,
      createdAt: now,
    });

    // Update provider's average rating
    const provider = await providerRepo.getById(booking.providerId);
    if (provider) {
      const totalJobs = provider.totalJobsCompleted || 1;
      const newRating = ((provider.rating * totalJobs) + body.rating) / (totalJobs + 1);
      await providerRepo.update(booking.providerId, {
        rating: Math.round(newRating * 100) / 100,
      });
    }

    return successResponse({
      bookingId: body.bookingId,
      rating: body.rating,
      createdAt: now,
    }, 201);
  } catch (error) {
    return handleError(error);
  }
}
