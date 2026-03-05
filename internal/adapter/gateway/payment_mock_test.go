package gateway

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMockPaymentGateway_Charge(t *testing.T) {
	gw := NewMockPaymentGateway("test-secret")

	result, err := gw.Charge(context.Background(), "PAY-2026-ABC", 100_000, "PHP", "gcash")

	require.NoError(t, err)
	assert.True(t, strings.HasPrefix(result.GatewayRef, "mock-"), "gateway ref should start with mock-")
	assert.Greater(t, len(result.GatewayRef), len("mock-"), "gateway ref should have random suffix")
}

func TestMockPaymentGateway_Refund(t *testing.T) {
	gw := NewMockPaymentGateway("test-secret")

	result, err := gw.Refund(context.Background(), "mock-abc123", 50_000)

	require.NoError(t, err)
	assert.True(t, strings.HasPrefix(result.GatewayRef, "mock-refund-"), "refund ref should start with mock-refund-")
}

func TestMockPaymentGateway_VerifyWebhookSignature_Valid(t *testing.T) {
	secret := "my-webhook-secret"
	gw := NewMockPaymentGateway(secret)

	payload := []byte(`{"paymentId":"PAY-2026-001","event":"payment.captured"}`)

	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write(payload)
	validSig := hex.EncodeToString(mac.Sum(nil))

	err := gw.VerifyWebhookSignature(payload, validSig)
	assert.NoError(t, err)
}

func TestMockPaymentGateway_VerifyWebhookSignature_Invalid(t *testing.T) {
	gw := NewMockPaymentGateway("my-webhook-secret")

	payload := []byte(`{"paymentId":"PAY-2026-001","event":"payment.captured"}`)

	err := gw.VerifyWebhookSignature(payload, "bad-signature")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid webhook signature")
}

func TestMockPaymentGateway_SignPayload(t *testing.T) {
	secret := "sign-test-secret"
	gw := NewMockPaymentGateway(secret)

	payload := []byte(`test-payload`)
	sig := gw.SignPayload(payload)

	// Verify SignPayload produces a valid signature.
	err := gw.VerifyWebhookSignature(payload, sig)
	assert.NoError(t, err)
}

func TestMockPaymentGateway_MultipleChargesUnique(t *testing.T) {
	gw := NewMockPaymentGateway("test-secret")
	ctx := context.Background()

	r1, err := gw.Charge(ctx, "PAY-1", 100, "PHP", "gcash")
	require.NoError(t, err)

	r2, err := gw.Charge(ctx, "PAY-2", 200, "PHP", "maya")
	require.NoError(t, err)

	assert.NotEqual(t, r1.GatewayRef, r2.GatewayRef, "each charge should produce a unique ref")
}
