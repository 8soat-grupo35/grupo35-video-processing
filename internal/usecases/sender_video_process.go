package usecases

import (
	"encoding/json"
	"fmt"

	"github.com/8soat-grupo35/grupo35-video-processing/internal/entities"
	"github.com/8soat-grupo35/grupo35-video-processing/internal/interfaces/repository"
)

type VideoProcessSender struct {
	SQS repository.SQS
}

func NewVideoProcessSender(sqs repository.SQS) VideoProcessSender {
	return VideoProcessSender{
		SQS: sqs,
	}
}

func (v VideoProcessSender) Send(video entities.Video) error {
	fmt.Println("Sending video to process")
	// Serializa os dados do vídeo para JSON
	message, err := json.Marshal(video)
	if err != nil {
		return fmt.Errorf("failed to marshal video: %w", err)
	}
	return v.SQS.SendMessage(message)
}
