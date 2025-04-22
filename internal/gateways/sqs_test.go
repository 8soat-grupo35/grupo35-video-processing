package gateways

import (
	"errors"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	wrappers_mock "github.com/8soat-grupo35/grupo35-video-processing/internal/adapters/wrappers/mock"
)

func TestConsumeMessages_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := wrappers_mock.NewMockISQSClient(ctrl)
	queueName := "test-queue"
	maxNumberOfMessages := int32(10)

	mockMessages := []types.Message{
		{MessageId: aws.String("1"), ReceiptHandle: aws.String("handle-1")},
		{MessageId: aws.String("2"), ReceiptHandle: aws.String("handle-2")},
	}

	mockClient.EXPECT().ReceiveMessage(gomock.Any(), &sqs.ReceiveMessageInput{
		QueueUrl:            aws.String(queueName),
		MaxNumberOfMessages: maxNumberOfMessages,
	}).Return(&sqs.ReceiveMessageOutput{Messages: mockMessages}, nil)

	mockClient.EXPECT().DeleteMessage(gomock.Any(), &sqs.DeleteMessageInput{
		QueueUrl:      aws.String(queueName),
		ReceiptHandle: aws.String("handle-1"),
	}).Return(nil, nil)

	mockClient.EXPECT().DeleteMessage(gomock.Any(), &sqs.DeleteMessageInput{
		QueueUrl:      aws.String(queueName),
		ReceiptHandle: aws.String("handle-2"),
	}).Return(nil, nil)

	gateway := SQSGateway{
		client:              mockClient,
		queueName:           queueName,
		maxNumberOfMessages: maxNumberOfMessages,
	}

	var consumedMessages []types.Message
	err := gateway.ConsumeMessages(func(message types.Message) {
		consumedMessages = append(consumedMessages, message)
	})

	assert.NoError(t, err)
	assert.Equal(t, mockMessages, consumedMessages)
}

func TestConsumeMessages_ReceiveMessageError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := wrappers_mock.NewMockISQSClient(ctrl)
	queueName := "test-queue"
	maxNumberOfMessages := int32(10)

	mockClient.EXPECT().ReceiveMessage(gomock.Any(), &sqs.ReceiveMessageInput{
		QueueUrl:            aws.String(queueName),
		MaxNumberOfMessages: maxNumberOfMessages,
	}).Return(nil, errors.New("receive error"))

	gateway := SQSGateway{
		client:              mockClient,
		queueName:           queueName,
		maxNumberOfMessages: maxNumberOfMessages,
	}

	err := gateway.ConsumeMessages(func(message types.Message) {})

	assert.Error(t, err)
	assert.Equal(t, "receive error", err.Error())
}

func TestConsumeMessages_DeleteMessageError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockClient := wrappers_mock.NewMockISQSClient(ctrl)
	queueName := "test-queue"
	maxNumberOfMessages := int32(10)

	mockMessages := []types.Message{
		{MessageId: aws.String("1"), ReceiptHandle: aws.String("handle-1")},
	}

	mockClient.EXPECT().ReceiveMessage(gomock.Any(), &sqs.ReceiveMessageInput{
		QueueUrl:            aws.String(queueName),
		MaxNumberOfMessages: maxNumberOfMessages,
	}).Return(&sqs.ReceiveMessageOutput{Messages: mockMessages}, nil)

	mockClient.EXPECT().DeleteMessage(gomock.Any(), &sqs.DeleteMessageInput{
		QueueUrl:      aws.String(queueName),
		ReceiptHandle: aws.String("handle-1"),
	}).Return(nil, errors.New("delete error"))

	gateway := SQSGateway{
		client:              mockClient,
		queueName:           queueName,
		maxNumberOfMessages: maxNumberOfMessages,
	}

	err := gateway.ConsumeMessages(func(message types.Message) {})

	assert.Error(t, err)
	assert.Equal(t, "delete error", err.Error())
}
