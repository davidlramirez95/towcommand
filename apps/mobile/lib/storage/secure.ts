import * as SecureStore from 'expo-secure-store';

const TOKEN_KEY = 'auth_tokens';

interface AuthTokens {
  accessToken: string;
  idToken: string;
  refreshToken: string;
}

export async function saveTokens(tokens: AuthTokens): Promise<void> {
  await SecureStore.setItemAsync(TOKEN_KEY, JSON.stringify(tokens));
}

export async function getTokens(): Promise<AuthTokens | null> {
  const raw = await SecureStore.getItemAsync(TOKEN_KEY);
  if (!raw) return null;
  return JSON.parse(raw) as AuthTokens;
}

export async function clearTokens(): Promise<void> {
  await SecureStore.deleteItemAsync(TOKEN_KEY);
}
