/**
 * BackHeader Component Tests — 2nd Order Logic
 *
 * 2nd order concerns:
 * - Back button MUST fire onBack (prevents user from getting stuck on a screen)
 * - Title must be visible (prevents confusion about which screen user is on)
 * - Right slot must render when provided (prevents missing action buttons like AI badge)
 */
import React from 'react';
import { Text } from 'react-native';
import { render, fireEvent } from '@testing-library/react-native';
import { BackHeader } from '@/components/ui/BackHeader';

describe('BackHeader', () => {
  it('renders title text', () => {
    const { getByText } = render(<BackHeader title="Vehicle Details" onBack={jest.fn()} />);
    expect(getByText('Vehicle Details')).toBeTruthy();
  });

  it('calls onBack when back button is pressed', () => {
    const onBack = jest.fn();
    const { getByLabelText } = render(<BackHeader title="Test" onBack={onBack} />);

    fireEvent.press(getByLabelText('Go back'));
    expect(onBack).toHaveBeenCalledTimes(1);
  });

  it('renders right slot when provided', () => {
    const { getByText } = render(
      <BackHeader title="Test" onBack={jest.fn()} right={<Text>AI-POWERED</Text>} />,
    );
    expect(getByText('AI-POWERED')).toBeTruthy();
  });

  it('does not render right slot when not provided', () => {
    const { queryByText } = render(<BackHeader title="Test" onBack={jest.fn()} />);
    expect(queryByText('AI-POWERED')).toBeNull();
  });
});
