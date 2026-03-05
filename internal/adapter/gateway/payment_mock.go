package gateway

import (
	"context"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"

	"github.com/davidlramirez95/towcommand/internal/usecase/port"
)

// MockPaymentGateway is a fake payment gateway for development and testing.
// It auto-succeeds on Charge and Refund, and verifies webhook signatures
// using HMAC-SHA256 with a configurable secret.
type MockPaymentGateway struct {
	webhookSecret string
}

// NewMockPaymentGateway constructs a MockPaymentGateway with the given webhook secret.
func NewMockPaymentGateway(webhookSecret string) *MockPaymentGateway {
	return &MockPaymentGateway{webhookSecret: webhookSecret}
}

// Charge simulates a successful charge and returns a mock gateway reference.
func (g *MockPaymentGateway) Charge(_ context.Context, _ string, _ int64, _, _ string) (*port.ChargeResult, error) {
	ref, err := randomHex(16)
	if err != nil {
		return nil, fmt.Errorf("generating mock ref: %w", err)
	}
	return &port.ChargeResult{GatewayRef: "mock-" + ref}, nil
}

// Refund simulates a successful refund and returns a mock gateway reference.
func (g *MockPaymentGateway) Refund(_ context.Context, _ string, _ int64) (*port.RefundResult, error) {
	ref, err := randomHex(16)
	if err != nil {
		return nil, fmt.Errorf("generating mock refund ref: %w", err)
	}
	return &port.RefundResult{GatewayRef: "mock-refund-" + ref}, nil
}

// VerifyWebhookSignature validates that the signature matches the HMAC-SHA256
// of the payload using the configured webhook secret.
func (g *MockPaymentGateway) VerifyWebhookSignature(payload []byte, signature string) error {
	mac := hmac.New(sha256.New, []byte(g.webhookSecret))
	mac.Write(payload)
	expected := hex.EncodeToString(mac.Sum(nil))
	if !hmac.Equal([]byte(expected), []byte(signature)) {
		return fmt.Errorf("invalid webhook signature")
	}
	return nil
}

// SignPayload creates an HMAC-SHA256 signature for a payload using the
// gateway's webhook secret. This is intended for testing webhook flows.
func (g *MockPaymentGateway) SignPayload(payload []byte) string {
	mac := hmac.New(sha256.New, []byte(g.webhookSecret))
	mac.Write(payload)
	return hex.EncodeToString(mac.Sum(nil))
}

// randomHex generates n random bytes and returns them as a hex string.
func randomHex(n int) (string, error) {
	b := make([]byte, n)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}
