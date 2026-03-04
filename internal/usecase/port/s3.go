package port

import (
	"context"
	"time"
)

// PresignedURLGenerator generates presigned URLs for S3 object uploads.
type PresignedURLGenerator interface {
	GenerateUploadURL(ctx context.Context, key, contentType string, expiry time.Duration) (string, error)
}
