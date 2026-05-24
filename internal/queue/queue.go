package queue

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
)

type SqsAPI interface {
	SendMessage(ctx context.Context, params *sqs.SendMessageInput, optFns ...func(*sqs.Options)) (*sqs.SendMessageOutput, error)
}

type SqsQueue struct {
	client   SqsAPI
	queueURL string
}

func New(client SqsAPI, queueURL string) *SqsQueue {
	return &SqsQueue{client: client, queueURL: queueURL}
}

func (q *SqsQueue) Send(ctx context.Context, msg string) error {
	_, err := q.client.SendMessage(ctx, &sqs.SendMessageInput{
		MessageBody: aws.String(msg),
		QueueUrl:    aws.String(q.queueURL),
	})
	if err != nil {
		return fmt.Errorf("sqs send: %w", err)
	}
	return nil
}
