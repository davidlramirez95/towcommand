import { create } from 'zustand';
import { persist } from 'zustand/middleware';
import { mmkvStorage } from '@/lib/storage/mmkv';

export interface Vehicle {
  id: string;
  make: string;
  model: string;
  year: number;
  plate: string;
  type: 'sedan' | 'suv' | 'pickup' | 'van' | 'motorcycle';
  color: string;
}

interface VehicleState {
  vehicles: Vehicle[];
  selectedVehicleId: string | null;
  selectedCondition: string | null;

  addVehicle: (vehicle: Vehicle) => void;
  removeVehicle: (id: string) => void;
  selectVehicle: (id: string) => void;
  selectCondition: (condition: string) => void;
  reset: () => void;
}

export const useVehicleStore = create<VehicleState>()(
  persist(
    (set) => ({
      vehicles: [],
      selectedVehicleId: null,
      selectedCondition: null,

      addVehicle: (vehicle) =>
        set((state) => ({ vehicles: [...state.vehicles, vehicle] })),

      removeVehicle: (id) =>
        set((state) => ({
          vehicles: state.vehicles.filter((v) => v.id !== id),
          selectedVehicleId: state.selectedVehicleId === id ? null : state.selectedVehicleId,
        })),

      selectVehicle: (id) => set({ selectedVehicleId: id }),

      selectCondition: (condition) => set({ selectedCondition: condition }),

      reset: () => set({ selectedVehicleId: null, selectedCondition: null }),
    }),
    {
      name: 'vehicle-storage',
      storage: mmkvStorage,
      partialize: (state) => ({ vehicles: state.vehicles }),
    },
  ),
);
