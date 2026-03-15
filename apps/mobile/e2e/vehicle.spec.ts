import { test, expect } from '@playwright/test';
import { navigateTo, expectText, takeEvidence } from './helpers';

test.describe('Vehicle Screen', () => {
  test.beforeEach(async ({ page }) => {
    await navigateTo(page, '/booking/vehicle');
  });

  test('renders vehicle selection with saved vehicles', async ({ page }) => {
    await expectText(page, 'Vehicle Details');
    await expectText(page, 'SAVED VEHICLES');
    await expectText(page, 'Montero GLS');
    await expectText(page, 'ABC 1234');
    await expectText(page, 'Add New Vehicle');

    await takeEvidence(page, 'vehicle-screen');
  });

  test('renders vehicle condition options', async ({ page }) => {
    await expectText(page, 'VEHICLE CONDITION');
    await expectText(page, "Engine won't start");
    await expectText(page, 'Flat tire(s)');
    await expectText(page, 'Accident damage');
    await expectText(page, 'Other / Not sure');
  });

  test('has continue button', async ({ page }) => {
    await expectText(page, 'Continue');
  });
});
