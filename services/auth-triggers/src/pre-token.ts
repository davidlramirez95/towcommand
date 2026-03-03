import type { TokenGenerationTriggerEvent } from 'aws-lambda';
import { UserRepository, ProviderRepository } from '@towcommand/db';
import { SessionCache } from '@towcommand/cache';

const userRepo = new UserRepository();
const providerRepo = new ProviderRepository();
const sessionCache = new SessionCache();

export async function handler(event: TokenGenerationTriggerEvent): Promise<TokenGenerationTriggerEvent> {
  try {
    const userId = event.request.userAttributes.sub;

    // Fetch user from DynamoDB to get current role and trust tier
    const user = await userRepo.getById(userId);

    if (user) {
      const claims: Record<string, string> = {
        'custom:user_type': user.userType,
        'custom:trust_tier': user.trustTier,
        'custom:status': user.status,
        'custom:userId': user.userId,
      };

      // If user is a provider, inject provider-specific claims
      if (user.userType === 'provider') {
        const provider = await providerRepo.getById(userId);
        if (provider) {
          claims['custom:providerId'] = provider.providerId;
          claims['custom:provider_status'] = provider.status;
        }
      }

      event.response.claimsOverrideDetails = {
        claimsToAddOrOverride: claims,
      };

      // Cache claims for fast access in JWT authorizer
      await sessionCache.setUserClaims(userId, claims);
    } else {
      // User not yet in DynamoDB (race condition with post-confirmation)
      event.response.claimsOverrideDetails = {
        claimsToAddOrOverride: {
          'custom:user_type': 'customer',
          'custom:trust_tier': 'basic',
          'custom:userId': userId,
        },
      };
    }

    return event;
  } catch (error) {
    console.error('Pre-token error:', error);
    // Don't throw - return with default claims
    event.response.claimsOverrideDetails = {
      claimsToAddOrOverride: {
        'custom:user_type': 'customer',
        'custom:trust_tier': 'basic',
      },
    };
    return event;
  }
}
