package bookinguc

import (
	"context"

	"github.com/davidlramirez95/towcommand/internal/domain/booking"
)

// ListBookingsInput holds the data needed to list bookings.
type ListBookingsInput struct {
	CallerID     string
	CallerType   string
	Limit        int32
	StatusFilter string
}

// ListBookingsOutput is the response for a list operation.
type ListBookingsOutput struct {
	Items []booking.Booking `json:"items"`
	Count int               `json:"count"`
}

// BookingLister combines the read interfaces needed by the list use case.
type BookingLister interface {
	BookingByUserLister
	BookingByStatusLister
}

// ListBookingsUseCase orchestrates listing bookings with role-based filtering.
type ListBookingsUseCase struct {
	repo BookingLister
}

// NewListBookingsUseCase constructs a ListBookingsUseCase with its dependencies.
func NewListBookingsUseCase(repo BookingLister) *ListBookingsUseCase {
	return &ListBookingsUseCase{repo: repo}
}

// Execute lists bookings based on the caller's role and optional filters.
// Admins with a status filter query across all bookings; others see only their own.
func (uc *ListBookingsUseCase) Execute(ctx context.Context, input ListBookingsInput) (*ListBookingsOutput, error) {
	limit := input.Limit
	if limit <= 0 {
		limit = 25
	}
	if limit > 100 {
		limit = 100
	}

	var (
		bookings []booking.Booking
		err      error
	)

	if isAdmin(input.CallerType) && input.StatusFilter != "" {
		bookings, err = uc.repo.FindByStatus(ctx, booking.BookingStatus(input.StatusFilter), limit)
	} else {
		bookings, err = uc.repo.FindByUser(ctx, input.CallerID, limit)
		// Client-side status filter for non-admin users.
		if err == nil && input.StatusFilter != "" && !isAdmin(input.CallerType) {
			bookings = filterByStatus(bookings, booking.BookingStatus(input.StatusFilter))
		}
	}
	if err != nil {
		return nil, err
	}

	return &ListBookingsOutput{
		Items: bookings,
		Count: len(bookings),
	}, nil
}

func filterByStatus(bookings []booking.Booking, status booking.BookingStatus) []booking.Booking {
	filtered := make([]booking.Booking, 0, len(bookings))
	for i := range bookings {
		if bookings[i].Status == status {
			filtered = append(filtered, bookings[i])
		}
	}
	return filtered
}
