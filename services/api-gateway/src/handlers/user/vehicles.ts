import type { APIGatewayProxyEvent, APIGatewayProxyResult } from 'aws-lambda';
import { WeightClass } from '@towcommand/core';
import { UserRepository } from '@towcommand/db';
import { handleError, successResponse } from '../../middleware/error-handler';
import { z } from 'zod';
import { ulid } from 'ulid';

const userRepo = new UserRepository();

const addVehicleSchema = z.object({
  make: z.string().min(1).max(100),
  model: z.string().min(1).max(100),
  year: z.number().int().min(1980).max(new Date().getFullYear() + 1),
  plateNumber: z.string().regex(/^[A-Z0-9]{3,4}-?[A-Z0-9]{3,4}$/i, 'Invalid plate number format'),
  weightClass: z.nativeEnum(WeightClass),
  color: z.string().min(1).max(50),
  isDefault: z.boolean().optional(),
});

export async function handler(event: APIGatewayProxyEvent): Promise<APIGatewayProxyResult> {
  try {
    const userId = event.requestContext.authorizer?.userId as string;
    const method = event.httpMethod;

    if (method === 'GET') {
      const vehicles = await userRepo.getVehicles(userId);
      return successResponse({ vehicles, count: vehicles.length });
    }

    if (method === 'POST') {
      const body = addVehicleSchema.parse(JSON.parse(event.body ?? '{}'));

      const vehicle = {
        vehicleId: `VEH-${ulid()}`,
        userId,
        make: body.make,
        model: body.model,
        year: body.year,
        plateNumber: body.plateNumber.toUpperCase(),
        weightClass: body.weightClass,
        color: body.color,
        isDefault: body.isDefault ?? false,
        createdAt: new Date().toISOString(),
      };

      await userRepo.addVehicle(vehicle);

      return successResponse(vehicle, 201);
    }

    return successResponse({ message: 'Method not allowed' }, 405);
  } catch (error) {
    return handleError(error);
  }
}
