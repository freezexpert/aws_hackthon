package models

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type DynamoDBClient struct {
	client *dynamodb.Client
}

func NewDynamoDBClient() (*DynamoDBClient, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return nil, err
	}

	client := dynamodb.NewFromConfig(cfg)
	return &DynamoDBClient{client: client}, nil
}

func (d *DynamoDBClient) GetHistory(userID string) (*History, error) {
	result, err := d.client.GetItem(context.TODO(), &dynamodb.GetItemInput{
		TableName: aws.String("History"),
		Key: map[string]types.AttributeValue{
			"user_id": &types.AttributeValueMemberS{Value: userID},
		},
	})
	if err != nil {
		return nil, err
	}

	if result.Item == nil {
		return nil, nil
	}

	var history History
	err = attributevalue.UnmarshalMap(result.Item, &history)
	if err != nil {
		return nil, err
	}

	return &history, nil
}

func (d *DynamoDBClient) CreateHistory(history History) error {
	item, err := attributevalue.MarshalMap(history)
	if err != nil {
		return err
	}

	_, err = d.client.PutItem(context.TODO(), &dynamodb.PutItemInput{
		TableName: aws.String("History"),
		Item:      item,
	})
	return err
}

func (d *DynamoDBClient) UpdateHistory(userID string, history History) error {
	item, err := attributevalue.MarshalMap(history)
	if err != nil {
		return err
	}

	_, err = d.client.PutItem(context.TODO(), &dynamodb.PutItemInput{
		TableName: aws.String("History"),
		Item:      item,
	})
	return err
}
