package auth

import (
	"context"
	"log/slog"

	"github.com/aws/aws-lambda-go/events"

	"github.com/davidlramirez95/towcommand/internal/adapter/handler"
	"github.com/davidlramirez95/towcommand/internal/domain/user"
)

// PreTokenGenerationUseCase handles the Cognito PreTokenGeneration trigger.
// It injects custom claims (user_type, trust_tier, status, and optionally
// providerId) into the ID token.
type PreTokenGenerationUseCase struct {
	users     UserFinder
	providers ProviderByCognitoSubFinder
}

// NewPreTokenGenerationUseCase creates a PreTokenGenerationUseCase with its dependencies.
func NewPreTokenGenerationUseCase(users UserFinder, providers ProviderByCognitoSubFinder) *PreTokenGenerationUseCase {
	return &PreTokenGenerationUseCase{
		users:     users,
		providers: providers,
	}
}

// Execute looks up the user and injects custom claims into the token.
// On any error (DDB failure, user not found), it falls back to safe defaults
// so authentication is never blocked by infrastructure failures.
func (uc *PreTokenGenerationUseCase) Execute(ctx context.Context, event *events.CognitoEventUserPoolsPreTokenGen) (events.CognitoEventUserPoolsPreTokenGen, error) {
	cognitoSub := handler.CognitoUserName(&event.CognitoEventUserPoolsHeader)

	slog.Default().InfoContext(ctx, "pre-token-generation trigger",
		slog.String("cognitoSub", cognitoSub),
	)

	claims := map[string]string{
		"custom:user_type":  string(user.UserTypeCustomer),
		"custom:trust_tier": string(user.TrustTierBasic),
		"custom:status":     string(user.UserStatusActive),
	}

	u, err := uc.users.FindByID(ctx, cognitoSub)
	if err != nil {
		slog.Default().WarnContext(ctx, "failed to lookup user for token claims, using defaults",
			slog.String("cognitoSub", cognitoSub),
			slog.Any("error", err),
		)
		return handler.AddClaimsToToken(event, claims), nil
	}

	if u == nil {
		slog.Default().WarnContext(ctx, "user not found for token claims, using defaults",
			slog.String("cognitoSub", cognitoSub),
		)
		return handler.AddClaimsToToken(event, claims), nil
	}

	claims["custom:user_type"] = string(u.UserType)
	claims["custom:trust_tier"] = string(u.TrustTier)
	claims["custom:status"] = string(u.Status)

	// If the user is a provider, look up the provider to inject providerId.
	if u.UserType == user.UserTypeProvider {
		p, err := uc.providers.FindByCognitoSub(ctx, cognitoSub)
		if err != nil {
			slog.Default().WarnContext(ctx, "failed to lookup provider for token claims",
				slog.String("cognitoSub", cognitoSub),
				slog.Any("error", err),
			)
		} else if p != nil {
			claims["custom:provider_id"] = p.ProviderID
		}
	}

	slog.Default().InfoContext(ctx, "injecting custom claims into token",
		slog.String("cognitoSub", cognitoSub),
		slog.String("userType", claims["custom:user_type"]),
		slog.String("trustTier", claims["custom:trust_tier"]),
	)

	return handler.AddClaimsToToken(event, claims), nil
}
