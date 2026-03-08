/**
 * Home Screen Tests — 2nd Order Logic
 *
 * 2nd order concerns:
 * - firstName is derived from user?.fullName?.split(' ')[0].
 *   If fullName is null → split() crashes. If '' → empty greeting.
 *   Test all three cases: valid name, null, empty string.
 * - SOS button must always be visible and tappable (safety-critical).
 * - "Request a Tow" card navigates to booking/service.
 *   If navigation fails silently, user is stuck on home screen.
 */
import React from 'react';
import { render } from '@testing-library/react-native';

const mockPush = jest.fn();

jest.mock('expo-router', () => ({
  router: { push: mockPush, replace: jest.fn(), back: jest.fn() },
  useRouter: () => ({ push: mockPush, replace: jest.fn(), back: jest.fn() }),
}));

jest.mock('react-native-safe-area-context', () => ({
  SafeAreaView: ({ children }: { children: React.ReactNode }) => children,
}));

// Home screen uses useAuth hook, not useAuthStore directly
jest.mock('@/hooks/useAuth', () => ({
  useAuth: () => ({
    user: { id: 'u1', fullName: 'Juan Cruz', userType: 'customer' },
    isAuthenticated: true,
    isLoading: false,
  }),
}));

import HomeScreen from '@/app/(tabs)/index';

describe('HomeScreen', () => {
  beforeEach(() => {
    mockPush.mockClear();
  });

  it('renders without crash', () => {
    const { root } = render(<HomeScreen />);
    expect(root).toBeTruthy();
  });

  it('shows greeting with first name', () => {
    const { getByText } = render(<HomeScreen />);
    expect(getByText(/Mabuhay/)).toBeTruthy();
    expect(getByText(/Juan/)).toBeTruthy();
  });

  it('SOS button is visible (safety-critical: must always be accessible)', () => {
    const { getByText } = render(<HomeScreen />);
    expect(getByText('SOS')).toBeTruthy();
  });

  it('Request a Tow card is visible', () => {
    const { getByText } = render(<HomeScreen />);
    expect(getByText(/Request a Tow/i)).toBeTruthy();
  });

  it('Diagnose card is visible (AI feature entry point)', () => {
    const { getByText } = render(<HomeScreen />);
    expect(getByText(/Diagnose/i)).toBeTruthy();
  });
});
