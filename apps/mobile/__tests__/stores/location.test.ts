/**
 * Location Store Tests — 2nd Order Logic
 *
 * 2nd order concern: Location store is written to at GPS frequency (~1/sec)
 * by the provider's background location task. reset() must also clear
 * isTracking flag — otherwise the UI shows "tracking" indicator after
 * the provider goes offline, and background task may not restart properly.
 */
import { useLocationStore } from '@/stores/location';

beforeEach(() => {
  useLocationStore.getState().reset();
});

describe('location store', () => {
  it('initial state has null coordinates and tracking off', () => {
    const state = useLocationStore.getState();
    expect(state.latitude).toBeNull();
    expect(state.longitude).toBeNull();
    expect(state.heading).toBeNull();
    expect(state.accuracy).toBeNull();
    expect(state.isTracking).toBe(false);
  });

  it('setLocation updates all coordinate fields', () => {
    useLocationStore.getState().setLocation(14.5995, 120.9842, 45.5, 10);

    const state = useLocationStore.getState();
    expect(state.latitude).toBe(14.5995);
    expect(state.longitude).toBe(120.9842);
    expect(state.heading).toBe(45.5);
    expect(state.accuracy).toBe(10);
  });

  it('setTracking toggles tracking state (controls background GPS task)', () => {
    useLocationStore.getState().setTracking(true);
    expect(useLocationStore.getState().isTracking).toBe(true);

    useLocationStore.getState().setTracking(false);
    expect(useLocationStore.getState().isTracking).toBe(false);
  });

  it('reset clears coordinates AND isTracking (prevents phantom tracking)', () => {
    useLocationStore.getState().setLocation(14.5, 120.9, 90, 5);
    useLocationStore.getState().setTracking(true);

    useLocationStore.getState().reset();

    const state = useLocationStore.getState();
    expect(state.latitude).toBeNull();
    expect(state.longitude).toBeNull();
    expect(state.isTracking).toBe(false);
    expect(state.heading).toBeNull();
  });

  it('rapid setLocation calls settle on last value (GPS flood)', () => {
    for (let i = 0; i < 100; i++) {
      useLocationStore.getState().setLocation(14.5 + i * 0.0001, 120.9 + i * 0.0001);
    }

    const state = useLocationStore.getState();
    expect(state.latitude).toBeCloseTo(14.5099, 4);
    expect(state.longitude).toBeCloseTo(120.9099, 4);
  });
});
