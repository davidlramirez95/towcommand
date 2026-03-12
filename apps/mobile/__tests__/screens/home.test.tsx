/**
 * Home Screen Tests — 2nd Order Logic
 *
 * 2nd order concerns:
 * - firstName is derived from user?.fullName?.split(' ')[0].
 *   If fullName is null → split() crashes. If '' → empty greeting.
 * - SOS button must always be visible and tappable (safety-critical).
 * - AI diagnosis card must be prominent (primary user action).
 * - Quick services grid must render all 4 services.
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
    expect(getByText(/Juan/)).toBeTruthy();
  });

  it('SOS button is visible (safety-critical: must always be accessible)', () => {
    const { getByText } = render(<HomeScreen />);
    expect(getByText('SOS')).toBeTruthy();
  });

  it('AI diagnosis card is visible (primary action)', () => {
    const { getByText } = render(<HomeScreen />);
    expect(getByText(/wrong with your car/i)).toBeTruthy();
  });

  it('Quick services grid renders all 4 services', () => {
    const { getByText } = render(<HomeScreen />);
    expect(getByText('Tow')).toBeTruthy();
    expect(getByText('Fuel')).toBeTruthy();
    expect(getByText('Jumpstart')).toBeTruthy();
    expect(getByText('Mechanic')).toBeTruthy();
  });
});
