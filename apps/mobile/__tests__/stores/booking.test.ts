/**
 * Booking Store Tests — 2nd Order Logic
 *
 * 2nd order concern: Booking store is updated from WebSocket messages
 * (high frequency location updates) AND user actions (status changes).
 * If updateProviderLocation creates new object refs each time, every
 * subscriber re-renders 1x/second. If reset() misses a field, stale
 * provider data shows on the NEXT booking.
 */
import { useBookingStore } from '@/stores/booking';

beforeEach(() => {
  useBookingStore.getState().reset();
});

describe('booking store', () => {
  const mockBooking = {
    id: 'BK-001',
    status: 'MATCHED' as const,
    providerName: 'Juan Driver',
    providerPhone: '+639171111111',
    vehiclePlate: 'ABC 1234',
    pickupLat: 14.5995,
    pickupLng: 120.9842,
    estimatedCost: 250000,
  };

  it('initial state has no active booking', () => {
    expect(useBookingStore.getState().activeBooking).toBeNull();
  });

  it('setActiveBooking stores full booking data', () => {
    useBookingStore.getState().setActiveBooking(mockBooking);
    const booking = useBookingStore.getState().activeBooking;

    expect(booking?.id).toBe('BK-001');
    expect(booking?.status).toBe('MATCHED');
    expect(booking?.providerName).toBe('Juan Driver');
  });

  it('updateStatus transitions correctly (WS booking_status event)', () => {
    useBookingStore.getState().setActiveBooking(mockBooking);
    useBookingStore.getState().updateStatus('EN_ROUTE');

    expect(useBookingStore.getState().activeBooking?.status).toBe('EN_ROUTE');
  });

  it('updateStatus with no active booking is a no-op (prevents null crash)', () => {
    // No booking set — updateStatus should not throw
    expect(() => {
      useBookingStore.getState().updateStatus('EN_ROUTE');
    }).not.toThrow();
    expect(useBookingStore.getState().activeBooking).toBeNull();
  });

  it('updateProviderLocation stores lat/lng (WS location_update, ~1/sec)', () => {
    useBookingStore.getState().setActiveBooking(mockBooking);
    useBookingStore.getState().updateProviderLocation(14.6000, 120.9850);

    const booking = useBookingStore.getState().activeBooking;
    expect(booking?.providerLat).toBe(14.6000);
    expect(booking?.providerLng).toBe(120.9850);
  });

  it('updateProviderLocation with no booking is a no-op', () => {
    expect(() => {
      useBookingStore.getState().updateProviderLocation(14.6, 120.98);
    }).not.toThrow();
  });

  it('updateETA stores ETA value (WS eta_update event)', () => {
    useBookingStore.getState().setActiveBooking(mockBooking);
    useBookingStore.getState().updateETA(8);

    expect(useBookingStore.getState().activeBooking?.eta).toBe(8);
  });

  it('updateETA with no booking is a no-op', () => {
    expect(() => {
      useBookingStore.getState().updateETA(5);
    }).not.toThrow();
  });

  it('reset() clears ALL booking fields — prevents stale provider data on next booking', () => {
    useBookingStore.getState().setActiveBooking({
      ...mockBooking,
      providerLat: 14.6,
      providerLng: 120.98,
      eta: 5,
      dropoffLat: 14.55,
      dropoffLng: 120.99,
    });

    useBookingStore.getState().reset();

    expect(useBookingStore.getState().activeBooking).toBeNull();
  });

  it('multiple rapid location updates dont lose data (WS flood scenario)', () => {
    useBookingStore.getState().setActiveBooking(mockBooking);

    // Simulate 5 rapid GPS updates
    for (let i = 0; i < 5; i++) {
      useBookingStore.getState().updateProviderLocation(14.5 + i * 0.001, 120.9 + i * 0.001);
    }

    const booking = useBookingStore.getState().activeBooking;
    // Last update wins
    expect(booking?.providerLat).toBeCloseTo(14.504, 3);
    expect(booking?.providerLng).toBeCloseTo(120.904, 3);
  });

  it('status update preserves other booking fields (partial update safety)', () => {
    useBookingStore.getState().setActiveBooking(mockBooking);
    useBookingStore.getState().updateProviderLocation(14.6, 120.98);
    useBookingStore.getState().updateETA(7);

    // Status update should NOT wipe location/ETA
    useBookingStore.getState().updateStatus('ARRIVED');

    const booking = useBookingStore.getState().activeBooking;
    expect(booking?.status).toBe('ARRIVED');
    expect(booking?.providerLat).toBe(14.6);
    expect(booking?.eta).toBe(7);
    expect(booking?.providerName).toBe('Juan Driver');
  });
});
