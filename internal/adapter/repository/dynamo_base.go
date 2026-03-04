// Package repository implements DynamoDB persistence adapters for the
// TowCommand single-table design. It translates between Go domain entities
// and the DynamoDB key patterns established in the legacy TypeScript codebase.
package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

// Key prefixes matching the legacy TypeScript table design (table-design.ts).
const (
	PrefixUser        = "USER#"
	PrefixProvider    = "PROV#"
	PrefixJob         = "JOB#"
	PrefixVehicle     = "VEH#"
	PrefixTransaction = "TXN#"
	PrefixOTP         = "OTP#"
	PrefixStatus      = "STATUS#"
	PrefixTier        = "TIER#"
	PrefixPhone       = "PHONE#"
	PrefixEmail       = "EMAIL#"
	PrefixEvidence    = "EVIDENCE#"
	PrefixMedia       = "MEDIA#"
	PrefixDoc         = "DOC#"
	PrefixPayment     = "PAY#"
	PrefixRate        = "RATE#"

	SKProfile = "PROFILE"
	SKDetails = "DETAILS"
	SKRating  = "RATING"
)

// GSI names as defined in the DynamoDB table design.
const (
	GSI1Name = "GSI1-UserJobs"
	GSI2Name = "GSI2-StatusJobs"
	GSI3Name = "GSI3-ProviderByTier"
	GSI4Name = "GSI4-DisputeByStatus"
	GSI5Name = "GSI5-PhoneIndex"
)

// DynamoDBAPI defines the subset of DynamoDB operations used by repositories.
// The concrete *dynamodb.Client satisfies this interface.
type DynamoDBAPI interface {
	GetItem(ctx context.Context, params *dynamodb.GetItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.GetItemOutput, error)
	PutItem(ctx context.Context, params *dynamodb.PutItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.PutItemOutput, error)
	UpdateItem(ctx context.Context, params *dynamodb.UpdateItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.UpdateItemOutput, error)
	Query(ctx context.Context, params *dynamodb.QueryInput, optFns ...func(*dynamodb.Options)) (*dynamodb.QueryOutput, error)
	TransactWriteItems(ctx context.Context, params *dynamodb.TransactWriteItemsInput, optFns ...func(*dynamodb.Options)) (*dynamodb.TransactWriteItemsOutput, error)
}

// baseRepository provides shared DynamoDB helpers for all repository implementations.
type baseRepository struct {
	client    DynamoDBAPI
	tableName string
}

// compositeKey builds a DynamoDB primary key from partition and sort key strings.
func compositeKey(pk, sk string) map[string]types.AttributeValue {
	return map[string]types.AttributeValue{
		"PK": stringAttr(pk),
		"SK": stringAttr(sk),
	}
}

// stringAttr creates a DynamoDB string attribute value.
func stringAttr(v string) *types.AttributeValueMemberS {
	return &types.AttributeValueMemberS{Value: v}
}

// marshalItem converts a Go struct to a DynamoDB attribute map using json struct tags.
func marshalItem(v any) (map[string]types.AttributeValue, error) {
	return attributevalue.MarshalMapWithOptions(v, func(o *attributevalue.EncoderOptions) {
		o.TagKey = "json"
	})
}

// unmarshalItem converts a DynamoDB attribute map to a Go struct using json struct tags.
func unmarshalItem(item map[string]types.AttributeValue, out any) error {
	return attributevalue.UnmarshalMapWithOptions(item, out, func(o *attributevalue.DecoderOptions) {
		o.TagKey = "json"
	})
}

// unmarshalItems converts a slice of DynamoDB items to a typed Go slice.
func unmarshalItems[T any](items []map[string]types.AttributeValue) ([]T, error) {
	result := make([]T, 0, len(items))
	for _, item := range items {
		var v T
		if err := unmarshalItem(item, &v); err != nil {
			return nil, err
		}
		result = append(result, v)
	}
	return result, nil
}

// getItem retrieves a single item by composite key. Returns false if the item does not exist.
func (r *baseRepository) getItem(ctx context.Context, pk, sk string, out any) (bool, error) {
	resp, err := r.client.GetItem(ctx, &dynamodb.GetItemInput{
		TableName: &r.tableName,
		Key:       compositeKey(pk, sk),
	})
	if err != nil {
		return false, fmt.Errorf("dynamodb GetItem: %w", err)
	}
	if resp.Item == nil {
		return false, nil
	}
	if err := unmarshalItem(resp.Item, out); err != nil {
		return false, fmt.Errorf("unmarshal item: %w", err)
	}
	return true, nil
}

// putItem stores an item in DynamoDB.
func (r *baseRepository) putItem(ctx context.Context, item map[string]types.AttributeValue) error {
	_, err := r.client.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: &r.tableName,
		Item:      item,
	})
	if err != nil {
		return fmt.Errorf("dynamodb PutItem: %w", err)
	}
	return nil
}

// updateItem applies a set of attribute updates to an existing item.
func (r *baseRepository) updateItem(ctx context.Context, pk, sk string, updates map[string]any) error {
	var upd expression.UpdateBuilder
	first := true
	for k, v := range updates {
		if first {
			upd = expression.Set(expression.Name(k), expression.Value(v))
			first = false
		} else {
			upd = upd.Set(expression.Name(k), expression.Value(v))
		}
	}

	expr, err := expression.NewBuilder().WithUpdate(upd).Build()
	if err != nil {
		return fmt.Errorf("build update expression: %w", err)
	}

	_, err = r.client.UpdateItem(ctx, &dynamodb.UpdateItemInput{
		TableName:                 &r.tableName,
		Key:                       compositeKey(pk, sk),
		UpdateExpression:          expr.Update(),
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
	})
	if err != nil {
		return fmt.Errorf("dynamodb UpdateItem: %w", err)
	}
	return nil
}

// queryItems executes a DynamoDB query and returns raw items.
func (r *baseRepository) queryItems(ctx context.Context, input *dynamodb.QueryInput) ([]map[string]types.AttributeValue, error) {
	input.TableName = &r.tableName
	resp, err := r.client.Query(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("dynamodb Query: %w", err)
	}
	return resp.Items, nil
}

// formatTime formats a time value for use in DynamoDB sort keys.
func formatTime(t time.Time) string {
	return t.Format(time.RFC3339)
}
