/**
 * Service Selection Screen Tests — 2nd Order Logic
 *
 * 2nd order concerns:
 * - All 8 service types must render. If one is missing, customers
 *   can't select that service and lose revenue.
 * - Selecting a service should navigate to the next booking step.
 *   If the navigation callback is broken, the booking flow is stuck.
 * - Service selection must be responsive (no accidental double-tap
 *   selecting wrong service).
 */
import React from 'react';
import { render, fireEvent } from '@testing-library/react-native';

jest.mock('react-native-safe-area-context', () => ({
  SafeAreaView: ({ children }: { children: React.ReactNode }) => children,
}));

// Use global expo-router mock from __mocks__/expo-router.js
// Per-test jest.mock('expo-router', factory) doesn't override moduleNameMapper
// eslint-disable-next-line @typescript-eslint/no-var-requires
const { router } = require('expo-router');

import ServiceScreen from '@/app/booking/service';

describe('ServiceScreen', () => {
  beforeEach(() => {
    router.push.mockClear();
  });

  it('renders without crash', () => {
    const { root } = render(<ServiceScreen />);
    expect(root).toBeTruthy();
  });

  it('shows screen title', () => {
    const { getByText } = render(<ServiceScreen />);
    expect(getByText('Select Service')).toBeTruthy();
  });

  it('displays service type options', () => {
    const { getByText } = render(<ServiceScreen />);
    expect(getByText('Flatbed Towing')).toBeTruthy();
  });

  it('displays all 8 service types (complete catalog)', () => {
    const { getByText } = render(<ServiceScreen />);
    expect(getByText('Flatbed Towing')).toBeTruthy();
    expect(getByText('Wheel Lift Towing')).toBeTruthy();
    expect(getByText('Motorcycle Towing')).toBeTruthy();
    expect(getByText('Jumpstart')).toBeTruthy();
    expect(getByText('Tire Change')).toBeTruthy();
    expect(getByText('Lockout Service')).toBeTruthy();
    expect(getByText('Fuel Delivery')).toBeTruthy();
    expect(getByText('Winch Recovery')).toBeTruthy();
  });

  it('selecting a service navigates to vehicle screen with service param', () => {
    const { getByText } = render(<ServiceScreen />);
    fireEvent.press(getByText('Flatbed Towing'));

    expect(router.push).toHaveBeenCalledTimes(1);
    expect(router.push).toHaveBeenCalledWith(
      expect.objectContaining({ pathname: '/booking/vehicle', params: { service: 'FLATBED_TOWING' } }),
    );
  });
});
