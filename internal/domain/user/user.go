package user

import "time"

// UserType represents the role of a user in the system.
type UserType string

const (
	UserTypeCustomer     UserType = "customer"
	UserTypeProvider     UserType = "provider"
	UserTypeFleetManager UserType = "fleet_manager"
	UserTypeOpsAgent     UserType = "ops_agent"
	UserTypeAdmin        UserType = "admin"
)

// TrustTier represents the loyalty/trust level of a user.
type TrustTier string

const (
	TrustTierBasic      TrustTier = "basic"
	TrustTierVerified   TrustTier = "verified"
	TrustTierSukiSilver TrustTier = "suki_silver"
	TrustTierSukiGold   TrustTier = "suki_gold"
	TrustTierSukiElite  TrustTier = "suki_elite"
)

// WeightClass represents the weight category of a vehicle.
type WeightClass string

const (
	WeightClassMotorcycle WeightClass = "motorcycle"
	WeightClassLight      WeightClass = "light"
	WeightClassMedium     WeightClass = "medium"
	WeightClassHeavy      WeightClass = "heavy"
	WeightClassSuperHeavy WeightClass = "super_heavy"
)

// UserStatus represents the account status of a user.
type UserStatus string

const (
	UserStatusActive    UserStatus = "active"
	UserStatusSuspended UserStatus = "suspended"
	UserStatusBanned    UserStatus = "banned"
)

// Language represents a supported UI language.
type Language string

const (
	LanguageEnglish  Language = "en"
	LanguageFilipino Language = "fil"
)

// User represents a platform user (customer, provider, admin, etc.).
type User struct {
	UserID     string     `json:"userId" validate:"required"`
	CognitoSub string     `json:"cognitoSub" validate:"required"`
	Email      string     `json:"email" validate:"required,email"`
	Phone      string     `json:"phone" validate:"required"`
	Name       string     `json:"name" validate:"required"`
	UserType   UserType   `json:"userType" validate:"required"`
	TrustTier  TrustTier  `json:"trustTier" validate:"required"`
	Language   Language   `json:"language" validate:"required,oneof=en fil"`
	Status     UserStatus `json:"status" validate:"required,oneof=active suspended banned"`
	CreatedAt  time.Time  `json:"createdAt"`
	UpdatedAt  time.Time  `json:"updatedAt"`
}

// UserVehicle represents a vehicle registered by a user.
type UserVehicle struct {
	VehicleID   string      `json:"vehicleId" validate:"required"`
	UserID      string      `json:"userId" validate:"required"`
	Make        string      `json:"make" validate:"required"`
	Model       string      `json:"model" validate:"required"`
	Year        int         `json:"year" validate:"required,min=1900"`
	PlateNumber string      `json:"plateNumber" validate:"required"`
	WeightClass WeightClass `json:"weightClass" validate:"required"`
	Color       string      `json:"color" validate:"required"`
	PhotoURL    string      `json:"photoUrl,omitempty"`
	IsDefault   bool        `json:"isDefault"`
	CreatedAt   time.Time   `json:"createdAt"`
}
