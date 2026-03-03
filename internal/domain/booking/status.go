package booking

// BookingStatus represents the current state in the booking lifecycle.
type BookingStatus string

const (
	BookingStatusPending         BookingStatus = "PENDING"
	BookingStatusMatched         BookingStatus = "MATCHED"
	BookingStatusEnRoute         BookingStatus = "EN_ROUTE"
	BookingStatusArrived         BookingStatus = "ARRIVED"
	BookingStatusConditionReport BookingStatus = "CONDITION_REPORT"
	BookingStatusOTPVerified     BookingStatus = "OTP_VERIFIED"
	BookingStatusLoading         BookingStatus = "LOADING"
	BookingStatusInTransit       BookingStatus = "IN_TRANSIT"
	BookingStatusArrivedDropoff  BookingStatus = "ARRIVED_DROPOFF"
	BookingStatusOTPDropoff      BookingStatus = "OTP_DROPOFF"
	BookingStatusCompleted       BookingStatus = "COMPLETED"
	BookingStatusCancelled       BookingStatus = "CANCELLED"
)

// ValidStatusTransitions defines the allowed state transitions for a booking.
// The linear flow is: PENDING → MATCHED → EN_ROUTE → ARRIVED → CONDITION_REPORT
// → OTP_VERIFIED → LOADING → IN_TRANSIT → ARRIVED_DROPOFF → OTP_DROPOFF → COMPLETED.
// CANCELLED is reachable from PENDING, MATCHED, and EN_ROUTE.
var ValidStatusTransitions = map[BookingStatus][]BookingStatus{
	BookingStatusPending:         {BookingStatusMatched, BookingStatusCancelled},
	BookingStatusMatched:         {BookingStatusEnRoute, BookingStatusCancelled},
	BookingStatusEnRoute:         {BookingStatusArrived, BookingStatusCancelled},
	BookingStatusArrived:         {BookingStatusConditionReport},
	BookingStatusConditionReport: {BookingStatusOTPVerified},
	BookingStatusOTPVerified:     {BookingStatusLoading},
	BookingStatusLoading:         {BookingStatusInTransit},
	BookingStatusInTransit:       {BookingStatusArrivedDropoff},
	BookingStatusArrivedDropoff:  {BookingStatusOTPDropoff},
	BookingStatusOTPDropoff:      {BookingStatusCompleted},
	BookingStatusCompleted:       {},
	BookingStatusCancelled:       {},
}

// CanTransition reports whether a booking can transition from one status to another.
func CanTransition(from, to BookingStatus) bool {
	targets, ok := ValidStatusTransitions[from]
	if !ok {
		return false
	}
	for _, t := range targets {
		if t == to {
			return true
		}
	}
	return false
}

// IsFinal reports whether the given status is a terminal state.
func IsFinal(status BookingStatus) bool {
	return status == BookingStatusCompleted || status == BookingStatusCancelled
}
