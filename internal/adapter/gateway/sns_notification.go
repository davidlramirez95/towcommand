package gateway

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sns"
	snstypes "github.com/aws/aws-sdk-go-v2/service/sns/types"

	"github.com/davidlramirez95/towcommand/internal/domain/errors"
	"github.com/davidlramirez95/towcommand/internal/usecase/port"
)

// SNSAPI is the subset of the SNS client needed by the SMS sender.
type SNSAPI interface {
	Publish(ctx context.Context, params *sns.PublishInput, optFns ...func(*sns.Options)) (*sns.PublishOutput, error)
}

// SNSNotificationSender sends SMS messages via AWS SNS.
// It implements the port.SMSSender interface.
type SNSNotificationSender struct {
	client SNSAPI
}

// NewSNSNotificationSender creates a new SNS-backed SMS sender.
func NewSNSNotificationSender(client SNSAPI) *SNSNotificationSender {
	return &SNSNotificationSender{client: client}
}

// Compile-time interface check.
var _ port.SMSSender = (*SNSNotificationSender)(nil)

// SendSMS publishes an SMS message to the given phone number via SNS.
// The message is sent as a Transactional SMS to ensure high deliverability.
func (s *SNSNotificationSender) SendSMS(ctx context.Context, phoneNumber, message string) error {
	_, err := s.client.Publish(ctx, &sns.PublishInput{
		PhoneNumber: aws.String(phoneNumber),
		Message:     aws.String(message),
		MessageAttributes: map[string]snstypes.MessageAttributeValue{
			"AWS.SNS.SMS.SMSType": {
				DataType:    aws.String("String"),
				StringValue: aws.String("Transactional"),
			},
		},
	})
	if err != nil {
		return errors.NewExternalServiceError("SNS", fmt.Errorf("sending SMS to %s: %w", phoneNumber, err))
	}
	return nil
}
