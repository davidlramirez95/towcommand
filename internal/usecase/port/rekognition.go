package port

import "context"

// ImageValidationResult holds the outcome of an image validation check.
type ImageValidationResult struct {
	IsValid bool
	Labels  []string
	Reason  string
}

// ImageValidator validates images against domain-specific criteria.
type ImageValidator interface {
	ValidateVehiclePhoto(ctx context.Context, s3Bucket, s3Key string) (*ImageValidationResult, error)
}
