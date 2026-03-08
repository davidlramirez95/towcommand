package repository

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"

	"github.com/davidlramirez95/towcommand/internal/domain/booking"
)

// DynamoDBScanAPI defines the Scan operation for DynamoDB. The concrete
// *dynamodb.Client satisfies this interface. It is separate from DynamoDBAPI
// to avoid modifying existing interfaces.
type DynamoDBScanAPI interface {
	Scan(ctx context.Context, params *dynamodb.ScanInput, optFns ...func(*dynamodb.Options)) (*dynamodb.ScanOutput, error)
}

// FindByProvider lists all bookings assigned to a given provider using a Scan
// with a filter on the providerId attribute. This is acceptable at MVP scale
// where providers have a bounded number of bookings. A dedicated GSI should be
// added for production workloads.
//
// The method performs a runtime type assertion to access the Scan operation on
// the underlying DynamoDB client. If the client does not support Scan (e.g. in
// tests using a minimal mock), an error is returned.
func (r *DynamoBookingRepository) FindByProvider(ctx context.Context, providerID string) ([]booking.Booking, error) {
	scanner, ok := r.client.(DynamoDBScanAPI)
	if !ok {
		return nil, fmt.Errorf("DynamoDB client does not support Scan operation")
	}

	filt := expression.Name("providerId").Equal(expression.Value(providerID)).
		And(expression.Name("entityType").Equal(expression.Value("Booking")))

	expr, err := expression.NewBuilder().WithFilter(filt).Build()
	if err != nil {
		return nil, fmt.Errorf("build filter expression: %w", err)
	}

	var allItems []map[string]types.AttributeValue
	input := &dynamodb.ScanInput{
		TableName:                 &r.tableName,
		FilterExpression:          expr.Filter(),
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
	}

	for {
		resp, scanErr := scanner.Scan(ctx, input)
		if scanErr != nil {
			return nil, fmt.Errorf("dynamodb Scan: %w", scanErr)
		}
		allItems = append(allItems, resp.Items...)
		if resp.LastEvaluatedKey == nil {
			break
		}
		input.ExclusiveStartKey = resp.LastEvaluatedKey
	}

	return unmarshalItems[booking.Booking](allItems)
}
