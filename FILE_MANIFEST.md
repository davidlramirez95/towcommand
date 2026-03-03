# TowCommand PH - Complete File Manifest

## Total Files Created: 66

### packages/cache (9 files)
```
packages/cache/
├── package.json
├── tsconfig.json
└── src/
    ├── index.ts (exports all modules)
    ├── client.ts (Redis singleton with lazy connection)
    ├── keys.ts (Cache key generators and TTL constants)
    └── patterns/
        ├── geo-cache.ts (Geographic indexing, geoadd/georadius)
        ├── session.ts (User claims and WebSocket connection caching)
        ├── rate-limiter.ts (Rate limiting and job locking)
        └── surge-pricing.ts (Surge multiplier caching)
```

**Key Features:**
- Redis client singleton pattern
- Geospatial indexing for provider locations (within 5km radius)
- Session caching with TTL (5 min user claims, 30 min job locks)
- Rate limiting per user (configurable max requests)
- Idempotent job locking mechanism
- Surge pricing multiplier clamped 1.0-1.5x

### packages/auth (7 files)
```
packages/auth/
├── package.json
├── tsconfig.json
└── src/
    ├── index.ts (exports all modules)
    ├── cognito-client.ts (Cognito client singleton)
    └── middleware/
        ├── jwt-verify.ts (API Gateway authorizer)
        ├── rbac.ts (Role-based access control)
        └── ban-check.ts (Ban/suspension checking)
```

**Key Features:**
- Cognito client singleton
- JWT token decoding and validation
- Role-based access control (5 roles: customer, provider, fleet_manager, ops_agent, admin)
- Owner-or-role authorization pattern
- Pre-authentication ban/suspension checks
- Custom claim mapping

### services/api-gateway (17 files)
```
services/api-gateway/
├── package.json
├── tsconfig.json
└── src/
    ├── middleware/
    │   ├── error-handler.ts (Unified error response formatter)
    │   ├── cors.ts (CORS headers and preflight)
    │   ├── request-logger.ts (Structured logging with correlation IDs)
    │   └── validation.ts (Zod schema validation wrapper)
    └── handlers/
        ├── booking/
        │   ├── create.ts (POST /bookings - create with event publishing)
        │   ├── cancel.ts (DELETE /bookings/:id)
        │   ├── get.ts (GET /bookings/:id)
        │   ├── list.ts (GET /bookings - user or provider perspective)
        │   └── update-status.ts (PATCH /bookings/:id/status)
        ├── provider/
        │   ├── register.ts (POST /providers - with document validation)
        │   ├── update-location.ts (POST /providers/location - geo cache update)
        │   ├── toggle-availability.ts (POST /providers/availability)
        │   └── get-nearby.ts (GET /providers/nearby?lat=X&lng=Y)
        ├── user/
        │   ├── profile.ts (GET/PATCH /users/profile)
        │   └── vehicles.ts (GET/POST /users/vehicles)
        ├── diagnosis/
        │   └── analyze.ts (POST /diagnosis/analyze - OBD parsing)
        ├── payment/
        │   ├── initiate.ts (POST /payments - payment intent creation)
        │   └── webhook.ts (POST /webhooks/payment - payment provider)
        └── rating/
            ├── submit.ts (POST /ratings - post-booking)
            └── get.ts (GET /ratings/:userId)
```

**Key Features:**
- Unified error handling with proper HTTP status codes
- CORS preflight and headers
- Correlation ID tracking for debugging
- Zod schema validation
- Event publishing for async processing
- Geo-spatial queries
- Role-based access control

### services/websocket (7 files)
```
services/websocket/
├── package.json
├── tsconfig.json
└── src/
    ├── handlers/
    │   ├── connect.ts (WebSocket $connect route)
    │   ├── disconnect.ts (WebSocket $disconnect route)
    │   ├── location-update.ts (Real-time provider location)
    │   ├── booking-status.ts (Booking status broadcasts)
    │   └── chat-message.ts (In-app messaging)
    └── lib/
        ├── connection-manager.ts (ApiGatewayManagementApi wrapper)
        └── broadcast.ts (User, provider, and booking-level broadcasts)
```

**Key Features:**
- Real-time location streaming
- Broadcast patterns (to user, to providers, to booking parties)
- Connection lifecycle management
- Chat message persistence
- Status update notifications
- ETA calculations

### services/matching (8 files)
```
services/matching/
├── package.json
├── tsconfig.json
└── src/
    ├── handler.ts (EventBridge subscriber for BookingCreated)
    ├── algorithms/
    │   ├── nearest.ts (Distance-based sorting)
    │   ├── weighted-score.ts (Multi-factor: distance, rating, acceptance)
    │   └── surge-aware.ts (Surge pricing adjustment)
    └── lib/
        ├── geo-search.ts (Redis geo-spatial queries)
        └── timeout.ts (Offer expiration management)
```

**Key Features:**
- Multiple matching algorithms
- Geospatial search within configurable radius
- Rating and acceptance rate factoring
- Surge pricing awareness
- Offer timeout management (30 seconds default)
- Provider availability filtering
- Push notification triggering

### services/notifications (9 files)
```
services/notifications/
├── package.json
├── tsconfig.json
└── src/
    ├── handler.ts (EventBridge event processor)
    ├── channels/
    │   ├── sms.ts (AWS SNS integration)
    │   ├── push.ts (AWS Pinpoint integration)
    │   └── email.ts (AWS SES integration)
    └── templates/
        ├── booking-confirmed.ts (SMS + Email templates)
        ├── driver-arriving.ts (ETA and driver info)
        ├── otp-code.ts (Verification code)
        └── sos-alert.ts (Emergency alerts)
```

**Key Features:**
- Multi-channel notifications (SMS, push, email)
- Template rendering for each channel
- EventBridge integration for event-driven notifications
- Support for OTP, booking confirmations, driver arrivals
- SOS alert routing
- Message localization support

### services/auth-triggers (5 files)
```
services/auth-triggers/
├── package.json
├── tsconfig.json
└── src/
    ├── pre-signup.ts (Auto-confirmation for trusted users)
    ├── post-confirmation.ts (Create user records in DB)
    ├── pre-token.ts (Add custom claims to JWT)
    ├── custom-message.ts (SMS/email message customization)
    └── pre-authentication.ts (Ban checks)
```

**Key Features:**
- User signup validation
- Auto-confirmation support
- Post-signup user record creation
- Custom JWT claims injection
- Message customization and localization
- Ban/suspension enforcement

### services/analytics (8 files)
```
services/analytics/
├── package.json
├── tsconfig.json
└── src/
    ├── handler.ts (EventBridge event aggregator)
    ├── lib/
    │   ├── pg-client.ts (PostgreSQL connection pool)
    │   └── schema.sql (Analytics table definitions)
    └── queries/
        ├── revenue-report.ts (Daily/periodic revenue analytics)
        ├── provider-performance.ts (Provider metrics and leaderboards)
        └── demand-heatmap.ts (Geographic demand visualization)
```

**Key Features:**
- Event-driven analytics aggregation
- PostgreSQL time-series storage
- Revenue reporting and breakdown
- Provider performance metrics
- Completion rates and ratings
- Geographic demand heatmaps
- Leaderboard generation

## Dependencies Summary

### packages/cache
- `ioredis@5.3.2` - Redis client
- `@towcommand/core@workspace:*`

### packages/auth
- `@aws-sdk/client-cognito-identity-provider@3.461.0`
- `jsonwebtoken@9.1.2`
- `jwks-rsa@3.0.1`
- `@towcommand/core@workspace:*`

### services/api-gateway
- `@towcommand/{core,db,events,cache,auth}@workspace:*`
- `@aws-sdk/client-s3@3.461.0`
- `zod@3.22.4`
- `ulid@2.3.0`
- `pino@8.17.2`

### services/websocket
- `@towcommand/{core,db,cache}@workspace:*`
- `@aws-sdk/client-apigatewaymanagementapi@3.461.0`
- `pino@8.17.2`

### services/matching
- `@towcommand/{core,db,cache,events}@workspace:*`
- `pino@8.17.2`

### services/notifications
- `@towcommand/{core,events}@workspace:*`
- `@aws-sdk/{client-sns,client-ses,client-pinpoint}@3.461.0`
- `pino@8.17.2`

### services/auth-triggers
- `@towcommand/{core,db}@workspace:*`
- `@aws-sdk/client-dynamodb@3.461.0`
- `pino@8.17.2`

### services/analytics
- `@towcommand/{core,events}@workspace:*`
- `pg@8.10.0`
- `pino@8.17.2`

## Key Patterns & Architecture

### Singleton Pattern
- Redis client (lazy connect)
- Cognito client
- Database connections

### Repository Pattern
- BookingRepository
- ProviderRepository
- UserRepository
- VehicleRepository
- RatingRepository
- ChatRepository

### Strategy Pattern
- Matching algorithms (nearest, weighted-score, surge-aware)

### Event-Driven Architecture
- EventBridge for async processing
- Booking, matching, notification, analytics services
- Event sourcing friendly

### Caching Strategy
- Geospatial: Provider locations (60s TTL)
- Session: User claims (300s TTL), WS connections (persistent)
- Rate limiting: 60s sliding window
- Job locks: 30s for idempotency
- Pricing: 30min for surge multipliers

### Error Handling
- AppError base class with status codes
- ValidationError for schema validation
- Zod integration
- Structured error responses

### Authorization
- JWT validation in authorizer
- Role-based access control (5 roles)
- Owner-or-role pattern
- Custom claims in tokens

## File Sizes & Complexity

Most files are focused and under 250 lines:
- Client/singleton files: ~20 lines
- Pattern implementations: ~50-80 lines
- Middleware: ~60-80 lines
- Handler stubs: ~30-50 lines (with TODO comments)
- Repository patterns: ~60-100 lines

## Next Implementation Steps

1. Implement business logic in all TODO handlers
2. Add integration tests for each service
3. Create SAM/CloudFormation templates
4. Set up CI/CD pipeline
5. Configure environment variables
6. Create API documentation
7. Set up monitoring/alerts
8. Performance testing and optimization

## Code Quality Considerations

- All files use TypeScript with strict mode
- Consistent error handling
- Environment variable configuration
- Structured logging
- Type-safe event contracts
- Async/await patterns
- Database transaction support ready
