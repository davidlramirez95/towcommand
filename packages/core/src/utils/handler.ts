/**
 * Base handler pattern for TowCommand Lambda functions
 * Routes requests based on HANDLER_TYPE environment variable
 * Pattern adapted from gutguard-ai handler routing
 */

import type { APIGatewayProxyEvent, APIGatewayProxyResult } from 'aws-lambda';
import { ResponseBuilder, errors } from './response.js';
import { Logger } from './logger.js';

export type HandlerFunction = (event: APIGatewayProxyEvent) => Promise<APIGatewayProxyResult>;

export interface HandlerMap {
  [handlerType: string]: HandlerFunction;
}

/**
 * Base class for Lambda handlers that route based on HANDLER_TYPE env var.
 * Each service extends this and registers its handler functions.
 *
 * Usage:
 * ```typescript
 * class BookingHandler extends BaseHandler {
 *   constructor() {
 *     super('booking');
 *     this.register('create', this.createBooking.bind(this));
 *     this.register('cancel', this.cancelBooking.bind(this));
 *     this.register('get', this.getBooking.bind(this));
 *   }
 * }
 * ```
 */
export abstract class BaseHandler {
  protected readonly logger: Logger;
  private readonly handlers: HandlerMap = {};
  private readonly serviceName: string;

  constructor(serviceName: string, logger?: Logger) {
    this.serviceName = serviceName;
    this.logger = logger || Logger.getInstance(serviceName);
  }

  /** Register a handler for a specific HANDLER_TYPE value */
  protected register(handlerType: string, fn: HandlerFunction): void {
    this.handlers[handlerType] = fn;
  }

  /** Parse JSON body from event, returning empty object on failure */
  protected parseBody(event: APIGatewayProxyEvent): Record<string, unknown> {
    if (!event.body) return {};
    try {
      return JSON.parse(event.body);
    } catch {
      return {};
    }
  }

  /** Extract user context from authorizer */
  protected getUserContext(event: APIGatewayProxyEvent): {
    userId: string;
    email: string;
    userType: string;
    groups: string[];
  } {
    const auth = event.requestContext?.authorizer || {};
    let groups: string[] = [];
    try {
      groups = auth.groups ? JSON.parse(auth.groups as string) : [];
    } catch {
      groups = [];
    }

    return {
      userId: (auth.userId as string) || '',
      email: (auth.email as string) || '',
      userType: (auth.userType as string) || 'CUSTOMER',
      groups,
    };
  }

  /** Main handler entry point */
  public async handle(event: APIGatewayProxyEvent): Promise<APIGatewayProxyResult> {
    const handlerType = process.env.HANDLER_TYPE;

    this.logger.debug('Request received', {
      handlerType,
      path: event.path,
      method: event.httpMethod,
      service: this.serviceName,
    });

    if (!handlerType) {
      this.logger.error('HANDLER_TYPE not set');
      return errors.serverError('Handler configuration error');
    }

    const fn = this.handlers[handlerType];
    if (!fn) {
      this.logger.error('Unknown handler type', { handlerType, available: Object.keys(this.handlers) });
      return errors.badRequest(`Unknown handler type: ${handlerType}`);
    }

    try {
      return await fn(event);
    } catch (err) {
      const error = err as Error;
      this.logger.error('Handler error', {
        handlerType,
        error: error.message,
        stack: error.stack,
      });
      return errors.serverError('An unexpected error occurred');
    }
  }
}
