package usecases

import (
	"github.com/8soat-grupo35/grupo35-video-processing/internal/entities"
	"github.com/8soat-grupo35/grupo35-video-processing/internal/interfaces/repository"
)

type VideoRepository struct {
	Dynamo repository.Dynamo
}

func NewVideoRepository(dynamo repository.Dynamo) VideoRepository {
	return VideoRepository{
		Dynamo: dynamo,
	}
}
func (v VideoRepository) GetByUserId(userId string) ([]entities.Video, error) {
	videos, err := v.Dynamo.GetByUserId(userId)
	if err != nil {
		return nil, err
	}
	return videos, nil
}
func (v VideoRepository) Create(video entities.Video) (*entities.Video, error) {
	createdVideo, err := v.Dynamo.Create(video)
	if err != nil {
		return nil, err
	}
	return createdVideo, nil
}
func (v VideoRepository) Update(video entities.Video) (*entities.Video, error) {
	updatedVideo, err := v.Dynamo.Update(video)
	if err != nil {
		return nil, err
	}
	return updatedVideo, nil
}
