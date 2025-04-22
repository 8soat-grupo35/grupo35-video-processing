package gateways

import (
	"bytes"
	"context"
	"errors"
	"fmt"

	"github.com/8soat-grupo35/grupo35-video-processing/internal/adapters/wrappers"
	"github.com/8soat-grupo35/grupo35-video-processing/internal/interfaces/repository"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type S3Gateway struct {
	client     wrappers.IS3Client
	bucketName *string
}

func NewS3Gateway(client wrappers.IS3Client) repository.S3 {
	return &S3Gateway{
		client: client,
	}
}

func (s *S3Gateway) SetBucketName(bucketName string) {
	s.bucketName = &bucketName
}

func (s *S3Gateway) UploadFile(key string, fileData []byte) error {
	if s.bucketName == nil {
		return errors.New("bucket name is not set")
	}

	fmt.Printf("Uploading file to S3 with key: %s\n", key)

	_, err := s.client.Upload(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String(*s.bucketName),
		Key:    aws.String(key),
		Body:   bytes.NewReader(fileData),
	})

	if err != nil {
		return fmt.Errorf("failed to upload file: %w", err)
	}

	return nil
}
