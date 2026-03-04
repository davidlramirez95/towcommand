package handler_test

import (
	"testing"

	"github.com/aws/aws-lambda-go/events"
	"github.com/stretchr/testify/assert"

	"github.com/davidlramirez95/towcommand/internal/adapter/handler"
)

// ---------------------------------------------------------------------------
// AutoConfirmUser
// ---------------------------------------------------------------------------

func TestAutoConfirmUser(t *testing.T) {
	tests := []struct {
		name            string
		autoVerifyEmail bool
	}{
		{"with email verification", true},
		{"without email verification", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			event := &events.CognitoEventUserPoolsPreSignup{}
			result := handler.AutoConfirmUser(event, tt.autoVerifyEmail)

			assert.True(t, result.Response.AutoConfirmUser)
			assert.Equal(t, tt.autoVerifyEmail, result.Response.AutoVerifyEmail)
		})
	}
}

func TestAutoConfirmUser_MutatesInput(t *testing.T) {
	event := &events.CognitoEventUserPoolsPreSignup{}
	_ = handler.AutoConfirmUser(event, true)
	// Pointer-based: the original event IS mutated (by design for pointer receivers)
	assert.True(t, event.Response.AutoConfirmUser)
}

// ---------------------------------------------------------------------------
// PreSignUpUserAttributes
// ---------------------------------------------------------------------------

func TestPreSignUpUserAttributes(t *testing.T) {
	event := &events.CognitoEventUserPoolsPreSignup{
		Request: events.CognitoEventUserPoolsPreSignupRequest{
			UserAttributes: map[string]string{
				"email": "user@example.com",
				"phone": "+639171234567",
			},
		},
	}
	attrs := handler.PreSignUpUserAttributes(event)
	assert.Equal(t, "user@example.com", attrs["email"])
	assert.Equal(t, "+639171234567", attrs["phone"])
}

// ---------------------------------------------------------------------------
// PostConfirmationUserAttributes
// ---------------------------------------------------------------------------

func TestPostConfirmationUserAttributes(t *testing.T) {
	event := &events.CognitoEventUserPoolsPostConfirmation{
		Request: events.CognitoEventUserPoolsPostConfirmationRequest{
			UserAttributes: map[string]string{
				"sub":   "user-123",
				"email": "user@example.com",
			},
		},
	}
	attrs := handler.PostConfirmationUserAttributes(event)
	assert.Equal(t, "user-123", attrs["sub"])
	assert.Equal(t, "user@example.com", attrs["email"])
}

// ---------------------------------------------------------------------------
// AddClaimsToToken
// ---------------------------------------------------------------------------

func TestAddClaimsToToken(t *testing.T) {
	event := &events.CognitoEventUserPoolsPreTokenGen{}
	claims := map[string]string{
		"custom:role":       "admin",
		"custom:trust_tier": "suki_gold",
	}

	result := handler.AddClaimsToToken(event, claims)

	assert.Equal(t, "admin", result.Response.ClaimsOverrideDetails.ClaimsToAddOrOverride["custom:role"])
	assert.Equal(t, "suki_gold", result.Response.ClaimsOverrideDetails.ClaimsToAddOrOverride["custom:trust_tier"])
}

func TestAddClaimsToToken_MutatesInput(t *testing.T) {
	event := &events.CognitoEventUserPoolsPreTokenGen{}
	_ = handler.AddClaimsToToken(event, map[string]string{"key": "val"})
	// Pointer-based: the original event IS mutated
	assert.Equal(t, "val", event.Response.ClaimsOverrideDetails.ClaimsToAddOrOverride["key"])
}

// ---------------------------------------------------------------------------
// SuppressClaims
// ---------------------------------------------------------------------------

func TestSuppressClaims(t *testing.T) {
	event := &events.CognitoEventUserPoolsPreTokenGen{}
	result := handler.SuppressClaims(event, []string{"email", "phone"})
	assert.Equal(t, []string{"email", "phone"}, result.Response.ClaimsOverrideDetails.ClaimsToSuppress)
}

// ---------------------------------------------------------------------------
// PreAuthUserAttributes
// ---------------------------------------------------------------------------

func TestPreAuthUserAttributes(t *testing.T) {
	event := &events.CognitoEventUserPoolsPreAuthentication{
		Request: events.CognitoEventUserPoolsPreAuthenticationRequest{
			UserAttributes: map[string]string{
				"email": "auth@example.com",
			},
		},
	}
	attrs := handler.PreAuthUserAttributes(event)
	assert.Equal(t, "auth@example.com", attrs["email"])
}

// ---------------------------------------------------------------------------
// CustomizeMessage
// ---------------------------------------------------------------------------

func TestCustomizeMessage(t *testing.T) {
	tests := []struct {
		name         string
		emailSubject string
		emailMessage string
		smsMessage   string
		wantSubject  string
		wantEmail    string
		wantSMS      string
	}{
		{
			name:         "all fields set",
			emailSubject: "Welcome to TowCommand",
			emailMessage: "Your code is {####}",
			smsMessage:   "Code: {####}",
			wantSubject:  "Welcome to TowCommand",
			wantEmail:    "Your code is {####}",
			wantSMS:      "Code: {####}",
		},
		{
			name:         "only email subject",
			emailSubject: "Verify your account",
			wantSubject:  "Verify your account",
			wantEmail:    "Default Body",
			wantSMS:      "Default SMS",
		},
		{
			name:        "empty strings preserve defaults",
			wantSubject: "Default Subject",
			wantEmail:   "Default Body",
			wantSMS:     "Default SMS",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			event := &events.CognitoEventUserPoolsCustomMessage{
				Response: events.CognitoEventUserPoolsCustomMessageResponse{
					EmailSubject: "Default Subject",
					EmailMessage: "Default Body",
					SMSMessage:   "Default SMS",
				},
			}
			result := handler.CustomizeMessage(event, tt.emailSubject, tt.emailMessage, tt.smsMessage)

			assert.Equal(t, tt.wantSubject, result.Response.EmailSubject)
			assert.Equal(t, tt.wantEmail, result.Response.EmailMessage)
			assert.Equal(t, tt.wantSMS, result.Response.SMSMessage)
		})
	}
}

// ---------------------------------------------------------------------------
// CognitoUserName / CognitoTriggerSource
// ---------------------------------------------------------------------------

func TestCognitoUserName(t *testing.T) {
	header := &events.CognitoEventUserPoolsHeader{UserName: "john.doe"}
	assert.Equal(t, "john.doe", handler.CognitoUserName(header))
}

func TestCognitoTriggerSource(t *testing.T) {
	header := &events.CognitoEventUserPoolsHeader{TriggerSource: handler.TriggerPreSignUp}
	assert.Equal(t, handler.TriggerPreSignUp, handler.CognitoTriggerSource(header))
}

// ---------------------------------------------------------------------------
// Trigger constants
// ---------------------------------------------------------------------------

func TestTriggerConstants(t *testing.T) {
	assert.Equal(t, "PreSignUp_SignUp", handler.TriggerPreSignUp)
	assert.Equal(t, "PostConfirmation_ConfirmSignUp", handler.TriggerPostConfirmation)
	assert.Equal(t, "TokenGeneration_HostedAuth", handler.TriggerPreTokenGeneration)
	assert.Equal(t, "PreAuthentication_Authentication", handler.TriggerPreAuthentication)
	assert.Equal(t, "CustomMessage_SignUp", handler.TriggerCustomMessage)
}
