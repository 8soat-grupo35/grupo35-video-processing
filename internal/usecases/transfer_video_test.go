package usecases

import (
	"errors"
	"testing"

	mock_repository "github.com/8soat-grupo35/grupo35-video-processing/internal/interfaces/repository/mock"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestTransferVideoUseCase_UploadVideo_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockS3 := mock_repository.NewMockS3(ctrl)
	key := "test-key"
	videoData := []byte("video-bytes")
	contentType := "video/mp4"
	bucketName := "grupo35-video-uploaded"

	mockS3.EXPECT().SetBucketName(bucketName).Times(1)
	mockS3.EXPECT().UploadFile(key, videoData, contentType).Return(nil).Times(1)

	useCase := NewTransferVideo(mockS3)
	err := useCase.UploadVideo(key, videoData, contentType)
	assert.NoError(t, err)
}

func TestTransferVideoUseCase_UploadVideo_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockS3 := mock_repository.NewMockS3(ctrl)
	key := "test-key"
	videoData := []byte("video-bytes")
	contentType := "video/mp4"
	bucketName := "grupo35-video-uploaded"
	expectedErr := errors.New("upload error")

	mockS3.EXPECT().SetBucketName(bucketName).Times(1)
	mockS3.EXPECT().UploadFile(key, videoData, contentType).Return(expectedErr).Times(1)

	useCase := NewTransferVideo(mockS3)
	err := useCase.UploadVideo(key, videoData, contentType)
	assert.EqualError(t, err, expectedErr.Error())
}
