package gateway

import (
	"context"
	"testing"
	"time"

	v4 "github.com/aws/aws-sdk-go-v2/aws/signer/v4"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// --- Mocks ---

type mockS3Presigner struct{ mock.Mock }

func (m *mockS3Presigner) PresignPutObject(ctx context.Context, params *s3.PutObjectInput, optFns ...func(*s3.PresignOptions)) (*v4.PresignedHTTPRequest, error) {
	args := m.Called(ctx, params, optFns)
	if v := args.Get(0); v != nil {
		return v.(*v4.PresignedHTTPRequest), args.Error(1)
	}
	return nil, args.Error(1)
}

// --- Tests ---

func TestS3EvidenceAdapter_GenerateUploadURL_Success(t *testing.T) {
	presigner := new(mockS3Presigner)
	adapter := NewS3EvidenceAdapter(presigner, "test-bucket")

	expectedURL := "https://test-bucket.s3.amazonaws.com/evidence/booking-123/pickup/FRONT_1234.jpg?X-Amz-Signature=abc"

	presigner.On("PresignPutObject", mock.Anything, mock.MatchedBy(func(input *s3.PutObjectInput) bool {
		return *input.Bucket == "test-bucket" &&
			*input.Key == "evidence/booking-123/pickup/FRONT_1234.jpg" &&
			*input.ContentType == "image/jpeg"
	}), mock.Anything).Return(&v4.PresignedHTTPRequest{
		URL: expectedURL,
	}, nil)

	url, err := adapter.GenerateUploadURL(
		context.Background(),
		"evidence/booking-123/pickup/FRONT_1234.jpg",
		"image/jpeg",
		15*time.Minute,
	)

	require.NoError(t, err)
	assert.Equal(t, expectedURL, url)
	presigner.AssertExpectations(t)
}

func TestS3EvidenceAdapter_GenerateUploadURL_Error(t *testing.T) {
	presigner := new(mockS3Presigner)
	adapter := NewS3EvidenceAdapter(presigner, "test-bucket")

	presigner.On("PresignPutObject", mock.Anything, mock.Anything, mock.Anything).
		Return(nil, assert.AnError)

	url, err := adapter.GenerateUploadURL(
		context.Background(),
		"evidence/booking-123/pickup/FRONT_1234.jpg",
		"image/jpeg",
		15*time.Minute,
	)

	assert.Error(t, err)
	assert.Empty(t, url)
	assert.Contains(t, err.Error(), "presigning put object")
	presigner.AssertExpectations(t)
}
