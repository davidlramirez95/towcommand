import { test, expect } from '@playwright/test';
import { navigateTo, expectText, takeEvidence } from './helpers';

test.describe('Home Screen', () => {
  test.beforeEach(async ({ page }) => {
    await navigateTo(page, '/(tabs)');
  });

  test('renders home screen with key elements', async ({ page }) => {
    // Greeting (defaults to "there" when no user)
    await expectText(page, 'Mabuhay');
    await expectText(page, 'Need help on the road?');

    // SOS button
    await expect(page.getByLabel('Emergency SOS')).toBeVisible();

    // Map placeholder
    await expectText(page, 'Map loads here');

    // Quick actions
    await expectText(page, 'Request a Tow');
    await expectText(page, 'Get help in minutes');
    await expectText(page, 'Diagnose');
    await expectText(page, 'History');

    await takeEvidence(page, 'home-screen');
  });

  test('SOS button is tappable', async ({ page }) => {
    const sosButton = page.getByLabel('Emergency SOS');
    await expect(sosButton).toBeVisible();
    await expect(sosButton).toBeEnabled();
  });

  test('tab navigation is visible', async ({ page }) => {
    // Tab bar should be present with Home, Activity, Account
    await expectText(page, 'Home');
    await expectText(page, 'Activity');
    await expectText(page, 'Account');

    await takeEvidence(page, 'home-tab-bar');
  });
});
