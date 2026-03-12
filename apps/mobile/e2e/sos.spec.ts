import { test, expect } from '@playwright/test';
import { navigateTo, expectText, takeEvidence } from './helpers';

test.describe('SOS Screen', () => {
  test.beforeEach(async ({ page }) => {
    await navigateTo(page, '/sos');
  });

  test('renders emergency SOS screen', async ({ page }) => {
    await expectText(page, 'Emergency SOS');
    await expectText(page, 'Press the button below to alert');
    await expectText(page, 'SOS');
    await expectText(page, 'Press to alert');

    // Close button
    await expect(page.getByLabel('Close SOS')).toBeVisible();

    // Disclaimer
    await expectText(page, 'call 911 directly');

    await takeEvidence(page, 'sos-screen');
  });

  test('SOS trigger button is visible and tappable', async ({ page }) => {
    const sosButton = page.getByLabel('Trigger SOS Alert');
    await expect(sosButton).toBeVisible();
    await expect(sosButton).toBeEnabled();
  });

  test('close button is accessible', async ({ page }) => {
    const closeButton = page.getByLabel('Close SOS');
    await expect(closeButton).toBeVisible();

    await takeEvidence(page, 'sos-close-button');
  });
});
