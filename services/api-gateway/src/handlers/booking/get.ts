import type { APIGatewayProxyEvent, APIGatewayProxyResult } from 'aws-lambda';
import { AppError, ErrorCode, UserType } from '@towcommand/core';
import { BookingRepository } from '@towcommand/db';
import { handleError, successResponse } from '../../middleware/error-handler';

const bookingRepo = new BookingRepository();

export async function handler(event: APIGatewayProxyEvent): Promise<APIGatewayProxyResult> {
  try {
    const userId = event.requestContext.authorizer?.userId as string;
    const userType = event.requestContext.authorizer?.userType as string;
    const bookingId = event.pathParameters?.id as string;

    const booking = await bookingRepo.getById(bookingId);

    if (!booking) {
      throw AppError.notFound('Booking');
    }

    // Customers can only view their own bookings; providers can view assigned bookings; admins see all
    const isOwner = booking.customerId === userId;
    const isAssignedProvider = booking.providerId === userId;
    const isAdmin = userType === UserType.ADMIN || userType === UserType.OPS_AGENT;

    if (!isOwner && !isAssignedProvider && !isAdmin) {
      throw AppError.forbidden('You do not have access to this booking');
    }

    return successResponse(booking);
  } catch (error) {
    return handleError(error);
  }
}
