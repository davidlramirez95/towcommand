import type { APIGatewayProxyEvent, APIGatewayProxyResult } from 'aws-lambda';
import { AppError, PaymentStatus } from '@towcommand/core';
import { PaymentRepository } from '@towcommand/db';
import { publishEvent, EVENT_CATALOG } from '@towcommand/events';
import { handleError, successResponse } from '../../middleware/error-handler';
import { createHmac } from 'crypto';

const paymentRepo = new PaymentRepository();

function verifyWebhookSignature(payload: string, signature: string): boolean {
  const secret = process.env.PAYMENT_WEBHOOK_SECRET;
  if (!secret) return false;
  const expected = createHmac('sha256', secret).update(payload).digest('hex');
  return expected === signature;
}

export async function handler(event: APIGatewayProxyEvent): Promise<APIGatewayProxyResult> {
  try {
    const rawBody = event.body ?? '';
    const signature = event.headers['x-webhook-signature'] ?? event.headers['X-Webhook-Signature'] ?? '';

    if (!verifyWebhookSignature(rawBody, signature)) {
      throw AppError.unauthorized('Invalid webhook signature');
    }

    const body = JSON.parse(rawBody);
    const { paymentId, status, gatewayRef } = body;

    const payment = await paymentRepo.getById(paymentId);

    if (!payment) {
      throw AppError.notFound('Payment');
    }

    const now = new Date().toISOString();

    if (status === 'success') {
      await paymentRepo.update(paymentId, {
        status: PaymentStatus.CAPTURED,
        gatewayRef,
        capturedAt: now,
      });

      await publishEvent(
        EVENT_CATALOG.payment.source,
        EVENT_CATALOG.payment.events.PaymentCompleted,
        {
          paymentId,
          bookingId: payment.bookingId,
          amount: payment.amount,
          method: payment.method,
          status: PaymentStatus.CAPTURED,
          gatewayRef,
        },
      );
    } else if (status === 'failed') {
      await paymentRepo.update(paymentId, {
        status: PaymentStatus.FAILED,
        gatewayRef,
      });

      await publishEvent(
        EVENT_CATALOG.payment.source,
        EVENT_CATALOG.payment.events.PaymentFailed,
        {
          paymentId,
          bookingId: payment.bookingId,
          amount: payment.amount,
          method: payment.method,
          reason: body.reason,
        },
      );
    }

    return successResponse({ received: true });
  } catch (error) {
    return handleError(error);
  }
}
