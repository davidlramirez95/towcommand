// Package main is the composition root for the trigger-matching EventBridge subscriber Lambda.
// It listens for BookingCreated events and runs the cascade matching algorithm to find
// available providers, then sends offers via WebSocket.
package main

import (
	"context"
	"errors"
	"log/slog"
	"runtime/debug"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"

	"github.com/davidlramirez95/towcommand/internal/adapter/cache"
	"github.com/davidlramirez95/towcommand/internal/adapter/gateway"
	"github.com/davidlramirez95/towcommand/internal/adapter/handler"
	"github.com/davidlramirez95/towcommand/internal/adapter/repository"
	"github.com/davidlramirez95/towcommand/internal/domain/booking"
	domainerrors "github.com/davidlramirez95/towcommand/internal/domain/errors"
	"github.com/davidlramirez95/towcommand/internal/domain/event"
	"github.com/davidlramirez95/towcommand/internal/platform/awsclient"
	"github.com/davidlramirez95/towcommand/internal/platform/config"
	"github.com/davidlramirez95/towcommand/internal/platform/logger"
	"github.com/davidlramirez95/towcommand/internal/usecase/matching"
	"github.com/davidlramirez95/towcommand/internal/usecase/port"
)

// deps bundles all handler dependencies wired at cold-start.
type deps struct {
	matcher     *matching.MatchBookingUseCase
	bookings    port.BookingStatusUpdater
	sessions    port.SessionCache
	wsPoster    handler.ConnectionPoster
	events      port.EventPublisher
	rateLimiter port.RateLimiter
}

// BookingCreatedDetail is the expected detail shape for BookingCreated events.
type BookingCreatedDetail struct {
	BookingID   string  `json:"bookingId" validate:"required"`
	CustomerID  string  `json:"customerId" validate:"required"`
	PickupLat   float64 `json:"pickupLat" validate:"required"`
	PickupLng   float64 `json:"pickupLng" validate:"required"`
	ServiceType string  `json:"serviceType" validate:"required"`
	WeightClass string  `json:"weightClass"`
}

// matchOfferMessage is the WebSocket payload sent to the provider.
type matchOfferMessage struct {
	Action    string              `json:"action"`
	BookingID string              `json:"bookingId"`
	Scores    []matchScoreSummary `json:"scores"`
	SurgeMode bool                `json:"surgeMode"`
}

type matchScoreSummary struct {
	ProviderID string  `json:"providerId"`
	TotalScore float64 `json:"totalScore"`
	DistanceKm float64 `json:"distanceKm"`
}

func main() {
	cfg := config.Load()
	log := logger.New(cfg.Stage, cfg.FunctionName, cfg.FunctionVersion, cfg.LogLevel)
	slog.SetDefault(log)

	ddb := awsclient.DynamoDBClient(cfg)
	eb := awsclient.EventBridgeClient(cfg)

	bookingRepo := repository.NewBookingRepository(ddb, cfg.DynamoDBTable)
	providerRepo := repository.NewProviderRepository(ddb, cfg.DynamoDBTable)
	evPublisher := gateway.NewEventBridgePublisher(eb, cfg.EventBusName, log)

	redisClient := cache.NewRedisClient(cache.Options{
		Host:     cfg.RedisHost,
		Port:     cfg.RedisPort,
		Password: cfg.RedisPassword,
	})
	geoCache := cache.NewRedisGeoCache(redisClient)
	surgeCache := cache.NewRedisSurgeCache(redisClient)
	sessionCacheInst := cache.NewRedisSessionCache(redisClient)
	rateLimiterInst := cache.NewRedisRateLimiter(redisClient)

	matcher := matching.NewMatchBookingUseCase(bookingRepo, providerRepo, geoCache, surgeCache)

	wsClient := awsclient.APIGatewayManagementClient(cfg, "")

	d := &deps{
		matcher:     matcher,
		bookings:    bookingRepo,
		sessions:    sessionCacheInst,
		wsPoster:    wsClient,
		events:      evPublisher,
		rateLimiter: rateLimiterInst,
	}

	lambda.Start(d.handleEvent)
}

//nolint:gocritic // Lambda SDK requires value receiver for CloudWatchEvent
func (d *deps) handleEvent(ctx context.Context, evt events.CloudWatchEvent) error {
	defer func() {
		if r := recover(); r != nil {
			slog.ErrorContext(ctx, "panic recovered in trigger-matching",
				"panic", r,
				"stack", string(debug.Stack()),
			)
		}
	}()

	correlationID := handler.ExtractCorrelationID(&evt)
	ctx = logger.SetCorrelationID(ctx, correlationID)

	slog.InfoContext(ctx, "trigger-matching invoked",
		"detail_type", evt.DetailType,
		"event_id", evt.ID,
	)

	detail, err := handler.ParseEventDetail[BookingCreatedDetail](&evt)
	if err != nil {
		slog.ErrorContext(ctx, "failed to parse BookingCreated detail", "error", err)
		return nil // do not retry invalid events
	}

	ctx = logger.SetBookingID(ctx, detail.BookingID)

	// Rate-limit lock to prevent duplicate matching.
	key := "match:" + detail.BookingID
	allowed, _, err := d.rateLimiter.CheckRateLimit(ctx, key, 1, 60)
	if err != nil {
		slog.WarnContext(ctx, "rate limiter error, proceeding anyway", "error", err)
	} else if !allowed {
		slog.WarnContext(ctx, "duplicate matching attempt blocked", "booking_id", detail.BookingID)
		return nil
	}

	result, err := d.matcher.Execute(ctx, detail.BookingID)
	if err != nil {
		var appErr *domainerrors.AppError
		if errors.As(err, &appErr) && appErr.Code == domainerrors.CodeProviderUnavailable {
			slog.WarnContext(ctx, "no providers available for booking", "booking_id", detail.BookingID)
			_ = d.events.Publish(ctx, event.SourceMatching, event.MatchingFailed, map[string]any{
				"bookingId":  detail.BookingID,
				"customerId": detail.CustomerID,
				"reason":     "no_providers_available",
			}, nil)
			return nil
		}
		slog.ErrorContext(ctx, "matching failed", "error", err)
		return nil // never fail the EventBridge trigger
	}

	// Update booking status to MATCHED.
	if err := d.bookings.UpdateStatus(ctx, detail.BookingID, booking.BookingStatusMatched, map[string]any{
		"matchedProviders": len(result.Scores),
		"surgeMode":        result.SurgeMode,
	}); err != nil {
		slog.ErrorContext(ctx, "failed to update booking status to MATCHED", "error", err)
	}

	// Send offer to top provider via WebSocket.
	if len(result.Scores) > 0 {
		topProvider := result.Scores[0]
		connID, connErr := d.sessions.GetConnection(ctx, topProvider.ProviderID)
		if connErr != nil {
			slog.WarnContext(ctx, "could not get provider connection", "provider_id", topProvider.ProviderID, "error", connErr)
		} else if connID != "" {
			scores := make([]matchScoreSummary, 0, len(result.Scores))
			for _, s := range result.Scores {
				scores = append(scores, matchScoreSummary{
					ProviderID: s.ProviderID,
					TotalScore: s.TotalScore,
					DistanceKm: s.DistanceKm,
				})
			}
			offer := matchOfferMessage{
				Action:    "MATCH_OFFER",
				BookingID: detail.BookingID,
				Scores:    scores,
				SurgeMode: result.SurgeMode,
			}
			if sendErr := handler.SendToConnection(ctx, d.wsPoster, connID, offer); sendErr != nil {
				slog.WarnContext(ctx, "failed to send match offer via WebSocket", "error", sendErr, "connection_id", connID)
			}
		}
	}

	// Publish MatchingCompleted event.
	_ = d.events.Publish(ctx, event.SourceMatching, event.MatchingCompleted, map[string]any{
		"bookingId":  detail.BookingID,
		"customerId": detail.CustomerID,
		"scores":     result.Scores,
		"surgeMode":  result.SurgeMode,
	}, nil)

	slog.InfoContext(ctx, "matching completed",
		"booking_id", detail.BookingID,
		"matched_providers", len(result.Scores),
		"surge_mode", result.SurgeMode,
	)

	return nil
}
