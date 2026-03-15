# /integrate-mobile — Mobile ↔ Backend Integration

## Purpose

Wire the TowCommand PH mobile app screens to the real Go backend API. Replaces mock/hardcoded data with real API calls, WebSocket connections, and Cognito auth.

## Usage

```
/integrate-mobile [phase]
```

## Arguments

- `$ARGUMENTS`: Phase to execute — `auth`, `booking`, `provider`, `tracking`, `payment`, or `all`

---

Execute mobile-backend integration for: **$ARGUMENTS**

## Context

The API client (`lib/api/client.ts`) and WebSocket client (`lib/ws/client.ts`) are already fully implemented with real Cognito auth, auto-reconnect, and heartbeat. The 6 Zustand stores have the right shapes. The gap is: **screens use hardcoded mock data instead of calling the real clients.**

This is a **wiring problem**, not a building problem.

## Architecture Reference

### Go Backend API Surface (27 REST + 5 WebSocket)

**Booking:** POST/GET /bookings, GET /bookings/{id}, PATCH /bookings/{id}/status, POST /bookings/{id}/cancel
**Payment:** POST /bookings/{id}/payments, POST /payments/{id}/capture, POST /payments/{id}/refund, POST /bookings/{id}/cancel-fee, POST /payments/webhook
**Rating:** POST /bookings/{id}/rating, GET /bookings/{id}/rating
**Safety:** POST /bookings/{id}/sos, POST /sos/{id}/resolve, GET /admin/sos/active
**Evidence:** POST /bookings/{id}/evidence-upload, POST /bookings/{id}/condition-report
**Diagnosis:** POST /diagnosis
**OTP:** POST /bookings/{id}/otp/generate, POST /bookings/{id}/otp/verify
**Provider:** POST /providers, PATCH /providers/{id}/availability, PATCH /providers/{id}/location, GET /providers/nearby, GET /providers/{id}/earnings
**Push:** POST /users/{id}/push-token
**Admin:** GET /admin/stats/bookings

**WebSocket:** $connect, $disconnect, locationUpdate, chatMessage, bookingStatus

### Zustand Stores (6)
- `authStore` — user, isAuthenticated, isLoading (persisted)
- `bookingStore` — activeBooking, priceBreakdown, paymentMethod, matchingState, matchedProvider, otp (in-memory)
- `locationStore` — lat, lng, heading, accuracy, isTracking (in-memory)
- `vehicleStore` — vehicles[], selectedVehicleId, selectedCondition (persisted)
- `diagnosisStore` — step, selectedSymptoms[], diagnosisResult (in-memory)
- `notificationStore` — notifications[], pushToken, unreadCount (in-memory)

### Booking State Machine (13 states, linear)
```
PENDING → MATCHED → EN_ROUTE → ARRIVED → CONDITION_REPORT → OTP_VERIFIED → LOADING → IN_TRANSIT → ARRIVED_DROPOFF → OTP_DROPOFF → COMPLETED
Early cancel from: PENDING, MATCHED, EN_ROUTE → CANCELLED
```

### EventBridge Events (subscribed via WebSocket on mobile)
- `BookingMatched` — {bookingId, providerId, providerName, ETA}
- `BookingStatusChanged` — {bookingId, oldStatus, newStatus}
- `LocationUpdated` — {providerId, bookingId, lat, lng}
- `PaymentCaptured` — {paymentId, bookingId, amount}
- `SOSTriggered` — {alertId, bookingId, userId, triggerType}

---

## Execution Plan — 5 Flows, 3 Phases

### Dependency Graph
```
Phase 1: Auth ──── MUST BE FIRST ────→ All other flows need JWT
                                        │
Phase 2: Booking Creation  ←─ independent ─→  Provider
                │                                │
                ↓ (booking created)              ↓ (provider accepts)
Phase 3: Live Tracking ←── WebSocket sync ──────┘
                │
                ↓ (arrived at dropoff)
         Payment + Completion
```

---

### Phase 1: Auth Flow

**Screens:** `(auth)/login.tsx`, `(auth)/signup.tsx`
**Store:** `authStore`
**Dependency:** None — do this first

**Tasks:**
1. Replace mock `setUser()` in login screen with real `signIn()` from `@aws-amplify/auth`
2. Replace mock signup with real `signUp()` + `confirmSignUp()` (OTP verification)
3. On auth success: call `GET /users/me` (or decode JWT claims) to populate authStore
4. Add token refresh handling — Amplify handles this automatically, but verify it works
5. Add sign-out that calls `signOut()` AND `reset()` on ALL stores
6. Add auth guard in `(tabs)/_layout.tsx` — redirect to login if not authenticated
7. Handle auth errors: wrong password, user not confirmed, account locked
8. Write TanStack Query hook: `useCurrentUser()` that fetches user profile on app start

**Env vars needed:**
```
EXPO_PUBLIC_COGNITO_USER_POOL_ID=
EXPO_PUBLIC_COGNITO_CLIENT_ID=
EXPO_PUBLIC_API_URL=https://api.towcommand.ph
EXPO_PUBLIC_WS_URL=wss://ws.towcommand.ph
```

**E2E update:** Update `auth-login.spec.ts` and `auth-signup.spec.ts` to test error states with mock API responses.

**Definition of Done:**
- Real Cognito sign-in/sign-up works
- JWT stored in secure storage
- Auth guard redirects unauthenticated users
- Sign-out clears all stores
- Token refresh works silently

---

### Phase 2A: Booking Creation Flow

**Screens:** home → service → vehicle → diagnose → condition → dropoff → price → matching
**Stores:** bookingStore, vehicleStore, diagnosisStore
**Dependency:** Auth (Phase 1)

**Tasks:**
1. **Home screen** — "Tow" quick action navigates to `/booking/service` (already works)
2. **Service selection** — store selected service type in bookingStore
3. **Vehicle selection** — wire to vehicleStore (already persisted); add `POST /users/{id}/vehicles` if backend supports it, otherwise keep local-only
4. **Diagnose screen** — replace mock with real `POST /diagnosis` via TanStack Query mutation:
   ```tsx
   const diagnose = useMutation({
     mutationFn: (data) => api.post('/diagnosis', data),
     onSuccess: (result) => diagnosisStore.setResults(result),
   });
   ```
5. **Condition report** — wire camera to `POST /bookings/{id}/evidence-upload` for presigned URL, then upload directly to S3
6. **Dropoff screen** — integrate Mapbox GL for pickup/dropoff location selection with geocoding
7. **Price screen** — call `POST /bookings` with full booking payload:
   ```json
   {
     "vehicleId": "from vehicleStore",
     "serviceType": "from bookingStore",
     "locations": [{"lat": x, "lng": y, "type": "pickup"}, {"lat": x, "lng": y, "type": "dropoff"}],
     "notes": "optional"
   }
   ```
   Display returned price breakdown in the UI
8. **Matching screen** — CRITICAL: replace the 3-second `setTimeout` with a real WebSocket listener:
   ```tsx
   useEffect(() => {
     const ws = getWebSocketClient();
     ws.onMessage('bookingStatus', (data) => {
       if (data.status === 'MATCHED') {
         bookingStore.setMatchingState('found');
         bookingStore.setMatchedProvider(data.provider);
         bookingStore.setOtp(data.otp);
         router.replace('/booking/matched');
       }
     });
     // Timeout after 120 seconds
     const timeout = setTimeout(() => {
       bookingStore.setMatchingState('timeout');
       router.replace('/(tabs)');
     }, 120_000);
     return () => clearTimeout(timeout);
   }, []);
   ```

**2nd-Order Risks:**
- The matching screen timeout (120s) must show a countdown or progress indicator
- If WebSocket disconnects during matching, the user is stuck — need reconnection + re-subscribe
- `POST /bookings` might fail if the provider location cache (Redis) is cold — handle gracefully

**E2E update:** Matching test needs to handle the WebSocket timeout instead of 3s redirect.

---

### Phase 2B: Provider Flow

**Screens:** `provider/dashboard.tsx`, `provider/earnings.tsx`
**Store:** authStore (provider role)
**Dependency:** Auth (Phase 1)

**Tasks:**
1. **Dashboard** — wire stats to `GET /providers/{id}/earnings`:
   ```tsx
   const { data } = useQuery({
     queryKey: ['provider', 'earnings', providerId],
     queryFn: () => api.get(`/providers/${providerId}/earnings`),
   });
   ```
2. **Online/Offline toggle** — `PATCH /providers/{id}/availability` with optimistic update
3. **Location broadcasting** — start WebSocket `locationUpdate` when provider goes online:
   ```tsx
   useEffect(() => {
     if (!isOnline) return;
     const interval = setInterval(async () => {
       const location = await Location.getCurrentPositionAsync();
       wsClient.send('locationUpdate', { lat: location.coords.latitude, lng: location.coords.longitude });
     }, 10_000); // Every 10 seconds
     return () => clearInterval(interval);
   }, [isOnline]);
   ```
4. **Earnings screen** — wire period tabs (daily/weekly/monthly/all-time) to earnings endpoint with date filters

**2nd-Order Risks:**
- Background location drains battery — throttle to 10s intervals, stop when app is backgrounded
- Provider must be distinguished from customer via authStore.user.userType

---

### Phase 3A: Live Tracking Flow

**Screens:** `booking/[id].tsx` (tracking), `booking/chat.tsx`
**Stores:** bookingStore, locationStore
**Dependency:** Booking Creation (Phase 2A)

**Tasks:**
1. **Tracking screen** — subscribe to WebSocket events:
   - `locationUpdate` → update provider marker on Mapbox GL map
   - `bookingStatus` → update status badge, trigger screen transitions
   - `ETAUpdated` → update ETA display
2. **Mapbox GL integration** — render real map with:
   - User's pickup location (static pin)
   - Provider's real-time location (animated marker)
   - Route line between them (Mapbox Directions API)
3. **Chat screen** — wire to WebSocket `chatMessage`:
   ```tsx
   // Send
   const sendMessage = (text) => {
     wsClient.send('chatMessage', {
       bookingId: activeBooking.id,
       message: text,
       recipientId: activeBooking.provider.id,
       senderId: authStore.user.id,
     });
   };
   // Receive
   wsClient.onMessage('chatMessage', (data) => {
     addMessageToChat(data);
   });
   ```
4. **Status transitions** — when provider updates status via their app, customer sees real-time changes:
   ```
   EN_ROUTE → show "Driver is on the way" + ETA
   ARRIVED → show "Driver has arrived" + prompt condition report
   IN_TRANSIT → show "En route to dropoff" + live tracking
   ARRIVED_DROPOFF → show "Arrived at dropoff" + prompt OTP
   ```
5. **Connection indicator** — show "Reconnecting..." banner when WebSocket drops

**2nd-Order Risks:**
- Mapbox GL on Expo Web may not work for E2E tests — need web fallback for map component
- WebSocket reconnection during tracking shows stale provider location — clear old markers on reconnect
- Chat messages received while app is backgrounded need to trigger push notification

---

### Phase 3B: Payment + Completion Flow

**Screens:** price → (WebView) → `booking/complete.tsx`, `booking/rate.tsx`
**Store:** bookingStore
**Dependency:** Booking Creation (Phase 2A)

**Tasks:**
1. **Payment initiation** — on "Book Now" in price screen:
   ```tsx
   const payment = await api.post(`/bookings/${bookingId}/payments`, {
     method: selectedPaymentMethod,
   });
   if (payment.method === 'gcash' || payment.method === 'maya') {
     // Open WebView with checkout URL
     router.push({ pathname: '/booking/payment-webview', params: { url: payment.checkoutUrl } });
   }
   ```
2. **Payment WebView** — new screen `booking/payment-webview.tsx`:
   - Load GCash/Maya checkout URL
   - Detect completion via redirect URL pattern (e.g., `towcommand://payment-success`)
   - Handle WebView crashes gracefully
3. **OTP verification** — at pickup: `POST /bookings/{id}/otp/generate` → show OTP to customer → provider enters → `POST /bookings/{id}/otp/verify`
4. **Complete screen** — wire to real booking data:
   ```tsx
   const { data: booking } = useQuery({
     queryKey: ['booking', bookingId],
     queryFn: () => api.get(`/bookings/${bookingId}`),
   });
   ```
5. **Rating** — `POST /bookings/{id}/rating` with mutation:
   ```tsx
   const submitRating = useMutation({
     mutationFn: (data) => api.post(`/bookings/${bookingId}/rating`, data),
     onSuccess: () => router.replace('/(tabs)'),
   });
   ```

**2nd-Order Risks:**
- GCash WebView on MIUI/Realme may intercept and open in external browser — detect and handle
- Payment webhook may arrive before WebView redirect — poll payment status as backup
- OTP SMS delivery in PH can take 5-30 seconds — show countdown, allow resend

---

## Spawn Commands (copy-paste ready)

```bash
# Phase 1 — do first
/spawn "Wire auth flow: replace mock login/signup with real Cognito signIn/signUp from @aws-amplify/auth. Populate authStore on success. Add auth guard layout. Handle errors. Add sign-out that clears all stores. Screens: (auth)/login.tsx, (auth)/signup.tsx. Store: authStore. Test: update auth E2E specs."

# Phase 2 — after auth merges (parallel)
/spawn "Wire booking creation flow: service selection → vehicle → diagnose (POST /diagnosis) → condition (POST /evidence-upload + S3) → dropoff (Mapbox picker) → price (POST /bookings) → matching (replace setTimeout with WebSocket BookingMatched listener, 120s timeout). Screens: 8 booking screens. Stores: bookingStore, vehicleStore, diagnosisStore. Test: update booking E2E specs."

/spawn "Wire provider flow: dashboard stats from GET /providers/{id}/earnings, online/offline toggle via PATCH /providers/{id}/availability, location broadcasting via WebSocket locationUpdate every 10s when online. Screens: provider/dashboard.tsx, provider/earnings.tsx. Store: authStore. Test: update provider E2E specs."

# Phase 3 — after booking creation merges (parallel)
/spawn "Wire live tracking: WebSocket locationUpdate → Mapbox GL map with animated provider marker + route line. WebSocket bookingStatus → status badge transitions. Chat screen → WebSocket chatMessage send/receive. Add reconnection banner. Screens: booking/[id].tsx, booking/chat.tsx. Stores: bookingStore, locationStore."

/spawn "Wire payment + completion: POST /payments on Book Now, GCash/Maya WebView checkout, OTP generate/verify at pickup+dropoff, complete screen with real booking data, rating submission via POST /rating. Create new booking/payment-webview.tsx screen. Stores: bookingStore."
```

## E2E Evidence Requirement

After each phase, follow the `e2e-evidence-per-pr` skill:
1. Run `pnpm test:e2e` — all tests must pass
2. Capture screenshots via `takeEvidence()`
3. Upload to `e2e-evidence-pr{N}` branch
4. Post PR comment with embedded screenshot tables

## Ship

After each phase completes and E2E passes:
```bash
/ship "feat(mobile): wire {phase-name} to real backend API"
```
