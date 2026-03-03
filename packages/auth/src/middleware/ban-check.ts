import type { PreAuthenticationTriggerEvent } from 'aws-lambda';

export async function banCheckHandler(event: PreAuthenticationTriggerEvent): Promise<PreAuthenticationTriggerEvent> {
  // Check DynamoDB for banned/suspended status
  // In production, query user record and check status
  const userAttributes = event.request.userAttributes;

  if (userAttributes['custom:status'] === 'banned') {
    throw new Error('User account has been banned. Contact support.');
  }

  if (userAttributes['custom:status'] === 'suspended') {
    throw new Error('User account is suspended. Contact support.');
  }

  return event;
}
