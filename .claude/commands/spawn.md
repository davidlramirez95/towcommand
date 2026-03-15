# /spawn — Parallel Agent Dispatch

## Purpose

Launch background agents for parallel execution. Formalizes TowCommand's sprint batch-spawn workflow into a managed command.

## Usage

```
/spawn "[task description]"
/spawn --list
/spawn --collect
/spawn --cancel [id]
```

## Arguments

- `$ARGUMENTS`: Task description (quoted), or management flag

---

Manage parallel agents: **$ARGUMENTS**

## Operations

### Launch Task

```bash
/spawn "implement booking cancel handler with tests"
```

**Process:**
1. Analyze task for parallelizability — verify it's independent of other active tasks
2. Select the right agent type:
   - Go backend work → `golang-pro`
   - Mobile screens/components → `expo-mobile-dev`
   - Full feature from roadmap → `full-stack-implementer`
   - Research/exploration → `Explore` subagent
3. Launch agent in worktree isolation when the task involves file changes
4. Return task ID for tracking

**2nd-Order Dispatch Check:**
Before spawning, ask:
- Does this task share files with any other running task? → Sequential, not parallel
- Could this task's changes conflict with another agent's changes? → Use worktree isolation
- Does this task depend on output from another running task? → Wait, don't spawn yet
- Will this task modify shared types or interfaces? → Flag for merge conflict risk

### List Tasks

```bash
/spawn --list
```

Shows all active and recently completed agent tasks with status.

### Collect Results

```bash
/spawn --collect
```

Aggregates results from all completed agents. For each:
1. Summarize changes made
2. List files modified
3. Report test results
4. Flag any merge conflicts with other agent branches

### Cancel Task

```bash
/spawn --cancel [id]
```

## Agent Selection Matrix

| Task Pattern | Agent | Isolation |
|---|---|---|
| Go handler, repo, use case, test | `golang-pro` | worktree |
| Expo screen, component, hook, store | `expo-mobile-dev` | worktree |
| Full feature from roadmap/issue | `full-stack-implementer` | worktree |
| Research, exploration, analysis | `Explore` | none |
| Code review, security audit | `general-purpose` | none |

## Sprint Execution Pattern

For batch sprint work (our most common pattern):

```bash
# Phase 1: Dispatch all independent PRs
/spawn "PR #X: implement payment webhook handler"
/spawn "PR #Y: add rating submission screen"
/spawn "PR #Z: create SOS alert component"

# Phase 2: Monitor
/spawn --list

# Phase 3: Collect and integrate
/spawn --collect

# Phase 4: Ship each
/ship "feat(payment): add webhook handler"
```

## Conflict Prevention Rules

1. **One agent per package** — never spawn two agents that touch the same Go package
2. **One agent per screen** — never spawn two agents building the same Expo screen
3. **Shared types are sequential** — if a task modifies `internal/domain/`, do it first, then spawn dependents
4. **Always use worktree** for file-changing tasks — prevents direct branch conflicts

## Best Practices

- Spawn 3-5 agents max concurrently (context + resource limits)
- Use clear, specific task descriptions with scope boundaries
- Collect results promptly — don't let agent branches drift
- After `/spawn --collect`, run full test suite before merging any branch
