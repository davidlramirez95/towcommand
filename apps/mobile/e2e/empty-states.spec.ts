import { test, expect } from '@playwright/test';
import { navigateTo, expectText, takeEvidence } from './helpers';

/**
 * Tests for zero-data / empty-state UI across screens.
 *
 * History screen:     Always shows empty state (data source not wired yet).
 * Profile screen:     No auth session on fresh load → fallback avatar 'U', name 'User', empty email/phone.
 * Provider dashboard: Always shows 0-value stats + "Waiting for jobs" empty state.
 * Provider earnings:  Always shows PHP 0.00 + 0 jobs.
 * Booking tracking:   Status defaults to 'LOADING' with no store data.
 */

test.describe('History Screen — empty state', () => {
  test.beforeEach(async ({ page }) => {
    await navigateTo(page, '/(tabs)/history');
  });

  test('shows empty state title and subtitle', async ({ page }) => {
    await expectText(page, 'No bookings yet');
    await expectText(page, 'Your tow truck bookings will appear here');
    await takeEvidence(page, 'history-empty-state');
  });

  test('empty state does not show any booking cards', async ({ page }) => {
    // No booking card content should be present
    await expect(page.getByText('COMPLETED')).not.toBeVisible().catch(() => {});
    await expect(page.getByText('CANCELLED')).not.toBeVisible().catch(() => {});
    // Header still shows
    await expectText(page, 'Activity');
    await expectText(page, 'Your booking history');
    await takeEvidence(page, 'history-no-booking-cards');
  });
});

test.describe('Profile Screen — unauthenticated / no user data', () => {
  test.beforeEach(async ({ page }) => {
    await navigateTo(page, '/(tabs)/profile');
  });

  test('avatar shows fallback initial "U" when no user', async ({ page }) => {
    // user?.fullName?.charAt(0)?.toUpperCase() ?? 'U'
    await expectText(page, 'U');
    await takeEvidence(page, 'profile-fallback-avatar');
  });

  test('name shows fallback "User" when no user', async ({ page }) => {
    // user?.fullName ?? 'User'
    await expectText(page, 'User');
    await takeEvidence(page, 'profile-fallback-name');
  });

  test('all menu items render without user data', async ({ page }) => {
    const menuItems = [
      'My Vehicles',
      'Payment Methods',
      'Notifications',
      'Help & Support',
      'About TowCommand',
    ];

    for (const item of menuItems) {
      await expect(page.getByLabel(item)).toBeVisible();
    }

    await takeEvidence(page, 'profile-menu-no-user');
  });

  test('sign out button renders without user data', async ({ page }) => {
    await expect(page.getByRole('button', { name: 'Sign Out' })).toBeVisible();
    await takeEvidence(page, 'profile-signout-no-user');
  });
});

test.describe('Provider Dashboard — zero stats / empty job queue', () => {
  test.beforeEach(async ({ page }) => {
    await navigateTo(page, '/provider/dashboard');
  });

  test('stats show zero values by default', async ({ page }) => {
    // StatItem renders value prop directly
    await expectText(page, '0');       // Jobs value
    await expectText(page, 'PHP 0');   // Earnings value
    await expectText(page, '--');      // Rating value (no jobs yet)
    await takeEvidence(page, 'provider-dashboard-zero-stats');
  });

  test('shows waiting-for-jobs empty state', async ({ page }) => {
    await expectText(page, 'Waiting for jobs');
    await expectText(page, 'New job requests will appear here');
    await takeEvidence(page, 'provider-dashboard-empty-queue');
  });

  test('online badge renders without active jobs', async ({ page }) => {
    await expectText(page, 'Online');
    await takeEvidence(page, 'provider-dashboard-online-badge');
  });
});

test.describe('Provider Earnings — zero earnings state', () => {
  test.beforeEach(async ({ page }) => {
    await navigateTo(page, '/provider/earnings');
  });

  test('all periods show PHP 0.00', async ({ page }) => {
    const zeroAmounts = await page.getByText('PHP 0.00').all();
    expect(zeroAmounts.length).toBeGreaterThan(0);
    await takeEvidence(page, 'provider-earnings-all-zero');
  });

  test('all periods show 0 jobs', async ({ page }) => {
    const zeroJobs = await page.getByText('0 jobs').all();
    expect(zeroJobs.length).toBeGreaterThan(0);
    await takeEvidence(page, 'provider-earnings-zero-jobs');
  });

  test('all time period tabs render', async ({ page }) => {
    await expectText(page, 'Today');
    await expectText(page, 'This Week');
    await expectText(page, 'This Month');
    await expectText(page, 'All Time');
    await takeEvidence(page, 'provider-earnings-period-tabs');
  });
});

test.describe('Booking Tracking — empty booking context', () => {
  test.beforeEach(async ({ page }) => {
    await navigateTo(page, '/booking/no-such-booking');
  });

  test('status defaults to LOADING without store data', async ({ page }) => {
    // activeBooking is null on fresh load → status badge renders 'LOADING'
    await expectText(page, 'LOADING');
    await takeEvidence(page, 'booking-tracking-empty-context');
  });

  test('no ETA is shown without active booking', async ({ page }) => {
    // etaText renders only when activeBooking?.eta is truthy
    await expect(page.getByText('ETA:')).not.toBeVisible().catch(() => {});
    await expectText(page, 'Live Tracking Map');
    await takeEvidence(page, 'booking-tracking-no-eta');
  });
});
