import { Page, expect } from '@playwright/test';

/**
 * Shared E2E helpers for TowCommand mobile app tests.
 */

/**
 * Wait for the Expo web app to finish loading.
 * Expo Router renders the root layout once hydration is complete.
 */
export async function waitForAppReady(page: Page) {
  // Wait for any RN-rendered text to appear (app is hydrated)
  await page.waitForSelector('[data-testid]', { timeout: 15_000 }).catch(() => {
    // Fallback: wait for any visible text content
  });
  // Give React time to settle
  await page.waitForTimeout(1000);
}

/**
 * Navigate to a specific route via the browser URL.
 * Expo Router maps file-based routes to URLs on web.
 */
export async function navigateTo(page: Page, route: string) {
  await page.goto(route);
  await waitForAppReady(page);
}

/**
 * Fill an Input component by its label text.
 * Our Input component uses accessibilityLabel={label} which maps to aria-label on web.
 */
export async function fillInput(page: Page, label: string, value: string) {
  const input = page.getByLabel(label);
  await input.fill(value);
}

/**
 * Tap/click a button by its accessible name.
 */
export async function tapButton(page: Page, name: string) {
  await page.getByRole('button', { name }).click();
}

/**
 * Assert visible text exists on screen.
 * Uses .first() to handle strict mode when multiple elements match.
 */
export async function expectText(page: Page, text: string) {
  await expect(page.getByText(text, { exact: false }).first()).toBeVisible();
}

/**
 * Take a named screenshot for PR evidence.
 */
export async function takeEvidence(page: Page, name: string) {
  await page.screenshot({ path: `e2e-results/${name}.png`, fullPage: true });
}
