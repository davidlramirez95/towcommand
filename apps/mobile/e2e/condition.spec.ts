import { test, expect } from '@playwright/test';
import { navigateTo, expectText, takeEvidence } from './helpers';

test.describe('Condition Report Screen', () => {
  test.beforeEach(async ({ page }) => {
    await navigateTo(page, '/booking/condition');
  });

  test('renders photo grid with 8 slots', async ({ page }) => {
    await expectText(page, 'Pre-Tow Condition Report');
    await expectText(page, '8 photos required before towing');

    // Photo angle labels
    await expectText(page, 'Front');
    await expectText(page, 'Rear');
    await expectText(page, 'Left');
    await expectText(page, 'Right');
    await expectText(page, 'FL Tire');

    await takeEvidence(page, 'condition-screen');
  });

  test('shows 360 video option', async ({ page }) => {
    await expectText(page, '360');
    await expectText(page, 'Walk-Around');
  });

  test('shows tamper-proof evidence note', async ({ page }) => {
    await expectText(page, 'tamper-proof');
    await expectText(page, 'SHA-256');
  });

  test('submit button shows photo progress', async ({ page }) => {
    await expectText(page, 'Submit Report');
    await expectText(page, '3/8');
  });
});
