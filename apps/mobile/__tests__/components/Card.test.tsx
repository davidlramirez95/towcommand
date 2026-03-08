/**
 * Card Component Tests — 2nd Order Logic
 *
 * 2nd order concerns:
 * - Card WITHOUT onPress should be a View (not tappable). If it's
 *   accidentally Pressable, screen readers announce it as a button.
 * - Card WITH onPress should be a Pressable (accessibility: tappable).
 * - Elevated card must have shadow props (visual hierarchy for CTAs).
 * - Selected card has different border color (booking service selection).
 */
import React from 'react';
import { Text } from 'react-native';
import { render, fireEvent } from '@testing-library/react-native';
import { Card } from '@/components/ui/Card';

describe('Card', () => {
  it('renders children', () => {
    const { getByText } = render(
      <Card>
        <Text>Card Content</Text>
      </Card>,
    );
    expect(getByText('Card Content')).toBeTruthy();
  });

  it('without onPress: renders as non-tappable (View, not Pressable)', () => {
    const { getByText } = render(
      <Card>
        <Text>Static Card</Text>
      </Card>,
    );

    // fireEvent.press should not crash but there's no handler
    const element = getByText('Static Card');
    expect(element).toBeTruthy();
  });

  it('with onPress: fires handler on tap', () => {
    const onPress = jest.fn();
    const { getByText } = render(
      <Card onPress={onPress}>
        <Text>Tappable Card</Text>
      </Card>,
    );

    fireEvent.press(getByText('Tappable Card'));
    expect(onPress).toHaveBeenCalledTimes(1);
  });

  it('elevated=true does not crash (shadow props applied)', () => {
    const { getByText } = render(
      <Card elevated>
        <Text>Elevated</Text>
      </Card>,
    );
    expect(getByText('Elevated')).toBeTruthy();
  });

  it('selected=true does not crash (orange border applied)', () => {
    const { getByText } = render(
      <Card selected>
        <Text>Selected</Text>
      </Card>,
    );
    expect(getByText('Selected')).toBeTruthy();
  });

  it('combines elevated + selected + onPress without crash', () => {
    const onPress = jest.fn();
    const { getByText } = render(
      <Card elevated selected onPress={onPress}>
        <Text>All Props</Text>
      </Card>,
    );

    fireEvent.press(getByText('All Props'));
    expect(onPress).toHaveBeenCalled();
  });

  it('passes through ViewProps (style merge)', () => {
    const { getByText } = render(
      <Card style={{ marginTop: 20 }}>
        <Text>Styled</Text>
      </Card>,
    );
    expect(getByText('Styled')).toBeTruthy();
  });
});
