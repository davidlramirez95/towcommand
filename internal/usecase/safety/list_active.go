package safetyuc

import (
	"context"

	"github.com/davidlramirez95/towcommand/internal/domain/safety"
)

// defaultActiveLimit is the default number of active alerts to return.
const defaultActiveLimit int32 = 50

// ListActiveSOSUseCase retrieves all active (unresolved) SOS alerts.
type ListActiveSOSUseCase struct {
	lister SOSActiveLister
}

// NewListActiveSOSUseCase constructs a ListActiveSOSUseCase with its dependencies.
func NewListActiveSOSUseCase(lister SOSActiveLister) *ListActiveSOSUseCase {
	return &ListActiveSOSUseCase{lister: lister}
}

// Execute returns active SOS alerts up to the given limit. If limit is 0 the
// default of 50 is used.
func (uc *ListActiveSOSUseCase) Execute(ctx context.Context, limit int32) ([]safety.SOSAlert, error) {
	if limit <= 0 {
		limit = defaultActiveLimit
	}
	return uc.lister.FindActive(ctx, limit)
}
