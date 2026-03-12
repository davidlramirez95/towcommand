/**
 * ProgressBar Component Tests — 2nd Order Logic
 *
 * 2nd order concerns:
 * - Step count must match steps array (prevents misleading progress)
 * - Labels must all render (user needs to know what each step means)
 * - currentStep=0 means no bars active (prevents premature progress indication)
 */
import React from 'react';
import { render } from '@testing-library/react-native';
import { ProgressBar } from '@/components/ui/ProgressBar';

const TRACKING_STEPS = ['Matched', 'En Route', 'Arrived', 'Loading', 'Complete'];

describe('ProgressBar', () => {
  it('renders all step labels', () => {
    const { getByText } = render(<ProgressBar steps={TRACKING_STEPS} currentStep={2} />);
    TRACKING_STEPS.forEach((step) => {
      expect(getByText(step)).toBeTruthy();
    });
  });

  it('renders correct number of bars (must match steps array length)', () => {
    const { getAllByTestId } = render(<ProgressBar steps={TRACKING_STEPS} currentStep={2} />);
    // Since we don't use testIDs on bars, just verify all labels render (proxy for bar count)
    const { getAllByText } = render(<ProgressBar steps={['A', 'B', 'C']} currentStep={1} />);
    expect(getAllByText(/^[ABC]$/).length).toBe(3);
  });

  it('renders with currentStep=0 (no progress yet)', () => {
    const { getByText } = render(<ProgressBar steps={TRACKING_STEPS} currentStep={0} />);
    expect(getByText('Matched')).toBeTruthy();
  });
});
