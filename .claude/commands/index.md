# /index — Project Structure Index

## Purpose

Generate a comprehensive project structure snapshot for faster navigation and context loading. Creates `PROJECT_INDEX.md` at the project root.

## Usage

```
/index
/index [directory]
```

## Arguments

- `$ARGUMENTS`: Optional directory to scope the index (default: entire project)

---

Generate project structure index for: **$ARGUMENTS**

## Process

1. **Scan project structure** — walk the directory tree, excluding:
   - `node_modules/`, `.git/`, `legacy/`, `vendor/`, `.claude/worktrees/`
   - Build artifacts, coverage reports, lock files

2. **Categorize by domain**:

   **Go Backend:**
   - `cmd/` — Lambda entry points (handler → which API endpoint)
   - `internal/domain/` — entities with key fields
   - `internal/usecase/` — business logic packages with method signatures
   - `internal/adapter/` — DynamoDB repos, Redis cache, EventBridge publisher
   - `internal/handler/` — HTTP/WS handler functions
   - `internal/platform/` — shared infrastructure (errors, config, middleware)

   **Mobile App:**
   - `apps/mobile/app/` — route structure (screens, tabs, groups)
   - `apps/mobile/components/` — UI components by category
   - `apps/mobile/hooks/` — custom hooks
   - `apps/mobile/stores/` — Zustand stores
   - `apps/mobile/lib/` — API client, WebSocket, theme, storage

   **Infrastructure:**
   - `infra/` — Terraform modules
   - `.github/workflows/` — CI/CD pipelines
   - `Taskfile.yml` — build commands

3. **Generate `PROJECT_INDEX.md`** with:
   - Directory tree (depth 3)
   - Package summaries (one line each)
   - API endpoint map (handler → route → method)
   - DynamoDB entity key schemas
   - Mobile screen inventory
   - Store → screen dependency map

4. **2nd-Order Index Value:**
   - Which packages have the most cross-cutting dependencies? (change risk)
   - Which screens consume which stores? (state coupling)
   - Which handlers share middleware? (blast radius)

## Output Format

```markdown
# TowCommand Project Index
Generated: 2026-03-15

## Go Backend (cmd/ + internal/)
### Lambda Handlers (41)
| Handler | Route | Method | Package |
|---------|-------|--------|---------|
| booking-create | /bookings | POST | booking |
...

### Domain Entities (8)
| Entity | PK | SK | GSIs |
|--------|----|----|------|
| Booking | BOOKING#id | BOOKING#id | GSI1, GSI2 |
...

### Use Cases (15 packages)
...

## Mobile App (apps/mobile/)
### Screens (18)
...

### Stores (6)
...

## Dependency Graph
### High-Risk Packages (most dependents)
...
```

## Best Practices

- Run `/index` at the start of major planning sessions
- Re-run after sprint completion to capture new packages
- Use the index to scope `/spawn` tasks and identify conflict zones
