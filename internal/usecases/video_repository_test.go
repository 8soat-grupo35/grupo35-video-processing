package usecases

import (
	"errors"
	"testing"

	"github.com/8soat-grupo35/grupo35-video-processing/internal/entities"
	repository_mocks "github.com/8soat-grupo35/grupo35-video-processing/internal/interfaces/repository/mock"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestVideoRepository_GetByUserId(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDynamo := repository_mocks.NewMockDynamo(ctrl)
	videoRepo := NewVideoRepository(mockDynamo)

	userId := "test-user-id"
	expectedVideos := []entities.Video{
		{
			UserID: userId,
			Path:   "video1.mp4",
		},
		{
			UserID: userId,
			Path:   "video2.mp4",
		},
	}

	mockDynamo.EXPECT().GetByUserId(userId).Return(expectedVideos, nil)

	videos, err := videoRepo.GetByUserId(userId)

	assert.NoError(t, err)
	assert.Equal(t, expectedVideos, videos)
}

func TestVideoRepository_GetByUserId_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDynamo := repository_mocks.NewMockDynamo(ctrl)
	videoRepo := NewVideoRepository(mockDynamo)

	userId := "test-user-id"
	mockDynamo.EXPECT().GetByUserId(userId).Return(nil, errors.New("some error"))

	videos, err := videoRepo.GetByUserId(userId)

	assert.Error(t, err)
	assert.Nil(t, videos)
}

func TestVideoRepository_Create(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDynamo := repository_mocks.NewMockDynamo(ctrl)
	videoRepo := NewVideoRepository(mockDynamo)

	video := entities.Video{UserID: "test-user-id", Path: "video.mp4"}
	mockDynamo.EXPECT().Create(video).Return(&video, nil)

	createdVideo, err := videoRepo.Create(video)

	assert.NoError(t, err)
	assert.Equal(t, &video, createdVideo)
}

func TestVideoRepository_Create_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDynamo := repository_mocks.NewMockDynamo(ctrl)
	videoRepo := NewVideoRepository(mockDynamo)

	video := entities.Video{UserID: "test-user-id", Path: "video.mp4"}
	mockDynamo.EXPECT().Create(video).Return(nil, errors.New("some error"))

	createdVideo, err := videoRepo.Create(video)

	assert.Error(t, err)
	assert.Nil(t, createdVideo)
}

func TestVideoRepository_Update(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDynamo := repository_mocks.NewMockDynamo(ctrl)
	videoRepo := NewVideoRepository(mockDynamo)

	video := entities.Video{UserID: "test-user-id", Path: "video.mp4"}
	mockDynamo.EXPECT().Update(video).Return(&video, nil)

	updatedVideo, err := videoRepo.Update(video)

	assert.NoError(t, err)
	assert.Equal(t, &video, updatedVideo)
}

func TestVideoRepository_Update_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockDynamo := repository_mocks.NewMockDynamo(ctrl)
	videoRepo := NewVideoRepository(mockDynamo)

	video := entities.Video{UserID: "test-user-id", Path: "video.mp4"}
	mockDynamo.EXPECT().Update(video).Return(nil, errors.New("some error"))

	updatedVideo, err := videoRepo.Update(video)

	assert.Error(t, err)
	assert.Nil(t, updatedVideo)
}
