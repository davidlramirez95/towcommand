# TowCommand PH — Backend Monorepo Architecture

## Overview

Serverless-first AWS backend using pnpm workspaces monorepo. Event-driven architecture via EventBridge, DynamoDB single-table design, deployed with Terraform IaC.

**Proven at scale:** 65M monthly requests, 20K concurrent RPS architecture patterns.

---

## Monorepo Structure

```
towcommand-backend/
├── package.json                    # Root workspace config
├── pnpm-workspace.yaml             # pnpm workspace definition
├── turbo.json                      # Turborepo pipeline config
├── tsconfig.base.json              # Shared TypeScript config
├── .env.example                    # Environment template
├── .eslintrc.js                    # Shared ESLint rules
├── .prettierrc                     # Code formatting
├── docker-compose.yml              # LocalStack + Redis + PostgreSQL
│
├── packages/
│   ├── core/                       # Shared business logic & types
│   │   ├── package.json
│   │   ├── tsconfig.json
│   │   └── src/
│   │       ├── index.ts
│   │       ├── types/              # Shared TypeScript interfaces
│   │       │   ├── booking.ts
│   │       │   ├── user.ts
│   │       │   ├── provider.ts
│   │       │   ├── vehicle.ts
│   │       │   ├── payment.ts
│   │       │   ├── events.ts       # EventBridge event schemas
│   │       │   └── index.ts
│   │       ├── constants/
│   │       │   ├── service-types.ts
│   │       │   ├── booking-status.ts
│   │       │   ├── regions.ts      # PH coverage areas
│   │       │   └── index.ts
│   │       ├── errors/             # Custom error classes
│   │       │   ├── app-error.ts
│   │       │   ├── validation-error.ts
│   │       │   └── index.ts
│   │       └── utils/
│   │           ├── geo.ts          # Haversine, geofencing
│   │           ├── pricing.ts      # Fare calculation engine
│   │           ├── otp.ts          # Digital Padala OTP gen
│   │           ├── validators.ts   # Zod schemas
│   │           └── index.ts
│   │
│   ├── db/                         # DynamoDB single-table + access patterns
│   │   ├── package.json
│   │   ├── tsconfig.json
│   │   └── src/
│   │       ├── index.ts
│   │       ├── client.ts           # DynamoDB Document Client singleton
│   │       ├── table-design.ts     # PK/SK patterns, GSI definitions
│   │       ├── entities/
│   │       │   ├── user.ts
│   │       │   ├── provider.ts
│   │       │   ├── booking.ts
│   │       │   ├── vehicle.ts
│   │       │   ├── rating.ts
│   │       │   ├── payment.ts
│   │       │   ├── suki-tier.ts    # Loyalty program
│   │       │   └── index.ts
│   │       ├── repositories/       # Data access layer
│   │       │   ├── base.repo.ts
│   │       │   ├── user.repo.ts
│   │       │   ├── provider.repo.ts
│   │       │   ├── booking.repo.ts
│   │       │   ├── vehicle.repo.ts
│   │       │   ├── rating.repo.ts
│   │       │   └── index.ts
│   │       └── migrations/         # Table schema versioning
│   │           └── v1-create-table.ts
│   │
│   ├── events/                     # EventBridge event catalog
│   │   ├── package.json
│   │   ├── tsconfig.json
│   │   └── src/
│   │       ├── index.ts
│   │       ├── publisher.ts        # EventBridge put helper
│   │       ├── schemas/            # Event JSON schemas
│   │       │   ├── booking.schema.ts
│   │       │   ├── provider.schema.ts
│   │       │   ├── payment.schema.ts
│   │       │   ├── tracking.schema.ts
│   │       │   └── notification.schema.ts
│   │       └── catalog.ts          # Event type registry
│   │
│   ├── cache/                      # Redis/ElastiCache wrapper
│   │   ├── package.json
│   │   ├── tsconfig.json
│   │   └── src/
│   │       ├── index.ts
│   │       ├── client.ts           # Redis connection
│   │       ├── patterns/
│   │       │   ├── geo-cache.ts    # Provider location caching
│   │       │   ├── session.ts      # User session cache
│   │       │   ├── rate-limiter.ts # API rate limiting
│   │       │   └── surge-pricing.ts
│   │       └── keys.ts             # Key naming conventions
│   │
│   └── auth/                       # Cognito helpers & middleware
│       ├── package.json
│       ├── tsconfig.json
│       └── src/
│           ├── index.ts
│           ├── cognito-client.ts
│           ├── middleware/
│           │   ├── jwt-verify.ts   # API Gateway authorizer
│           │   ├── rbac.ts         # Role-based access
│           │   └── ban-check.ts
│           └── utils/
│               ├── token-claims.ts
│               └── social-link.ts  # Account linking logic
│
├── services/                       # Lambda function services
│   ├── api-gateway/                # REST API handlers
│   │   ├── package.json
│   │   ├── tsconfig.json
│   │   ├── serverless.yml          # or SAM template
│   │   └── src/
│   │       ├── handlers/
│   │       │   ├── booking/
│   │       │   │   ├── create.ts
│   │       │   │   ├── cancel.ts
│   │       │   │   ├── get.ts
│   │       │   │   ├── list.ts
│   │       │   │   └── update-status.ts
│   │       │   ├── provider/
│   │       │   │   ├── register.ts
│   │       │   │   ├── update-location.ts
│   │       │   │   ├── toggle-availability.ts
│   │       │   │   └── get-nearby.ts
│   │       │   ├── user/
│   │       │   │   ├── profile.ts
│   │       │   │   ├── vehicles.ts
│   │       │   │   └── preferences.ts
│   │       │   ├── diagnosis/
│   │       │   │   ├── analyze.ts  # AI symptom matching
│   │       │   │   └── history.ts
│   │       │   ├── payment/
│   │       │   │   ├── initiate.ts
│   │       │   │   ├── webhook.ts  # GCash/Maya callbacks
│   │       │   │   └── receipt.ts
│   │       │   └── rating/
│   │       │       ├── submit.ts
│   │       │       └── get.ts
│   │       └── middleware/
│   │           ├── error-handler.ts
│   │           ├── cors.ts
│   │           ├── request-logger.ts
│   │           └── validation.ts
│   │
│   ├── websocket/                  # Real-time WebSocket API
│   │   ├── package.json
│   │   ├── tsconfig.json
│   │   └── src/
│   │       ├── handlers/
│   │       │   ├── connect.ts
│   │       │   ├── disconnect.ts
│   │       │   ├── location-update.ts  # Driver GPS stream
│   │       │   ├── booking-status.ts   # Status push
│   │       │   └── chat-message.ts     # In-app messaging
│   │       └── lib/
│   │           ├── connection-manager.ts
│   │           └── broadcast.ts
│   │
│   ├── matching/                   # Provider matching engine
│   │   ├── package.json
│   │   ├── tsconfig.json
│   │   └── src/
│   │       ├── handler.ts          # EventBridge subscriber
│   │       ├── algorithms/
│   │       │   ├── nearest.ts      # Haversine + availability
│   │       │   ├── weighted-score.ts  # Rating + distance + price
│   │       │   └── surge-aware.ts  # Typhoon mode pricing
│   │       └── lib/
│   │           ├── geo-search.ts   # Redis GEOSEARCH
│   │           └── timeout.ts      # Match timeout + escalation
│   │
│   ├── notifications/              # Push, SMS, Email
│   │   ├── package.json
│   │   ├── tsconfig.json
│   │   └── src/
│   │       ├── handler.ts          # EventBridge subscriber
│   │       ├── channels/
│   │       │   ├── sms.ts          # AWS SNS / Semaphore API
│   │       │   ├── push.ts         # FCM / APNs via SNS
│   │       │   └── email.ts        # SES
│   │       └── templates/
│   │           ├── booking-confirmed.ts
│   │           ├── driver-arriving.ts
│   │           ├── otp-code.ts     # Filipino SMS template
│   │           └── sos-alert.ts
│   │
│   ├── auth-triggers/              # Cognito Lambda triggers
│   │   ├── package.json
│   │   ├── tsconfig.json
│   │   └── src/
│   │       ├── pre-signup.ts       # Auto-link social accounts
│   │       ├── post-confirmation.ts # DynamoDB user sync
│   │       ├── pre-token.ts        # RBAC claims injection
│   │       ├── custom-message.ts   # Filipino SMS templates
│   │       └── pre-authentication.ts # Ban checking
│   │
│   └── analytics/                  # PostgreSQL analytics sidecar
│       ├── package.json
│       ├── tsconfig.json
│       └── src/
│           ├── handler.ts          # EventBridge → PostgreSQL sink
│           ├── queries/
│           │   ├── revenue-report.ts
│           │   ├── provider-performance.ts
│           │   └── demand-heatmap.ts
│           └── lib/
│               ├── pg-client.ts    # PostgreSQL connection
│               └── schema.sql      # Analytics tables
│
├── infra/                          # Terraform IaC
│   ├── modules/
│   │   ├── dynamodb/
│   │   │   ├── main.tf             # Single table + GSIs
│   │   │   ├── variables.tf
│   │   │   └── outputs.tf
│   │   ├── cognito/
│   │   │   ├── main.tf             # User pool + identity providers
│   │   │   ├── triggers.tf         # Lambda trigger associations
│   │   │   └── variables.tf
│   │   ├── api-gateway/
│   │   │   ├── rest.tf             # REST API
│   │   │   ├── websocket.tf        # WebSocket API
│   │   │   ├── authorizer.tf
│   │   │   └── variables.tf
│   │   ├── lambda/
│   │   │   ├── main.tf             # Function definitions (arm64 Graviton)
│   │   │   ├── layers.tf           # Shared Lambda layers
│   │   │   └── variables.tf
│   │   ├── eventbridge/
│   │   │   ├── main.tf             # Event bus + rules
│   │   │   ├── schemas.tf          # Event schema registry
│   │   │   └── variables.tf
│   │   ├── elasticache/
│   │   │   ├── main.tf             # Redis cluster
│   │   │   └── variables.tf
│   │   ├── rds/
│   │   │   ├── main.tf             # PostgreSQL (analytics)
│   │   │   └── variables.tf
│   │   ├── s3/
│   │   │   ├── main.tf             # Uploads, photos, documents
│   │   │   └── variables.tf
│   │   ├── monitoring/
│   │   │   ├── cloudwatch.tf       # Alarms, dashboards
│   │   │   ├── xray.tf             # Distributed tracing
│   │   │   └── variables.tf
│   │   └── vpc/
│   │       ├── main.tf
│   │       └── variables.tf
│   │
│   └── environments/
│       ├── dev/
│       │   ├── main.tf
│       │   ├── terraform.tfvars
│       │   └── backend.tf          # S3 state backend
│       ├── staging/
│       │   ├── main.tf
│       │   ├── terraform.tfvars
│       │   └── backend.tf
│       └── prod/
│           ├── main.tf
│           ├── terraform.tfvars
│           └── backend.tf
│
├── scripts/
│   ├── seed-db.ts                  # Dev data seeding
│   ├── deploy.sh                   # CI/CD deployment
│   ├── local-setup.sh              # LocalStack bootstrap
│   └── generate-event-docs.ts      # Auto-gen event catalog docs
│
└── tests/
    ├── unit/                       # Vitest unit tests
    ├── integration/                # DynamoDB + EventBridge integration
    └── e2e/                        # API endpoint tests
```

---

## Tech Stack Summary

| Layer | Technology | Justification |
|-------|-----------|---------------|
| Runtime | Node.js 20 (TypeScript) | Primary API/business logic |
| ML Services | Python 3.12 | Future AI diagnosis, risk scoring |
| Database | DynamoDB (single-table) | Scales to 65M+ requests, <10ms reads |
| Cache | ElastiCache Redis 7 | Geo queries, session, rate limiting |
| Analytics DB | PostgreSQL 16 (RDS) | Complex queries, reporting |
| Auth | Cognito + Social SSO | Google/Facebook/Apple + Phone OTP |
| Events | EventBridge | Decoupled, extensible event bus |
| Real-time | API Gateway WebSocket | Live GPS tracking, chat |
| Storage | S3 + CloudFront | Photos, documents, static assets |
| IaC | Terraform | Multi-env, modular, state management |
| Compute | Lambda arm64 (Graviton) | 34% better price-performance |
| Monorepo | pnpm + Turborepo | Fast builds, dependency deduplication |
| Testing | Vitest + Supertest | Unit, integration, e2e |
| CI/CD | GitHub Actions | Automated deploy pipeline |

---

## DynamoDB Single-Table Design

### Table: `TowCommand-{env}`

**14 Entities, 5 GSIs**

| Entity | PK | SK | GSI1-PK | GSI1-SK |
|--------|----|----|---------|---------|
| User | `USER#{userId}` | `PROFILE` | `EMAIL#{email}` | `USER` |
| Provider | `PROV#{providerId}` | `PROFILE` | `STATUS#{status}` | `REGION#{regionCode}` |
| Vehicle | `USER#{userId}` | `VEH#{vehicleId}` | | |
| Booking | `BOOK#{bookingId}` | `META` | `USER#{userId}` | `BOOK#{createdAt}` |
| BookingProvider | `BOOK#{bookingId}` | `PROV#{providerId}` | `PROV#{providerId}` | `BOOK#{createdAt}` |
| Rating | `BOOK#{bookingId}` | `RATING` | `PROV#{providerId}` | `RATE#{createdAt}` |
| Payment | `BOOK#{bookingId}` | `PAY#{paymentId}` | `PAY_STATUS#{status}` | `PAY#{createdAt}` |
| SukiTier | `USER#{userId}` | `SUKI` | | |
| OTP | `BOOK#{bookingId}` | `OTP` | | |
| SOSAlert | `SOS#{alertId}` | `META` | `REGION#{regionCode}` | `SOS#{createdAt}` |
| ChatMessage | `BOOK#{bookingId}` | `MSG#{timestamp}` | | |
| ProviderDoc | `PROV#{providerId}` | `DOC#{docType}` | | |
| ServiceArea | `REGION#{regionCode}` | `AREA` | | |
| AuditLog | `AUDIT#{entityId}` | `LOG#{timestamp}` | | |

### GSI Definitions

| GSI | Purpose | PK | SK |
|-----|---------|----|----|
| GSI1 | User lookups by email, booking by user | GSI1PK | GSI1SK |
| GSI2 | Provider by status + region | GSI2PK | GSI2SK |
| GSI3 | Payment by status | GSI3PK | GSI3SK |
| GSI4 | Rating by provider | GSI4PK | GSI4SK |
| GSI5 | SOS by region (emergency queries) | GSI5PK | GSI5SK |

---

## EventBridge Event Catalog

### Event Bus: `towcommand-{env}`

| Source | Detail Type | Triggers |
|--------|------------|----------|
| `tc.booking` | `BookingCreated` | Matching engine, notifications |
| `tc.booking` | `BookingAccepted` | User notification, tracking init |
| `tc.booking` | `BookingCancelled` | Provider release, refund flow |
| `tc.booking` | `BookingCompleted` | Payment capture, rating prompt, Suki points |
| `tc.matching` | `ProviderMatched` | OTP generation, WebSocket push |
| `tc.matching` | `MatchTimeout` | Escalation, expanded search |
| `tc.tracking` | `LocationUpdated` | WebSocket broadcast, ETA recalc |
| `tc.tracking` | `DriverArrived` | User notification, OTP verify prompt |
| `tc.payment` | `PaymentInitiated` | Payment gateway call |
| `tc.payment` | `PaymentCompleted` | Receipt generation, provider payout |
| `tc.payment` | `PaymentFailed` | Retry logic, user notification |
| `tc.sos` | `SOSActivated` | Ops center alert, PNP-HPG notify, nearby providers |
| `tc.auth` | `UserRegistered` | Welcome flow, Suki tier init |
| `tc.provider` | `ProviderOnline` | Availability index update |
| `tc.provider` | `ProviderOffline` | Remove from matching pool |

### V1.1 Extension Points (Zero Code Changes)

New subscribers to existing events:
- `BookingCompleted` → Blockchain receipt writer
- `BookingCreated` → ML risk scoring
- `SOSActivated` → Typhoon Mode orchestrator
- `PaymentCompleted` → BNPL provider integration
- `ProviderMatched` → Mechanic marketplace router

---

## API Endpoints

### REST API: `api.towcommand.ph/v1`

**Auth** — Cognito JWT required (except public routes)

| Method | Path | Handler | Auth |
|--------|------|---------|------|
| POST | `/bookings` | Create booking | User |
| GET | `/bookings/:id` | Get booking details | User/Provider |
| PATCH | `/bookings/:id/status` | Update status | Provider/Admin |
| DELETE | `/bookings/:id` | Cancel booking | User |
| GET | `/bookings` | List user bookings | User |
| POST | `/diagnosis` | AI symptom analysis | User |
| GET | `/providers/nearby` | Get nearby providers | User |
| PATCH | `/providers/location` | Update GPS position | Provider |
| PATCH | `/providers/availability` | Toggle online/offline | Provider |
| POST | `/providers/register` | Provider onboarding | Public |
| GET | `/users/profile` | Get own profile | User |
| PATCH | `/users/profile` | Update profile | User |
| POST | `/users/vehicles` | Add vehicle | User |
| POST | `/payments/initiate` | Start payment | User |
| POST | `/payments/webhook` | Gateway callback | Internal |
| POST | `/ratings` | Submit rating | User |
| POST | `/sos` | Activate SOS | User |

### WebSocket API: `wss://ws.towcommand.ph`

| Route | Direction | Payload |
|-------|-----------|---------|
| `$connect` | Client→Server | JWT auth token |
| `$disconnect` | Client→Server | Cleanup |
| `location.update` | Provider→Server | `{lat, lng, heading, speed}` |
| `location.broadcast` | Server→Client | `{providerId, lat, lng, eta}` |
| `booking.status` | Server→Client | `{bookingId, status, metadata}` |
| `chat.send` | Client→Server | `{bookingId, message}` |
| `chat.receive` | Server→Client | `{bookingId, senderId, message}` |
| `sos.alert` | Server→Client | `{alertId, location, type}` |

---

## Cost Estimate (MVP)

| Resource | Dev | Staging | Prod (1K users) | Prod (10K users) |
|----------|-----|---------|-----------------|------------------|
| Lambda | Free tier | ~$20 | ~$50 | ~$200 |
| DynamoDB | Free tier | ~$15 | ~$40 | ~$150 |
| ElastiCache | ~$15 | ~$30 | ~$60 | ~$120 |
| RDS PostgreSQL | ~$15 | ~$30 | ~$50 | ~$100 |
| API Gateway | Free tier | ~$10 | ~$30 | ~$100 |
| Cognito | Free tier | Free tier | ~$20 | ~$50 |
| S3 + CloudFront | ~$5 | ~$10 | ~$20 | ~$50 |
| EventBridge | Free tier | ~$5 | ~$10 | ~$30 |
| **Total** | **~$35/mo** | **~$120/mo** | **~$280/mo** | **~$800/mo** |
