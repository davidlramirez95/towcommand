import { BookingStatus, GeoLocation, PriceBreakdown, ServiceType } from './booking';
import { PaymentMethod, PaymentStatus } from './payment';

export type EventSource =
  | 'tc.booking'
  | 'tc.matching'
  | 'tc.tracking'
  | 'tc.payment'
  | 'tc.sos'
  | 'tc.auth'
  | 'tc.provider'
  | 'tc.evidence'
  | 'tc.notification';

export interface TowCommandEvent<T = unknown> {
  source: EventSource;
  detailType: string;
  detail: T;
  metadata: EventMetadata;
}

export interface EventMetadata {
  eventId: string;
  correlationId: string;
  timestamp: string;
  version: string;
  actor?: { userId: string; userType: string };
}

// Booking Events
export interface BookingCreatedEvent {
  bookingId: string;
  customerId: string;
  serviceType: ServiceType;
  pickupLocation: GeoLocation;
  dropoffLocation: GeoLocation;
  price: PriceBreakdown;
}

export interface BookingStatusChangedEvent {
  bookingId: string;
  previousStatus: BookingStatus;
  newStatus: BookingStatus;
  changedBy: string;
  metadata?: Record<string, unknown>;
}

export interface BookingCompletedEvent {
  bookingId: string;
  customerId: string;
  providerId: string;
  price: PriceBreakdown;
  duration: number;
  distanceKm: number;
}

// Matching Events
export interface ProviderMatchedEvent {
  bookingId: string;
  providerId: string;
  providerName: string;
  providerPhone: string;
  truckPlate: string;
  eta: number;
  score: number;
}

export interface MatchTimeoutEvent {
  bookingId: string;
  cascade: number;
  attemptedProviders: string[];
}

// Tracking Events
export interface LocationUpdatedEvent {
  providerId: string;
  bookingId?: string;
  lat: number;
  lng: number;
  heading: number;
  speed: number;
  timestamp: string;
}

// Payment Events
export interface PaymentCompletedEvent {
  paymentId: string;
  bookingId: string;
  amount: number;
  method: PaymentMethod;
  status: PaymentStatus;
  gatewayRef: string;
}

// SOS Events
export interface SOSActivatedEvent {
  alertId: string;
  bookingId?: string;
  triggeredBy: string;
  triggerType: 'manual' | 'auto_route_deviation' | 'auto_panic';
  location: GeoLocation;
  timestamp: string;
}

// Auth Events
export interface UserRegisteredEvent {
  userId: string;
  email: string;
  phone: string;
  authProvider: 'google' | 'facebook' | 'apple' | 'phone';
}

// Provider Events
export interface ProviderOnlineEvent {
  providerId: string;
  lat: number;
  lng: number;
  serviceArea: string;
}
