package usecases

import (
	"fmt"
	"time"

	"github.com/8soat-grupo35/grupo35-video-processing/internal/adapters"
	"github.com/8soat-grupo35/grupo35-video-processing/internal/entities"
	"github.com/8soat-grupo35/grupo35-video-processing/internal/interfaces/repository"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
)

type StatusVideoConsumer struct {
	sqs    repository.SQS
	dynamo repository.Dynamo
}

func NewStatusVideoConsumer(sqs repository.SQS, dynamo repository.Dynamo) StatusVideoConsumer {
	return StatusVideoConsumer{
		sqs:    sqs,
		dynamo: dynamo,
	}
}

func (s StatusVideoConsumer) ConsumeMessages() {
	for {
		s.sqs.ConsumeMessages(func(message types.Message) {
			statusToProcess, err := adapters.NewStatusToProcessFromSQSMessage(message)

			if err != nil {
				fmt.Println(err)
				return
			}

			s.dynamo.Update(entities.Video{
				UserID:    statusToProcess.User.ID,
				Path:      statusToProcess.VideoPath,
				Status:    statusToProcess.Status,
				UpdatedAt: time.Now().Format(time.RFC3339),
			})
		})
	}
}
