package repository

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"

	"github.com/davidlramirez95/towcommand/internal/domain/booking"
)

const (
	chatMessagePrefix = "CHAT#MESSAGE#"
)

// DynamoChatRepository implements chat message persistence against DynamoDB.
// PK: JOB#<bookingId>, SK: CHAT#MESSAGE#<timestamp>
type DynamoChatRepository struct {
	baseRepository
}

// NewChatRepository creates a new DynamoDB-backed chat repository.
func NewChatRepository(client DynamoDBAPI, tableName string) *DynamoChatRepository {
	return &DynamoChatRepository{
		baseRepository: baseRepository{client: client, tableName: tableName},
	}
}

// Save persists a chat message under the booking's partition key.
func (r *DynamoChatRepository) Save(ctx context.Context, msg *booking.ChatMessage) error {
	item, err := marshalItem(msg)
	if err != nil {
		return fmt.Errorf("marshal chat message: %w", err)
	}

	item["PK"] = stringAttr(PrefixJob + msg.BookingID)
	item["SK"] = stringAttr(chatMessagePrefix + formatTime(msg.CreatedAt))
	item["entityType"] = stringAttr("ChatMessage")

	return r.putItem(ctx, item)
}

// FindByBooking retrieves chat messages for a booking, ordered by timestamp ascending.
func (r *DynamoChatRepository) FindByBooking(ctx context.Context, bookingID string, limit int32) ([]booking.ChatMessage, error) {
	keyCond := expression.Key("PK").Equal(expression.Value(PrefixJob + bookingID)).
		And(expression.Key("SK").BeginsWith(chatMessagePrefix))

	expr, err := expression.NewBuilder().WithKeyCondition(keyCond).Build()
	if err != nil {
		return nil, fmt.Errorf("build key expression: %w", err)
	}

	items, err := r.queryItems(ctx, &dynamodb.QueryInput{
		KeyConditionExpression:    expr.KeyCondition(),
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		ScanIndexForward:          aws.Bool(true),
		Limit:                     &limit,
	})
	if err != nil {
		return nil, err
	}
	return unmarshalItems[booking.ChatMessage](items)
}
