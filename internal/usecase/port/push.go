package port

import (
	"context"
	"time"
)

// PushPlatform identifies the mobile push platform.
type PushPlatform string

const (
	// PushPlatformFCM represents Firebase Cloud Messaging (Android).
	PushPlatformFCM PushPlatform = "FCM"
	// PushPlatformAPNS represents Apple Push Notification Service (iOS).
	PushPlatformAPNS PushPlatform = "APNS"
)

// PushToken represents a registered push notification token for a device.
type PushToken struct {
	UserID      string       `json:"userId"`
	Token       string       `json:"token"`
	Platform    PushPlatform `json:"platform"`
	DeviceID    string       `json:"deviceId"`
	EndpointArn string       `json:"endpointArn"`
	CreatedAt   time.Time    `json:"createdAt"`
	UpdatedAt   time.Time    `json:"updatedAt"`
}

// PushTokenRegistrar manages push notification token storage.
type PushTokenRegistrar interface {
	Register(ctx context.Context, token *PushToken) error
	FindByUserID(ctx context.Context, userID string) ([]PushToken, error)
	Delete(ctx context.Context, userID, deviceID string) error
}

// PushTokenFinder retrieves push tokens for a user.
type PushTokenFinder interface {
	FindByUserID(ctx context.Context, userID string) ([]PushToken, error)
}

// PushSender sends push notifications to mobile devices.
type PushSender interface {
	SendPush(ctx context.Context, endpointArn, title, message string, data map[string]string) error
}

// PushEndpointCreator creates SNS platform endpoints for push tokens.
type PushEndpointCreator interface {
	CreateEndpoint(ctx context.Context, platform PushPlatform, token string) (endpointArn string, err error)
}
