import type { APIGatewayProxyResult } from 'aws-lambda';
import { AppError, ValidationError } from '@towcommand/core';
import { ZodError } from 'zod';

export function handleError(error: unknown): APIGatewayProxyResult {
  if (error instanceof ValidationError) {
    return jsonResponse(400, error.toJSON());
  }

  if (error instanceof ZodError) {
    const ve = new ValidationError(error);
    return jsonResponse(400, ve.toJSON());
  }

  if (error instanceof AppError) {
    return jsonResponse(error.statusCode, error.toJSON());
  }

  console.error('Unhandled error:', error);
  return jsonResponse(500, { error: 'INTERNAL_ERROR', message: 'Internal server error' });
}

export function jsonResponse(statusCode: number, body: unknown): APIGatewayProxyResult {
  return {
    statusCode,
    headers: {
      'Content-Type': 'application/json',
      'Access-Control-Allow-Origin': '*',
      'Access-Control-Allow-Headers': 'Content-Type,Authorization',
      'Access-Control-Allow-Methods': 'GET,POST,PATCH,DELETE,OPTIONS',
    },
    body: JSON.stringify(body),
  };
}

export function successResponse(data: unknown, statusCode = 200): APIGatewayProxyResult {
  return jsonResponse(statusCode, { success: true, data });
}
