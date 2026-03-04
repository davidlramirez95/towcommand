package notification

import "fmt"

// BookingMatchedSMS returns the SMS text sent to a customer when their booking is matched.
func BookingMatchedSMS(bookingID, providerName string, eta int) string {
	return fmt.Sprintf(
		"Mabuhay! Ang iyong tow truck ay paparating na. Driver: %s, ETA: %d min. Booking #%s",
		providerName, eta, bookingID,
	)
}

// BookingCancelledSMS returns the SMS text sent when a booking is cancelled.
func BookingCancelledSMS(bookingID string) string {
	return fmt.Sprintf("Ang booking #%s ay na-cancel.", bookingID)
}

// BookingCompletedSMS returns the SMS text sent to a customer when their booking completes.
func BookingCompletedSMS(bookingID string) string {
	return fmt.Sprintf("Salamat! Ang iyong booking #%s ay tapos na. Maraming salamat sa paggamit ng TowCommand!", bookingID)
}

// BookingStatusChangedSMS returns the SMS text for a generic booking status change.
func BookingStatusChangedSMS(bookingID, status string) string {
	msg, ok := statusMessages[status]
	if !ok {
		return fmt.Sprintf("Booking #%s: status updated to %s.", bookingID, status)
	}
	return fmt.Sprintf("Booking #%s: %s", bookingID, msg)
}

// statusMessages maps booking statuses to Filipino-friendly descriptions.
var statusMessages = map[string]string{
	"EN_ROUTE":         "Paparating na ang iyong tow truck!",
	"ARRIVED":          "Nandito na ang tow truck sa pickup location.",
	"CONDITION_REPORT": "Sinusuri ang kondisyon ng sasakyan.",
	"OTP_VERIFIED":     "Na-verify na ang OTP. Magsisimula na ang loading.",
	"LOADING":          "Nilo-load na ang iyong sasakyan.",
	"IN_TRANSIT":       "Nasa daan na ang iyong sasakyan papunta sa destinasyon.",
	"ARRIVED_DROPOFF":  "Dumating na sa drop-off location.",
	"OTP_DROPOFF":      "Kailangan ng OTP para sa drop-off verification.",
}

// SOSAlertSMS returns the SMS text sent to operations when an SOS is triggered.
func SOSAlertSMS(bookingID string, lat, lng float64) string {
	return fmt.Sprintf(
		"[SOS ALERT] Booking #%s needs immediate assistance! Location: %.6f, %.6f. Respond ASAP.",
		bookingID, lat, lng,
	)
}

// SOSAlertEmailSubject returns the email subject for an SOS alert.
func SOSAlertEmailSubject(bookingID string) string {
	return fmt.Sprintf("[URGENT] SOS Alert - Booking #%s", bookingID)
}

// SOSAlertEmailBody returns the HTML email body for an SOS alert.
func SOSAlertEmailBody(bookingID string, lat, lng float64, severity string) string {
	return fmt.Sprintf(
		`<h2>SOS Alert Triggered</h2>
<p><strong>Booking:</strong> %s</p>
<p><strong>Location:</strong> %.6f, %.6f</p>
<p><strong>Severity:</strong> %s</p>
<p>Please investigate immediately and coordinate with the operations team.</p>`,
		bookingID, lat, lng, severity,
	)
}

// PaymentCapturedSMS returns the SMS text sent when a payment is captured.
func PaymentCapturedSMS(amountCentavos int64) string {
	pesos := float64(amountCentavos) / 100.0
	return fmt.Sprintf("Natanggap ang bayad: PHP %.2f. Salamat sa TowCommand!", pesos)
}

// WelcomeSMS returns the SMS text sent when a new user registers.
func WelcomeSMS() string {
	return "Welcome sa TowCommand! Mabuhay! Ang Grab ng Towing, handang tumulong sa inyo."
}
