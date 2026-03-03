// Package port defines interfaces (ports) that adapters must implement.
package port

import "context"

// Actor identifies who triggered an event.
type Actor struct {
	UserID   string
	UserType string
}

// EventPublisher publishes domain events to an event bus.
type EventPublisher interface {
	Publish(ctx context.Context, source, detailType string, detail any, actor *Actor) error
}
