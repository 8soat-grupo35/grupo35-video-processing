package usecases

import "github.com/8soat-grupo35/grupo35-video-processing/internal/entities"

type VideoRepository interface {
	GetByUserId(userId string) ([]entities.Video, error)
	Create(video entities.Video) (*entities.Video, error)
	Update(video entities.Video) (*entities.Video, error)
}
