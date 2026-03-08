package notification

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/davidlramirez95/towcommand/internal/domain/event"
	"github.com/davidlramirez95/towcommand/internal/domain/user"
	"github.com/davidlramirez95/towcommand/internal/usecase/port"
)

// --- Push Mocks ---

type mockPushSender struct{ mock.Mock }

func (m *mockPushSender) SendPush(ctx context.Context, endpointArn, title, message string, data map[string]string) error {
	args := m.Called(ctx, endpointArn, title, message, data)
	return args.Error(0)
}

type mockPushTokenFinder struct{ mock.Mock }

func (m *mockPushTokenFinder) FindByUserID(ctx context.Context, userID string) ([]port.PushToken, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]port.PushToken), args.Error(1)
}

func newPushTestRouter(
	sms *mockSMSSender,
	email *mockEmailSender,
	users *mockUserFinder,
	bookings *mockBookingFinder,
	push *mockPushSender,
	pushTokens *mockPushTokenFinder,
) *PushNotificationRouter {
	base := NewNotificationRouter(sms, email, users, bookings, "+639170000000", "safety@towcommand.ph")
	return NewPushNotificationRouter(base, push, pushTokens)
}

// --- Tests ---

func TestPushRoute_BookingMatched(t *testing.T) {
	sms := new(mockSMSSender)
	email := new(mockEmailSender)
	users := new(mockUserFinder)
	bookings := new(mockBookingFinder)
	push := new(mockPushSender)
	pushTokens := new(mockPushTokenFinder)
	router := newPushTestRouter(sms, email, users, bookings, push, pushTokens)

	users.On("FindByID", mock.Anything, "cust-1").Return(&user.User{
		Phone: "+639171234567",
	}, nil)
	sms.On("SendSMS", mock.Anything, "+639171234567", mock.Anything).Return(nil)

	pushTokens.On("FindByUserID", mock.Anything, "cust-1").Return([]port.PushToken{
		{EndpointArn: "arn:endpoint-1"},
		{EndpointArn: "arn:endpoint-2"},
	}, nil)
	push.On("SendPush", mock.Anything, "arn:endpoint-1", "Driver On The Way", mock.Anything, mock.Anything).Return(nil)
	push.On("SendPush", mock.Anything, "arn:endpoint-2", "Driver On The Way", mock.Anything, mock.Anything).Return(nil)

	detail := mustJSON(t, bookingMatchedDetail{
		BookingID:    "BK-001",
		ProviderID:   "prov-1",
		ProviderName: "Juan",
		CustomerID:   "cust-1",
		ETA:          10,
	})

	err := router.Route(context.Background(), event.BookingMatched, detail)
	require.NoError(t, err)
	sms.AssertExpectations(t)
	push.AssertExpectations(t)
	pushTokens.AssertExpectations(t)
}

func TestPushRoute_BookingCancelled_BothParties(t *testing.T) {
	sms := new(mockSMSSender)
	email := new(mockEmailSender)
	users := new(mockUserFinder)
	bookings := new(mockBookingFinder)
	push := new(mockPushSender)
	pushTokens := new(mockPushTokenFinder)
	router := newPushTestRouter(sms, email, users, bookings, push, pushTokens)

	users.On("FindByID", mock.Anything, "cust-1").Return(&user.User{Phone: "+63900001111"}, nil)
	users.On("FindByID", mock.Anything, "prov-1").Return(&user.User{Phone: "+63900002222"}, nil)
	sms.On("SendSMS", mock.Anything, "+63900001111", mock.Anything).Return(nil)
	sms.On("SendSMS", mock.Anything, "+63900002222", mock.Anything).Return(nil)

	pushTokens.On("FindByUserID", mock.Anything, "cust-1").Return([]port.PushToken{
		{EndpointArn: "arn:cust-ep"},
	}, nil)
	pushTokens.On("FindByUserID", mock.Anything, "prov-1").Return([]port.PushToken{
		{EndpointArn: "arn:prov-ep"},
	}, nil)
	push.On("SendPush", mock.Anything, "arn:cust-ep", "Booking Cancelled", mock.Anything, mock.Anything).Return(nil)
	push.On("SendPush", mock.Anything, "arn:prov-ep", "Booking Cancelled", mock.Anything, mock.Anything).Return(nil)

	detail := mustJSON(t, bookingCancelledDetail{
		BookingID:  "BK-002",
		CustomerID: "cust-1",
		ProviderID: "prov-1",
	})

	err := router.Route(context.Background(), event.BookingCancelled, detail)
	require.NoError(t, err)
	push.AssertNumberOfCalls(t, "SendPush", 2)
}

func TestPushRoute_BookingStatusChanged(t *testing.T) {
	sms := new(mockSMSSender)
	email := new(mockEmailSender)
	users := new(mockUserFinder)
	bookings := new(mockBookingFinder)
	push := new(mockPushSender)
	pushTokens := new(mockPushTokenFinder)
	router := newPushTestRouter(sms, email, users, bookings, push, pushTokens)

	users.On("FindByID", mock.Anything, "cust-1").Return(&user.User{Phone: "+63912345678"}, nil)
	sms.On("SendSMS", mock.Anything, "+63912345678", mock.Anything).Return(nil)

	pushTokens.On("FindByUserID", mock.Anything, "cust-1").Return([]port.PushToken{
		{EndpointArn: "arn:ep-status"},
	}, nil)
	push.On("SendPush", mock.Anything, "arn:ep-status", "Booking Update", mock.Anything, mock.Anything).Return(nil)

	detail := mustJSON(t, bookingStatusChangedDetail{
		BookingID:  "BK-003",
		CustomerID: "cust-1",
		Status:     "EN_ROUTE",
	})

	err := router.Route(context.Background(), event.BookingStatusChanged, detail)
	require.NoError(t, err)
	push.AssertExpectations(t)
}

func TestPushRoute_BookingCompleted(t *testing.T) {
	sms := new(mockSMSSender)
	email := new(mockEmailSender)
	users := new(mockUserFinder)
	bookings := new(mockBookingFinder)
	push := new(mockPushSender)
	pushTokens := new(mockPushTokenFinder)
	router := newPushTestRouter(sms, email, users, bookings, push, pushTokens)

	users.On("FindByID", mock.Anything, "cust-1").Return(&user.User{Phone: "+63999888777"}, nil)
	sms.On("SendSMS", mock.Anything, "+63999888777", mock.Anything).Return(nil)

	pushTokens.On("FindByUserID", mock.Anything, "cust-1").Return([]port.PushToken{
		{EndpointArn: "arn:ep-done"},
	}, nil)
	push.On("SendPush", mock.Anything, "arn:ep-done", "Booking Completed", mock.Anything, mock.Anything).Return(nil)

	detail := mustJSON(t, bookingCompletedDetail{
		BookingID:  "BK-004",
		CustomerID: "cust-1",
	})

	err := router.Route(context.Background(), event.BookingCompleted, detail)
	require.NoError(t, err)
	push.AssertExpectations(t)
}

func TestPushRoute_PaymentCaptured(t *testing.T) {
	sms := new(mockSMSSender)
	email := new(mockEmailSender)
	users := new(mockUserFinder)
	bookings := new(mockBookingFinder)
	push := new(mockPushSender)
	pushTokens := new(mockPushTokenFinder)
	router := newPushTestRouter(sms, email, users, bookings, push, pushTokens)

	users.On("FindByID", mock.Anything, "cust-pay").Return(&user.User{Phone: "+63900009999"}, nil)
	sms.On("SendSMS", mock.Anything, "+63900009999", mock.Anything).Return(nil)

	pushTokens.On("FindByUserID", mock.Anything, "cust-pay").Return([]port.PushToken{
		{EndpointArn: "arn:ep-pay"},
	}, nil)
	push.On("SendPush", mock.Anything, "arn:ep-pay", "Payment Received", mock.Anything, mock.Anything).Return(nil)

	detail := mustJSON(t, paymentCapturedDetail{
		BookingID:      "BK-PAY",
		CustomerID:     "cust-pay",
		AmountCentavos: 250000,
	})

	err := router.Route(context.Background(), event.PaymentCaptured, detail)
	require.NoError(t, err)
	push.AssertExpectations(t)
}

func TestPushRoute_UserRegistered(t *testing.T) {
	sms := new(mockSMSSender)
	email := new(mockEmailSender)
	users := new(mockUserFinder)
	bookings := new(mockBookingFinder)
	push := new(mockPushSender)
	pushTokens := new(mockPushTokenFinder)
	router := newPushTestRouter(sms, email, users, bookings, push, pushTokens)

	sms.On("SendSMS", mock.Anything, "+63917WELCOME", mock.Anything).Return(nil)

	// No push tokens registered yet for new user — no push expected.
	pushTokens.On("FindByUserID", mock.Anything, "user-new").Return([]port.PushToken{}, nil)

	detail := mustJSON(t, userRegisteredDetail{
		UserID: "user-new",
		Phone:  "+63917WELCOME",
	})

	err := router.Route(context.Background(), event.UserRegistered, detail)
	require.NoError(t, err)
	sms.AssertExpectations(t)
	push.AssertNotCalled(t, "SendPush")
}

func TestPushRoute_SOSTriggered_NoPush(t *testing.T) {
	sms := new(mockSMSSender)
	email := new(mockEmailSender)
	users := new(mockUserFinder)
	bookings := new(mockBookingFinder)
	push := new(mockPushSender)
	pushTokens := new(mockPushTokenFinder)
	router := newPushTestRouter(sms, email, users, bookings, push, pushTokens)

	sms.On("SendSMS", mock.Anything, "+639170000000", mock.Anything).Return(nil)
	email.On("SendEmail", mock.Anything, "safety@towcommand.ph", mock.Anything, mock.Anything).Return(nil)

	detail := mustJSON(t, sosTriggeredDetail{
		BookingID:  "BK-SOS",
		CustomerID: "cust-sos",
		Lat:        14.5995,
		Lng:        120.9842,
		Severity:   "HIGH",
	})

	err := router.Route(context.Background(), event.SOSTriggered, detail)
	require.NoError(t, err)
	push.AssertNotCalled(t, "SendPush")
}

func TestPushRoute_NilPushSender_NoError(t *testing.T) {
	sms := new(mockSMSSender)
	email := new(mockEmailSender)
	users := new(mockUserFinder)
	bookings := new(mockBookingFinder)

	base := NewNotificationRouter(sms, email, users, bookings, "+639170000000", "safety@towcommand.ph")
	// Create with nil push sender and token finder.
	router := NewPushNotificationRouter(base, nil, nil)

	users.On("FindByID", mock.Anything, "cust-1").Return(&user.User{
		Phone: "+639171234567",
	}, nil)
	sms.On("SendSMS", mock.Anything, "+639171234567", mock.Anything).Return(nil)

	detail := mustJSON(t, bookingMatchedDetail{
		BookingID:    "BK-NIL",
		ProviderID:   "prov-1",
		ProviderName: "Driver",
		CustomerID:   "cust-1",
		ETA:          5,
	})

	// Should work fine without push — just SMS.
	err := router.Route(context.Background(), event.BookingMatched, detail)
	require.NoError(t, err)
	sms.AssertExpectations(t)
}

func TestPushRoute_PushSendError_NotPropagated(t *testing.T) {
	sms := new(mockSMSSender)
	email := new(mockEmailSender)
	users := new(mockUserFinder)
	bookings := new(mockBookingFinder)
	push := new(mockPushSender)
	pushTokens := new(mockPushTokenFinder)
	router := newPushTestRouter(sms, email, users, bookings, push, pushTokens)

	users.On("FindByID", mock.Anything, "cust-err").Return(&user.User{Phone: "+63900009999"}, nil)
	sms.On("SendSMS", mock.Anything, "+63900009999", mock.Anything).Return(nil)

	pushTokens.On("FindByUserID", mock.Anything, "cust-err").Return([]port.PushToken{
		{EndpointArn: "arn:ep-err"},
	}, nil)
	push.On("SendPush", mock.Anything, "arn:ep-err", "Driver On The Way", mock.Anything, mock.Anything).Return(assert.AnError)

	detail := mustJSON(t, bookingMatchedDetail{
		BookingID:    "BK-ERR",
		CustomerID:   "cust-err",
		ProviderName: "Driver",
		ETA:          5,
	})

	// Push error is logged but does not propagate.
	err := router.Route(context.Background(), event.BookingMatched, detail)
	assert.NoError(t, err)
	push.AssertExpectations(t)
}

func TestPushRoute_UnhandledEvent_NoPush(t *testing.T) {
	sms := new(mockSMSSender)
	email := new(mockEmailSender)
	users := new(mockUserFinder)
	bookings := new(mockBookingFinder)
	push := new(mockPushSender)
	pushTokens := new(mockPushTokenFinder)
	router := newPushTestRouter(sms, email, users, bookings, push, pushTokens)

	err := router.Route(context.Background(), "SomeUnknownEvent", json.RawMessage(`{}`))
	assert.NoError(t, err)
	push.AssertNotCalled(t, "SendPush")
	pushTokens.AssertNotCalled(t, "FindByUserID")
}

func TestPushRoute_InvalidJSON_NoPanic(t *testing.T) {
	sms := new(mockSMSSender)
	email := new(mockEmailSender)
	users := new(mockUserFinder)
	bookings := new(mockBookingFinder)
	push := new(mockPushSender)
	pushTokens := new(mockPushTokenFinder)
	router := newPushTestRouter(sms, email, users, bookings, push, pushTokens)

	// BookingMatched with invalid JSON — base router returns error, push silently skips.
	err := router.Route(context.Background(), event.BookingMatched, json.RawMessage(`{invalid`))
	assert.Error(t, err) // from base router
	push.AssertNotCalled(t, "SendPush")
}

func TestPushRoute_TokenFinderError_NoPush(t *testing.T) {
	sms := new(mockSMSSender)
	email := new(mockEmailSender)
	users := new(mockUserFinder)
	bookings := new(mockBookingFinder)
	push := new(mockPushSender)
	pushTokens := new(mockPushTokenFinder)
	router := newPushTestRouter(sms, email, users, bookings, push, pushTokens)

	users.On("FindByID", mock.Anything, "cust-1").Return(&user.User{Phone: "+639171234567"}, nil)
	sms.On("SendSMS", mock.Anything, "+639171234567", mock.Anything).Return(nil)

	// Token finder returns error — push silently skipped.
	pushTokens.On("FindByUserID", mock.Anything, "cust-1").Return(nil, assert.AnError)

	detail := mustJSON(t, bookingMatchedDetail{
		BookingID:    "BK-TKERR",
		ProviderID:   "prov-1",
		ProviderName: "Juan",
		CustomerID:   "cust-1",
		ETA:          10,
	})

	err := router.Route(context.Background(), event.BookingMatched, detail)
	require.NoError(t, err)
	push.AssertNotCalled(t, "SendPush")
}
