package handler_test

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"testing"

	"github.com/aws/aws-lambda-go/events"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/davidlramirez95/towcommand/internal/adapter/handler"
	"github.com/davidlramirez95/towcommand/internal/domain/booking"
	"github.com/davidlramirez95/towcommand/internal/domain/evidence"
	evidenceuc "github.com/davidlramirez95/towcommand/internal/usecase/evidence"
	"github.com/davidlramirez95/towcommand/internal/usecase/port"
)

// ---------------------------------------------------------------------------
// Mock implementations for condition report use case dependencies
// ---------------------------------------------------------------------------

type mockEvidenceBookingFinder struct {
	FindByIDFunc func(ctx context.Context, bookingID string) (*booking.Booking, error)
}

func (m *mockEvidenceBookingFinder) FindByID(ctx context.Context, bookingID string) (*booking.Booking, error) {
	if m.FindByIDFunc != nil {
		return m.FindByIDFunc(ctx, bookingID)
	}
	return nil, nil
}

type mockEvidenceSaver struct {
	SaveFunc func(ctx context.Context, r *evidence.ConditionReport) error
}

func (m *mockEvidenceSaver) Save(ctx context.Context, r *evidence.ConditionReport) error {
	if m.SaveFunc != nil {
		return m.SaveFunc(ctx, r)
	}
	return nil
}

type mockEvidenceEventPublisher struct {
	PublishFunc func(ctx context.Context, source, detailType string, detail any, actor *port.Actor) error
}

func (m *mockEvidenceEventPublisher) Publish(ctx context.Context, source, detailType string, detail any, actor *port.Actor) error {
	if m.PublishFunc != nil {
		return m.PublishFunc(ctx, source, detailType, detail, actor)
	}
	return nil
}

// ---------------------------------------------------------------------------
// CreateConditionReportHandler tests
// ---------------------------------------------------------------------------

func TestCreateConditionReportHandler(t *testing.T) {
	validBooking := &booking.Booking{
		BookingID:  "bk-1",
		CustomerID: "cust-1",
		ProviderID: "prov-1",
		Status:     booking.BookingStatusConditionReport,
	}

	tests := []struct {
		name        string
		event       *events.APIGatewayProxyRequest
		setupMocks  func(bf *mockEvidenceBookingFinder, es *mockEvidenceSaver, ep *mockEvidenceEventPublisher)
		wantStatus  int
		wantErrCode string
		checkBody   func(t *testing.T, body string)
	}{
		{
			name: "success - pickup phase",
			event: func() *events.APIGatewayProxyRequest {
				e := apiEventWithAuth("prov-1")
				e.PathParameters = map[string]string{"id": "bk-1"}
				e.Body = `{"phase":"pickup","notes":"Minor scratch on left door"}`
				return e
			}(),
			setupMocks: func(bf *mockEvidenceBookingFinder, es *mockEvidenceSaver, ep *mockEvidenceEventPublisher) {
				bf.FindByIDFunc = func(_ context.Context, _ string) (*booking.Booking, error) {
					return validBooking, nil
				}
			},
			wantStatus: http.StatusCreated,
			checkBody: func(t *testing.T, body string) {
				t.Helper()
				var r evidence.ConditionReport
				require.NoError(t, json.Unmarshal([]byte(body), &r))
				assert.Equal(t, "bk-1", r.BookingID)
				assert.Equal(t, "prov-1", r.ProviderID)
				assert.Equal(t, "pickup", r.Phase)
				assert.Equal(t, "Minor scratch on left door", r.Notes)
				assert.NotEmpty(t, r.ReportID)
			},
		},
		{
			name: "success - dropoff phase with no notes",
			event: func() *events.APIGatewayProxyRequest {
				e := apiEventWithAuth("prov-1")
				e.PathParameters = map[string]string{"id": "bk-1"}
				e.Body = `{"phase":"dropoff"}`
				return e
			}(),
			setupMocks: func(bf *mockEvidenceBookingFinder, es *mockEvidenceSaver, ep *mockEvidenceEventPublisher) {
				bf.FindByIDFunc = func(_ context.Context, _ string) (*booking.Booking, error) {
					return validBooking, nil
				}
			},
			wantStatus: http.StatusCreated,
			checkBody: func(t *testing.T, body string) {
				t.Helper()
				var r evidence.ConditionReport
				require.NoError(t, json.Unmarshal([]byte(body), &r))
				assert.Equal(t, "dropoff", r.Phase)
				assert.Empty(t, r.Notes)
			},
		},
		{
			name: "unauthorized - no user ID",
			event: &events.APIGatewayProxyRequest{
				PathParameters: map[string]string{"id": "bk-1"},
				Body:           `{"phase":"pickup"}`,
			},
			wantStatus:  http.StatusUnauthorized,
			wantErrCode: "UNAUTHORIZED",
		},
		{
			name: "missing booking ID path param",
			event: func() *events.APIGatewayProxyRequest {
				e := apiEventWithAuth("prov-1")
				e.Body = `{"phase":"pickup"}`
				return e
			}(),
			wantStatus:  http.StatusBadRequest,
			wantErrCode: "VALIDATION_ERROR",
		},
		{
			name: "invalid JSON body",
			event: func() *events.APIGatewayProxyRequest {
				e := apiEventWithAuth("prov-1")
				e.PathParameters = map[string]string{"id": "bk-1"}
				e.Body = `{not json}`
				return e
			}(),
			wantStatus:  http.StatusBadRequest,
			wantErrCode: "VALIDATION_ERROR",
		},
		{
			name: "invalid phase value",
			event: func() *events.APIGatewayProxyRequest {
				e := apiEventWithAuth("prov-1")
				e.PathParameters = map[string]string{"id": "bk-1"}
				e.Body = `{"phase":"invalid"}`
				return e
			}(),
			wantStatus:  http.StatusBadRequest,
			wantErrCode: "VALIDATION_ERROR",
		},
		{
			name: "missing required phase",
			event: func() *events.APIGatewayProxyRequest {
				e := apiEventWithAuth("prov-1")
				e.PathParameters = map[string]string{"id": "bk-1"}
				e.Body = `{"notes":"some notes"}`
				return e
			}(),
			wantStatus:  http.StatusBadRequest,
			wantErrCode: "VALIDATION_ERROR",
		},
		{
			name: "use case error - booking not found",
			event: func() *events.APIGatewayProxyRequest {
				e := apiEventWithAuth("prov-1")
				e.PathParameters = map[string]string{"id": "bk-missing"}
				e.Body = `{"phase":"pickup"}`
				return e
			}(),
			setupMocks: func(bf *mockEvidenceBookingFinder, es *mockEvidenceSaver, ep *mockEvidenceEventPublisher) {
				bf.FindByIDFunc = func(_ context.Context, _ string) (*booking.Booking, error) {
					return nil, nil // not found
				}
			},
			wantStatus:  http.StatusNotFound,
			wantErrCode: "NOT_FOUND",
		},
		{
			name: "use case error - internal error from repo",
			event: func() *events.APIGatewayProxyRequest {
				e := apiEventWithAuth("prov-1")
				e.PathParameters = map[string]string{"id": "bk-1"}
				e.Body = `{"phase":"pickup"}`
				return e
			}(),
			setupMocks: func(bf *mockEvidenceBookingFinder, es *mockEvidenceSaver, ep *mockEvidenceEventPublisher) {
				bf.FindByIDFunc = func(_ context.Context, _ string) (*booking.Booking, error) {
					return nil, errors.New("dynamo timeout")
				}
			},
			wantStatus:  http.StatusInternalServerError,
			wantErrCode: "INTERNAL_ERROR",
		},
		{
			name: "use case error - evidence save fails",
			event: func() *events.APIGatewayProxyRequest {
				e := apiEventWithAuth("prov-1")
				e.PathParameters = map[string]string{"id": "bk-1"}
				e.Body = `{"phase":"pickup"}`
				return e
			}(),
			setupMocks: func(bf *mockEvidenceBookingFinder, es *mockEvidenceSaver, ep *mockEvidenceEventPublisher) {
				bf.FindByIDFunc = func(_ context.Context, _ string) (*booking.Booking, error) {
					return validBooking, nil
				}
				es.SaveFunc = func(_ context.Context, _ *evidence.ConditionReport) error {
					return errors.New("dynamo put failed")
				}
			},
			wantStatus:  http.StatusInternalServerError,
			wantErrCode: "INTERNAL_ERROR",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bf := &mockEvidenceBookingFinder{}
			es := &mockEvidenceSaver{}
			ep := &mockEvidenceEventPublisher{}

			if tt.setupMocks != nil {
				tt.setupMocks(bf, es, ep)
			}

			uc := evidenceuc.NewCreateConditionReportUseCase(bf, es, ep)
			h := handler.NewCreateConditionReportHandler(uc)

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
