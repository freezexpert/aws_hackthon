package models

import (
	"context"
	"log"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type History struct {
	UserID      string    `json:"user_id" dynamodbav:"user_id"`
	Type        string    `json:"type" dynamodbav:"type"`
	Chats       []Chat    `json:"chats" dynamodbav:"chats"`
	VoiceID     string    `json:"voice_id" dynamodbav:"voice_id"`
	LastUpdated time.Time `json:"last_updated" dynamodbav:"last_updated"`
}

type Chat struct {
	Role      string    `json:"role" dynamodbav:"role"`
	Content   string    `json:"content" dynamodbav:"content"`
	Time      string    `json:"time" dynamodbav:"time"`
	AudioURL  string    `json:"audio_url" dynamodbav:"audio_url"`
	Timestamp time.Time `json:"timestamp" dynamodbav:"timestamp"`
}

type HistoryService interface {
	Search_chat(id string) (bool, []Chat)
	Create_chat(his History) error
	Insert_chat(id string, chats []Chat) error
	// Delete_db(id string) error

}

func (t *controllerOps) Search_chat(id string) (bool, []Chat) {
	result, err := t.Client.GetItem(context.TODO(), &dynamodb.GetItemInput{
		TableName: aws.String("History"),
		Key: map[string]types.AttributeValue{
			"user_id": &types.AttributeValueMemberS{Value: id},
		},
	})
	if err != nil {
		log.Printf("Error getting item: %v", err)
		return false, nil
	}

	if result.Item == nil {
		return false, nil
	}

	var history History
	err = attributevalue.UnmarshalMap(result.Item, &history)
	if err != nil {
		log.Printf("Error unmarshaling item: %v", err)
		return false, nil
	}

	return true, history.Chats
}

func (t *controllerOps) Create_chat(his History) error {
	item, err := attributevalue.MarshalMap(his)
	if err != nil {
		return err
	}

	_, err = t.Client.PutItem(context.TODO(), &dynamodb.PutItemInput{
		TableName: aws.String("History"),
		Item:      item,
	})
	return err
}

func (t *controllerOps) Insert_chat(id string, chats []Chat) error {
	history, err := t.getHistory(id)
	if err != nil {
		return err
	}

	history.Chats = chats
	return t.updateHistory(history)
}

func (t *controllerOps) getHistory(id string) (*History, error) {
	result, err := t.Client.GetItem(context.TODO(), &dynamodb.GetItemInput{
		TableName: aws.String("History"),
		Key: map[string]types.AttributeValue{
			"user_id": &types.AttributeValueMemberS{Value: id},
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

func (t *controllerOps) updateHistory(history *History) error {
	item, err := attributevalue.MarshalMap(history)
	if err != nil {
		return err
	}

	_, err = t.Client.PutItem(context.TODO(), &dynamodb.PutItemInput{
		TableName: aws.String("History"),
		Item:      item,
	})
	return err
}
