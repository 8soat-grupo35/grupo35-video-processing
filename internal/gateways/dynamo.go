package gateways

import (
	"fmt"

	"github.com/8soat-grupo35/grupo35-video-processing/internal/entities"
	"github.com/8soat-grupo35/grupo35-video-processing/internal/external"
)

type DynamoGateway struct {
	dynamo external.DynamoAdapter
}

func NewDynamoGateway(orm external.DynamoAdapter) *DynamoGateway {
	orm.SetTable("videos")
	return &DynamoGateway{
		dynamo: orm,
	}
}
func (d DynamoGateway) GetByUserId(userId string) (videos []entities.Video, err error) {
	value, err := d.dynamo.GetListByKey("user_id", userId)
	if err != nil {
		return nil, err
	}
	for _, item := range value {
		videos = append(videos, *d.convertDynamoToEntity(item))
	}
	return videos, nil
}

func (d DynamoGateway) Create(video entities.Video) (*entities.Video, error) {
	fmt.Printf("Creating video with userId: %s, status: %s, path: %s\n", video.UserID, video.Status, video.Path)
	err := d.dynamo.Create(video)
	if err != nil {
		return nil, err
	}
	fmt.Printf("Video created with userId: %s, status: %s, path: %s\n", video.UserID, video.Status, video.Path)
	return &video, nil
}
func (d DynamoGateway) Update(video entities.Video) (*entities.Video, error) {
	fmt.Printf("Updating video with userId: %s, path %s, status: %s\n", video.UserID, video.Path, video.Status)
	value, err := d.dynamo.UpdateValue("user_id", video.UserID, "path", video.Path, "status", video.Status)
	if err != nil {
		return nil, err
	}
	fmt.Printf("Video updated with userId: %s, status: %s, path: %s\n", video.UserID, video.Status, video.Path)
	return d.convertDynamoToEntity(value), nil
}

func (d DynamoGateway) convertDynamoToEntity(item map[string]any) *entities.Video {
	video := &entities.Video{}

	if userID, ok := item["user_id"].(string); ok {
		video.UserID = userID
	}
	if status, ok := item["status"].(string); ok {
		video.Status = status
	}
	if path, ok := item["path"].(string); ok {
		video.Path = path
	}
	if createdAt, ok := item["created_at"].(string); ok {
		video.CreatedAt = createdAt
	}
	if updatedAt, ok := item["updated_at"].(string); ok {
		video.UpdatedAt = updatedAt
	}

	return video
}
