package gateway

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sns"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/davidlramirez95/towcommand/internal/usecase/port"
)

// mockSNSPushClient implements SNSPushAPI for testing.
type mockSNSPushClient struct {
	publishFunc              func(ctx context.Context, params *sns.PublishInput, optFns ...func(*sns.Options)) (*sns.PublishOutput, error)
	createPlatformEndpointFn func(ctx context.Context, params *sns.CreatePlatformEndpointInput, optFns ...func(*sns.Options)) (*sns.CreatePlatformEndpointOutput, error)
}

func (m *mockSNSPushClient) Publish(ctx context.Context, params *sns.PublishInput, optFns ...func(*sns.Options)) (*sns.PublishOutput, error) {
	return m.publishFunc(ctx, params, optFns...)
}

func (m *mockSNSPushClient) CreatePlatformEndpoint(ctx context.Context, params *sns.CreatePlatformEndpointInput, optFns ...func(*sns.Options)) (*sns.CreatePlatformEndpointOutput, error) {
	return m.createPlatformEndpointFn(ctx, params, optFns...)
}

func TestSNSPushSender_SendPush(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name           string
		endpointArn    string
		title          string
		message        string
		data           map[string]string
		mockOutput     *sns.PublishOutput
		mockErr        error
		wantErr        bool
		wantErrContain string
		validate       func(t *testing.T, input *sns.PublishInput)
	}{
		{
			name:        "successful push send",
			endpointArn: "arn:aws:sns:ap-southeast-1:123:endpoint/GCM/App/abc",
			title:       "Booking Update",
			message:     "Your driver is on the way!",
			data:        map[string]string{"bookingId": "bk-1"},
			mockOutput:  &sns.PublishOutput{MessageId: aws.String("msg-001")},
			validate: func(t *testing.T, input *sns.PublishInput) {
				t.Helper()
				assert.Equal(t, "arn:aws:sns:ap-southeast-1:123:endpoint/GCM/App/abc", *input.TargetArn)
				assert.Equal(t, "json", *input.MessageStructure)

				// Parse the SNS message and verify it has GCM and APNS keys.
				var snsMsg map[string]string
				require.NoError(t, json.Unmarshal([]byte(*input.Message), &snsMsg))
				assert.Contains(t, snsMsg, "GCM")
				assert.Contains(t, snsMsg, "APNS")

				// Verify GCM payload structure.
				var gcm gcmPayload
				require.NoError(t, json.Unmarshal([]byte(snsMsg["GCM"]), &gcm))
				assert.Equal(t, "Booking Update", gcm.Notification.Title)
				assert.Equal(t, "Your driver is on the way!", gcm.Notification.Body)
				assert.Equal(t, "bk-1", gcm.Data["bookingId"])

				// Verify APNS payload structure.
				var apns apnsPayload
				require.NoError(t, json.Unmarshal([]byte(snsMsg["APNS"]), &apns))
				assert.Equal(t, "Booking Update", apns.APS.Alert.Title)
				assert.Equal(t, "Your driver is on the way!", apns.APS.Alert.Body)
				assert.Equal(t, "bk-1", apns.Data["bookingId"])
			},
		},
		{
			name:        "successful push send with nil data",
			endpointArn: "arn:aws:sns:ap-southeast-1:123:endpoint/GCM/App/abc",
			title:       "Welcome",
			message:     "Welcome to TowCommand!",
			data:        nil,
			mockOutput:  &sns.PublishOutput{MessageId: aws.String("msg-002")},
			validate: func(t *testing.T, input *sns.PublishInput) {
				t.Helper()
				var snsMsg map[string]string
				require.NoError(t, json.Unmarshal([]byte(*input.Message), &snsMsg))

				var gcm gcmPayload
				require.NoError(t, json.Unmarshal([]byte(snsMsg["GCM"]), &gcm))
				assert.Equal(t, "Welcome", gcm.Notification.Title)
				assert.Nil(t, gcm.Data)
			},
		},
		{
			name:           "SNS API error returns external service error",
			endpointArn:    "arn:aws:sns:ap-southeast-1:123:endpoint/GCM/App/abc",
			title:          "Test",
			message:        "Test message",
			mockErr:        fmt.Errorf("endpoint disabled"),
			wantErr:        true,
			wantErrContain: "SNS",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			var capturedInput *sns.PublishInput
			mock := &mockSNSPushClient{
				publishFunc: func(_ context.Context, params *sns.PublishInput, _ ...func(*sns.Options)) (*sns.PublishOutput, error) {
					capturedInput = params
					if tt.mockErr != nil {
						return nil, tt.mockErr
					}
					return tt.mockOutput, nil
				},
				createPlatformEndpointFn: func(_ context.Context, _ *sns.CreatePlatformEndpointInput, _ ...func(*sns.Options)) (*sns.CreatePlatformEndpointOutput, error) {
					return nil, nil
				},
			}

			sender := NewSNSPushSenderWithARNs(mock, "arn:fcm", "arn:apns")
			err := sender.SendPush(context.Background(), tt.endpointArn, tt.title, tt.message, tt.data)

			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.wantErrContain)
				return
			}

			require.NoError(t, err)
			if tt.validate != nil {
				tt.validate(t, capturedInput)
			}
		})
	}
}

func TestSNSPushSender_CreateEndpoint(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name           string
		platform       port.PushPlatform
		token          string
		fcmArn         string
		apnsArn        string
		mockOutput     *sns.CreatePlatformEndpointOutput
		mockErr        error
		wantArn        string
		wantErr        bool
		wantErrContain string
		validate       func(t *testing.T, input *sns.CreatePlatformEndpointInput)
	}{
		{
			name:     "FCM endpoint creation",
			platform: port.PushPlatformFCM,
			token:    "fcm-device-token-123",
			fcmArn:   "arn:aws:sns:ap-southeast-1:123:app/GCM/TowCommand",
			apnsArn:  "arn:aws:sns:ap-southeast-1:123:app/APNS/TowCommand",
			mockOutput: &sns.CreatePlatformEndpointOutput{
				EndpointArn: aws.String("arn:aws:sns:ap-southeast-1:123:endpoint/GCM/TowCommand/new-endpoint"),
			},
			wantArn: "arn:aws:sns:ap-southeast-1:123:endpoint/GCM/TowCommand/new-endpoint",
			validate: func(t *testing.T, input *sns.CreatePlatformEndpointInput) {
				t.Helper()
				assert.Equal(t, "arn:aws:sns:ap-southeast-1:123:app/GCM/TowCommand", *input.PlatformApplicationArn)
				assert.Equal(t, "fcm-device-token-123", *input.Token)
			},
		},
		{
			name:     "APNS endpoint creation",
			platform: port.PushPlatformAPNS,
			token:    "apns-device-token-456",
			fcmArn:   "arn:aws:sns:ap-southeast-1:123:app/GCM/TowCommand",
			apnsArn:  "arn:aws:sns:ap-southeast-1:123:app/APNS/TowCommand",
			mockOutput: &sns.CreatePlatformEndpointOutput{
				EndpointArn: aws.String("arn:aws:sns:ap-southeast-1:123:endpoint/APNS/TowCommand/new-endpoint"),
			},
			wantArn: "arn:aws:sns:ap-southeast-1:123:endpoint/APNS/TowCommand/new-endpoint",
			validate: func(t *testing.T, input *sns.CreatePlatformEndpointInput) {
				t.Helper()
				assert.Equal(t, "arn:aws:sns:ap-southeast-1:123:app/APNS/TowCommand", *input.PlatformApplicationArn)
				assert.Equal(t, "apns-device-token-456", *input.Token)
			},
		},
		{
			name:           "SNS API error",
			platform:       port.PushPlatformFCM,
			token:          "bad-token",
			fcmArn:         "arn:aws:sns:ap-southeast-1:123:app/GCM/TowCommand",
			mockErr:        fmt.Errorf("invalid token"),
			wantErr:        true,
			wantErrContain: "SNS",
		},
		{
			name:           "missing platform ARN",
			platform:       port.PushPlatformFCM,
			token:          "some-token",
			fcmArn:         "",
			apnsArn:        "",
			wantErr:        true,
			wantErrContain: "not configured",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			var capturedInput *sns.CreatePlatformEndpointInput
			mock := &mockSNSPushClient{
				publishFunc: func(_ context.Context, _ *sns.PublishInput, _ ...func(*sns.Options)) (*sns.PublishOutput, error) {
					return nil, nil
				},
				createPlatformEndpointFn: func(_ context.Context, params *sns.CreatePlatformEndpointInput, _ ...func(*sns.Options)) (*sns.CreatePlatformEndpointOutput, error) {
					capturedInput = params
					if tt.mockErr != nil {
						return nil, tt.mockErr
					}
					return tt.mockOutput, nil
				},
			}

			sender := NewSNSPushSenderWithARNs(mock, tt.fcmArn, tt.apnsArn)
			arn, err := sender.CreateEndpoint(context.Background(), tt.platform, tt.token)

			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.wantErrContain)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tt.wantArn, arn)
			if tt.validate != nil {
				tt.validate(t, capturedInput)
			}
		})
	}
}

func TestSNSPushSender_ImplementsPorts(t *testing.T) {
	t.Parallel()
	var _ port.PushSender = (*SNSPushSender)(nil)
	var _ port.PushEndpointCreator = (*SNSPushSender)(nil)
}
