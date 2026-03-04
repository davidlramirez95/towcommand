package repository

import (
	"context"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// mockDynamoDBForAnalytics wraps mock.Mock to satisfy DynamoDBAPI.
type mockDynamoDBForAnalytics struct{ mock.Mock }

func (m *mockDynamoDBForAnalytics) GetItem(ctx context.Context, params *dynamodb.GetItemInput, _ ...func(*dynamodb.Options)) (*dynamodb.GetItemOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*dynamodb.GetItemOutput), args.Error(1)
}

func (m *mockDynamoDBForAnalytics) PutItem(ctx context.Context, params *dynamodb.PutItemInput, _ ...func(*dynamodb.Options)) (*dynamodb.PutItemOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*dynamodb.PutItemOutput), args.Error(1)
}

func (m *mockDynamoDBForAnalytics) UpdateItem(ctx context.Context, params *dynamodb.UpdateItemInput, _ ...func(*dynamodb.Options)) (*dynamodb.UpdateItemOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*dynamodb.UpdateItemOutput), args.Error(1)
}

func (m *mockDynamoDBForAnalytics) Query(ctx context.Context, params *dynamodb.QueryInput, _ ...func(*dynamodb.Options)) (*dynamodb.QueryOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*dynamodb.QueryOutput), args.Error(1)
}

func (m *mockDynamoDBForAnalytics) TransactWriteItems(ctx context.Context, params *dynamodb.TransactWriteItemsInput, _ ...func(*dynamodb.Options)) (*dynamodb.TransactWriteItemsOutput, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*dynamodb.TransactWriteItemsOutput), args.Error(1)
}

func TestIncrementDailyCounter(t *testing.T) {
	ddb := new(mockDynamoDBForAnalytics)
	repo := NewDynamoAnalyticsRepository(ddb, "tc-test-table")

	ddb.On("UpdateItem", mock.Anything, mock.MatchedBy(func(input *dynamodb.UpdateItemInput) bool {
		pk, ok := input.Key["PK"].(*types.AttributeValueMemberS)
		if !ok {
			return false
		}
		sk, ok := input.Key["SK"].(*types.AttributeValueMemberS)
		if !ok {
			return false
		}
		return pk.Value == "ANALYTICS#DAILY#2026-03-04" && sk.Value == "SUMMARY"
	})).Return(&dynamodb.UpdateItemOutput{}, nil)

	err := repo.IncrementDailyCounter(context.Background(), "2026-03-04", "totalBookings", 1)
	require.NoError(t, err)
	ddb.AssertExpectations(t)
}

func TestIncrementDailyCounter_Error(t *testing.T) {
	ddb := new(mockDynamoDBForAnalytics)
	repo := NewDynamoAnalyticsRepository(ddb, "tc-test-table")

	ddb.On("UpdateItem", mock.Anything, mock.Anything).Return(&dynamodb.UpdateItemOutput{}, assert.AnError)

	err := repo.IncrementDailyCounter(context.Background(), "2026-03-04", "totalBookings", 1)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "daily counter")
}

func TestIncrementHeatmapCell(t *testing.T) {
	ddb := new(mockDynamoDBForAnalytics)
	repo := NewDynamoAnalyticsRepository(ddb, "tc-test-table")

	ddb.On("UpdateItem", mock.Anything, mock.MatchedBy(func(input *dynamodb.UpdateItemInput) bool {
		pk, ok := input.Key["PK"].(*types.AttributeValueMemberS)
		if !ok {
			return false
		}
		sk, ok := input.Key["SK"].(*types.AttributeValueMemberS)
		if !ok {
			return false
		}
		return pk.Value == "ANALYTICS#HEATMAP#2026-03-04" && sk.Value == "CELL#14.600,120.984"
	})).Return(&dynamodb.UpdateItemOutput{}, nil)

	err := repo.IncrementHeatmapCell(context.Background(), "2026-03-04", "14.600,120.984", 14.5995, 120.9842)
	require.NoError(t, err)
	ddb.AssertExpectations(t)
}

func TestIncrementHeatmapCell_Error(t *testing.T) {
	ddb := new(mockDynamoDBForAnalytics)
	repo := NewDynamoAnalyticsRepository(ddb, "tc-test-table")

	ddb.On("UpdateItem", mock.Anything, mock.Anything).Return(&dynamodb.UpdateItemOutput{}, assert.AnError)

	err := repo.IncrementHeatmapCell(context.Background(), "2026-03-04", "14.600,120.984", 14.5995, 120.9842)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "heatmap")
}

func TestIncrementProviderCounter(t *testing.T) {
	ddb := new(mockDynamoDBForAnalytics)
	repo := NewDynamoAnalyticsRepository(ddb, "tc-test-table")

	ddb.On("UpdateItem", mock.Anything, mock.MatchedBy(func(input *dynamodb.UpdateItemInput) bool {
		pk, ok := input.Key["PK"].(*types.AttributeValueMemberS)
		if !ok {
			return false
		}
		sk, ok := input.Key["SK"].(*types.AttributeValueMemberS)
		if !ok {
			return false
		}
		return pk.Value == "ANALYTICS#PROVIDER#prov-1" && sk.Value == "DAILY#2026-03-04"
	})).Return(&dynamodb.UpdateItemOutput{}, nil)

	err := repo.IncrementProviderCounter(context.Background(), "prov-1", "2026-03-04", "completedJobs", 1)
	require.NoError(t, err)
	ddb.AssertExpectations(t)
}

func TestIncrementProviderCounter_Error(t *testing.T) {
	ddb := new(mockDynamoDBForAnalytics)
	repo := NewDynamoAnalyticsRepository(ddb, "tc-test-table")

	ddb.On("UpdateItem", mock.Anything, mock.Anything).Return(&dynamodb.UpdateItemOutput{}, assert.AnError)

	err := repo.IncrementProviderCounter(context.Background(), "prov-1", "2026-03-04", "completedJobs", 1)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "provider counter")
}
