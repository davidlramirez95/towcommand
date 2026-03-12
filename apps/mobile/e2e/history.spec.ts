import { test, expect } from '@playwright/test';
import { navigateTo, expectText, takeEvidence } from './helpers';

test.describe('History Screen', () => {
  test.beforeEach(async ({ page }) => {
    await navigateTo(page, '/(tabs)/history');
  });

  test('renders history screen with empty state', async ({ page }) => {
    await expectText(page, 'Activity');
    await expectText(page, 'Your booking history');

    // Empty state
    await expectText(page, 'No bookings yet');
    await expectText(page, 'Your tow truck bookings will appear here');

    await takeEvidence(page, 'history-screen');
  });
});
