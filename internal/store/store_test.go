package store

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/eugenepok/go-aws-url-shortener/pkg/models"
	"github.com/stretchr/testify/require"
)

type fakePutItemAPI struct {
	gotInput *dynamodb.PutItemInput
	err      error
}

func (f *fakePutItemAPI) PutItem(_ context.Context, in *dynamodb.PutItemInput, _ ...func(*dynamodb.Options)) (*dynamodb.PutItemOutput, error) {
	f.gotInput = in
	return &dynamodb.PutItemOutput{}, f.err
}

func TestSave_MarshalsAndCallsPutItem(t *testing.T) {
	fp := &fakePutItemAPI{}
	s := New(fp, "test-table")
	err := s.Save(context.Background(), &models.UrlData{
		FullURL:   "https://www.google.com",
		ShortURL:  "eswp9Xga",
		CreatedAt: time.Now(),
	})
	require.NoError(t, err)
	require.Equal(t, "test-table", *fp.gotInput.TableName)
	require.Contains(t, fp.gotInput.Item, "full_url")
	require.Contains(t, fp.gotInput.Item, "short_url")
}

func TestSave_PropagatesError(t *testing.T) {
	s := New(&fakePutItemAPI{err: errors.New("boom")}, "t")
	err := s.Save(context.Background(), &models.UrlData{})
	require.ErrorContains(t, err, "boom")
}
