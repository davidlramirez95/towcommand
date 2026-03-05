package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"

	"github.com/davidlramirez95/towcommand/internal/domain/safety"
)

const (
	prefixSOS       = "SOS#"
	skMetadata      = "METADATA"
	prefixSOSStatus = "SOSSTATUS#"
)

// DynamoSOSRepository implements SOS alert persistence against DynamoDB.
type DynamoSOSRepository struct {
	baseRepository
}

// NewSOSRepository creates a new DynamoDB-backed SOS repository.
func NewSOSRepository(client DynamoDBAPI, tableName string) *DynamoSOSRepository {
	return &DynamoSOSRepository{
		baseRepository: baseRepository{client: client, tableName: tableName},
	}
}

// Save persists an SOS alert with all key attributes.
// PK: SOS#<alertId>, SK: METADATA
// GSI2PK: SOSSTATUS#ACTIVE (or SOSSTATUS#RESOLVED), GSI2SK: <ISO timestamp>
func (r *DynamoSOSRepository) Save(ctx context.Context, alert *safety.SOSAlert) error {
	item, err := marshalItem(alert)
	if err != nil {
		return fmt.Errorf("marshal SOS alert: %w", err)
	}

	item["PK"] = stringAttr(prefixSOS + alert.AlertID)
	item["SK"] = stringAttr(skMetadata)

	statusKey := prefixSOSStatus + "ACTIVE"
	if alert.Resolved {
		statusKey = prefixSOSStatus + "RESOLVED"
	}
	item["GSI2PK"] = stringAttr(statusKey)
	item["GSI2SK"] = stringAttr(formatTime(alert.Timestamp))
	item["entityType"] = stringAttr("SOSAlert")

	return r.putItem(ctx, item)
}

// FindByID retrieves an SOS alert by its ID. Returns nil if not found.
func (r *DynamoSOSRepository) FindByID(ctx context.Context, alertID string) (*safety.SOSAlert, error) {
	var alert safety.SOSAlert
	found, err := r.getItem(ctx, prefixSOS+alertID, skMetadata, &alert)
	if err != nil {
		return nil, err
	}
	if !found {
		return nil, nil
	}
	return &alert, nil
}

// Resolve marks an SOS alert as resolved by updating its status fields and GSI2 key.
func (r *DynamoSOSRepository) Resolve(ctx context.Context, alertID string, resolvedBy string, resolvedAt time.Time) error {
	return r.updateItem(ctx, prefixSOS+alertID, skMetadata, map[string]any{
		"resolved":   true,
		"resolvedAt": resolvedAt,
		"resolvedBy": resolvedBy,
		"GSI2PK":     prefixSOSStatus + "RESOLVED",
	})
}

// FindActive queries GSI2 for active SOS alerts, ordered by timestamp descending.
func (r *DynamoSOSRepository) FindActive(ctx context.Context, limit int32) ([]safety.SOSAlert, error) {
	keyCond := expression.Key("GSI2PK").Equal(expression.Value(prefixSOSStatus + "ACTIVE"))

	expr, err := expression.NewBuilder().WithKeyCondition(keyCond).Build()
	if err != nil {
		return nil, fmt.Errorf("build key expression: %w", err)
	}

	items, err := r.queryItems(ctx, &dynamodb.QueryInput{
		IndexName:                 aws.String(GSI2Name),
		KeyConditionExpression:    expr.KeyCondition(),
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		ScanIndexForward:          aws.Bool(false),
		Limit:                     &limit,
	})
	if err != nil {
		return nil, err
	}
	return unmarshalItems[safety.SOSAlert](items)
}
