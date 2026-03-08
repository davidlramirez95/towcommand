package repository

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"

	"github.com/davidlramirez95/towcommand/internal/usecase/port"
)

// Key prefix for push token items.
const (
	PrefixPush = "PUSH#"
	PrefixDev  = "DEV#"
)

// DynamoDBDeleteAPI extends DynamoDBAPI with DeleteItem support.
type DynamoDBDeleteAPI interface {
	DynamoDBAPI
	DeleteItem(ctx context.Context, params *dynamodb.DeleteItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.DeleteItemOutput, error)
}

// DynamoPushRepository implements push token persistence against DynamoDB.
type DynamoPushRepository struct {
	baseRepository
	deleteClient DynamoDBDeleteAPI
}

// NewPushRepository creates a new DynamoDB-backed push token repository.
// The client must implement DynamoDBDeleteAPI (the standard *dynamodb.Client does).
func NewPushRepository(client DynamoDBDeleteAPI, tableName string) *DynamoPushRepository {
	return &DynamoPushRepository{
		baseRepository: baseRepository{client: client, tableName: tableName},
		deleteClient:   client,
	}
}

// Compile-time interface checks.
var (
	_ port.PushTokenRegistrar = (*DynamoPushRepository)(nil)
	_ port.PushTokenFinder    = (*DynamoPushRepository)(nil)
)

// Register persists a push token to DynamoDB.
// PK: PUSH#{userId}, SK: DEV#{deviceId}
// GSI1PK: PUSH#{userId}, GSI1SK: DEV#{createdAt}
func (r *DynamoPushRepository) Register(ctx context.Context, token *port.PushToken) error {
	item, err := marshalItem(token)
	if err != nil {
		return fmt.Errorf("marshal push token: %w", err)
	}

	item["PK"] = stringAttr(PrefixPush + token.UserID)
	item["SK"] = stringAttr(PrefixDev + token.DeviceID)
	item["GSI1PK"] = stringAttr(PrefixPush + token.UserID)
	item["GSI1SK"] = stringAttr(PrefixDev + formatTime(token.CreatedAt))
	item["entityType"] = stringAttr("PushToken")

	return r.putItem(ctx, item)
}

// FindByUserID retrieves all push tokens for a user via GSI1.
func (r *DynamoPushRepository) FindByUserID(ctx context.Context, userID string) ([]port.PushToken, error) {
	keyCond := expression.Key("GSI1PK").Equal(expression.Value(PrefixPush + userID)).
		And(expression.Key("GSI1SK").BeginsWith(PrefixDev))

	expr, err := expression.NewBuilder().WithKeyCondition(keyCond).Build()
	if err != nil {
		return nil, fmt.Errorf("build key expression: %w", err)
	}

	items, err := r.queryItems(ctx, &dynamodb.QueryInput{
		IndexName:                 aws.String(GSI1Name),
		KeyConditionExpression:    expr.KeyCondition(),
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
	})
	if err != nil {
		return nil, err
	}

	return unmarshalItems[port.PushToken](items)
}

// Delete removes a push token by userID and deviceID.
func (r *DynamoPushRepository) Delete(ctx context.Context, userID, deviceID string) error {
	_, err := r.deleteClient.DeleteItem(ctx, &dynamodb.DeleteItemInput{
		TableName: &r.tableName,
		Key:       compositeKey(PrefixPush+userID, PrefixDev+deviceID),
	})
	if err != nil {
		return fmt.Errorf("dynamodb DeleteItem: %w", err)
	}
	return nil
}
