import { test, expect } from '@playwright/test';
import { navigateTo, expectText, takeEvidence } from './helpers';

test.describe('Home Screen', () => {
  test.beforeEach(async ({ page }) => {
    await navigateTo(page, '/(tabs)');
  });

  test('renders home screen with key elements', async ({ page }) => {
    // Greeting
    await expectText(page, 'Magandang hapon');

    // SOS button
    await expect(page.getByLabel('Emergency SOS')).toBeVisible();

    // AI Diagnosis card
    await expectText(page, 'wrong with your car');
    await expectText(page, 'AI will diagnose');

    // Quick services
    await expectText(page, 'Quick Services');
    await expectText(page, 'Tow');
    await expectText(page, 'Fuel');

    // Map with trucks badge
    await expectText(page, '3 trucks nearby');

    await takeEvidence(page, 'home-screen');
  });

  test('SOS button is tappable', async ({ page }) => {
    const sosButton = page.getByLabel('Emergency SOS');
    await expect(sosButton).toBeVisible();
    await expect(sosButton).toBeEnabled();
  });

  test('tab navigation includes suki tab', async ({ page }) => {
    await expectText(page, 'Home');
    await expectText(page, 'Activity');
    await expectText(page, 'Suki');
    await expectText(page, 'Account');

    await takeEvidence(page, 'home-tab-bar');
  });
});
