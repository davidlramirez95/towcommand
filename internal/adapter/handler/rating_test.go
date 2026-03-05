package handler_test

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"testing"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/davidlramirez95/towcommand/internal/adapter/handler"
	"github.com/davidlramirez95/towcommand/internal/domain/booking"
	"github.com/davidlramirez95/towcommand/internal/domain/provider"
	"github.com/davidlramirez95/towcommand/internal/domain/rating"
	"github.com/davidlramirez95/towcommand/internal/usecase/port"
	ratinguc "github.com/davidlramirez95/towcommand/internal/usecase/rating"
)

// ---------------------------------------------------------------------------
// Mock implementations for rating use case dependencies
// ---------------------------------------------------------------------------

type mockRatingSaver struct {
	SaveFunc func(ctx context.Context, r *rating.Rating) error
}

func (m *mockRatingSaver) Save(ctx context.Context, r *rating.Rating) error {
	if m.SaveFunc != nil {
		return m.SaveFunc(ctx, r)
	}
	return nil
}

type mockRatingByBookingFinder struct {
	FindByBookingFunc func(ctx context.Context, bookingID string) (*rating.Rating, error)
}

func (m *mockRatingByBookingFinder) FindByBooking(ctx context.Context, bookingID string) (*rating.Rating, error) {
	if m.FindByBookingFunc != nil {
		return m.FindByBookingFunc(ctx, bookingID)
	}
	return nil, nil
}

type mockRatingByProviderLister struct {
	FindByProviderFunc func(ctx context.Context, providerID string, limit int32) ([]rating.Rating, error)
}

func (m *mockRatingByProviderLister) FindByProvider(ctx context.Context, providerID string, limit int32) ([]rating.Rating, error) {
	if m.FindByProviderFunc != nil {
		return m.FindByProviderFunc(ctx, providerID, limit)
	}
	return nil, nil
}

type mockBookingFinder struct {
	FindByIDFunc func(ctx context.Context, bookingID string) (*booking.Booking, error)
}

func (m *mockBookingFinder) FindByID(ctx context.Context, bookingID string) (*booking.Booking, error) {
	if m.FindByIDFunc != nil {
		return m.FindByIDFunc(ctx, bookingID)
	}
	return nil, nil
}

type mockProviderFinder struct {
	FindByIDFunc func(ctx context.Context, providerID string) (*provider.Provider, error)
}

func (m *mockProviderFinder) FindByID(ctx context.Context, providerID string) (*provider.Provider, error) {
	if m.FindByIDFunc != nil {
		return m.FindByIDFunc(ctx, providerID)
	}
	return nil, nil
}

type mockProviderSaver struct {
	SaveFunc func(ctx context.Context, p *provider.Provider) error
}

func (m *mockProviderSaver) Save(ctx context.Context, p *provider.Provider) error {
	if m.SaveFunc != nil {
		return m.SaveFunc(ctx, p)
	}
	return nil
}

type mockEventPublisher struct {
	PublishFunc func(ctx context.Context, source, detailType string, detail any, actor *port.Actor) error
}

func (m *mockEventPublisher) Publish(ctx context.Context, source, detailType string, detail any, actor *port.Actor) error {
	if m.PublishFunc != nil {
		return m.PublishFunc(ctx, source, detailType, detail, actor)
	}
	return nil
}

// ---------------------------------------------------------------------------
// Helper: build API Gateway event with Cognito auth
// ---------------------------------------------------------------------------

func apiEventWithAuth(userID string) *events.APIGatewayProxyRequest {
	return &events.APIGatewayProxyRequest{
		RequestContext: events.APIGatewayProxyRequestContext{
			Authorizer: map[string]interface{}{
				"claims": map[string]interface{}{
					"sub": userID,
				},
			},
		},
	}
}

// ---------------------------------------------------------------------------
// Helper: parse error response body
// ---------------------------------------------------------------------------

type errorBody struct {
	Error struct {
		Code    string `json:"code"`
		Message string `json:"message"`
	} `json:"error"`
}

func parseErrorBody(t *testing.T, body string) errorBody {
	t.Helper()
	var eb errorBody
	require.NoError(t, json.Unmarshal([]byte(body), &eb))
	return eb
}

// ---------------------------------------------------------------------------
// SubmitRatingHandler tests
// ---------------------------------------------------------------------------

func TestSubmitRatingHandler(t *testing.T) {
	now := time.Now().UTC()

	completedBooking := &booking.Booking{
		BookingID:  "bk-1",
		CustomerID: "user-1",
		ProviderID: "prov-1",
		Status:     booking.BookingStatusCompleted,
	}

	testProvider := &provider.Provider{
		ProviderID: "prov-1",
		Rating:     4.0,
	}

	tests := []struct {
		name        string
		event       *events.APIGatewayProxyRequest
		setupMocks  func(rs *mockRatingSaver, rbf *mockRatingByBookingFinder, rpl *mockRatingByProviderLister, bf *mockBookingFinder, pf *mockProviderFinder, ps *mockProviderSaver, ep *mockEventPublisher)
		wantStatus  int
		wantErrCode string
		checkBody   func(t *testing.T, body string)
	}{
		{
			name: "success",
			event: func() *events.APIGatewayProxyRequest {
				e := apiEventWithAuth("user-1")
				e.PathParameters = map[string]string{"id": "bk-1"}
				e.Body = `{"score":5,"comment":"Great service!","tags":["fast","friendly"]}`
				return e
			}(),
			setupMocks: func(rs *mockRatingSaver, rbf *mockRatingByBookingFinder, rpl *mockRatingByProviderLister, bf *mockBookingFinder, pf *mockProviderFinder, ps *mockProviderSaver, ep *mockEventPublisher) {
				bf.FindByIDFunc = func(_ context.Context, _ string) (*booking.Booking, error) {
					return completedBooking, nil
				}
				rbf.FindByBookingFunc = func(_ context.Context, _ string) (*rating.Rating, error) {
					return nil, nil // no existing rating
				}
				pf.FindByIDFunc = func(_ context.Context, _ string) (*provider.Provider, error) {
					return testProvider, nil
				}
				rpl.FindByProviderFunc = func(_ context.Context, _ string, _ int32) ([]rating.Rating, error) {
					return []rating.Rating{
						{Score: 5, CreatedAt: now},
					}, nil
				}
			},
			wantStatus: http.StatusCreated,
			checkBody: func(t *testing.T, body string) {
				t.Helper()
				var r rating.Rating
				require.NoError(t, json.Unmarshal([]byte(body), &r))
				assert.Equal(t, "bk-1", r.BookingID)
				assert.Equal(t, "user-1", r.CustomerID)
				assert.Equal(t, "prov-1", r.ProviderID)
				assert.Equal(t, 5, r.Score)
				assert.Equal(t, "Great service!", r.Comment)
				assert.Equal(t, []string{"fast", "friendly"}, r.Tags)
			},
		},
		{
			name: "unauthorized - no user ID",
			event: &events.APIGatewayProxyRequest{
				PathParameters: map[string]string{"id": "bk-1"},
				Body:           `{"score":5}`,
			},
			wantStatus:  http.StatusUnauthorized,
			wantErrCode: "UNAUTHORIZED",
		},
		{
			name: "missing booking ID",
			event: func() *events.APIGatewayProxyRequest {
				e := apiEventWithAuth("user-1")
				e.Body = `{"score":5}`
				return e
			}(),
			wantStatus:  http.StatusBadRequest,
			wantErrCode: "VALIDATION_ERROR",
		},
		{
			name: "invalid body - bad JSON",
			event: func() *events.APIGatewayProxyRequest {
				e := apiEventWithAuth("user-1")
				e.PathParameters = map[string]string{"id": "bk-1"}
				e.Body = `{not json}`
				return e
			}(),
			wantStatus:  http.StatusBadRequest,
			wantErrCode: "VALIDATION_ERROR",
		},
		{
			name: "invalid body - score out of range",
			event: func() *events.APIGatewayProxyRequest {
				e := apiEventWithAuth("user-1")
				e.PathParameters = map[string]string{"id": "bk-1"}
				e.Body = `{"score":0}`
				return e
			}(),
			wantStatus:  http.StatusBadRequest,
			wantErrCode: "VALIDATION_ERROR",
		},
		{
			name: "invalid body - missing required score",
			event: func() *events.APIGatewayProxyRequest {
				e := apiEventWithAuth("user-1")
				e.PathParameters = map[string]string{"id": "bk-1"}
				e.Body = `{"comment":"nice"}`
				return e
			}(),
			wantStatus:  http.StatusBadRequest,
			wantErrCode: "VALIDATION_ERROR",
		},
		{
			name: "use case error - booking not found",
			event: func() *events.APIGatewayProxyRequest {
				e := apiEventWithAuth("user-1")
				e.PathParameters = map[string]string{"id": "bk-missing"}
				e.Body = `{"score":4}`
				return e
			}(),
			setupMocks: func(rs *mockRatingSaver, rbf *mockRatingByBookingFinder, rpl *mockRatingByProviderLister, bf *mockBookingFinder, pf *mockProviderFinder, ps *mockProviderSaver, ep *mockEventPublisher) {
				bf.FindByIDFunc = func(_ context.Context, _ string) (*booking.Booking, error) {
					return nil, nil // not found
				}
			},
			wantStatus:  http.StatusNotFound,
			wantErrCode: "NOT_FOUND",
		},
		{
			name: "use case error - booking not completed",
			event: func() *events.APIGatewayProxyRequest {
				e := apiEventWithAuth("user-1")
				e.PathParameters = map[string]string{"id": "bk-1"}
				e.Body = `{"score":4}`
				return e
			}(),
			setupMocks: func(rs *mockRatingSaver, rbf *mockRatingByBookingFinder, rpl *mockRatingByProviderLister, bf *mockBookingFinder, pf *mockProviderFinder, ps *mockProviderSaver, ep *mockEventPublisher) {
				bf.FindByIDFunc = func(_ context.Context, _ string) (*booking.Booking, error) {
					return &booking.Booking{
						BookingID:  "bk-1",
						CustomerID: "user-1",
						Status:     booking.BookingStatusPending,
					}, nil
				}
			},
			wantStatus:  http.StatusConflict,
			wantErrCode: "CONFLICT",
		},
		{
			name: "use case error - duplicate rating",
			event: func() *events.APIGatewayProxyRequest {
				e := apiEventWithAuth("user-1")
				e.PathParameters = map[string]string{"id": "bk-1"}
				e.Body = `{"score":4}`
				return e
			}(),
			setupMocks: func(rs *mockRatingSaver, rbf *mockRatingByBookingFinder, rpl *mockRatingByProviderLister, bf *mockBookingFinder, pf *mockProviderFinder, ps *mockProviderSaver, ep *mockEventPublisher) {
				bf.FindByIDFunc = func(_ context.Context, _ string) (*booking.Booking, error) {
					return completedBooking, nil
				}
				rbf.FindByBookingFunc = func(_ context.Context, _ string) (*rating.Rating, error) {
					return &rating.Rating{BookingID: "bk-1", Score: 3}, nil // already exists
				}
			},
			wantStatus:  http.StatusConflict,
			wantErrCode: "CONFLICT",
		},
		{
			name: "use case error - internal error from repo",
			event: func() *events.APIGatewayProxyRequest {
				e := apiEventWithAuth("user-1")
				e.PathParameters = map[string]string{"id": "bk-1"}
				e.Body = `{"score":4}`
				return e
			}(),
			setupMocks: func(rs *mockRatingSaver, rbf *mockRatingByBookingFinder, rpl *mockRatingByProviderLister, bf *mockBookingFinder, pf *mockProviderFinder, ps *mockProviderSaver, ep *mockEventPublisher) {
				bf.FindByIDFunc = func(_ context.Context, _ string) (*booking.Booking, error) {
					return nil, errors.New("dynamo timeout")
				}
			},
			wantStatus:  http.StatusInternalServerError,
			wantErrCode: "INTERNAL_ERROR",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rs := &mockRatingSaver{}
			rbf := &mockRatingByBookingFinder{}
			rpl := &mockRatingByProviderLister{}
			bf := &mockBookingFinder{}
			pf := &mockProviderFinder{}
			ps := &mockProviderSaver{}
			ep := &mockEventPublisher{}

			if tt.setupMocks != nil {
				tt.setupMocks(rs, rbf, rpl, bf, pf, ps, ep)
			}

			uc := ratinguc.NewSubmitRatingUseCase(rs, rbf, rpl, bf, pf, ps, ep)
			h := handler.NewSubmitRatingHandler(uc)

			resp, err := h.Handle(context.Background(), tt.event)

			require.NoError(t, err)
			assert.Equal(t, tt.wantStatus, resp.StatusCode)

			if tt.wantErrCode != "" {
				eb := parseErrorBody(t, resp.Body)
				assert.Equal(t, tt.wantErrCode, eb.Error.Code)
			}

			if tt.checkBody != nil {
				tt.checkBody(t, resp.Body)
			}
		})
	}
}

// ---------------------------------------------------------------------------
// GetRatingHandler tests
// ---------------------------------------------------------------------------

func TestGetRatingHandler(t *testing.T) {
	now := time.Now().UTC()

	existingRating := &rating.Rating{
		BookingID:  "bk-1",
		CustomerID: "user-1",
		ProviderID: "prov-1",
		Score:      5,
		Comment:    "Excellent",
		Tags:       []string{"fast"},
		CreatedAt:  now,
	}

	tests := []struct {
		name        string
		event       *events.APIGatewayProxyRequest
		setupMock   func(rbf *mockRatingByBookingFinder)
		wantStatus  int
		wantErrCode string
		checkBody   func(t *testing.T, body string)
	}{
		{
			name: "success",
			event: func() *events.APIGatewayProxyRequest {
				e := apiEventWithAuth("user-1")
				e.PathParameters = map[string]string{"id": "bk-1"}
				return e
			}(),
			setupMock: func(rbf *mockRatingByBookingFinder) {
				rbf.FindByBookingFunc = func(_ context.Context, _ string) (*rating.Rating, error) {
					return existingRating, nil
				}
			},
			wantStatus: http.StatusOK,
			checkBody: func(t *testing.T, body string) {
				t.Helper()
				var r rating.Rating
				require.NoError(t, json.Unmarshal([]byte(body), &r))
				assert.Equal(t, "bk-1", r.BookingID)
				assert.Equal(t, 5, r.Score)
				assert.Equal(t, "Excellent", r.Comment)
				assert.Equal(t, []string{"fast"}, r.Tags)
			},
		},
		{
			name: "unauthorized - no user ID",
			event: &events.APIGatewayProxyRequest{
				PathParameters: map[string]string{"id": "bk-1"},
			},
			wantStatus:  http.StatusUnauthorized,
			wantErrCode: "UNAUTHORIZED",
		},
		{
			name: "missing booking ID",
			event: func() *events.APIGatewayProxyRequest {
				return apiEventWithAuth("user-1")
			}(),
			wantStatus:  http.StatusBadRequest,
			wantErrCode: "VALIDATION_ERROR",
		},
		{
			name: "use case error - rating not found",
			event: func() *events.APIGatewayProxyRequest {
				e := apiEventWithAuth("user-1")
				e.PathParameters = map[string]string{"id": "bk-missing"}
				return e
			}(),
			setupMock: func(rbf *mockRatingByBookingFinder) {
				rbf.FindByBookingFunc = func(_ context.Context, _ string) (*rating.Rating, error) {
					return nil, nil // not found -> use case returns NotFoundError
				}
			},
			wantStatus:  http.StatusNotFound,
			wantErrCode: "NOT_FOUND",
		},
		{
			name: "use case error - internal error from repo",
			event: func() *events.APIGatewayProxyRequest {
				e := apiEventWithAuth("user-1")
				e.PathParameters = map[string]string{"id": "bk-1"}
				return e
			}(),
			setupMock: func(rbf *mockRatingByBookingFinder) {
				rbf.FindByBookingFunc = func(_ context.Context, _ string) (*rating.Rating, error) {
					return nil, errors.New("dynamo timeout")
				}
			},
			wantStatus:  http.StatusInternalServerError,
			wantErrCode: "INTERNAL_ERROR",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rbf := &mockRatingByBookingFinder{}

			if tt.setupMock != nil {
				tt.setupMock(rbf)
			}

			uc := ratinguc.NewGetBookingRatingUseCase(rbf)
			h := handler.NewGetRatingHandler(uc)

			resp, err := h.Handle(context.Background(), tt.event)

			require.NoError(t, err)
			assert.Equal(t, tt.wantStatus, resp.StatusCode)

			if tt.wantErrCode != "" {
				eb := parseErrorBody(t, resp.Body)
				assert.Equal(t, tt.wantErrCode, eb.Error.Code)
			}

			if tt.checkBody != nil {
				tt.checkBody(t, resp.Body)
			}
		})
	}
}
