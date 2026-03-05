package notification

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"

	"github.com/davidlramirez95/towcommand/internal/domain/event"
)

// NotificationRouter dispatches domain events to the appropriate notification
// channels (SMS, email) based on event type.
type NotificationRouter struct {
	sms         SMSSender
	email       EmailSender
	users       UserFinder
	bookings    BookingFinder
	opsPhone    string
	safetyEmail string
}

// NewNotificationRouter creates a NotificationRouter with its dependencies.
func NewNotificationRouter(
	sms SMSSender,
	email EmailSender,
	users UserFinder,
	bookings BookingFinder,
	opsPhone string,
	safetyEmail string,
) *NotificationRouter {
	return &NotificationRouter{
		sms:         sms,
		email:       email,
		users:       users,
		bookings:    bookings,
		opsPhone:    opsPhone,
		safetyEmail: safetyEmail,
	}
}

// Route dispatches a domain event to the correct notification channel based on eventType.
// Errors are logged but never propagated -- notifications are best-effort.
func (r *NotificationRouter) Route(ctx context.Context, eventType string, detail json.RawMessage) error {
	switch eventType {
	case event.BookingMatched:
		return r.handleBookingMatched(ctx, detail)
	case event.BookingCancelled:
		return r.handleBookingCancelled(ctx, detail)
	case event.BookingStatusChanged:
		return r.handleBookingStatusChanged(ctx, detail)
	case event.BookingCompleted:
		return r.handleBookingCompleted(ctx, detail)
	case event.SOSTriggered:
		return r.handleSOSTriggered(ctx, detail)
	case event.PaymentCaptured:
		return r.handlePaymentCaptured(ctx, detail)
	case event.UserRegistered:
		return r.handleUserRegistered(ctx, detail)
	default:
		slog.WarnContext(ctx, "unhandled event type for notification", "event_type", eventType)
		return nil
	}
}

// bookingMatchedDetail is the expected detail shape for BookingMatched events.
type bookingMatchedDetail struct {
	BookingID    string `json:"bookingId"`
	ProviderID   string `json:"providerId"`
	ProviderName string `json:"providerName"`
	CustomerID   string `json:"customerId"`
	ETA          int    `json:"eta"`
}

func (r *NotificationRouter) handleBookingMatched(ctx context.Context, detail json.RawMessage) error {
	var d bookingMatchedDetail
	if err := json.Unmarshal(detail, &d); err != nil {
		return fmt.Errorf("unmarshalling BookingMatched detail: %w", err)
	}

	u, err := r.users.FindByID(ctx, d.CustomerID)
	if err != nil || u == nil {
		slog.WarnContext(ctx, "could not find customer for notification", "customer_id", d.CustomerID, "error", err)
		return nil
	}

	msg := BookingMatchedSMS(d.BookingID, d.ProviderName, d.ETA)
	if err := r.sms.SendSMS(ctx, u.Phone, msg); err != nil {
		slog.ErrorContext(ctx, "failed to send BookingMatched SMS", "error", err, "phone", u.Phone)
	}
	return nil
}

// bookingCancelledDetail is the expected detail shape for BookingCancelled events.
type bookingCancelledDetail struct {
	BookingID  string `json:"bookingId"`
	CustomerID string `json:"customerId"`
	ProviderID string `json:"providerId"`
}

func (r *NotificationRouter) handleBookingCancelled(ctx context.Context, detail json.RawMessage) error {
	var d bookingCancelledDetail
	if err := json.Unmarshal(detail, &d); err != nil {
		return fmt.Errorf("unmarshalling BookingCancelled detail: %w", err)
	}

	msg := BookingCancelledSMS(d.BookingID)

	// Notify customer.
	if u, err := r.users.FindByID(ctx, d.CustomerID); err == nil && u != nil {
		if err := r.sms.SendSMS(ctx, u.Phone, msg); err != nil {
			slog.ErrorContext(ctx, "failed to send cancellation SMS to customer", "error", err)
		}
	}

	// Notify provider if assigned.
	if d.ProviderID != "" {
		if u, err := r.users.FindByID(ctx, d.ProviderID); err == nil && u != nil {
			if err := r.sms.SendSMS(ctx, u.Phone, msg); err != nil {
				slog.ErrorContext(ctx, "failed to send cancellation SMS to provider", "error", err)
			}
		}
	}

	return nil
}

// bookingStatusChangedDetail is the expected detail shape for BookingStatusChanged events.
type bookingStatusChangedDetail struct {
	BookingID  string `json:"bookingId"`
	CustomerID string `json:"customerId"`
	Status     string `json:"status"`
}

func (r *NotificationRouter) handleBookingStatusChanged(ctx context.Context, detail json.RawMessage) error {
	var d bookingStatusChangedDetail
	if err := json.Unmarshal(detail, &d); err != nil {
		return fmt.Errorf("unmarshalling BookingStatusChanged detail: %w", err)
	}

	u, err := r.users.FindByID(ctx, d.CustomerID)
	if err != nil || u == nil {
		slog.WarnContext(ctx, "could not find customer for status notification", "customer_id", d.CustomerID)
		return nil
	}

	msg := BookingStatusChangedSMS(d.BookingID, d.Status)
	if err := r.sms.SendSMS(ctx, u.Phone, msg); err != nil {
		slog.ErrorContext(ctx, "failed to send status change SMS", "error", err)
	}
	return nil
}

// bookingCompletedDetail is the expected detail shape for BookingCompleted events.
type bookingCompletedDetail struct {
	BookingID  string `json:"bookingId"`
	CustomerID string `json:"customerId"`
}

func (r *NotificationRouter) handleBookingCompleted(ctx context.Context, detail json.RawMessage) error {
	var d bookingCompletedDetail
	if err := json.Unmarshal(detail, &d); err != nil {
		return fmt.Errorf("unmarshalling BookingCompleted detail: %w", err)
	}

	u, err := r.users.FindByID(ctx, d.CustomerID)
	if err != nil || u == nil {
		slog.WarnContext(ctx, "could not find customer for completion notification", "customer_id", d.CustomerID)
		return nil
	}

	msg := BookingCompletedSMS(d.BookingID)
	if err := r.sms.SendSMS(ctx, u.Phone, msg); err != nil {
		slog.ErrorContext(ctx, "failed to send BookingCompleted SMS", "error", err)
	}
	return nil
}

// sosTriggeredDetail is the expected detail shape for SOSTriggered events.
type sosTriggeredDetail struct {
	BookingID  string  `json:"bookingId"`
	CustomerID string  `json:"customerId"`
	Lat        float64 `json:"lat"`
	Lng        float64 `json:"lng"`
	Severity   string  `json:"severity"`
}

func (r *NotificationRouter) handleSOSTriggered(ctx context.Context, detail json.RawMessage) error {
	var d sosTriggeredDetail
	if err := json.Unmarshal(detail, &d); err != nil {
		return fmt.Errorf("unmarshalling SOSTriggered detail: %w", err)
	}

	// SMS to ops phone.
	smsMsg := SOSAlertSMS(d.BookingID, d.Lat, d.Lng)
	if err := r.sms.SendSMS(ctx, r.opsPhone, smsMsg); err != nil {
		slog.ErrorContext(ctx, "failed to send SOS SMS to ops", "error", err)
	}

	// Email to safety team.
	subject := SOSAlertEmailSubject(d.BookingID)
	body := SOSAlertEmailBody(d.BookingID, d.Lat, d.Lng, d.Severity)
	if err := r.email.SendEmail(ctx, r.safetyEmail, subject, body); err != nil {
		slog.ErrorContext(ctx, "failed to send SOS email to safety team", "error", err)
	}

	return nil
}

// paymentCapturedDetail is the expected detail shape for PaymentCaptured events.
type paymentCapturedDetail struct {
	BookingID      string `json:"bookingId"`
	CustomerID     string `json:"customerId"`
	AmountCentavos int64  `json:"amountCentavos"`
}

func (r *NotificationRouter) handlePaymentCaptured(ctx context.Context, detail json.RawMessage) error {
	var d paymentCapturedDetail
	if err := json.Unmarshal(detail, &d); err != nil {
		return fmt.Errorf("unmarshalling PaymentCaptured detail: %w", err)
	}

	u, err := r.users.FindByID(ctx, d.CustomerID)
	if err != nil || u == nil {
		slog.WarnContext(ctx, "could not find customer for payment notification", "customer_id", d.CustomerID)
		return nil
	}

	msg := PaymentCapturedSMS(d.AmountCentavos)
	if err := r.sms.SendSMS(ctx, u.Phone, msg); err != nil {
		slog.ErrorContext(ctx, "failed to send PaymentCaptured SMS", "error", err)
	}
	return nil
}

// userRegisteredDetail is the expected detail shape for UserRegistered events.
type userRegisteredDetail struct {
	UserID string `json:"userId"`
	Phone  string `json:"phone"`
}

func (r *NotificationRouter) handleUserRegistered(ctx context.Context, detail json.RawMessage) error {
	var d userRegisteredDetail
	if err := json.Unmarshal(detail, &d); err != nil {
		return fmt.Errorf("unmarshalling UserRegistered detail: %w", err)
	}

	phone := d.Phone
	if phone == "" {
		// Fallback: look up user to get phone.
		u, err := r.users.FindByID(ctx, d.UserID)
		if err != nil || u == nil {
			slog.WarnContext(ctx, "could not find user for welcome notification", "user_id", d.UserID)
			return nil
		}
		phone = u.Phone
	}

	msg := WelcomeSMS()
	if err := r.sms.SendSMS(ctx, phone, msg); err != nil {
		slog.ErrorContext(ctx, "failed to send welcome SMS", "error", err)
	}
	return nil
}
