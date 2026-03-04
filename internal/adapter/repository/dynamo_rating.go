package repository

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"

	"github.com/davidlramirez95/towcommand/internal/domain/rating"
)

// DynamoRatingRepository implements rating persistence against DynamoDB.
type DynamoRatingRepository struct {
	baseRepository
}

// NewRatingRepository creates a new DynamoDB-backed rating repository.
func NewRatingRepository(client DynamoDBAPI, tableName string) *DynamoRatingRepository {
	return &DynamoRatingRepository{
		baseRepository: baseRepository{client: client, tableName: tableName},
	}
}

// Save persists a rating.
// PK: JOB#<bookingId>, SK: RATING
// GSI1: PROV#<providerId> / RATE#<createdAt>
func (r *DynamoRatingRepository) Save(ctx context.Context, rt *rating.Rating) error {
	item, err := marshalItem(rt)
	if err != nil {
		return fmt.Errorf("marshal rating: %w", err)
	}

	item["PK"] = stringAttr(PrefixJob + rt.BookingID)
	item["SK"] = stringAttr(SKRating)
	item["GSI1PK"] = stringAttr(PrefixProvider + rt.ProviderID)
	item["GSI1SK"] = stringAttr(PrefixRate + formatTime(rt.CreatedAt))
	item["entityType"] = stringAttr("Rating")

	return r.putItem(ctx, item)
}

// FindByBooking retrieves the rating for a booking. Returns nil if not found.
func (r *DynamoRatingRepository) FindByBooking(ctx context.Context, bookingID string) (*rating.Rating, error) {
	var rt rating.Rating
	found, err := r.getItem(ctx, PrefixJob+bookingID, SKRating, &rt)
	if err != nil {
		return nil, err
	}
	if !found {
		return nil, nil
	}
	return &rt, nil
}

// FindByProvider lists ratings for a provider via GSI1, ordered by date descending.
func (r *DynamoRatingRepository) FindByProvider(ctx context.Context, providerID string, limit int32) ([]rating.Rating, error) {
	keyCond := expression.Key("GSI1PK").Equal(expression.Value(PrefixProvider + providerID)).
		And(expression.Key("GSI1SK").BeginsWith(PrefixRate))

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
	return unmarshalItems[rating.Rating](items)
}
