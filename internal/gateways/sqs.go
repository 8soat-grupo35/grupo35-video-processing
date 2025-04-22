package gateways

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/8soat-grupo35/grupo35-video-processing/internal/adapters/wrappers"
	"github.com/8soat-grupo35/grupo35-video-processing/internal/interfaces/repository"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
)

type SQSGateway struct {
	client              wrappers.ISQSClient
	queueName           string
	maxNumberOfMessages int32
}

func NewSQSGateway(client wrappers.ISQSClient, queueName string, maxNumberOfMessages int32) repository.SQS {
	return SQSGateway{
		client:              client,
		queueName:           queueName,
		maxNumberOfMessages: maxNumberOfMessages,
	}
}

func (s SQSGateway) ConsumeMessages(consumeFn func(message types.Message)) error {
	output, err := s.client.ReceiveMessage(context.TODO(), &sqs.ReceiveMessageInput{
		QueueUrl:            aws.String(s.queueName),
		MaxNumberOfMessages: s.maxNumberOfMessages,
	})

	if err != nil {
		return err
	}

	fmt.Println("Getting messages from SQS length: ", len(output.Messages))

	for _, message := range output.Messages {
		consumeFn(message)

		_, err = s.client.DeleteMessage(context.TODO(), &sqs.DeleteMessageInput{
			QueueUrl:      aws.String(s.queueName),
			ReceiptHandle: message.ReceiptHandle,
		})

		if err != nil {
			return err
		}
	}

	return nil
}

func (s SQSGateway) SendMessage(message any) error {
	convertedMessage, err := json.Marshal(message)

	if err != nil {
		return err
	}

	fmt.Printf("sending message to sqs: %v\n", message)

	_, err = s.client.SendMessage(context.TODO(), &sqs.SendMessageInput{
		QueueUrl:    aws.String(s.queueName),
		MessageBody: aws.String(string(convertedMessage)),
	})

	if err != nil {
		return err
	}

	fmt.Printf("Message sent to SQS: %v\n", message)

	return nil
}
