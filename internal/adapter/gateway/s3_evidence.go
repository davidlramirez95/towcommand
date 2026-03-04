package gateway

import (
	"context"
	"fmt"
	"time"

	v4 "github.com/aws/aws-sdk-go-v2/aws/signer/v4"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

// S3PresignAPI is the subset of the S3 presign client needed by the evidence adapter.
type S3PresignAPI interface {
	PresignPutObject(ctx context.Context, params *s3.PutObjectInput, optFns ...func(*s3.PresignOptions)) (*v4.PresignedHTTPRequest, error)
}

// S3EvidenceAdapter generates presigned upload URLs for evidence media.
// It implements the port.PresignedURLGenerator interface.
type S3EvidenceAdapter struct {
	presigner S3PresignAPI
	bucket    string
}

// NewS3EvidenceAdapter creates a new S3 evidence adapter.
func NewS3EvidenceAdapter(presigner S3PresignAPI, bucket string) *S3EvidenceAdapter {
	return &S3EvidenceAdapter{
		presigner: presigner,
		bucket:    bucket,
	}
}

// GenerateUploadURL returns a presigned PUT URL for uploading evidence media to S3.
func (a *S3EvidenceAdapter) GenerateUploadURL(ctx context.Context, key, contentType string, expiry time.Duration) (string, error) {
	result, err := a.presigner.PresignPutObject(ctx, &s3.PutObjectInput{
		Bucket:      &a.bucket,
		Key:         &key,
		ContentType: &contentType,
	}, func(opts *s3.PresignOptions) {
		opts.Expires = expiry
	})
	if err != nil {
		return "", fmt.Errorf("presigning put object for key %s: %w", key, err)
	}
	return result.URL, nil
}
