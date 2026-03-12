/**
 * Vehicle Store Tests — 2nd Order Logic
 *
 * 2nd order concerns:
 * - removeVehicle must clear selectedVehicleId if it was selected (prevents ghost selection)
 * - reset must NOT clear vehicles array (persisted data survives session reset)
 * - reset must clear selectedVehicleId and selectedCondition (per-session state)
 * - addVehicle must append, not replace (prevents losing saved vehicles)
 */
import { useVehicleStore } from '@/stores/vehicle';

const mockVehicle = {
  id: 'v1',
  make: 'Mitsubishi',
  model: 'Montero GLS',
  year: 2026,
  plate: 'ABC 1234',
  type: 'suv' as const,
  color: 'White',
};

const mockVehicle2 = {
  id: 'v2',
  make: 'Toyota',
  model: 'Vios',
  year: 2025,
  plate: 'XYZ 5678',
  type: 'sedan' as const,
  color: 'Silver',
};

// Reset store before each test
beforeEach(() => {
  useVehicleStore.setState({ vehicles: [], selectedVehicleId: null, selectedCondition: null });
});

describe('VehicleStore', () => {
  it('starts with empty vehicles', () => {
    const state = useVehicleStore.getState();
    expect(state.vehicles).toEqual([]);
    expect(state.selectedVehicleId).toBeNull();
    expect(state.selectedCondition).toBeNull();
  });

  it('addVehicle appends to vehicles array', () => {
    useVehicleStore.getState().addVehicle(mockVehicle);
    useVehicleStore.getState().addVehicle(mockVehicle2);
    expect(useVehicleStore.getState().vehicles).toHaveLength(2);
    expect(useVehicleStore.getState().vehicles[0].plate).toBe('ABC 1234');
    expect(useVehicleStore.getState().vehicles[1].plate).toBe('XYZ 5678');
  });

  it('removeVehicle removes by id', () => {
    useVehicleStore.getState().addVehicle(mockVehicle);
    useVehicleStore.getState().addVehicle(mockVehicle2);
    useVehicleStore.getState().removeVehicle('v1');
    expect(useVehicleStore.getState().vehicles).toHaveLength(1);
    expect(useVehicleStore.getState().vehicles[0].id).toBe('v2');
  });

  it('removeVehicle clears selectedVehicleId if removed vehicle was selected (prevents ghost selection)', () => {
    useVehicleStore.getState().addVehicle(mockVehicle);
    useVehicleStore.getState().selectVehicle('v1');
    expect(useVehicleStore.getState().selectedVehicleId).toBe('v1');

    useVehicleStore.getState().removeVehicle('v1');
    expect(useVehicleStore.getState().selectedVehicleId).toBeNull();
  });

  it('removeVehicle does NOT clear selectedVehicleId if different vehicle removed', () => {
    useVehicleStore.getState().addVehicle(mockVehicle);
    useVehicleStore.getState().addVehicle(mockVehicle2);
    useVehicleStore.getState().selectVehicle('v1');

    useVehicleStore.getState().removeVehicle('v2');
    expect(useVehicleStore.getState().selectedVehicleId).toBe('v1');
  });

  it('selectVehicle sets selectedVehicleId', () => {
    useVehicleStore.getState().selectVehicle('v1');
    expect(useVehicleStore.getState().selectedVehicleId).toBe('v1');
  });

  it('selectCondition sets selectedCondition', () => {
    useVehicleStore.getState().selectCondition("Engine won't start");
    expect(useVehicleStore.getState().selectedCondition).toBe("Engine won't start");
  });

  it('reset clears session state but NOT vehicles (persisted data survives)', () => {
    useVehicleStore.getState().addVehicle(mockVehicle);
    useVehicleStore.getState().selectVehicle('v1');
    useVehicleStore.getState().selectCondition('Flat tire(s)');

    useVehicleStore.getState().reset();
    const state = useVehicleStore.getState();
    expect(state.vehicles).toHaveLength(1); // vehicles preserved
    expect(state.selectedVehicleId).toBeNull(); // session state cleared
    expect(state.selectedCondition).toBeNull(); // session state cleared
  });
});
