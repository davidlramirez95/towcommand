export enum UserType {
  CUSTOMER = 'customer',
  PROVIDER = 'provider',
  FLEET_MANAGER = 'fleet_manager',
  OPS_AGENT = 'ops_agent',
  ADMIN = 'admin',
}

export enum TrustTier {
  BASIC = 'basic',
  VERIFIED = 'verified',
  SUKI_SILVER = 'suki_silver',
  SUKI_GOLD = 'suki_gold',
  SUKI_ELITE = 'suki_elite',
}

export interface User {
  userId: string;
  cognitoSub: string;
  email: string;
  phone: string;
  name: string;
  userType: UserType;
  trustTier: TrustTier;
  language: 'en' | 'fil';
  status: 'active' | 'suspended' | 'banned';
  createdAt: string;
  updatedAt: string;
}

export interface UserVehicle {
  vehicleId: string;
  userId: string;
  make: string;
  model: string;
  year: number;
  plateNumber: string;
  weightClass: WeightClass;
  color: string;
  photoUrl?: string;
  isDefault: boolean;
  createdAt: string;
}

export enum WeightClass {
  MOTORCYCLE = 'motorcycle',
  LIGHT = 'light',
  MEDIUM = 'medium',
  HEAVY = 'heavy',
  SUPER_HEAVY = 'super_heavy',
}
