package bookinguc

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/davidlramirez95/towcommand/internal/domain/booking"
	domainerrors "github.com/davidlramirez95/towcommand/internal/domain/errors"
)

// --- Mocks ---

type mockBookingSaver struct{ mock.Mock }

func (m *mockBookingSaver) Save(ctx context.Context, b *booking.Booking) error {
	args := m.Called(ctx, b)
	return args.Error(0)
}

type mockEventPublisher struct{ mock.Mock }

func (m *mockEventPublisher) Publish(ctx context.Context, source, detailType string, detail any, actor *Actor) error {
	args := m.Called(ctx, source, detailType, detail, actor)
	return args.Error(0)
}

// --- Tests ---

func TestCreateBookingUseCase_Execute_Success(t *testing.T) {
	repoMock := new(mockBookingSaver)
	eventsMock := new(mockEventPublisher)

	uc := NewCreateBookingUseCase(repoMock, eventsMock)
	uc.idGen = func() string { return "TC-2026-TESTID" }
	fixedTime := time.Date(2026, 3, 4, 10, 0, 0, 0, time.UTC) // daytime
	uc.now = func() time.Time { return fixedTime }

	input := CreateBookingInput{
		CustomerID:      "user-123",
		VehicleID:       "veh-456",
		ServiceType:     booking.ServiceTypeFlatbedTow,
		PickupLocation:  booking.GeoLocation{Lat: 14.5995, Lng: 120.9842, Address: "Manila"},
		DropoffLocation: booking.GeoLocation{Lat: 14.5547, Lng: 121.0244, Address: "Makati"},
		EstimateID:      "est-789",
		Notes:           "careful with the bumper",
	}

	repoMock.On("Save", mock.Anything, mock.AnythingOfType("*booking.Booking")).Return(nil)
	eventsMock.On("Publish", mock.Anything, eventSourceBooking, eventBookingCreated, mock.Anything, mock.Anything).Return(nil)

	result, err := uc.Execute(context.Background(), &input)

	require.NoError(t, err)
	assert.Equal(t, "TC-2026-TESTID", result.BookingID)
	assert.Equal(t, "user-123", result.CustomerID)
	assert.Equal(t, "veh-456", result.VehicleID)
	assert.Equal(t, booking.ServiceTypeFlatbedTow, result.ServiceType)
	assert.Equal(t, booking.BookingStatusPending, result.Status)
	assert.Equal(t, "est-789", result.EstimateID)
	assert.Equal(t, "careful with the bumper", result.Notes)
	assert.Equal(t, fixedTime, result.CreatedAt)
	assert.Equal(t, fixedTime, result.UpdatedAt)
	assert.Equal(t, "PHP", result.Price.Currency)
	assert.Greater(t, result.Price.Total, int64(0))

	repoMock.AssertExpectations(t)
	eventsMock.AssertExpectations(t)
}

func TestCreateBookingUseCase_Execute_NightSurcharge(t *testing.T) {
	repoMock := new(mockBookingSaver)
	eventsMock := new(mockEventPublisher)

	uc := NewCreateBookingUseCase(repoMock, eventsMock)
	uc.idGen = func() string { return "TC-2026-NIGHT" }
	// 23:00 PHT = 15:00 UTC
	nightTime := time.Date(2026, 3, 4, 15, 0, 0, 0, time.UTC)
	uc.now = func() time.Time { return nightTime }

	input := CreateBookingInput{
		CustomerID:      "user-123",
		VehicleID:       "veh-456",
		ServiceType:     booking.ServiceTypeFlatbedTow,
		PickupLocation:  booking.GeoLocation{Lat: 14.5995, Lng: 120.9842},
		DropoffLocation: booking.GeoLocation{Lat: 14.5547, Lng: 121.0244},
		EstimateID:      "est-789",
	}

	repoMock.On("Save", mock.Anything, mock.Anything).Return(nil)
	eventsMock.On("Publish", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)

	result, err := uc.Execute(context.Background(), &input)

	require.NoError(t, err)
	assert.Greater(t, result.Price.TimeSurcharge, int64(0))
}

func TestCreateBookingUseCase_Execute_RepoError(t *testing.T) {
	repoMock := new(mockBookingSaver)
	eventsMock := new(mockEventPublisher)

	uc := NewCreateBookingUseCase(repoMock, eventsMock)
	uc.idGen = func() string { return "TC-2026-ERR" }
	uc.now = func() time.Time { return time.Now().UTC() }

	repoMock.On("Save", mock.Anything, mock.Anything).Return(domainerrors.NewInternalError("db error"))

	_, err := uc.Execute(context.Background(), &CreateBookingInput{
		CustomerID:      "user-123",
		VehicleID:       "veh-456",
		ServiceType:     booking.ServiceTypeJumpstart,
		PickupLocation:  booking.GeoLocation{Lat: 14.5, Lng: 121.0},
		DropoffLocation: booking.GeoLocation{Lat: 14.6, Lng: 121.1},
		EstimateID:      "est-001",
	})

	assert.Error(t, err)
	repoMock.AssertExpectations(t)
	eventsMock.AssertNotCalled(t, "Publish")
}

func TestCreateBookingUseCase_Execute_EventPublishErrorDoesNotFail(t *testing.T) {
	repoMock := new(mockBookingSaver)
	eventsMock := new(mockEventPublisher)

	uc := NewCreateBookingUseCase(repoMock, eventsMock)
	uc.idGen = func() string { return "TC-2026-EVT" }
	uc.now = func() time.Time { return time.Now().UTC() }

	repoMock.On("Save", mock.Anything, mock.Anything).Return(nil)
	eventsMock.On("Publish", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
		Return(domainerrors.NewExternalServiceError("EventBridge", nil))

	result, err := uc.Execute(context.Background(), &CreateBookingInput{
		CustomerID:      "user-123",
		VehicleID:       "veh-456",
		ServiceType:     booking.ServiceTypeLockout,
		PickupLocation:  booking.GeoLocation{Lat: 14.5, Lng: 121.0},
		DropoffLocation: booking.GeoLocation{Lat: 14.6, Lng: 121.1},
		EstimateID:      "est-002",
	})

	require.NoError(t, err)
	assert.NotEmpty(t, result.BookingID)
}

func TestHaversineDistance(t *testing.T) {
	// Manila to Makati: approximately 6.5 km
	d := haversineDistance(14.5995, 120.9842, 14.5547, 121.0244)
	assert.InDelta(t, 6.0, d, 2.0) // within 2km tolerance
}

func TestIsNightTimePHT(t *testing.T) {
	tests := []struct {
		name string
		utc  time.Time
		want bool
	}{
		{"10am PHT (2am UTC)", time.Date(2026, 1, 1, 2, 0, 0, 0, time.UTC), false},
		{"11pm PHT (3pm UTC)", time.Date(2026, 1, 1, 15, 0, 0, 0, time.UTC), true},
		{"5am PHT (9pm UTC prev)", time.Date(2025, 12, 31, 21, 0, 0, 0, time.UTC), true},
		{"7am PHT (11pm UTC prev)", time.Date(2025, 12, 31, 23, 0, 0, 0, time.UTC), false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, isNightTimePHT(tt.utc))
		})
	}
}

func TestCalculatePrice_AllServiceTypes(t *testing.T) {
	serviceTypes := []booking.ServiceType{
		booking.ServiceTypeFlatbedTow,
		booking.ServiceTypeWheelLift,
		booking.ServiceTypeJumpstart,
		booking.ServiceTypeTireChange,
		booking.ServiceTypeFuelDelivery,
		booking.ServiceTypeLockout,
		booking.ServiceTypeAccidentRecovery,
	}
	for _, st := range serviceTypes {
		t.Run(string(st), func(t *testing.T) {
			p := calculatePrice(st, "light", 10.0, false)
			assert.Greater(t, p.Total, int64(0))
			assert.Equal(t, "PHP", p.Currency)
			assert.Equal(t, p.Base+p.Distance+p.Weight+p.TimeSurcharge+p.SurgePricing, p.Total)
		})
	}
}
