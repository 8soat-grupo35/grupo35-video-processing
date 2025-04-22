package wrappers

import (
	"github.com/8soat-grupo35/grupo35-video-processing/internal/entities"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

type IDynamoClient interface {
	Save(video entities.Video) error
	FindByUserID(userID string) ([]entities.Video, error)
}

type DynamoClient struct {
	client *dynamodb.Client
}

func NewDynamoClient(client *dynamodb.Client) IDynamoClient {
	return DynamoClient{
		client: client,
	}
}

func (d DynamoClient) Save(video entities.Video) error {
	// Implement the logic to save the video to DynamoDB
	// This is a placeholder implementation
	return nil
}

func (d DynamoClient) FindByUserID(userID string) ([]entities.Video, error) {
	// Implement the logic to find videos by user ID in DynamoDB
	// This is a placeholder implementation
	return nil, nil
}
