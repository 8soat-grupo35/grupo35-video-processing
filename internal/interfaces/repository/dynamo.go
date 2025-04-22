package repository

import "github.com/8soat-grupo35/grupo35-video-processing/internal/entities"

//go:generate mockgen -source=dynamo.go -destination=mock/dynamo.go
type Dynamo interface {
	GetByUserId(userId string) ([]entities.Video, error)
	Create(video entities.Video) (*entities.Video, error)
	Update(video entities.Video) (*entities.Video, error)
}
