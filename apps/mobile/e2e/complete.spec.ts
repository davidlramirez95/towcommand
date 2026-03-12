import { test, expect } from '@playwright/test';
import { navigateTo, expectText, takeEvidence } from './helpers';

test.describe('Complete Screen', () => {
  test.beforeEach(async ({ page }) => {
    await navigateTo(page, '/booking/complete');
  });

  test('renders job completion summary', async ({ page }) => {
    await expectText(page, 'Tow Complete!');
    await expectText(page, 'delivered safely');
    await expectText(page, 'TC-2026-00847');
    await expectText(page, 'Flatbed Tow');
    await expectText(page, '12.4 km');
    await expectText(page, '45 min');

    await takeEvidence(page, 'complete-screen');
  });

  test('shows total paid with payment method', async ({ page }) => {
    await expectText(page, 'Total Paid');
    await expectText(page, '1,850');
    await expectText(page, 'GCash');
  });

  test('shows suki points earned', async ({ page }) => {
    await expectText(page, '+1 Suki Point earned!');
    await expectText(page, 'Gold tier');
  });

  test('has rate and home buttons', async ({ page }) => {
    await expectText(page, 'Rate Your Experience');
    await expectText(page, 'Back to Home');
  });
});
