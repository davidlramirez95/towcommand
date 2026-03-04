//go:build integration

package repository

import (
	"context"
	"os"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

const testTableName = "TowCommand-test"

// newTestClient creates a DynamoDB client pointing at LocalStack and ensures
// the test table exists. It registers a cleanup function that deletes the table.
func newTestClient(t *testing.T) *dynamodb.Client {
	t.Helper()

	endpoint := os.Getenv("DYNAMODB_ENDPOINT")
	if endpoint == "" {
		endpoint = "http://localhost:4566"
	}

	cfg, err := config.LoadDefaultConfig(context.Background(),
		config.WithRegion("ap-southeast-1"),
	)
	if err != nil {
		t.Fatalf("load aws config: %v", err)
	}

	client := dynamodb.NewFromConfig(cfg, func(o *dynamodb.Options) {
		o.BaseEndpoint = aws.String(endpoint)
	})

	createTestTable(t, client)
	return client
}

func createTestTable(t *testing.T, client *dynamodb.Client) {
	t.Helper()

	// Delete table if it exists from a previous run.
	_, _ = client.DeleteTable(context.Background(), &dynamodb.DeleteTableInput{
		TableName: aws.String(testTableName),
	})
	waiter := dynamodb.NewTableNotExistsWaiter(client)
	_ = waiter.Wait(context.Background(), &dynamodb.DescribeTableInput{
		TableName: aws.String(testTableName),
	}, nil)

	_, err := client.CreateTable(context.Background(), &dynamodb.CreateTableInput{
		TableName: aws.String(testTableName),
		KeySchema: []types.KeySchemaElement{
			{AttributeName: aws.String("PK"), KeyType: types.KeyTypeHash},
			{AttributeName: aws.String("SK"), KeyType: types.KeyTypeRange},
		},
		AttributeDefinitions: []types.AttributeDefinition{
			{AttributeName: aws.String("PK"), AttributeType: types.ScalarAttributeTypeS},
			{AttributeName: aws.String("SK"), AttributeType: types.ScalarAttributeTypeS},
			{AttributeName: aws.String("GSI1PK"), AttributeType: types.ScalarAttributeTypeS},
			{AttributeName: aws.String("GSI1SK"), AttributeType: types.ScalarAttributeTypeS},
			{AttributeName: aws.String("GSI2PK"), AttributeType: types.ScalarAttributeTypeS},
			{AttributeName: aws.String("GSI2SK"), AttributeType: types.ScalarAttributeTypeS},
			{AttributeName: aws.String("GSI3PK"), AttributeType: types.ScalarAttributeTypeS},
			{AttributeName: aws.String("GSI3SK"), AttributeType: types.ScalarAttributeTypeS},
			{AttributeName: aws.String("GSI5PK"), AttributeType: types.ScalarAttributeTypeS},
			{AttributeName: aws.String("GSI5SK"), AttributeType: types.ScalarAttributeTypeS},
		},
		GlobalSecondaryIndexes: []types.GlobalSecondaryIndex{
			{
				IndexName: aws.String(GSI1Name),
				KeySchema: []types.KeySchemaElement{
					{AttributeName: aws.String("GSI1PK"), KeyType: types.KeyTypeHash},
					{AttributeName: aws.String("GSI1SK"), KeyType: types.KeyTypeRange},
				},
				Projection: &types.Projection{ProjectionType: types.ProjectionTypeAll},
			},
			{
				IndexName: aws.String(GSI2Name),
				KeySchema: []types.KeySchemaElement{
					{AttributeName: aws.String("GSI2PK"), KeyType: types.KeyTypeHash},
					{AttributeName: aws.String("GSI2SK"), KeyType: types.KeyTypeRange},
				},
				Projection: &types.Projection{ProjectionType: types.ProjectionTypeAll},
			},
			{
				IndexName: aws.String(GSI3Name),
				KeySchema: []types.KeySchemaElement{
					{AttributeName: aws.String("GSI3PK"), KeyType: types.KeyTypeHash},
					{AttributeName: aws.String("GSI3SK"), KeyType: types.KeyTypeRange},
				},
				Projection: &types.Projection{ProjectionType: types.ProjectionTypeAll},
			},
			{
				IndexName: aws.String(GSI5Name),
				KeySchema: []types.KeySchemaElement{
					{AttributeName: aws.String("GSI5PK"), KeyType: types.KeyTypeHash},
					{AttributeName: aws.String("GSI5SK"), KeyType: types.KeyTypeRange},
				},
				Projection: &types.Projection{ProjectionType: types.ProjectionTypeAll},
			},
		},
		BillingMode: types.BillingModePayPerRequest,
	})
	if err != nil {
		t.Fatalf("create test table: %v", err)
	}

	waiterActive := dynamodb.NewTableExistsWaiter(client)
	if err := waiterActive.Wait(context.Background(), &dynamodb.DescribeTableInput{
		TableName: aws.String(testTableName),
	}, nil); err != nil {
		t.Fatalf("wait for table active: %v", err)
	}

	t.Cleanup(func() {
		_, _ = client.DeleteTable(context.Background(), &dynamodb.DeleteTableInput{
			TableName: aws.String(testTableName),
		})
	})
}
