/**
 * Shared types between Go backend and mobile app.
 * These will eventually be generated from OpenAPI spec.
 * For now, hand-written to match Go domain entities.
 */

// === Booking ===

export type BookingStatus =
  | 'PENDING'
  | 'MATCHING'
  | 'MATCHED'
  | 'EN_ROUTE'
  | 'ARRIVED'
  | 'CONDITION_REPORT'
  | 'OTP_VERIFIED'
  | 'LOADING'
  | 'IN_TRANSIT'
  | 'ARRIVED_DROPOFF'
  | 'OTP_DROPOFF'
  | 'COMPLETED'
  | 'CANCELLED';

export type ServiceType =
  | 'FLATBED_TOWING'
  | 'WHEEL_LIFT_TOWING'
  | 'MOTORCYCLE_TOWING'
  | 'JUMPSTART'
  | 'TIRE_CHANGE'
  | 'LOCKOUT'
  | 'FUEL_DELIVERY'
  | 'WINCH_RECOVERY';

export interface Booking {
  bookingId: string;
  customerId: string;
  providerId?: string;
  serviceType: ServiceType;
  status: BookingStatus;
  pickupLat: number;
  pickupLng: number;
  pickupAddress: string;
  dropoffLat?: number;
  dropoffLng?: number;
  dropoffAddress?: string;
  estimatedCost: number; // centavos
  finalCost?: number; // centavos
  vehiclePlate: string;
  vehicleType: string;
  notes?: string;
  createdAt: string;
  updatedAt: string;
}

export interface CreateBookingRequest {
  serviceType: ServiceType;
  pickupLat: number;
  pickupLng: number;
  pickupAddress: string;
  dropoffLat?: number;
  dropoffLng?: number;
  dropoffAddress?: string;
  vehiclePlate: string;
  vehicleType: string;
  notes?: string;
}

export interface PriceEstimate {
  baseFare: number; // centavos
  distanceFare: number;
  surgeMultiplier: number;
  totalEstimate: number;
  currency: 'PHP';
}

// === Payment ===

export type PaymentMethod = 'gcash' | 'maya' | 'card' | 'cash' | 'corporate';
export type PaymentStatus = 'pending' | 'held' | 'captured' | 'refunded' | 'failed' | 'cancelled';

export interface Payment {
  paymentId: string;
  bookingId: string;
  userId: string;
  amount: number; // centavos
  currency: string;
  method: PaymentMethod;
  status: PaymentStatus;
  createdAt: string;
}

// === Rating ===

export interface Rating {
  ratingId: string;
  bookingId: string;
  customerId: string;
  providerId: string;
  score: number; // 1-5
  comment?: string;
  tags?: string[];
  createdAt: string;
}

export interface SubmitRatingRequest {
  score: number;
  comment?: string;
  tags?: string[];
}

// === User ===

export type UserType = 'customer' | 'provider' | 'admin';

export interface User {
  userId: string;
  email: string;
  phone: string;
  fullName: string;
  userType: UserType;
  avatarUrl?: string;
  createdAt: string;
}

// === Provider ===

export interface Provider {
  providerId: string;
  userId: string;
  businessName: string;
  vehicleType: string;
  vehiclePlate: string;
  isAvailable: boolean;
  lat?: number;
  lng?: number;
  averageRating: number;
  totalJobs: number;
}

// === Evidence ===

export type PhotoPosition =
  | 'FRONT' | 'REAR' | 'LEFT' | 'RIGHT'
  | 'FRONT_LEFT' | 'FRONT_RIGHT' | 'REAR_LEFT' | 'REAR_RIGHT';

export interface ConditionReport {
  reportId: string;
  bookingId: string;
  providerId: string;
  phase: 'pickup' | 'dropoff';
  media: MediaItem[];
  notes?: string;
  createdAt: string;
}

export interface MediaItem {
  mediaId: string;
  s3Key: string;
  position: PhotoPosition;
  mimeType: string;
  capturedAt: string;
}

// === Diagnosis ===

export interface DiagnosisResult {
  recommendedService: ServiceType;
  urgencyLevel: 'LOW' | 'MEDIUM' | 'HIGH' | 'CRITICAL';
  estimatedCostMin: number; // centavos
  estimatedCostMax: number;
  description: string;
  safetyWarnings: string[];
}

export interface DiagnoseRequest {
  description: string;
  photoUrls?: string[];
  vehicleType?: string;
  lat?: number;
  lng?: number;
}

// === Earnings ===

export interface EarningsPeriod {
  grossAmount: number; // centavos
  commission: number;
  netAmount: number;
  bookingCount: number;
}

export interface EarningsOutput {
  providerId: string;
  today: EarningsPeriod;
  thisWeek: EarningsPeriod;
  thisMonth: EarningsPeriod;
  allTime: EarningsPeriod;
}

// === API Error ===

export interface APIErrorResponse {
  error: {
    code: string;
    message: string;
  };
}
