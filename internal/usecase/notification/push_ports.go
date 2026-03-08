package notification

import (
	"context"

	"github.com/davidlramirez95/towcommand/internal/usecase/port"
)

// PushSender sends push notifications to mobile devices.
type PushSender interface {
	SendPush(ctx context.Context, endpointArn, title, message string, data map[string]string) error
}

// PushTokenFinder retrieves push tokens for a user.
type PushTokenFinder interface {
	FindByUserID(ctx context.Context, userID string) ([]port.PushToken, error)
}
