package analytics

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/davidlramirez95/towcommand/internal/domain/event"
)

// --- Mocks ---

type mockAnalyticsRecorder struct{ mock.Mock }

func (m *mockAnalyticsRecorder) IncrementDailyCounter(ctx context.Context, date, field string, delta int64) error {
	args := m.Called(ctx, date, field, delta)
	return args.Error(0)
}

func (m *mockAnalyticsRecorder) IncrementHeatmapCell(ctx context.Context, date, geohash string, lat, lng float64) error {
	args := m.Called(ctx, date, geohash, lat, lng)
	return args.Error(0)
}

func (m *mockAnalyticsRecorder) IncrementProviderCounter(ctx context.Context, providerID, date, field string, delta int64) error {
	args := m.Called(ctx, providerID, date, field, delta)
	return args.Error(0)
}

func mustJSONRaw(t *testing.T, v any) json.RawMessage {
	t.Helper()
	b, err := json.Marshal(v)
	require.NoError(t, err)
	return b
}

var testTime = time.Date(2026, 3, 4, 12, 0, 0, 0, time.UTC)

// --- Tests ---

func TestRecord_BookingCreated(t *testing.T) {
	repo := new(mockAnalyticsRecorder)
	recorder := NewEventRecorder(repo)

	repo.On("IncrementDailyCounter", mock.Anything, "2026-03-04", "totalBookings", int64(1)).Return(nil)
	repo.On("IncrementHeatmapCell", mock.Anything, "2026-03-04", "14.600,120.984", 14.5995, 120.9842).Return(nil)

	detail := mustJSONRaw(t, bookingCreatedDetail{
		BookingID: "BK-001",
		PickupLat: 14.5995,
		PickupLng: 120.9842,
	})

	err := recorder.Record(context.Background(), event.BookingCreated, detail, testTime)
	require.NoError(t, err)
	repo.AssertExpectations(t)
}

func TestRecord_BookingCompleted(t *testing.T) {
	repo := new(mockAnalyticsRecorder)
	recorder := NewEventRecorder(repo)

	repo.On("IncrementDailyCounter", mock.Anything, "2026-03-04", "completedBookings", int64(1)).Return(nil)
	repo.On("IncrementProviderCounter", mock.Anything, "prov-1", "2026-03-04", "completedJobs", int64(1)).Return(nil)

	detail := mustJSONRaw(t, bookingCompletedAnalyticsDetail{
		BookingID:  "BK-002",
		ProviderID: "prov-1",
	})

	err := recorder.Record(context.Background(), event.BookingCompleted, detail, testTime)
	require.NoError(t, err)
	repo.AssertExpectations(t)
}

func TestRecord_BookingCompleted_NoProvider(t *testing.T) {
	repo := new(mockAnalyticsRecorder)
	recorder := NewEventRecorder(repo)

	repo.On("IncrementDailyCounter", mock.Anything, "2026-03-04", "completedBookings", int64(1)).Return(nil)

	detail := mustJSONRaw(t, bookingCompletedAnalyticsDetail{
		BookingID:  "BK-003",
		ProviderID: "",
	})

	err := recorder.Record(context.Background(), event.BookingCompleted, detail, testTime)
	require.NoError(t, err)
	repo.AssertExpectations(t)
	repo.AssertNotCalled(t, "IncrementProviderCounter")
}

func TestRecord_BookingCancelled(t *testing.T) {
	repo := new(mockAnalyticsRecorder)
	recorder := NewEventRecorder(repo)

	repo.On("IncrementDailyCounter", mock.Anything, "2026-03-04", "cancelledBookings", int64(1)).Return(nil)
	repo.On("IncrementProviderCounter", mock.Anything, "prov-2", "2026-03-04", "cancelledJobs", int64(1)).Return(nil)

	detail := mustJSONRaw(t, bookingCancelledAnalyticsDetail{
		BookingID:  "BK-004",
		ProviderID: "prov-2",
	})

	err := recorder.Record(context.Background(), event.BookingCancelled, detail, testTime)
	require.NoError(t, err)
	repo.AssertExpectations(t)
}

func TestRecord_PaymentCaptured(t *testing.T) {
	repo := new(mockAnalyticsRecorder)
	recorder := NewEventRecorder(repo)

	repo.On("IncrementDailyCounter", mock.Anything, "2026-03-04", "totalRevenueCentavos", int64(250000)).Return(nil)
	repo.On("IncrementProviderCounter", mock.Anything, "prov-3", "2026-03-04", "totalRevenueCentavos", int64(250000)).Return(nil)

	detail := mustJSONRaw(t, paymentCapturedAnalyticsDetail{
		BookingID:      "BK-005",
		ProviderID:     "prov-3",
		AmountCentavos: 250000,
	})

	err := recorder.Record(context.Background(), event.PaymentCaptured, detail, testTime)
	require.NoError(t, err)
	repo.AssertExpectations(t)
}

func TestRecord_UnhandledEvent(t *testing.T) {
	repo := new(mockAnalyticsRecorder)
	recorder := NewEventRecorder(repo)

	err := recorder.Record(context.Background(), "SomeOtherEvent", json.RawMessage(`{}`), testTime)
	assert.NoError(t, err)
	repo.AssertNotCalled(t, "IncrementDailyCounter")
}

func TestRecord_InvalidJSON(t *testing.T) {
	repo := new(mockAnalyticsRecorder)
	recorder := NewEventRecorder(repo)

	err := recorder.Record(context.Background(), event.BookingCreated, json.RawMessage(`{invalid`), testTime)
	assert.Error(t, err)
}

func TestRecord_DailyCounterError(t *testing.T) {
	repo := new(mockAnalyticsRecorder)
	recorder := NewEventRecorder(repo)

	repo.On("IncrementDailyCounter", mock.Anything, "2026-03-04", "totalBookings", int64(1)).
		Return(assert.AnError)

	detail := mustJSONRaw(t, bookingCreatedDetail{
		BookingID: "BK-ERR",
		PickupLat: 14.5,
		PickupLng: 121.0,
	})

	err := recorder.Record(context.Background(), event.BookingCreated, detail, testTime)
	assert.Error(t, err)
}

func TestRecord_BookingCancelled_NoProvider(t *testing.T) {
	repo := new(mockAnalyticsRecorder)
	recorder := NewEventRecorder(repo)

	repo.On("IncrementDailyCounter", mock.Anything, "2026-03-04", "cancelledBookings", int64(1)).Return(nil)

	detail := mustJSONRaw(t, bookingCancelledAnalyticsDetail{
		BookingID:  "BK-NOPROV",
		ProviderID: "",
	})

	err := recorder.Record(context.Background(), event.BookingCancelled, detail, testTime)
	require.NoError(t, err)
	repo.AssertExpectations(t)
	repo.AssertNotCalled(t, "IncrementProviderCounter")
}

func TestRecord_PaymentCaptured_NoProvider(t *testing.T) {
	repo := new(mockAnalyticsRecorder)
	recorder := NewEventRecorder(repo)

	repo.On("IncrementDailyCounter", mock.Anything, "2026-03-04", "totalRevenueCentavos", int64(100000)).Return(nil)

	detail := mustJSONRaw(t, paymentCapturedAnalyticsDetail{
		BookingID:      "BK-NOPROV-PAY",
		ProviderID:     "",
		AmountCentavos: 100000,
	})

	err := recorder.Record(context.Background(), event.PaymentCaptured, detail, testTime)
	require.NoError(t, err)
	repo.AssertExpectations(t)
	repo.AssertNotCalled(t, "IncrementProviderCounter")
}

func TestRecord_BookingCompleted_DailyCounterError(t *testing.T) {
	repo := new(mockAnalyticsRecorder)
	recorder := NewEventRecorder(repo)

	repo.On("IncrementDailyCounter", mock.Anything, "2026-03-04", "completedBookings", int64(1)).
		Return(assert.AnError)

	detail := mustJSONRaw(t, bookingCompletedAnalyticsDetail{
		BookingID:  "BK-COMP-ERR",
		ProviderID: "prov-1",
	})

	err := recorder.Record(context.Background(), event.BookingCompleted, detail, testTime)
	assert.Error(t, err)
}

func TestRecord_BookingCancelled_InvalidJSON(t *testing.T) {
	repo := new(mockAnalyticsRecorder)
	recorder := NewEventRecorder(repo)

	err := recorder.Record(context.Background(), event.BookingCancelled, json.RawMessage(`{bad`), testTime)
	assert.Error(t, err)
}

func TestRecord_PaymentCaptured_InvalidJSON(t *testing.T) {
	repo := new(mockAnalyticsRecorder)
	recorder := NewEventRecorder(repo)

	err := recorder.Record(context.Background(), event.PaymentCaptured, json.RawMessage(`{bad`), testTime)
	assert.Error(t, err)
}

func TestRecord_BookingCompleted_InvalidJSON(t *testing.T) {
	repo := new(mockAnalyticsRecorder)
	recorder := NewEventRecorder(repo)

	err := recorder.Record(context.Background(), event.BookingCompleted, json.RawMessage(`{bad`), testTime)
	assert.Error(t, err)
}

func TestRecord_HeatmapCellError(t *testing.T) {
	repo := new(mockAnalyticsRecorder)
	recorder := NewEventRecorder(repo)

	repo.On("IncrementDailyCounter", mock.Anything, "2026-03-04", "totalBookings", int64(1)).Return(nil)
	repo.On("IncrementHeatmapCell", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
		Return(assert.AnError)

	detail := mustJSONRaw(t, bookingCreatedDetail{
		BookingID: "BK-HEAT-ERR",
		PickupLat: 14.5,
		PickupLng: 121.0,
	})

	err := recorder.Record(context.Background(), event.BookingCreated, detail, testTime)
	assert.Error(t, err)
}

func TestRecord_ProviderCounterError(t *testing.T) {
	repo := new(mockAnalyticsRecorder)
	recorder := NewEventRecorder(repo)

	repo.On("IncrementDailyCounter", mock.Anything, "2026-03-04", "completedBookings", int64(1)).Return(nil)
	repo.On("IncrementProviderCounter", mock.Anything, "prov-err", "2026-03-04", "completedJobs", int64(1)).
		Return(assert.AnError)

	detail := mustJSONRaw(t, bookingCompletedAnalyticsDetail{
		BookingID:  "BK-PROV-ERR",
		ProviderID: "prov-err",
	})

	err := recorder.Record(context.Background(), event.BookingCompleted, detail, testTime)
	assert.Error(t, err)
}

func TestGeohash6(t *testing.T) {
	tests := []struct {
		name string
		lat  float64
		lng  float64
		want string
	}{
		{"Manila", 14.5995, 120.9842, "14.600,120.984"},
		{"exact values", 14.500, 121.000, "14.500,121.000"},
		{"negative coords", -6.2088, 106.8456, "-6.209,106.846"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Geohash6(tt.lat, tt.lng)
			assert.Equal(t, tt.want, got)
		})
	}
}
