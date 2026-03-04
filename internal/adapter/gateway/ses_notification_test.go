package gateway

import (
	"context"
	"fmt"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/ses"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/davidlramirez95/towcommand/internal/usecase/port"
)

// mockSESClient implements SESAPI for testing.
type mockSESClient struct {
	sendEmailFunc func(ctx context.Context, params *ses.SendEmailInput, optFns ...func(*ses.Options)) (*ses.SendEmailOutput, error)
}

func (m *mockSESClient) SendEmail(ctx context.Context, params *ses.SendEmailInput, optFns ...func(*ses.Options)) (*ses.SendEmailOutput, error) {
	return m.sendEmailFunc(ctx, params, optFns...)
}

func TestSESNotificationSender_SendEmail(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name           string
		to             string
		subject        string
		htmlBody       string
		mockOutput     *ses.SendEmailOutput
		mockErr        error
		wantErr        bool
		wantErrContain string
		validate       func(t *testing.T, input *ses.SendEmailInput)
	}{
		{
			name:       "successful email send",
			to:         "customer@example.com",
			subject:    "Your OTP Code",
			htmlBody:   "<h1>123456</h1>",
			mockOutput: &ses.SendEmailOutput{MessageId: strPtr("ses-msg-001")},
			validate: func(t *testing.T, input *ses.SendEmailInput) {
				t.Helper()
				assert.Equal(t, "noreply@towcommand.ph", *input.Source)
				require.Len(t, input.Destination.ToAddresses, 1)
				assert.Equal(t, "customer@example.com", input.Destination.ToAddresses[0])
				assert.Equal(t, "Your OTP Code", *input.Message.Subject.Data)
				assert.Equal(t, "UTF-8", *input.Message.Subject.Charset)
				assert.Equal(t, "<h1>123456</h1>", *input.Message.Body.Html.Data)
				assert.Equal(t, "UTF-8", *input.Message.Body.Html.Charset)
			},
		},
		{
			name:           "SES API error returns external service error",
			to:             "customer@example.com",
			subject:        "OTP",
			htmlBody:       "<p>code</p>",
			mockErr:        fmt.Errorf("message rejected"),
			wantErr:        true,
			wantErrContain: "SES",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			var capturedInput *ses.SendEmailInput
			mock := &mockSESClient{
				sendEmailFunc: func(_ context.Context, params *ses.SendEmailInput, _ ...func(*ses.Options)) (*ses.SendEmailOutput, error) {
					capturedInput = params
					if tt.mockErr != nil {
						return nil, tt.mockErr
					}
					return tt.mockOutput, nil
				},
			}

			sender := NewSESNotificationSender(mock)
			err := sender.SendEmail(context.Background(), tt.to, tt.subject, tt.htmlBody)

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

func TestSESNotificationSender_ImplementsPort(t *testing.T) {
	t.Parallel()
	var _ port.EmailSender = (*SESNotificationSender)(nil)
}
