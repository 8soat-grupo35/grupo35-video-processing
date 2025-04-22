package wrappers

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

//go:generate mockgen -source=s3_client.go -destination=mock/s3_client.go
type IS3Client interface {
	Upload(ctx context.Context, input *s3.PutObjectInput, opts ...func(*manager.Uploader)) (*manager.UploadOutput, error)
}

type S3Client struct {
	client *s3.Client
}

func NewS3Client(cfg aws.Config) IS3Client {
	return S3Client{
		client: s3.NewFromConfig(cfg),
	}
}

// Upload implements IS3Client.
func (s S3Client) Upload(ctx context.Context, input *s3.PutObjectInput, opts ...func(*manager.Uploader)) (*manager.UploadOutput, error) {
	uploader := manager.NewUploader(s.client)
	return uploader.Upload(ctx, input, opts...)
}
