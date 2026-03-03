import { ServiceType, WeightClass } from '../types';

export interface ServiceConfig {
  type: ServiceType;
  label: string;
  labelFil: string;
  description: string;
  requiresTruck: boolean;
  compatibleWeightClasses: WeightClass[];
  estimatedDurationMinutes: { min: number; max: number };
}

export const SERVICE_CONFIGS: Record<ServiceType, ServiceConfig> = {
  [ServiceType.FLATBED_TOW]: {
    type: ServiceType.FLATBED_TOW,
    label: 'Flatbed Tow',
    labelFil: 'Flatbed na Towing',
    description: 'Vehicle loaded onto flatbed truck for safe transport',
    requiresTruck: true,
    compatibleWeightClasses: [WeightClass.LIGHT, WeightClass.MEDIUM, WeightClass.HEAVY],
    estimatedDurationMinutes: { min: 30, max: 90 },
  },
  [ServiceType.WHEEL_LIFT]: {
    type: ServiceType.WHEEL_LIFT,
    label: 'Wheel Lift Tow',
    labelFil: 'Wheel Lift na Towing',
    description: 'Vehicle towed with front or rear wheels lifted',
    requiresTruck: true,
    compatibleWeightClasses: [WeightClass.LIGHT, WeightClass.MEDIUM],
    estimatedDurationMinutes: { min: 25, max: 75 },
  },
  [ServiceType.JUMPSTART]: {
    type: ServiceType.JUMPSTART,
    label: 'Jumpstart',
    labelFil: 'Jumpstart ng Baterya',
    description: 'Battery jumpstart service',
    requiresTruck: false,
    compatibleWeightClasses: [WeightClass.LIGHT, WeightClass.MEDIUM, WeightClass.HEAVY],
    estimatedDurationMinutes: { min: 15, max: 30 },
  },
  [ServiceType.TIRE_CHANGE]: {
    type: ServiceType.TIRE_CHANGE,
    label: 'Tire Change',
    labelFil: 'Pagpapalit ng Gulong',
    description: 'Flat tire replacement with spare',
    requiresTruck: false,
    compatibleWeightClasses: [WeightClass.MOTORCYCLE, WeightClass.LIGHT, WeightClass.MEDIUM],
    estimatedDurationMinutes: { min: 20, max: 45 },
  },
  [ServiceType.FUEL_DELIVERY]: {
    type: ServiceType.FUEL_DELIVERY,
    label: 'Fuel Delivery',
    labelFil: 'Hatid Gasolina',
    description: 'Emergency fuel delivery',
    requiresTruck: false,
    compatibleWeightClasses: [WeightClass.MOTORCYCLE, WeightClass.LIGHT, WeightClass.MEDIUM, WeightClass.HEAVY],
    estimatedDurationMinutes: { min: 20, max: 40 },
  },
  [ServiceType.LOCKOUT]: {
    type: ServiceType.LOCKOUT,
    label: 'Lockout Assistance',
    labelFil: 'Naka-lock Out',
    description: 'Vehicle lockout assistance',
    requiresTruck: false,
    compatibleWeightClasses: [WeightClass.LIGHT, WeightClass.MEDIUM],
    estimatedDurationMinutes: { min: 15, max: 45 },
  },
  [ServiceType.ACCIDENT_RECOVERY]: {
    type: ServiceType.ACCIDENT_RECOVERY,
    label: 'Accident Recovery',
    labelFil: 'Aksidente Recovery',
    description: 'Vehicle recovery after accident',
    requiresTruck: true,
    compatibleWeightClasses: [WeightClass.LIGHT, WeightClass.MEDIUM, WeightClass.HEAVY, WeightClass.SUPER_HEAVY],
    estimatedDurationMinutes: { min: 45, max: 120 },
  },
};
