# TowCommand PH — Development Guide

## Quick Start

### Local Setup
```bash
bash scripts/local-setup.sh
```

This will:
- Verify Node.js 20+, Docker, and pnpm are installed
- Create `.env` from `.env.example`
- Start Docker services (LocalStack, Redis, PostgreSQL)
- Wait for all services to be ready
- Run database migrations
- Install dependencies
- Build all packages

### Available Services
- **LocalStack**: http://localhost:4566 (AWS emulation)
- **Redis**: localhost:6379
- **Redis Commander**: http://localhost:8081
- **PostgreSQL**: localhost:5432

## Scripts

### `scripts/local-setup.sh`
Initializes the complete local development environment with all Docker services and dependencies.

### `scripts/deploy.sh <stage>`
Deploys to dev, staging, or prod environments.
- Stages: `dev`, `staging`, `prod`
- Runs full test suite before deployment
- Requires AWS credentials via `AWS_ROLE_ARN`

### `scripts/seed-db.ts`
Seeds DynamoDB with test data:
- Sample customer user
- Sample vehicles
- Sample providers (verified & basic tiers)
- Service areas

Run with: `pnpm run seed-db`

### `scripts/generate-event-docs.ts`
Auto-generates event catalog documentation from event sources.

Run with: `pnpm run generate-event-docs`

## CI/CD Workflows

### `.github/workflows/ci.yml`
Main continuous integration pipeline that runs on:
- Pull requests to `main` or `develop`
- Pushes to `main`

Jobs:
1. **Lint & Type Check** — pnpm lint, typecheck
2. **Unit Tests** — Full coverage report
3. **Integration Tests** — With LocalStack and Redis
4. **Security Scan** — Dependency audit
5. **Deploy to Dev** — Automatic on main push (if all jobs pass)

### `.github/workflows/deploy-staging.yml`
Manual staging deployments via:
- Workflow dispatch (manual trigger)
- Tags matching `v*-rc*` (release candidates)

## Test Structure

### `tests/unit/`
Unit tests for individual functions/modules.
- Run: `pnpm run test:unit`
- Coverage: `pnpm run test:unit -- --coverage`

### `tests/integration/`
Integration tests with real services:
- LocalStack (DynamoDB, S3, SQS, SNS, EventBridge)
- Redis
- PostgreSQL (when added)
- Run: `pnpm run test:integration`

### `tests/e2e/`
End-to-end tests for complete workflows:
- API endpoint tests
- Full booking flow
- Payment processing
- Run: `pnpm run test:e2e`

## Event Catalog

Events are organized by domain:
- `tc.booking` — Booking lifecycle events
- `tc.matching` — Provider matching events
- `tc.tracking` — Location and tracking events
- `tc.payment` — Payment processing events
- `tc.sos` — Emergency SOS events
- `tc.auth` — Authentication events
- `tc.provider` — Provider lifecycle events
- `tc.evidence` — Evidence and dispute events

View full catalog: `docs/event-catalog.md`

## Database Schema

Using DynamoDB with complex GSI structure:
- Primary table: `TowCommand-<stage>`
- Entities: User, Provider, Booking, Payment, Rating, SukiTier
- Multiple GSIs for efficient querying by email, phone, tier, location, etc.

See: `packages/db/src/table-design.ts`

## Environment Variables

Copy `.env.example` to `.env`:
```bash
# AWS Configuration
AWS_REGION=ap-southeast-1
DYNAMODB_ENDPOINT=http://localhost:4566
DYNAMODB_TABLE_NAME=TowCommand-dev

# Redis
REDIS_HOST=localhost
REDIS_PORT=6379

# Services
EVENT_BUS_NAME=towcommand-dev
STAGE=dev
```

## Development Workflow

1. Create feature branch from `develop`
2. Make changes and commit
3. Push branch and create PR to `develop`
4. CI runs automatically
5. After approval, merge to `develop`
6. Create RC tag (e.g., `v1.2.3-rc1`) for staging deployment
7. After QA, merge to `main` for dev deployment
8. After validation, use `scripts/deploy.sh prod` for production

## Troubleshooting

### Services not starting
```bash
# Check Docker
docker ps
docker logs <service-name>

# Clean up and restart
docker-compose down
docker-compose up -d
```

### LocalStack issues
```bash
# Check LocalStack health
docker-compose exec localstack awslocal sts get-caller-identity

# Create tables
pnpm run db:migrate
```

### Redis issues
```bash
# Check Redis
docker-compose exec redis redis-cli ping

# Clear cache
docker-compose exec redis redis-cli FLUSHALL
```

## Contributing

- Use TypeScript for all new code
- Follow ESLint rules (checked in CI)
- Write tests for new features
- Update DEVELOPMENT.md if adding new scripts/workflows
- Ensure all tests pass before pushing
