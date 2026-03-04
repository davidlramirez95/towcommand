// Package websocket implements WebSocket use cases following CLEAN architecture.
// Each use case declares only the port interfaces it needs (ISP).
package websocket

import (
	"context"
	"time"

	"github.com/davidlramirez95/towcommand/internal/domain/booking"
	"github.com/davidlramirez95/towcommand/internal/usecase/port"
)

// ConnectionMapper stores forward and reverse WebSocket connection mappings.
type ConnectionMapper interface {
	MapConnection(ctx context.Context, userID, connectionID string, ttl time.Duration) error
	MapReverseConnection(ctx context.Context, connectionID, userID string, ttl time.Duration) error
}

// ConnectionRemover removes forward and reverse WebSocket connection mappings.
type ConnectionRemover interface {
	GetUserByConnection(ctx context.Context, connectionID string) (string, error)
	RemoveConnection(ctx context.Context, userID string) error
	RemoveReverseConnection(ctx context.Context, connectionID string) error
}

// ConnectionLookup resolves a user ID to a WebSocket connection ID.
type ConnectionLookup interface {
	GetConnection(ctx context.Context, userID string) (string, error)
}

// GeoUpdater updates a provider's geospatial position.
type GeoUpdater interface {
	AddProviderLocation(ctx context.Context, providerID string, lat, lng float64) error
}

// EventPublisher publishes domain events to an event bus.
type EventPublisher interface {
	Publish(ctx context.Context, source, detailType string, detail any, actor *port.Actor) error
}

// ChatSaver persists a chat message.
type ChatSaver interface {
	Save(ctx context.Context, msg *booking.ChatMessage) error
}

// ConnectionPoster sends data to a WebSocket connection via API Gateway Management API.
type ConnectionPoster interface {
	PostToConnection(ctx context.Context, connectionID string, data any) error
}

// Actor is a type alias for port.Actor to avoid import stuttering in use case code.
type Actor = port.Actor
