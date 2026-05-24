package store

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/eugenepok/go-aws-url-shortener/pkg/models"
)

type PutItemAPI interface {
	PutItem(ctx context.Context, in *dynamodb.PutItemInput, opts ...func(*dynamodb.Options)) (*dynamodb.PutItemOutput, error)
}

type DynamoStore struct {
	client PutItemAPI
	table  string
}

func New(client PutItemAPI, table string) *DynamoStore {
	return &DynamoStore{client: client, table: table}
}

func (s *DynamoStore) Save(ctx context.Context, m *models.UrlData) error {
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
