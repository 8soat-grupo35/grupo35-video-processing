package usecases

import (
	"fmt"

	"github.com/8soat-grupo35/grupo35-video-processing/internal/entities"
	"github.com/8soat-grupo35/grupo35-video-processing/internal/interfaces/repository"
)

type ProcessVideoSender struct {
	SQS repository.SQS
}

func NewProcessVideoSender(sqs repository.SQS) ProcessVideoSender {
	return ProcessVideoSender{
		SQS: sqs,
	}
}

func (p ProcessVideoSender) Send(video entities.VideoMessage) error {
	fmt.Println("Sending video to process")
	return p.SQS.SendMessage(video)
}
