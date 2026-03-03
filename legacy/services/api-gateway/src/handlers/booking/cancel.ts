import type { APIGatewayProxyEvent, APIGatewayProxyResult } from 'aws-lambda';
import { AppError, ErrorCode, BookingStatus, VALID_STATUS_TRANSITIONS } from '@towcommand/core';
import { calculateCancellationFee } from '@towcommand/core';
import { BookingRepository } from '@towcommand/db';
import { publishEvent, EVENT_CATALOG } from '@towcommand/events';
import { handleError, successResponse } from '../../middleware/error-handler';

const bookingRepo = new BookingRepository();

export async function handler(event: APIGatewayProxyEvent): Promise<APIGatewayProxyResult> {
  try {
    const userId = event.requestContext.authorizer?.userId as string;
    const bookingId = event.pathParameters?.id as string;
    const body = JSON.parse(event.body ?? '{}');
    const reason = body.reason as string | undefined;

    const booking = await bookingRepo.getById(bookingId);

    if (!booking) {
      throw AppError.notFound('Booking');
    }

    if (booking.customerId !== userId) {
      throw AppError.forbidden('You can only cancel your own bookings');
    }

    const allowedTransitions = VALID_STATUS_TRANSITIONS[booking.status];
    if (!allowedTransitions.includes(BookingStatus.CANCELLED)) {
      throw AppError.badRequest(
        ErrorCode.CANNOT_CANCEL,
        `Cannot cancel booking in ${booking.status} status`,
      );
    }

    const cancellationFee = calculateCancellationFee(booking.status, 0);

    await bookingRepo.updateStatus(bookingId, BookingStatus.CANCELLED, {
      cancelledBy: userId,
      reason,
      cancellationFee,
    });

    await publishEvent(
      EVENT_CATALOG.booking.source,
      EVENT_CATALOG.booking.events.BookingCancelled,
      {
        bookingId,
        customerId: userId,
        providerId: booking.providerId,
        previousStatus: booking.status,
        reason,
        cancellationFee,
      },
      { userId, userType: 'customer' },
    );

    return successResponse({
      bookingId,
      status: BookingStatus.CANCELLED,
      cancellationFee,
    });
  } catch (error) {
    return handleError(error);
  }
}
