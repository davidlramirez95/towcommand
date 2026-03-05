package notification

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/davidlramirez95/towcommand/internal/domain/booking"
	"github.com/davidlramirez95/towcommand/internal/domain/event"
	"github.com/davidlramirez95/towcommand/internal/domain/user"
)

// --- Mocks ---

type mockSMSSender struct{ mock.Mock }

func (m *mockSMSSender) SendSMS(ctx context.Context, phoneNumber, message string) error {
	args := m.Called(ctx, phoneNumber, message)
	return args.Error(0)
}

type mockEmailSender struct{ mock.Mock }

func (m *mockEmailSender) SendEmail(ctx context.Context, to, subject, htmlBody string) error {
	args := m.Called(ctx, to, subject, htmlBody)
	return args.Error(0)
}

type mockUserFinder struct{ mock.Mock }

func (m *mockUserFinder) FindByID(ctx context.Context, userID string) (*user.User, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*user.User), args.Error(1)
}

type mockBookingFinder struct{ mock.Mock }

func (m *mockBookingFinder) FindByID(ctx context.Context, bookingID string) (*booking.Booking, error) {
	args := m.Called(ctx, bookingID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*booking.Booking), args.Error(1)
}

func newTestRouter(sms *mockSMSSender, email *mockEmailSender, users *mockUserFinder, bookings *mockBookingFinder) *NotificationRouter {
	return NewNotificationRouter(sms, email, users, bookings, "+639170000000", "safety@towcommand.ph")
}

func mustJSON(t *testing.T, v any) json.RawMessage {
	t.Helper()
	b, err := json.Marshal(v)
	require.NoError(t, err)
	return b
}

// --- Tests ---

func TestRoute_BookingMatched(t *testing.T) {
	sms := new(mockSMSSender)
	email := new(mockEmailSender)
	users := new(mockUserFinder)
	bookings := new(mockBookingFinder)
	router := newTestRouter(sms, email, users, bookings)

	users.On("FindByID", mock.Anything, "cust-1").Return(&user.User{
		Phone: "+639171234567",
	}, nil)
	sms.On("SendSMS", mock.Anything, "+639171234567", mock.MatchedBy(func(msg string) bool {
		return msg != ""
	})).Return(nil)

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
	users.AssertExpectations(t)
}

func TestRoute_BookingCancelled_BothParties(t *testing.T) {
	sms := new(mockSMSSender)
	email := new(mockEmailSender)
	users := new(mockUserFinder)
	bookings := new(mockBookingFinder)
	router := newTestRouter(sms, email, users, bookings)

	users.On("FindByID", mock.Anything, "cust-1").Return(&user.User{Phone: "+63900001111"}, nil)
	users.On("FindByID", mock.Anything, "prov-1").Return(&user.User{Phone: "+63900002222"}, nil)
	sms.On("SendSMS", mock.Anything, "+63900001111", mock.Anything).Return(nil)
	sms.On("SendSMS", mock.Anything, "+63900002222", mock.Anything).Return(nil)

	detail := mustJSON(t, bookingCancelledDetail{
		BookingID:  "BK-002",
		CustomerID: "cust-1",
		ProviderID: "prov-1",
	})

	err := router.Route(context.Background(), event.BookingCancelled, detail)
	require.NoError(t, err)
	sms.AssertNumberOfCalls(t, "SendSMS", 2)
}

func TestRoute_BookingStatusChanged(t *testing.T) {
	sms := new(mockSMSSender)
	email := new(mockEmailSender)
	users := new(mockUserFinder)
	bookings := new(mockBookingFinder)
	router := newTestRouter(sms, email, users, bookings)

	users.On("FindByID", mock.Anything, "cust-1").Return(&user.User{Phone: "+63912345678"}, nil)
	sms.On("SendSMS", mock.Anything, "+63912345678", mock.MatchedBy(func(msg string) bool {
		return msg != ""
	})).Return(nil)

	detail := mustJSON(t, bookingStatusChangedDetail{
		BookingID:  "BK-003",
		CustomerID: "cust-1",
		Status:     "EN_ROUTE",
	})

	err := router.Route(context.Background(), event.BookingStatusChanged, detail)
	require.NoError(t, err)
	sms.AssertExpectations(t)
}

func TestRoute_BookingCompleted(t *testing.T) {
	sms := new(mockSMSSender)
	email := new(mockEmailSender)
	users := new(mockUserFinder)
	bookings := new(mockBookingFinder)
	router := newTestRouter(sms, email, users, bookings)

	users.On("FindByID", mock.Anything, "cust-1").Return(&user.User{Phone: "+63999888777"}, nil)
	sms.On("SendSMS", mock.Anything, "+63999888777", mock.MatchedBy(func(msg string) bool {
		return msg != ""
	})).Return(nil)

	detail := mustJSON(t, bookingCompletedDetail{
		BookingID:  "BK-004",
		CustomerID: "cust-1",
	})

	err := router.Route(context.Background(), event.BookingCompleted, detail)
	require.NoError(t, err)
	sms.AssertExpectations(t)
}

func TestRoute_SOSTriggered(t *testing.T) {
	sms := new(mockSMSSender)
	email := new(mockEmailSender)
	users := new(mockUserFinder)
	bookings := new(mockBookingFinder)
	router := newTestRouter(sms, email, users, bookings)

	sms.On("SendSMS", mock.Anything, "+639170000000", mock.MatchedBy(func(msg string) bool {
		return msg != ""
	})).Return(nil)
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
	sms.AssertExpectations(t)
	email.AssertExpectations(t)
}

func TestRoute_PaymentCaptured(t *testing.T) {
	sms := new(mockSMSSender)
	email := new(mockEmailSender)
	users := new(mockUserFinder)
	bookings := new(mockBookingFinder)
	router := newTestRouter(sms, email, users, bookings)

	users.On("FindByID", mock.Anything, "cust-pay").Return(&user.User{Phone: "+63900009999"}, nil)
	sms.On("SendSMS", mock.Anything, "+63900009999", mock.MatchedBy(func(msg string) bool {
		return msg != ""
	})).Return(nil)

	detail := mustJSON(t, paymentCapturedDetail{
		BookingID:      "BK-PAY",
		CustomerID:     "cust-pay",
		AmountCentavos: 250000,
	})

	err := router.Route(context.Background(), event.PaymentCaptured, detail)
	require.NoError(t, err)
	sms.AssertExpectations(t)
}

func TestRoute_UserRegistered(t *testing.T) {
	sms := new(mockSMSSender)
	email := new(mockEmailSender)
	users := new(mockUserFinder)
	bookings := new(mockBookingFinder)
	router := newTestRouter(sms, email, users, bookings)

	sms.On("SendSMS", mock.Anything, "+63917WELCOME", mock.MatchedBy(func(msg string) bool {
		return msg != ""
	})).Return(nil)

	detail := mustJSON(t, userRegisteredDetail{
		UserID: "user-new",
		Phone:  "+63917WELCOME",
	})

	err := router.Route(context.Background(), event.UserRegistered, detail)
	require.NoError(t, err)
	sms.AssertExpectations(t)
}

func TestRoute_UserRegistered_FallbackLookup(t *testing.T) {
	sms := new(mockSMSSender)
	email := new(mockEmailSender)
	users := new(mockUserFinder)
	bookings := new(mockBookingFinder)
	router := newTestRouter(sms, email, users, bookings)

	users.On("FindByID", mock.Anything, "user-nophone").Return(&user.User{Phone: "+63918FALLBACK"}, nil)
	sms.On("SendSMS", mock.Anything, "+63918FALLBACK", mock.Anything).Return(nil)

	detail := mustJSON(t, userRegisteredDetail{
		UserID: "user-nophone",
		Phone:  "",
	})

	err := router.Route(context.Background(), event.UserRegistered, detail)
	require.NoError(t, err)
	sms.AssertExpectations(t)
	users.AssertExpectations(t)
}

func TestRoute_UnhandledEvent(t *testing.T) {
	sms := new(mockSMSSender)
	email := new(mockEmailSender)
	users := new(mockUserFinder)
	bookings := new(mockBookingFinder)
	router := newTestRouter(sms, email, users, bookings)

	err := router.Route(context.Background(), "SomeUnknownEvent", json.RawMessage(`{}`))
	assert.NoError(t, err)
	sms.AssertNotCalled(t, "SendSMS")
	email.AssertNotCalled(t, "SendEmail")
}

func TestRoute_InvalidJSON(t *testing.T) {
	sms := new(mockSMSSender)
	email := new(mockEmailSender)
	users := new(mockUserFinder)
	bookings := new(mockBookingFinder)
	router := newTestRouter(sms, email, users, bookings)

	err := router.Route(context.Background(), event.BookingMatched, json.RawMessage(`{invalid`))
	assert.Error(t, err)
}

func TestRoute_CustomerNotFound_NoError(t *testing.T) {
	sms := new(mockSMSSender)
	email := new(mockEmailSender)
	users := new(mockUserFinder)
	bookings := new(mockBookingFinder)
	router := newTestRouter(sms, email, users, bookings)

	users.On("FindByID", mock.Anything, "ghost-user").Return(nil, nil)

	detail := mustJSON(t, bookingMatchedDetail{
		BookingID:    "BK-GHOST",
		CustomerID:   "ghost-user",
		ProviderName: "Driver",
		ETA:          5,
	})

	err := router.Route(context.Background(), event.BookingMatched, detail)
	assert.NoError(t, err)
	sms.AssertNotCalled(t, "SendSMS")
}

func TestRoute_BookingCancelled_NoProvider(t *testing.T) {
	sms := new(mockSMSSender)
	email := new(mockEmailSender)
	users := new(mockUserFinder)
	bookings := new(mockBookingFinder)
	router := newTestRouter(sms, email, users, bookings)

	users.On("FindByID", mock.Anything, "cust-solo").Return(&user.User{Phone: "+63900001111"}, nil)
	sms.On("SendSMS", mock.Anything, "+63900001111", mock.Anything).Return(nil)

	detail := mustJSON(t, bookingCancelledDetail{
		BookingID:  "BK-SOLO",
		CustomerID: "cust-solo",
		ProviderID: "",
	})

	err := router.Route(context.Background(), event.BookingCancelled, detail)
	require.NoError(t, err)
	sms.AssertNumberOfCalls(t, "SendSMS", 1) // only customer notified
}

func TestRoute_BookingCancelled_InvalidJSON(t *testing.T) {
	sms := new(mockSMSSender)
	email := new(mockEmailSender)
	users := new(mockUserFinder)
	bookings := new(mockBookingFinder)
	router := newTestRouter(sms, email, users, bookings)

	err := router.Route(context.Background(), event.BookingCancelled, json.RawMessage(`{bad`))
	assert.Error(t, err)
}

func TestRoute_BookingStatusChanged_InvalidJSON(t *testing.T) {
	sms := new(mockSMSSender)
	email := new(mockEmailSender)
	users := new(mockUserFinder)
	bookings := new(mockBookingFinder)
	router := newTestRouter(sms, email, users, bookings)

	err := router.Route(context.Background(), event.BookingStatusChanged, json.RawMessage(`{bad`))
	assert.Error(t, err)
}

func TestRoute_BookingStatusChanged_CustomerNotFound(t *testing.T) {
	sms := new(mockSMSSender)
	email := new(mockEmailSender)
	users := new(mockUserFinder)
	bookings := new(mockBookingFinder)
	router := newTestRouter(sms, email, users, bookings)

	users.On("FindByID", mock.Anything, "ghost-cust").Return(nil, nil)

	detail := mustJSON(t, bookingStatusChangedDetail{
		BookingID:  "BK-GHOST",
		CustomerID: "ghost-cust",
		Status:     "ARRIVED",
	})

	err := router.Route(context.Background(), event.BookingStatusChanged, detail)
	assert.NoError(t, err)
	sms.AssertNotCalled(t, "SendSMS")
}

func TestRoute_BookingCompleted_InvalidJSON(t *testing.T) {
	sms := new(mockSMSSender)
	email := new(mockEmailSender)
	users := new(mockUserFinder)
	bookings := new(mockBookingFinder)
	router := newTestRouter(sms, email, users, bookings)

	err := router.Route(context.Background(), event.BookingCompleted, json.RawMessage(`{bad`))
	assert.Error(t, err)
}

func TestRoute_BookingCompleted_CustomerNotFound(t *testing.T) {
	sms := new(mockSMSSender)
	email := new(mockEmailSender)
	users := new(mockUserFinder)
	bookings := new(mockBookingFinder)
	router := newTestRouter(sms, email, users, bookings)

	users.On("FindByID", mock.Anything, "gone-cust").Return(nil, nil)

	detail := mustJSON(t, bookingCompletedDetail{
		BookingID:  "BK-NOCUST",
		CustomerID: "gone-cust",
	})

	err := router.Route(context.Background(), event.BookingCompleted, detail)
	assert.NoError(t, err)
	sms.AssertNotCalled(t, "SendSMS")
}

func TestRoute_SOSTriggered_InvalidJSON(t *testing.T) {
	sms := new(mockSMSSender)
	email := new(mockEmailSender)
	users := new(mockUserFinder)
	bookings := new(mockBookingFinder)
	router := newTestRouter(sms, email, users, bookings)

	err := router.Route(context.Background(), event.SOSTriggered, json.RawMessage(`{bad`))
	assert.Error(t, err)
}

func TestRoute_PaymentCaptured_InvalidJSON(t *testing.T) {
	sms := new(mockSMSSender)
	email := new(mockEmailSender)
	users := new(mockUserFinder)
	bookings := new(mockBookingFinder)
	router := newTestRouter(sms, email, users, bookings)

	err := router.Route(context.Background(), event.PaymentCaptured, json.RawMessage(`{bad`))
	assert.Error(t, err)
}

func TestRoute_PaymentCaptured_CustomerNotFound(t *testing.T) {
	sms := new(mockSMSSender)
	email := new(mockEmailSender)
	users := new(mockUserFinder)
	bookings := new(mockBookingFinder)
	router := newTestRouter(sms, email, users, bookings)

	users.On("FindByID", mock.Anything, "no-cust").Return(nil, nil)

	detail := mustJSON(t, paymentCapturedDetail{
		BookingID:      "BK-NOPAY",
		CustomerID:     "no-cust",
		AmountCentavos: 100000,
	})

	err := router.Route(context.Background(), event.PaymentCaptured, detail)
	assert.NoError(t, err)
	sms.AssertNotCalled(t, "SendSMS")
}

func TestRoute_UserRegistered_InvalidJSON(t *testing.T) {
	sms := new(mockSMSSender)
	email := new(mockEmailSender)
	users := new(mockUserFinder)
	bookings := new(mockBookingFinder)
	router := newTestRouter(sms, email, users, bookings)

	err := router.Route(context.Background(), event.UserRegistered, json.RawMessage(`{bad`))
	assert.Error(t, err)
}

func TestRoute_UserRegistered_NoPhoneNoUser(t *testing.T) {
	sms := new(mockSMSSender)
	email := new(mockEmailSender)
	users := new(mockUserFinder)
	bookings := new(mockBookingFinder)
	router := newTestRouter(sms, email, users, bookings)

	users.On("FindByID", mock.Anything, "ghost-new").Return(nil, nil)

	detail := mustJSON(t, userRegisteredDetail{
		UserID: "ghost-new",
		Phone:  "",
	})

	err := router.Route(context.Background(), event.UserRegistered, detail)
	assert.NoError(t, err)
	sms.AssertNotCalled(t, "SendSMS")
}

func TestRoute_BookingMatched_SMSSendError(t *testing.T) {
	sms := new(mockSMSSender)
	email := new(mockEmailSender)
	users := new(mockUserFinder)
	bookings := new(mockBookingFinder)
	router := newTestRouter(sms, email, users, bookings)

	users.On("FindByID", mock.Anything, "cust-err").Return(&user.User{Phone: "+63900009999"}, nil)
	sms.On("SendSMS", mock.Anything, "+63900009999", mock.Anything).Return(assert.AnError)

	detail := mustJSON(t, bookingMatchedDetail{
		BookingID:    "BK-ERR",
		CustomerID:   "cust-err",
		ProviderName: "Driver",
		ETA:          5,
	})

	// SMS error is logged but does not propagate.
	err := router.Route(context.Background(), event.BookingMatched, detail)
	assert.NoError(t, err)
	sms.AssertExpectations(t)
}

func TestRoute_SOSTriggered_SMSAndEmailErrors(t *testing.T) {
	sms := new(mockSMSSender)
	email := new(mockEmailSender)
	users := new(mockUserFinder)
	bookings := new(mockBookingFinder)
	router := newTestRouter(sms, email, users, bookings)

	sms.On("SendSMS", mock.Anything, "+639170000000", mock.Anything).Return(assert.AnError)
	email.On("SendEmail", mock.Anything, "safety@towcommand.ph", mock.Anything, mock.Anything).Return(assert.AnError)

	detail := mustJSON(t, sosTriggeredDetail{
		BookingID:  "BK-SOS-ERR",
		CustomerID: "cust-sos",
		Lat:        14.5995,
		Lng:        120.9842,
		Severity:   "CRITICAL",
	})

	// Errors are logged but do not propagate.
	err := router.Route(context.Background(), event.SOSTriggered, detail)
	assert.NoError(t, err)
	sms.AssertExpectations(t)
	email.AssertExpectations(t)
}
