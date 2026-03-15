import { test, expect } from '@playwright/test';
import { navigateTo, expectText, takeEvidence } from './helpers';

test.describe('Chat Screen', () => {
  test.beforeEach(async ({ page }) => {
    await navigateTo(page, '/booking/chat');
  });

  test('renders chat with driver messages', async ({ page }) => {
    await expectText(page, 'Juan Reyes');
    await expectText(page, 'Online');
    await expectText(page, 'Magandang hapon po');
    await expectText(page, "McDonald's EDSA");
    await expectText(page, 'White Montero');

    await takeEvidence(page, 'chat-screen');
  });

  test('has message input and send button', async ({ page }) => {
    await expect(page.getByLabel('Message input')).toBeVisible();
    await expect(page.getByLabel('Send message')).toBeVisible();
  });

  test('has quick action buttons (location, photo, voice)', async ({ page }) => {
    await expectText(page, '📍');
    await expectText(page, '📷');
    await expectText(page, '🎤');
  });
});
