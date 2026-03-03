import { TrustTier } from './user';

export enum ProviderStatus {
  PENDING_VERIFICATION = 'pending_verification',
  ACTIVE = 'active',
  SUSPENDED = 'suspended',
  DEACTIVATED = 'deactivated',
}

export enum TruckType {
  FLATBED = 'flatbed',
  WHEEL_LIFT = 'wheel_lift',
  BOOM = 'boom',
  MOTORCYCLE_CARRIER = 'motorcycle_carrier',
}

export interface Provider {
  providerId: string;
  cognitoSub: string;
  name: string;
  phone: string;
  email: string;
  status: ProviderStatus;
  trustTier: TrustTier;
  truckType: TruckType;
  maxWeightCapacityKg: number;
  plateNumber: string;
  ltoRegistration: string;
  nbiClearanceStatus: 'pending' | 'approved' | 'expired';
  drugTestStatus: 'pending' | 'approved' | 'expired';
  mmadAccredited: boolean;
  rating: number;
  totalJobsCompleted: number;
  acceptanceRate: number;
  isOnline: boolean;
  currentLat?: number;
  currentLng?: number;
  lastLocationUpdate?: string;
  serviceAreas: string[];
  createdAt: string;
  updatedAt: string;
}

export interface ProviderDoc {
  providerId: string;
  docType: 'nbi_clearance' | 'lto_registration' | 'drug_test' | 'vehicle_inspection' | 'insurance';
  s3Key: string;
  status: 'pending' | 'approved' | 'rejected';
  reviewedBy?: string;
  reviewedAt?: string;
  expiresAt?: string;
  uploadedAt: string;
}

export interface ProviderLocation {
  providerId: string;
  lat: number;
  lng: number;
  heading: number;
  speed: number;
  timestamp: string;
}
