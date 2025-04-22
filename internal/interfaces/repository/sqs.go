package repository

import "github.com/aws/aws-sdk-go-v2/service/sqs/types"

//go:generate mockgen -source=sqs.go -destination=mock/sqs.go
type SQS interface {
	ConsumeMessages(consumeFn func(message types.Message)) error
	SendMessage(message any) error
}
