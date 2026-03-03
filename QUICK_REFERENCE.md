# TowCommand PH - Quick Reference Guide

## Directory Tree
```
/sessions/awesome-gallant-planck/mnt/towcommand/
├── packages/
│   ├── cache/           # Redis wrapper with patterns
│   └── auth/            # Cognito helpers & middleware
├── services/
│   ├── api-gateway/     # REST API (15 handlers)
│   ├── websocket/       # Real-time (5 handlers)
│   ├── matching/        # Provider matching engine
│   ├── notifications/   # Multi-channel notifications
│   ├── auth-triggers/   # Cognito Lambda triggers
│   └── analytics/       # PostgreSQL analytics
└── docs/                # Documentation files
```

## Key Files to Implement First

1. **packages/cache/src/client.ts** - Already done (Redis singleton)
2. **packages/auth/src/cognito-client.ts** - Already done (Cognito singleton)
3. **services/api-gateway/src/handlers/booking/create.ts** - Has example, needs completion
4. **services/matching/src/handler.ts** - EventBridge entry point
5. **services/notifications/src/handler.ts** - Event processor

## Common Patterns

### Using Redis Cache
```typescript
import { GeoCache, SessionCache, RateLimiter } from '@towcommand/cache';

// Geo operations
const geo = new GeoCache();
await geo.updateProviderLocation('PROV-001', 14.5995, 121.0137);
const nearby = await geo.getNearbyProviders(14.5995, 121.0137, 5);

// Session caching
const session = new SessionCache();
await session.setUserClaims('sub-123', { role: 'customer' });

// Rate limiting
const limiter = new RateLimiter();
const result = await limiter.checkLimit('USER-001', 100);
```

### Using Auth
```typescript
import { requireRole, jwtAuthorizer } from '@towcommand/auth';

// In a handler
requireRole('customer', 'provider')(userType);

// As API Gateway authorizer
export const handler = jwtAuthorizer;
```

### Handler Template
```typescript
import type { APIGatewayProxyEvent, APIGatewayProxyResult } from 'aws-lambda';
import { handleError, successResponse } from '../../middleware/error-handler';

export async function handler(event: APIGatewayProxyEvent): Promise<APIGatewayProxyResult> {
  try {
    const userId = event.requestContext.authorizer?.userId as string;
    // TODO: Implement logic
    return successResponse({ message: 'Success' });
  } catch (error) {
    return handleError(error);
  }
}
```

## Environment Variables

### Development
```bash
REDIS_HOST=localhost
COGNITO_REGION=ap-southeast-1
DATABASE_HOST=localhost
AUTO_CONFIRM_ENABLED=true
LOG_LEVEL=debug
```

### Production
```bash
REDIS_HOST=elasticache-endpoint
REDIS_PASSWORD=from-secrets-manager
COGNITO_USER_POOL_ID=pool-id
DATABASE_PASSWORD=from-secrets-manager
```

## Testing a Handler Locally

```bash
cd services/api-gateway
npm run build

# Test with AWS SAM
sam local start-api

# Or invoke directly
sam local invoke CreateBooking -e event.json
```

## Package Build
```bash
cd packages/cache
npm install
npm run build
```

## Common Tasks

### Add a new handler
1. Create file in `services/api-gateway/src/handlers/[domain]/[action].ts`
2. Import `handleError`, `successResponse`
3. Implement async handler function
4. Add route to API Gateway template

### Add a new event type
1. Add schema to `@towcommand/events` if not exists
2. Use `publishEvent()` in handler
3. Create EventBridge rule for processor
4. Implement handler in relevant service

### Debug a handler
1. Check logs in CloudWatch
2. Use correlation ID to trace request
3. Check cache/database state
4. Verify environment variables

## Cache Keys Pattern
```
provider:location:{id}           # 60s TTL
provider:available:{city}        # Geospatial index
user:claims:{sub}                # 300s TTL
rate:user:{userId}               # 60s TTL
job:lock:{jobId}                 # 30s TTL
surge:{region}                   # 1800s TTL
```

## Common Errors

### "Cannot find module '@towcommand/cache'"
- Solution: Build packages first: `npm run build` in packages/cache

### "Redis connection timeout"
- Solution: Check REDIS_HOST and REDIS_PORT env vars
- Dev: Start Redis locally `docker run -p 6379:6379 redis`

### "Cognito token invalid"
- Solution: Check COGNITO_USER_POOL_ID is correct
- Verify JWT hasn't expired

### "Rate limit exceeded"
- Solution: Check RateLimiter.checkLimit() result
- Return 429 Too Many Requests

## File Organization Rules

1. **One responsibility per file**
2. **Handlers are thin** - logic in service/repository
3. **Types defined inline** or in types/ folder
4. **Tests alongside source** (test files separate)
5. **No circular imports**
6. **Exports from index.ts** for each package

## Deployment Checklist

- [ ] All environment variables configured
- [ ] All TODO items completed
- [ ] Tests passing (unit and integration)
- [ ] TypeScript compiles without errors
- [ ] CloudWatch alarms configured
- [ ] Database migrations run
- [ ] Cache TTLs tuned for workload
- [ ] Error handling tested
- [ ] Rate limits appropriate
- [ ] CORS origins whitelisted

## Performance Tuning

### Redis
- Increase MAX_CONNECTIONS for high load
- Monitor cache hit ratio
- Set appropriate TTLs (shorter = more requests)

### Database
- Add indexes on frequently queried columns
- Use connection pooling (configured in pg-client.ts)
- Monitor query execution time

### Lambda
- Set memory appropriately (256MB-3GB)
- Adjust timeout (30s default)
- Use Lambda provisioned concurrency for burst

### API Gateway
- Set throttling limits
- Enable caching for GET requests
- Use compression

## Useful Commands

```bash
# Build all packages
npm run build

# Type check
npm run typecheck

# List all handlers
find services -name "*.ts" -path "*/handlers/*" | sort

# Count lines of code
find . -name "*.ts" | xargs wc -l

# Check package sizes
npm ls --depth=0
```

## Documentation Files

- **STRUCTURE_CREATED.md** - Architecture overview
- **FILE_MANIFEST.md** - Detailed file breakdown
- **USAGE_EXAMPLES.md** - Code samples
- **COMPLETION_SUMMARY.md** - Implementation status
- **QUICK_REFERENCE.md** - This file

## Contact/Support

For implementation help:
1. Check USAGE_EXAMPLES.md for code samples
2. Review existing handlers for patterns
3. Check @towcommand/core for error classes
4. Look at test files for usage patterns

## Version Info

- TypeScript: 5.3.3
- Node.js: 18+ (Lambda)
- ioredis: 5.3.2
- Zod: 3.22.4
- AWS SDK: 3.461.0
