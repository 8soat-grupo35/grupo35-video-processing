package usecases

import (
	"github.com/8soat-grupo35/grupo35-video-processing/internal/entities"
)

//go:generate mockgen -source=sender_video_process.go -destination=mock/video_process_sender.go
type VideoProcessSender interface {
	Send(videos entities.Video) error
}
