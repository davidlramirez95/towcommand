package booking

import (
	"testing"
)

func TestCanTransition_ValidTransitions(t *testing.T) {
	tests := []struct {
		name string
		from BookingStatus
		to   BookingStatus
		want bool
	}{
		// Happy path: full linear flow
		{"pending to matched", BookingStatusPending, BookingStatusMatched, true},
		{"matched to en_route", BookingStatusMatched, BookingStatusEnRoute, true},
		{"en_route to arrived", BookingStatusEnRoute, BookingStatusArrived, true},
		{"arrived to condition_report", BookingStatusArrived, BookingStatusConditionReport, true},
		{"condition_report to otp_verified", BookingStatusConditionReport, BookingStatusOTPVerified, true},
		{"otp_verified to loading", BookingStatusOTPVerified, BookingStatusLoading, true},
		{"loading to in_transit", BookingStatusLoading, BookingStatusInTransit, true},
		{"in_transit to arrived_dropoff", BookingStatusInTransit, BookingStatusArrivedDropoff, true},
		{"arrived_dropoff to otp_dropoff", BookingStatusArrivedDropoff, BookingStatusOTPDropoff, true},
		{"otp_dropoff to completed", BookingStatusOTPDropoff, BookingStatusCompleted, true},

		// Cancellations
		{"pending to cancelled", BookingStatusPending, BookingStatusCancelled, true},
		{"matched to cancelled", BookingStatusMatched, BookingStatusCancelled, true},
		{"en_route to cancelled", BookingStatusEnRoute, BookingStatusCancelled, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CanTransition(tt.from, tt.to); got != tt.want {
				t.Errorf("CanTransition(%q, %q) = %v, want %v", tt.from, tt.to, got, tt.want)
			}
		})
	}
}

func TestCanTransition_InvalidTransitions(t *testing.T) {
	tests := []struct {
		name string
		from BookingStatus
		to   BookingStatus
	}{
		// Skip states
		{"pending to arrived (skip)", BookingStatusPending, BookingStatusArrived},
		{"matched to loading (skip)", BookingStatusMatched, BookingStatusLoading},
		{"en_route to completed (skip)", BookingStatusEnRoute, BookingStatusCompleted},

		// Backward transitions
		{"matched to pending (backward)", BookingStatusMatched, BookingStatusPending},
		{"completed to pending (backward)", BookingStatusCompleted, BookingStatusPending},
		{"in_transit to en_route (backward)", BookingStatusInTransit, BookingStatusEnRoute},

		// Cancel from non-cancellable states
		{"arrived to cancelled", BookingStatusArrived, BookingStatusCancelled},
		{"condition_report to cancelled", BookingStatusConditionReport, BookingStatusCancelled},
		{"loading to cancelled", BookingStatusLoading, BookingStatusCancelled},
		{"in_transit to cancelled", BookingStatusInTransit, BookingStatusCancelled},
		{"arrived_dropoff to cancelled", BookingStatusArrivedDropoff, BookingStatusCancelled},
		{"otp_dropoff to cancelled", BookingStatusOTPDropoff, BookingStatusCancelled},

		// Terminal states
		{"completed to matched", BookingStatusCompleted, BookingStatusMatched},
		{"cancelled to pending", BookingStatusCancelled, BookingStatusPending},

		// Self-transitions
		{"pending to pending", BookingStatusPending, BookingStatusPending},
		{"completed to completed", BookingStatusCompleted, BookingStatusCompleted},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if CanTransition(tt.from, tt.to) {
				t.Errorf("CanTransition(%q, %q) = true, want false", tt.from, tt.to)
			}
		})
	}
}

func TestCanTransition_UnknownStatus(t *testing.T) {
	if CanTransition("NONEXISTENT", BookingStatusPending) {
		t.Error("CanTransition from unknown status should return false")
	}
}

func TestIsFinal(t *testing.T) {
	tests := []struct {
		status BookingStatus
		want   bool
	}{
		{BookingStatusPending, false},
		{BookingStatusMatched, false},
		{BookingStatusEnRoute, false},
		{BookingStatusArrived, false},
		{BookingStatusConditionReport, false},
		{BookingStatusOTPVerified, false},
		{BookingStatusLoading, false},
		{BookingStatusInTransit, false},
		{BookingStatusArrivedDropoff, false},
		{BookingStatusOTPDropoff, false},
		{BookingStatusCompleted, true},
		{BookingStatusCancelled, true},
	}

	for _, tt := range tests {
		t.Run(string(tt.status), func(t *testing.T) {
			if got := IsFinal(tt.status); got != tt.want {
				t.Errorf("IsFinal(%q) = %v, want %v", tt.status, got, tt.want)
			}
		})
	}
}

func TestValidStatusTransitions_AllStatusesHaveEntry(t *testing.T) {
	allStatuses := []BookingStatus{
		BookingStatusPending,
		BookingStatusMatched,
		BookingStatusEnRoute,
		BookingStatusArrived,
		BookingStatusConditionReport,
		BookingStatusOTPVerified,
		BookingStatusLoading,
		BookingStatusInTransit,
		BookingStatusArrivedDropoff,
		BookingStatusOTPDropoff,
		BookingStatusCompleted,
		BookingStatusCancelled,
	}

	for _, s := range allStatuses {
		if _, ok := ValidStatusTransitions[s]; !ok {
			t.Errorf("ValidStatusTransitions missing entry for %q", s)
		}
	}

	if len(ValidStatusTransitions) != len(allStatuses) {
		t.Errorf("ValidStatusTransitions has %d entries, want %d", len(ValidStatusTransitions), len(allStatuses))
	}
}

func TestTerminalStatesHaveNoTransitions(t *testing.T) {
	for _, status := range []BookingStatus{BookingStatusCompleted, BookingStatusCancelled} {
		if targets := ValidStatusTransitions[status]; len(targets) != 0 {
			t.Errorf("terminal status %q has %d transitions, want 0", status, len(targets))
		}
	}
}
