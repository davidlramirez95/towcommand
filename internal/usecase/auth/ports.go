// Package auth implements Cognito trigger use cases following CLEAN architecture.
// Each use case declares only the port interfaces it needs (ISP).
package auth

import (
	"context"

	"github.com/davidlramirez95/towcommand/internal/domain/provider"
	"github.com/davidlramirez95/towcommand/internal/domain/user"
	"github.com/davidlramirez95/towcommand/internal/usecase/port"
)

// UserSaver persists a new user created during post-confirmation.
type UserSaver interface {
	Save(ctx context.Context, u *user.User) error
}

// UserFinder retrieves a user by their ID (CognitoSub).
type UserFinder interface {
	FindByID(ctx context.Context, userID string) (*user.User, error)
}

// ProviderByCognitoSubFinder finds a provider by their Cognito sub (user identity).
type ProviderByCognitoSubFinder interface {
	FindByCognitoSub(ctx context.Context, cognitoSub string) (*provider.Provider, error)
}

// EventPublisher publishes domain events to an event bus.
type EventPublisher interface {
	Publish(ctx context.Context, source, detailType string, detail any, actor *port.Actor) error
}

// Actor is a type alias for port.Actor to avoid import stuttering in use case code.
type Actor = port.Actor
