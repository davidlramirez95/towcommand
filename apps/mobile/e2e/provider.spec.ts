import { test, expect } from '@playwright/test';
import { navigateTo, expectText, takeEvidence } from './helpers';

test.describe('Provider Dashboard', () => {
  test.beforeEach(async ({ page }) => {
    await navigateTo(page, '/provider/dashboard');
  });

  test('renders provider dashboard with stats', async ({ page }) => {
    await expectText(page, 'Provider Dashboard');
    await expectText(page, 'Online');

    // Stats card
    await expectText(page, "Today's Summary");
    await expectText(page, 'Jobs');
    await expectText(page, 'Earnings');
    await expectText(page, 'Rating');

    // Empty state
    await expectText(page, 'Waiting for jobs');
    await expectText(page, 'New job requests will appear here');

    await takeEvidence(page, 'provider-dashboard');
  });
});

test.describe('Provider Earnings', () => {
  test.beforeEach(async ({ page }) => {
    await navigateTo(page, '/provider/earnings');
  });

  test('renders all earnings periods', async ({ page }) => {
    await expectText(page, 'Earnings');
    await expectText(page, 'Today');
    await expectText(page, 'This Week');
    await expectText(page, 'This Month');
    await expectText(page, 'All Time');

    // PHP amounts
    await expectText(page, 'PHP 0.00');

    // Job counts
    await expectText(page, '0 jobs');

    await takeEvidence(page, 'provider-earnings');
  });
});
