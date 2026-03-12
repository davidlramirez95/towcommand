import { test, expect } from '@playwright/test';
import { navigateTo, expectText, takeEvidence } from './helpers';

test.describe('Booking Tracking Screen', () => {
  test.beforeEach(async ({ page }) => {
    await navigateTo(page, '/booking/abc123');
  });

  test('renders tracking screen with booking ID', async ({ page }) => {
    await expectText(page, 'Live Tracking Map');
    await expectText(page, 'Booking #abc123');

    // Status badge (default = LOADING when no active booking)
    await expectText(page, 'LOADING');

    await takeEvidence(page, 'booking-tracking');
  });
});
