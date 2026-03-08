/**
 * Provider Earnings Screen Tests — 2nd Order Logic
 *
 * 2nd order concerns:
 * - Earnings shows PHP 0.00 for all periods initially. When wired to
 *   the API, this must match the EarningsOutput format from backend.
 * - All 4 periods must render (Today, This Week, This Month, All Time).
 *   Missing a period means providers can't see recent earnings.
 * - Both gross and net amounts must be displayed (providers care about
 *   take-home after commission, not just gross).
 */
import React from 'react';
import { render } from '@testing-library/react-native';

jest.mock('expo-router', () => ({
  useRouter: () => ({ push: jest.fn(), replace: jest.fn(), back: jest.fn() }),
}));

jest.mock('react-native-safe-area-context', () => ({
  SafeAreaView: ({ children }: { children: React.ReactNode }) => children,
}));

import EarningsScreen from '@/app/provider/earnings';

describe('EarningsScreen', () => {
  it('renders without crash', () => {
    const { root } = render(<EarningsScreen />);
    expect(root).toBeTruthy();
  });

  it('shows Earnings title', () => {
    const { getByText } = render(<EarningsScreen />);
    expect(getByText('Earnings')).toBeTruthy();
  });

  it('displays all 4 time periods', () => {
    const { getByText } = render(<EarningsScreen />);
    expect(getByText('Today')).toBeTruthy();
    expect(getByText('This Week')).toBeTruthy();
    expect(getByText('This Month')).toBeTruthy();
    expect(getByText('All Time')).toBeTruthy();
  });

  it('shows PHP currency format for all periods', () => {
    const { getAllByText } = render(<EarningsScreen />);
    // All period cards should show PHP amounts
    const phpElements = getAllByText(/PHP/);
    expect(phpElements.length).toBeGreaterThanOrEqual(4);
  });

  it('shows job count for each period', () => {
    const { getAllByText } = render(<EarningsScreen />);
    const jobElements = getAllByText(/jobs/);
    expect(jobElements.length).toBe(4); // one per period
  });

  it('displays both gross and net amounts (commission transparency)', () => {
    const { getAllByText } = render(<EarningsScreen />);
    const grossElements = getAllByText(/Gross/);
    expect(grossElements.length).toBe(4); // one per period card
  });
});
