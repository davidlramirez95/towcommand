package auth

import (
	"context"
	"log/slog"
	"strings"

	"github.com/aws/aws-lambda-go/events"

	"github.com/davidlramirez95/towcommand/internal/adapter/handler"
)

// PreSignUpUseCase handles the Cognito PreSignUp trigger.
// It auto-confirms social logins and dev/local stage accounts.
type PreSignUpUseCase struct {
	stage string
}

// NewPreSignUpUseCase creates a PreSignUpUseCase for the given deployment stage.
func NewPreSignUpUseCase(stage string) *PreSignUpUseCase {
	return &PreSignUpUseCase{stage: stage}
}

// Execute processes the PreSignUp event. It auto-confirms users from external
// identity providers (social login) and accounts in dev/local stages for
// easier testing. Normal sign-ups are returned unchanged.
func (uc *PreSignUpUseCase) Execute(ctx context.Context, event *events.CognitoEventUserPoolsPreSignup) (events.CognitoEventUserPoolsPreSignup, error) {
	triggerSource := handler.CognitoTriggerSource(&event.CognitoEventUserPoolsHeader)
	userName := handler.CognitoUserName(&event.CognitoEventUserPoolsHeader)

	slog.Default().InfoContext(ctx, "pre-signup trigger",
		slog.String("triggerSource", triggerSource),
		slog.String("userName", userName),
		slog.String("stage", uc.stage),
	)

	// Social login: auto-confirm and auto-verify email.
	if strings.Contains(triggerSource, "PreSignUp_ExternalProvider") {
		slog.Default().InfoContext(ctx, "auto-confirming social login user",
			slog.String("userName", userName),
		)
		return handler.AutoConfirmUser(event, true), nil
	}

	// Dev/local stage: auto-confirm for easier testing.
	if uc.stage == "dev" || uc.stage == "local" {
		slog.Default().InfoContext(ctx, "auto-confirming dev/local user",
			slog.String("userName", userName),
			slog.String("stage", uc.stage),
		)
		return handler.AutoConfirmUser(event, true), nil
	}

	// Normal flow: return event unchanged, Cognito handles verification.
	return *event, nil
}
