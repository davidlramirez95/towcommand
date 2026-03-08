# TowCommand PH — Frontend Tech Stack Reference

## Overview

TowCommand is a Philippine tow truck / roadside assistance platform. The mobile app is the primary user interface — customers request tows, providers accept jobs, both track in real-time. The backend (Go + AWS serverless) is already built. This document explains every technology chosen for the mobile frontend, why it was chosen, and what it does.

---

## Core Framework

### Expo SDK 55+ (React Native)
**What it is:** Expo is a platform built on top of React Native that simplifies mobile app development. React Native (created by Meta) lets you write mobile apps in JavaScript/TypeScript that compile to native iOS and Android components — not web views, actual native UI.

**Why we chose it:**
- **One codebase → two platforms.** Write code once in TypeScript, Expo produces both an Android .apk/.aab and an iOS .ipa. 95%+ code is shared between platforms.
- **Managed workflow.** Expo handles native build tooling (Xcode, Android Studio, Gradle) so you don't need to configure them manually. `eas build` compiles everything in the cloud.
- **OTA updates.** Push JavaScript-only updates instantly to users without going through App Store / Google Play review. Critical for bug fixes.
- **PH market fit.** 85% of Philippine mobile users are on Android. Expo's Android support is first-class, and the base app size is ~15MB (fits the <50MB target for 3G networks).
- **Mature ecosystem.** Used by Shopify, Discord, Microsoft Teams. Not experimental.

**What it replaces:** Writing two separate apps in Swift (iOS) and Kotlin (Android), which would double development time.

---

### Expo Router v7 (File-Based Navigation)
**What it is:** The navigation/routing library for the app. It uses a file-system based approach — each file in the `app/` directory becomes a screen/route, similar to how Next.js handles web pages.

**Why we chose it:**
- **Automatic deep linking.** Every screen gets a URL automatically (e.g., `towcommand://booking/bk-123`). Push notifications can link directly to a booking screen.
- **Type-safe routes.** TypeScript knows which routes exist and what parameters they accept — catches navigation bugs at compile time.
- **Layouts and groups.** Supports nested layouts like `(auth)/` for login screens and `(tabs)/` for the main tab navigator — cleaner code organization.
- **SEO-ready.** If we ever add a web version (Expo supports web targets), routes are already defined.

**What it replaces:** React Navigation (manual route configuration with large boilerplate objects). Expo Router is built on React Navigation under the hood but with better DX.

---

## Language

### TypeScript (Strict Mode)
**What it is:** A typed superset of JavaScript. Every variable, function parameter, and return value has a declared type. The compiler catches errors before the code runs.

**Why we chose it:**
- **Catch bugs early.** If the backend returns a `Booking` object with a `status` field, TypeScript ensures the mobile app handles all possible status values — missing a case is a compile error, not a runtime crash.
- **Shared types with backend.** We generate TypeScript types from the Go backend's API spec (OpenAPI). The mobile app and backend share the same data contract.
- **AI agent productivity.** Claude Code agents write more accurate code when they have type information — fewer bugs, fewer iterations.
- **Industry standard.** React Native + TypeScript is the default for new projects.

**What it replaces:** Plain JavaScript, which allows any variable to be any type (leading to "undefined is not a function" runtime crashes).

---

## State Management

### Zustand
**What it is:** A minimal state management library for React. "State" means the data your app keeps in memory — logged-in user, active booking, provider location, etc.

**Why we chose it:**
- **Tiny (1.2KB).** Doesn't bloat the app bundle.
- **Simple API.** Create a store in ~10 lines. No boilerplate (unlike Redux which requires actions, reducers, selectors, middleware).
- **Works with WebSockets.** When a WebSocket message arrives (e.g., provider location update), you can update the store directly. Components subscribed to that data re-render automatically.
- **Persist to disk.** Combined with MMKV (see below), stores survive app restarts — critical for offline mode.

**Example:**
```typescript
// stores/booking.ts
const useBookingStore = create((set) => ({
  activeBooking: null,
  setActiveBooking: (booking) => set({ activeBooking: booking }),
}));

// Any component can read/update:
const booking = useBookingStore((s) => s.activeBooking);
```

**What it replaces:** Redux (complex, verbose), MobX (magic decorators), or React Context (causes unnecessary re-renders at scale).

### MMKV (via react-native-mmkv)
**What it is:** A fast key-value storage engine originally built by WeChat for 1 billion+ users. Used to persist Zustand stores to disk.

**Why we chose it:**
- **30x faster than AsyncStorage** (the default React Native storage). Synchronous reads — no `await` needed.
- **Offline support.** When the user loses connectivity (common in PH provinces), the app still works because data is stored locally.
- **Encryption support.** Can encrypt stored data at rest.

**What it replaces:** AsyncStorage (slow, async-only, no encryption).

---

## API & Data Fetching

### TanStack Query v5 (React Query)
**What it is:** A server-state management library. It handles fetching data from the backend, caching responses, retrying failed requests, and keeping the UI in sync with the server.

**Why we chose it:**
- **Automatic caching.** Fetch a booking once, and every screen that needs it reads from cache — no duplicate API calls.
- **Optimistic updates.** When a user submits a rating, show success immediately while the API call happens in the background. If it fails, roll back automatically.
- **Retry logic.** If an API call fails (network glitch), TanStack Query retries 3 times with exponential backoff — no manual retry code.
- **Offline queue.** When offline, mutations (create booking, submit rating) are queued and executed when connectivity returns.
- **Refetch on focus.** When the user switches back to the app, stale data is automatically refreshed.

**Example:**
```typescript
// Fetch booking — cached, retried, auto-refreshed
const { data: booking, isLoading } = useQuery({
  queryKey: ['booking', bookingId],
  queryFn: () => api.getBooking(bookingId),
});
```

**What it replaces:** Manual `fetch()` calls with `useState` + `useEffect` + loading/error state management (error-prone, no caching, no retry).

---

## Authentication

### @aws-amplify/auth (Standalone)
**What it is:** The authentication module from AWS Amplify, extracted as a standalone package. Handles Cognito user pools — signup, login, token refresh, social login (Google, Facebook, Apple).

**Why we chose it:**
- **Direct Cognito integration.** Our backend uses AWS Cognito for auth (5 Lambda triggers already implemented). This library speaks Cognito's protocol natively.
- **Standalone import.** We import ONLY the auth module (~50KB), not the full Amplify framework (~500KB+). Keeps the bundle small.
- **Token management.** Automatically refreshes JWT tokens before they expire. The mobile app never sees expired tokens.
- **Social login.** Built-in support for Google/Facebook/Apple Sign-In — important for PH market where social login is expected.

**What it replaces:** Building a custom OAuth2/OIDC flow from scratch (complex, security-sensitive) or importing the full Amplify library (bloated).

---

## Maps & Location

### @rnmapbox/maps (Mapbox GL Native)
**What it is:** A React Native wrapper around Mapbox's native map SDKs. Provides interactive maps with vector tiles, real-time markers, route drawing, and geocoding.

**Why we chose it:**
- **Offline map tiles.** Users in PH provinces with spotty connectivity can download map regions. Google Maps requires constant internet.
- **Better performance.** Vector tiles (Mapbox) render faster than raster tiles (Google Maps) — smoother panning/zooming on budget Android devices common in PH.
- **Custom styling.** We can match the map to TowCommand's brand (navy/teal color scheme).
- **Free tier.** 25,000 monthly active users free — sufficient well beyond MVP.
- **Geocoding + Directions.** Built-in APIs for address search and route calculation.

**What it replaces:** Google Maps (requires constant internet, $200/month beyond free tier for our usage pattern, raster tiles slower on low-end devices).

### expo-location
**What it is:** Expo's location module. Provides GPS coordinates, background location tracking, geofencing.

**Why we chose it:**
- **Background tracking.** Providers need to broadcast their location even when the app is backgrounded. expo-location supports this on both iOS and Android.
- **Permission handling.** Automatically shows the correct permission dialogs per platform (iOS is stricter than Android).
- **Battery optimization.** Supports "significant change" mode that uses less battery than continuous GPS polling.

---

## Real-Time Communication

### WebSocket (Built-in)
**What it is:** React Native includes a WebSocket client natively. We connect to the existing API Gateway WebSocket endpoints that the Go backend already serves.

**Why we chose it:**
- **Already built.** The backend has 5 WebSocket handlers (connect, disconnect, location-update, chat-message, booking-status). We just need to connect to them.
- **Bidirectional.** Provider sends location → server broadcasts to customer. Customer sends chat → server routes to provider. All over one persistent connection.
- **Low latency.** Real-time tracking needs sub-second updates. WebSocket delivers this. REST polling would add 1-3 second delays.

**What we build on top:**
- Auto-reconnect with exponential backoff (connection drops are common on mobile)
- Message type routing (location updates → Zustand location store, chat → chat store, status → booking store)
- Heartbeat ping/pong to detect dead connections

**What it replaces:** Polling (wasteful, laggy), Server-Sent Events (unidirectional — we need bidirectional), Firebase Realtime Database (vendor lock-in, not compatible with our DynamoDB backend).

---

## Push Notifications

### expo-notifications + FCM/APNs
**What it is:** Expo's notification module that interfaces with Firebase Cloud Messaging (Android) and Apple Push Notification Service (iOS).

**Why we chose it:**
- **Cross-platform.** One API to send push notifications to both Android and iOS.
- **Background delivery.** Notifications arrive even when the app is closed — "Your tow truck is 2 minutes away."
- **Deep link support.** Tapping a notification opens the relevant screen (e.g., booking tracking).
- **Backend integration.** Sprint 5B adds push token registration to our Go backend. The notification router (already built) will dispatch via SNS → FCM/APNs.

---

## Camera & Media

### expo-camera + expo-image-picker
**What it is:** Expo modules for camera access and photo selection from the device gallery.

**Why we chose it:**
- **Evidence capture.** The condition report workflow requires 8 photos (front, back, left, right, damage close-ups). expo-camera provides a custom camera UI.
- **Image picker.** Users can select existing photos from their gallery instead of taking new ones.
- **Permission handling.** Automatic runtime permission requests per platform.
- **Image compression.** Built-in quality/size options to keep uploads under bandwidth limits (important for 3G PH networks).

---

## Secure Storage

### expo-secure-store
**What it is:** Encrypted key-value storage using the platform's native secure storage (iOS Keychain, Android Keystore).

**Why we chose it:**
- **Token storage.** Cognito JWT tokens must be stored securely — not in plain text AsyncStorage.
- **Hardware-backed encryption.** On devices with secure enclaves (most modern phones), keys never leave the hardware.
- **Small data only.** Designed for secrets (tokens, API keys), not bulk data (use MMKV for that).

---

## Payments

### react-native-webview
**What it is:** A WebView component that renders web pages inside the app.

**Why we chose it for payments:**
- **GCash/Maya integration.** Both Philippine e-wallets provide web-based payment pages. The user taps "Pay with GCash" → WebView opens the GCash payment page → user authenticates in GCash → webhook confirms payment to our backend.
- **No native SDK needed.** GCash and Maya don't offer React Native SDKs. WebView is the standard integration pattern for mobile.
- **PCI compliance.** We never handle card numbers — the payment provider's web page does. Reduces our security scope.
- **Replaceable.** When PayMongo (our planned payment gateway) ships their React Native SDK, we swap the WebView for their native component. Zero backend changes needed.

---

## Type Generation (Backend → Frontend)

### @hey-api/openapi-ts
**What it is:** A code generator that reads an OpenAPI specification (describing our REST API endpoints, request/response shapes) and produces TypeScript types and an API client.

**Why we chose it:**
- **Single source of truth.** The Go backend defines the API. We generate an OpenAPI spec from it. This tool generates TypeScript types from that spec. No manual type duplication.
- **Type-safe API calls.** The generated client knows that `POST /bookings` expects a `CreateBookingRequest` body and returns a `Booking` response. TypeScript enforces this.
- **Auto-updated.** When the backend adds a new field, regenerate types → TypeScript immediately flags every place in the mobile app that needs updating.

**Pipeline:**
```
Go handlers → OpenAPI 3.1 spec (YAML) → @hey-api/openapi-ts → TypeScript types + client
```

---

## Development & Build Tools

### EAS (Expo Application Services)
**What it is:** Expo's cloud build and deployment service.

- **EAS Build:** Compiles the app in the cloud (no local Xcode/Android Studio needed). Produces .apk (Android) and .ipa (iOS).
- **EAS Update:** Pushes OTA (over-the-air) JavaScript updates to users without app store review.
- **EAS Submit:** Uploads builds to Google Play / App Store.

**Free tier:** 30 builds/month — sufficient for development. Production plan is $99/month when needed.

### pnpm + Turborepo (If monorepo)
**What it is:** Package manager (pnpm) and build system (Turborepo) for managing multiple packages in one repository.

- **pnpm:** Faster than npm, uses hard links to save disk space, strict dependency resolution prevents phantom dependencies.
- **Turborepo:** Caches build outputs across packages. If `shared-types` hasn't changed, it skips rebuilding it.

Only needed if the mobile app lives in the same repo as the Go backend (monorepo approach).

---

## Summary Table

| Category | Technology | Bundle Impact | Purpose |
|----------|-----------|---------------|---------|
| Framework | Expo SDK 55 | Base ~15MB | Cross-platform mobile runtime |
| Navigation | Expo Router v7 | Included | File-based routing + deep links |
| Language | TypeScript | 0 (compile-time) | Type safety |
| State | Zustand | 1.2KB | In-memory state management |
| Persistence | MMKV | 50KB native | Fast disk storage, offline support |
| API Client | TanStack Query v5 | 12KB | Caching, retry, offline queue |
| Auth | @aws-amplify/auth | ~50KB | Cognito login/signup/tokens |
| Maps | @rnmapbox/maps | ~5MB native | Interactive maps, offline tiles |
| Location | expo-location | Included | GPS tracking, background mode |
| Real-time | WebSocket (built-in) | 0 | Live tracking, chat, status updates |
| Push | expo-notifications | Included | FCM/APNs push notifications |
| Camera | expo-camera | Included | Evidence photo capture |
| Gallery | expo-image-picker | Included | Photo selection |
| Secure Storage | expo-secure-store | Included | Encrypted token storage |
| Payments | react-native-webview | 200KB | GCash/Maya payment pages |
| Type Gen | @hey-api/openapi-ts | 0 (dev-only) | Backend → TypeScript types |
| Build | EAS Build | N/A (cloud) | Cloud builds for iOS + Android |

**Estimated total app size:** ~25-30MB (well under 50MB target for PH 3G networks)

---

## Cost Summary (MVP Phase)

| Service | Cost | When Needed |
|---------|------|-------------|
| Google Play Developer | $25 one-time | Before first Android release |
| Apple Developer Program | $99/year | Before first iOS release (can defer) |
| Mapbox | Free (25K MAU) | From day one |
| EAS Build | Free (30 builds/month) | From day one |
| EAS Update | Free (1K users) | From day one |
| Expo Push | Free | From day one |
| **Total to launch on Android** | **$25** | |
| **Total to launch on both** | **$124** | |
