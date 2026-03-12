import { test, expect } from '@playwright/test';
import { navigateTo, expectText, takeEvidence } from './helpers';

test.describe('Typhoon Mode Screen', () => {
  test.beforeEach(async ({ page }) => {
    await navigateTo(page, '/typhoon');
  });

  test('renders typhoon alert with flood info', async ({ page }) => {
    await expectText(page, 'TYPHOON ALERT');
    await expectText(page, 'Typhoon Mode Active');
    await expectText(page, 'Signal #3');

    await expectText(page, 'Flood Level');
    await expectText(page, 'Knee-Deep');

    await takeEvidence(page, 'typhoon-screen');
  });

  test('shows surge pricing with MMDA regulation', async ({ page }) => {
    await expectText(page, 'Surge');
    await expectText(page, '1.5');
    await expectText(page, 'MMDA-regulated');
  });

  test('shows safety guarantee', async ({ page }) => {
    await expectText(page, 'Safety Guaranteed');
    await expectText(page, 'insurance');
  });

  test('has book emergency tow button', async ({ page }) => {
    await expectText(page, 'Book Emergency Tow');
  });
});
