/**
 * Login Screen Tests — 2nd Order Logic
 *
 * 2nd order concerns:
 * - Empty email/password submission: should the UI prevent it or let
 *   the server reject? Currently no client-side validation → test that
 *   the flow doesn't crash on empty submit.
 * - Social login buttons: signInWithRedirect requires deep link setup.
 *   If it throws, the UI shouldn't crash — should show error.
 * - After successful login, the screen should trigger navigation away.
 *   If auth store isn't updated, user is stuck on login screen.
 */
import React from 'react';
import { render, fireEvent } from '@testing-library/react-native';

// Mock the stores and hooks before importing the component
jest.mock('@/hooks/useAuth', () => ({
  useAuth: () => ({
    user: null,
    isAuthenticated: false,
    isLoading: false,
    signIn: jest.fn(() => Promise.resolve()),
    signUp: jest.fn(() => Promise.resolve()),
    signOut: jest.fn(() => Promise.resolve()),
    socialSignIn: jest.fn(() => Promise.resolve()),
    checkAuth: jest.fn(() => Promise.resolve()),
  }),
}));

jest.mock('expo-router', () => ({
  router: { replace: jest.fn(), push: jest.fn(), back: jest.fn() },
  useRouter: () => ({
    push: jest.fn(),
    replace: jest.fn(),
    back: jest.fn(),
  }),
  Link: ({ children }: { children: React.ReactNode }) => children,
}));

jest.mock('react-native-safe-area-context', () => ({
  SafeAreaView: ({ children }: { children: React.ReactNode }) => children,
}));

import LoginScreen from '@/app/(auth)/login';

describe('LoginScreen', () => {
  it('renders without crash', () => {
    const { root } = render(<LoginScreen />);
    expect(root).toBeTruthy();
  });

  it('shows brand name', () => {
    const { getByText } = render(<LoginScreen />);
    expect(getByText(/TowCommand/i)).toBeTruthy();
  });

  it('has email and password inputs', () => {
    const { getByLabelText } = render(<LoginScreen />);
    // Input component sets accessibilityLabel={label}
    expect(getByLabelText('Email')).toBeTruthy();
    expect(getByLabelText('Password')).toBeTruthy();
  });

  it('has sign in button', () => {
    const { getByText } = render(<LoginScreen />);
    expect(getByText(/Sign In/)).toBeTruthy();
  });

  it('has social login buttons (Google, Facebook, Apple)', () => {
    const { getByText } = render(<LoginScreen />);
    expect(getByText(/Google/)).toBeTruthy();
    expect(getByText(/Facebook/)).toBeTruthy();
    expect(getByText(/Apple/)).toBeTruthy();
  });

  it('has sign up navigation link', () => {
    const { getByText } = render(<LoginScreen />);
    expect(getByText(/Sign Up/)).toBeTruthy();
  });

  it('email input accepts text and updates value', () => {
    const { getByLabelText, getByDisplayValue } = render(<LoginScreen />);
    const input = getByLabelText('Email');
    fireEvent.changeText(input, 'juan@test.ph');
    expect(getByDisplayValue('juan@test.ph')).toBeTruthy();
  });
});
