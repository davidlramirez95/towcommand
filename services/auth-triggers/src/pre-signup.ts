import type { PreSignUpTriggerEvent } from 'aws-lambda';

export async function handler(event: PreSignUpTriggerEvent): Promise<PreSignUpTriggerEvent> {
  try {
    const { email, phone_number: phone } = event.request.userAttributes;
    const triggerSource = event.triggerSource;

    // Auto-link social provider accounts (Google, Facebook, Apple)
    // If a user signs up with social and an account with the same email exists,
    // Cognito will link them automatically when autoConfirmUser is true
    if (triggerSource === 'PreSignUp_ExternalProvider') {
      event.response.autoConfirmUser = true;
      event.response.autoVerifiedEmail = true;
      if (phone) {
        event.response.autoVerifyPhone = true;
      }
      return event;
    }

    // For standard sign-up: auto-confirm in dev, require verification in prod
    if (process.env.AUTO_CONFIRM_ENABLED === 'true') {
      event.response.autoConfirmUser = true;
      event.response.autoVerifiedEmail = true;
    }

    return event;
  } catch (error) {
    console.error('Pre-signup error:', error);
    throw error;
  }
}
