package wrappers

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
)

// ISQSClient interface para interagir com o SQS.
type ISQSClient interface {
	ReceiveMessage(ctx context.Context, params *sqs.ReceiveMessageInput, optFns ...func(*sqs.Options)) (*sqs.ReceiveMessageOutput, error)
	DeleteMessage(ctx context.Context, params *sqs.DeleteMessageInput, optFns ...func(*sqs.Options)) (*sqs.DeleteMessageOutput, error)
	SendMessage(ctx context.Context, params *sqs.SendMessageInput, optFns ...func(*sqs.Options)) (*sqs.SendMessageOutput, error)
}

// SQSClient estrutura que implementa ISQSClient.
type SQSClient struct {
	client *sqs.Client
}

// ReceiveMessage implementa ISQSClient.
func (s *SQSClient) ReceiveMessage(ctx context.Context, params *sqs.ReceiveMessageInput, optFns ...func(*sqs.Options)) (*sqs.ReceiveMessageOutput, error) {
	return s.client.ReceiveMessage(ctx, params, optFns...)
}

// DeleteMessage implementa ISQSClient.
func (s *SQSClient) DeleteMessage(ctx context.Context, params *sqs.DeleteMessageInput, optFns ...func(*sqs.Options)) (*sqs.DeleteMessageOutput, error) {
	return s.client.DeleteMessage(ctx, params, optFns...)
}

// SendMessage implementa ISQSClient.
func (s *SQSClient) SendMessage(ctx context.Context, params *sqs.SendMessageInput, optFns ...func(*sqs.Options)) (*sqs.SendMessageOutput, error) {
	return s.client.SendMessage(ctx, params, optFns...)
}

// NewSQSClient cria uma nova instância de SQSClient.
func NewSQSClient(cfg aws.Config, useLocalStack bool) ISQSClient {
	if useLocalStack {
		client := sqs.NewFromConfig(cfg, func(o *sqs.Options) {
			o.EndpointResolver = sqs.EndpointResolverFunc(func(region string, options sqs.EndpointResolverOptions) (aws.Endpoint, error) {
				return aws.Endpoint{
					URL: "http://localhost:4566",
				}, nil
			})
		})
		return &SQSClient{client: client}
	}
	return &SQSClient{
		client: sqs.NewFromConfig(cfg),
	}
}
