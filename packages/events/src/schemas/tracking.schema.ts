// Tracking event detail schemas for EventBridge
// Zod schemas can be added here as needed for event validation

export interface LocationUpdatedDetail {
  bookingId: string;
  providerId: string;
  location: { lat: number; lng: number };
  timestamp: string;
}

export interface DriverArrivedDetail {
  bookingId: string;
  providerId: string;
  arrivedAt: string;
}

export interface RouteDeviationDetail {
  bookingId: string;
  providerId: string;
  deviationReason: string;
  detectedAt: string;
}
