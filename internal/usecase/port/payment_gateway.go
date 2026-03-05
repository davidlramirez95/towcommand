package port

import "context"

// ChargeResult contains the gateway reference after a successful charge.
type ChargeResult struct {
	GatewayRef string
}

// RefundResult contains the gateway reference after a successful refund.
type RefundResult struct {
	GatewayRef string
}

// PaymentGateway defines the interface for external payment providers (PayMongo, etc.).
type PaymentGateway interface {
	Charge(ctx context.Context, paymentID string, amountCentavos int64, currency, method string) (*ChargeResult, error)
	Refund(ctx context.Context, gatewayRef string, amountCentavos int64) (*RefundResult, error)
	VerifyWebhookSignature(payload []byte, signature string) error
}
