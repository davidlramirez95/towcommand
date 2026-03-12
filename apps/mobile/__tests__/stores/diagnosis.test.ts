/**
 * Diagnosis Store Tests — 2nd Order Logic
 *
 * 2nd order concerns:
 * - toggleSymptom must add AND remove (prevents stuck selections)
 * - startAnalysis must change step (prevents UI from staying on selector)
 * - reset must clear ALL fields (prevents stale diagnosis leaking into next booking)
 * - Empty symptoms + startAnalysis = still changes step (UI handles empty state)
 */
import { useDiagnosisStore } from '@/stores/diagnosis';

// Reset store before each test
beforeEach(() => {
  useDiagnosisStore.getState().reset();
});

describe('DiagnosisStore', () => {
  it('starts with select step and empty symptoms', () => {
    const state = useDiagnosisStore.getState();
    expect(state.step).toBe('select');
    expect(state.selectedSymptoms).toEqual([]);
    expect(state.diagnosisResult).toBeNull();
  });

  it('toggleSymptom adds a symptom', () => {
    useDiagnosisStore.getState().toggleSymptom('tire');
    expect(useDiagnosisStore.getState().selectedSymptoms).toEqual(['tire']);
  });

  it('toggleSymptom removes an already-selected symptom', () => {
    useDiagnosisStore.getState().toggleSymptom('tire');
    useDiagnosisStore.getState().toggleSymptom('tire');
    expect(useDiagnosisStore.getState().selectedSymptoms).toEqual([]);
  });

  it('toggleSymptom handles multiple symptoms', () => {
    const { toggleSymptom } = useDiagnosisStore.getState();
    toggleSymptom('tire');
    toggleSymptom('battery');
    toggleSymptom('fuel');
    expect(useDiagnosisStore.getState().selectedSymptoms).toEqual(['tire', 'battery', 'fuel']);

    toggleSymptom('battery');
    expect(useDiagnosisStore.getState().selectedSymptoms).toEqual(['tire', 'fuel']);
  });

  it('startAnalysis changes step to analyzing', () => {
    useDiagnosisStore.getState().startAnalysis();
    expect(useDiagnosisStore.getState().step).toBe('analyzing');
  });

  it('setResults sets step to results and stores diagnosis', () => {
    const result = {
      needsTow: false,
      confidence: 92,
      recommendedService: 'On-Site Mechanic',
      estimatedSavings: 2500,
      alternatives: [{ service: 'Full Tow', priceRange: '₱1,800–₱3,500', eta: '12 min' }],
    };

    useDiagnosisStore.getState().setResults(result);
    const state = useDiagnosisStore.getState();
    expect(state.step).toBe('results');
    expect(state.diagnosisResult).toEqual(result);
  });

  it('reset clears ALL fields (prevents stale data in next booking)', () => {
    const { toggleSymptom, startAnalysis, setResults, reset } = useDiagnosisStore.getState();
    toggleSymptom('tire');
    startAnalysis();
    setResults({
      needsTow: true,
      confidence: 85,
      recommendedService: 'Flatbed Tow',
      estimatedSavings: 0,
      alternatives: [],
    });

    reset();
    const state = useDiagnosisStore.getState();
    expect(state.step).toBe('select');
    expect(state.selectedSymptoms).toEqual([]);
    expect(state.diagnosisResult).toBeNull();
  });
});
