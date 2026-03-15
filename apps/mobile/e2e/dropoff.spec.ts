import { test, expect } from '@playwright/test';
import { navigateTo, expectText, takeEvidence } from './helpers';

test.describe('Dropoff Screen', () => {
  test.beforeEach(async ({ page }) => {
    await navigateTo(page, '/booking/dropoff');
  });

  test('renders pickup and dropoff locations', async ({ page }) => {
    await expectText(page, 'Drop-off Location');
    await expectText(page, 'PICKUP LOCATION');
    await expectText(page, 'Current Location');
    await expectText(page, 'EDSA cor. Ayala Ave, Makati City');
    await expectText(page, 'DROP-OFF LOCATION');
    await expectText(page, 'Toyota Shaw, Mandaluyong');

    await takeEvidence(page, 'dropoff-screen');
  });

  test('renders recent locations', async ({ page }) => {
    await expectText(page, 'RECENT LOCATIONS');
    await expectText(page, 'BF Homes');
    await expectText(page, 'Mitsubishi Ortigas');
    await expectText(page, 'AutoHub SLEX');
  });

  test('has confirm route button', async ({ page }) => {
    await expectText(page, 'Confirm Route');
  });
});
