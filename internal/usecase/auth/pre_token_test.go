package auth

import (
	"context"
	"errors"
	"testing"

	"github.com/aws/aws-lambda-go/events"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/davidlramirez95/towcommand/internal/domain/provider"
	"github.com/davidlramirez95/towcommand/internal/domain/user"
)

// --- Mocks for PreToken ---

type mockUserFinder struct{ mock.Mock }

func (m *mockUserFinder) FindByID(ctx context.Context, userID string) (*user.User, error) {
	args := m.Called(ctx, userID)
	u, _ := args.Get(0).(*user.User)
	return u, args.Error(1)
}

type mockProviderByCognitoSubFinder struct{ mock.Mock }

func (m *mockProviderByCognitoSubFinder) FindByCognitoSub(ctx context.Context, cognitoSub string) (*provider.Provider, error) {
	args := m.Called(ctx, cognitoSub)
	p, _ := args.Get(0).(*provider.Provider)
	return p, args.Error(1)
}

// --- Tests ---

func TestPreTokenGenerationUseCase_Execute_CustomerUser(t *testing.T) {
	users := new(mockUserFinder)
	providers := new(mockProviderByCognitoSubFinder)

	uc := NewPreTokenGenerationUseCase(users, providers)

	event := &events.CognitoEventUserPoolsPreTokenGen{
		CognitoEventUserPoolsHeader: events.CognitoEventUserPoolsHeader{
			UserName: "cognito-sub-cust",
		},
	}

	users.On("FindByID", mock.Anything, "cognito-sub-cust").Return(&user.User{
		UserID:    "cognito-sub-cust",
		UserType:  user.UserTypeCustomer,
		TrustTier: user.TrustTierSukiGold,
		Status:    user.UserStatusActive,
	}, nil)

	result, err := uc.Execute(context.Background(), event)

	require.NoError(t, err)
	claims := result.Response.ClaimsOverrideDetails.ClaimsToAddOrOverride
	assert.Equal(t, "customer", claims["custom:user_type"])
	assert.Equal(t, "suki_gold", claims["custom:trust_tier"])
	assert.Equal(t, "active", claims["custom:status"])
	assert.Empty(t, claims["custom:provider_id"])
	users.AssertExpectations(t)
	providers.AssertNotCalled(t, "FindByCognitoSub")
}

func TestPreTokenGenerationUseCase_Execute_ProviderUser(t *testing.T) {
	users := new(mockUserFinder)
	providers := new(mockProviderByCognitoSubFinder)

	uc := NewPreTokenGenerationUseCase(users, providers)

	event := &events.CognitoEventUserPoolsPreTokenGen{
		CognitoEventUserPoolsHeader: events.CognitoEventUserPoolsHeader{
			UserName: "cognito-sub-prov",
		},
	}

	users.On("FindByID", mock.Anything, "cognito-sub-prov").Return(&user.User{
		UserID:    "cognito-sub-prov",
		UserType:  user.UserTypeProvider,
		TrustTier: user.TrustTierVerified,
		Status:    user.UserStatusActive,
	}, nil)

	providers.On("FindByCognitoSub", mock.Anything, "cognito-sub-prov").Return(&provider.Provider{
		ProviderID: "PROV-abc123",
	}, nil)

	result, err := uc.Execute(context.Background(), event)

	require.NoError(t, err)
	claims := result.Response.ClaimsOverrideDetails.ClaimsToAddOrOverride
	assert.Equal(t, "provider", claims["custom:user_type"])
	assert.Equal(t, "verified", claims["custom:trust_tier"])
	assert.Equal(t, "active", claims["custom:status"])
	assert.Equal(t, "PROV-abc123", claims["custom:provider_id"])
	users.AssertExpectations(t)
	providers.AssertExpectations(t)
}

func TestPreTokenGenerationUseCase_Execute_UserNotFound_UsesDefaults(t *testing.T) {
	users := new(mockUserFinder)
	providers := new(mockProviderByCognitoSubFinder)

	uc := NewPreTokenGenerationUseCase(users, providers)

	event := &events.CognitoEventUserPoolsPreTokenGen{
		CognitoEventUserPoolsHeader: events.CognitoEventUserPoolsHeader{
			UserName: "cognito-sub-missing",
		},
	}

	users.On("FindByID", mock.Anything, "cognito-sub-missing").Return(nil, nil)

	result, err := uc.Execute(context.Background(), event)

	require.NoError(t, err)
	claims := result.Response.ClaimsOverrideDetails.ClaimsToAddOrOverride
	assert.Equal(t, "customer", claims["custom:user_type"])
	assert.Equal(t, "basic", claims["custom:trust_tier"])
	assert.Equal(t, "active", claims["custom:status"])
	users.AssertExpectations(t)
	providers.AssertNotCalled(t, "FindByCognitoSub")
}

func TestPreTokenGenerationUseCase_Execute_UserFinderError_UsesDefaults(t *testing.T) {
	users := new(mockUserFinder)
	providers := new(mockProviderByCognitoSubFinder)

	uc := NewPreTokenGenerationUseCase(users, providers)

	event := &events.CognitoEventUserPoolsPreTokenGen{
		CognitoEventUserPoolsHeader: events.CognitoEventUserPoolsHeader{
			UserName: "cognito-sub-err",
		},
	}

	users.On("FindByID", mock.Anything, "cognito-sub-err").Return(nil, errors.New("dynamo timeout"))

	result, err := uc.Execute(context.Background(), event)

	require.NoError(t, err)
	claims := result.Response.ClaimsOverrideDetails.ClaimsToAddOrOverride
	assert.Equal(t, "customer", claims["custom:user_type"])
	assert.Equal(t, "basic", claims["custom:trust_tier"])
	assert.Equal(t, "active", claims["custom:status"])
	users.AssertExpectations(t)
	providers.AssertNotCalled(t, "FindByCognitoSub")
}

func TestPreTokenGenerationUseCase_Execute_ProviderLookupError_SkipsProviderID(t *testing.T) {
	users := new(mockUserFinder)
	providers := new(mockProviderByCognitoSubFinder)

	uc := NewPreTokenGenerationUseCase(users, providers)

	event := &events.CognitoEventUserPoolsPreTokenGen{
		CognitoEventUserPoolsHeader: events.CognitoEventUserPoolsHeader{
			UserName: "cognito-sub-prov-err",
		},
	}

	users.On("FindByID", mock.Anything, "cognito-sub-prov-err").Return(&user.User{
		UserID:    "cognito-sub-prov-err",
		UserType:  user.UserTypeProvider,
		TrustTier: user.TrustTierBasic,
		Status:    user.UserStatusActive,
	}, nil)

	providers.On("FindByCognitoSub", mock.Anything, "cognito-sub-prov-err").
		Return(nil, errors.New("scan failed"))

	result, err := uc.Execute(context.Background(), event)

	require.NoError(t, err)
	claims := result.Response.ClaimsOverrideDetails.ClaimsToAddOrOverride
	assert.Equal(t, "provider", claims["custom:user_type"])
	assert.Empty(t, claims["custom:provider_id"])
	users.AssertExpectations(t)
	providers.AssertExpectations(t)
}

func TestPreTokenGenerationUseCase_Execute_ProviderNotFound_SkipsProviderID(t *testing.T) {
	users := new(mockUserFinder)
	providers := new(mockProviderByCognitoSubFinder)

	uc := NewPreTokenGenerationUseCase(users, providers)

	event := &events.CognitoEventUserPoolsPreTokenGen{
		CognitoEventUserPoolsHeader: events.CognitoEventUserPoolsHeader{
			UserName: "cognito-sub-prov-nil",
		},
	}

	users.On("FindByID", mock.Anything, "cognito-sub-prov-nil").Return(&user.User{
		UserID:    "cognito-sub-prov-nil",
		UserType:  user.UserTypeProvider,
		TrustTier: user.TrustTierBasic,
		Status:    user.UserStatusActive,
	}, nil)

	providers.On("FindByCognitoSub", mock.Anything, "cognito-sub-prov-nil").
		Return(nil, nil)

	result, err := uc.Execute(context.Background(), event)

	require.NoError(t, err)
	claims := result.Response.ClaimsOverrideDetails.ClaimsToAddOrOverride
	assert.Equal(t, "provider", claims["custom:user_type"])
	assert.Empty(t, claims["custom:provider_id"])
	users.AssertExpectations(t)
	providers.AssertExpectations(t)
}
