package gateway

import (
	"context"
	"fmt"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/sns"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/davidlramirez95/towcommand/internal/usecase/port"
)

// mockSNSClient implements SNSAPI for testing.
type mockSNSClient struct {
	publishFunc func(ctx context.Context, params *sns.PublishInput, optFns ...func(*sns.Options)) (*sns.PublishOutput, error)
}

func (m *mockSNSClient) Publish(ctx context.Context, params *sns.PublishInput, optFns ...func(*sns.Options)) (*sns.PublishOutput, error) {
	return m.publishFunc(ctx, params, optFns...)
}

func TestSNSNotificationSender_SendSMS(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name           string
		phoneNumber    string
		message        string
		mockOutput     *sns.PublishOutput
		mockErr        error
		wantErr        bool
		wantErrContain string
		validate       func(t *testing.T, input *sns.PublishInput)
	}{
		{
			name:        "successful SMS send",
			phoneNumber: "+639171234567",
			message:     "Your TowCommand verification code: 123456",
			mockOutput:  &sns.PublishOutput{MessageId: strPtr("msg-001")},
			validate: func(t *testing.T, input *sns.PublishInput) {
				t.Helper()
				assert.Equal(t, "+639171234567", *input.PhoneNumber)
				assert.Equal(t, "Your TowCommand verification code: 123456", *input.Message)

				smsType, ok := input.MessageAttributes["AWS.SNS.SMS.SMSType"]
				require.True(t, ok)
				assert.Equal(t, "String", *smsType.DataType)
				assert.Equal(t, "Transactional", *smsType.StringValue)
			},
		},
		{
			name:           "SNS API error returns external service error",
			phoneNumber:    "+639171234567",
			message:        "code: 654321",
			mockErr:        fmt.Errorf("throttling exception"),
			wantErr:        true,
			wantErrContain: "SNS",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			var capturedInput *sns.PublishInput
			mock := &mockSNSClient{
				publishFunc: func(_ context.Context, params *sns.PublishInput, _ ...func(*sns.Options)) (*sns.PublishOutput, error) {
					capturedInput = params
					if tt.mockErr != nil {
						return nil, tt.mockErr
					}
					return tt.mockOutput, nil
				},
			}

			sender := NewSNSNotificationSender(mock)
			err := sender.SendSMS(context.Background(), tt.phoneNumber, tt.message)

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

func TestSNSNotificationSender_ImplementsPort(t *testing.T) {
	t.Parallel()
	var _ port.SMSSender = (*SNSNotificationSender)(nil)
}
