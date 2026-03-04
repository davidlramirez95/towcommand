package repository

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"

	"github.com/davidlramirez95/towcommand/internal/domain/user"
)

// DynamoUserRepository implements user persistence against DynamoDB.
type DynamoUserRepository struct {
	baseRepository
}

// NewUserRepository creates a new DynamoDB-backed user repository.
func NewUserRepository(client DynamoDBAPI, tableName string) *DynamoUserRepository {
	return &DynamoUserRepository{
		baseRepository: baseRepository{client: client, tableName: tableName},
	}
}

// Save persists a user with all key attributes.
// PK: USER#<userId>, SK: PROFILE
// GSI1: EMAIL#<email> / USER
// GSI5: PHONE#<phone> / PROFILE
func (r *DynamoUserRepository) Save(ctx context.Context, u *user.User) error {
	item, err := marshalItem(u)
	if err != nil {
		return fmt.Errorf("marshal user: %w", err)
	}

	item["PK"] = stringAttr(PrefixUser + u.UserID)
	item["SK"] = stringAttr(SKProfile)
	item["GSI1PK"] = stringAttr(PrefixEmail + u.Email)
	item["GSI1SK"] = stringAttr("USER")
	item["GSI5PK"] = stringAttr(PrefixPhone + u.Phone)
	item["GSI5SK"] = stringAttr(SKProfile)
	item["entityType"] = stringAttr("User")

	return r.putItem(ctx, item)
}

// FindByID retrieves a user by their ID. Returns nil if not found.
func (r *DynamoUserRepository) FindByID(ctx context.Context, userID string) (*user.User, error) {
	var u user.User
	found, err := r.getItem(ctx, PrefixUser+userID, SKProfile, &u)
	if err != nil {
		return nil, err
	}
	if !found {
		return nil, nil
	}
	return &u, nil
}

// FindByEmail retrieves a user by their email address via GSI1. Returns nil if not found.
func (r *DynamoUserRepository) FindByEmail(ctx context.Context, email string) (*user.User, error) {
	keyCond := expression.Key("GSI1PK").Equal(expression.Value(PrefixEmail + email)).
		And(expression.Key("GSI1SK").Equal(expression.Value("USER")))

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
	if len(items) == 0 {
		return nil, nil
	}

	var u user.User
	if err := unmarshalItem(items[0], &u); err != nil {
		return nil, fmt.Errorf("unmarshal user: %w", err)
	}
	return &u, nil
}

// FindByPhone retrieves a user by their phone number via GSI5. Returns nil if not found.
func (r *DynamoUserRepository) FindByPhone(ctx context.Context, phone string) (*user.User, error) {
	keyCond := expression.Key("GSI5PK").Equal(expression.Value(PrefixPhone + phone)).
		And(expression.Key("GSI5SK").Equal(expression.Value(SKProfile)))

	expr, err := expression.NewBuilder().WithKeyCondition(keyCond).Build()
	if err != nil {
		return nil, fmt.Errorf("build key expression: %w", err)
	}

	items, err := r.queryItems(ctx, &dynamodb.QueryInput{
		IndexName:                 aws.String(GSI5Name),
		KeyConditionExpression:    expr.KeyCondition(),
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
	})
	if err != nil {
		return nil, err
	}
	if len(items) == 0 {
		return nil, nil
	}

	var u user.User
	if err := unmarshalItem(items[0], &u); err != nil {
		return nil, fmt.Errorf("unmarshal user: %w", err)
	}
	return &u, nil
}

// AddVehicle persists a vehicle under its owner's partition.
// PK: USER#<userId>, SK: VEH#<vehicleId>
func (r *DynamoUserRepository) AddVehicle(ctx context.Context, v *user.UserVehicle) error {
	item, err := marshalItem(v)
	if err != nil {
		return fmt.Errorf("marshal vehicle: %w", err)
	}

	item["PK"] = stringAttr(PrefixUser + v.UserID)
	item["SK"] = stringAttr(PrefixVehicle + v.VehicleID)
	item["entityType"] = stringAttr("UserVehicle")

	return r.putItem(ctx, item)
}

// GetVehicles lists all vehicles for a user.
func (r *DynamoUserRepository) GetVehicles(ctx context.Context, userID string) ([]user.UserVehicle, error) {
	keyCond := expression.Key("PK").Equal(expression.Value(PrefixUser + userID)).
		And(expression.Key("SK").BeginsWith(PrefixVehicle))

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
	return unmarshalItems[user.UserVehicle](items)
}
