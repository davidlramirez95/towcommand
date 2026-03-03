package provider

import (
	"time"

	"github.com/davidlramirez95/towcommand/internal/domain/user"
)

// ProviderStatus represents the verification/account status of a provider.
type ProviderStatus string

const (
	ProviderStatusPendingVerification ProviderStatus = "pending_verification"
	ProviderStatusActive              ProviderStatus = "active"
	ProviderStatusSuspended           ProviderStatus = "suspended"
	ProviderStatusDeactivated         ProviderStatus = "deactivated"
)

// TruckType represents the type of tow truck a provider operates.
type TruckType string

const (
	TruckTypeFlatbed           TruckType = "flatbed"
	TruckTypeWheelLift         TruckType = "wheel_lift"
	TruckTypeBoom              TruckType = "boom"
	TruckTypeMotorcycleCarrier TruckType = "motorcycle_carrier"
)

// ClearanceStatus represents the approval state of a KYC document.
type ClearanceStatus string

const (
	ClearanceStatusPending  ClearanceStatus = "pending"
	ClearanceStatusApproved ClearanceStatus = "approved"
	ClearanceStatusExpired  ClearanceStatus = "expired"
)

// DocType represents the type of provider document.
type DocType string

const (
	DocTypeNBIClearance      DocType = "nbi_clearance"
	DocTypeLTORegistration   DocType = "lto_registration"
	DocTypeDrugTest          DocType = "drug_test"
	DocTypeVehicleInspection DocType = "vehicle_inspection"
	DocTypeInsurance         DocType = "insurance"
)

// DocStatus represents the review status of a provider document.
type DocStatus string

const (
	DocStatusPending  DocStatus = "pending"
	DocStatusApproved DocStatus = "approved"
	DocStatusRejected DocStatus = "rejected"
)

// Provider represents a tow truck operator/service provider.
type Provider struct {
	ProviderID          string          `json:"providerId" validate:"required"`
	CognitoSub          string          `json:"cognitoSub" validate:"required"`
	Name                string          `json:"name" validate:"required"`
	Phone               string          `json:"phone" validate:"required"`
	Email               string          `json:"email" validate:"required,email"`
	Status              ProviderStatus  `json:"status" validate:"required"`
	TrustTier           user.TrustTier  `json:"trustTier" validate:"required"`
	TruckType           TruckType       `json:"truckType" validate:"required"`
	MaxWeightCapacityKg int             `json:"maxWeightCapacityKg" validate:"required,gt=0"`
	PlateNumber         string          `json:"plateNumber" validate:"required"`
	LTORegistration     string          `json:"ltoRegistration" validate:"required"`
	NBIClearanceStatus  ClearanceStatus `json:"nbiClearanceStatus" validate:"required"`
	DrugTestStatus      ClearanceStatus `json:"drugTestStatus" validate:"required"`
	MMADAccredited      bool            `json:"mmadAccredited"`
	Rating              float64         `json:"rating" validate:"min=0,max=5"`
	TotalJobsCompleted  int             `json:"totalJobsCompleted" validate:"min=0"`
	AcceptanceRate      float64         `json:"acceptanceRate" validate:"min=0,max=1"`
	IsOnline            bool            `json:"isOnline"`
	CurrentLat          *float64        `json:"currentLat,omitempty"`
	CurrentLng          *float64        `json:"currentLng,omitempty"`
	LastLocationUpdate  *time.Time      `json:"lastLocationUpdate,omitempty"`
	ServiceAreas        []string        `json:"serviceAreas"`
	CreatedAt           time.Time       `json:"createdAt"`
	UpdatedAt           time.Time       `json:"updatedAt"`
}

// ProviderDoc represents a KYC/compliance document uploaded by a provider.
type ProviderDoc struct {
	ProviderID string     `json:"providerId" validate:"required"`
	DocType    DocType    `json:"docType" validate:"required"`
	S3Key      string     `json:"s3Key" validate:"required"`
	Status     DocStatus  `json:"status" validate:"required"`
	ReviewedBy string     `json:"reviewedBy,omitempty"`
	ReviewedAt *time.Time `json:"reviewedAt,omitempty"`
	ExpiresAt  *time.Time `json:"expiresAt,omitempty"`
	UploadedAt time.Time  `json:"uploadedAt"`
}

// ProviderLocation represents a real-time GPS ping from a provider.
type ProviderLocation struct {
	ProviderID string    `json:"providerId" validate:"required"`
	Lat        float64   `json:"lat" validate:"required,latitude"`
	Lng        float64   `json:"lng" validate:"required,longitude"`
	Heading    float64   `json:"heading" validate:"min=0,max=360"`
	Speed      float64   `json:"speed" validate:"min=0"`
	Timestamp  time.Time `json:"timestamp"`
}
