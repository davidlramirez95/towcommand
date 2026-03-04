package auth

import (
	"context"
	"testing"

	"github.com/aws/aws-lambda-go/events"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPreSignUpUseCase_Execute(t *testing.T) {
	tests := []struct {
		name              string
		stage             string
		triggerSource     string
		wantAutoConfirm   bool
		wantAutoVerifyEml bool
	}{
		{
			name:              "social login auto-confirms and auto-verifies email",
			stage:             "prod",
			triggerSource:     "PreSignUp_ExternalProvider",
			wantAutoConfirm:   true,
			wantAutoVerifyEml: true,
		},
		{
			name:              "social login with Google prefix",
			stage:             "prod",
			triggerSource:     "PreSignUp_ExternalProvider_Google",
			wantAutoConfirm:   true,
			wantAutoVerifyEml: true,
		},
		{
			name:              "dev stage auto-confirms",
			stage:             "dev",
			triggerSource:     "PreSignUp_SignUp",
			wantAutoConfirm:   true,
			wantAutoVerifyEml: true,
		},
		{
			name:              "local stage auto-confirms",
			stage:             "local",
			triggerSource:     "PreSignUp_SignUp",
			wantAutoConfirm:   true,
			wantAutoVerifyEml: true,
		},
		{
			name:              "prod stage normal signup is unchanged",
			stage:             "prod",
			triggerSource:     "PreSignUp_SignUp",
			wantAutoConfirm:   false,
			wantAutoVerifyEml: false,
		},
		{
			name:              "staging stage normal signup is unchanged",
			stage:             "staging",
			triggerSource:     "PreSignUp_SignUp",
			wantAutoConfirm:   false,
			wantAutoVerifyEml: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			uc := NewPreSignUpUseCase(tt.stage)

			event := &events.CognitoEventUserPoolsPreSignup{
				CognitoEventUserPoolsHeader: events.CognitoEventUserPoolsHeader{
					TriggerSource: tt.triggerSource,
					UserName:      "test-user-123",
				},
				Request: events.CognitoEventUserPoolsPreSignupRequest{
					UserAttributes: map[string]string{
						"email": "test@example.com",
					},
				},
			}

			result, err := uc.Execute(context.Background(), event)

			require.NoError(t, err)
			assert.Equal(t, tt.wantAutoConfirm, result.Response.AutoConfirmUser)
			assert.Equal(t, tt.wantAutoVerifyEml, result.Response.AutoVerifyEmail)
		})
	}
}
