# Dispatching Parallel Agents

## Description

Pattern for handling multiple independent tasks by dispatching concurrent agents. Use when 3+ independent work items exist across different domains that can be solved in parallel without merge conflicts.

## When to Use

- Sprint execution with multiple independent PRs
- Multiple subsystems broken independently
- Batch test writing across unrelated packages
- Parallel research on independent topics

## When NOT to Use

- Tasks share Go packages or mobile screens
- Tasks modify shared types in `internal/domain/`
- Sequential dependency exists (one task's output feeds another)
- Tasks modify the same DynamoDB entity schema
- Tasks touch shared middleware or platform code

---

## Core Principle

**"One agent per independent domain boundary. Verify independence before dispatching."**

## 2nd-Order Dispatch Analysis

Before parallel dispatch, run this checklist:

### Independence Verification Matrix

| Check | If YES → | If NO → |
|---|---|---|
| Do tasks share any Go packages? | Sequential | Parallel OK |
| Do tasks modify shared interfaces? | Sequential, shared first | Parallel OK |
| Do tasks touch the same mobile screens? | Sequential | Parallel OK |
| Do tasks modify DynamoDB key schemas? | Sequential | Parallel OK |
| Could task A's changes make task B's tests fail? | Sequential | Parallel OK |
| Do tasks share Zustand stores? | Sequential or careful scoping | Parallel OK |

### Cascade Risk Assessment

For each parallel task, ask:
1. **If this agent makes a mistake, does it block the other agents?** → If yes, do this one first
2. **If this agent's changes are rejected in review, do other agents' changes still make sense?** → If no, they're not truly independent
3. **Does this agent need to read files that another agent is writing?** → If yes, sequential

---

## TowCommand Dispatch Patterns

### Pattern: Sprint PR Batch

```
# Safe to parallelize — different domain boundaries:
Agent 1: "Implement payment webhook handler" → touches cmd/payment-webhook/, internal/usecase/payment/
Agent 2: "Implement rating submission screen" → touches apps/mobile/app/booking/rate.tsx, stores/rating.ts
Agent 3: "Add SOS alert component" → touches apps/mobile/components/safety/, stores/safety.ts

# NOT safe — shared dependency:
Agent 1: "Add cancelledAt field to Booking entity" → touches internal/domain/booking.go
Agent 2: "Implement booking cancel handler" → depends on Agent 1's entity change
→ Do Agent 1 first, then Agent 2
```

### Pattern: Test Writing Batch

```
# Safe — tests for independent packages:
Agent 1: "Write tests for internal/usecase/payment/"
Agent 2: "Write tests for internal/usecase/rating/"
Agent 3: "Write tests for internal/usecase/safety/"

# NOT safe — tests share mocks:
Agent 1: "Write tests for internal/adapter/dynamodb/booking_repo.go"
Agent 2: "Write tests for internal/adapter/dynamodb/provider_repo.go"
→ These may share DynamoDB test helpers — check first
```

### Pattern: Backend + Mobile Parallel

```
# Usually safe — different languages, different directories:
Agent 1 (golang-pro): "Implement Go handler for /ratings endpoint"
Agent 2 (expo-mobile-dev): "Build rating submission screen"
→ Safe if they agree on the API contract first
→ 2nd-order risk: if agent 1 changes the response schema, agent 2's types break
→ Mitigation: define the contract in shared-types first, then spawn both
```

---

## Agent Task Template

Each spawned agent receives:

```markdown
## Task: [specific deliverable]

**Scope**: Only modify files in [directories]
**Agent type**: golang-pro | expo-mobile-dev | full-stack-implementer

**Deliverable**:
1. [file/feature to implement]
2. [tests to write]

**Constraints**:
- Do NOT modify files outside [scope]
- Do NOT change shared types in internal/domain/ or packages/shared-types/
- Do NOT modify Taskfile.yml, go.mod, or package.json
- Follow existing patterns in the target package

**Definition of Done**:
- Implementation complete
- Tests pass (go test -race or pnpm test)
- Lint clean (golangci-lint or pnpm lint)
- Files modified list returned
```

---

## Integration After Parallel Completion

1. **Collect all results** — `/spawn --collect`
2. **Check for file conflicts** — any overlapping modifications?
3. **Run full test suite** — `task test-unit && cd apps/mobile && pnpm test`
4. **Run lint** — `task lint && cd apps/mobile && pnpm lint`
5. **Merge branches** — one at a time, re-test after each merge
6. **Ship** — `/ship` for each PR

---

## Anti-Patterns

| Anti-Pattern | Why It Fails | Instead Do |
|---|---|---|
| Spawn 10+ agents at once | Resource exhaustion, context limits | Max 3-5 concurrent |
| Spawn without independence check | Merge conflicts, wasted work | Run verification matrix first |
| Skip post-merge testing | Agents' changes interact in ways neither predicted | Always run full suite after merging |
| Share Go packages between agents | Race to modify same files | One agent per package boundary |
| Forget to define API contracts first | Backend and mobile agents disagree on types | Define contract, then spawn both |
