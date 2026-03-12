/**
 * Badge Component Tests — 2nd Order Logic
 *
 * 2nd order concerns:
 * - All 6 variants map to specific background colors. If a variant
 *   maps to undefined, the badge is invisible against the card background.
 * - Badge is used in booking status display — every booking status
 *   should render a valid badge variant.
 * - Default variant is the fallback — must render something visible.
 */
import React from 'react';
import { render } from '@testing-library/react-native';
import { Badge } from '@/components/ui/Badge';

describe('Badge', () => {
  it('renders text content', () => {
    const { getByText } = render(<Badge>EN_ROUTE</Badge>);
    expect(getByText('EN_ROUTE')).toBeTruthy();
  });

  it('all 6 variants render without crash', () => {
    const variants = ['default', 'success', 'warning', 'error', 'info', 'premium'] as const;

    for (const variant of variants) {
      const { unmount, getByText } = render(
        <Badge variant={variant}>{variant.toUpperCase()}</Badge>,
      );
      expect(getByText(variant.toUpperCase())).toBeTruthy();
      unmount();
    }
  });

  it('default variant renders when variant prop is omitted', () => {
    const { getByText } = render(<Badge>PENDING</Badge>);
    expect(getByText('PENDING')).toBeTruthy();
  });

  it('renders booking status labels (real-world usage)', () => {
    const statuses = ['PENDING', 'MATCHED', 'EN_ROUTE', 'ARRIVED', 'COMPLETED', 'CANCELLED'];

    for (const status of statuses) {
      const { unmount, getByText } = render(<Badge>{status}</Badge>);
      expect(getByText(status)).toBeTruthy();
      unmount();
    }
  });

  it('handles long text without crash (provider trust tier labels)', () => {
    const { getByText } = render(<Badge variant="premium">SUKI ELITE PROVIDER</Badge>);
    expect(getByText('SUKI ELITE PROVIDER')).toBeTruthy();
  });
});
