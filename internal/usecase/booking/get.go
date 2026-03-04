package bookinguc

import (
	"context"

	"github.com/davidlramirez95/towcommand/internal/domain/booking"
	domainerrors "github.com/davidlramirez95/towcommand/internal/domain/errors"
	"github.com/davidlramirez95/towcommand/internal/domain/user"
)

// GetBookingInput holds the data needed to retrieve a booking.
type GetBookingInput struct {
	BookingID  string
	CallerID   string
	CallerType string
}

// GetBookingUseCase orchestrates retrieving a single booking with authorization checks.
type GetBookingUseCase struct {
	repo BookingFinder
}

// NewGetBookingUseCase constructs a GetBookingUseCase with its dependencies.
func NewGetBookingUseCase(repo BookingFinder) *GetBookingUseCase {
	return &GetBookingUseCase{repo: repo}
}

// Execute retrieves a booking by ID, enforcing authorization.
// Customers can view their own bookings; providers can view assigned bookings; admins see all.
func (uc *GetBookingUseCase) Execute(ctx context.Context, input GetBookingInput) (*booking.Booking, error) {
	b, err := uc.repo.FindByID(ctx, input.BookingID)
	if err != nil {
		return nil, err
	}
	if b == nil {
		return nil, domainerrors.NewNotFoundError("Booking", input.BookingID)
	}

	if !canAccessBooking(b, input.CallerID, input.CallerType) {
		return nil, domainerrors.NewForbiddenError("You do not have access to this booking")
	}

	return b, nil
}

// canAccessBooking checks whether a caller can view a booking.
func canAccessBooking(b *booking.Booking, callerID, callerType string) bool {
	if b.CustomerID == callerID {
		return true
	}
	if b.ProviderID == callerID {
		return true
	}
	return isAdmin(callerType)
}

// isAdmin returns true if the caller type is admin or ops_agent.
func isAdmin(callerType string) bool {
	return callerType == string(user.UserTypeAdmin) || callerType == string(user.UserTypeOpsAgent)
}
