package auth

import (
	"context"
	"errors"
	"testing"

	"github.com/aws/aws-lambda-go/events"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/davidlramirez95/towcommand/internal/domain/user"
)

func TestPreAuthenticationUseCase_Execute(t *testing.T) {
	tests := []struct {
		name       string
		userResult *user.User
		userErr    error
		wantErr    bool
		errMsg     string
	}{
		{
			name: "active user passes pre-auth",
			userResult: &user.User{
				UserID: "cognito-sub-active",
				Status: user.UserStatusActive,
			},
			wantErr: false,
		},
		{
			name: "banned user is blocked",
			userResult: &user.User{
				UserID: "cognito-sub-banned",
				Status: user.UserStatusBanned,
			},
			wantErr: true,
			errMsg:  "user account is banned",
		},
		{
			name: "suspended user is blocked",
			userResult: &user.User{
				UserID: "cognito-sub-suspended",
				Status: user.UserStatusSuspended,
			},
			wantErr: true,
			errMsg:  "user account is suspended",
		},
		{
			name:       "user not found allows auth (new user)",
			userResult: nil,
			userErr:    nil,
			wantErr:    false,
		},
		{
			name:       "DDB error allows auth (graceful degradation)",
			userResult: nil,
			userErr:    errors.New("dynamo timeout"),
			wantErr:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			userMock := new(mockUserFinder)
			uc := NewPreAuthenticationUseCase(userMock)

			cognitoSub := "cognito-sub-test"
			if tt.userResult != nil {
				cognitoSub = tt.userResult.UserID
			}

			event := &events.CognitoEventUserPoolsPreAuthentication{
				CognitoEventUserPoolsHeader: events.CognitoEventUserPoolsHeader{
					UserName: cognitoSub,
				},
				Request: events.CognitoEventUserPoolsPreAuthenticationRequest{
					UserAttributes: map[string]string{
						"email": "user@example.com",
					},
				},
			}

			userMock.On("FindByID", mock.Anything, cognitoSub).Return(tt.userResult, tt.userErr)

			result, err := uc.Execute(context.Background(), event)

			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
			} else {
				require.NoError(t, err)
			}

			assert.Equal(t, cognitoSub, result.UserName)
			userMock.AssertExpectations(t)
		})
	}
}
