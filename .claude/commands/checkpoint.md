# /checkpoint — Session State Management

## Purpose

Save and restore working context for multi-session complex tasks. Uses git stash + metadata for state preservation across Claude Code sessions.

## Usage

```
/checkpoint save [name]
/checkpoint list
/checkpoint restore [name]
/checkpoint delete [name]
```

## Arguments

- `$ARGUMENTS`: Operation (save/list/restore/delete) + checkpoint name

---

Manage session checkpoints: **$ARGUMENTS**

## Operations

### Save Checkpoint

```bash
/checkpoint save feature-auth
```

**Process:**
1. Capture current git state (branch, stash uncommitted changes)
2. Record working context metadata to `.claude/checkpoints/[name].json`:
   ```json
   {
     "name": "feature-auth",
     "created": "2026-03-15T10:30:00Z",
     "branch": "feat/auth-module",
     "git_stash": "stash@{0}",
     "files_in_progress": ["internal/usecase/auth/login.go", "cmd/auth-login/main.go"],
     "current_task": "Implementing JWT refresh token rotation",
     "tests_status": "3/5 passing",
     "next_steps": ["Fix token expiry edge case", "Add integration test"],
     "notes": ""
   }
   ```
3. Confirm checkpoint saved

**2nd-Order Save Check:**
- Am I checkpointing *enough* context? Will future-me understand where I left off?
- Are there uncommitted changes that would be lost without the stash?
- Is the `next_steps` field specific enough to resume without re-reading everything?

### List Checkpoints

```bash
/checkpoint list
```

Shows all saved checkpoints with age, branch, and task summary.

### Restore Checkpoint

```bash
/checkpoint restore feature-auth
```

**Process:**
1. Read checkpoint metadata
2. Switch to the saved branch
3. Apply git stash if one was saved
4. Display context summary: what was being worked on, what's next
5. Ready to continue

### Delete Checkpoint

```bash
/checkpoint delete feature-auth
```

Removes metadata file and drops associated git stash.

## When to Checkpoint

- Before switching to a different feature/PR
- Before attempting a risky refactor
- At natural breakpoints in multi-phase work
- Before ending a session on incomplete work
- Before running destructive operations

## Storage

```
.claude/checkpoints/
  feature-auth.json
  booking-refactor.json
```

## Best Practices

1. **Name descriptively** — `payment-webhook-wip` not `checkpoint-1`
2. **Include next steps** — the most valuable part of a checkpoint
3. **Checkpoint before context switches** — cheaper than re-deriving context
4. **Clean up old checkpoints** — delete after the work is merged
5. **Commit before checkpointing when possible** — stashes are fragile
