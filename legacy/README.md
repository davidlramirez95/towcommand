# Legacy TypeScript Codebase (Archived)

This directory contains the original TypeScript/Node.js implementation of TowCommand PH. It is preserved as a **read-only reference specification** for the Go migration.

## Purpose

- Reference for translating domain entities, business logic, and API contracts to Go
- Key types are in `packages/core/src/types/` — use these when building Go domain structs
- Service implementations in `services/` show the existing API behavior to replicate

## Structure

```
legacy/
  packages/        # Shared libraries (core types, DB, auth, cache, events)
  services/        # Lambda handlers (api-gateway, matching, notifications, etc.)
  scripts/         # Dev scripts (seeding, event docs)
  tests/           # Unit, integration, and e2e tests
```

## Important

- **DO NOT modify** files in this directory
- **DO NOT install** dependencies or run builds from here
- This directory will be deleted once the Go migration is complete and verified
