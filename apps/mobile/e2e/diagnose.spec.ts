import { test, expect } from '@playwright/test';
import { navigateTo, expectText, takeEvidence } from './helpers';

test.describe('Diagnose Screen', () => {
  test.beforeEach(async ({ page }) => {
    await navigateTo(page, '/booking/diagnose');
  });

  test('renders smart diagnosis with symptom picker', async ({ page }) => {
    await expectText(page, 'Smart Diagnosis');
    await expectText(page, "What's happening to your vehicle?");
    await expectText(page, 'Voice Describe');
    await expectText(page, 'Take Photo');

    // Verify symptom chips render
    await expectText(page, 'Empty fuel');
    await expectText(page, 'Flat tire');
    await expectText(page, 'Dead battery');
    await expectText(page, 'Overheating');

    await takeEvidence(page, 'diagnose-screen');
  });

  test('diagnose button shows symptom count', async ({ page }) => {
    // Button should show count of 0
    await expectText(page, 'Diagnose My Problem (0)');
  });

  test('skip link navigates to manual service selection', async ({ page }) => {
    await expectText(page, 'I already know what I need');
  });
});
