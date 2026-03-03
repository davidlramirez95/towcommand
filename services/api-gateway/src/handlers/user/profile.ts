import type { APIGatewayProxyEvent, APIGatewayProxyResult } from 'aws-lambda';
import { AppError } from '@towcommand/core';
import { UserRepository } from '@towcommand/db';
import { handleError, successResponse } from '../../middleware/error-handler';

const userRepo = new UserRepository();

const ALLOWED_PROFILE_FIELDS = new Set(['name', 'phone', 'language']);

export async function handler(event: APIGatewayProxyEvent): Promise<APIGatewayProxyResult> {
  try {
    const userId = event.requestContext.authorizer?.userId as string;
    const method = event.httpMethod;

    if (method === 'GET') {
      const user = await userRepo.getById(userId);

      if (!user) {
        throw AppError.notFound('User');
      }

      return successResponse(user);
    }

    if (method === 'PATCH') {
      const body = JSON.parse(event.body ?? '{}');

      // Only allow whitelisted fields to be updated
      const updates: Record<string, unknown> = {};
      for (const [key, value] of Object.entries(body)) {
        if (ALLOWED_PROFILE_FIELDS.has(key)) {
          updates[key] = value;
        }
      }

      if (Object.keys(updates).length === 0) {
        throw AppError.badRequest('VALIDATION_ERROR' as any, 'No valid fields to update');
      }

      await userRepo.update(userId, updates as any);

      const updated = await userRepo.getById(userId);
      return successResponse(updated);
    }

    return successResponse({ message: 'Method not allowed' }, 405);
  } catch (error) {
    return handleError(error);
  }
}
