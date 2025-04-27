package models

import (
	"context"
	"log"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

type Service interface {
	HistoryService
	BedrockService
	TTSService
}

type service struct {
	*controllerOps
	bedrockService BedrockService
	ttsService     TTSService
}

type controllerOps struct {
	*dynamodb.Client
}

// New returns a Service instance for operating all model service.
func New() (Service, error) {
	client, err := GetDynamoDBClient()
	if err != nil {
		return nil, err
	}

	bedrockService, err := NewBedrockService()
	if err != nil {
		return nil, err
	}

	ttsService, err := NewTTSService()
	if err != nil {
		return nil, err
	}

	serv := &service{
		controllerOps:  &controllerOps{client},
		bedrockService: bedrockService,
		ttsService:     ttsService,
	}

	return serv, nil
}

func (s *service) GenerateResponse(prompt string) (string, error) {
	return s.bedrockService.GenerateResponse(prompt)
}

func (s *service) GenerateSpeech(text string, model_id int, speaker_name string) (string, error) {
	return s.ttsService.GenerateSpeech(text, model_id, speaker_name)
}

func GetDynamoDBClient() (*dynamodb.Client, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return nil, err
	}

	client := dynamodb.NewFromConfig(cfg)
	log.Println("Successfully connected to DynamoDB!")
	return client, nil
}
