package port

import (
	"context"

	"github.com/davidlramirez95/towcommand/internal/domain/provider"
	"github.com/davidlramirez95/towcommand/internal/domain/user"
)

// ProviderSaver persists a new provider.
type ProviderSaver interface {
	Save(ctx context.Context, p *provider.Provider) error
}

// ProviderFinder retrieves a provider by their ID.
type ProviderFinder interface {
	FindByID(ctx context.Context, providerID string) (*provider.Provider, error)
}

// ProviderByTierLister lists providers by trust tier and city via GSI3.
type ProviderByTierLister interface {
	FindByTierAndCity(ctx context.Context, tier user.TrustTier, city string, limit int32) ([]provider.Provider, error)
}

// ProviderLocationUpdater updates a provider's GPS coordinates.
type ProviderLocationUpdater interface {
	UpdateLocation(ctx context.Context, providerID string, lat, lng float64) error
}

// ProviderAvailabilityUpdater toggles a provider's online status.
type ProviderAvailabilityUpdater interface {
	UpdateAvailability(ctx context.Context, providerID string, isOnline bool) error
}

// ProviderDocSaver persists a provider KYC document.
type ProviderDocSaver interface {
	UploadDoc(ctx context.Context, doc *provider.ProviderDoc) error
}

// ProviderDocLister lists all documents for a provider.
type ProviderDocLister interface {
	GetDocs(ctx context.Context, providerID string) ([]provider.ProviderDoc, error)
}
