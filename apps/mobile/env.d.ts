declare const process: {
  env: {
    EXPO_PUBLIC_API_URL?: string;
    EXPO_PUBLIC_WS_URL?: string;
    EXPO_PUBLIC_COGNITO_USER_POOL_ID?: string;
    EXPO_PUBLIC_COGNITO_CLIENT_ID?: string;
    EXPO_PUBLIC_COGNITO_DOMAIN?: string;
    [key: string]: string | undefined;
  };
};
