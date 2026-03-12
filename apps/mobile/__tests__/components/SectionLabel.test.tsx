/**
 * SectionLabel Component Tests — 2nd Order Logic
 *
 * 2nd order concerns:
 * - Text MUST be uppercased (consistency across all section headers)
 * - Missing uppercase would break visual hierarchy and look like body text
 */
import React from 'react';
import { render } from '@testing-library/react-native';
import { SectionLabel } from '@/components/ui/SectionLabel';

describe('SectionLabel', () => {
  it('renders text in uppercase', () => {
    const { getByText } = render(<SectionLabel>saved vehicles</SectionLabel>);
    expect(getByText('SAVED VEHICLES')).toBeTruthy();
  });

  it('renders already-uppercase text unchanged', () => {
    const { getByText } = render(<SectionLabel>PAYMENT METHOD</SectionLabel>);
    expect(getByText('PAYMENT METHOD')).toBeTruthy();
  });
});
