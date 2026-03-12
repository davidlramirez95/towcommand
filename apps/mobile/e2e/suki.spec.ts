import { test, expect } from '@playwright/test';
import { navigateTo, expectText, takeEvidence } from './helpers';

test.describe('Suki Rewards Screen', () => {
  test.beforeEach(async ({ page }) => {
    await navigateTo(page, '/(tabs)/suki');
  });

  test('renders loyalty tier with progress', async ({ page }) => {
    await expectText(page, 'Suki Rewards');
    await expectText(page, 'Silver Member');
    await expectText(page, '4 of 6 bookings to Gold');

    await takeEvidence(page, 'suki-screen');
  });

  test('renders benefits list', async ({ page }) => {
    await expectText(page, 'YOUR BENEFITS');
    await expectText(page, '5% off all services');
    await expectText(page, 'Priority matching');
    await expectText(page, 'VIP support');
  });

  test('renders points history', async ({ page }) => {
    await expectText(page, 'POINTS HISTORY');
    await expectText(page, 'Flatbed Tow');
    await expectText(page, '+1 point');
  });
});
