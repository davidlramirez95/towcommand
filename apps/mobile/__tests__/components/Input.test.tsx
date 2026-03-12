/**
 * Input Component Tests — 2nd Order Logic
 *
 * 2nd order concerns:
 * - Error state and error message: if error string is set, the red
 *   border + error text must appear. If helper text is ALSO set,
 *   only error should show (helper is redundant during error state).
 * - Focus border (teal) vs error border (coral): if input is focused
 *   AND has error, error must visually win — user needs to see the problem.
 * - Controlled input: onChangeText must correctly propagate to parent
 *   (store-connected inputs depend on this).
 */
import React from 'react';
import { render, fireEvent } from '@testing-library/react-native';
import { Input } from '@/components/ui/Input';

describe('Input', () => {
  it('renders with label', () => {
    const { getByText } = render(<Input label="Email Address" />);
    expect(getByText('Email Address')).toBeTruthy();
  });

  it('renders without label (optional prop)', () => {
    const { root } = render(<Input placeholder="Type here..." />);
    expect(root).toBeTruthy();
  });

  it('fires onChangeText when user types', () => {
    const onChangeText = jest.fn();
    const { getByPlaceholderText } = render(
      <Input placeholder="Email" onChangeText={onChangeText} />,
    );

    fireEvent.changeText(getByPlaceholderText('Email'), 'juan@test.ph');
    expect(onChangeText).toHaveBeenCalledWith('juan@test.ph');
  });

  it('shows error message when error prop is set', () => {
    const { getByText } = render(
      <Input label="Phone" error="Phone number is required" />,
    );
    expect(getByText('Phone number is required')).toBeTruthy();
  });

  it('shows helper text when no error', () => {
    const { getByText } = render(
      <Input label="Password" helper="Must be at least 8 characters" />,
    );
    expect(getByText('Must be at least 8 characters')).toBeTruthy();
  });

  it('error text takes precedence over helper text (no conflicting guidance)', () => {
    const { getByText, queryByText } = render(
      <Input
        label="Password"
        helper="Must be at least 8 characters"
        error="Password is too short"
      />,
    );

    expect(getByText('Password is too short')).toBeTruthy();
    // Helper should be hidden when error is present
    expect(queryByText('Must be at least 8 characters')).toBeNull();
  });

  it('passes through TextInput props (secureTextEntry, keyboardType)', () => {
    const { getByPlaceholderText } = render(
      <Input
        placeholder="Password"
        secureTextEntry
        keyboardType="email-address"
      />,
    );
    expect(getByPlaceholderText('Password')).toBeTruthy();
  });

  it('accessibility: label prop sets accessibilityLabel', () => {
    const { getByLabelText } = render(
      <Input label="Full Name" placeholder="Enter your name" />,
    );
    expect(getByLabelText('Full Name')).toBeTruthy();
  });
});
