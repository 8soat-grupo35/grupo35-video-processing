package gateways

import (
	"errors"
	"testing"

	"github.com/8soat-grupo35/grupo35-video-processing/internal/entities"
	mock_external "github.com/8soat-grupo35/grupo35-video-processing/internal/external/mock"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestDynamoGateway_GetByUserId_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAdapter := mock_external.NewMockDynamoAdapter(ctrl)
	userId := "user-1"
	mockResult := []map[string]any{
		{
			"user_id":    userId,
			"path":       "video.mp4",
			"status":     entities.VideoStatusInProcessing,
			"created_at": "2024-01-01T00:00:00Z",
			"updated_at": "2024-01-01T00:00:00Z",
		},
	}
	mockAdapter.EXPECT().
		GetListByKey("user_id", userId).
		Return(mockResult, nil)

	mockAdapter.EXPECT().SetTable("videos").Times(1)
	gateway := NewDynamoGateway(mockAdapter)
	videos, err := gateway.GetByUserId(userId)
	assert.NoError(t, err)
	assert.Len(t, videos, 1)
	assert.Equal(t, userId, videos[0].UserID)
	assert.Equal(t, "video.mp4", videos[0].Path)
	assert.Equal(t, entities.VideoStatusInProcessing, videos[0].Status)
}

func TestDynamoGateway_GetByUserId_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAdapter := mock_external.NewMockDynamoAdapter(ctrl)
	userId := "user-1"
	mockAdapter.EXPECT().
		GetListByKey("user_id", userId).
		Return(nil, errors.New("dynamo error"))

	mockAdapter.EXPECT().SetTable("videos").Times(1)
	gateway := NewDynamoGateway(mockAdapter)
	videos, err := gateway.GetByUserId(userId)
	assert.Error(t, err)
	assert.Nil(t, videos)
}

func TestDynamoGateway_Create_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAdapter := mock_external.NewMockDynamoAdapter(ctrl)
	video := entities.Video{
		UserID:    "user-1",
		Path:      "video.mp4",
		Status:    entities.VideoStatusInProcessing,
		CreatedAt: "2024-01-01T00:00:00Z",
	}
	mockAdapter.EXPECT().
		Create(video).
		Return(nil)
	mockAdapter.EXPECT().SetTable("videos").Times(1)
	gateway := NewDynamoGateway(mockAdapter)
	result, err := gateway.Create(video)
	assert.NoError(t, err)
	assert.Equal(t, &video, result)
}

func TestDynamoGateway_Create_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAdapter := mock_external.NewMockDynamoAdapter(ctrl)
	video := entities.Video{
		UserID:    "user-1",
		Path:      "video.mp4",
		Status:    entities.VideoStatusInProcessing,
		CreatedAt: "2024-01-01T00:00:00Z",
	}
	mockAdapter.EXPECT().
		Create(video).
		Return(errors.New("create error"))

	mockAdapter.EXPECT().SetTable("videos").Times(1)
	gateway := NewDynamoGateway(mockAdapter)
	result, err := gateway.Create(video)
	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestDynamoGateway_Update_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAdapter := mock_external.NewMockDynamoAdapter(ctrl)
	video := entities.Video{
		UserID: "user-1",
		Path:   "video.mp4",
		Status: entities.VideoStatusSuccess,
	}
	mockResult := map[string]any{
		"user_id": video.UserID,
		"path":    video.Path,
		"status":  video.Status,
	}
	mockAdapter.EXPECT().
		UpdateValue("user_id", video.UserID, "path", video.Path, "status", video.Status).
		Return(mockResult, nil)

	mockAdapter.EXPECT().SetTable("videos").Times(1)
	gateway := NewDynamoGateway(mockAdapter)
	result, err := gateway.Update(video)
	assert.NoError(t, err)
	assert.Equal(t, video.UserID, result.UserID)
	assert.Equal(t, video.Path, result.Path)
	assert.Equal(t, video.Status, result.Status)
}

func TestDynamoGateway_Update_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockAdapter := mock_external.NewMockDynamoAdapter(ctrl)
	video := entities.Video{
		UserID: "user-1",
		Path:   "video.mp4",
		Status: entities.VideoStatusSuccess,
	}
	mockAdapter.EXPECT().
		UpdateValue("user_id", video.UserID, "path", video.Path, "status", video.Status).
		Return(nil, errors.New("update error"))

	mockAdapter.EXPECT().SetTable("videos").Times(1)
	gateway := NewDynamoGateway(mockAdapter)
	result, err := gateway.Update(video)
	assert.Error(t, err)
	assert.Nil(t, result)
}
