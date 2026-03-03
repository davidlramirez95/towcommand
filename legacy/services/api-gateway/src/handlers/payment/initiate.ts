import type { APIGatewayProxyEvent, APIGatewayProxyResult } from 'aws-lambda';
import { initiatePaymentSchema, AppError, ErrorCode, PaymentStatus, BookingStatus } from '@towcommand/core';
import { BookingRepository, PaymentRepository } from '@towcommand/db';
import { publishEvent, EVENT_CATALOG } from '@towcommand/events';
import { handleError, successResponse } from '../../middleware/error-handler';
import { ulid } from 'ulid';

const bookingRepo = new BookingRepository();
const paymentRepo = new PaymentRepository();

export async function handler(event: APIGatewayProxyEvent): Promise<APIGatewayProxyResult> {
  try {
    const userId = event.requestContext.authorizer?.userId as string;
    const body = initiatePaymentSchema.parse(JSON.parse(event.body ?? '{}'));

    const booking = await bookingRepo.getById(body.bookingId);

    if (!booking) {
      throw AppError.notFound('Booking');
    }

    if (booking.customerId !== userId) {
      throw AppError.forbidden('You can only pay for your own bookings');
    }

    // Payment can be initiated for completed or in-transit bookings
    if (booking.status !== BookingStatus.COMPLETED && booking.status !== BookingStatus.IN_TRANSIT) {
      throw AppError.badRequest(
        ErrorCode.PAYMENT_FAILED,
        `Cannot initiate payment for booking in ${booking.status} status`,
      );
    }

    const paymentId = `PAY-${ulid()}`;
    const now = new Date().toISOString();

    const payment = {
      paymentId,
      bookingId: body.bookingId,
      userId,
      amount: body.amount,
      currency: 'PHP' as const,
      method: body.method,
      status: PaymentStatus.PENDING,
      createdAt: now,
      updatedAt: now,
    };

    await paymentRepo.create(payment);

    await publishEvent(
      EVENT_CATALOG.payment.source,
      EVENT_CATALOG.payment.events.PaymentInitiated,
      {
        paymentId,
        bookingId: body.bookingId,
        amount: body.amount,
        method: body.method,
      },
      { userId, userType: 'customer' },
    );

    return successResponse({
      paymentId,
      status: PaymentStatus.PENDING,
      amount: body.amount,
      method: body.method,
    }, 201);
  } catch (error) {
    return handleError(error);
  }
}
