import { test, expect } from '@playwright/test';
import { navigateTo, expectText, takeEvidence } from './helpers';

test.describe('Matched Screen', () => {
  test.beforeEach(async ({ page }) => {
    await navigateTo(page, '/booking/matched');
  });

  test('renders driver info with avatar and stats', async ({ page }) => {
    await expectText(page, 'Driver is on the way');
    await expectText(page, 'Juan Reyes');
    await expectText(page, '4.9');
    await expectText(page, '847 jobs');

    await takeEvidence(page, 'matched-screen');
  });

  test('displays OTP digits', async ({ page }) => {
    await expectText(page, 'Digital Padala OTP');
    await expectText(page, 'Share this code when driver arrives');
    // OTP digits: 4, 8, 2, 9, 1, 7
    await expectText(page, '4');
    await expectText(page, '8');
  });

  test('has message and track action buttons', async ({ page }) => {
    await expectText(page, 'Message');
    await expectText(page, 'Track');
  });

  test('has emergency SOS button', async ({ page }) => {
    await expectText(page, 'Emergency SOS');
  });
});
