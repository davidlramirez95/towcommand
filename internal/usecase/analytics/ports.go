// Package analytics implements event-driven analytics recording for the TowCommand platform.
// It processes domain events and updates DynamoDB atomic counters and heatmap cells.
package analytics

import "context"

// AnalyticsRecorder persists analytics counter updates to the data store.
type AnalyticsRecorder interface {
	IncrementDailyCounter(ctx context.Context, date, field string, delta int64) error
	IncrementHeatmapCell(ctx context.Context, date, geohash string, lat, lng float64) error
	IncrementProviderCounter(ctx context.Context, providerID, date, field string, delta int64) error
}
