package notification

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBookingMatchedSMS(t *testing.T) {
	msg := BookingMatchedSMS("TC-2026-ABC", "Juan dela Cruz", 12)
	assert.Contains(t, msg, "TC-2026-ABC")
	assert.Contains(t, msg, "Juan dela Cruz")
	assert.Contains(t, msg, "12 min")
	assert.Contains(t, msg, "Mabuhay")
}

func TestBookingCancelledSMS(t *testing.T) {
	msg := BookingCancelledSMS("TC-2026-XYZ")
	assert.Contains(t, msg, "TC-2026-XYZ")
	assert.Contains(t, msg, "na-cancel")
}

func TestBookingCompletedSMS(t *testing.T) {
	msg := BookingCompletedSMS("TC-2026-DONE")
	assert.Contains(t, msg, "TC-2026-DONE")
	assert.Contains(t, msg, "Salamat")
	assert.Contains(t, msg, "tapos na")
}

func TestBookingStatusChangedSMS(t *testing.T) {
	tests := []struct {
		name     string
		status   string
		contains string
	}{
		{"EN_ROUTE", "EN_ROUTE", "Paparating"},
		{"ARRIVED", "ARRIVED", "Nandito na"},
		{"LOADING", "LOADING", "Nilo-load"},
		{"IN_TRANSIT", "IN_TRANSIT", "Nasa daan"},
		{"unknown status", "SOME_NEW_STATUS", "SOME_NEW_STATUS"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			msg := BookingStatusChangedSMS("BK-001", tt.status)
			assert.Contains(t, msg, "BK-001")
			assert.Contains(t, msg, tt.contains)
		})
	}
}

func TestSOSAlertSMS(t *testing.T) {
	msg := SOSAlertSMS("TC-2026-SOS", 14.5995, 120.9842)
	assert.Contains(t, msg, "SOS ALERT")
	assert.Contains(t, msg, "TC-2026-SOS")
	assert.Contains(t, msg, "14.599500")
	assert.Contains(t, msg, "120.984200")
}

func TestSOSAlertEmail(t *testing.T) {
	subject := SOSAlertEmailSubject("TC-2026-SOS")
	assert.Contains(t, subject, "URGENT")
	assert.Contains(t, subject, "TC-2026-SOS")

	body := SOSAlertEmailBody("TC-2026-SOS", 14.5995, 120.9842, "HIGH")
	assert.Contains(t, body, "TC-2026-SOS")
	assert.Contains(t, body, "HIGH")
	assert.Contains(t, body, "14.599500")
}

func TestPaymentCapturedSMS(t *testing.T) {
	msg := PaymentCapturedSMS(150000)
	assert.Contains(t, msg, "PHP 1500.00")
	assert.Contains(t, msg, "Natanggap")
}

func TestWelcomeSMS(t *testing.T) {
	msg := WelcomeSMS()
	assert.Contains(t, msg, "Welcome sa TowCommand")
	assert.Contains(t, msg, "Mabuhay")
}
