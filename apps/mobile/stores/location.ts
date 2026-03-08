import { create } from 'zustand';

interface LocationState {
  latitude: number | null;
  longitude: number | null;
  heading: number | null;
  accuracy: number | null;
  isTracking: boolean;

  setLocation: (lat: number, lng: number, heading?: number, accuracy?: number) => void;
  setTracking: (tracking: boolean) => void;
  reset: () => void;
}

export const useLocationStore = create<LocationState>()((set) => ({
  latitude: null,
  longitude: null,
  heading: null,
  accuracy: null,
  isTracking: false,

  setLocation: (latitude, longitude, heading, accuracy) =>
    set({ latitude, longitude, heading: heading ?? null, accuracy: accuracy ?? null }),
  setTracking: (isTracking) => set({ isTracking }),
  reset: () => set({ latitude: null, longitude: null, heading: null, accuracy: null, isTracking: false }),
}));
