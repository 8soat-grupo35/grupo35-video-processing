package usecases

import (
	"github.com/8soat-grupo35/grupo35-video-processing/internal/entities"
)

//go:generate mockgen -source=process_video_sender.go -destination=mock/process_video_sender.go
type ProcessVideoSender interface {
	Send(videos entities.VideoMessage) error
}
