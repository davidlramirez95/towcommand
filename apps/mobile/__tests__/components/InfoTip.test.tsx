/**
 * InfoTip Component Tests — 2nd Order Logic
 *
 * 2nd order concerns:
 * - Icon + text must both render (missing icon = user doesn't notice the tip)
 * - Supports both string and JSX children (flexibility for bold text within tips)
 */
import React from 'react';
import { Text } from 'react-native';
import { render } from '@testing-library/react-native';
import { InfoTip } from '@/components/ui/InfoTip';

describe('InfoTip', () => {
  it('renders icon and string children', () => {
    const { getByText } = render(<InfoTip icon="💡">Save up to PHP 2,100</InfoTip>);
    expect(getByText('💡')).toBeTruthy();
    expect(getByText('Save up to PHP 2,100')).toBeTruthy();
  });

  it('renders JSX children', () => {
    const { getByText } = render(
      <InfoTip icon="🔒">
        <Text>Tamper-proof evidence</Text>
      </InfoTip>,
    );
    expect(getByText('Tamper-proof evidence')).toBeTruthy();
  });
});
