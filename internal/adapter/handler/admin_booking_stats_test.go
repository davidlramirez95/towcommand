package handler_test

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"testing"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/davidlramirez95/towcommand/internal/adapter/handler"
)

// ---------------------------------------------------------------------------
// Mock DynamoDB reader for BookingStatsHandler
// ---------------------------------------------------------------------------

type mockDynamoDBReader struct {
	GetItemFunc func(ctx context.Context, params *dynamodb.GetItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.GetItemOutput, error)
}

func (m *mockDynamoDBReader) GetItem(ctx context.Context, params *dynamodb.GetItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.GetItemOutput, error) {
	if m.GetItemFunc != nil {
		return m.GetItemFunc(ctx, params, optFns...)
	}
	return &dynamodb.GetItemOutput{}, nil
}

// ---------------------------------------------------------------------------
// BookingStatsHandler tests
// ---------------------------------------------------------------------------

func TestBookingStatsHandler(t *testing.T) {
	today := time.Now().UTC().Format("2006-01-02")

	tests := []struct {
		name        string
		event       *events.APIGatewayProxyRequest
		setupMock   func(m *mockDynamoDBReader)
		wantStatus  int
		wantErrCode string
		checkBody   func(t *testing.T, body string)
	}{
		{
			name: "success - with date param",
			event: func() *events.APIGatewayProxyRequest {
				e := apiEventWithAuth("admin-1")
				e.QueryStringParameters = map[string]string{"date": "2026-03-05"}
				return e
			}(),
			setupMock: func(m *mockDynamoDBReader) {
				m.GetItemFunc = func(_ context.Context, params *dynamodb.GetItemInput, _ ...func(*dynamodb.Options)) (*dynamodb.GetItemOutput, error) {
					// Verify correct key was used.
					pk := params.Key["PK"].(*types.AttributeValueMemberS).Value
					assert.Equal(t, "ANALYTICS#DAILY#2026-03-05", pk)

					return &dynamodb.GetItemOutput{
						Item: map[string]types.AttributeValue{
							"PK":               &types.AttributeValueMemberS{Value: "ANALYTICS#DAILY#2026-03-05"},
							"SK":               &types.AttributeValueMemberS{Value: "SUMMARY"},
							"bookingsCreated":  &types.AttributeValueMemberN{Value: "42"},
							"bookingsComplete": &types.AttributeValueMemberN{Value: "30"},
							"sosTriggered":     &types.AttributeValueMemberN{Value: "3"},
						},
					}, nil
				}
			},
			wantStatus: http.StatusOK,
			checkBody: func(t *testing.T, body string) {
				t.Helper()
				var result struct {
					Date  string           `json:"date"`
					Stats map[string]int64 `json:"stats"`
				}
				require.NoError(t, json.Unmarshal([]byte(body), &result))
				assert.Equal(t, "2026-03-05", result.Date)
				assert.Equal(t, int64(42), result.Stats["bookingsCreated"])
				assert.Equal(t, int64(30), result.Stats["bookingsComplete"])
				assert.Equal(t, int64(3), result.Stats["sosTriggered"])
			},
		},
		{
			name: "success - default date (today)",
			event: func() *events.APIGatewayProxyRequest {
				return apiEventWithAuth("admin-1")
			}(),
			setupMock: func(m *mockDynamoDBReader) {
				m.GetItemFunc = func(_ context.Context, params *dynamodb.GetItemInput, _ ...func(*dynamodb.Options)) (*dynamodb.GetItemOutput, error) {
					pk := params.Key["PK"].(*types.AttributeValueMemberS).Value
					assert.Equal(t, "ANALYTICS#DAILY#"+today, pk)
					return &dynamodb.GetItemOutput{
						Item: map[string]types.AttributeValue{
							"PK":              &types.AttributeValueMemberS{Value: "ANALYTICS#DAILY#" + today},
							"SK":              &types.AttributeValueMemberS{Value: "SUMMARY"},
							"bookingsCreated": &types.AttributeValueMemberN{Value: "5"},
						},
					}, nil
				}
			},
			wantStatus: http.StatusOK,
			checkBody: func(t *testing.T, body string) {
				t.Helper()
				var result struct {
					Date  string           `json:"date"`
					Stats map[string]int64 `json:"stats"`
				}
				require.NoError(t, json.Unmarshal([]byte(body), &result))
				assert.Equal(t, today, result.Date)
				assert.Equal(t, int64(5), result.Stats["bookingsCreated"])
			},
		},
		{
			name: "success - no data for date",
			event: func() *events.APIGatewayProxyRequest {
				e := apiEventWithAuth("admin-1")
				e.QueryStringParameters = map[string]string{"date": "2020-01-01"}
				return e
			}(),
			setupMock: func(m *mockDynamoDBReader) {
				m.GetItemFunc = func(_ context.Context, _ *dynamodb.GetItemInput, _ ...func(*dynamodb.Options)) (*dynamodb.GetItemOutput, error) {
					return &dynamodb.GetItemOutput{Item: nil}, nil
				}
			},
			wantStatus: http.StatusOK,
			checkBody: func(t *testing.T, body string) {
				t.Helper()
				var result struct {
					Date  string           `json:"date"`
					Stats map[string]int64 `json:"stats"`
				}
				require.NoError(t, json.Unmarshal([]byte(body), &result))
				assert.Equal(t, "2020-01-01", result.Date)
				assert.Empty(t, result.Stats)
			},
		},
		{
			name:        "unauthorized - no user ID",
			event:       &events.APIGatewayProxyRequest{},
			wantStatus:  http.StatusUnauthorized,
			wantErrCode: "UNAUTHORIZED",
		},
		{
			name: "DynamoDB error",
			event: func() *events.APIGatewayProxyRequest {
				return apiEventWithAuth("admin-1")
			}(),
			setupMock: func(m *mockDynamoDBReader) {
				m.GetItemFunc = func(_ context.Context, _ *dynamodb.GetItemInput, _ ...func(*dynamodb.Options)) (*dynamodb.GetItemOutput, error) {
					return nil, errors.New("dynamo timeout")
				}
			},
			wantStatus:  http.StatusInternalServerError,
			wantErrCode: "INTERNAL_ERROR",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ddb := &mockDynamoDBReader{}

			if tt.setupMock != nil {
				tt.setupMock(ddb)
			}

			h := handler.NewBookingStatsHandler(ddb, "test-table")

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
