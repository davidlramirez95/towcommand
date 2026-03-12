/**
 * Button Component Tests — 2nd Order Logic
 *
 * 2nd order concerns:
 * - Loading state MUST block onPress to prevent double-submission
 *   (e.g., user taps "Pay PHP 2,500" twice → double charge)
 * - Disabled state must have visual indicator (opacity) AND block press
 * - Ghost variant on dark backgrounds: must still have visible text
 * - All 5 variants must render without crashing (regression guard)
 */
import React from 'react';
import { render, fireEvent } from '@testing-library/react-native';
import { Button } from '@/components/ui/Button';

describe('Button', () => {
  it('renders children text', () => {
    const { getByText } = render(<Button onPress={jest.fn()}>Request a Tow</Button>);
    expect(getByText('Request a Tow')).toBeTruthy();
  });

  it('calls onPress when tapped', () => {
    const onPress = jest.fn();
    const { getByText } = render(<Button onPress={onPress}>Tap Me</Button>);

    fireEvent.press(getByText('Tap Me'));
    expect(onPress).toHaveBeenCalledTimes(1);
  });

  it('loading=true BLOCKS onPress (prevents double-charge on payment)', () => {
    const onPress = jest.fn();
    // When loading, the Button replaces children text with ActivityIndicator.
    // Use accessibilityLabel (defaults to children prop) to find the button.
    const { getByLabelText } = render(
      <Button onPress={onPress} loading>
        Pay PHP 2,500
      </Button>,
    );

    fireEvent.press(getByLabelText('Pay PHP 2,500'));
    expect(onPress).not.toHaveBeenCalled();
  });

  it('loading=true replaces text with ActivityIndicator', () => {
    const { queryByText } = render(
      <Button onPress={jest.fn()} loading>
        Submit
      </Button>,
    );

    // When loading, text is replaced by ActivityIndicator — text not rendered
    expect(queryByText('Submit')).toBeNull();
  });

  it('disabled=true blocks onPress', () => {
    const onPress = jest.fn();
    const { getByText } = render(
      <Button onPress={onPress} disabled>
        Disabled
      </Button>,
    );

    fireEvent.press(getByText('Disabled'));
    expect(onPress).not.toHaveBeenCalled();
  });

  it('all 5 variants render without crash (regression guard)', () => {
    const variants = ['primary', 'secondary', 'teal', 'danger', 'ghost'] as const;

    for (const variant of variants) {
      const { unmount } = render(
        <Button onPress={jest.fn()} variant={variant}>
          {variant}
        </Button>,
      );
      unmount();
    }
  });

  it('small variant renders without crash', () => {
    const { getByText } = render(
      <Button onPress={jest.fn()} small>
        Small
      </Button>,
    );
    expect(getByText('Small')).toBeTruthy();
  });

  it('fullWidth variant renders without crash', () => {
    const { getByText } = render(
      <Button onPress={jest.fn()} fullWidth>
        Full Width
      </Button>,
    );
    expect(getByText('Full Width')).toBeTruthy();
  });

  it('accessibilityLabel is set when provided', () => {
    const { getByLabelText } = render(
      <Button onPress={jest.fn()} accessibilityLabel="Submit booking request">
        Submit
      </Button>,
    );
    expect(getByLabelText('Submit booking request')).toBeTruthy();
  });
});
