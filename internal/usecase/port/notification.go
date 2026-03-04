package port

import "context"

// SMSSender sends SMS messages to phone numbers.
type SMSSender interface {
	SendSMS(ctx context.Context, phoneNumber, message string) error
}

// EmailSender sends email messages.
type EmailSender interface {
	SendEmail(ctx context.Context, to, subject, htmlBody string) error
}
