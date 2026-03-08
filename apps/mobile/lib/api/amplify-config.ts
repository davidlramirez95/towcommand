import { Amplify } from 'aws-amplify';

/**
 * Configure Amplify auth to use the existing Cognito User Pool.
 * Values come from environment variables set in app.json / eas.json.
 */
export function configureAmplify() {
  Amplify.configure({
    Auth: {
      Cognito: {
        userPoolId: process.env.EXPO_PUBLIC_COGNITO_USER_POOL_ID ?? '',
        userPoolClientId: process.env.EXPO_PUBLIC_COGNITO_CLIENT_ID ?? '',
        loginWith: {
          oauth: {
            domain: process.env.EXPO_PUBLIC_COGNITO_DOMAIN ?? '',
            scopes: ['openid', 'email', 'profile'],
            redirectSignIn: ['towcommand://auth/callback'],
            redirectSignOut: ['towcommand://auth/signout'],
            responseType: 'code',
            providers: ['Google', 'Facebook', 'Apple'],
          },
        },
      },
    },
  });
}
