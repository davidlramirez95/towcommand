import { WeightClass } from './user';

export enum BookingStatus {
  PENDING = 'PENDING',
  MATCHED = 'MATCHED',
  EN_ROUTE = 'EN_ROUTE',
  ARRIVED = 'ARRIVED',
  CONDITION_REPORT = 'CONDITION_REPORT',
  OTP_VERIFIED = 'OTP_VERIFIED',
  LOADING = 'LOADING',
  IN_TRANSIT = 'IN_TRANSIT',
  ARRIVED_DROPOFF = 'ARRIVED_DROPOFF',
  OTP_DROPOFF = 'OTP_DROPOFF',
  COMPLETED = 'COMPLETED',
  CANCELLED = 'CANCELLED',
}

export enum ServiceType {
  FLATBED_TOW = 'FLATBED_TOW',
  WHEEL_LIFT = 'WHEEL_LIFT',
  JUMPSTART = 'JUMPSTART',
  TIRE_CHANGE = 'TIRE_CHANGE',
  FUEL_DELIVERY = 'FUEL_DELIVERY',
  LOCKOUT = 'LOCKOUT',
  ACCIDENT_RECOVERY = 'ACCIDENT_RECOVERY',
}

export interface GeoLocation {
  lat: number;
  lng: number;
  address?: string;
}

export interface Booking {
  bookingId: string;
  customerId: string;
  providerId?: string;
  vehicleId: string;
  serviceType: ServiceType;
  status: BookingStatus;
  pickupLocation: GeoLocation;
  dropoffLocation: GeoLocation;
  weightClass: WeightClass;
  price: PriceBreakdown;
  estimateId: string;
  notes?: string;
  cancellationReason?: string;
  cancellationFee?: number;
  matchedAt?: string;
  completedAt?: string;
  createdAt: string;
  updatedAt: string;
}

export interface PriceBreakdown {
  base: number;
  distance: number;
  weight: number;
  timeSurcharge: number;
  surgePricing: number;
  total: number;
  currency: 'PHP';
  surgeMultiplier?: number;
}

export interface BookingEstimate {
  estimateId: string;
  pickupLocation: GeoLocation;
  dropoffLocation: GeoLocation;
  serviceType: ServiceType;
  weightClass: WeightClass;
  distanceKm: number;
  price: PriceBreakdown;
  availableProviders: number;
  estimatedEtaMinutes: number;
  expiresAt: string;
  createdAt: string;
}

export const VALID_STATUS_TRANSITIONS: Record<BookingStatus, BookingStatus[]> = {
  [BookingStatus.PENDING]: [BookingStatus.MATCHED, BookingStatus.CANCELLED],
  [BookingStatus.MATCHED]: [BookingStatus.EN_ROUTE, BookingStatus.CANCELLED],
  [BookingStatus.EN_ROUTE]: [BookingStatus.ARRIVED, BookingStatus.CANCELLED],
  [BookingStatus.ARRIVED]: [BookingStatus.CONDITION_REPORT],
  [BookingStatus.CONDITION_REPORT]: [BookingStatus.OTP_VERIFIED],
  [BookingStatus.OTP_VERIFIED]: [BookingStatus.LOADING],
  [BookingStatus.LOADING]: [BookingStatus.IN_TRANSIT],
  [BookingStatus.IN_TRANSIT]: [BookingStatus.ARRIVED_DROPOFF],
  [BookingStatus.ARRIVED_DROPOFF]: [BookingStatus.OTP_DROPOFF],
  [BookingStatus.OTP_DROPOFF]: [BookingStatus.COMPLETED],
  [BookingStatus.COMPLETED]: [],
  [BookingStatus.CANCELLED]: [],
};
