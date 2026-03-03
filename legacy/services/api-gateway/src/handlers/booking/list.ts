import type { APIGatewayProxyEvent, APIGatewayProxyResult } from 'aws-lambda';
import { BookingStatus, UserType } from '@towcommand/core';
import { BookingRepository } from '@towcommand/db';
import { handleError, successResponse } from '../../middleware/error-handler';

const bookingRepo = new BookingRepository();

export async function handler(event: APIGatewayProxyEvent): Promise<APIGatewayProxyResult> {
  try {
    const userId = event.requestContext.authorizer?.userId as string;
    const userType = event.requestContext.authorizer?.userType as string;
    const params = event.queryStringParameters ?? {};
    const limit = Math.min(parseInt(params.limit ?? '25', 10), 100);
    const statusFilter = params.status as BookingStatus | undefined;

    let bookings;

    if (userType === UserType.ADMIN || userType === UserType.OPS_AGENT) {
      // Admins can filter by status across all bookings
      if (statusFilter) {
        bookings = await bookingRepo.getByStatus(statusFilter, limit);
      } else {
        bookings = await bookingRepo.getByUser(userId, limit);
      }
    } else {
      // Customers and providers see only their own bookings
      bookings = await bookingRepo.getByUser(userId, limit);
    }

    // Client-side status filter for non-admin users
    if (statusFilter && userType !== UserType.ADMIN && userType !== UserType.OPS_AGENT) {
      bookings = bookings.filter((b) => b.status === statusFilter);
    }

    return successResponse({ items: bookings, count: bookings.length });
  } catch (error) {
    return handleError(error);
  }
}
