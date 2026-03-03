import { CognitoIdentityProviderClient } from '@aws-sdk/client-cognito-identity-provider';

let client: CognitoIdentityProviderClient | null = null;

export function getCognitoClient(): CognitoIdentityProviderClient {
  if (!client) {
    client = new CognitoIdentityProviderClient({
      region: process.env.COGNITO_REGION ?? 'ap-southeast-1',
    });
  }
  return client;
}

export function getUserPoolId(): string {
  return process.env.COGNITO_USER_POOL_ID ?? '';
}

export function getClientId(): string {
  return process.env.COGNITO_CLIENT_ID ?? '';
}
