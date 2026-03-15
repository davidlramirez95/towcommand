import { test, expect } from '@playwright/test';
import { navigateTo, expectText, takeEvidence } from './helpers';

test.describe('Price Screen', () => {
  test.beforeEach(async ({ page }) => {
    await navigateTo(page, '/booking/price');
  });

  test('renders MMDA-compliant price breakdown', async ({ page }) => {
    await expectText(page, 'Price Estimate');
    await expectText(page, 'Estimated Total');
    await expectText(page, '1,850');
    await expectText(page, 'MMDA Reg. 24-004 compliant pricing');

    // Line items
    await expectText(page, 'Base fare');
    await expectText(page, 'Distance');
    await expectText(page, 'Weight surcharge');
    await expectText(page, 'Platform fee');

    await takeEvidence(page, 'price-screen');
  });

  test('renders payment method', async ({ page }) => {
    await expectText(page, 'PAYMENT METHOD');
    await expectText(page, 'GCash');
    await expectText(page, '8847');
  });

  test('has book now button with price', async ({ page }) => {
    await expectText(page, 'Book Now');
  });
});
