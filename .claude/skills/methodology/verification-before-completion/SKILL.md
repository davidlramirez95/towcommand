# Verification Before Completion

## Description

Mandatory evidence-based verification before claiming any task is complete. Prevents false completion claims, reduces rework, and enforces TowCommand's E2E evidence requirement.

## When to Use

- Before reporting any implementation as "done"
- Before marking any task/TODO as complete
- Before creating a commit or PR
- Before reporting test results

---

## Core Principle

**"Never claim completion without evidence. Evidence means tool output, not memory."**

### Why This Matters (2nd-Order)

1st-order: Prevents "I think it works" claims that waste reviewer time
2nd-order: Builds trust in agent output over time → reviewer spends less time re-checking → faster PR merges → faster sprint velocity

---

## Verification Protocol

### Level 1: Code Compilation

**Claim**: "The code is implemented"
**Required evidence**:
```bash
# Go
go build ./...           # Must show: no output (success)

# Mobile
cd apps/mobile && npx tsc --noEmit   # Must show: no errors
```

**NOT acceptable**: "I've written the code" without running the build

### Level 2: Tests Pass

**Claim**: "Tests are passing"
**Required evidence**:
```bash
# Go
go test -race ./internal/usecase/[package]/...   # Must show: ok, PASS

# Mobile
cd apps/mobile && pnpm test -- --filter=[test]   # Must show: Tests passed
```

**NOT acceptable**: "I've written the tests" without running them

### Level 3: Lint Clean

**Claim**: "Code is clean"
**Required evidence**:
```bash
# Go
golangci-lint run ./...   # Must show: no issues

# Mobile
cd apps/mobile && pnpm lint   # Must show: no errors
```

### Level 4: E2E Verified

**Claim**: "Feature works end-to-end"
**Required evidence**: Actual test output showing:
- HTTP status codes from API calls
- Screenshot or test trace from Playwright
- Database state verification (if applicable)

### Level 5: No Regressions

**Claim**: "Nothing else broke"
**Required evidence**:
```bash
# Full test suite
task test-unit   # All existing tests still pass
cd apps/mobile && pnpm test   # All existing tests still pass
```

---

## Completion Checklist

Before saying "done", verify ALL applicable items:

```markdown
### Implementation Verification
- [ ] `go build ./...` or `tsc --noEmit` succeeds (show output)
- [ ] New code has tests (show test file)
- [ ] Tests pass with `-race` flag (show output)
- [ ] Lint passes (show output)
- [ ] No regressions in existing tests (show full suite output)

### Integration Verification
- [ ] API contract matches what mobile app expects
- [ ] Event schema matches what subscribers expect
- [ ] DynamoDB key patterns are consistent with existing access patterns
- [ ] State machine transitions are valid for new states

### 2nd-Order Verification
- [ ] Change doesn't break downstream consumers I haven't tested
- [ ] Change doesn't introduce new failure modes without error handling
- [ ] Change doesn't create data that existing queries can't find (GSI coverage)
- [ ] Change doesn't add latency to hot paths without justification
```

---

## Anti-Patterns

| Anti-Pattern | Evidence of Violation | Correct Behavior |
|---|---|---|
| "Tests should pass" | No test output shown | Run tests, paste output |
| "I've fixed the lint error" | No lint output shown | Run lint, show clean output |
| "The build succeeds" | No build output shown | Run build, show success |
| "All tests pass" after modifying one file | Only ran tests for that file | Run full suite |
| "E2E verified" | No screenshots or HTTP responses | Show actual E2E output |
| "No regressions" | Only ran new tests | Run ALL existing tests |

---

## Recovery When Verification Fails

1. **Build fails** → Read the error, fix the code, re-verify from Level 1
2. **Tests fail** → Analyze failure (code bug vs test bug), fix, re-run
3. **Lint fails** → Fix lint issues (don't disable rules), re-run
4. **E2E fails** → Determine if it's code, deployment, or environment issue
5. **Regression found** → Fix the regression first, then re-verify everything

**Maximum 3 fix-and-retry cycles** before reporting the issue to the user with:
- What was attempted
- What the error is
- What the suspected root cause is
- Suggested next steps
