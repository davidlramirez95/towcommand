package repository

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"

	"github.com/davidlramirez95/towcommand/internal/domain/evidence"
)

// DynamoEvidenceRepository implements evidence/condition report persistence against DynamoDB.
type DynamoEvidenceRepository struct {
	baseRepository
}

// NewEvidenceRepository creates a new DynamoDB-backed evidence repository.
func NewEvidenceRepository(client DynamoDBAPI, tableName string) *DynamoEvidenceRepository {
	return &DynamoEvidenceRepository{
		baseRepository: baseRepository{client: client, tableName: tableName},
	}
}

// Save persists a condition report.
// PK: JOB#<bookingId>, SK: EVIDENCE#<reportId>
func (r *DynamoEvidenceRepository) Save(ctx context.Context, cr *evidence.ConditionReport) error {
	item, err := marshalItem(cr)
	if err != nil {
		return fmt.Errorf("marshal condition report: %w", err)
	}

	item["PK"] = stringAttr(PrefixJob + cr.BookingID)
	item["SK"] = stringAttr(PrefixEvidence + cr.ReportID)
	item["entityType"] = stringAttr("ConditionReport")

	return r.putItem(ctx, item)
}

// FindByBooking lists condition reports for a booking.
func (r *DynamoEvidenceRepository) FindByBooking(ctx context.Context, bookingID string) ([]evidence.ConditionReport, error) {
	keyCond := expression.Key("PK").Equal(expression.Value(PrefixJob + bookingID)).
		And(expression.Key("SK").BeginsWith(PrefixEvidence))

	expr, err := expression.NewBuilder().WithKeyCondition(keyCond).Build()
	if err != nil {
		return nil, fmt.Errorf("build key expression: %w", err)
	}

	items, err := r.queryItems(ctx, &dynamodb.QueryInput{
		KeyConditionExpression:    expr.KeyCondition(),
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
	})
	if err != nil {
		return nil, err
	}
	return unmarshalItems[evidence.ConditionReport](items)
}

// AddMediaItem persists an individual media item under the booking's partition.
// PK: JOB#<bookingId>, SK: MEDIA#<mediaId>
func (r *DynamoEvidenceRepository) AddMediaItem(ctx context.Context, bookingID string, mi *evidence.MediaItem) error {
	item, err := marshalItem(mi)
	if err != nil {
		return fmt.Errorf("marshal media item: %w", err)
	}

	item["PK"] = stringAttr(PrefixJob + bookingID)
	item["SK"] = stringAttr(PrefixMedia + mi.MediaID)
	item["entityType"] = stringAttr("MediaItem")

	return r.putItem(ctx, item)
}
