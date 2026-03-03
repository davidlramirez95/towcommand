import { WeightClass } from './user';

export interface VehicleDatabase {
  make: string;
  model: string;
  yearRange: { from: number; to: number };
  defaultWeightClass: WeightClass;
  estimatedWeightKg: number;
}

export const COMMON_PH_VEHICLES: VehicleDatabase[] = [
  { make: 'Toyota', model: 'Vios', yearRange: { from: 2013, to: 2026 }, defaultWeightClass: WeightClass.LIGHT, estimatedWeightKg: 1050 },
  { make: 'Toyota', model: 'Innova', yearRange: { from: 2016, to: 2026 }, defaultWeightClass: WeightClass.LIGHT, estimatedWeightKg: 1800 },
  { make: 'Toyota', model: 'Fortuner', yearRange: { from: 2016, to: 2026 }, defaultWeightClass: WeightClass.MEDIUM, estimatedWeightKg: 2200 },
  { make: 'Toyota', model: 'Hilux', yearRange: { from: 2015, to: 2026 }, defaultWeightClass: WeightClass.MEDIUM, estimatedWeightKg: 2100 },
  { make: 'Mitsubishi', model: 'Montero Sport', yearRange: { from: 2016, to: 2026 }, defaultWeightClass: WeightClass.MEDIUM, estimatedWeightKg: 2100 },
  { make: 'Mitsubishi', model: 'Mirage', yearRange: { from: 2012, to: 2026 }, defaultWeightClass: WeightClass.LIGHT, estimatedWeightKg: 900 },
  { make: 'Honda', model: 'City', yearRange: { from: 2014, to: 2026 }, defaultWeightClass: WeightClass.LIGHT, estimatedWeightKg: 1100 },
  { make: 'Honda', model: 'Civic', yearRange: { from: 2016, to: 2026 }, defaultWeightClass: WeightClass.LIGHT, estimatedWeightKg: 1300 },
  { make: 'Nissan', model: 'Navara', yearRange: { from: 2015, to: 2026 }, defaultWeightClass: WeightClass.MEDIUM, estimatedWeightKg: 2050 },
  { make: 'Suzuki', model: 'Ertiga', yearRange: { from: 2019, to: 2026 }, defaultWeightClass: WeightClass.LIGHT, estimatedWeightKg: 1135 },
  { make: 'Ford', model: 'Ranger', yearRange: { from: 2015, to: 2026 }, defaultWeightClass: WeightClass.MEDIUM, estimatedWeightKg: 2200 },
  { make: 'Isuzu', model: 'D-Max', yearRange: { from: 2012, to: 2026 }, defaultWeightClass: WeightClass.MEDIUM, estimatedWeightKg: 1900 },
  { make: 'Honda', model: 'Click 125i', yearRange: { from: 2018, to: 2026 }, defaultWeightClass: WeightClass.MOTORCYCLE, estimatedWeightKg: 110 },
  { make: 'Yamaha', model: 'NMAX', yearRange: { from: 2020, to: 2026 }, defaultWeightClass: WeightClass.MOTORCYCLE, estimatedWeightKg: 131 },
];
