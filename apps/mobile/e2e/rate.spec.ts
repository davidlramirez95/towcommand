import { test, expect } from '@playwright/test';
import { navigateTo, expectText, takeEvidence } from './helpers';

test.describe('Rate Screen', () => {
  test.beforeEach(async ({ page }) => {
    await navigateTo(page, '/booking/rate');
  });

  test('renders rating form with driver info', async ({ page }) => {
    await expectText(page, 'Rate Your Experience');
    await expectText(page, 'Juan Reyes');
    await expectText(page, 'Flatbed Tow');

    await takeEvidence(page, 'rate-screen');
  });

  test('shows 5 star rating buttons', async ({ page }) => {
    const stars = page.getByRole('radio');
    await expect(stars).toHaveCount(5);
  });

  test('has submit review button', async ({ page }) => {
    await expectText(page, 'Submit Review');
  });
});
