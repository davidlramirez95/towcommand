/**
 * Storage Adapter Tests — 2nd Order Logic
 *
 * 2nd order concerns:
 * - MMKV StateStorage adapter is the bridge between Zustand persist
 *   middleware and MMKV. If getItem returns undefined instead of null
 *   for missing keys, Zustand's rehydration logic behaves differently.
 * - SecureStore wraps expo-secure-store. If getTokens receives malformed
 *   JSON, it currently throws (no try/catch) — callers must handle this.
 */
import { mmkvStorage } from '@/lib/storage/mmkv';
import { saveTokens, getTokens, clearTokens } from '@/lib/storage/secure';
import * as SecureStore from 'expo-secure-store';

describe('mmkvStorage (Zustand adapter)', () => {
  it('setItem + getItem roundtrips correctly', () => {
    mmkvStorage.setItem('test-key', JSON.stringify({ foo: 'bar' }));
    const result = mmkvStorage.getItem('test-key');
    expect(JSON.parse(result!)).toEqual({ foo: 'bar' });
  });

  it('getItem returns null for missing keys (Zustand expects null, not undefined)', () => {
    const result = mmkvStorage.getItem('nonexistent-key');
    expect(result).toBeNull();
  });

  it('removeItem deletes the key', () => {
    mmkvStorage.setItem('to-delete', 'value');
    mmkvStorage.removeItem('to-delete');
    const result = mmkvStorage.getItem('to-delete');
    expect(result).toBeNull();
  });
});

describe('secureStore (token storage)', () => {
  beforeEach(() => {
    jest.clearAllMocks();
  });

  it('saveTokens + getTokens roundtrips', async () => {
    const tokens = { accessToken: 'abc', refreshToken: 'xyz', idToken: 'ijk' };
    await saveTokens(tokens);
    const result = await getTokens();

    expect(result).toEqual(tokens);
    expect(SecureStore.setItemAsync).toHaveBeenCalledWith(
      'auth_tokens',
      JSON.stringify(tokens),
    );
  });

  it('getTokens returns null when no tokens stored', async () => {
    (SecureStore.getItemAsync as jest.Mock).mockResolvedValueOnce(null);
    const result = await getTokens();
    expect(result).toBeNull();
  });

  it('clearTokens removes the key', async () => {
    await clearTokens();
    expect(SecureStore.deleteItemAsync).toHaveBeenCalledWith('auth_tokens');
  });

  it('getTokens throws on corrupt JSON (callers must handle — no try/catch in impl)', async () => {
    // 2nd order: getTokens does JSON.parse without try/catch.
    // If SecureStore returns corrupt data (e.g., partial write during crash),
    // this will throw SyntaxError. Callers (useAuth) must wrap in try/catch.
    (SecureStore.getItemAsync as jest.Mock).mockResolvedValueOnce('{invalid json');
    await expect(getTokens()).rejects.toThrow(SyntaxError);
  });
});
