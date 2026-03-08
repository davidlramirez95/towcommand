import { create } from 'zustand';

type BookingStatus =
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

interface ActiveBooking {
  bookingId: string;
  status: BookingStatus;
  providerName?: string;
  providerPhone?: string;
  vehiclePlate?: string;
  providerLat?: number;
  providerLng?: number;
  pickupLat: number;
  pickupLng: number;
  dropoffLat?: number;
  dropoffLng?: number;
  estimatedCost?: number;
  eta?: number;
}

interface BookingState {
  activeBooking: ActiveBooking | null;

  setActiveBooking: (booking: ActiveBooking | null) => void;
  updateStatus: (status: BookingStatus) => void;
  updateProviderLocation: (lat: number, lng: number) => void;
  updateETA: (eta: number) => void;
  reset: () => void;
}

export const useBookingStore = create<BookingState>()((set) => ({
  activeBooking: null,

  setActiveBooking: (booking) => set({ activeBooking: booking }),
  updateStatus: (status) =>
    set((state) =>
      state.activeBooking ? { activeBooking: { ...state.activeBooking, status } } : state,
    ),
  updateProviderLocation: (lat, lng) =>
    set((state) =>
      state.activeBooking
        ? { activeBooking: { ...state.activeBooking, providerLat: lat, providerLng: lng } }
        : state,
    ),
  updateETA: (eta) =>
    set((state) =>
      state.activeBooking ? { activeBooking: { ...state.activeBooking, eta } } : state,
    ),
  reset: () => set({ activeBooking: null }),
}));
