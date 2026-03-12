import { create } from 'zustand';

interface DiagnosisResult {
  needsTow: boolean;
  confidence: number;
  recommendedService: string;
  estimatedSavings: number;
  alternatives: Array<{ service: string; priceRange: string; eta: string }>;
}

interface DiagnosisState {
  step: 'select' | 'analyzing' | 'results';
  selectedSymptoms: string[];
  diagnosisResult: DiagnosisResult | null;

  toggleSymptom: (symptomId: string) => void;
  startAnalysis: () => void;
  setResults: (result: DiagnosisResult) => void;
  reset: () => void;
}

export const useDiagnosisStore = create<DiagnosisState>()((set) => ({
  step: 'select',
  selectedSymptoms: [],
  diagnosisResult: null,

  toggleSymptom: (symptomId) =>
    set((state) => ({
      selectedSymptoms: state.selectedSymptoms.includes(symptomId)
        ? state.selectedSymptoms.filter((id) => id !== symptomId)
        : [...state.selectedSymptoms, symptomId],
    })),

  startAnalysis: () => set({ step: 'analyzing' }),

  setResults: (result) => set({ step: 'results', diagnosisResult: result }),

  reset: () => set({ step: 'select', selectedSymptoms: [], diagnosisResult: null }),
}));
