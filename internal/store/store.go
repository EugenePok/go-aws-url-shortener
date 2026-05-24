package store

import (
	"context"
	"errors"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/eugenepok/go-aws-url-shortener/pkg/models"
)

// ErrNotFound is returned when no item exists for the given key.
var ErrNotFound = errors.New("store: url not found")

type DynamoAPI interface {
	PutItem(ctx context.Context, in *dynamodb.PutItemInput, opts ...func(*dynamodb.Options)) (*dynamodb.PutItemOutput, error)
	GetItemByKey(ctx context.Context, in *dynamodb.GetItemInput, opts ...func(*dynamodb.Options)) (*dynamodb.GetItemOutput, error)
}

type DynamoStore struct {
	client DynamoAPI
	table  string
}

func New(client DynamoAPI, table string) *DynamoStore {
	return &DynamoStore{client: client, table: table}
}

func (s *DynamoStore) AddShortURL(ctx context.Context, m *models.UrlData) error {
	item, err := attributevalue.MarshalMap(m)
	if err != nil {
		return fmt.Errorf("marshal metadata: %w", err)
	}
	_, err = s.client.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: &s.table,
		Item:      item,
	})
	if err != nil {
		return fmt.Errorf("dynamodb put: %w", err)
	}
	return nil
}

func (s *DynamoStore) GetFullURL(ctx context.Context, shortURL string) (*models.UrlData, error) {
	selectedKey, err := attributevalue.MarshalMap(&models.UrlData{ShortURL: shortURL})
	if err != nil {
		return nil, fmt.Errorf("dynamodb get marshal: %w", err)
	}
	resp, err := s.client.GetItemByKey(ctx, &dynamodb.GetItemInput{
		TableName: &s.table,
		Key:       selectedKey,
	})
	if err != nil {
		return nil, fmt.Errorf("dynamodb get item: %w", err)
	}

	if resp.Item == nil {
		return nil, ErrNotFound
	}

	var item models.UrlData
	if err := attributevalue.UnmarshalMap(resp.Item, &item); err != nil {
		return nil, fmt.Errorf("dynamodb get unmarshal: %w", err)
	}
	return &item, nil
}
