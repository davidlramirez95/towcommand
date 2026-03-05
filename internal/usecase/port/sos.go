package port

import (
	"context"
	"time"

	"github.com/davidlramirez95/towcommand/internal/domain/safety"
)

// SOSSaver persists a new SOS alert.
type SOSSaver interface {
	Save(ctx context.Context, alert *safety.SOSAlert) error
}

// SOSFinder retrieves an SOS alert by its ID.
type SOSFinder interface {
	FindByID(ctx context.Context, alertID string) (*safety.SOSAlert, error)
}

// SOSResolver marks an SOS alert as resolved.
type SOSResolver interface {
	Resolve(ctx context.Context, alertID, resolvedBy string, resolvedAt time.Time) error
}

// SOSActiveLister queries for active (unresolved) SOS alerts.
type SOSActiveLister interface {
	FindActive(ctx context.Context, limit int32) ([]safety.SOSAlert, error)
}
