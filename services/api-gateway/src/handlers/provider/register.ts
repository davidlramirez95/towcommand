import type { APIGatewayProxyEvent, APIGatewayProxyResult } from 'aws-lambda';
import { providerRegistrationSchema, ProviderStatus, TrustTier } from '@towcommand/core';
import { ProviderRepository } from '@towcommand/db';
import { publishEvent, EVENT_CATALOG } from '@towcommand/events';
import { handleError, successResponse } from '../../middleware/error-handler';
import { ulid } from 'ulid';

const providerRepo = new ProviderRepository();

export async function handler(event: APIGatewayProxyEvent): Promise<APIGatewayProxyResult> {
  try {
    const cognitoSub = event.requestContext.authorizer?.sub as string ?? '';
    const body = providerRegistrationSchema.parse(JSON.parse(event.body ?? '{}'));

    const providerId = `PROV-${ulid()}`;
    const now = new Date().toISOString();

    const provider = {
      providerId,
      cognitoSub,
      name: body.name,
      phone: body.phone,
      email: body.email,
      status: ProviderStatus.PENDING_VERIFICATION,
      trustTier: TrustTier.BASIC,
      truckType: body.truckType as any,
      maxWeightCapacityKg: body.maxWeightCapacityKg,
      plateNumber: body.plateNumber,
      ltoRegistration: body.ltoRegistration,
      nbiClearanceStatus: 'pending' as const,
      drugTestStatus: 'pending' as const,
      mmadAccredited: false,
      rating: 0,
      totalJobsCompleted: 0,
      acceptanceRate: 100,
      isOnline: false,
      serviceAreas: body.serviceAreas,
      createdAt: now,
      updatedAt: now,
    };

    await providerRepo.create(provider, body.serviceAreas[0] ?? 'NCR');

    await publishEvent(
      EVENT_CATALOG.provider.source,
      EVENT_CATALOG.provider.events.ProviderVerified,
      {
        providerId,
        name: body.name,
        phone: body.phone,
        status: ProviderStatus.PENDING_VERIFICATION,
      },
      { userId: providerId, userType: 'provider' },
    );

    return successResponse({ providerId, status: ProviderStatus.PENDING_VERIFICATION }, 201);
  } catch (error) {
    return handleError(error);
  }
}
