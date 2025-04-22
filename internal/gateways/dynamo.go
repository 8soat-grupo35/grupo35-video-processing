package gateways

import (
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
	value, err := d.dynamo.GetListByKey("userId", userId)

	if err != nil {
		return nil, err
	}

	for _, item := range value {
		videos = append(videos, *d.convertDynamoToEntity(item))
	}

	return videos, nil
}

func (d DynamoGateway) Create(video entities.Video) (*entities.Video, error) {
	err := d.dynamo.Create(video)
	if err != nil {
		return nil, err
	}
	return &video, err
}
func (d DynamoGateway) Update(video entities.Video) (*entities.Video, error) {
	value, err := d.dynamo.UpdateValue("userId", video.UserID, "Status", video.Status)
	if err != nil {
		return nil, err
	}
	return d.convertDynamoToEntity(value), nil
}

func (d DynamoGateway) convertDynamoToEntity(item map[string]interface{}) *entities.Video {
	return &entities.Video{
		UserID:    item["user_id"].(string),
		Status:    item["status"].(string),
		Path:      item["path"].(string),
		CreatedAt: item["created_at"].(string),
		UpdatedAt: item["updated_at"].(string),
	}
}
