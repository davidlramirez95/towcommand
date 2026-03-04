package provider

import (
	"context"
	"log/slog"

	"github.com/davidlramirez95/towcommand/internal/domain/errors"
	"github.com/davidlramirez95/towcommand/internal/domain/event"
	provdomain "github.com/davidlramirez95/towcommand/internal/domain/provider"
	"github.com/davidlramirez95/towcommand/internal/usecase/port"
)

// ToggleAvailabilityInput carries the validated fields for toggling availability.
type ToggleAvailabilityInput struct {
	ProviderID string `json:"providerId" validate:"required"`
	Online     bool   `json:"online"`
}

// ToggleAvailabilityOutput is the response after toggling availability.
type ToggleAvailabilityOutput struct {
	ProviderID string `json:"providerId"`
	Online     bool   `json:"online"`
}

// ToggleAvailabilityUseCase orchestrates provider availability toggling.
type ToggleAvailabilityUseCase struct {
	finder    port.ProviderFinder
	updater   port.ProviderAvailabilityUpdater
	geo       port.GeoCache
	publisher port.EventPublisher
	logger    *slog.Logger
}

// NewToggleAvailabilityUseCase creates a new ToggleAvailabilityUseCase.
func NewToggleAvailabilityUseCase(
	finder port.ProviderFinder,
	updater port.ProviderAvailabilityUpdater,
	geo port.GeoCache,
	publisher port.EventPublisher,
	logger *slog.Logger,
) *ToggleAvailabilityUseCase {
	return &ToggleAvailabilityUseCase{
		finder: finder, updater: updater, geo: geo, publisher: publisher, logger: logger,
	}
}

// Execute toggles a provider's online status. Only active providers may go online.
func (uc *ToggleAvailabilityUseCase) Execute(ctx context.Context, input ToggleAvailabilityInput) (*ToggleAvailabilityOutput, error) {
	p, err := uc.finder.FindByID(ctx, input.ProviderID)
	if err != nil {
		return nil, errors.NewInternalError("failed to find provider").WithCause(err)
	}
	if p == nil {
		return nil, errors.NewNotFoundError("Provider", input.ProviderID)
	}

	if input.Online && p.Status != provdomain.ProviderStatusActive {
		return nil, errors.NewValidationError("only verified providers can go online")
	}

	if err := uc.updater.UpdateAvailability(ctx, input.ProviderID, input.Online); err != nil {
		return nil, errors.NewInternalError("failed to update availability").WithCause(err)
	}

	if input.Online {
		if p.CurrentLat != nil && p.CurrentLng != nil {
			if err := uc.geo.AddProviderLocation(ctx, input.ProviderID, *p.CurrentLat, *p.CurrentLng); err != nil {
				uc.logger.WarnContext(ctx, "failed to add provider to geo index", "error", err)
			}
		}

		if err := uc.publisher.Publish(ctx, event.SourceProvider, event.ProviderOnline, map[string]any{
			"providerId": input.ProviderID,
			"lat":        p.CurrentLat,
			"lng":        p.CurrentLng,
		}, &port.Actor{UserID: input.ProviderID, UserType: "provider"}); err != nil {
			uc.logger.WarnContext(ctx, "failed to publish ProviderOnline event", "error", err)
		}
	} else {
		if err := uc.geo.RemoveProvider(ctx, input.ProviderID); err != nil {
			uc.logger.WarnContext(ctx, "failed to remove provider from geo index", "error", err)
		}

		if err := uc.publisher.Publish(ctx, event.SourceProvider, event.ProviderOffline, map[string]any{
			"providerId": input.ProviderID,
		}, &port.Actor{UserID: input.ProviderID, UserType: "provider"}); err != nil {
			uc.logger.WarnContext(ctx, "failed to publish ProviderOffline event", "error", err)
		}
	}

	return &ToggleAvailabilityOutput{
		ProviderID: input.ProviderID,
		Online:     input.Online,
	}, nil
}
