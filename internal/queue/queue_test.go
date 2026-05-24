package queue

import (
	"context"
	"errors"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/stretchr/testify/require"
)

type fakeSendQueueAPI struct {
	gotInput *sqs.SendMessageInput
	err      error
}

func (f *fakeSendQueueAPI) SendMessage(ctx context.Context, params *sqs.SendMessageInput, optFns ...func(*sqs.Options)) (*sqs.SendMessageOutput, error) {
	f.gotInput = params
	return &sqs.SendMessageOutput{}, f.err
}

func TestSendMessage(t *testing.T) {
	fp := &fakeSendQueueAPI{}
	s := New(fp, "test-queue")
	err := s.Send(context.Background(), "abcdef")
	require.NoError(t, err)
	require.Equal(t, "test-queue", *fp.gotInput.QueueUrl)
	require.Equal(t, "abcdef", *fp.gotInput.MessageBody)
}

func TestSendMessage_PropagatesError(t *testing.T) {
	s := New(&fakeSendQueueAPI{err: errors.New("boom")}, "t")
	err := s.Send(context.Background(), "")
	require.ErrorContains(t, err, "boom")
}
