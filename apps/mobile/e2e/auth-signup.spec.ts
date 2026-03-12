import { test, expect } from '@playwright/test';
import { navigateTo, fillInput, expectText, takeEvidence } from './helpers';

test.describe('Signup Screen', () => {
  test.beforeEach(async ({ page }) => {
    await navigateTo(page, '/(auth)/signup');
  });

  test('renders signup form with all fields', async ({ page }) => {
    await expectText(page, 'Create Account');
    await expectText(page, 'Join TowCommand PH for roadside help');

    await expect(page.getByLabel('Full Name')).toBeVisible();
    await expect(page.getByLabel('Email')).toBeVisible();
    await expect(page.getByLabel('Phone Number')).toBeVisible();
    await expect(page.getByLabel('Password')).toBeVisible();
    await expect(page.getByRole('button', { name: 'Create Account' })).toBeVisible();

    // Helper text
    await expectText(page, 'At least 8 characters with a number');

    // Sign in link
    await expectText(page, 'Already have an account?');
    await expectText(page, 'Sign In');

    await takeEvidence(page, 'signup-screen');
  });

  test('all input fields accept text', async ({ page }) => {
    await fillInput(page, 'Full Name', 'Juan dela Cruz');
    await fillInput(page, 'Email', 'juan@example.com');
    await fillInput(page, 'Phone Number', '+639171234567');
    await fillInput(page, 'Password', 'SecurePass123');

    await expect(page.getByLabel('Full Name')).toHaveValue('Juan dela Cruz');
    await expect(page.getByLabel('Email')).toHaveValue('juan@example.com');
    await expect(page.getByLabel('Phone Number')).toHaveValue('+639171234567');
    await expect(page.getByLabel('Password')).toHaveValue('SecurePass123');

    await takeEvidence(page, 'signup-form-filled');
  });

  test('sign in link navigates back to login', async ({ page }) => {
    await page.getByText('Sign In').click();
    // Expo Router uses client-side routing; wait for login screen content instead of URL
    await expect(page.getByText('Get help on the road', { exact: false }).first()).toBeVisible();
  });
});
