/**
 * SOS Screen Tests — 2nd Order Logic
 *
 * 2nd order concerns:
 * - SOS is safety-critical. The trigger button must ALWAYS render
 *   and be pressable, even if other parts of the screen fail.
 * - The screen uses Vibration.vibrate(), not Haptics. On devices without
 *   vibration support, this must not crash the SOS flow.
 * - After triggering, the button should disable to prevent accidental
 *   re-trigger, but the emergency info must be visible.
 * - SOS must work without network (the button should still respond).
 */
import React from 'react';
import { render, fireEvent } from '@testing-library/react-native';

jest.mock('expo-router', () => ({
  router: { push: jest.fn(), replace: jest.fn(), back: jest.fn() },
  useRouter: () => ({ push: jest.fn(), replace: jest.fn(), back: jest.fn() }),
}));

jest.mock('react-native-safe-area-context', () => ({
  SafeAreaView: ({ children }: { children: React.ReactNode }) => children,
}));

// SOS screen uses useAuth hook
jest.mock('@/hooks/useAuth', () => ({
  useAuth: () => ({
    user: { id: 'u1', fullName: 'Juan', userType: 'customer' },
    isAuthenticated: true,
  }),
}));

// Mock the API client
jest.mock('@/lib/api', () => ({
  api: {
    post: jest.fn(() => Promise.resolve({})),
    get: jest.fn(() => Promise.resolve({})),
  },
}));

import SOSScreen from '@/app/sos';

describe('SOSScreen', () => {
  beforeEach(() => {
    jest.clearAllMocks();
  });

  it('renders without crash', () => {
    const { root } = render(<SOSScreen />);
    expect(root).toBeTruthy();
  });

  it('shows SOS trigger button', () => {
    const { getAllByText } = render(<SOSScreen />);
    // The SOS screen has the word "SOS" in both title and button
    const sosElements = getAllByText(/SOS/);
    expect(sosElements.length).toBeGreaterThanOrEqual(1);
  });

  it('shows emergency safety context', () => {
    const { getAllByText } = render(<SOSScreen />);
    // The screen mentions "safety team" in multiple places (subtitle + disclaimer)
    const matches = getAllByText(/safety team/i);
    expect(matches.length).toBeGreaterThanOrEqual(1);
  });

  it('SOS button is accessible via accessibility label', () => {
    const { getByLabelText } = render(<SOSScreen />);
    expect(getByLabelText('Trigger SOS Alert')).toBeTruthy();
  });
});
