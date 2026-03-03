# TowCommand PH

Philippine tow truck and roadside assistance platform — serverless backend.

## Architecture

Go backend following Clean Architecture + SOLID + 12-Factor principles, deployed as AWS Lambda functions.

```
cmd/                    # Lambda entry points (composition roots)
internal/
  domain/               # Entities + value objects (zero external deps)
  usecase/              # Application services + port interfaces
  adapter/              # DynamoDB repos, Redis cache, EventBridge, handlers
  platform/             # Config, AWS clients, logger
infra/                  # Terraform IaC
legacy/                 # Archived TypeScript codebase (read-only reference)
```

## Tech Stack

- **Runtime:** Go 1.22+, AWS Lambda (provided.al2023, arm64)
- **Database:** DynamoDB (single-table design)
- **Cache:** Redis (ElastiCache)
- **Events:** EventBridge
- **Auth:** Cognito
- **IaC:** Terraform

## Development

```bash
# Build all Lambda functions
GOOS=linux GOARCH=arm64 CGO_ENABLED=0 go build -tags lambda.norpc -o bootstrap ./cmd/<func>/

# Run tests
go test ./...

# Lint
golangci-lint run ./...

# Task runner
task build   # see Taskfile.yml
```

## Legacy Reference

The `legacy/` directory contains the original TypeScript implementation, preserved as a read-only reference for the Go migration. See [legacy/README.md](legacy/README.md) for details.
