package auth

import (
	"context"
	"testing"

	"github.com/aws/aws-lambda-go/events"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCustomMessageUseCase_Execute(t *testing.T) {
	tests := []struct {
		name             string
		triggerSource    string
		wantSubject      string
		wantEmailContain string
		wantSMSContain   string
		wantDefault      bool
	}{
		{
			name:             "sign up verification",
			triggerSource:    "CustomMessage_SignUp",
			wantSubject:      "Mabuhay! Welcome to TowCommand - Verify Your Account",
			wantEmailContain: "Mabuhay! Welcome to TowCommand",
			wantSMSContain:   "Mabuhay! TowCommand verification code: {####}",
		},
		{
			name:             "forgot password",
			triggerSource:    "CustomMessage_ForgotPassword",
			wantSubject:      "TowCommand - Password Reset Code",
			wantEmailContain: "password reset code is {####}",
			wantSMSContain:   "password reset code: {####}",
		},
		{
			name:             "resend code",
			triggerSource:    "CustomMessage_ResendCode",
			wantSubject:      "TowCommand - Verification Code",
			wantEmailContain: "verification code is {####}",
			wantSMSContain:   "verification code: {####}",
		},
		{
			name:             "admin create user",
			triggerSource:    "CustomMessage_AdminCreateUser",
			wantSubject:      "Mabuhay! Welcome to TowCommand",
			wantEmailContain: "temporary password is {####}",
			wantSMSContain:   "temporary password: {####}",
		},
		{
			name:          "unknown trigger uses defaults",
			triggerSource: "CustomMessage_Unknown",
			wantDefault:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			uc := NewCustomMessageUseCase()

			event := &events.CognitoEventUserPoolsCustomMessage{
				CognitoEventUserPoolsHeader: events.CognitoEventUserPoolsHeader{
					TriggerSource: tt.triggerSource,
					UserName:      "test-user",
				},
				Response: events.CognitoEventUserPoolsCustomMessageResponse{
					EmailSubject: "Default Subject",
					EmailMessage: "Default Body",
					SMSMessage:   "Default SMS",
				},
			}

			result, err := uc.Execute(context.Background(), event)

			require.NoError(t, err)

			if tt.wantDefault {
				assert.Equal(t, "Default Subject", result.Response.EmailSubject)
				assert.Equal(t, "Default Body", result.Response.EmailMessage)
				assert.Equal(t, "Default SMS", result.Response.SMSMessage)
				return
			}

			assert.Equal(t, tt.wantSubject, result.Response.EmailSubject)
			assert.Contains(t, result.Response.EmailMessage, tt.wantEmailContain)
			assert.Contains(t, result.Response.SMSMessage, tt.wantSMSContain)
			// All templates must contain the code placeholder.
			assert.Contains(t, result.Response.EmailMessage, "{####}")
			assert.Contains(t, result.Response.SMSMessage, "{####}")
		})
	}
}
