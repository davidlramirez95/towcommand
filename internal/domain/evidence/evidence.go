package evidence

import "time"

// PhotoPosition represents the required camera angle for condition documentation.
type PhotoPosition string

const (
	PhotoPositionFront      PhotoPosition = "FRONT"
	PhotoPositionRear       PhotoPosition = "REAR"
	PhotoPositionLeft       PhotoPosition = "LEFT"
	PhotoPositionRight      PhotoPosition = "RIGHT"
	PhotoPositionFrontLeft  PhotoPosition = "FRONT_LEFT"
	PhotoPositionFrontRight PhotoPosition = "FRONT_RIGHT"
	PhotoPositionRearLeft   PhotoPosition = "REAR_LEFT"
	PhotoPositionRearRight  PhotoPosition = "REAR_RIGHT"
)

// AllPhotoPositions is the complete set of required photo angles.
var AllPhotoPositions = []PhotoPosition{
	PhotoPositionFront,
	PhotoPositionRear,
	PhotoPositionLeft,
	PhotoPositionRight,
	PhotoPositionFrontLeft,
	PhotoPositionFrontRight,
	PhotoPositionRearLeft,
	PhotoPositionRearRight,
}

// HashIntegrity holds the cryptographic hash for tamper detection of media.
type HashIntegrity struct {
	Algorithm string `json:"algorithm" validate:"required"`
	Hash      string `json:"hash" validate:"required"`
}

// MediaItem represents a single photo or video captured during a condition report.
type MediaItem struct {
	MediaID    string        `json:"mediaId" validate:"required"`
	S3Key      string        `json:"s3Key" validate:"required"`
	Position   PhotoPosition `json:"position" validate:"required"`
	MimeType   string        `json:"mimeType" validate:"required"`
	Integrity  HashIntegrity `json:"integrity"`
	CapturedAt time.Time     `json:"capturedAt"`
}

// ConditionReport documents the vehicle state at pickup and/or dropoff.
type ConditionReport struct {
	ReportID   string      `json:"reportId" validate:"required"`
	BookingID  string      `json:"bookingId" validate:"required"`
	ProviderID string      `json:"providerId" validate:"required"`
	Phase      string      `json:"phase" validate:"required,oneof=pickup dropoff"`
	Media      []MediaItem `json:"media"`
	Notes      string      `json:"notes,omitempty"`
	CreatedAt  time.Time   `json:"createdAt"`
}

// IsComplete reports whether the condition report has all 8 required photo positions.
func (r *ConditionReport) IsComplete() bool {
	if len(r.Media) < len(AllPhotoPositions) {
		return false
	}
	covered := make(map[PhotoPosition]bool, len(AllPhotoPositions))
	for _, m := range r.Media {
		covered[m.Position] = true
	}
	for _, pos := range AllPhotoPositions {
		if !covered[pos] {
			return false
		}
	}
	return true
}
