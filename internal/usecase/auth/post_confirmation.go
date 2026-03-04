package auth

import (
	"context"
	"log/slog"
	"time"

	"github.com/aws/aws-lambda-go/events"

	"github.com/davidlramirez95/towcommand/internal/adapter/handler"
	domevent "github.com/davidlramirez95/towcommand/internal/domain/event"
	"github.com/davidlramirez95/towcommand/internal/domain/user"
)

// PostConfirmationUseCase handles the Cognito PostConfirmation trigger.
// It creates a user record in DynamoDB and publishes a UserRegistered event.
type PostConfirmationUseCase struct {
	saver     UserSaver
	publisher EventPublisher
	now       func() time.Time
}

// NewPostConfirmationUseCase creates a PostConfirmationUseCase with its dependencies.
func NewPostConfirmationUseCase(saver UserSaver, publisher EventPublisher) *PostConfirmationUseCase {
	return &PostConfirmationUseCase{
		saver:     saver,
		publisher: publisher,
		now:       func() time.Time { return time.Now().UTC() },
	}
}

// Execute creates a new user from the confirmed Cognito account.
// Errors are logged but never returned — Cognito blocks signup if a
// PostConfirmation trigger returns an error.
func (uc *PostConfirmationUseCase) Execute(ctx context.Context, event *events.CognitoEventUserPoolsPostConfirmation) (events.CognitoEventUserPoolsPostConfirmation, error) {
	cognitoSub := handler.CognitoUserName(&event.CognitoEventUserPoolsHeader)
	attrs := handler.PostConfirmationUserAttributes(event)

	slog.Default().InfoContext(ctx, "post-confirmation trigger",
		slog.String("cognitoSub", cognitoSub),
		slog.String("email", attrs["email"]),
	)

	now := uc.now()
	u := &user.User{
		UserID:     cognitoSub,
		CognitoSub: cognitoSub,
		Email:      attrs["email"],
		Phone:      attrs["phone_number"],
		Name:       attrs["name"],
		UserType:   user.UserTypeCustomer,
		TrustTier:  user.TrustTierBasic,
		Language:   user.LanguageEnglish,
		Status:     user.UserStatusActive,
		CreatedAt:  now,
		UpdatedAt:  now,
	}

	if err := uc.saver.Save(ctx, u); err != nil {
		slog.Default().ErrorContext(ctx, "failed to save user in post-confirmation",
			slog.String("cognitoSub", cognitoSub),
			slog.Any("error", err),
		)
		// Return event without error to avoid blocking Cognito signup flow.
		return *event, nil
	}

	if err := uc.publisher.Publish(ctx, domevent.SourceAuth, domevent.UserRegistered, map[string]any{
		"userId":   cognitoSub,
		"email":    attrs["email"],
		"phone":    attrs["phone_number"],
		"name":     attrs["name"],
		"userType": string(user.UserTypeCustomer),
	}, &Actor{UserID: cognitoSub, UserType: string(user.UserTypeCustomer)}); err != nil {
		slog.Default().WarnContext(ctx, "failed to publish UserRegistered event",
			slog.String("cognitoSub", cognitoSub),
			slog.Any("error", err),
		)
	}

	slog.Default().InfoContext(ctx, "user created from post-confirmation",
		slog.String("userId", cognitoSub),
		slog.String("email", attrs["email"]),
	)

	return *event, nil
}
