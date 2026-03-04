package websocket

import (
	"context"
	"log/slog"
	"time"

	"github.com/davidlramirez95/towcommand/internal/domain/event"
	"github.com/davidlramirez95/towcommand/internal/usecase/port"
)

// LocationUpdateInput carries the validated fields from a WebSocket location message.
type LocationUpdateInput struct {
	ProviderID string  `json:"providerId" validate:"required"`
	Lat        float64 `json:"lat" validate:"required"`
	Lng        float64 `json:"lng" validate:"required"`
	Heading    float64 `json:"heading" validate:"min=0,max=360"`
	Speed      float64 `json:"speed" validate:"min=0"`
}

// LocationUpdateUseCase handles real-time provider location updates via WebSocket.
type LocationUpdateUseCase struct {
	geo       GeoUpdater
	publisher EventPublisher
	logger    *slog.Logger
}

// NewLocationUpdateUseCase creates a new LocationUpdateUseCase.
func NewLocationUpdateUseCase(geo GeoUpdater, publisher EventPublisher, logger *slog.Logger) *LocationUpdateUseCase {
	return &LocationUpdateUseCase{geo: geo, publisher: publisher, logger: logger}
}

// Execute updates the provider's geospatial position and publishes a LocationUpdated event.
func (uc *LocationUpdateUseCase) Execute(ctx context.Context, input LocationUpdateInput) error {
	if err := uc.geo.AddProviderLocation(ctx, input.ProviderID, input.Lat, input.Lng); err != nil {
		return err
	}

	if err := uc.publisher.Publish(ctx, event.SourceTracking, event.LocationUpdated, map[string]any{
		"providerId": input.ProviderID,
		"lat":        input.Lat,
		"lng":        input.Lng,
		"heading":    input.Heading,
		"speed":      input.Speed,
		"timestamp":  time.Now().UTC().Format(time.RFC3339),
	}, &port.Actor{UserID: input.ProviderID, UserType: "provider"}); err != nil {
		uc.logger.WarnContext(ctx, "failed to publish LocationUpdated event",
			slog.String("providerId", input.ProviderID),
			slog.String("error", err.Error()),
		)
	}

	return nil
}
