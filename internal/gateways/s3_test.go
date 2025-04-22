package gateways

import (
	"bytes"
	"context"
	"errors"
	"testing"

	wrappers_mock "github.com/8soat-grupo35/grupo35-video-processing/internal/adapters/wrappers/mock"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestSetBucketName(t *testing.T) {
	s3Gateway := NewS3Gateway(nil)
	bucketName := "bucket-test"
	s3Gateway.SetBucketName(bucketName)
	assert.Equal(t, bucketName, *s3Gateway.(*S3Gateway).bucketName)
}

func TestUploadFile_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockS3Client := wrappers_mock.NewMockIS3Client(ctrl)
	bucketName := "bucket-test"
	key := "test/key"
	fileData := []byte("file content")
	contentType := "video/mp4"

	mockS3Client.EXPECT().
		Upload(context.TODO(), &s3.PutObjectInput{
			Bucket:      aws.String(bucketName),
			Key:         aws.String(key),
			Body:        bytes.NewReader(fileData),
			ContentType: aws.String(contentType),
		}).
		Return(nil, nil).
		Times(1)

	s3Gateway := NewS3Gateway(mockS3Client)
	s3Gateway.SetBucketName(bucketName)
	err := s3Gateway.UploadFile(key, fileData, contentType)
	assert.NoError(t, err)
}

func TestUploadFile_Error_BucketNotSet(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockS3Client := wrappers_mock.NewMockIS3Client(ctrl)
	key := "test/key"
	fileData := []byte("file content")
	contentType := "video/mp4"

	s3Gateway := NewS3Gateway(mockS3Client)
	err := s3Gateway.UploadFile(key, fileData, contentType)
	assert.EqualError(t, err, "bucket name is not set")
}

func TestUploadFile_Error_Upload(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockS3Client := wrappers_mock.NewMockIS3Client(ctrl)
	bucketName := "bucket-test"
	key := "test/key"
	fileData := []byte("file content")
	contentType := "video/mp4"
	expectedErr := errors.New("upload error")

	mockS3Client.EXPECT().
		Upload(context.TODO(), &s3.PutObjectInput{
			Bucket:      aws.String(bucketName),
			Key:         aws.String(key),
			Body:        bytes.NewReader(fileData),
			ContentType: aws.String(contentType),
		}).
		Return(nil, expectedErr).
		Times(1)

	s3Gateway := NewS3Gateway(mockS3Client)
	s3Gateway.SetBucketName(bucketName)
	err := s3Gateway.UploadFile(key, fileData, contentType)
	assert.EqualError(t, err, expectedErr.Error())
}
