# /ship — Ship Code to PR

## Purpose

Complete workflow to lint, test, E2E verify, commit, and create a PR ready for merge. Enforces TowCommand's mandatory E2E evidence requirement.

## Usage

```
/ship [commit message or 'quick']
```

## Arguments

- `$ARGUMENTS`: Commit message, or `quick` to auto-generate + skip review

---

Ship the current changes: **$ARGUMENTS**

## Workflow

### Phase 1: Pre-Ship Audit

1. **Inventory changes**
   ```bash
   git status
   git diff --staged
   git diff
   ```

2. **2nd-Order Check** — Before shipping, ask:
   - What could this change break downstream that isn't obvious from the diff?
   - Are there event consumers, Lambda subscribers, or mobile screens that depend on the changed contracts?
   - Does this change affect DynamoDB key schemas, GSI projections, or access patterns?
   - Could this cause a state machine transition failure in the booking flow?

3. **Quick validation**
   - No secrets, API keys, or `.env` files in changes
   - No debug statements (`fmt.Println`, `console.log`)
   - No commented-out code blocks
   - No `any` types in TypeScript files

### Phase 2: Lint & Test (Go Backend)

If changes touch `cmd/`, `internal/`, or `*.go` files:

```bash
task lint          # golangci-lint run ./...
task test-unit     # go test -race -cover ./...
```

Fix any failures before proceeding. Do NOT skip lint warnings.

### Phase 3: Lint & Test (Mobile)

If changes touch `apps/mobile/`:

```bash
cd apps/mobile && pnpm lint
cd apps/mobile && pnpm tsc --noEmit
cd apps/mobile && pnpm test
```

### Phase 4: Build

```bash
# Go backend
task build         # or: task build-func

# Mobile (if changed)
cd apps/mobile && pnpm expo export --platform web
```

### Phase 5: E2E Evidence (MANDATORY)

**This is a hard requirement. Do NOT skip this phase.**
**Skill reference:** `.claude/skills/methodology/e2e-evidence-per-pr/SKILL.md`

#### Backend E2E
```bash
task test-e2e      # or run targeted integration tests
```

#### Mobile E2E
```bash
cd apps/mobile && pnpm test:e2e   # Playwright against Expo web
```

All E2E specs must call `takeEvidence(page, 'name')` for screenshot capture.

#### Upload Screenshots & Post to PR

```bash
# Push screenshots to evidence branch
git stash
git checkout --orphan e2e-evidence-pr{NUMBER}
git rm -rf .
cp apps/mobile/e2e-results/*.png .
git add *.png
git commit -m "E2E evidence screenshots for PR #{NUMBER}"
git push -f origin e2e-evidence-pr{NUMBER}
git checkout {original-branch}
git stash pop

# Post PR comment with embedded screenshots (see skill for template)
gh pr comment {NUMBER} --body "..."
```

Every PR comment MUST include:
1. Pass/fail summary with test count and runtime
2. Embedded screenshots organized by screen flow (tables with inline images)
3. Both device profiles represented (iPhone 14 + Pixel 7)
4. New test suites highlighted

### Phase 6: Commit

1. Stage specific files (never `git add -A` blindly):
   ```bash
   git add [specific files]
   ```

2. Generate conventional commit message:
   ```bash
   git commit -m "$(cat <<'EOF'
   type(scope): subject

   - Change detail 1
   - Change detail 2

   Co-Authored-By: Claude Opus 4.6 <noreply@anthropic.com>
   EOF
   )"
   ```

### Phase 7: Push & PR

1. Push with tracking:
   ```bash
   git push -u origin $(git branch --show-current)
   ```

2. Create PR with E2E evidence:
   ```bash
   gh pr create --title "type(scope): description" --body "$(cat <<'EOF'
   ## Summary
   - Change 1
   - Change 2

   ## Test Plan
   - [ ] Unit tests pass
   - [ ] E2E tests pass
   - [ ] Manual verification

   ## E2E Evidence
   [Paste test results, screenshots, or trace output here]

   🤖 Generated with [Claude Code](https://claude.com/claude-code)
   EOF
   )"
   ```

3. Post E2E results as PR comment if not in body.

## Quick Ship Mode

When using `/ship quick`:
- Auto-generate commit message from diff analysis
- Skip detailed self-review
- Still run ALL tests and E2E (never skip)
- Minimal output

## Pre-Ship Checklist

- [ ] All changes staged intentionally
- [ ] No secrets or debug code
- [ ] Lint passes (Go + mobile)
- [ ] Unit tests pass
- [ ] Build succeeds
- [ ] E2E evidence captured
- [ ] Commit message follows conventional format
- [ ] PR includes E2E results
