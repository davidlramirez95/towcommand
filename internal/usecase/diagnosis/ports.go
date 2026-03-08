// Package diagnosisuc implements the AI-powered vehicle diagnosis use case
// following CLEAN architecture. Each use case declares only the port
// interfaces it needs (ISP).
package diagnosisuc

import "context"

// DiagnosisInput holds the customer's vehicle issue description and optional
// context used to produce an AI diagnosis.
type DiagnosisInput struct {
	Description string   `json:"description" validate:"required,min=10,max=1000"`
	PhotoURLs   []string `json:"photoUrls"`
	VehicleType string   `json:"vehicleType"`
	Location    *LatLng  `json:"location"`
}

// LatLng represents a geographic coordinate.
type LatLng struct {
	Lat float64 `json:"lat"`
	Lng float64 `json:"lng"`
}

// DiagnosisResult contains the AI-generated service recommendation, urgency
// assessment, cost estimate, and any safety warnings.
type DiagnosisResult struct {
	RecommendedService string   `json:"recommendedService"`
	UrgencyLevel       string   `json:"urgencyLevel"`
	EstimatedCostMin   int64    `json:"estimatedCostMin"`
	EstimatedCostMax   int64    `json:"estimatedCostMax"`
	Description        string   `json:"description"`
	SafetyWarnings     []string `json:"safetyWarnings"`
}

// DiagnosisEngine analyses vehicle issue descriptions (and optionally photos)
// to produce a structured diagnosis with service recommendation.
type DiagnosisEngine interface {
	Diagnose(ctx context.Context, input *DiagnosisInput) (*DiagnosisResult, error)
}
