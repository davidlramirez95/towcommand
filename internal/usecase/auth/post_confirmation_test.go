package auth

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/davidlramirez95/towcommand/internal/domain/user"
	"github.com/davidlramirez95/towcommand/internal/usecase/port"
)

// --- Mocks for PostConfirmation ---

type mockUserSaver struct{ mock.Mock }

func (m *mockUserSaver) Save(ctx context.Context, u *user.User) error {
	args := m.Called(ctx, u)
	return args.Error(0)
}

type mockEventPublisher struct{ mock.Mock }

func (m *mockEventPublisher) Publish(ctx context.Context, source, detailType string, detail any, actor *port.Actor) error {
	args := m.Called(ctx, source, detailType, detail, actor)
	return args.Error(0)
}

// --- Tests ---

func TestPostConfirmationUseCase_Execute_Success(t *testing.T) {
	saver := new(mockUserSaver)
	pub := new(mockEventPublisher)

	uc := NewPostConfirmationUseCase(saver, pub)
	fixedTime := time.Date(2026, 3, 4, 10, 0, 0, 0, time.UTC)
	uc.now = func() time.Time { return fixedTime }

	event := &events.CognitoEventUserPoolsPostConfirmation{
		CognitoEventUserPoolsHeader: events.CognitoEventUserPoolsHeader{
			TriggerSource: "PostConfirmation_ConfirmSignUp",
			UserName:      "cognito-sub-abc",
		},
		Request: events.CognitoEventUserPoolsPostConfirmationRequest{
			UserAttributes: map[string]string{
				"email":        "juan@example.com",
				"phone_number": "+639171234567",
				"name":         "Juan Dela Cruz",
			},
		},
	}

	saver.On("Save", mock.Anything, mock.MatchedBy(func(u *user.User) bool {
		return u.UserID == "cognito-sub-abc" &&
			u.CognitoSub == "cognito-sub-abc" &&
			u.Email == "juan@example.com" &&
			u.Phone == "+639171234567" &&
			u.Name == "Juan Dela Cruz" &&
			u.UserType == user.UserTypeCustomer &&
			u.TrustTier == user.TrustTierBasic &&
			u.Language == user.LanguageEnglish &&
			u.Status == user.UserStatusActive &&
			u.CreatedAt.Equal(fixedTime) &&
			u.UpdatedAt.Equal(fixedTime)
	})).Return(nil)

	pub.On("Publish", mock.Anything, "tc.auth", "UserRegistered", mock.Anything, mock.Anything).Return(nil)

	result, err := uc.Execute(context.Background(), event)

	require.NoError(t, err)
	assert.Equal(t, "cognito-sub-abc", result.UserName)
	saver.AssertExpectations(t)
	pub.AssertExpectations(t)
}

func TestPostConfirmationUseCase_Execute_SaveError_DoesNotReturnError(t *testing.T) {
	saver := new(mockUserSaver)
	pub := new(mockEventPublisher)

	uc := NewPostConfirmationUseCase(saver, pub)
	uc.now = func() time.Time { return time.Now().UTC() }

	event := &events.CognitoEventUserPoolsPostConfirmation{
		CognitoEventUserPoolsHeader: events.CognitoEventUserPoolsHeader{
			TriggerSource: "PostConfirmation_ConfirmSignUp",
			UserName:      "cognito-sub-err",
		},
		Request: events.CognitoEventUserPoolsPostConfirmationRequest{
			UserAttributes: map[string]string{
				"email": "fail@example.com",
			},
		},
	}

	saver.On("Save", mock.Anything, mock.Anything).Return(errors.New("dynamo timeout"))

	result, err := uc.Execute(context.Background(), event)

	// Must not return error even if save fails.
	require.NoError(t, err)
	assert.Equal(t, "cognito-sub-err", result.UserName)
	saver.AssertExpectations(t)
	// Publish should not be called if save fails.
	pub.AssertNotCalled(t, "Publish")
}

func TestPostConfirmationUseCase_Execute_PublishError_DoesNotReturnError(t *testing.T) {
	saver := new(mockUserSaver)
	pub := new(mockEventPublisher)

	uc := NewPostConfirmationUseCase(saver, pub)
	uc.now = func() time.Time { return time.Now().UTC() }

	event := &events.CognitoEventUserPoolsPostConfirmation{
		CognitoEventUserPoolsHeader: events.CognitoEventUserPoolsHeader{
			TriggerSource: "PostConfirmation_ConfirmSignUp",
			UserName:      "cognito-sub-pub-err",
		},
		Request: events.CognitoEventUserPoolsPostConfirmationRequest{
			UserAttributes: map[string]string{
				"email": "pubfail@example.com",
			},
		},
	}

	saver.On("Save", mock.Anything, mock.Anything).Return(nil)
	pub.On("Publish", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
		Return(errors.New("eventbridge error"))

	result, err := uc.Execute(context.Background(), event)

	// Must not return error even if publish fails.
	require.NoError(t, err)
	assert.Equal(t, "cognito-sub-pub-err", result.UserName)
	saver.AssertExpectations(t)
	pub.AssertExpectations(t)
}

func TestPostConfirmationUseCase_Execute_EmptyAttributes(t *testing.T) {
	saver := new(mockUserSaver)
	pub := new(mockEventPublisher)

	uc := NewPostConfirmationUseCase(saver, pub)
	uc.now = func() time.Time { return time.Now().UTC() }

	event := &events.CognitoEventUserPoolsPostConfirmation{
		CognitoEventUserPoolsHeader: events.CognitoEventUserPoolsHeader{
			TriggerSource: "PostConfirmation_ConfirmSignUp",
			UserName:      "cognito-sub-empty",
		},
		Request: events.CognitoEventUserPoolsPostConfirmationRequest{
			UserAttributes: map[string]string{},
		},
	}

	saver.On("Save", mock.Anything, mock.MatchedBy(func(u *user.User) bool {
		return u.UserID == "cognito-sub-empty" &&
			u.Email == "" &&
			u.Phone == "" &&
			u.Name == ""
	})).Return(nil)

	pub.On("Publish", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)

	result, err := uc.Execute(context.Background(), event)

	require.NoError(t, err)
	assert.Equal(t, "cognito-sub-empty", result.UserName)
	saver.AssertExpectations(t)
}
