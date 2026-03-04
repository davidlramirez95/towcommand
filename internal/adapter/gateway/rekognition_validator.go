package gateway

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/rekognition"
	rektypes "github.com/aws/aws-sdk-go-v2/service/rekognition/types"

	"github.com/davidlramirez95/towcommand/internal/usecase/port"
)

// vehicleLabels is the set of Rekognition labels that indicate a vehicle is present.
var vehicleLabels = map[string]bool{
	"Car":            true,
	"Vehicle":        true,
	"Truck":          true,
	"Automobile":     true,
	"Transportation": true,
	"Van":            true,
	"Motorcycle":     true,
	"Wheel":          true,
}

// RekognitionAPI is the subset of the Rekognition client needed by the validator.
type RekognitionAPI interface {
	DetectLabels(ctx context.Context, params *rekognition.DetectLabelsInput, optFns ...func(*rekognition.Options)) (*rekognition.DetectLabelsOutput, error)
	DetectModerationLabels(ctx context.Context, params *rekognition.DetectModerationLabelsInput, optFns ...func(*rekognition.Options)) (*rekognition.DetectModerationLabelsOutput, error)
}

// RekognitionValidator validates vehicle photos using AWS Rekognition.
// It implements the port.ImageValidator interface.
type RekognitionValidator struct {
	client RekognitionAPI
}

// NewRekognitionValidator creates a new Rekognition-based image validator.
func NewRekognitionValidator(client RekognitionAPI) *RekognitionValidator {
	return &RekognitionValidator{client: client}
}

// ValidateVehiclePhoto checks an S3 image for inappropriate content and vehicle presence.
// It returns IsValid=false if moderation labels are found or no vehicle is detected.
func (v *RekognitionValidator) ValidateVehiclePhoto(ctx context.Context, s3Bucket, s3Key string) (*port.ImageValidationResult, error) {
	s3Image := &rektypes.S3Object{
		Bucket: &s3Bucket,
		Name:   &s3Key,
	}

	// Step 1: Check for inappropriate content.
	modOutput, err := v.client.DetectModerationLabels(ctx, &rekognition.DetectModerationLabelsInput{
		Image: &rektypes.Image{
			S3Object: s3Image,
		},
	})
	if err != nil {
		return nil, fmt.Errorf("detecting moderation labels for %s/%s: %w", s3Bucket, s3Key, err)
	}

	if len(modOutput.ModerationLabels) > 0 {
		labels := make([]string, 0, len(modOutput.ModerationLabels))
		for _, ml := range modOutput.ModerationLabels {
			if ml.Name != nil {
				labels = append(labels, *ml.Name)
			}
		}
		return &port.ImageValidationResult{
			IsValid: false,
			Labels:  labels,
			Reason:  "inappropriate content detected",
		}, nil
	}

	// Step 2: Check for vehicle labels.
	minConfidence := float32(80.0)
	labelOutput, err := v.client.DetectLabels(ctx, &rekognition.DetectLabelsInput{
		Image: &rektypes.Image{
			S3Object: s3Image,
		},
		MinConfidence: aws.Float32(minConfidence),
	})
	if err != nil {
		return nil, fmt.Errorf("detecting labels for %s/%s: %w", s3Bucket, s3Key, err)
	}

	detectedLabels := make([]string, 0, len(labelOutput.Labels))
	vehicleFound := false
	for _, lbl := range labelOutput.Labels {
		if lbl.Name != nil {
			detectedLabels = append(detectedLabels, *lbl.Name)
			if vehicleLabels[*lbl.Name] {
				vehicleFound = true
			}
		}
	}

	if !vehicleFound {
		return &port.ImageValidationResult{
			IsValid: false,
			Labels:  detectedLabels,
			Reason:  "no vehicle detected",
		}, nil
	}

	return &port.ImageValidationResult{
		IsValid: true,
		Labels:  detectedLabels,
		Reason:  "",
	}, nil
}
