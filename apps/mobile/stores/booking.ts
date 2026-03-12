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

interface PriceBreakdown {
  baseFare: number;
  distanceCharge: number;
  weightSurcharge: number;
  timeSurcharge: number;
  platformFee: number;
  total: number;
}

interface PaymentMethod {
  type: 'gcash' | 'maya' | 'card' | 'cash';
  label: string;
  last4?: string;
}

interface MatchedProvider {
  id: string;
  name: string;
  rating: number;
  jobCount: number;
  plate: string;
  eta: number;
  verified: boolean;
}

type MatchingState = 'idle' | 'searching' | 'found' | 'timeout';

interface BookingState {
  activeBooking: ActiveBooking | null;
  priceBreakdown: PriceBreakdown | null;
  paymentMethod: PaymentMethod | null;
  matchingState: MatchingState;
  matchedProvider: MatchedProvider | null;
  otp: string | null;

  setActiveBooking: (booking: ActiveBooking | null) => void;
  updateStatus: (status: BookingStatus) => void;
  updateProviderLocation: (lat: number, lng: number) => void;
  updateETA: (eta: number) => void;
  setPriceBreakdown: (breakdown: PriceBreakdown | null) => void;
  setPaymentMethod: (method: PaymentMethod | null) => void;
  setMatchingState: (state: MatchingState) => void;
  setMatchedProvider: (provider: MatchedProvider | null) => void;
  setOtp: (otp: string | null) => void;
  reset: () => void;
}

export const useBookingStore = create<BookingState>()((set) => ({
  activeBooking: null,
  priceBreakdown: null,
  paymentMethod: null,
  matchingState: 'idle',
  matchedProvider: null,
  otp: null,

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
  setPriceBreakdown: (breakdown) => set({ priceBreakdown: breakdown }),
  setPaymentMethod: (method) => set({ paymentMethod: method }),
  setMatchingState: (matchingState) => set({ matchingState }),
  setMatchedProvider: (provider) => set({ matchedProvider: provider }),
  setOtp: (otp) => set({ otp }),
  reset: () =>
    set({
      activeBooking: null,
      priceBreakdown: null,
      paymentMethod: null,
      matchingState: 'idle',
      matchedProvider: null,
      otp: null,
    }),
}));
