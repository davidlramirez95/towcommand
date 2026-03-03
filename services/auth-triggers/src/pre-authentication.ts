import type { PreAuthenticationTriggerEvent } from 'aws-lambda';
import { UserRepository } from '@towcommand/db';

const userRepo = new UserRepository();

export async function handler(event: PreAuthenticationTriggerEvent): Promise<PreAuthenticationTriggerEvent> {
  try {
    const userId = event.request.userAttributes.sub;

    // Check user status in DynamoDB (more up-to-date than Cognito attributes)
    const user = await userRepo.getById(userId);

    if (user) {
      if (user.status === 'banned') {
        throw new Error(
          'Your account has been permanently banned due to policy violations. ' +
          'Contact support@towcommand.ph for assistance.',
        );
      }

      if (user.status === 'suspended') {
        throw new Error(
          'Your account is temporarily suspended. ' +
          'Contact support@towcommand.ph for assistance.',
        );
      }
    }

    // Fallback: check Cognito custom attributes
    const customStatus = event.request.userAttributes['custom:status'];
    if (customStatus === 'banned') {
      throw new Error('User account has been banned. Contact support.');
    }
    if (customStatus === 'suspended') {
      throw new Error('User account is suspended. Contact support.');
    }

    return event;
  } catch (error) {
    if (error instanceof Error && (
      error.message.includes('banned') || error.message.includes('suspended')
    )) {
      throw error;
    }
    console.error('Pre-authentication error:', error);
    // Don't block auth for unexpected errors
    return event;
  }
}
