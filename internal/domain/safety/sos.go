package safety

import "time"

// TriggerType represents how an SOS alert was activated.
type TriggerType string

const (
	TriggerTypeTripleTap TriggerType = "TRIPLE_TAP"
	TriggerTypeShake     TriggerType = "SHAKE"
	TriggerTypeCodeWord  TriggerType = "CODE_WORD"
	TriggerTypeButton    TriggerType = "BUTTON"
)

// RiskScore quantifies the threat level of an SOS situation.
type RiskScore struct {
	Score      int       `json:"score" validate:"min=0,max=100"`
	Level      string    `json:"level" validate:"required,oneof=low medium high critical"`
	Factors    []string  `json:"factors,omitempty"`
	AssessedAt time.Time `json:"assessedAt"`
}

// ComputeLevel derives the risk level from the numeric score.
func (r *RiskScore) ComputeLevel() string {
	switch {
	case r.Score >= 80:
		return "critical"
	case r.Score >= 60:
		return "high"
	case r.Score >= 30:
		return "medium"
	default:
		return "low"
	}
}

// SOSAlert represents an emergency distress signal from a user.
type SOSAlert struct {
	AlertID     string      `json:"alertId" validate:"required"`
	BookingID   string      `json:"bookingId,omitempty"`
	TriggeredBy string      `json:"triggeredBy" validate:"required"`
	TriggerType TriggerType `json:"triggerType" validate:"required"`
	Lat         float64     `json:"lat" validate:"required,latitude"`
	Lng         float64     `json:"lng" validate:"required,longitude"`
	Risk        RiskScore   `json:"risk"`
	Resolved    bool        `json:"resolved"`
	ResolvedAt  *time.Time  `json:"resolvedAt,omitempty"`
	ResolvedBy  string      `json:"resolvedBy,omitempty"`
	Timestamp   time.Time   `json:"timestamp"`
}
