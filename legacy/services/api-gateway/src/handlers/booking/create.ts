import type { APIGatewayProxyEvent, APIGatewayProxyResult } from 'aws-lambda';
import { createBookingSchema, BookingStatus, WeightClass } from '@towcommand/core';
import { haversineDistance, isNightTime, calculatePrice } from '@towcommand/core';
import { BookingRepository } from '@towcommand/db';
import { publishEvent, EVENT_CATALOG } from '@towcommand/events';
import { handleError, successResponse } from '../../middleware/error-handler';
import { ulid } from 'ulid';

const bookingRepo = new BookingRepository();

export async function handler(event: APIGatewayProxyEvent): Promise<APIGatewayProxyResult> {
  try {
    const userId = event.requestContext.authorizer?.userId as string;
    const body = createBookingSchema.parse(JSON.parse(event.body ?? '{}'));

    const bookingId = `TC-${new Date().getFullYear()}-${ulid()}`;

    const distanceKm = haversineDistance(
      body.pickupLocation.lat, body.pickupLocation.lng,
      body.dropoffLocation.lat, body.dropoffLocation.lng,
    );

    // Default to LIGHT if weight class not derivable from estimate
    const weightClass = WeightClass.LIGHT;

    const price = calculatePrice(body.serviceType, weightClass, distanceKm, {
      isNightTime: isNightTime(),
    });

    const now = new Date().toISOString();

    const booking = {
      bookingId,
      customerId: userId,
      vehicleId: body.vehicleId,
      serviceType: body.serviceType,
      status: BookingStatus.PENDING,
      pickupLocation: body.pickupLocation,
      dropoffLocation: body.dropoffLocation,
      weightClass,
      price,
      estimateId: body.estimateId,
      notes: body.notes,
      createdAt: now,
      updatedAt: now,
    };

    await bookingRepo.create(booking);

    await publishEvent(
      EVENT_CATALOG.booking.source,
      EVENT_CATALOG.booking.events.BookingCreated,
      {
        bookingId,
        customerId: userId,
        serviceType: body.serviceType,
        pickupLocation: body.pickupLocation,
        dropoffLocation: body.dropoffLocation,
        price,
      },
      { userId, userType: 'customer' },
    );

    return successResponse(booking, 201);
  } catch (error) {
    return handleError(error);
  }
}
