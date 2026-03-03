import type { APIGatewayProxyEvent, APIGatewayProxyResult } from 'aws-lambda';
import { AppError, ErrorCode, BookingStatus, UserType, VALID_STATUS_TRANSITIONS, updateStatusSchema } from '@towcommand/core';
import { BookingRepository } from '@towcommand/db';
import { publishEvent, EVENT_CATALOG } from '@towcommand/events';
import { handleError, successResponse } from '../../middleware/error-handler';

const bookingRepo = new BookingRepository();

// Statuses that only providers can set
const PROVIDER_STATUSES = new Set<BookingStatus>([
  BookingStatus.EN_ROUTE,
  BookingStatus.ARRIVED,
  BookingStatus.CONDITION_REPORT,
  BookingStatus.LOADING,
  BookingStatus.IN_TRANSIT,
  BookingStatus.ARRIVED_DROPOFF,
  BookingStatus.COMPLETED,
]);

export async function handler(event: APIGatewayProxyEvent): Promise<APIGatewayProxyResult> {
  try {
    const userId = event.requestContext.authorizer?.userId as string;
    const userType = event.requestContext.authorizer?.userType as string;
    const bookingId = event.pathParameters?.id as string;
    const body = updateStatusSchema.parse(JSON.parse(event.body ?? '{}'));

    const booking = await bookingRepo.getById(bookingId);

    if (!booking) {
      throw AppError.notFound('Booking');
    }

    // Authorization: providers can update their assigned bookings, admins can update any
    const isAssignedProvider = booking.providerId === userId;
    const isAdmin = userType === UserType.ADMIN || userType === UserType.OPS_AGENT;

    if (PROVIDER_STATUSES.has(body.status) && !isAssignedProvider && !isAdmin) {
      throw AppError.forbidden('Only the assigned provider can update this status');
    }

    // Validate status transition
    const allowedTransitions = VALID_STATUS_TRANSITIONS[booking.status];
    if (!allowedTransitions.includes(body.status)) {
      throw AppError.badRequest(
        ErrorCode.INVALID_STATUS_TRANSITION,
        `Cannot transition from ${booking.status} to ${body.status}`,
        { currentStatus: booking.status, requestedStatus: body.status, allowed: allowedTransitions },
      );
    }

    await bookingRepo.updateStatus(bookingId, body.status, {
      changedBy: userId,
      ...body.metadata,
    });

    await publishEvent(
      EVENT_CATALOG.booking.source,
      EVENT_CATALOG.booking.events.BookingStatusChanged,
      {
        bookingId,
        previousStatus: booking.status,
        newStatus: body.status,
        changedBy: userId,
        metadata: body.metadata,
      },
      { userId, userType },
    );

    // Publish completion event for downstream processing
    if (body.status === BookingStatus.COMPLETED && booking.providerId) {
      await publishEvent(
        EVENT_CATALOG.booking.source,
        EVENT_CATALOG.booking.events.BookingCompleted,
        {
          bookingId,
          customerId: booking.customerId,
          providerId: booking.providerId,
          price: booking.price,
        },
        { userId, userType },
      );
    }

    return successResponse({
      bookingId,
      previousStatus: booking.status,
      status: body.status,
    });
  } catch (error) {
    return handleError(error);
  }
}
