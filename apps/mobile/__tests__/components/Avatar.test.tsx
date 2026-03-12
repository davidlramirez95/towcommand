/**
 * Avatar Component Tests — 2nd Order Logic
 *
 * 2nd order concerns:
 * - Initials must extract correctly from any name format (single, two, three words)
 * - Max 2 initials (prevents overflow in small circle)
 * - Empty name must not crash (defensive against bad API data)
 */
import React from 'react';
import { render } from '@testing-library/react-native';
import { Avatar } from '@/components/ui/Avatar';

describe('Avatar', () => {
  it('renders two initials for two-word name', () => {
    const { getByText } = render(<Avatar name="Juan Reyes" />);
    expect(getByText('JR')).toBeTruthy();
  });

  it('renders one initial for single-word name', () => {
    const { getByText } = render(<Avatar name="David" />);
    expect(getByText('D')).toBeTruthy();
  });

  it('renders max 2 initials for three-word name (prevents overflow)', () => {
    const { getByText } = render(<Avatar name="Juan Dela Cruz" />);
    expect(getByText('JD')).toBeTruthy();
  });

  it('has accessibility label with full name', () => {
    const { getByLabelText } = render(<Avatar name="Juan Reyes" />);
    expect(getByLabelText('Avatar for Juan Reyes')).toBeTruthy();
  });
});
