package diagnosisuc

import (
	"context"
	"fmt"

	domainerrors "github.com/davidlramirez95/towcommand/internal/domain/errors"
)

// validServiceTypes enumerates the service types the AI may recommend.
var validServiceTypes = map[string]bool{
	"FLATBED_TOWING":    true,
	"WHEEL_LIFT_TOWING": true,
	"MOTORCYCLE_TOWING": true,
	"JUMPSTART":         true,
	"TIRE_CHANGE":       true,
	"LOCKOUT":           true,
	"FUEL_DELIVERY":     true,
	"WINCH_RECOVERY":    true,
}

// validUrgencyLevels enumerates the urgency levels the AI may assign.
var validUrgencyLevels = map[string]bool{
	"LOW":      true,
	"MEDIUM":   true,
	"HIGH":     true,
	"CRITICAL": true,
}

// DiagnoseUseCase orchestrates an AI-powered vehicle issue diagnosis.
type DiagnoseUseCase struct {
	engine DiagnosisEngine
}

// NewDiagnoseUseCase constructs a DiagnoseUseCase with its dependencies.
func NewDiagnoseUseCase(engine DiagnosisEngine) *DiagnoseUseCase {
	return &DiagnoseUseCase{engine: engine}
}

// Execute validates the input, calls the AI engine, and sanitises the result.
func (uc *DiagnoseUseCase) Execute(ctx context.Context, input *DiagnosisInput) (*DiagnosisResult, error) {
	// 1. Validate description length.
	if len(input.Description) < 10 {
		return nil, domainerrors.NewValidationError("description must be at least 10 characters")
	}
	if len(input.Description) > 1000 {
		return nil, domainerrors.NewValidationError("description must not exceed 1000 characters")
	}

	// 2. Call the AI engine.
	result, err := uc.engine.Diagnose(ctx, input)
	if err != nil {
		return nil, domainerrors.NewExternalServiceError("bedrock", err)
	}

	// 3. Sanitise / validate AI output.
	if !validServiceTypes[result.RecommendedService] {
		result.RecommendedService = "FLATBED_TOWING" // safe default
	}
	if !validUrgencyLevels[result.UrgencyLevel] {
		result.UrgencyLevel = "MEDIUM" // safe default
	}
	if result.EstimatedCostMin < 0 {
		result.EstimatedCostMin = 0
	}
	if result.EstimatedCostMax < result.EstimatedCostMin {
		result.EstimatedCostMax = result.EstimatedCostMin
	}
	if result.Description == "" {
		result.Description = fmt.Sprintf("AI diagnosis for: %s", input.Description)
	}
	if result.SafetyWarnings == nil {
		result.SafetyWarnings = []string{}
	}

	return result, nil
}
