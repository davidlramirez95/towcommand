package auth

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/aws/aws-lambda-go/events"

	"github.com/davidlramirez95/towcommand/internal/adapter/handler"
	"github.com/davidlramirez95/towcommand/internal/domain/user"
)

// PreAuthenticationUseCase handles the Cognito PreAuthentication trigger.
// It blocks banned or suspended users from authenticating.
type PreAuthenticationUseCase struct {
	users UserFinder
}

// NewPreAuthenticationUseCase creates a PreAuthenticationUseCase with its dependencies.
func NewPreAuthenticationUseCase(users UserFinder) *PreAuthenticationUseCase {
	return &PreAuthenticationUseCase{users: users}
}

// Execute checks whether the authenticating user is banned or suspended.
// If the user is blocked, it returns an error to deny authentication.
// On infrastructure errors, authentication proceeds to avoid blocking
// users due to transient failures.
func (uc *PreAuthenticationUseCase) Execute(ctx context.Context, event *events.CognitoEventUserPoolsPreAuthentication) (events.CognitoEventUserPoolsPreAuthentication, error) {
	cognitoSub := handler.CognitoUserName(&event.CognitoEventUserPoolsHeader)
	_ = handler.PreAuthUserAttributes(event)

	slog.Default().InfoContext(ctx, "pre-authentication trigger",
		slog.String("cognitoSub", cognitoSub),
	)

	u, err := uc.users.FindByID(ctx, cognitoSub)
	if err != nil {
		slog.Default().WarnContext(ctx, "failed to lookup user for pre-auth check, allowing auth",
			slog.String("cognitoSub", cognitoSub),
			slog.Any("error", err),
		)
		return *event, nil
	}

	if u == nil {
		slog.Default().InfoContext(ctx, "user not found in pre-auth, allowing auth for new user",
			slog.String("cognitoSub", cognitoSub),
		)
		return *event, nil
	}

	if u.Status == user.UserStatusBanned {
		slog.Default().WarnContext(ctx, "blocked banned user from authenticating",
			slog.String("cognitoSub", cognitoSub),
		)
		return *event, fmt.Errorf("user account is banned")
	}

	if u.Status == user.UserStatusSuspended {
		slog.Default().WarnContext(ctx, "blocked suspended user from authenticating",
			slog.String("cognitoSub", cognitoSub),
		)
		return *event, fmt.Errorf("user account is suspended")
	}

	slog.Default().InfoContext(ctx, "pre-auth check passed",
		slog.String("cognitoSub", cognitoSub),
		slog.String("status", string(u.Status)),
	)

	return *event, nil
}
