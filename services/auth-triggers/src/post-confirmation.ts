import type { PostConfirmationTriggerEvent } from 'aws-lambda';
import { UserType, TrustTier } from '@towcommand/core';
import { UserRepository } from '@towcommand/db';
import { publishEvent, EVENT_CATALOG } from '@towcommand/events';

const userRepo = new UserRepository();

export async function handler(event: PostConfirmationTriggerEvent): Promise<PostConfirmationTriggerEvent> {
  try {
    const userAttributes = event.request.userAttributes;
    const userId = userAttributes.sub;
    const email = userAttributes.email ?? '';
    const phone = userAttributes.phone_number ?? '';
    const name = userAttributes.name ?? userAttributes.email?.split('@')[0] ?? 'User';

    // Determine auth provider from trigger source
    const triggerSource = event.triggerSource;
    const isExternal = triggerSource === 'PostConfirmation_ConfirmForgotPassword'
      ? undefined
      : event.request.userAttributes.identities;

    let authProvider: 'google' | 'facebook' | 'apple' | 'phone' = 'phone';
    if (isExternal) {
      try {
        const identities = JSON.parse(isExternal);
        const providerName = identities[0]?.providerName?.toLowerCase() ?? '';
        if (providerName.includes('google')) authProvider = 'google';
        else if (providerName.includes('facebook')) authProvider = 'facebook';
        else if (providerName.includes('apple')) authProvider = 'apple';
      } catch {
        // Default to phone
      }
    }

    const now = new Date().toISOString();

    // Create user record in DynamoDB
    await userRepo.create({
      userId,
      cognitoSub: userId,
      email,
      phone,
      name,
      userType: UserType.CUSTOMER,
      trustTier: TrustTier.BASIC,
      language: 'en',
      status: 'active',
      createdAt: now,
      updatedAt: now,
    });

    // Publish user registration event (triggers welcome flow, suki tier init)
    await publishEvent(
      EVENT_CATALOG.auth.source,
      EVENT_CATALOG.auth.events.UserRegistered,
      { userId, email, phone, authProvider },
      { userId, userType: 'customer' },
    );

    return event;
  } catch (error) {
    console.error('Post-confirmation error:', error);
    // Don't throw - allow sign-up to proceed even if DB write fails
    // A reconciliation job can catch missed records
    return event;
  }
}
