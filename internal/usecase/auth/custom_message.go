package auth

import (
	"context"
	"log/slog"

	"github.com/aws/aws-lambda-go/events"

	"github.com/davidlramirez95/towcommand/internal/adapter/handler"
)

// CustomMessageUseCase handles the Cognito CustomMessage trigger.
// It provides Filipino/English message templates for verification codes,
// password resets, and code resends.
type CustomMessageUseCase struct{}

// NewCustomMessageUseCase creates a CustomMessageUseCase.
func NewCustomMessageUseCase() *CustomMessageUseCase {
	return &CustomMessageUseCase{}
}

// Execute customizes the email and SMS messages for the given trigger source.
func (uc *CustomMessageUseCase) Execute(ctx context.Context, event *events.CognitoEventUserPoolsCustomMessage) (events.CognitoEventUserPoolsCustomMessage, error) {
	triggerSource := handler.CognitoTriggerSource(&event.CognitoEventUserPoolsHeader)
	userName := handler.CognitoUserName(&event.CognitoEventUserPoolsHeader)

	slog.Default().InfoContext(ctx, "custom-message trigger",
		slog.String("triggerSource", triggerSource),
		slog.String("userName", userName),
	)

	var emailSubject, emailMessage, smsMessage string

	switch triggerSource {
	case "CustomMessage_SignUp":
		emailSubject = "Mabuhay! Welcome to TowCommand - Verify Your Account"
		emailMessage = "Mabuhay! Welcome to TowCommand. Your verification code is {####}. " +
			"Please enter this code to verify your account."
		smsMessage = "Mabuhay! TowCommand verification code: {####}"

	case "CustomMessage_ForgotPassword":
		emailSubject = "TowCommand - Password Reset Code"
		emailMessage = "Your TowCommand password reset code is {####}. " +
			"If you did not request this, please ignore this message."
		smsMessage = "TowCommand password reset code: {####}"

	case "CustomMessage_ResendCode":
		emailSubject = "TowCommand - Verification Code"
		emailMessage = "Your TowCommand verification code is {####}. " +
			"Please enter this code to verify your account."
		smsMessage = "TowCommand verification code: {####}"

	case "CustomMessage_AdminCreateUser":
		emailSubject = "Mabuhay! Welcome to TowCommand"
		emailMessage = "Mabuhay! Your TowCommand account has been created. " +
			"Your temporary password is {####}. Please log in and change it."
		smsMessage = "Mabuhay! TowCommand temporary password: {####}"

	default:
		slog.Default().InfoContext(ctx, "no custom template for trigger source, using defaults",
			slog.String("triggerSource", triggerSource),
		)
		return *event, nil
	}

	return handler.CustomizeMessage(event, emailSubject, emailMessage, smsMessage), nil
}
