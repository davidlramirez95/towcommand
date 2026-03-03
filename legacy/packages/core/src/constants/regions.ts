export interface ServiceRegion {
  code: string;
  name: string;
  nameFil: string;
  defaultRadiusKm: number;
  surgeMultiplierCap: number;
  isActive: boolean;
}

export const PH_REGIONS: Record<string, ServiceRegion> = {
  NCR: { code: 'NCR', name: 'National Capital Region', nameFil: 'Kalakhang Maynila', defaultRadiusKm: 10, surgeMultiplierCap: 1.5, isActive: true },
  CALABARZON: { code: 'CALABARZON', name: 'CALABARZON', nameFil: 'CALABARZON', defaultRadiusKm: 25, surgeMultiplierCap: 1.5, isActive: true },
  CENTRAL_LUZON: { code: 'CENTRAL_LUZON', name: 'Central Luzon', nameFil: 'Gitnang Luzon', defaultRadiusKm: 25, surgeMultiplierCap: 1.5, isActive: false },
};

export const MMDA_SECTORS: Record<string, string[]> = {
  NORTH: ['Caloocan', 'Malabon', 'Navotas', 'Valenzuela', 'Quezon City (North)'],
  EAST: ['Marikina', 'San Juan', 'Mandaluyong', 'Pasig', 'Quezon City (East)'],
  WEST: ['Manila', 'Quezon City (West)'],
  SOUTH: ['Makati', 'Taguig', 'Paranaque', 'Las Pinas', 'Muntinlupa'],
  CENTRAL: ['Pasay', 'Pateros'],
};
