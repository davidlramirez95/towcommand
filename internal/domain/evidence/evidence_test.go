package evidence

import (
	"testing"
	"time"
)

func TestIsComplete_AllPositions(t *testing.T) {
	now := time.Now()
	report := &ConditionReport{
		ReportID:   "report-1",
		BookingID:  "booking-1",
		ProviderID: "provider-1",
		Phase:      "pickup",
		CreatedAt:  now,
	}

	for _, pos := range AllPhotoPositions {
		report.Media = append(report.Media, MediaItem{
			MediaID:    "media-" + string(pos),
			S3Key:      "photos/" + string(pos) + ".jpg",
			Position:   pos,
			MimeType:   "image/jpeg",
			Integrity:  HashIntegrity{Algorithm: "sha256", Hash: "abc123"},
			CapturedAt: now,
		})
	}

	if !report.IsComplete() {
		t.Error("report with all 8 positions should be complete")
	}
}

func TestIsComplete_MissingPositions(t *testing.T) {
	now := time.Now()
	report := &ConditionReport{
		ReportID:   "report-1",
		BookingID:  "booking-1",
		ProviderID: "provider-1",
		Phase:      "pickup",
		CreatedAt:  now,
	}

	// Add only 4 of 8 positions
	partial := []PhotoPosition{PhotoPositionFront, PhotoPositionRear, PhotoPositionLeft, PhotoPositionRight}
	for _, pos := range partial {
		report.Media = append(report.Media, MediaItem{
			MediaID:    "media-" + string(pos),
			S3Key:      "photos/" + string(pos) + ".jpg",
			Position:   pos,
			MimeType:   "image/jpeg",
			Integrity:  HashIntegrity{Algorithm: "sha256", Hash: "abc123"},
			CapturedAt: now,
		})
	}

	if report.IsComplete() {
		t.Error("report missing 4 positions should not be complete")
	}
}

func TestIsComplete_EmptyMedia(t *testing.T) {
	report := &ConditionReport{
		ReportID:   "report-1",
		BookingID:  "booking-1",
		ProviderID: "provider-1",
		Phase:      "pickup",
		CreatedAt:  time.Now(),
	}

	if report.IsComplete() {
		t.Error("report with no media should not be complete")
	}
}

func TestIsComplete_DuplicatePositions(t *testing.T) {
	now := time.Now()
	report := &ConditionReport{
		ReportID:   "report-1",
		BookingID:  "booking-1",
		ProviderID: "provider-1",
		Phase:      "pickup",
		CreatedAt:  now,
	}

	// Add 8 items but all same position
	for i := 0; i < 8; i++ {
		report.Media = append(report.Media, MediaItem{
			MediaID:    "media-dup",
			S3Key:      "photos/front.jpg",
			Position:   PhotoPositionFront,
			MimeType:   "image/jpeg",
			Integrity:  HashIntegrity{Algorithm: "sha256", Hash: "abc123"},
			CapturedAt: now,
		})
	}

	if report.IsComplete() {
		t.Error("report with duplicate positions should not be complete")
	}
}

func TestAllPhotoPositions_Count(t *testing.T) {
	if len(AllPhotoPositions) != 8 {
		t.Errorf("AllPhotoPositions has %d entries, want 8", len(AllPhotoPositions))
	}
}

func TestAllPhotoPositions_Unique(t *testing.T) {
	seen := make(map[PhotoPosition]bool)
	for _, pos := range AllPhotoPositions {
		if seen[pos] {
			t.Errorf("duplicate position: %q", pos)
		}
		seen[pos] = true
	}
}
