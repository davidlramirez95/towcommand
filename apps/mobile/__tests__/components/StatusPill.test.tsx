/**
 * StatusPill Component Tests — 2nd Order Logic
 *
 * 2nd order concerns:
 * - All 5 variants must render without crashing (regression guard)
 * - Label must be uppercased (consistency with design system)
 * - Wrong variant color = user misreads booking status = confusion
 */
import React from 'react';
import { render } from '@testing-library/react-native';
import { StatusPill } from '@/components/ui/StatusPill';

describe('StatusPill', () => {
  const variants = ['success', 'warning', 'danger', 'info', 'neutral'] as const;

  variants.forEach((variant) => {
    it(`renders ${variant} variant without crashing`, () => {
      const { getByText } = render(<StatusPill label="verified" variant={variant} />);
      expect(getByText('VERIFIED')).toBeTruthy();
    });
  });

  it('uppercases the label text', () => {
    const { getByText } = render(<StatusPill label="completed" variant="success" />);
    expect(getByText('COMPLETED')).toBeTruthy();
  });
});
