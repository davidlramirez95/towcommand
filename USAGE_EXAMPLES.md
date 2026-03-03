# TowCommand PH - Usage Examples

## Cache Package Usage

### GeoCache - Provider Location Management
```typescript
import { GeoCache } from '@towcommand/cache';

const geoCache = new GeoCache();

// Update provider location
await geoCache.updateProviderLocation('PROV-001', 14.5995, 121.0137, 'NCR');

// Find nearby providers within 5km
const nearbyProviders = await geoCache.getNearbyProviders(
  14.5995,
  121.0137,
  5,
  'NCR',
  20
);
// Returns: [{ providerId: 'PROV-001', distance: 0.5 }, ...]

// Get provider location
const location = await geoCache.getProviderLocation('PROV-001');
// Returns: { lat: 14.5995, lng: 121.0137, updatedAt: 1234567890 }
```

### SessionCache - User Session Management
```typescript
import { SessionCache } from '@towcommand/cache';

const sessionCache = new SessionCache();

// Cache user claims from Cognito token
await sessionCache.setUserClaims('sub-12345', {
  sub: 'sub-12345',
  email: 'user@example.com',
  'custom:user_type': 'customer',
  'custom:trust_tier': 'gold',
});

// Retrieve cached claims (5 min TTL)
const claims = await sessionCache.getUserClaims('sub-12345');
```

### RateLimiter - Request Rate Limiting
```typescript
import { RateLimiter } from '@towcommand/cache';

const limiter = new RateLimiter();

// Check rate limit (100 requests per minute)
const result = await limiter.checkLimit('USER-001', 100);
// Returns: { allowed: true, remaining: 99, resetIn: 45 }

// Job locking for idempotency
const acquired = await limiter.acquireJobLock('BOOKING-123');
if (acquired) {
  // Process booking
  await limiter.releaseJobLock('BOOKING-123');
}
```

## Auth Package Usage

### JWT Verification in API Gateway
```typescript
import { jwtAuthorizer } from '@towcommand/auth';

export const handler = jwtAuthorizer;
// Automatically validates JWT and populates authorizer context
```

### Role-Based Access Control
```typescript
import { requireRole, type Role } from '@towcommand/auth';

export async function handler(event: APIGatewayProxyEvent) {
  const userType = event.requestContext.authorizer?.userType as string;
  
  try {
    // Require customer or admin role
    requireRole('customer', 'admin')(userType);
    
    // Proceed with handler logic
    return successResponse({ success: true });
  } catch (error) {
    return handleError(error); // Returns 403 Forbidden
  }
}
```

## API Gateway Service Usage

### Create Booking Handler
```typescript
// POST /bookings
import { handler as createBooking } from './handlers/booking/create';

const event = {
  body: JSON.stringify({
    vehicleId: 'VEH-001',
    serviceType: 'towing',
    pickupLocation: { lat: 14.5995, lng: 121.0137 },
    dropoffLocation: { lat: 14.6091, lng: 121.0223 },
    estimateId: 'EST-001',
  }),
  requestContext: {
    authorizer: { userId: 'USER-001' }
  }
};

const response = await createBooking(event);
// Returns: { statusCode: 201, body: '{"success":true,"data":{...}}' }
```

### Update Location Handler
```typescript
// POST /providers/location
import { handler as updateLocation } from './handlers/provider/update-location';

const event = {
  body: JSON.stringify({
    lat: 14.5995,
    lng: 121.0137,
    city: 'NCR',
    accuracy: 10,
  }),
  requestContext: {
    authorizer: { providerId: 'PROV-001' }
  }
};

const response = await updateLocation(event);
// Updates Redis geo cache and triggers matching if needed
```

## WebSocket Service Usage

### Location Streaming
```typescript
// WebSocket route: location-update
// Sent by provider app

{
  "type": "location-update",
  "bookingId": "TC-2025-01HJPZED72PFK5D",
  "lat": 14.5995,
  "lng": 121.0137,
  "accuracy": 5,
  "heading": 45
}

// Broadcasting to customer
{
  "type": "location-update",
  "providerId": "PROV-001",
  "lat": 14.5995,
  "lng": 121.0137,
  "eta": 300  // seconds
}
```

## Matching Service

### EventBridge Event Processing
```typescript
// BookingCreated event from api-gateway
{
  "source": "booking",
  "detail-type": "BookingCreated",
  "detail": {
    "bookingId": "TC-2025-01HJPZED72PFK5D",
    "customerId": "USER-001",
    "pickupLocation": { "lat": 14.5995, "lng": 121.0137 },
    "serviceType": "towing",
    "estimateId": "EST-001"
  }
}

// Matching service:
// 1. Finds nearby providers using GeoCache
// 2. Scores using WeightedScoreAlgorithm (distance, rating, acceptance rate)
// 3. Adjusts for surge pricing using SurgePricingCache
// 4. Sends push notifications with 30-second expiration timeout
// 5. Publishes MatchingStarted event
```

## Notifications Service

### Template Usage
```typescript
import { renderBookingConfirmedSMS, renderBookingConfirmedEmail } from 
  './templates/booking-confirmed';

const sms = renderBookingConfirmedSMS({
  bookingId: 'TC-2025-01HJPZED72PFK5D',
  customerName: 'Juan',
  providerName: 'Maria',
  pickupAddress: '123 Bonifacio Ave, Makati',
  estimatedArrival: '8 minutes',
  totalCost: 450,
  providerPhone: '+63912345678',
});
// "Hi Juan, your booking TC-2025-01HJPZED72PFK5D is confirmed..."

const email = renderBookingConfirmedEmail({...});
// HTML email with booking details, driver info, cost breakdown
```

## Auth Triggers

### Pre-Token Generation
```typescript
// Event: TokenGenerationTriggerEvent
{
  "request": {
    "userAttributes": {
      "sub": "12345678-1234-1234-1234-123456789012",
      "email": "user@example.com",
      "custom:user_type": "provider",
      "custom:provider_id": "PROV-001"
    }
  }
}

// Response adds custom claims
{
  "claimsOverrideDetails": {
    "claimsToAddOrOverride": {
      "custom:user_type": "provider",
      "custom:trust_tier": "premium",
      "custom:provider_id": "PROV-001"
    }
  }
}
```

## Analytics Service

### Analytics Queries
```typescript
import { getRevenueReport } from './queries/revenue-report';
import { getProviderPerformance, getTopProviders } from 
  './queries/provider-performance';
import { getDemandHeatmap } from './queries/demand-heatmap';

// Daily revenue report
const report = await getRevenueReport('2025-01-01', '2025-01-31');
// [{ date: '2025-01-01', totalRevenue: 15000, bookingCount: 30, ... }]

// Provider stats
const stats = await getProviderPerformance('PROV-001');
// { totalBookings: 150, completionRate: 98, averageRating: 4.8, ... }

// Top performers
const topProviders = await getTopProviders(10);
// Ranked by rating, then booking count

// Demand heatmap for visualization
const heatmap = await getDemandHeatmap('2025-01-15');
// [{ gridCellId: 'h3-123', demandLevel: 4, bookingCount: 45, ... }]
```

## Environment Configuration

### Development (.env.local)
```bash
# Redis
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=

# Cognito
COGNITO_REGION=ap-southeast-1
COGNITO_USER_POOL_ID=ap-southeast-1_xxxxx
COGNITO_CLIENT_ID=xxxxx

# Database
DATABASE_HOST=localhost
DATABASE_PORT=5432
DATABASE_NAME=towcommand_dev
DATABASE_USER=postgres
DATABASE_PASSWORD=postgres

# AWS
AWS_REGION=ap-southeast-1
AUTO_CONFIRM_ENABLED=true
LOG_LEVEL=debug
```

### Production (Lambda Environment Variables)
```bash
REDIS_HOST=elasticache-cluster.xxxxx.ng.0001.apse1.cache.amazonaws.com
REDIS_PORT=6379
REDIS_PASSWORD=(from Secrets Manager)

COGNITO_USER_POOL_ID=ap-southeast-1_PROD123
DATABASE_HOST=towcommand-prod.xxxxx.rds.amazonaws.com
DATABASE_PASSWORD=(from Secrets Manager)

WEBSOCKET_ENDPOINT=https://xxxxx.execute-api.ap-southeast-1.amazonaws.com
PINPOINT_PROJECT_ID=xxxxx
```

## Error Handling Examples

### AppError Usage
```typescript
import { AppError } from '@towcommand/core';

// Throw validation error
throw AppError.badRequest('Invalid booking details');

// Throw not found
throw AppError.notFound('Booking not found');

// Throw forbidden
throw AppError.forbidden('You cannot cancel this booking');

// Throw conflict
throw AppError.conflict('Booking already completed');

// Returns automatic error response:
// { statusCode: 400, body: '{"error":"BAD_REQUEST","message":"..."}' }
```

### Validation Example
```typescript
import { z } from 'zod';
import { validateRequest } from '@towcommand/api-gateway';

const createBookingSchema = z.object({
  vehicleId: z.string().min(1),
  serviceType: z.enum(['towing', 'roadside-assistance', 'recovery']),
  pickupLocation: z.object({
    lat: z.number().min(-90).max(90),
    lng: z.number().min(-180).max(180),
  }),
});

try {
  const data = await validateRequest(
    JSON.parse(event.body),
    createBookingSchema
  );
  // Proceed with validated data
} catch (error) {
  return handleError(error); // Returns 400 with validation errors
}
```
