/**
 * Auth Store Tests — 2nd Order Logic
 *
 * 2nd order concern: Auth store is the single source of truth for "who is
 * logged in." Every screen checks isAuthenticated. Every API call reads
 * the user. If reset() leaves stale data, the NEXT user session sees
 * PREVIOUS user's data. If setUser(null) crashes, logout flow breaks.
 */
import { useAuthStore } from '@/stores/auth';

// Reset store between tests to prevent state leakage
beforeEach(() => {
  useAuthStore.getState().reset();
});

describe('auth store', () => {
  it('initial state is logged out after reset', () => {
    const state = useAuthStore.getState();
    expect(state.user).toBeNull();
    expect(state.isAuthenticated).toBe(false);
    // reset() sets isLoading to false (app needs explicit setLoading(true) for loading state)
    expect(state.isLoading).toBe(false);
  });

  it('setUser marks authenticated and stores user data', () => {
    useAuthStore.getState().setUser({
      id: 'user-1',
      email: 'juan@test.ph',
      phone: '+639171234567',
      fullName: 'Juan Cruz',
      userType: 'customer',
    });

    const state = useAuthStore.getState();
    expect(state.isAuthenticated).toBe(true);
    expect(state.user?.id).toBe('user-1');
    expect(state.user?.fullName).toBe('Juan Cruz');
    // setUser also sets isLoading to false (auth check complete)
    expect(state.isLoading).toBe(false);
  });

  it('setUser(null) transitions to logged-out (logout flow)', () => {
    // Login first
    useAuthStore.getState().setUser({
      id: 'user-1',
      email: 'test@test.ph',
      phone: '+639170000000',
      fullName: 'Test',
      userType: 'customer',
    });
    expect(useAuthStore.getState().isAuthenticated).toBe(true);

    // Logout
    useAuthStore.getState().setUser(null);
    expect(useAuthStore.getState().isAuthenticated).toBe(false);
    expect(useAuthStore.getState().user).toBeNull();
  });

  it('reset() clears ALL fields — prevents data leaking to next session', () => {
    // Simulate full session
    useAuthStore.getState().setUser({
      id: 'user-old',
      email: 'old@test.ph',
      phone: '+639170000000',
      fullName: 'Old User',
      userType: 'provider',
      avatarUrl: 'https://example.com/avatar.jpg',
    });
    useAuthStore.getState().setLoading(false);

    // Reset (called during signOut)
    useAuthStore.getState().reset();

    const state = useAuthStore.getState();
    expect(state.user).toBeNull();
    expect(state.isAuthenticated).toBe(false);
    // reset() sets isLoading to false per the store definition
    expect(state.isLoading).toBe(false);
  });

  it('setLoading toggles loading state independently of auth', () => {
    useAuthStore.getState().setLoading(false);
    expect(useAuthStore.getState().isLoading).toBe(false);

    useAuthStore.getState().setLoading(true);
    expect(useAuthStore.getState().isLoading).toBe(true);
  });

  it('multiple rapid setUser calls settle on last value (race condition safety)', () => {
    const store = useAuthStore.getState();
    store.setUser({ id: 'a', email: 'a@test.ph', phone: '+639170000001', fullName: 'A', userType: 'customer' });
    store.setUser({ id: 'b', email: 'b@test.ph', phone: '+639170000002', fullName: 'B', userType: 'provider' });
    store.setUser({ id: 'c', email: 'c@test.ph', phone: '+639170000003', fullName: 'C', userType: 'customer' });

    expect(useAuthStore.getState().user?.id).toBe('c');
    expect(useAuthStore.getState().isAuthenticated).toBe(true);
  });

  it('provider user type is preserved (provider dashboard routing depends on this)', () => {
    useAuthStore.getState().setUser({
      id: 'prov-1',
      email: 'driver@test.ph',
      phone: '+639170000000',
      fullName: 'Provider Juan',
      userType: 'provider',
    });

    expect(useAuthStore.getState().user?.userType).toBe('provider');
  });
});
