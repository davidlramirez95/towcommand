package provider

import (
	"context"
	"log/slog"
	"time"

	"github.com/davidlramirez95/towcommand/internal/domain/errors"
	"github.com/davidlramirez95/towcommand/internal/domain/event"
	"github.com/davidlramirez95/towcommand/internal/usecase/port"
)

// UpdateLocationInput carries the validated fields for a location update.
type UpdateLocationInput struct {
	ProviderID string  `json:"providerId" validate:"required"`
	Lat        float64 `json:"lat" validate:"required"`
	Lng        float64 `json:"lng" validate:"required"`
	Heading    float64 `json:"heading" validate:"min=0,max=360"`
	Speed      float64 `json:"speed" validate:"min=0"`
}

// UpdateLocationOutput is the response after a successful location update.
type UpdateLocationOutput struct {
	ProviderID string  `json:"providerId"`
	Lat        float64 `json:"lat"`
	Lng        float64 `json:"lng"`
}

// UpdateLocationUseCase orchestrates provider location updates.
type UpdateLocationUseCase struct {
	repo      port.ProviderLocationUpdater
	geo       port.GeoCache
	publisher port.EventPublisher
	logger    *slog.Logger
}

// NewUpdateLocationUseCase creates a new UpdateLocationUseCase.
func NewUpdateLocationUseCase(
	repo port.ProviderLocationUpdater,
	geo port.GeoCache,
	publisher port.EventPublisher,
	logger *slog.Logger,
) *UpdateLocationUseCase {
	return &UpdateLocationUseCase{repo: repo, geo: geo, publisher: publisher, logger: logger}
}

// Execute updates the provider's location in both the geo cache and DynamoDB.
func (uc *UpdateLocationUseCase) Execute(ctx context.Context, input UpdateLocationInput) (*UpdateLocationOutput, error) {
	if !isValidPhilippineCoordinate(input.Lat, input.Lng) {
		return nil, errors.NewValidationError("coordinates must be within the Philippines")
	}

	if err := uc.geo.AddProviderLocation(ctx, input.ProviderID, input.Lat, input.Lng); err != nil {
		return nil, errors.NewExternalServiceError("Redis", err)
	}

	if err := uc.repo.UpdateLocation(ctx, input.ProviderID, input.Lat, input.Lng); err != nil {
		return nil, errors.NewInternalError("failed to update provider location").WithCause(err)
	}

	if err := uc.publisher.Publish(ctx, event.SourceTracking, event.LocationUpdated, map[string]any{
		"providerId": input.ProviderID,
		"lat":        input.Lat,
		"lng":        input.Lng,
		"heading":    input.Heading,
		"speed":      input.Speed,
		"timestamp":  time.Now().UTC().Format(time.RFC3339),
	}, &port.Actor{UserID: input.ProviderID, UserType: "provider"}); err != nil {
		uc.logger.WarnContext(ctx, "failed to publish LocationUpdated event", "error", err)
	}

	return &UpdateLocationOutput{
		ProviderID: input.ProviderID,
		Lat:        input.Lat,
		Lng:        input.Lng,
	}, nil
}

// isValidPhilippineCoordinate checks whether coordinates fall within the Philippines bounding box.
// Approximate bounds: lat 4.5°N–21.5°N, lng 116°E–127°E.
func isValidPhilippineCoordinate(lat, lng float64) bool {
	return lat >= 4.5 && lat <= 21.5 && lng >= 116.0 && lng <= 127.0
}
