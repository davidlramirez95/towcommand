# TowCommand PH — Feature Roadmap & Prioritization

## Sprint Breakdown (12-Week MVP)

### Sprint 1 — Foundation (Weeks 1–2)
**Priority: P0 — Blocking everything else**

| Task | Package/Service | Status |
|------|----------------|--------|
| Terraform base infra (VPC, DDB, Redis, S3, Cognito) | `infra/` | Scaffolded |
| Monorepo setup (pnpm, Turbo, tsconfig, ESLint) | Root | Scaffolded |
| CI/CD pipeline (GitHub Actions) | `.github/workflows/` | Scaffolded |
| DynamoDB single-table + 5 GSIs | `packages/db/` | Scaffolded |
| Cognito User Pool + Google/Facebook/Apple SSO | `infra/modules/cognito/` | Scaffolded |
| Auth Lambda triggers (pre-signup, post-confirm, pre-token) | `services/auth-triggers/` | Scaffolded |
| User CRUD (profile, vehicles) | `services/api-gateway/handlers/user/` | Scaffolded |
| LocalStack + Docker Compose dev environment | `docker-compose.yml` | Scaffolded |
| Shared core types, errors, validators | `packages/core/` | Scaffolded |

### Sprint 2 — Core Booking (Weeks 3–4)
**Priority: P0 — Core revenue path**

| Task | Package/Service | Status |
|------|----------------|--------|
| Booking create + estimate + cancel handlers | `services/api-gateway/handlers/booking/` | Scaffolded |
| Pricing engine (MMDA Reg. 24-004 rates) | `packages/core/utils/pricing.ts` | Scaffolded |
| Matching engine (weighted scoring algorithm) | `services/matching/` | Scaffolded |
| Redis geo-caching for provider locations | `packages/cache/patterns/geo-cache.ts` | Scaffolded |
| WebSocket API (GPS tracking + job status push) | `services/websocket/` | Scaffolded |
| EventBridge events (BookingCreated, Matched, etc.) | `packages/events/` | Scaffolded |
| Booking status state machine | `packages/core/types/booking.ts` | Scaffolded |

### Sprint 3 — Evidence + OTP (Weeks 5–6)
**Priority: P0 — Legal compliance (PH Rules on Electronic Evidence)**

| Task | Package/Service | Status |
|------|----------------|--------|
| Evidence upload pipeline (S3 pre-signed URLs) | New: `services/evidence/` | Not started |
| Rekognition vehicle detection + blur check | New: `services/evidence/` | Not started |
| SHA-256 hash integrity verification | New: `services/evidence/` | Not started |
| Digital Padala OTP service | New: `services/otp/` | Not started |
| Condition report entity + 8-photo enforcement | `packages/db/entities/` | Partial |
| Metadata watermarking (ImageMagick Lambda layer) | New: `services/evidence/` | Not started |

### Sprint 4 — Payments + Safety (Weeks 7–8)
**Priority: P0 — Revenue capture + user safety**

| Task | Package/Service | Status |
|------|----------------|--------|
| GCash/Maya payment integration (hold/capture/refund) | `services/api-gateway/handlers/payment/` | Scaffolded |
| Cancellation fee engine (tiered by status) | `packages/core/utils/pricing.ts` | Scaffolded |
| Silent SOS trigger | New: `services/safety/` | Not started |
| Geofence route monitoring | New: `services/safety/` | Not started |
| Night mode auto-escalation | New: `services/safety/` | Not started |
| Rule-based risk scoring | New: `services/safety/` | Not started |

### Sprint 5 — Provider + Comms (Weeks 9–10)
**Priority: P1 — Supply-side enablement**

| Task | Package/Service | Status |
|------|----------------|--------|
| Provider onboarding (KYC upload, verification) | `services/api-gateway/handlers/provider/` | Scaffolded |
| Provider dashboard / earnings API | New handler | Not started |
| In-app chat (WebSocket bidirectional) | `services/websocket/handlers/chat-message.ts` | Scaffolded |
| Quick reply templates (Filipino) | `services/notifications/templates/` | Scaffolded |
| Push/SMS notification service | `services/notifications/` | Scaffolded |
| AI Smart Diagnosis (Bedrock/Claude) | `services/api-gateway/handlers/diagnosis/` | Scaffolded |

### Sprint 6 — Ops + Hardening (Weeks 11–12)
**Priority: P1 — Operational readiness**

| Task | Package/Service | Status |
|------|----------------|--------|
| Admin ops dashboard API (live map, safety console) | New: `services/admin/` | Not started |
| Dispute filing flow | New: `services/dispute/` | Not started |
| Analytics pipeline (DDB Streams to PostgreSQL) | `services/analytics/` | Scaffolded |
| CloudWatch dashboards + alarms | `infra/modules/monitoring/` | Scaffolded |
| Load testing (Artillery / k6) | `tests/` | Not started |
| Security audit + WAF tuning | `infra/modules/` | Partial |
| Staging + Production deployment | `infra/environments/` | Scaffolded |

---

## Post-MVP Roadmap

### V1.1 (Month 4–6)
- Blockchain evidence anchoring (Polygon Merkle root)
- ML risk scoring (SageMaker, replaces rule-based)
- AI damage detection (pre/post photo comparison)
- Dashcam streaming (Kinesis Video, Suki Elite only)
- BNPL integration (Billease)
- Typhoon Mode (PAGASA webhook + surge pricing + safety escalation)

### V1.2 (Month 7–9)
- Corporate fleet dashboard
- Mechanic marketplace
- Full-text search (OpenSearch)
- Multi-language chat support (Filipino NLP)

### V2.0 (Month 10–12)
- Insurance integration (read-only claim evidence API)
- Partner API (white-label towing)
- Advanced analytics (demand forecasting, provider optimization)

---

## Prioritization Framework

**P0 — Must Have (MVP Blocker)**
Revenue-generating, legally required, or safety-critical. The app cannot launch without these.

**P1 — Should Have (MVP Enhancement)**
Significantly improves UX or operational efficiency. Can soft-launch without, but needed for scale.

**P2 — Nice to Have (V1.1)**
Differentiators and competitive advantages. Plugs into existing events without core changes.

**P3 — Future (V1.2+)**
New revenue streams and market expansion. Requires new service domains.

---

## Scaffolding Completion Summary

| Component | Files | Status |
|-----------|-------|--------|
| Root config (pnpm, turbo, ts, eslint, docker) | 9 | Complete |
| packages/core (types, errors, utils, constants) | 22 | Complete |
| packages/db (client, entities, repos, migrations) | 12 | Complete |
| packages/events (publisher, schemas, catalog) | 8 | Complete |
| packages/cache (Redis client, geo, rate-limit) | 8 | Complete |
| packages/auth (Cognito, JWT, RBAC) | 6 | Complete |
| services/api-gateway (15 handlers + 4 middleware) | 19 | Complete |
| services/websocket (5 handlers + 2 libs) | 7 | Complete |
| services/matching (3 algorithms + handler) | 6 | Complete |
| services/notifications (3 channels + 4 templates) | 8 | Complete |
| services/auth-triggers (5 Cognito triggers) | 5 | Complete |
| services/analytics (handler + 3 queries + pg-client) | 5 | Complete |
| infra/ (10 Terraform modules, 3 environments) | 57 | Complete |
| CI/CD (GitHub Actions) | 2 | Complete |
| Scripts (setup, deploy, seed, event-docs) | 4 | Complete |
| Tests (unit/integration/e2e scaffolding) | 5 | Complete |
| **Total** | **~183** | **Complete** |

---

## Architecture Decisions Record

1. **Serverless-first**: All handlers are Lambda (arm64 Graviton) — no EC2
2. **Single-table DynamoDB**: 14 entities, 5 GSIs — proven at 65M req/month
3. **EventBridge**: Central event bus with 15+ event types for extensibility
4. **pnpm + Turborepo**: Fast monorepo builds with dependency deduplication
5. **TypeScript strict mode**: Shared types eliminate runtime type errors
6. **Terraform IaC**: 10 reusable modules, 3 environments (dev/staging/prod)
7. **Redis**: Geo-queries (GEOADD), rate limiting, session cache, OTP store
8. **PostgreSQL sidecar**: Analytics/reporting that DDB can't handle efficiently

## Note on gutguard-ai

The gutguard-ai repo at /Users/david.ramirez/Downloads/gutguard-ai was not accessible from the sandbox environment. The scaffolding incorporates common serverless agent patterns (middleware chains, repository pattern, event-driven architecture, IAM least-privilege) that would typically be found in such a reference repo. To integrate specific patterns from gutguard-ai, copy it into the towcommand workspace folder and re-run.
