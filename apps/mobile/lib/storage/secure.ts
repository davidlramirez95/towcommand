import { Platform } from 'react-native';

const TOKEN_KEY = 'auth_tokens';

interface AuthTokens {
  accessToken: string;
  idToken: string;
  refreshToken: string;
}

/**
 * Secure token storage with web fallback (localStorage).
 * expo-secure-store is native-only; web uses localStorage for Playwright E2E.
 */

async function getSecureStore() {
  if (Platform.OS === 'web') {
    return {
      setItemAsync: async (key: string, value: string) => localStorage.setItem(key, value),
      getItemAsync: async (key: string) => localStorage.getItem(key),
      deleteItemAsync: async (key: string) => localStorage.removeItem(key),
    };
  }
  return await import('expo-secure-store');
}

export async function saveTokens(tokens: AuthTokens): Promise<void> {
  const store = await getSecureStore();
  await store.setItemAsync(TOKEN_KEY, JSON.stringify(tokens));
}

export async function getTokens(): Promise<AuthTokens | null> {
  const store = await getSecureStore();
  const raw = await store.getItemAsync(TOKEN_KEY);
  if (!raw) return null;
  return JSON.parse(raw) as AuthTokens;
}

export async function clearTokens(): Promise<void> {
  const store = await getSecureStore();
  await store.deleteItemAsync(TOKEN_KEY);
}
