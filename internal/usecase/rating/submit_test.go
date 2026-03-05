package ratinguc

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/davidlramirez95/towcommand/internal/domain/booking"
	domainerrors "github.com/davidlramirez95/towcommand/internal/domain/errors"
	"github.com/davidlramirez95/towcommand/internal/domain/provider"
	"github.com/davidlramirez95/towcommand/internal/domain/rating"
	"github.com/davidlramirez95/towcommand/internal/usecase/port"
)

// --- Mocks ---

type mockRatingSaver struct{ mock.Mock }

func (m *mockRatingSaver) Save(ctx context.Context, r *rating.Rating) error {
	args := m.Called(ctx, r)
	return args.Error(0)
}

type mockRatingByBookingFinder struct{ mock.Mock }

func (m *mockRatingByBookingFinder) FindByBooking(ctx context.Context, bookingID string) (*rating.Rating, error) {
	args := m.Called(ctx, bookingID)
	if v := args.Get(0); v != nil {
		return v.(*rating.Rating), args.Error(1)
	}
	return nil, args.Error(1)
}

type mockRatingByProviderLister struct{ mock.Mock }

func (m *mockRatingByProviderLister) FindByProvider(ctx context.Context, providerID string, limit int32) ([]rating.Rating, error) {
	args := m.Called(ctx, providerID, limit)
	if v := args.Get(0); v != nil {
		return v.([]rating.Rating), args.Error(1)
	}
	return nil, args.Error(1)
}

type mockBookingFinder struct{ mock.Mock }

func (m *mockBookingFinder) FindByID(ctx context.Context, bookingID string) (*booking.Booking, error) {
	args := m.Called(ctx, bookingID)
	if v := args.Get(0); v != nil {
		return v.(*booking.Booking), args.Error(1)
	}
	return nil, args.Error(1)
}

type mockProviderFinder struct{ mock.Mock }

func (m *mockProviderFinder) FindByID(ctx context.Context, providerID string) (*provider.Provider, error) {
	args := m.Called(ctx, providerID)
	if v := args.Get(0); v != nil {
		return v.(*provider.Provider), args.Error(1)
	}
	return nil, args.Error(1)
}

type mockProviderSaver struct{ mock.Mock }

func (m *mockProviderSaver) Save(ctx context.Context, p *provider.Provider) error {
	args := m.Called(ctx, p)
	return args.Error(0)
}

type mockEventPublisher struct{ mock.Mock }

func (m *mockEventPublisher) Publish(ctx context.Context, source, detailType string, detail any, actor *port.Actor) error {
	args := m.Called(ctx, source, detailType, detail, actor)
	return args.Error(0)
}

// --- Helpers ---

func completedBooking() *booking.Booking {
	return &booking.Booking{
		BookingID:  "booking-001",
		CustomerID: "customer-123",
		ProviderID: "provider-456",
		Status:     booking.BookingStatusCompleted,
	}
}

func testProvider() *provider.Provider {
	return &provider.Provider{
		ProviderID:         "provider-456",
		Rating:             4.0,
		TotalJobsCompleted: 10,
	}
}

func validInput() *SubmitRatingInput {
	return &SubmitRatingInput{
		BookingID:  "booking-001",
		CustomerID: "customer-123",
		Score:      5,
		Comment:    "Great service!",
		Tags:       []string{"fast", "professional"},
	}
}

// --- Tests ---

func TestSubmitRatingUseCase_Execute(t *testing.T) {
	fixedTime := time.Date(2026, 3, 5, 10, 0, 0, 0, time.UTC)

	tests := []struct {
		name       string
		input      *SubmitRatingInput
		setup      func(rs *mockRatingSaver, rf *mockRatingByBookingFinder, rl *mockRatingByProviderLister, bf *mockBookingFinder, pf *mockProviderFinder, ps *mockProviderSaver, ep *mockEventPublisher)
		wantRating *rating.Rating
		wantErr    func(t *testing.T, err error)
	}{
		{
			name:  "successful submission with provider average update",
			input: validInput(),
			setup: func(rs *mockRatingSaver, rf *mockRatingByBookingFinder, rl *mockRatingByProviderLister, bf *mockBookingFinder, pf *mockProviderFinder, ps *mockProviderSaver, ep *mockEventPublisher) {
				bf.On("FindByID", mock.Anything, "booking-001").Return(completedBooking(), nil)
				rf.On("FindByBooking", mock.Anything, "booking-001").Return(nil, nil)
				pf.On("FindByID", mock.Anything, "provider-456").Return(testProvider(), nil)
				rs.On("Save", mock.Anything, mock.AnythingOfType("*rating.Rating")).Return(nil)
				rl.On("FindByProvider", mock.Anything, "provider-456", int32(1000)).Return([]rating.Rating{
					{Score: 4},
					{Score: 3},
					{Score: 5}, // the newly saved rating
				}, nil)
				ps.On("Save", mock.Anything, mock.AnythingOfType("*provider.Provider")).Return(nil)
				ep.On("Publish", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)
			},
			wantRating: &rating.Rating{
				BookingID:  "booking-001",
				CustomerID: "customer-123",
				ProviderID: "provider-456",
				Score:      5,
				Comment:    "Great service!",
				Tags:       []string{"fast", "professional"},
				CreatedAt:  fixedTime,
			},
			wantErr: nil,
		},
		{
			name:  "booking not found",
			input: validInput(),
			setup: func(rs *mockRatingSaver, rf *mockRatingByBookingFinder, rl *mockRatingByProviderLister, bf *mockBookingFinder, pf *mockProviderFinder, ps *mockProviderSaver, ep *mockEventPublisher) {
				bf.On("FindByID", mock.Anything, "booking-001").Return(nil, nil)
			},
			wantErr: func(t *testing.T, err error) {
				t.Helper()
				var appErr *domainerrors.AppError
				require.True(t, errors.As(err, &appErr))
				assert.Equal(t, domainerrors.CodeNotFound, appErr.Code)
			},
		},
		{
			name:  "booking not completed",
			input: validInput(),
			setup: func(rs *mockRatingSaver, rf *mockRatingByBookingFinder, rl *mockRatingByProviderLister, bf *mockBookingFinder, pf *mockProviderFinder, ps *mockProviderSaver, ep *mockEventPublisher) {
				b := completedBooking()
				b.Status = booking.BookingStatusEnRoute
				bf.On("FindByID", mock.Anything, "booking-001").Return(b, nil)
			},
			wantErr: func(t *testing.T, err error) {
				t.Helper()
				var appErr *domainerrors.AppError
				require.True(t, errors.As(err, &appErr))
				assert.Equal(t, domainerrors.CodeConflict, appErr.Code)
				assert.Contains(t, appErr.Message, "not completed")
			},
		},
		{
			name: "not the customer of the booking",
			input: &SubmitRatingInput{
				BookingID:  "booking-001",
				CustomerID: "someone-else",
				Score:      5,
			},
			setup: func(rs *mockRatingSaver, rf *mockRatingByBookingFinder, rl *mockRatingByProviderLister, bf *mockBookingFinder, pf *mockProviderFinder, ps *mockProviderSaver, ep *mockEventPublisher) {
				bf.On("FindByID", mock.Anything, "booking-001").Return(completedBooking(), nil)
			},
			wantErr: func(t *testing.T, err error) {
				t.Helper()
				var appErr *domainerrors.AppError
				require.True(t, errors.As(err, &appErr))
				assert.Equal(t, domainerrors.CodeForbidden, appErr.Code)
				assert.Contains(t, appErr.Message, "only the booking customer")
			},
		},
		{
			name:  "duplicate rating",
			input: validInput(),
			setup: func(rs *mockRatingSaver, rf *mockRatingByBookingFinder, rl *mockRatingByProviderLister, bf *mockBookingFinder, pf *mockProviderFinder, ps *mockProviderSaver, ep *mockEventPublisher) {
				bf.On("FindByID", mock.Anything, "booking-001").Return(completedBooking(), nil)
				rf.On("FindByBooking", mock.Anything, "booking-001").Return(&rating.Rating{
					BookingID: "booking-001",
					Score:     4,
				}, nil)
			},
			wantErr: func(t *testing.T, err error) {
				t.Helper()
				var appErr *domainerrors.AppError
				require.True(t, errors.As(err, &appErr))
				assert.Equal(t, domainerrors.CodeConflict, appErr.Code)
				assert.Contains(t, appErr.Message, "already submitted")
			},
		},
		{
			name: "score too low (0)",
			input: &SubmitRatingInput{
				BookingID:  "booking-001",
				CustomerID: "customer-123",
				Score:      0,
			},
			setup: func(rs *mockRatingSaver, rf *mockRatingByBookingFinder, rl *mockRatingByProviderLister, bf *mockBookingFinder, pf *mockProviderFinder, ps *mockProviderSaver, ep *mockEventPublisher) {
				bf.On("FindByID", mock.Anything, "booking-001").Return(completedBooking(), nil)
				rf.On("FindByBooking", mock.Anything, "booking-001").Return(nil, nil)
			},
			wantErr: func(t *testing.T, err error) {
				t.Helper()
				var appErr *domainerrors.AppError
				require.True(t, errors.As(err, &appErr))
				assert.Equal(t, domainerrors.CodeValidationError, appErr.Code)
				assert.Contains(t, appErr.Message, "score must be between 1 and 5")
			},
		},
		{
			name: "score too high (6)",
			input: &SubmitRatingInput{
				BookingID:  "booking-001",
				CustomerID: "customer-123",
				Score:      6,
			},
			setup: func(rs *mockRatingSaver, rf *mockRatingByBookingFinder, rl *mockRatingByProviderLister, bf *mockBookingFinder, pf *mockProviderFinder, ps *mockProviderSaver, ep *mockEventPublisher) {
				bf.On("FindByID", mock.Anything, "booking-001").Return(completedBooking(), nil)
				rf.On("FindByBooking", mock.Anything, "booking-001").Return(nil, nil)
			},
			wantErr: func(t *testing.T, err error) {
				t.Helper()
				var appErr *domainerrors.AppError
				require.True(t, errors.As(err, &appErr))
				assert.Equal(t, domainerrors.CodeValidationError, appErr.Code)
				assert.Contains(t, appErr.Message, "score must be between 1 and 5")
			},
		},
		{
			name: "negative score (-1)",
			input: &SubmitRatingInput{
				BookingID:  "booking-001",
				CustomerID: "customer-123",
				Score:      -1,
			},
			setup: func(rs *mockRatingSaver, rf *mockRatingByBookingFinder, rl *mockRatingByProviderLister, bf *mockBookingFinder, pf *mockProviderFinder, ps *mockProviderSaver, ep *mockEventPublisher) {
				bf.On("FindByID", mock.Anything, "booking-001").Return(completedBooking(), nil)
				rf.On("FindByBooking", mock.Anything, "booking-001").Return(nil, nil)
			},
			wantErr: func(t *testing.T, err error) {
				t.Helper()
				var appErr *domainerrors.AppError
				require.True(t, errors.As(err, &appErr))
				assert.Equal(t, domainerrors.CodeValidationError, appErr.Code)
			},
		},
		{
			name:  "booking finder returns error",
			input: validInput(),
			setup: func(rs *mockRatingSaver, rf *mockRatingByBookingFinder, rl *mockRatingByProviderLister, bf *mockBookingFinder, pf *mockProviderFinder, ps *mockProviderSaver, ep *mockEventPublisher) {
				bf.On("FindByID", mock.Anything, "booking-001").Return(nil, domainerrors.NewInternalError("db error"))
			},
			wantErr: func(t *testing.T, err error) {
				t.Helper()
				assert.Error(t, err)
			},
		},
		{
			name:  "provider average recalculated correctly",
			input: validInput(),
			setup: func(rs *mockRatingSaver, rf *mockRatingByBookingFinder, rl *mockRatingByProviderLister, bf *mockBookingFinder, pf *mockProviderFinder, ps *mockProviderSaver, ep *mockEventPublisher) {
				bf.On("FindByID", mock.Anything, "booking-001").Return(completedBooking(), nil)
				rf.On("FindByBooking", mock.Anything, "booking-001").Return(nil, nil)
				pf.On("FindByID", mock.Anything, "provider-456").Return(testProvider(), nil)
				rs.On("Save", mock.Anything, mock.AnythingOfType("*rating.Rating")).Return(nil)
				// Existing ratings: 4, 4, 4, 4 and the new one is 5
				// Average = (4+4+4+4+5) / 5 = 4.2
				rl.On("FindByProvider", mock.Anything, "provider-456", int32(1000)).Return([]rating.Rating{
					{Score: 4},
					{Score: 4},
					{Score: 4},
					{Score: 4},
					{Score: 5},
				}, nil)
				ps.On("Save", mock.Anything, mock.MatchedBy(func(p *provider.Provider) bool {
					return p.Rating == 4.2
				})).Return(nil)
				ep.On("Publish", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)
			},
			wantRating: &rating.Rating{
				BookingID:  "booking-001",
				CustomerID: "customer-123",
				ProviderID: "provider-456",
				Score:      5,
				Comment:    "Great service!",
				Tags:       []string{"fast", "professional"},
				CreatedAt:  fixedTime,
			},
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rs := new(mockRatingSaver)
			rf := new(mockRatingByBookingFinder)
			rl := new(mockRatingByProviderLister)
			bf := new(mockBookingFinder)
			pf := new(mockProviderFinder)
			ps := new(mockProviderSaver)
			ep := new(mockEventPublisher)

			tt.setup(rs, rf, rl, bf, pf, ps, ep)

			uc := NewSubmitRatingUseCase(rs, rf, rl, bf, pf, ps, ep)
			uc.now = func() time.Time { return fixedTime }

			result, err := uc.Execute(context.Background(), tt.input)

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

			rs.AssertExpectations(t)
			rf.AssertExpectations(t)
			bf.AssertExpectations(t)
			ep.AssertExpectations(t)
		})
	}
}
