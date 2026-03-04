package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"

	"github.com/davidlramirez95/towcommand/internal/domain/booking"
)

// DynamoBookingRepository implements booking persistence against DynamoDB.
type DynamoBookingRepository struct {
	baseRepository
}

// NewBookingRepository creates a new DynamoDB-backed booking repository.
func NewBookingRepository(client DynamoDBAPI, tableName string) *DynamoBookingRepository {
	return &DynamoBookingRepository{
		baseRepository: baseRepository{client: client, tableName: tableName},
	}
}

// Save persists a booking with all key attributes.
// PK: JOB#<bookingId>, SK: DETAILS
// GSI1: USER#<customerId> / JOB#<createdAt>
// GSI2: STATUS#<status> / <createdAt>
func (r *DynamoBookingRepository) Save(ctx context.Context, b *booking.Booking) error {
	item, err := marshalItem(b)
	if err != nil {
		return fmt.Errorf("marshal booking: %w", err)
	}

	item["PK"] = stringAttr(PrefixJob + b.BookingID)
	item["SK"] = stringAttr(SKDetails)
	item["GSI1PK"] = stringAttr(PrefixUser + b.CustomerID)
	item["GSI1SK"] = stringAttr(PrefixJob + formatTime(b.CreatedAt))
	item["GSI2PK"] = stringAttr(PrefixStatus + string(b.Status))
	item["GSI2SK"] = stringAttr(formatTime(b.CreatedAt))
	item["entityType"] = stringAttr("Booking")

	return r.putItem(ctx, item)
}

// FindByID retrieves a booking by its ID. Returns nil if not found.
func (r *DynamoBookingRepository) FindByID(ctx context.Context, bookingID string) (*booking.Booking, error) {
	var b booking.Booking
	found, err := r.getItem(ctx, PrefixJob+bookingID, SKDetails, &b)
	if err != nil {
		return nil, err
	}
	if !found {
		return nil, nil
	}
	return &b, nil
}

// FindByUser lists bookings for a user via GSI1, ordered by creation date descending.
func (r *DynamoBookingRepository) FindByUser(ctx context.Context, userID string, limit int32) ([]booking.Booking, error) {
	keyCond := expression.Key("GSI1PK").Equal(expression.Value(PrefixUser + userID)).
		And(expression.Key("GSI1SK").BeginsWith(PrefixJob))

	expr, err := expression.NewBuilder().WithKeyCondition(keyCond).Build()
	if err != nil {
		return nil, fmt.Errorf("build key expression: %w", err)
	}

	items, err := r.queryItems(ctx, &dynamodb.QueryInput{
		IndexName:                 aws.String(GSI1Name),
		KeyConditionExpression:    expr.KeyCondition(),
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		ScanIndexForward:          aws.Bool(false),
		Limit:                     &limit,
	})
	if err != nil {
		return nil, err
	}
	return unmarshalItems[booking.Booking](items)
}

// UpdateStatus changes a booking's status, updates the GSI2 keys, and records a history entry.
func (r *DynamoBookingRepository) UpdateStatus(ctx context.Context, bookingID string, status booking.BookingStatus, metadata map[string]any) error {
	now := time.Now().UTC()
	nowStr := formatTime(now)

	// Update the booking item.
	err := r.updateItem(ctx, PrefixJob+bookingID, SKDetails, map[string]any{
		"status":    string(status),
		"updatedAt": now,
		"GSI2PK":    PrefixStatus + string(status),
		"GSI2SK":    nowStr,
	})
	if err != nil {
		return err
	}

	// Write a history record: PK=JOB#<id>, SK=STATUS#<timestamp>
	historyItem := map[string]any{
		"entityType": "BookingHistory",
		"status":     string(status),
		"changedAt":  nowStr,
		"metadata":   metadata,
	}
	item, err := marshalItem(historyItem)
	if err != nil {
		return fmt.Errorf("marshal history: %w", err)
	}
	item["PK"] = stringAttr(PrefixJob + bookingID)
	item["SK"] = stringAttr(PrefixStatus + nowStr)

	return r.putItem(ctx, item)
}

// FindByStatus lists bookings with a given status via GSI2, ordered by creation date descending.
func (r *DynamoBookingRepository) FindByStatus(ctx context.Context, status booking.BookingStatus, limit int32) ([]booking.Booking, error) {
	keyCond := expression.Key("GSI2PK").Equal(expression.Value(PrefixStatus + string(status)))

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
	return unmarshalItems[booking.Booking](items)
}
