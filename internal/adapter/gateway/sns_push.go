package gateway

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sns"

	"github.com/davidlramirez95/towcommand/internal/domain/errors"
	"github.com/davidlramirez95/towcommand/internal/usecase/port"
)

// SNSPushAPI is the subset of the SNS client needed by the push sender.
type SNSPushAPI interface {
	Publish(ctx context.Context, params *sns.PublishInput, optFns ...func(*sns.Options)) (*sns.PublishOutput, error)
	CreatePlatformEndpoint(ctx context.Context, params *sns.CreatePlatformEndpointInput, optFns ...func(*sns.Options)) (*sns.CreatePlatformEndpointOutput, error)
}

// SNSPushSender sends push notifications via AWS SNS Platform Application.
// It implements port.PushSender and port.PushEndpointCreator.
type SNSPushSender struct {
	client      SNSPushAPI
	fcmPlatArn  string
	apnsPlatArn string
}

// NewSNSPushSender creates a new SNS-backed push notification sender.
// Platform application ARNs are read from environment variables:
// SNS_PLATFORM_APP_ARN_FCM and SNS_PLATFORM_APP_ARN_APNS.
func NewSNSPushSender(client SNSPushAPI) *SNSPushSender {
	return &SNSPushSender{
		client:      client,
		fcmPlatArn:  os.Getenv("SNS_PLATFORM_APP_ARN_FCM"),
		apnsPlatArn: os.Getenv("SNS_PLATFORM_APP_ARN_APNS"),
	}
}

// NewSNSPushSenderWithARNs creates a new SNS-backed push sender with explicit
// platform application ARNs. This is useful for testing and environments
// where env vars are not available.
func NewSNSPushSenderWithARNs(client SNSPushAPI, fcmArn, apnsArn string) *SNSPushSender {
	return &SNSPushSender{
		client:      client,
		fcmPlatArn:  fcmArn,
		apnsPlatArn: apnsArn,
	}
}

// Compile-time interface checks.
var (
	_ port.PushSender          = (*SNSPushSender)(nil)
	_ port.PushEndpointCreator = (*SNSPushSender)(nil)
)

// gcmPayload is the FCM message structure published via SNS.
type gcmPayload struct {
	Notification gcmNotification   `json:"notification"`
	Data         map[string]string `json:"data,omitempty"`
}

// gcmNotification holds the title and body for an FCM notification.
type gcmNotification struct {
	Title string `json:"title"`
	Body  string `json:"body"`
}

// apnsPayload is the APNs message structure published via SNS.
type apnsPayload struct {
	APS  apnsAPS           `json:"aps"`
	Data map[string]string `json:"data,omitempty"`
}

// apnsAPS holds the alert for an APNs notification.
type apnsAPS struct {
	Alert apnsAlert `json:"alert"`
}

// apnsAlert holds the title and body for an APNs alert.
type apnsAlert struct {
	Title string `json:"title"`
	Body  string `json:"body"`
}

// SendPush publishes a push notification to the given SNS platform endpoint.
// The message is formatted for both FCM and APNs using the SNS message structure.
func (s *SNSPushSender) SendPush(ctx context.Context, endpointArn, title, message string, data map[string]string) error {
	gcm := gcmPayload{
		Notification: gcmNotification{Title: title, Body: message},
		Data:         data,
	}
	gcmJSON, err := json.Marshal(gcm)
	if err != nil {
		return fmt.Errorf("marshalling GCM payload: %w", err)
	}

	apns := apnsPayload{
		APS:  apnsAPS{Alert: apnsAlert{Title: title, Body: message}},
		Data: data,
	}
	apnsJSON, err := json.Marshal(apns)
	if err != nil {
		return fmt.Errorf("marshalling APNS payload: %w", err)
	}

	snsMsg := map[string]string{
		"GCM":  string(gcmJSON),
		"APNS": string(apnsJSON),
	}
	snsMsgJSON, err := json.Marshal(snsMsg)
	if err != nil {
		return fmt.Errorf("marshalling SNS message: %w", err)
	}

	_, err = s.client.Publish(ctx, &sns.PublishInput{
		TargetArn:        aws.String(endpointArn),
		Message:          aws.String(string(snsMsgJSON)),
		MessageStructure: aws.String("json"),
	})
	if err != nil {
		return errors.NewExternalServiceError("SNS", fmt.Errorf("sending push to %s: %w", endpointArn, err))
	}
	return nil
}

// CreateEndpoint creates an SNS platform endpoint for the given device token.
// It returns the endpoint ARN that can be used to send push notifications.
func (s *SNSPushSender) CreateEndpoint(ctx context.Context, platform port.PushPlatform, token string) (string, error) {
	platArn := s.platformARN(platform)
	if platArn == "" {
		return "", errors.NewValidationError(
			fmt.Sprintf("platform application ARN not configured for %s", platform),
		)
	}

	out, err := s.client.CreatePlatformEndpoint(ctx, &sns.CreatePlatformEndpointInput{
		PlatformApplicationArn: aws.String(platArn),
		Token:                  aws.String(token),
	})
	if err != nil {
		return "", errors.NewExternalServiceError("SNS", fmt.Errorf("creating platform endpoint for %s: %w", platform, err))
	}
	return *out.EndpointArn, nil
}

// platformARN returns the SNS platform application ARN for the given platform.
func (s *SNSPushSender) platformARN(platform port.PushPlatform) string {
	switch platform {
	case port.PushPlatformFCM:
		return s.fcmPlatArn
	case port.PushPlatformAPNS:
		return s.apnsPlatArn
	default:
		return ""
	}
}
