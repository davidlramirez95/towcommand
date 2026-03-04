package port

import (
	"context"

	"github.com/davidlramirez95/towcommand/internal/domain/user"
)

// UserSaver persists a new user.
type UserSaver interface {
	Save(ctx context.Context, u *user.User) error
}

// UserFinder retrieves a user by their ID.
type UserFinder interface {
	FindByID(ctx context.Context, userID string) (*user.User, error)
}

// UserByEmailFinder retrieves a user by their email address via GSI1.
type UserByEmailFinder interface {
	FindByEmail(ctx context.Context, email string) (*user.User, error)
}

// UserByPhoneFinder retrieves a user by their phone number via GSI5.
type UserByPhoneFinder interface {
	FindByPhone(ctx context.Context, phone string) (*user.User, error)
}

// VehicleAdder adds a vehicle to a user's profile.
type VehicleAdder interface {
	AddVehicle(ctx context.Context, v *user.UserVehicle) error
}

// VehicleLister lists all vehicles registered under a user.
type VehicleLister interface {
	GetVehicles(ctx context.Context, userID string) ([]user.UserVehicle, error)
}
