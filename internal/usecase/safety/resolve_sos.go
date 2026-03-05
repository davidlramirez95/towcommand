package safetyuc

import (
	"context"
	"fmt"
	"time"

	domainerrors "github.com/davidlramirez95/towcommand/internal/domain/errors"
	"github.com/davidlramirez95/towcommand/internal/domain/event"
	"github.com/davidlramirez95/towcommand/internal/domain/safety"
)

// ResolveSOSInput holds the data needed to resolve an SOS alert.
type ResolveSOSInput struct {
	AlertID    string `json:"alertId" validate:"required"`
	ResolvedBy string `json:"resolvedBy" validate:"required"`
}

// ResolveSOSUseCase orchestrates the resolution of an SOS alert.
type ResolveSOSUseCase struct {
	finder   SOSFinder
	resolver SOSResolver
	events   EventPublisher
	now      func() time.Time
}

// NewResolveSOSUseCase constructs a ResolveSOSUseCase with its dependencies.
func NewResolveSOSUseCase(
	finder SOSFinder,
	resolver SOSResolver,
	events EventPublisher,
) *ResolveSOSUseCase {
	return &ResolveSOSUseCase{
		finder:   finder,
		resolver: resolver,
		events:   events,
		now:      func() time.Time { return time.Now().UTC() },
	}
}

// Execute resolves an SOS alert: verifies it exists and is not already resolved,
// marks it resolved in the repository, and publishes an SOSResolved event.
func (uc *ResolveSOSUseCase) Execute(ctx context.Context, input *ResolveSOSInput) (*safety.SOSAlert, error) {
	if input.AlertID == "" {
		return nil, domainerrors.NewValidationError("alertId is required")
	}
	if input.ResolvedBy == "" {
		return nil, domainerrors.NewValidationError("resolvedBy is required")
	}

	alert, err := uc.finder.FindByID(ctx, input.AlertID)
	if err != nil {
		return nil, fmt.Errorf("finding SOS alert %s: %w", input.AlertID, err)
	}
	if alert == nil {
		return nil, domainerrors.NewNotFoundError("SOS alert", input.AlertID)
	}

	if alert.Resolved {
		return nil, domainerrors.NewConflictError("SOS alert is already resolved")
	}

	resolvedAt := uc.now()

	if err := uc.resolver.Resolve(ctx, input.AlertID, input.ResolvedBy, resolvedAt); err != nil {
		return nil, fmt.Errorf("resolving SOS alert %s: %w", input.AlertID, err)
	}

	// Return the updated alert.
	alert.Resolved = true
	alert.ResolvedBy = input.ResolvedBy
	alert.ResolvedAt = &resolvedAt

	_ = uc.events.Publish(ctx, event.SourceSOS, event.SOSResolved, alert, &Actor{
		UserID:   input.ResolvedBy,
		UserType: "ops_agent",
	})

	return alert, nil
}
