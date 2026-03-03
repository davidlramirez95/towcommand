package payment

import "time"

// PaymentMethod represents a supported payment channel.
type PaymentMethod string

const (
	PaymentMethodGCash     PaymentMethod = "gcash"
	PaymentMethodMaya      PaymentMethod = "maya"
	PaymentMethodCard      PaymentMethod = "card"
	PaymentMethodCash      PaymentMethod = "cash"
	PaymentMethodCorporate PaymentMethod = "corporate"
)

// PaymentStatus represents the lifecycle state of a payment.
type PaymentStatus string

const (
	PaymentStatusPending   PaymentStatus = "pending"
	PaymentStatusHeld      PaymentStatus = "held"
	PaymentStatusCaptured  PaymentStatus = "captured"
	PaymentStatusRefunded  PaymentStatus = "refunded"
	PaymentStatusFailed    PaymentStatus = "failed"
	PaymentStatusCancelled PaymentStatus = "cancelled"
)

// PayoutStatus represents the lifecycle state of a provider payout.
type PayoutStatus string

const (
	PayoutStatusPending    PayoutStatus = "pending"
	PayoutStatusProcessing PayoutStatus = "processing"
	PayoutStatusCompleted  PayoutStatus = "completed"
	PayoutStatusFailed     PayoutStatus = "failed"
)

// Payment represents a customer payment for a booking. All amounts are in centavos.
type Payment struct {
	PaymentID    string        `json:"paymentId" validate:"required"`
	BookingID    string        `json:"bookingId" validate:"required"`
	UserID       string        `json:"userId" validate:"required"`
	Amount       int64         `json:"amount" validate:"required,gt=0"`
	Currency     string        `json:"currency" validate:"required,eq=PHP"`
	Method       PaymentMethod `json:"method" validate:"required"`
	Status       PaymentStatus `json:"status" validate:"required"`
	HoldAmount   int64         `json:"holdAmount,omitempty"`
	GatewayRef   string        `json:"gatewayRef,omitempty"`
	CapturedAt   *time.Time    `json:"capturedAt,omitempty"`
	RefundedAt   *time.Time    `json:"refundedAt,omitempty"`
	RefundReason string        `json:"refundReason,omitempty"`
	CreatedAt    time.Time     `json:"createdAt"`
	UpdatedAt    time.Time     `json:"updatedAt"`
}

// ProviderPayout represents a disbursement to a provider. All amounts are in centavos.
type ProviderPayout struct {
	PayoutID       string        `json:"payoutId" validate:"required"`
	ProviderID     string        `json:"providerId" validate:"required"`
	BookingID      string        `json:"bookingId" validate:"required"`
	GrossAmount    int64         `json:"grossAmount" validate:"required,gt=0"`
	Commission     int64         `json:"commission" validate:"min=0"`
	CommissionRate float64       `json:"commissionRate" validate:"min=0,max=1"`
	NetAmount      int64         `json:"netAmount" validate:"required,gt=0"`
	Status         PayoutStatus  `json:"status" validate:"required"`
	PayoutMethod   PaymentMethod `json:"payoutMethod" validate:"required"`
	ScheduledAt    time.Time     `json:"scheduledAt"`
	CompletedAt    *time.Time    `json:"completedAt,omitempty"`
}
