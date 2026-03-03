export enum ErrorCode {
  // Auth errors
  UNAUTHORIZED = 'UNAUTHORIZED',
  FORBIDDEN = 'FORBIDDEN',
  TOKEN_EXPIRED = 'TOKEN_EXPIRED',
  ACCOUNT_SUSPENDED = 'ACCOUNT_SUSPENDED',

  // Booking errors
  BOOKING_NOT_FOUND = 'BOOKING_NOT_FOUND',
  INVALID_STATUS_TRANSITION = 'INVALID_STATUS_TRANSITION',
  BOOKING_ALREADY_MATCHED = 'BOOKING_ALREADY_MATCHED',
  NO_PROVIDERS_AVAILABLE = 'NO_PROVIDERS_AVAILABLE',
  ESTIMATE_EXPIRED = 'ESTIMATE_EXPIRED',
  CANNOT_CANCEL = 'CANNOT_CANCEL',

  // Provider errors
  PROVIDER_NOT_FOUND = 'PROVIDER_NOT_FOUND',
  PROVIDER_OFFLINE = 'PROVIDER_OFFLINE',
  PROVIDER_NOT_VERIFIED = 'PROVIDER_NOT_VERIFIED',

  // Payment errors
  PAYMENT_FAILED = 'PAYMENT_FAILED',
  INSUFFICIENT_FUNDS = 'INSUFFICIENT_FUNDS',
  REFUND_FAILED = 'REFUND_FAILED',
  HOLD_EXPIRED = 'HOLD_EXPIRED',

  // OTP errors
  INVALID_OTP = 'INVALID_OTP',
  OTP_EXPIRED = 'OTP_EXPIRED',
  OTP_MAX_ATTEMPTS = 'OTP_MAX_ATTEMPTS',

  // Validation errors
  VALIDATION_ERROR = 'VALIDATION_ERROR',
  INVALID_LOCATION = 'INVALID_LOCATION',

  // Generic
  INTERNAL_ERROR = 'INTERNAL_ERROR',
  NOT_FOUND = 'NOT_FOUND',
  RATE_LIMITED = 'RATE_LIMITED',
  SERVICE_UNAVAILABLE = 'SERVICE_UNAVAILABLE',
}

export class AppError extends Error {
  public readonly code: ErrorCode;
  public readonly statusCode: number;
  public readonly isOperational: boolean;
  public readonly details?: Record<string, unknown>;

  constructor(
    code: ErrorCode,
    message: string,
    statusCode = 500,
    isOperational = true,
    details?: Record<string, unknown>,
  ) {
    super(message);
    this.code = code;
    this.statusCode = statusCode;
    this.isOperational = isOperational;
    this.details = details;
    Object.setPrototypeOf(this, AppError.prototype);
  }

  static unauthorized(message = 'Unauthorized') {
    return new AppError(ErrorCode.UNAUTHORIZED, message, 401);
  }

  static forbidden(message = 'Forbidden') {
    return new AppError(ErrorCode.FORBIDDEN, message, 403);
  }

  static notFound(resource: string) {
    return new AppError(ErrorCode.NOT_FOUND, `${resource} not found`, 404);
  }

  static badRequest(code: ErrorCode, message: string, details?: Record<string, unknown>) {
    return new AppError(code, message, 400, true, details);
  }

  static internal(message = 'Internal server error') {
    return new AppError(ErrorCode.INTERNAL_ERROR, message, 500, false);
  }

  static rateLimited(retryAfter = 60) {
    return new AppError(ErrorCode.RATE_LIMITED, 'Too many requests', 429, true, { retryAfter });
  }

  toJSON() {
    return {
      error: this.code,
      message: this.message,
      ...(this.details ? { details: this.details } : {}),
    };
  }
}
