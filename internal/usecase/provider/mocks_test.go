package provider_test

import (
	"context"

	"github.com/davidlramirez95/towcommand/internal/domain/provider"
	"github.com/davidlramirez95/towcommand/internal/usecase/port"
)

// --- ProviderSaver ---

type mockProviderSaver struct {
	saved *provider.Provider
	err   error
}

func (m *mockProviderSaver) Save(_ context.Context, p *provider.Provider) error {
	if m.err != nil {
		return m.err
	}
	m.saved = p
	return nil
}

// --- ProviderFinder (single provider) ---

type mockProviderFinder struct {
	provider *provider.Provider
	err      error
}

func (m *mockProviderFinder) FindByID(_ context.Context, _ string) (*provider.Provider, error) {
	return m.provider, m.err
}

// --- ProviderFinder (map-based, for get-nearby) ---

type mockProviderFinderMap struct {
	providers map[string]*provider.Provider
	err       error
}

func (m *mockProviderFinderMap) FindByID(_ context.Context, providerID string) (*provider.Provider, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.providers[providerID], nil
}

// --- ProviderLocationUpdater ---

type mockProviderLocationUpdater struct {
	err error
}

func (m *mockProviderLocationUpdater) UpdateLocation(_ context.Context, _ string, _, _ float64) error {
	return m.err
}

// --- ProviderAvailabilityUpdater ---

type mockProviderAvailabilityUpdater struct {
	err error
}

func (m *mockProviderAvailabilityUpdater) UpdateAvailability(_ context.Context, _ string, _ bool) error {
	return m.err
}

// --- GeoCache ---

type mockGeoCache struct {
	addErr      error
	findErr     error
	removeErr   error
	nearby      []port.ProviderDistance
	addCalled   bool
	removeCalled bool
}

func (m *mockGeoCache) AddProviderLocation(_ context.Context, _ string, _, _ float64) error {
	m.addCalled = true
	return m.addErr
}

func (m *mockGeoCache) FindNearbyProviders(_ context.Context, _, _, _ float64) ([]port.ProviderDistance, error) {
	return m.nearby, m.findErr
}

func (m *mockGeoCache) RemoveProvider(_ context.Context, _ string) error {
	m.removeCalled = true
	return m.removeErr
}

// --- EventPublisher ---

type mockEventPublisher struct {
	err    error
	called bool
}

func (m *mockEventPublisher) Publish(_ context.Context, _, _ string, _ any, _ *port.Actor) error {
	m.called = true
	return m.err
}
