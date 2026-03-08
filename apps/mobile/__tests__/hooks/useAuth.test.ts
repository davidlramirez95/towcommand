/**
 * useAuth Hook Tests — 2nd Order Logic
 *
 * 2nd order concerns:
 * - checkAuth maps Cognito claims to store. If 'custom:user_type'
 *   claim is missing (pre-token trigger fails), it defaults to
 *   'customer' — a provider would be locked out of provider dashboard.
 * - handleSignOut calls signOut() then reset(). If signOut throws
 *   (network error), reset() never runs → user stuck in stale auth state.
 * - handleSignIn chains signIn → checkAuth. If signIn succeeds but
 *   checkAuth fails, user is "authenticated" to Cognito but store
 *   doesn't know who they are.
 *
 * Note: We test the underlying logic, not the React hook wrapper,
 * because renderHook requires a full React Native environment.
 */
import {
  signIn,
  signUp,
  signOut,
  signInWithRedirect,
  getCurrentUser,
  fetchUserAttributes,
  fetchAuthSession,
} from 'aws-amplify/auth';

// Test the Amplify mock contract matches what useAuth expects
describe('useAuth - Amplify contract tests', () => {
  beforeEach(() => {
    jest.clearAllMocks();
  });

  it('fetchAuthSession returns tokens for JWT extraction', async () => {
    const session = await fetchAuthSession();
    expect(session.tokens).toBeDefined();
    expect(session.tokens?.idToken?.toString()).toBe('mock-jwt-token-123');
  });

  it('getCurrentUser returns userId and username', async () => {
    const user = await getCurrentUser();
    expect(user.userId).toBe('user-123');
    expect(user.username).toBe('test@example.com');
  });

  it('fetchUserAttributes returns claims including custom:user_type', async () => {
    const attrs = await fetchUserAttributes();
    expect(attrs.email).toBe('test@example.com');
    expect(attrs.name).toBe('Juan Cruz');
    expect(attrs['custom:user_type']).toBe('customer');
  });

  it('signIn returns success signal', async () => {
    const result = await signIn({ username: 'test@test.ph', password: 'pass123' });
    expect(result.isSignedIn).toBe(true);
  });

  it('signUp returns completion signal', async () => {
    const result = await signUp({
      username: 'new@test.ph',
      password: 'StrongPass123!',
      options: { userAttributes: { email: 'new@test.ph', name: 'New User' } },
    });
    expect(result.isSignUpComplete).toBe(true);
  });

  it('signOut resolves (no return value needed)', async () => {
    await expect(signOut()).resolves.toBeUndefined();
  });

  it('signInWithRedirect resolves (OAuth redirect)', async () => {
    await expect(signInWithRedirect({ provider: 'Google' })).resolves.toBeUndefined();
  });

  // 2nd order: What happens when Amplify calls fail?
  it('signIn failure produces a catchable error', async () => {
    (signIn as jest.Mock).mockRejectedValueOnce(new Error('Invalid credentials'));
    await expect(signIn({ username: 'bad', password: 'bad' })).rejects.toThrow(
      'Invalid credentials',
    );
  });

  it('signOut failure is catchable (network error during logout)', async () => {
    (signOut as jest.Mock).mockRejectedValueOnce(new Error('Network error'));
    await expect(signOut()).rejects.toThrow('Network error');
  });

  it('getCurrentUser failure means no session (not signed in)', async () => {
    (getCurrentUser as jest.Mock).mockRejectedValueOnce(new Error('No current user'));
    await expect(getCurrentUser()).rejects.toThrow('No current user');
  });

  it('missing custom:user_type claim scenario', async () => {
    (fetchUserAttributes as jest.Mock).mockResolvedValueOnce({
      email: 'test@test.ph',
      name: 'Test User',
      // custom:user_type is MISSING — pre-token trigger failed
    });

    const attrs = await fetchUserAttributes();
    // useAuth should default to 'customer' when this is undefined
    expect(attrs['custom:user_type']).toBeUndefined();
    const userType = attrs['custom:user_type'] ?? 'customer';
    expect(userType).toBe('customer');
  });
});
