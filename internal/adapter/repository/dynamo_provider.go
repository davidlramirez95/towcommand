package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"

	"github.com/davidlramirez95/towcommand/internal/domain/provider"
	"github.com/davidlramirez95/towcommand/internal/domain/user"
)

// DynamoProviderRepository implements provider persistence against DynamoDB.
type DynamoProviderRepository struct {
	baseRepository
}

// NewProviderRepository creates a new DynamoDB-backed provider repository.
func NewProviderRepository(client DynamoDBAPI, tableName string) *DynamoProviderRepository {
	return &DynamoProviderRepository{
		baseRepository: baseRepository{client: client, tableName: tableName},
	}
}

// providerCity extracts the primary city from service areas, defaulting to NCR.
func providerCity(p *provider.Provider) string {
	if len(p.ServiceAreas) > 0 {
		return p.ServiceAreas[0]
	}
	return "NCR"
}

// gsi3SK formats the rating for GSI3 sort key (zero-padded to 5 chars).
func gsi3SK(rating float64) string {
	return fmt.Sprintf("%05.0f", rating)
}

// Save persists a provider with all key attributes.
// PK: PROV#<providerId>, SK: PROFILE
// GSI3: TIER#<trustTier>#<city> / <rating padded>
func (r *DynamoProviderRepository) Save(ctx context.Context, p *provider.Provider) error {
	item, err := marshalItem(p)
	if err != nil {
		return fmt.Errorf("marshal provider: %w", err)
	}

	city := providerCity(p)

	item["PK"] = stringAttr(PrefixProvider + p.ProviderID)
	item["SK"] = stringAttr(SKProfile)
	item["GSI3PK"] = stringAttr(PrefixTier + string(p.TrustTier) + "#" + city)
	item["GSI3SK"] = stringAttr(gsi3SK(p.Rating))
	item["entityType"] = stringAttr("Provider")

	return r.putItem(ctx, item)
}

// FindByID retrieves a provider by their ID. Returns nil if not found.
func (r *DynamoProviderRepository) FindByID(ctx context.Context, providerID string) (*provider.Provider, error) {
	var p provider.Provider
	found, err := r.getItem(ctx, PrefixProvider+providerID, SKProfile, &p)
	if err != nil {
		return nil, err
	}
	if !found {
		return nil, nil
	}
	return &p, nil
}

// FindByTierAndCity lists providers by trust tier and city via GSI3, ordered by rating descending.
func (r *DynamoProviderRepository) FindByTierAndCity(ctx context.Context, tier user.TrustTier, city string, limit int32) ([]provider.Provider, error) {
	keyCond := expression.Key("GSI3PK").Equal(expression.Value(PrefixTier + string(tier) + "#" + city))

	expr, err := expression.NewBuilder().WithKeyCondition(keyCond).Build()
	if err != nil {
		return nil, fmt.Errorf("build key expression: %w", err)
	}

	items, err := r.queryItems(ctx, &dynamodb.QueryInput{
		IndexName:                 aws.String(GSI3Name),
		KeyConditionExpression:    expr.KeyCondition(),
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		ScanIndexForward:          aws.Bool(false),
		Limit:                     &limit,
	})
	if err != nil {
		return nil, err
	}
	return unmarshalItems[provider.Provider](items)
}

// UpdateLocation updates a provider's GPS coordinates and timestamp.
func (r *DynamoProviderRepository) UpdateLocation(ctx context.Context, providerID string, lat, lng float64) error {
	now := time.Now().UTC()
	return r.updateItem(ctx, PrefixProvider+providerID, SKProfile, map[string]any{
		"currentLat":         lat,
		"currentLng":         lng,
		"lastLocationUpdate": now,
		"updatedAt":          now,
	})
}

// UpdateAvailability toggles a provider's online status.
func (r *DynamoProviderRepository) UpdateAvailability(ctx context.Context, providerID string, isOnline bool) error {
	return r.updateItem(ctx, PrefixProvider+providerID, SKProfile, map[string]any{
		"isOnline":  isOnline,
		"updatedAt": time.Now().UTC(),
	})
}

// UploadDoc persists a provider KYC document.
// PK: PROV#<providerId>, SK: DOC#<docType>
func (r *DynamoProviderRepository) UploadDoc(ctx context.Context, doc *provider.ProviderDoc) error {
	item, err := marshalItem(doc)
	if err != nil {
		return fmt.Errorf("marshal provider doc: %w", err)
	}

	item["PK"] = stringAttr(PrefixProvider + doc.ProviderID)
	item["SK"] = stringAttr(PrefixDoc + string(doc.DocType))
	item["entityType"] = stringAttr("ProviderDoc")

	return r.putItem(ctx, item)
}

// GetDocs lists all KYC documents for a provider.
func (r *DynamoProviderRepository) GetDocs(ctx context.Context, providerID string) ([]provider.ProviderDoc, error) {
	keyCond := expression.Key("PK").Equal(expression.Value(PrefixProvider + providerID)).
		And(expression.Key("SK").BeginsWith(PrefixDoc))

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
	return unmarshalItems[provider.ProviderDoc](items)
}
