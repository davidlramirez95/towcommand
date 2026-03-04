package repository

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"

	"github.com/davidlramirez95/towcommand/internal/domain/provider"
)

// scannerAPI extends DynamoDBAPI with Scan support.
// The AWS SDK *dynamodb.Client satisfies this interface.
type scannerAPI interface {
	Scan(ctx context.Context, params *dynamodb.ScanInput, optFns ...func(*dynamodb.Options)) (*dynamodb.ScanOutput, error)
}

// FindByCognitoSub retrieves a provider by their Cognito sub using a filtered scan.
// This is used during token generation to resolve the provider ID from a user identity.
// Consider adding a GSI for cognitoSub if this becomes a hot path.
func (r *DynamoProviderRepository) FindByCognitoSub(ctx context.Context, cognitoSub string) (*provider.Provider, error) {
	scanner, ok := r.client.(scannerAPI)
	if !ok {
		return nil, fmt.Errorf("DynamoDB client does not support Scan operations")
	}

	filter := expression.Name("cognitoSub").Equal(expression.Value(cognitoSub)).
		And(expression.Name("entityType").Equal(expression.Value("Provider")))

	expr, err := expression.NewBuilder().WithFilter(filter).Build()
	if err != nil {
		return nil, fmt.Errorf("build filter expression: %w", err)
	}

	resp, err := scanner.Scan(ctx, &dynamodb.ScanInput{
		TableName:                 &r.tableName,
		FilterExpression:          expr.Filter(),
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
	})
	if err != nil {
		return nil, fmt.Errorf("dynamodb Scan: %w", err)
	}
	if len(resp.Items) == 0 {
		return nil, nil
	}

	var p provider.Provider
	if err := unmarshalItem(resp.Items[0], &p); err != nil {
		return nil, fmt.Errorf("unmarshal provider: %w", err)
	}
	return &p, nil
}
