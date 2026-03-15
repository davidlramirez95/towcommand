# E2E Evidence Per PR

## Description

Mandatory skill that enforces screenshot-backed E2E test evidence on every PR. No PR ships without headless Playwright screenshots posted as a comment. This is a hard gate — not optional, not "nice to have."

## When to Use

- Every time a PR is created or updated
- Every time `/ship` is invoked
- When any agent (golang-pro, expo-mobile-dev, full-stack-implementer) completes work that will become a PR
- When reviewing someone else's PR for evidence completeness

## When NOT to Skip

- "It's just a docs change" — still needs evidence that existing E2E didn't break
- "It's just a config change" — config changes can break rendering
- "Tests already passed in CI" — CI results are not screenshots; reviewers need visual proof
- "I ran it locally" — unverifiable claim; show the output

---

## Core Principle

**"If there's no screenshot, it didn't happen."**

### Why (2nd-Order)

1st-order: Screenshots prove the app renders correctly after changes.
2nd-order: Over time, screenshot evidence creates a **visual changelog** across PRs. Reviewers can compare screenshots between PRs to spot unintended visual regressions that no test assertion would catch — a button that shifted 10px, a color that got slightly wrong, a layout that broke on Pixel 7 but not iPhone 14.

---

## Evidence Capture Protocol

### Step 1: Run Full E2E Suite

```bash
cd apps/mobile && pnpm test:e2e
```

This runs all Playwright specs across both device profiles (iPhone 14 + Pixel 7). Every spec that calls `takeEvidence(page, 'name')` writes a screenshot to `e2e-results/`.

**Result required:** `N passed (X.Xm)` with zero failures. If any test fails, fix it before proceeding.

### Step 2: Capture Key Screenshots

The `takeEvidence()` helper in `e2e/helpers.ts` captures full-page screenshots:

```typescript
export async function takeEvidence(page: Page, name: string) {
  await page.screenshot({ path: `e2e-results/${name}.png`, fullPage: true });
}
```

**Every E2E spec MUST call `takeEvidence()` at least once** with a descriptive name. If you're writing a new spec, include evidence capture:

```typescript
test('renders booking form with all fields', async ({ page }) => {
  await expectText(page, 'Select Service');
  await expectText(page, 'Choose Vehicle');
  // ... assertions ...

  await takeEvidence(page, 'booking-form');  // ← MANDATORY
});
```

### Step 3: Upload Screenshots to Evidence Branch

Push screenshots to an orphan branch so they can be embedded in PR comments via raw GitHub URLs:

```bash
# From project root
git stash
git checkout --orphan e2e-evidence-pr{NUMBER}
git rm -rf .
cp apps/mobile/e2e-results/*.png .
git add *.png
git commit -m "E2E evidence screenshots for PR #{NUMBER}"
git push -f origin e2e-evidence-pr{NUMBER}
git checkout {original-branch}
git stash pop
```

**Base URL pattern:**
```
https://raw.githubusercontent.com/davidlramirez95/towcommand/e2e-evidence-pr{NUMBER}/{screenshot-name}.png
```

### Step 4: Post PR Comment with Embedded Screenshots

Post a structured comment using `gh pr comment`. The comment MUST include:

1. **Test result summary** — total passed/failed, runtime, device profiles
2. **Screenshots organized by flow** — tables with inline images
3. **New test callouts** — highlight any new test suites added in this PR

#### Comment Template

```markdown
## 📸 E2E Screenshot Evidence — Headless Playwright ({commit-hash})

**{N}/{N} tests passed** | iPhone 14 + Pixel 7 | {X} test suites | {X.X} min runtime

---

### 🏠 Core Screens

| Home | Login | Profile | History |
|:---:|:---:|:---:|:---:|
| ![home]({BASE}/home-screen.png) | ![login]({BASE}/login-screen.png) | ![profile]({BASE}/profile-screen.png) | ![history]({BASE}/history-empty-state.png) |

### 🚗 Booking Flow

| Service Selection | AI Diagnosis | Condition Report | Price Breakdown |
|:---:|:---:|:---:|:---:|
| ![service]({BASE}/service-selection.png) | ![diagnose]({BASE}/diagnose-screen.png) | ![condition]({BASE}/condition-screen.png) | ![price]({BASE}/price-screen.png) |

| Matching → Matched | OTP | Chat | Complete |
|:---:|:---:|:---:|:---:|
| ![matching]({BASE}/matching-auto-redirect.png) | ![otp]({BASE}/matching-otp-result.png) | ![chat]({BASE}/chat-screen.png) | ![complete]({BASE}/complete-screen.png) |

### 🚨 Safety & Rewards

| SOS | Typhoon Mode | Suki Rewards |
|:---:|:---:|:---:|
| ![sos]({BASE}/sos-screen.png) | ![typhoon]({BASE}/typhoon-screen.png) | ![suki]({BASE}/suki-screen.png) |

### 🏢 Provider

| Dashboard | Earnings |
|:---:|:---:|
| ![dashboard]({BASE}/provider-dashboard.png) | ![earnings]({BASE}/provider-earnings.png) |

### ❌ Error / Empty States (if applicable)

| Login Empty | Booking No Context |
|:---:|:---:|
| ![login-empty]({BASE}/login-empty-fields.png) | ![tracking-empty]({BASE}/booking-tracking-no-context-status.png) |

---

### Test Results

```
{N} passed ({X.X}m)
0 failed
```

🤖 Generated with [Claude Code](https://claude.com/claude-code)
```

### Step 5: Verify Screenshots Render in PR

After posting, open the PR comment in a browser and verify the images actually render. GitHub raw URLs sometimes have a brief caching delay.

---

## Screenshot Naming Convention

Screenshots must use kebab-case names that describe what's shown:

| Pattern | Example |
|---|---|
| `{screen-name}.png` | `home-screen.png`, `login-screen.png` |
| `{screen}-{state}.png` | `history-empty-state.png`, `login-empty-fields.png` |
| `{flow}-{step}.png` | `matching-auto-redirect.png`, `matching-otp-result.png` |
| `{screen}-{feature}.png` | `provider-earnings-all-zero.png` |

---

## Writing New E2E Specs — Evidence Checklist

When creating a new E2E spec file, ensure:

- [ ] At least one `takeEvidence()` call per test suite
- [ ] Screenshot name is descriptive (not `test-1.png`)
- [ ] Tests run on both device profiles (Playwright config handles this)
- [ ] Happy path AND error/empty states covered
- [ ] Screenshot captures the full page (`fullPage: true`)

---

## Evidence for Backend-Only PRs

Even if the PR only touches Go code (`cmd/`, `internal/`):

1. Run `task test-unit` and show output in PR body
2. If the change affects an API endpoint, run the mobile E2E suite to verify the frontend still works against the contract
3. If the change is purely internal (no API contract change), post unit test output as evidence — screenshots are not required

**2nd-Order reasoning**: A "backend-only" change that modifies an API response schema will break mobile screens. The E2E suite catches this. Always ask: "does this backend change affect any API response that a mobile screen consumes?"

---

## Cleanup

After a PR is merged, the evidence branch can be deleted:

```bash
git push origin --delete e2e-evidence-pr{NUMBER}
```

Keep evidence branches for the last 5 merged PRs for historical reference, then clean up older ones.

---

## Anti-Patterns

| Anti-Pattern | Why It Fails | Correct Behavior |
|---|---|---|
| "Tests pass" with no output | Unverifiable claim | Show `pnpm test:e2e` output |
| Screenshots only from one device | Pixel 7 layout bugs missed on iPhone | Run on both device profiles |
| Evidence from a previous PR | Stale — doesn't reflect current changes | Run fresh E2E on the current commit |
| Screenshots in PR body (not comment) | Gets buried in description edits | Post as a separate comment for visibility |
| Skipping evidence for "small changes" | Small changes cause big regressions | Every PR, every time |
| No `takeEvidence()` in new specs | New screens have no visual proof | Add at least one per new spec |
| Evidence branch left forever | Branch clutter | Clean up after merge |
