package repository

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"

	"github.com/davidlramirez95/towcommand/internal/usecase/analytics"
)

// DynamoAnalyticsRepository implements analytics.AnalyticsRecorder using DynamoDB
// atomic counter operations (UpdateItem with ADD).
type DynamoAnalyticsRepository struct {
	baseRepository
}

// NewDynamoAnalyticsRepository creates a new analytics repository.
func NewDynamoAnalyticsRepository(client DynamoDBAPI, tableName string) *DynamoAnalyticsRepository {
	return &DynamoAnalyticsRepository{
		baseRepository: baseRepository{
			client:    client,
			tableName: tableName,
		},
	}
}

// Compile-time interface check.
var _ analytics.AnalyticsRecorder = (*DynamoAnalyticsRepository)(nil)

// IncrementDailyCounter atomically increments a daily analytics counter.
// Key: PK=ANALYTICS#DAILY#{date}, SK=SUMMARY.
func (r *DynamoAnalyticsRepository) IncrementDailyCounter(ctx context.Context, date, field string, delta int64) error {
	pk := "ANALYTICS#DAILY#" + date
	sk := "SUMMARY"

	_, err := r.client.UpdateItem(ctx, &dynamodb.UpdateItemInput{
		TableName:        &r.tableName,
		Key:              compositeKey(pk, sk),
		UpdateExpression: aws.String("ADD #field :delta"),
		ExpressionAttributeNames: map[string]string{
			"#field": field,
		},
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":delta": &types.AttributeValueMemberN{Value: fmt.Sprintf("%d", delta)},
		},
	})
	if err != nil {
		return fmt.Errorf("dynamodb UpdateItem (daily counter %s/%s): %w", date, field, err)
	}
	return nil
}

// IncrementHeatmapCell atomically increments the demand count for a geohash cell.
// Key: PK=ANALYTICS#HEATMAP#{date}, SK=CELL#{geohash}.
// Also stores the lat/lng for display purposes using SET (only on first write).
func (r *DynamoAnalyticsRepository) IncrementHeatmapCell(ctx context.Context, date, geohash string, lat, lng float64) error {
	pk := "ANALYTICS#HEATMAP#" + date
	sk := "CELL#" + geohash

	_, err := r.client.UpdateItem(ctx, &dynamodb.UpdateItemInput{
		TableName:        &r.tableName,
		Key:              compositeKey(pk, sk),
		UpdateExpression: aws.String("ADD #demand :one SET #lat = if_not_exists(#lat, :lat), #lng = if_not_exists(#lng, :lng)"),
		ExpressionAttributeNames: map[string]string{
			"#demand": "demandCount",
			"#lat":    "lat",
			"#lng":    "lng",
		},
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":one": &types.AttributeValueMemberN{Value: "1"},
			":lat": &types.AttributeValueMemberN{Value: fmt.Sprintf("%f", lat)},
			":lng": &types.AttributeValueMemberN{Value: fmt.Sprintf("%f", lng)},
		},
	})
	if err != nil {
		return fmt.Errorf("dynamodb UpdateItem (heatmap %s/%s): %w", date, geohash, err)
	}
	return nil
}

// IncrementProviderCounter atomically increments a provider-level daily counter.
// Key: PK=ANALYTICS#PROVIDER#{providerId}, SK=DAILY#{date}.
func (r *DynamoAnalyticsRepository) IncrementProviderCounter(ctx context.Context, providerID, date, field string, delta int64) error {
	pk := "ANALYTICS#PROVIDER#" + providerID
	sk := "DAILY#" + date

	_, err := r.client.UpdateItem(ctx, &dynamodb.UpdateItemInput{
		TableName:        &r.tableName,
		Key:              compositeKey(pk, sk),
		UpdateExpression: aws.String("ADD #field :delta"),
		ExpressionAttributeNames: map[string]string{
			"#field": field,
		},
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":delta": &types.AttributeValueMemberN{Value: fmt.Sprintf("%d", delta)},
		},
	})
	if err != nil {
		return fmt.Errorf("dynamodb UpdateItem (provider counter %s/%s/%s): %w", providerID, date, field, err)
	}
	return nil
}
