package ratinguc

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	domainerrors "github.com/davidlramirez95/towcommand/internal/domain/errors"
	"github.com/davidlramirez95/towcommand/internal/domain/rating"
)

// mockRatingFinder is a test double for the RatingByBookingFinder interface.
type mockRatingFinder struct{ mock.Mock }

func (m *mockRatingFinder) FindByBooking(ctx context.Context, bookingID string) (*rating.Rating, error) {
	args := m.Called(ctx, bookingID)
	if v := args.Get(0); v != nil {
		return v.(*rating.Rating), args.Error(1)
	}
	return nil, args.Error(1)
}

func TestGetBookingRatingUseCase_Execute(t *testing.T) {
	fixedTime := time.Date(2026, 3, 5, 10, 0, 0, 0, time.UTC)

	tests := []struct {
		name       string
		bookingID  string
		setup      func(m *mockRatingFinder)
		wantRating *rating.Rating
		wantErr    func(t *testing.T, err error)
	}{
		{
			name:      "rating found",
			bookingID: "booking-001",
			setup: func(m *mockRatingFinder) {
				m.On("FindByBooking", mock.Anything, "booking-001").Return(&rating.Rating{
					BookingID:  "booking-001",
					CustomerID: "customer-123",
					ProviderID: "provider-456",
					Score:      5,
					Comment:    "Excellent!",
					Tags:       []string{"fast"},
					CreatedAt:  fixedTime,
				}, nil)
			},
			wantRating: &rating.Rating{
				BookingID:  "booking-001",
				CustomerID: "customer-123",
				ProviderID: "provider-456",
				Score:      5,
				Comment:    "Excellent!",
				Tags:       []string{"fast"},
				CreatedAt:  fixedTime,
			},
			wantErr: nil,
		},
		{
			name:      "rating not found",
			bookingID: "booking-999",
			setup: func(m *mockRatingFinder) {
				m.On("FindByBooking", mock.Anything, "booking-999").Return(nil, nil)
			},
			wantErr: func(t *testing.T, err error) {
				t.Helper()
				var appErr *domainerrors.AppError
				require.True(t, errors.As(err, &appErr))
				assert.Equal(t, domainerrors.CodeNotFound, appErr.Code)
			},
		},
		{
			name:      "repository returns error",
			bookingID: "booking-err",
			setup: func(m *mockRatingFinder) {
				m.On("FindByBooking", mock.Anything, "booking-err").Return(nil, domainerrors.NewInternalError("db failure"))
			},
			wantErr: func(t *testing.T, err error) {
				t.Helper()
				assert.Error(t, err)
				var appErr *domainerrors.AppError
				require.True(t, errors.As(err, &appErr))
				assert.Equal(t, domainerrors.CodeInternalError, appErr.Code)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			finder := new(mockRatingFinder)
			tt.setup(finder)

			uc := NewGetBookingRatingUseCase(finder)
			result, err := uc.Execute(context.Background(), tt.bookingID)

			if tt.wantErr != nil {
				tt.wantErr(t, err)
				assert.Nil(t, result)
				return
			}

			require.NoError(t, err)
			require.NotNil(t, result)
			assert.Equal(t, tt.wantRating.BookingID, result.BookingID)
			assert.Equal(t, tt.wantRating.CustomerID, result.CustomerID)
			assert.Equal(t, tt.wantRating.ProviderID, result.ProviderID)
			assert.Equal(t, tt.wantRating.Score, result.Score)
			assert.Equal(t, tt.wantRating.Comment, result.Comment)
			assert.Equal(t, tt.wantRating.Tags, result.Tags)
			assert.Equal(t, tt.wantRating.CreatedAt, result.CreatedAt)

			finder.AssertExpectations(t)
		})
	}
}
