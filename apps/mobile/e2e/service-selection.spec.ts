import { test, expect } from '@playwright/test';
import { navigateTo, expectText, takeEvidence } from './helpers';

test.describe('Service Selection Screen', () => {
  test.beforeEach(async ({ page }) => {
    await navigateTo(page, '/booking/service');
  });

  test('renders all 8 service types', async ({ page }) => {
    await expectText(page, 'Select Service');

    const services = [
      'Flatbed Towing',
      'Wheel Lift Towing',
      'Motorcycle Towing',
      'Jumpstart',
      'Tire Change',
      'Lockout Service',
      'Fuel Delivery',
      'Winch Recovery',
    ];

    for (const service of services) {
      await expect(page.getByText(service)).toBeVisible();
    }

    await takeEvidence(page, 'service-selection');
  });

  test('shows service descriptions', async ({ page }) => {
    await expectText(page, 'Best for sedans, SUVs, luxury vehicles');
    await expectText(page, 'Dead battery?');
    await expectText(page, 'Locked out of your vehicle?');
  });

  test('back button is visible', async ({ page }) => {
    // Back arrow
    await expect(page.getByText('←')).toBeVisible();
  });
});
