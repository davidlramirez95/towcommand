import type { APIGatewayProxyEvent } from 'aws-lambda';
import pino from 'pino';
import { ulid } from 'ulid';

export const logger = pino({
  level: process.env.LOG_LEVEL || 'info',
  transport: {
    target: 'pino-pretty',
    options: {
      colorize: true,
    },
  },
});

export function getCorrelationId(event: APIGatewayProxyEvent): string {
  return (event.headers?.['x-correlation-id'] as string) || ulid();
}

export function createRequestLogger(correlationId: string) {
  return logger.child({ correlationId });
}
