# TowCommand PH — Remaining Backend + Mobile App Plan

## Context

**Sprints 0-4 COMPLETE** (25 PRs merged, #30-#54). The Go backend has 37 Lambda handlers, 9 DynamoDB repos, 13 use case packages, 8 domain entities, full payment/rating/safety/admin systems. All P0 (MVP-blocking) backend work is done.

**What's needed now:**
1. A slim "Sprint 5 Backend" pass to fill gaps that the mobile app will hit immediately
2. A full mobile app (Expo + React Native) — the actual product users interact with
3. Custom subagent definitions to parallelize development as a solo dev + AI agents

**User preferences:**
- Solo developer with AI agents (Claude Code subagents)
- Defer admin dashboard, focus on mobile first
- 4-6 week timeline to working MVP
- Best flexibility/interoperability → Expo + React Native
- **Same monorepo** — mobile app in `apps/mobile/` alongside Go backend

---

## Part 1: Sprint 5 Backend (Week 1 — ~5 days)

### Why
The mobile app needs 4 backend capabilities that don't exist yet:
1. **AI Smart Diagnosis** — the "Diagnose" screen sends photos/description, gets service recommendation
2. **Push Notifications** — mobile needs FCM/APNs token registration + push delivery
3. **Condition Report Handler** — evidence upload pipeline exists but no API handler to create the condition report entity
4. **Provider Earnings Endpoint** — provider app needs to show earnings summary

### Sprint 5 Issues (4 parallel golang-pro agents)

#### Issue 5A: AI Smart Diagnosis (~400 LOC)
```
internal/usecase/diagnosis/ports.go          — DiagnosisEngine interface
internal/usecase/diagnosis/diagnose.go       — DiagnoseUseCase (orchestrates Bedrock call)
internal/usecase/diagnosis/diagnose_test.go
internal/adapter/gateway/bedrock_diagnosis.go — BedrockDiagnosisEngine (Claude Sonnet via AWS SDK)
internal/adapter/gateway/bedrock_diagnosis_test.go
cmd/api-diagnosis/main.go                    — POST /diagnosis
internal/adapter/handler/diagnosis.go        — DiagnoseHandler
internal/adapter/handler/diagnosis_test.go
```
**Flow:** Customer uploads photo(s) + text description -> Bedrock Claude Sonnet analyzes -> returns recommended service type, urgency level, estimated cost range.

#### Issue 5B: Push Notifications (~350 LOC)
```
internal/usecase/port/push.go               — PushTokenRegistrar, PushSender interfaces
internal/adapter/repository/dynamo_push.go   — DynamoDB push token store (PK: PUSH#{userId})
internal/adapter/repository/dynamo_push_test.go
internal/adapter/gateway/sns_push.go         — SNS Platform Application push sender
internal/adapter/gateway/sns_push_test.go
cmd/api-push-register/main.go               — POST /users/{id}/push-token
internal/adapter/handler/push_register.go    — RegisterPushTokenHandler
internal/adapter/handler/push_register_test.go
```
**Flow:** Mobile registers FCM/APNs token on login -> stored in DDB -> notification router dispatches to SNS platform endpoint -> push notification delivered.

**Modify:** `internal/usecase/notification/router.go` — add push notification channel alongside SMS/Email.

#### Issue 5C: Condition Report Handler (~200 LOC)
```
cmd/api-condition-report/main.go             — POST /bookings/{id}/condition-report
internal/adapter/handler/condition_report.go — CreateConditionReportHandler
internal/adapter/handler/condition_report_test.go
```
**Flow:** After photos uploaded via existing evidence-upload, customer/provider submits condition report with photo references -> validated -> saved -> publishes ConditionReportCreated event.

**Reuse:** Existing `evidence.ConditionReport` entity, `port.EvidenceRepo`, `evidence.ProcessEvidenceUseCase`.

#### Issue 5D: Provider Earnings Endpoint (~250 LOC)
```
internal/usecase/payment/earnings.go         — GetProviderEarningsUseCase
internal/usecase/payment/earnings_test.go
cmd/api-provider-earnings/main.go            — GET /providers/{id}/earnings
internal/adapter/handler/provider_earnings.go
internal/adapter/handler/provider_earnings_test.go
```
**Flow:** Provider requests earnings -> query payments by provider (existing GSI) -> aggregate totals (today, week, month, all-time) -> return summary.

**Reuse:** Existing `port.PaymentByBookingLister`, `dynamo_payment.go` repo (add FindByProvider query).

### Sprint 5 Execution
- All 4 issues are independent -> 4 parallel golang-pro agents
- Each creates a feature branch, implements, runs lint+test, pushes PR
- Merge order: 5C (smallest) -> 5B -> 5D -> 5A
- ~1,200 LOC total, 3-5 day turnaround with agents

---

## Part 2: Mobile App Architecture

### Tech Stack
| Layer | Choice | Why |
|-------|--------|-----|
| Framework | Expo SDK 55+ (React Native) | OTA updates, managed workflow, 85% Android PH market |
| Navigation | Expo Router v7 (file-based) | Next.js-like DX, deep links free |
| Language | TypeScript (strict) | Shared types with backend via OpenAPI codegen |
| State | Zustand + MMKV persist | Lightweight, WebSocket-compatible, offline-first |
| API Client | TanStack Query v5 | Auto-caching, retry, optimistic updates, offline queue |
| Auth | @aws-amplify/auth (standalone) | Cognito integration without full Amplify overhead |
| Maps | @rnmapbox/maps (Mapbox GL) | Offline tiles for PH provinces, better perf than Google Maps |
| Real-time | Built-in WebSocket | Connect to existing API Gateway WebSocket endpoints |
| Push | expo-notifications + FCM/APNs | Handled by backend SNS Platform Application |
| Camera | expo-camera + expo-image-picker | Evidence photos, condition reports |
| Storage | expo-secure-store | Tokens, sensitive data |
| Payments | react-native-webview | GCash/Maya payment pages (no native SDK needed for MVP) |
| Type Gen | @hey-api/openapi-ts | Go backend -> OpenAPI spec -> TypeScript client |

### Repo Decision: Same Monorepo
The mobile app lives in `apps/mobile/` alongside the Go backend. Rationale:
- Shared OpenAPI types between backend and frontend (single source of truth)
- Separate CI workflows (Go CI ignores `apps/`, Expo CI ignores `internal/`/`cmd/`)
- Atomic commits when backend API changes require mobile updates
- No cross-repo coordination overhead for a solo developer

### Auth: @aws-amplify/auth -> Existing Cognito Pool
The mobile app uses `@aws-amplify/auth` (standalone, ~50KB) as a **client** for the Cognito pool already deployed at `infra/modules/cognito/main.tf`. No new auth infrastructure needed:
- SRP auth + social login -> existing Cognito User Pool + mobile client
- Token refresh -> automatic (30-day refresh tokens configured)
- Pre-signup/post-confirmation/pre-token triggers -> already built as Go Lambdas
- API calls -> attach Cognito JWT -> existing API Gateway authorizer validates
- `ExtractUserID()` / `ExtractUserType()` in handler helpers work unchanged

### App Structure
```
app/                     <- Expo Router file-based routes
  (auth)/                <- Auth group (login, signup)
    login.tsx
    signup.tsx
  (tabs)/                <- Main tab navigator
    index.tsx            <- Home (request tow)
    history.tsx          <- Booking history
    profile.tsx          <- User profile
  booking/
    [id].tsx             <- Booking detail + tracking
    diagnose.tsx         <- AI diagnosis
    rate.tsx             <- Rating screen
  provider/
    dashboard.tsx        <- Provider home
    earnings.tsx         <- Earnings summary
  sos.tsx                <- SOS screen
  _layout.tsx            <- Root layout
components/              <- Shared UI components
  ui/                    <- Primitives (Button, Card, Input, Badge)
  booking/               <- Booking-specific (StatusBadge, PriceCard)
  map/                   <- Map components (TrackingMap, PickupMarker)
hooks/                   <- Custom hooks
  useAuth.ts             <- Cognito auth state
  useBooking.ts          <- Booking CRUD + state machine
  useWebSocket.ts        <- WS connection management
  useLocation.ts         <- GPS tracking
lib/                     <- Core utilities
  api/                   <- Generated API client (from OpenAPI)
  ws/                    <- WebSocket client + reconnect logic
  storage/               <- MMKV + Secure Store wrappers
  theme/                 <- Brand colors, typography (from UI mockup)
stores/                  <- Zustand stores
  auth.ts                <- User session, tokens
  booking.ts             <- Active booking state
  location.ts            <- Provider GPS state
  notifications.ts       <- Push notification state
assets/                  <- Fonts (Poppins), images, icons
```

### Brand Identity (from UI Mockup)
```
Colors:
  navy:   #0B1D33  (primary background)
  teal:   #00897B  (primary action)
  gold:   #F5A623  (accent, premium)
  orange: #FF6B35  (urgency, SOS)
  white:  #FFFFFF  (text on dark)

Typography: Poppins (Regular, SemiBold, Bold)
Border radius: 16px (cards), 12px (buttons), 24px (pills)
```

### Screen Inventory (22 screens from UI mockup)
| Screen | Route | Priority | Backend Dependency |
|--------|-------|----------|-------------------|
| Splash/Logo | (auth)/ | P0 | None |
| Login | (auth)/login | P0 | Cognito |
| Home (Request Tow) | (tabs)/index | P0 | Booking create |
| AI Diagnose | booking/diagnose | P1 | Issue 5A |
| Service Type Select | booking/service | P0 | Static data |
| Vehicle Select | booking/vehicle | P0 | User vehicles |
| Dropoff Location | booking/dropoff | P0 | Maps |
| Price Estimate | booking/price | P0 | Pricing engine |
| Matching | booking/matching | P0 | Matching + WS |
| Provider Matched | booking/matched | P0 | WS event |
| Live Tracking | booking/[id] | P0 | WS location |
| Chat | booking/chat | P1 | WS chat |
| Condition Report | booking/condition | P0 | Issue 5C + Evidence |
| Complete | booking/complete | P0 | Booking status |
| Rate | booking/rate | P0 | Rating submit |
| SOS | sos | P0 | SOS trigger |
| Provider Dashboard | provider/dashboard | P0 | Provider endpoints |
| Earnings | provider/earnings | P1 | Issue 5D |
| History | (tabs)/history | P1 | Booking list |
| Profile | (tabs)/profile | P1 | User CRUD |
| Typhoon Mode | booking/typhoon | P2 (defer) | Not implemented |
| Suki Program | profile/suki | P2 (defer) | Not implemented |

---

## Part 3: Mobile Sprint Plan (Weeks 2-6)

### Week 2: Foundation + Auth (Mobile Sprint 1)
**Goal:** App boots, user can log in, see home screen with map.

| Task | Files |
|------|-------|
| Expo project scaffold (SDK 55, Router v7, TypeScript) | app/, package.json, tsconfig |
| Theme system (colors, typography, spacing from brand guide) | lib/theme/ |
| Shared UI primitives (Button, Card, Input, Badge, StatusBar) | components/ui/ |
| Cognito auth flow (login, signup, token refresh) | hooks/useAuth.ts, stores/auth.ts |
| Auth screens (login, signup) | app/(auth)/ |
| Home screen with Mapbox map | app/(tabs)/index.tsx |
| Bottom tab navigator (Home, History, Profile) | app/(tabs)/_layout.tsx |
| API client setup (TanStack Query + generated types) | lib/api/ |

### Week 3: Core Booking Flow (Mobile Sprint 2)
**Goal:** Customer can request a tow from start to price estimate.

| Task | Files |
|------|-------|
| Service type selection screen | app/booking/service.tsx |
| Vehicle selection screen | app/booking/vehicle.tsx |
| Dropoff location picker (Mapbox geocoding) | app/booking/dropoff.tsx |
| Price estimate screen (calls pricing API) | app/booking/price.tsx |
| Booking create flow (orchestrates all screens) | hooks/useBooking.ts |
| WebSocket client with auto-reconnect | lib/ws/ |
| Matching screen (loading + WS events) | app/booking/matching.tsx |
| Provider matched screen | app/booking/matched.tsx |

### Week 4: Tracking + Safety (Mobile Sprint 3)
**Goal:** Live tracking works, SOS functional, provider can accept jobs.

| Task | Files |
|------|-------|
| Live tracking map (provider location via WS) | app/booking/[id].tsx, components/map/ |
| Provider dashboard (accept/reject jobs via WS) | app/provider/dashboard.tsx |
| Provider location broadcasting (background GPS) | hooks/useLocation.ts |
| SOS trigger screen (silent alarm) | app/sos.tsx |
| Push notification setup (expo-notifications) | hooks/usePushNotifications.ts |
| Booking status badge component | components/booking/StatusBadge.tsx |
| Complete booking screen | app/booking/complete.tsx |

### Week 5: Evidence + Rating + Polish (Mobile Sprint 4)
**Goal:** Full booking lifecycle complete, evidence capture works.

| Task | Files |
|------|-------|
| Condition report camera flow (8-photo) | app/booking/condition.tsx |
| Evidence upload (S3 presigned URLs) | hooks/useEvidence.ts |
| Rating screen (1-5 stars, tags, comment) | app/booking/rate.tsx |
| Chat screen (WS bidirectional) | app/booking/chat.tsx |
| History screen (paginated booking list) | app/(tabs)/history.tsx |
| Profile screen (edit profile, vehicles) | app/(tabs)/profile.tsx |
| Payment screen (GCash/Maya webview) | app/booking/payment.tsx |

### Week 6: Integration + Testing + Store Prep
**Goal:** App is testable end-to-end, ready for TestFlight/Internal Testing.

| Task | Files |
|------|-------|
| E2E test suite (Detox or Maestro) | e2e/ |
| Offline mode (queue actions, retry on reconnect) | lib/offline/ |
| Error boundaries + crash reporting | components/ErrorBoundary.tsx |
| App icon, splash screen, store metadata | assets/, app.json |
| EAS Build configuration (dev, preview, production) | eas.json |
| Performance optimization (lazy screens, image caching) | Various |
| Provider earnings screen | app/provider/earnings.tsx |
| AI diagnosis screen | app/booking/diagnose.tsx |

---

## Part 4: Type Sharing Pipeline

### Go -> OpenAPI -> TypeScript
```
1. Add swaggo/swag annotations to Go handlers (or hand-write openapi.yaml)
2. Generate OpenAPI 3.1 spec: `task generate-openapi`
3. Generate TypeScript client: `npx @hey-api/openapi-ts -i openapi.yaml -o packages/shared-types/`
4. Mobile app imports types: `import { Booking, Payment } from '@towcommand/shared-types'`
```

For Sprint 5 (Week 1), we hand-write a minimal `openapi.yaml` covering the endpoints the mobile app needs. Full swaggo annotation is a Week 6 polish task.

---

## Part 5: What Gets Deferred

| Item | Why Deferred | When |
|------|-------------|------|
| Admin dashboard (Next.js) | User decision — focus on mobile first | Post-MVP (Month 4) |
| Sprint 6 backend (dispute, analytics pipeline, load testing, security audit, staging/prod deploy) | Not needed for MVP internal testing | Post-MVP |
| Typhoon Mode screen | Backend not implemented (V1.1 feature) | Post-MVP |
| Suki Program screen | Backend loyalty system not scoped | Post-MVP |
| Real PayMongo integration | Mock gateway sufficient for testing | Pre-launch |
| Blockchain evidence anchoring | V1.1 feature | Post-MVP |
| ML risk scoring (replace rule-based) | V1.1 feature | Post-MVP |
| Geofence zone database (IsHighRiskZone) | Hardcoded false for MVP | Post-MVP |
| Route monitor -> WS integration | Implemented but not wired | Sprint 5 stretch |

---

## Part 6: Parallel Execution Strategy

```
WEEK 1 (Backend Sprint 5):
  +-- golang-pro Agent 1: Issue 5A (AI Diagnosis)
  +-- golang-pro Agent 2: Issue 5B (Push Notifications)
  +-- golang-pro Agent 3: Issue 5C (Condition Report Handler)
  +-- golang-pro Agent 4: Issue 5D (Provider Earnings)

WEEK 1 (Parallel — Mobile Setup):
  +-- expo-mobile-dev Agent: Scaffold Expo project, theme, primitives

WEEKS 2-5 (Mobile Development):
  +-- expo-mobile-dev Agent: One sprint per week (Sprints 1-4 above)

WEEK 6 (Integration + Polish):
  +-- expo-mobile-dev Agent: E2E tests, offline mode, store prep
  +-- golang-pro Agent: OpenAPI spec generation, any backend fixes found during integration
```

### Risk Mitigation
1. **WebSocket integration complexity** — Test WS connection early (Week 2). If Expo WS has issues, fall back to polling with TanStack Query refetchInterval.
2. **Mapbox licensing** — Free tier covers 25K monthly active users. Sufficient for MVP.
3. **Cognito token refresh** — @aws-amplify/auth handles this automatically, but test with expired tokens in Week 3.
4. **Background GPS (provider)** — Expo has `expo-location` with background task support on both platforms, but iOS background location requires specific permission justifications for App Store.
5. **Camera permissions on Android 13+** — expo-camera handles runtime permissions, but test on physical Android 13+ device.

---

## Verification Plan

### Backend (Sprint 5)
- `go build ./...` — compiles
- `golangci-lint run ./...` — lints clean
- `go test -race -count=1 ./...` — all tests pass
- E2E results posted on each PR before merge

### Mobile (Each Sprint)
- `npx expo lint` — no lint errors
- `npx tsc --noEmit` — TypeScript checks pass
- Manual testing on Expo Go (dev) -> EAS Build (preview)
- Screen-by-screen verification using Chrome browser automation (screenshot comparison to UI mockup)
- E2E tests with Maestro (Week 6)

### Integration
- Full booking flow: login -> request tow -> match -> track -> condition report -> complete -> rate -> pay
- SOS flow: trigger -> verify backend receives -> verify notification sent
- Provider flow: login -> go online -> receive job -> accept -> navigate -> complete
