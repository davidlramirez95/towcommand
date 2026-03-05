package repository

import (
	"context"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/davidlramirez95/towcommand/internal/domain/safety"
)

// mockDynamoDBForSOS wraps mock.Mock to satisfy DynamoDBAPI.
type mockDynamoDBForSOS struct{ mock.Mock }

func (m *mockDynamoDBForSOS) GetItem(ctx context.Context, params *dynamodb.GetItemInput, _ ...func(*dynamodb.Options)) (*dynamodb.GetItemOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*dynamodb.GetItemOutput), args.Error(1)
}

func (m *mockDynamoDBForSOS) PutItem(ctx context.Context, params *dynamodb.PutItemInput, _ ...func(*dynamodb.Options)) (*dynamodb.PutItemOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*dynamodb.PutItemOutput), args.Error(1)
}

func (m *mockDynamoDBForSOS) UpdateItem(ctx context.Context, params *dynamodb.UpdateItemInput, _ ...func(*dynamodb.Options)) (*dynamodb.UpdateItemOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*dynamodb.UpdateItemOutput), args.Error(1)
}

func (m *mockDynamoDBForSOS) Query(ctx context.Context, params *dynamodb.QueryInput, _ ...func(*dynamodb.Options)) (*dynamodb.QueryOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*dynamodb.QueryOutput), args.Error(1)
}

func (m *mockDynamoDBForSOS) TransactWriteItems(ctx context.Context, params *dynamodb.TransactWriteItemsInput, _ ...func(*dynamodb.Options)) (*dynamodb.TransactWriteItemsOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*dynamodb.TransactWriteItemsOutput), args.Error(1)
}

func TestSOSRepository_Save(t *testing.T) {
	tests := []struct {
		name      string
		alert     *safety.SOSAlert
		ddbErr    error
		wantErr   bool
		checkItem func(t *testing.T, input *dynamodb.PutItemInput)
	}{
		{
			name: "saves active alert with correct keys",
			alert: &safety.SOSAlert{
				AlertID:     "SOS-2026-abc123",
				BookingID:   "BK-001",
				TriggeredBy: "USR-001",
				TriggerType: safety.TriggerTypeButton,
				Lat:         14.5995,
				Lng:         120.9842,
				Resolved:    false,
				Timestamp:   time.Date(2026, 3, 5, 10, 0, 0, 0, time.UTC),
			},
			checkItem: func(t *testing.T, input *dynamodb.PutItemInput) {
				t.Helper()
				pk := input.Item["PK"].(*types.AttributeValueMemberS).Value
				sk := input.Item["SK"].(*types.AttributeValueMemberS).Value
				gsi2pk := input.Item["GSI2PK"].(*types.AttributeValueMemberS).Value
				entityType := input.Item["entityType"].(*types.AttributeValueMemberS).Value

				assert.Equal(t, "SOS#SOS-2026-abc123", pk)
				assert.Equal(t, "METADATA", sk)
				assert.Equal(t, "SOSSTATUS#ACTIVE", gsi2pk)
				assert.Equal(t, "SOSAlert", entityType)
			},
		},
		{
			name: "saves resolved alert with correct GSI2PK",
			alert: &safety.SOSAlert{
				AlertID:   "SOS-2026-def456",
				Resolved:  true,
				Timestamp: time.Date(2026, 3, 5, 12, 0, 0, 0, time.UTC),
			},
			checkItem: func(t *testing.T, input *dynamodb.PutItemInput) {
				t.Helper()
				gsi2pk := input.Item["GSI2PK"].(*types.AttributeValueMemberS).Value
				assert.Equal(t, "SOSSTATUS#RESOLVED", gsi2pk)
			},
		},
		{
			name: "returns error on DynamoDB failure",
			alert: &safety.SOSAlert{
				AlertID:   "SOS-2026-err",
				Timestamp: time.Date(2026, 3, 5, 10, 0, 0, 0, time.UTC),
			},
			ddbErr:  assert.AnError,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ddb := new(mockDynamoDBForSOS)
			repo := NewSOSRepository(ddb, "tc-test-table")

			if tt.checkItem != nil {
				ddb.On("PutItem", mock.Anything, mock.MatchedBy(func(input *dynamodb.PutItemInput) bool {
					tt.checkItem(t, input)
					return true
				})).Return(&dynamodb.PutItemOutput{}, nil)
			} else {
				ddb.On("PutItem", mock.Anything, mock.Anything).Return(&dynamodb.PutItemOutput{}, tt.ddbErr)
			}

			err := repo.Save(context.Background(), tt.alert)

			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)
			ddb.AssertExpectations(t)
		})
	}
}

func TestSOSRepository_FindByID(t *testing.T) {
	tests := []struct {
		name    string
		alertID string
		item    map[string]types.AttributeValue
		ddbErr  error
		wantNil bool
		wantErr bool
	}{
		{
			name:    "found",
			alertID: "SOS-2026-abc123",
			item: map[string]types.AttributeValue{
				"alertId":     &types.AttributeValueMemberS{Value: "SOS-2026-abc123"},
				"bookingId":   &types.AttributeValueMemberS{Value: "BK-001"},
				"triggeredBy": &types.AttributeValueMemberS{Value: "USR-001"},
				"triggerType": &types.AttributeValueMemberS{Value: "BUTTON"},
				"lat":         &types.AttributeValueMemberN{Value: "14.5995"},
				"lng":         &types.AttributeValueMemberN{Value: "120.9842"},
				"resolved":    &types.AttributeValueMemberBOOL{Value: false},
			},
		},
		{
			name:    "not found",
			alertID: "SOS-MISSING",
			item:    nil,
			wantNil: true,
		},
		{
			name:    "DynamoDB error",
			alertID: "SOS-ERR",
			ddbErr:  assert.AnError,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ddb := new(mockDynamoDBForSOS)
			repo := NewSOSRepository(ddb, "tc-test-table")

			ddb.On("GetItem", mock.Anything, mock.MatchedBy(func(input *dynamodb.GetItemInput) bool {
				pk := input.Key["PK"].(*types.AttributeValueMemberS).Value
				sk := input.Key["SK"].(*types.AttributeValueMemberS).Value
				return pk == "SOS#"+tt.alertID && sk == "METADATA"
			})).Return(&dynamodb.GetItemOutput{Item: tt.item}, tt.ddbErr)

			got, err := repo.FindByID(context.Background(), tt.alertID)

			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)
			if tt.wantNil {
				assert.Nil(t, got)
			} else {
				require.NotNil(t, got)
				assert.Equal(t, tt.alertID, got.AlertID)
			}
			ddb.AssertExpectations(t)
		})
	}
}

func TestSOSRepository_Resolve(t *testing.T) {
	tests := []struct {
		name    string
		ddbErr  error
		wantErr bool
	}{
		{
			name: "success",
		},
		{
			name:    "DynamoDB error",
			ddbErr:  assert.AnError,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ddb := new(mockDynamoDBForSOS)
			repo := NewSOSRepository(ddb, "tc-test-table")

			resolvedAt := time.Date(2026, 3, 5, 14, 0, 0, 0, time.UTC)

			ddb.On("UpdateItem", mock.Anything, mock.MatchedBy(func(input *dynamodb.UpdateItemInput) bool {
				pk := input.Key["PK"].(*types.AttributeValueMemberS).Value
				sk := input.Key["SK"].(*types.AttributeValueMemberS).Value
				return pk == "SOS#SOS-2026-abc" && sk == "METADATA"
			})).Return(&dynamodb.UpdateItemOutput{}, tt.ddbErr)

			err := repo.Resolve(context.Background(), "SOS-2026-abc", "ADMIN-001", resolvedAt)

			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)
			ddb.AssertExpectations(t)
		})
	}
}

func TestSOSRepository_FindActive(t *testing.T) {
	tests := []struct {
		name    string
		items   []map[string]types.AttributeValue
		ddbErr  error
		want    int
		wantErr bool
	}{
		{
			name: "returns active alerts",
			items: []map[string]types.AttributeValue{
				{
					"alertId":     &types.AttributeValueMemberS{Value: "SOS-1"},
					"bookingId":   &types.AttributeValueMemberS{Value: "BK-1"},
					"triggeredBy": &types.AttributeValueMemberS{Value: "USR-1"},
					"triggerType": &types.AttributeValueMemberS{Value: "BUTTON"},
					"lat":         &types.AttributeValueMemberN{Value: "14.5"},
					"lng":         &types.AttributeValueMemberN{Value: "121.0"},
					"resolved":    &types.AttributeValueMemberBOOL{Value: false},
				},
				{
					"alertId":     &types.AttributeValueMemberS{Value: "SOS-2"},
					"bookingId":   &types.AttributeValueMemberS{Value: "BK-2"},
					"triggeredBy": &types.AttributeValueMemberS{Value: "USR-2"},
					"triggerType": &types.AttributeValueMemberS{Value: "SHAKE"},
					"lat":         &types.AttributeValueMemberN{Value: "14.6"},
					"lng":         &types.AttributeValueMemberN{Value: "121.1"},
					"resolved":    &types.AttributeValueMemberBOOL{Value: false},
				},
			},
			want: 2,
		},
		{
			name:  "empty result",
			items: []map[string]types.AttributeValue{},
			want:  0,
		},
		{
			name:    "DynamoDB error",
			ddbErr:  assert.AnError,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ddb := new(mockDynamoDBForSOS)
			repo := NewSOSRepository(ddb, "tc-test-table")

			ddb.On("Query", mock.Anything, mock.MatchedBy(func(input *dynamodb.QueryInput) bool {
				return *input.IndexName == GSI2Name && *input.ScanIndexForward == false
			})).Return(&dynamodb.QueryOutput{Items: tt.items}, tt.ddbErr)

			got, err := repo.FindActive(context.Background(), 50)

			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.Len(t, got, tt.want)
			ddb.AssertExpectations(t)
		})
	}
}
