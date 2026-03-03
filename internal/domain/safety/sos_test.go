package safety

import (
	"testing"
)

func TestComputeLevel(t *testing.T) {
	tests := []struct {
		name  string
		score int
		want  string
	}{
		{"zero score", 0, "low"},
		{"low boundary", 29, "low"},
		{"medium boundary", 30, "medium"},
		{"mid medium", 45, "medium"},
		{"high boundary", 60, "high"},
		{"mid high", 70, "high"},
		{"critical boundary", 80, "critical"},
		{"max score", 100, "critical"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &RiskScore{Score: tt.score}
			if got := r.ComputeLevel(); got != tt.want {
				t.Errorf("ComputeLevel() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestTriggerTypes(t *testing.T) {
	triggers := []TriggerType{
		TriggerTypeTripleTap,
		TriggerTypeShake,
		TriggerTypeCodeWord,
		TriggerTypeButton,
	}

	seen := make(map[TriggerType]bool)
	for _, trigger := range triggers {
		if trigger == "" {
			t.Error("trigger type should not be empty")
		}
		if seen[trigger] {
			t.Errorf("duplicate trigger type: %q", trigger)
		}
		seen[trigger] = true
	}

	if len(triggers) != 4 {
		t.Errorf("expected 4 trigger types, got %d", len(triggers))
	}
}
