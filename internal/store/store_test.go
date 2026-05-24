package store

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/eugenepok/go-aws-url-shortener/pkg/models"
	"github.com/stretchr/testify/require"
)

type fakeDynamoAPI struct {
	gotPutInput  *dynamodb.PutItemInput
	gotGetInput  *dynamodb.GetItemInput
	gotGetOutput *dynamodb.GetItemOutput
	err          error
}

func (f *fakeDynamoAPI) PutItem(_ context.Context, in *dynamodb.PutItemInput, _ ...func(*dynamodb.Options)) (*dynamodb.PutItemOutput, error) {
	f.gotPutInput = in
	return &dynamodb.PutItemOutput{}, f.err
}

func (f *fakeDynamoAPI) GetItemByKey(ctx context.Context, in *dynamodb.GetItemInput, _ ...func(*dynamodb.Options)) (*dynamodb.GetItemOutput, error) {
	f.gotGetInput = in
	return f.gotGetOutput, f.err
}

func TestAddShortURL_MarshalsAndCallsPutItem(t *testing.T) {
	fp := &fakeDynamoAPI{}
	s := New(fp, "test-table")
	err := s.AddShortURL(context.Background(), &models.UrlData{
		FullURL:   "https://www.google.com",
		ShortURL:  "eswp9Xga",
		CreatedAt: time.Now(),
	})
	require.NoError(t, err)
	require.Equal(t, "test-table", *fp.gotPutInput.TableName)
	require.Contains(t, fp.gotPutInput.Item, "full_url")
	require.Contains(t, fp.gotPutInput.Item, "short_url")
}

func TestAddShortURL_PropagatesError(t *testing.T) {
	s := New(&fakeDynamoAPI{err: errors.New("boom")}, "t")
	err := s.AddShortURL(context.Background(), &models.UrlData{})
	require.ErrorContains(t, err, "boom")
}

func TestGetFullURL_Found(t *testing.T) {
	item, err := attributevalue.MarshalMap(&models.UrlData{
		ShortURL: "eswp9Xga",
		FullURL:  "https://www.google.com",
	})
	require.NoError(t, err)

	fp := &fakeDynamoAPI{gotGetOutput: &dynamodb.GetItemOutput{Item: item}}
	s := New(fp, "test-table")

	got, err := s.GetFullURL(context.Background(), "eswp9Xga")
	require.NoError(t, err)
	require.Equal(t, "https://www.google.com", got.FullURL)
	require.Equal(t, "eswp9Xga", got.ShortURL)

	// the right key was sent to dynamo
	require.Equal(t, "test-table", *fp.gotGetInput.TableName)
	require.Contains(t, fp.gotGetInput.Key, "short_url")
}

func TestGetFullURL_NotFound(t *testing.T) {
	// Item == nil means dynamo found no matching key.
	fp := &fakeDynamoAPI{gotGetOutput: &dynamodb.GetItemOutput{}}
	s := New(fp, "test-table")

	_, err := s.GetFullURL(context.Background(), "missing")
	require.ErrorIs(t, err, ErrNotFound)
}

func TestGetFullURL_PropagatesError(t *testing.T) {
	s := New(&fakeDynamoAPI{err: errors.New("boom")}, "t")
	_, err := s.GetFullURL(context.Background(), "x")
	require.ErrorContains(t, err, "boom")
}
