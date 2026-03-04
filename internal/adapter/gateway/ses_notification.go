package gateway

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ses"
	sestypes "github.com/aws/aws-sdk-go-v2/service/ses/types"

	"github.com/davidlramirez95/towcommand/internal/domain/errors"
	"github.com/davidlramirez95/towcommand/internal/usecase/port"
)

const sesSourceEmail = "noreply@towcommand.ph"

// SESAPI is the subset of the SES client needed by the email sender.
type SESAPI interface {
	SendEmail(ctx context.Context, params *ses.SendEmailInput, optFns ...func(*ses.Options)) (*ses.SendEmailOutput, error)
}

// SESNotificationSender sends email messages via AWS SES.
// It implements the port.EmailSender interface.
type SESNotificationSender struct {
	client SESAPI
}

// NewSESNotificationSender creates a new SES-backed email sender.
func NewSESNotificationSender(client SESAPI) *SESNotificationSender {
	return &SESNotificationSender{client: client}
}

// Compile-time interface check.
var _ port.EmailSender = (*SESNotificationSender)(nil)

// SendEmail sends an HTML email via SES from the TowCommand no-reply address.
func (s *SESNotificationSender) SendEmail(ctx context.Context, to, subject, htmlBody string) error {
	_, err := s.client.SendEmail(ctx, &ses.SendEmailInput{
		Source: aws.String(sesSourceEmail),
		Destination: &sestypes.Destination{
			ToAddresses: []string{to},
		},
		Message: &sestypes.Message{
			Subject: &sestypes.Content{
				Data:    aws.String(subject),
				Charset: aws.String("UTF-8"),
			},
			Body: &sestypes.Body{
				Html: &sestypes.Content{
					Data:    aws.String(htmlBody),
					Charset: aws.String("UTF-8"),
				},
			},
		},
	})
	if err != nil {
		return errors.NewExternalServiceError("SES", fmt.Errorf("sending email to %s: %w", to, err))
	}
	return nil
}
