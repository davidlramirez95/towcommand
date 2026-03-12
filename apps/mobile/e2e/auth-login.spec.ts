import { test, expect } from '@playwright/test';
import { navigateTo, fillInput, tapButton, expectText, takeEvidence } from './helpers';

test.describe('Login Screen', () => {
  test.beforeEach(async ({ page }) => {
    await navigateTo(page, '/(auth)/login');
  });

  test('renders login form with all elements', async ({ page }) => {
    await expectText(page, 'TowCommand');
    await expectText(page, 'PILIPINAS');
    await expectText(page, 'Get help on the road');
    await expectText(page, 'Sign in to book a tow truck');

    // Social login buttons
    await expect(page.getByText('Continue with Google')).toBeVisible();
    await expect(page.getByText('Continue with Facebook')).toBeVisible();
    await expect(page.getByText('Continue with Apple')).toBeVisible();

    // Email form
    await expect(page.getByLabel('Email')).toBeVisible();
    await expect(page.getByLabel('Password')).toBeVisible();
    await expect(page.getByRole('button', { name: 'Sign In' })).toBeVisible();

    // Sign up link
    await expectText(page, "Don't have an account?");
    await expectText(page, 'Sign Up');

    await takeEvidence(page, 'login-screen');
  });

  test('email and password inputs accept text', async ({ page }) => {
    await fillInput(page, 'Email', 'test@example.com');
    await fillInput(page, 'Password', 'Test123!');

    await expect(page.getByLabel('Email')).toHaveValue('test@example.com');
    await expect(page.getByLabel('Password')).toHaveValue('Test123!');
  });

  test('sign up link navigates to signup page', async ({ page }) => {
    await page.getByText('Sign Up').click();
    await page.waitForURL('**/signup');
    await expectText(page, 'Create Account');

    await takeEvidence(page, 'login-to-signup-navigation');
  });
});
