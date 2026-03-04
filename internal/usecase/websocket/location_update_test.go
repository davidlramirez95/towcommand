package websocket

import (
	"context"
	"log/slog"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	domainerrors "github.com/davidlramirez95/towcommand/internal/domain/errors"
	"github.com/davidlramirez95/towcommand/internal/usecase/port"
)

// --- Mocks ---

type mockGeoUpdater struct{ mock.Mock }

func (m *mockGeoUpdater) AddProviderLocation(ctx context.Context, providerID string, lat, lng float64) error {
	args := m.Called(ctx, providerID, lat, lng)
	return args.Error(0)
}

type mockEventPublisher struct{ mock.Mock }

func (m *mockEventPublisher) Publish(ctx context.Context, source, detailType string, detail any, actor *port.Actor) error {
	args := m.Called(ctx, source, detailType, detail, actor)
	return args.Error(0)
}

// --- Tests ---

func TestLocationUpdateUseCase_Execute_Success(t *testing.T) {
	geo := new(mockGeoUpdater)
	pub := new(mockEventPublisher)
	logger := slog.Default()
	uc := NewLocationUpdateUseCase(geo, pub, logger)

	input := LocationUpdateInput{
		ProviderID: "prov-1",
		Lat:        14.5995,
		Lng:        120.9842,
		Heading:    90.0,
		Speed:      30.0,
	}

	geo.On("AddProviderLocation", mock.Anything, "prov-1", 14.5995, 120.9842).Return(nil)
	pub.On("Publish", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)

	err := uc.Execute(context.Background(), input)

	require.NoError(t, err)
	geo.AssertExpectations(t)
	pub.AssertExpectations(t)
}

func TestLocationUpdateUseCase_Execute_GeoError(t *testing.T) {
	geo := new(mockGeoUpdater)
	pub := new(mockEventPublisher)
	logger := slog.Default()
	uc := NewLocationUpdateUseCase(geo, pub, logger)

	input := LocationUpdateInput{
		ProviderID: "prov-1",
		Lat:        14.5995,
		Lng:        120.9842,
	}

	geo.On("AddProviderLocation", mock.Anything, "prov-1", 14.5995, 120.9842).
		Return(domainerrors.NewExternalServiceError("Redis", nil))

	err := uc.Execute(context.Background(), input)

	assert.Error(t, err)
	pub.AssertNotCalled(t, "Publish")
}

func TestLocationUpdateUseCase_Execute_PublishError_NonFatal(t *testing.T) {
	geo := new(mockGeoUpdater)
	pub := new(mockEventPublisher)
	logger := slog.Default()
	uc := NewLocationUpdateUseCase(geo, pub, logger)

	input := LocationUpdateInput{
		ProviderID: "prov-1",
		Lat:        14.5995,
		Lng:        120.9842,
	}

	geo.On("AddProviderLocation", mock.Anything, "prov-1", 14.5995, 120.9842).Return(nil)
	pub.On("Publish", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
		Return(domainerrors.NewExternalServiceError("EventBridge", nil))

	// Publish error should not fail the use case.
	err := uc.Execute(context.Background(), input)

	require.NoError(t, err)
	geo.AssertExpectations(t)
	pub.AssertExpectations(t)
}
