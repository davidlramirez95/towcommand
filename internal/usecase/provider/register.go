// Package provider contains use cases for provider operations.
package provider

import (
	"context"
	"crypto/rand"
	"fmt"
	"log/slog"
	"time"

	"github.com/davidlramirez95/towcommand/internal/domain/errors"
	"github.com/davidlramirez95/towcommand/internal/domain/event"
	"github.com/davidlramirez95/towcommand/internal/domain/provider"
	"github.com/davidlramirez95/towcommand/internal/domain/user"
	"github.com/davidlramirez95/towcommand/internal/usecase/port"
)

// RegisterInput carries the validated fields for provider registration.
type RegisterInput struct {
	CognitoSub          string   `json:"cognitoSub" validate:"required"`
	Name                string   `json:"name" validate:"required"`
	Phone               string   `json:"phone" validate:"required"`
	Email               string   `json:"email" validate:"required,email"`
	TruckType           string   `json:"truckType" validate:"required,oneof=flatbed wheel_lift boom motorcycle_carrier"`
	MaxWeightCapacityKg int      `json:"maxWeightCapacityKg" validate:"required,gt=0"`
	PlateNumber         string   `json:"plateNumber" validate:"required"`
	LTORegistration     string   `json:"ltoRegistration" validate:"required"`
	ServiceAreas        []string `json:"serviceAreas" validate:"required,min=1"`
}

// RegisterOutput is the response after successful registration.
type RegisterOutput struct {
	ProviderID string                  `json:"providerId"`
	Status     provider.ProviderStatus `json:"status"`
}

// RegisterUseCase orchestrates provider registration.
type RegisterUseCase struct {
	saver     port.ProviderSaver
	publisher port.EventPublisher
	logger    *slog.Logger
}

// NewRegisterUseCase creates a new RegisterUseCase.
func NewRegisterUseCase(saver port.ProviderSaver, publisher port.EventPublisher, logger *slog.Logger) *RegisterUseCase {
	return &RegisterUseCase{saver: saver, publisher: publisher, logger: logger}
}

// Execute registers a new provider with pending verification status and BASIC trust tier.
func (uc *RegisterUseCase) Execute(ctx context.Context, input RegisterInput) (*RegisterOutput, error) {
	providerID := generateProviderID()
	now := time.Now().UTC()

	p := &provider.Provider{
		ProviderID:          providerID,
		CognitoSub:          input.CognitoSub,
		Name:                input.Name,
		Phone:               input.Phone,
		Email:               input.Email,
		Status:              provider.ProviderStatusPendingVerification,
		TrustTier:           user.TrustTierBasic,
		TruckType:           provider.TruckType(input.TruckType),
		MaxWeightCapacityKg: input.MaxWeightCapacityKg,
		PlateNumber:         input.PlateNumber,
		LTORegistration:     input.LTORegistration,
		NBIClearanceStatus:  provider.ClearanceStatusPending,
		DrugTestStatus:      provider.ClearanceStatusPending,
		MMADAccredited:      false,
		Rating:              0,
		TotalJobsCompleted:  0,
		AcceptanceRate:      1.0,
		IsOnline:            false,
		ServiceAreas:        input.ServiceAreas,
		CreatedAt:           now,
		UpdatedAt:           now,
	}

	if err := uc.saver.Save(ctx, p); err != nil {
		return nil, errors.NewInternalError("failed to save provider").WithCause(err)
	}

	if err := uc.publisher.Publish(ctx, event.SourceProvider, event.ProviderRegistered, map[string]any{
		"providerId": providerID,
		"name":       input.Name,
		"phone":      input.Phone,
		"status":     string(provider.ProviderStatusPendingVerification),
	}, &port.Actor{UserID: providerID, UserType: "provider"}); err != nil {
		uc.logger.WarnContext(ctx, "failed to publish ProviderRegistered event", "error", err)
	}

	return &RegisterOutput{
		ProviderID: providerID,
		Status:     provider.ProviderStatusPendingVerification,
	}, nil
}

// generateProviderID produces a unique provider ID with PROV- prefix.
func generateProviderID() string {
	b := make([]byte, 16)
	_, _ = rand.Read(b)
	return fmt.Sprintf("PROV-%x", b)
}
