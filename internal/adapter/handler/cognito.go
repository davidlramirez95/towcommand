package handler

import (
	"github.com/aws/aws-lambda-go/events"
)

// Cognito trigger source constants.
const (
	TriggerPreSignUp          = "PreSignUp_SignUp"
	TriggerPostConfirmation   = "PostConfirmation_ConfirmSignUp"
	TriggerPreTokenGeneration = "TokenGeneration_HostedAuth"
	TriggerPreAuthentication  = "PreAuthentication_Authentication"
	TriggerCustomMessage      = "CustomMessage_SignUp"
)

// AutoConfirmUser sets auto-confirm and auto-verify flags on a PreSignUp response.
func AutoConfirmUser(event *events.CognitoEventUserPoolsPreSignup, autoVerifyEmail bool) events.CognitoEventUserPoolsPreSignup {
	event.Response.AutoConfirmUser = true
	event.Response.AutoVerifyEmail = autoVerifyEmail
	return *event
}

// PreSignUpUserAttributes extracts user attributes from a PreSignUp event.
func PreSignUpUserAttributes(event *events.CognitoEventUserPoolsPreSignup) map[string]string {
	return event.Request.UserAttributes
}

// PostConfirmationUserAttributes extracts user attributes from a PostConfirmation event.
func PostConfirmationUserAttributes(event *events.CognitoEventUserPoolsPostConfirmation) map[string]string {
	return event.Request.UserAttributes
}

// AddClaimsToToken adds or overrides claims in the PreTokenGeneration response.
func AddClaimsToToken(event *events.CognitoEventUserPoolsPreTokenGen, claims map[string]string) events.CognitoEventUserPoolsPreTokenGen {
	event.Response.ClaimsOverrideDetails.ClaimsToAddOrOverride = claims
	return *event
}

// SuppressClaims suppresses claims in the PreTokenGeneration response.
func SuppressClaims(event *events.CognitoEventUserPoolsPreTokenGen, claims []string) events.CognitoEventUserPoolsPreTokenGen {
	event.Response.ClaimsOverrideDetails.ClaimsToSuppress = claims
	return *event
}

// PreAuthUserAttributes extracts user attributes from a PreAuthentication event.
func PreAuthUserAttributes(event *events.CognitoEventUserPoolsPreAuthentication) map[string]string {
	return event.Request.UserAttributes
}

// CustomizeMessage sets custom email and SMS content on a CustomMessage event.
func CustomizeMessage(event *events.CognitoEventUserPoolsCustomMessage, emailSubject, emailMessage, smsMessage string) events.CognitoEventUserPoolsCustomMessage {
	if emailSubject != "" {
		event.Response.EmailSubject = emailSubject
	}
	if emailMessage != "" {
		event.Response.EmailMessage = emailMessage
	}
	if smsMessage != "" {
		event.Response.SMSMessage = smsMessage
	}
	return *event
}

// CognitoUserName extracts the username from any Cognito trigger event header.
func CognitoUserName(header *events.CognitoEventUserPoolsHeader) string {
	return header.UserName
}

// CognitoTriggerSource extracts the trigger source from any Cognito trigger event header.
func CognitoTriggerSource(header *events.CognitoEventUserPoolsHeader) string {
	return header.TriggerSource
}
