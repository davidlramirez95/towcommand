import { test, expect } from '@playwright/test';
import { navigateTo, fillInput, tapButton, expectText, takeEvidence } from './helpers';

/**
 * Negative-path tests for auth screens.
 *
 * Validation in login.tsx and signup.tsx uses Alert.alert() for errors,
 * so there is no inline error text rendered in JSX to assert on.
 * These tests verify the form guard conditions:
 *   - Sign In button is present but does not navigate when fields are empty
 *     (handleLogin returns early, page stays on login)
 *   - Fields are individually clearable / left empty
 *   - Create Account button stays on signup when fields are empty
 */

test.describe('Login Screen — error paths', () => {
  test.beforeEach(async ({ page }) => {
    await navigateTo(page, '/(auth)/login');
  });

  test('sign in button is visible with empty fields', async ({ page }) => {
    // Inputs start empty, button must still render (guard is runtime, not render-time)
    await expect(page.getByLabel('Email')).toHaveValue('');
    await expect(page.getByLabel('Password')).toHaveValue('');
    await expect(page.getByRole('button', { name: 'Sign In' })).toBeVisible();

    await takeEvidence(page, 'login-empty-fields');
  });

  test('sign in button is visible with email only (no password)', async ({ page }) => {
    await fillInput(page, 'Email', 'juan@example.com');
    // password intentionally left empty
    await expect(page.getByLabel('Password')).toHaveValue('');
    await expect(page.getByRole('button', { name: 'Sign In' })).toBeVisible();
    await expect(page.getByRole('button', { name: 'Sign In' })).toBeEnabled();

    await takeEvidence(page, 'login-email-only');
  });

  test('sign in button is visible with password only (no email)', async ({ page }) => {
    // email intentionally left empty
    await fillInput(page, 'Password', 'SomePassword1');
    await expect(page.getByLabel('Email')).toHaveValue('');
    await expect(page.getByRole('button', { name: 'Sign In' })).toBeVisible();
    await expect(page.getByRole('button', { name: 'Sign In' })).toBeEnabled();

    await takeEvidence(page, 'login-password-only');
  });

  test('tapping sign in with empty fields stays on login screen', async ({ page }) => {
    // handleLogin returns early with Alert when email or password is empty
    await tapButton(page, 'Sign In');
    // Page should still show login screen content
    await expectText(page, 'Sign in to book a tow truck');
    await expect(page.getByRole('button', { name: 'Sign In' })).toBeVisible();

    await takeEvidence(page, 'login-empty-submit-stays');
  });
});

test.describe('Signup Screen — error paths', () => {
  test.beforeEach(async ({ page }) => {
    await navigateTo(page, '/(auth)/signup');
  });

  test('create account button is visible with all fields empty', async ({ page }) => {
    await expect(page.getByLabel('Full Name')).toHaveValue('');
    await expect(page.getByLabel('Email')).toHaveValue('');
    await expect(page.getByLabel('Phone Number')).toHaveValue('');
    await expect(page.getByLabel('Password')).toHaveValue('');
    await expect(page.getByRole('button', { name: 'Create Account' })).toBeVisible();

    await takeEvidence(page, 'signup-all-empty');
  });

  test('tapping create account with empty fields stays on signup screen', async ({ page }) => {
    // handleSignUp returns early with Alert when any field is missing
    await tapButton(page, 'Create Account');
    await expectText(page, 'Create Account');
    await expectText(page, 'Join TowCommand PH for roadside help');
    await expect(page.getByRole('button', { name: 'Create Account' })).toBeVisible();

    await takeEvidence(page, 'signup-empty-submit-stays');
  });

  test('tapping create account with partial fields stays on signup screen', async ({ page }) => {
    await fillInput(page, 'Full Name', 'Juan dela Cruz');
    await fillInput(page, 'Email', 'juan@example.com');
    // phone and password intentionally omitted
    await tapButton(page, 'Create Account');
    await expectText(page, 'Create Account');
    await expect(page.getByLabel('Phone Number')).toHaveValue('');

    await takeEvidence(page, 'signup-partial-submit-stays');
  });

  test('password helper text is always visible', async ({ page }) => {
    // The Input helper prop renders static text regardless of input value
    await expectText(page, 'At least 8 characters with a number');

    await takeEvidence(page, 'signup-password-helper');
  });
});
