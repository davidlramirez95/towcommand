import { test, expect } from '@playwright/test';
import { navigateTo, expectText, takeEvidence } from './helpers';

test.describe('Profile Screen', () => {
  test.beforeEach(async ({ page }) => {
    await navigateTo(page, '/(tabs)/profile');
  });

  test('renders profile screen with menu items', async ({ page }) => {
    await expectText(page, 'Account');

    // Default avatar (no user = "U")
    await expectText(page, 'U');
    await expectText(page, 'User');

    // Menu items
    await expectText(page, 'My Vehicles');
    await expectText(page, 'Payment Methods');
    await expectText(page, 'Notifications');
    await expectText(page, 'Help & Support');
    await expectText(page, 'About TowCommand');

    // Sign out button
    await expect(page.getByRole('button', { name: 'Sign Out' })).toBeVisible();

    await takeEvidence(page, 'profile-screen');
  });

  test('menu items are accessible buttons', async ({ page }) => {
    const menuItems = ['My Vehicles', 'Payment Methods', 'Notifications', 'Help & Support', 'About TowCommand'];

    for (const item of menuItems) {
      await expect(page.getByLabel(item)).toBeVisible();
    }
  });
});
