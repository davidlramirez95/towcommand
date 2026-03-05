package handler

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"

	domainerrors "github.com/davidlramirez95/towcommand/internal/domain/errors"
)

// DynamoDBReader defines the subset of DynamoDB operations needed by BookingStatsHandler.
type DynamoDBReader interface {
	GetItem(ctx context.Context, params *dynamodb.GetItemInput, optFns ...func(*dynamodb.Options)) (*dynamodb.GetItemOutput, error)
}

// BookingStatsHandler handles GET /admin/stats/bookings requests.
type BookingStatsHandler struct {
	ddb       DynamoDBReader
	tableName string
}

// NewBookingStatsHandler constructs a BookingStatsHandler.
func NewBookingStatsHandler(ddb DynamoDBReader, tableName string) *BookingStatsHandler {
	return &BookingStatsHandler{ddb: ddb, tableName: tableName}
}

// Handle processes an admin booking-stats API Gateway event.
func (h *BookingStatsHandler) Handle(ctx context.Context, event *events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	userID := ExtractUserID(event)
	if userID == "" {
		return ErrorResponse(domainerrors.NewUnauthorizedError()), nil
	}

	date := ParseQueryParam(event, "date")
	if date == "" {
		date = time.Now().UTC().Format("2006-01-02")
	}

	pk := "ANALYTICS#DAILY#" + date
	sk := "SUMMARY"

	resp, err := h.ddb.GetItem(ctx, &dynamodb.GetItemInput{
		TableName: &h.tableName,
		Key: map[string]types.AttributeValue{
			"PK": &types.AttributeValueMemberS{Value: pk},
			"SK": &types.AttributeValueMemberS{Value: sk},
		},
	})
	if err != nil {
		return ErrorResponse(domainerrors.NewInternalError(
			fmt.Sprintf("failed to read analytics: %s", err.Error()),
		)), nil
	}

	stats := make(map[string]int64)
	if resp.Item != nil {
		for k, v := range resp.Item {
			// Skip key attributes.
			if k == "PK" || k == "SK" || k == "entityType" {
				continue
			}
			if nv, ok := v.(*types.AttributeValueMemberN); ok {
				n, parseErr := strconv.ParseInt(nv.Value, 10, 64)
				if parseErr == nil {
					stats[k] = n
				}
			}
		}
	}

	return SuccessResponse(http.StatusOK, map[string]any{
		"date":  date,
		"stats": stats,
	}), nil
}
