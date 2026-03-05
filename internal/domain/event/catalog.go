// Package event defines domain event sources and detail types for the TowCommand platform.
package event

// Event sources identify the bounded context that originated the event.
const (
	SourceBooking  = "tc.booking"
	SourceMatching = "tc.matching"
	SourceTracking = "tc.tracking"
	SourcePayment  = "tc.payment"
	SourceSOS      = "tc.sos"
	SourceAuth     = "tc.auth"
	SourceEvidence = "tc.evidence"
)

// Booking detail types.
const (
	BookingCreated       = "BookingCreated"
	BookingMatched       = "BookingMatched"
	BookingCancelled     = "BookingCancelled"
	BookingStatusChanged = "BookingStatusChanged"
	BookingCompleted     = "BookingCompleted"
)

// Matching detail types.
const (
	MatchingStarted   = "MatchingStarted"
	MatchingCompleted = "MatchingCompleted"
	MatchingFailed    = "MatchingFailed"
	ProviderAccepted  = "ProviderAccepted"
	ProviderDeclined  = "ProviderDeclined"
)

// Tracking detail types.
const (
	LocationUpdated = "LocationUpdated"
	ETAUpdated      = "ETAUpdated"
	GeofenceEntered = "GeofenceEntered"
	GeofenceExited  = "GeofenceExited"
)

// Payment detail types.
const (
	PaymentInitiated = "PaymentInitiated"
	PaymentCaptured  = "PaymentCaptured"
	PaymentRefunded  = "PaymentRefunded"
	PaymentFailed    = "PaymentFailed"
)

// SOS detail types.
const (
	SOSTriggered   = "SOSTriggered"
	SOSResolved    = "SOSResolved"
	RouteDeviation = "RouteDeviation"
)

// Auth detail types.
const (
	UserRegistered    = "UserRegistered"
	UserVerified      = "UserVerified"
	ProviderOnboarded = "ProviderOnboarded"
)

// Provider detail types.
const (
	SourceProvider      = "tc.provider"
	ProviderRegistered  = "ProviderRegistered"
	ProviderOnline      = "ProviderOnline"
	ProviderOffline     = "ProviderOffline"
	AvailabilityChanged = "AvailabilityChanged"
)

// Evidence detail types.
const (
	EvidenceUploaded       = "EvidenceUploaded"
	ConditionReportCreated = "ConditionReportCreated"
	EvidenceValidated      = "EvidenceValidated"
)

// Rating detail types.
const (
	SourceRating    = "tc.rating"
	RatingSubmitted = "RatingSubmitted"
)

// OTP detail types.
const (
	SourceOTP    = "tc.otp"
	OTPGenerated = "OTPGenerated"
	OTPVerified  = "OTPVerified"
)
