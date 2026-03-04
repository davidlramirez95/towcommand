package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"

	"github.com/davidlramirez95/towcommand/internal/domain/payment"
)

// DynamoPaymentRepository implements payment persistence against DynamoDB.
type DynamoPaymentRepository struct {
	baseRepository
}

// NewPaymentRepository creates a new DynamoDB-backed payment repository.
func NewPaymentRepository(client DynamoDBAPI, tableName string) *DynamoPaymentRepository {
	return &DynamoPaymentRepository{
		baseRepository: baseRepository{client: client, tableName: tableName},
	}
}

// Save persists a payment with all key attributes.
// PK: TXN#<paymentId>, SK: DETAILS
// GSI1: JOB#<bookingId> / PAY#<createdAt>
func (r *DynamoPaymentRepository) Save(ctx context.Context, p *payment.Payment) error {
	item, err := marshalItem(p)
	if err != nil {
		return fmt.Errorf("marshal payment: %w", err)
	}

	item["PK"] = stringAttr(PrefixTransaction + p.PaymentID)
	item["SK"] = stringAttr(SKDetails)
	item["GSI1PK"] = stringAttr(PrefixJob + p.BookingID)
	item["GSI1SK"] = stringAttr(PrefixPayment + formatTime(p.CreatedAt))
	item["entityType"] = stringAttr("Payment")

	return r.putItem(ctx, item)
}

// FindByID retrieves a payment by its ID. Returns nil if not found.
func (r *DynamoPaymentRepository) FindByID(ctx context.Context, paymentID string) (*payment.Payment, error) {
	var p payment.Payment
	found, err := r.getItem(ctx, PrefixTransaction+paymentID, SKDetails, &p)
	if err != nil {
		return nil, err
	}
	if !found {
		return nil, nil
	}
	return &p, nil
}

// FindByBooking lists payments for a booking via GSI1, ordered by creation date descending.
func (r *DynamoPaymentRepository) FindByBooking(ctx context.Context, bookingID string) ([]payment.Payment, error) {
	keyCond := expression.Key("GSI1PK").Equal(expression.Value(PrefixJob + bookingID)).
		And(expression.Key("GSI1SK").BeginsWith(PrefixPayment))

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
	})
	if err != nil {
		return nil, err
	}
	return unmarshalItems[payment.Payment](items)
}

// UpdateStatus changes a payment's status.
func (r *DynamoPaymentRepository) UpdateStatus(ctx context.Context, paymentID string, status payment.PaymentStatus) error {
	return r.updateItem(ctx, PrefixTransaction+paymentID, SKDetails, map[string]any{
		"status":    string(status),
		"updatedAt": time.Now().UTC(),
	})
}
