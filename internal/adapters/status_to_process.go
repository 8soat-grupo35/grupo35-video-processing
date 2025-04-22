package adapters

import (
	"encoding/json"

	"github.com/8soat-grupo35/grupo35-video-processing/internal/entities"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
)

func NewStatusToProcessFromSQSMessage(message types.Message) (*entities.VideoMessage, error) {
	var convertedMessage entities.VideoMessage

	err := json.Unmarshal([]byte(*message.Body), &convertedMessage)

	if err != nil {
		return nil, err
	}

	return &convertedMessage, nil
}
