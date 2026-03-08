package notification

import (
	"context"
	"encoding/json"
	"log/slog"

	"github.com/davidlramirez95/towcommand/internal/domain/event"
)

// PushNotificationRouter extends NotificationRouter to also send push
// notifications alongside SMS/email. It decorates the base router's Route
// method to dispatch push messages after each event is handled.
type PushNotificationRouter struct {
	*NotificationRouter
	push       PushSender
	pushTokens PushTokenFinder
}

// NewPushNotificationRouter creates a PushNotificationRouter that wraps the
// given base router and adds push notifications. If push or pushTokens is nil,
// push notifications are silently skipped.
func NewPushNotificationRouter(
	base *NotificationRouter,
	push PushSender,
	pushTokens PushTokenFinder,
) *PushNotificationRouter {
	return &PushNotificationRouter{
		NotificationRouter: base,
		push:               push,
		pushTokens:         pushTokens,
	}
}

// Route dispatches the event to the base router and then sends push
// notifications to the appropriate users.
func (r *PushNotificationRouter) Route(ctx context.Context, eventType string, detail json.RawMessage) error {
	// Always run the base router first (SMS + email).
	err := r.NotificationRouter.Route(ctx, eventType, detail)

	// Send push notifications alongside SMS — best-effort, never fail.
	r.routePush(ctx, eventType, detail)

	return err
}

// routePush dispatches push notifications based on the event type.
func (r *PushNotificationRouter) routePush(ctx context.Context, eventType string, detail json.RawMessage) {
	if r.push == nil || r.pushTokens == nil {
		return
	}

	switch eventType {
	case event.BookingMatched:
		r.pushBookingMatched(ctx, detail)
	case event.BookingCancelled:
		r.pushBookingCancelled(ctx, detail)
	case event.BookingStatusChanged:
		r.pushBookingStatusChanged(ctx, detail)
	case event.BookingCompleted:
		r.pushBookingCompleted(ctx, detail)
	case event.PaymentCaptured:
		r.pushPaymentCaptured(ctx, detail)
	case event.UserRegistered:
		r.pushUserRegistered(ctx, detail)
		// SOSTriggered goes to ops, not to user devices — no push.
	}
}

func (r *PushNotificationRouter) pushBookingMatched(ctx context.Context, detail json.RawMessage) {
	var d bookingMatchedDetail
	if err := json.Unmarshal(detail, &d); err != nil {
		return
	}
	msg := BookingMatchedSMS(d.BookingID, d.ProviderName, d.ETA)
	r.sendPush(ctx, d.CustomerID, "Driver On The Way", msg)
}

func (r *PushNotificationRouter) pushBookingCancelled(ctx context.Context, detail json.RawMessage) {
	var d bookingCancelledDetail
	if err := json.Unmarshal(detail, &d); err != nil {
		return
	}
	msg := BookingCancelledSMS(d.BookingID)
	r.sendPush(ctx, d.CustomerID, "Booking Cancelled", msg)
	if d.ProviderID != "" {
		r.sendPush(ctx, d.ProviderID, "Booking Cancelled", msg)
	}
}

func (r *PushNotificationRouter) pushBookingStatusChanged(ctx context.Context, detail json.RawMessage) {
	var d bookingStatusChangedDetail
	if err := json.Unmarshal(detail, &d); err != nil {
		return
	}
	msg := BookingStatusChangedSMS(d.BookingID, d.Status)
	r.sendPush(ctx, d.CustomerID, "Booking Update", msg)
}

func (r *PushNotificationRouter) pushBookingCompleted(ctx context.Context, detail json.RawMessage) {
	var d bookingCompletedDetail
	if err := json.Unmarshal(detail, &d); err != nil {
		return
	}
	msg := BookingCompletedSMS(d.BookingID)
	r.sendPush(ctx, d.CustomerID, "Booking Completed", msg)
}

func (r *PushNotificationRouter) pushPaymentCaptured(ctx context.Context, detail json.RawMessage) {
	var d paymentCapturedDetail
	if err := json.Unmarshal(detail, &d); err != nil {
		return
	}
	msg := PaymentCapturedSMS(d.AmountCentavos)
	r.sendPush(ctx, d.CustomerID, "Payment Received", msg)
}

func (r *PushNotificationRouter) pushUserRegistered(ctx context.Context, detail json.RawMessage) {
	var d userRegisteredDetail
	if err := json.Unmarshal(detail, &d); err != nil {
		return
	}
	msg := WelcomeSMS()
	r.sendPush(ctx, d.UserID, "Welcome to TowCommand", msg)
}

// sendPush sends a push notification to all registered devices for a user.
func (r *PushNotificationRouter) sendPush(ctx context.Context, userID, title, message string) {
	tokens, err := r.pushTokens.FindByUserID(ctx, userID)
	if err != nil || len(tokens) == 0 {
		return
	}
	for _, t := range tokens {
		if err := r.push.SendPush(ctx, t.EndpointArn, title, message, nil); err != nil {
			slog.ErrorContext(ctx, "failed to send push notification",
				"error", err,
				"endpoint", t.EndpointArn,
				"user_id", userID,
			)
		}
	}
}
