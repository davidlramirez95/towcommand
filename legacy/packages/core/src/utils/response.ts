/**
 * HTTP Response utilities for TowCommand Lambda functions
 * Provides consistent response formatting with CORS headers
 * Pattern adapted from gutguard-ai ResponseBuilder
 */

import type { APIGatewayProxyResult } from 'aws-lambda';

export interface ErrorBody {
  error: {
    code: string;
    message: string;
    requestId?: string;
  };
}

export type ErrorCode =
  | 'BAD_REQUEST'
  | 'UNAUTHORIZED'
  | 'FORBIDDEN'
  | 'NOT_FOUND'
  | 'CONFLICT'
  | 'TOO_MANY_REQUESTS'
  | 'SERVICE_UNAVAILABLE'
  | 'INTERNAL_ERROR'
  | 'VALIDATION_ERROR'
  | 'BOOKING_ERROR'
  | 'PAYMENT_ERROR'
  | 'PROVIDER_UNAVAILABLE';

export class ResponseBuilder {
  private static readonly CORS_HEADERS: Record<string, string> = {
    'Content-Type': 'application/json',
    'Access-Control-Allow-Origin': '*',
    'Access-Control-Allow-Headers': 'Content-Type,Authorization,X-Request-Id',
    'Access-Control-Allow-Methods': 'GET,POST,PUT,PATCH,DELETE,OPTIONS',
  };

  public static success<T>(data: T, statusCode: number = 200): APIGatewayProxyResult {
    return {
      statusCode,
      headers: ResponseBuilder.CORS_HEADERS,
      body: JSON.stringify({ success: true, data }),
    };
  }

  public static created<T>(data: T): APIGatewayProxyResult {
    return ResponseBuilder.success(data, 201);
  }

  public static noContent(): APIGatewayProxyResult {
    return {
      statusCode: 204,
      headers: ResponseBuilder.CORS_HEADERS,
      body: '',
    };
  }

  public static error(
    message: string,
    statusCode: number = 400,
    errorCode: ErrorCode = 'BAD_REQUEST'
  ): APIGatewayProxyResult {
    const body: ErrorBody = {
      error: {
        code: errorCode,
        message,
      },
    };
    return {
      statusCode,
      headers: ResponseBuilder.CORS_HEADERS,
      body: JSON.stringify(body),
    };
  }

  public static badRequest(message: string): APIGatewayProxyResult {
    return ResponseBuilder.error(message, 400, 'BAD_REQUEST');
  }

  public static validationError(message: string): APIGatewayProxyResult {
    return ResponseBuilder.error(message, 400, 'VALIDATION_ERROR');
  }

  public static unauthorized(message: string = 'Unauthorized'): APIGatewayProxyResult {
    return ResponseBuilder.error(message, 401, 'UNAUTHORIZED');
  }

  public static forbidden(message: string = 'Forbidden'): APIGatewayProxyResult {
    return ResponseBuilder.error(message, 403, 'FORBIDDEN');
  }

  public static notFound(message: string = 'Resource not found'): APIGatewayProxyResult {
    return ResponseBuilder.error(message, 404, 'NOT_FOUND');
  }

  public static conflict(message: string): APIGatewayProxyResult {
    return ResponseBuilder.error(message, 409, 'CONFLICT');
  }

  public static tooManyRequests(message: string = 'Rate limit exceeded. Please try again later.'): APIGatewayProxyResult {
    return ResponseBuilder.error(message, 429, 'TOO_MANY_REQUESTS');
  }

  public static serviceUnavailable(message: string = 'Service temporarily unavailable'): APIGatewayProxyResult {
    return ResponseBuilder.error(message, 503, 'SERVICE_UNAVAILABLE');
  }

  public static serverError(message: string = 'An unexpected error occurred'): APIGatewayProxyResult {
    return ResponseBuilder.error(message, 500, 'INTERNAL_ERROR');
  }

  // TowCommand-specific error helpers
  public static bookingError(message: string): APIGatewayProxyResult {
    return ResponseBuilder.error(message, 422, 'BOOKING_ERROR');
  }

  public static paymentError(message: string): APIGatewayProxyResult {
    return ResponseBuilder.error(message, 422, 'PAYMENT_ERROR');
  }

  public static providerUnavailable(message: string = 'No providers available in your area'): APIGatewayProxyResult {
    return ResponseBuilder.error(message, 503, 'PROVIDER_UNAVAILABLE');
  }
}

// Export convenience functions
export const success = ResponseBuilder.success;
export const created = ResponseBuilder.created;
export const error = ResponseBuilder.error;

export const errors = {
  badRequest: ResponseBuilder.badRequest,
  validationError: ResponseBuilder.validationError,
  unauthorized: ResponseBuilder.unauthorized,
  forbidden: ResponseBuilder.forbidden,
  notFound: ResponseBuilder.notFound,
  conflict: ResponseBuilder.conflict,
  tooManyRequests: ResponseBuilder.tooManyRequests,
  serviceUnavailable: ResponseBuilder.serviceUnavailable,
  serverError: ResponseBuilder.serverError,
  bookingError: ResponseBuilder.bookingError,
  paymentError: ResponseBuilder.paymentError,
  providerUnavailable: ResponseBuilder.providerUnavailable,
};

export default { success, created, error, errors, ResponseBuilder };
