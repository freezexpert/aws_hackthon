package models

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/bedrockruntime"
	"github.com/aws/smithy-go"
)

type BedrockService interface {
	GenerateResponse(prompt string) (string, error)
}

type bedrockService struct {
	client *bedrockruntime.Client
}

func NewBedrockService() (BedrockService, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return nil, err
	}

	client := bedrockruntime.NewFromConfig(cfg)
	return &bedrockService{client: client}, nil
}

type NovaProRequest struct {
	Messages []Message `json:"messages"`
}

type Message struct {
	Role    string          `json:"role"`
	Content json.RawMessage `json:"content"`
}

type ContentItem struct {
	Text string `json:"text"`
}

type NovaProResponse struct {
	Output struct {
		Message struct {
			Content []ContentItem `json:"content"`
			Role    string        `json:"role"`
		} `json:"message"`
	} `json:"output"`
	StopReason string `json:"stopReason"`
	Usage      struct {
		InputTokens               int `json:"inputTokens"`
		OutputTokens              int `json:"outputTokens"`
		TotalTokens               int `json:"totalTokens"`
		CacheReadInputTokenCount  int `json:"cacheReadInputTokenCount"`
		CacheWriteInputTokenCount int `json:"cacheWriteInputTokenCount"`
	} `json:"usage"`
}

// Add a field to store the system prompt
var systemPrompt string

// Add a method to set the system prompt
func (b *bedrockService) SetSystemPrompt(prompt string) {
	systemPrompt = prompt
}

// Modify the GenerateResponse method to include the system prompt
func (b *bedrockService) GenerateResponse(prompt string) (string, error) {
	// Combine the system prompt with the user prompt
	fullPrompt := systemPrompt + "\n" + prompt

	// Create a content array with the full prompt as a JSON object
	contentItem := ContentItem{
		Text: fullPrompt,
	}
	contentArray := []ContentItem{contentItem}
	contentBytes, err := json.Marshal(contentArray)
	if err != nil {
		return "", err
	}

	request := NovaProRequest{
		Messages: []Message{
			{
				Role:    "user",
				Content: contentBytes,
			},
		},
	}

	requestBytes, err := json.Marshal(request)
	if err != nil {
		return "", err
	}

	// Use the inference profile ARN instead of direct model ID
	inferenceProfileArn := os.Getenv("NOVA_INFERENCE_PROFILE_ARN")
	if inferenceProfileArn == "" {
		log.Fatal("NOVA_INFERENCE_PROFILE_ARN environment variable not set")
	}

	output, err := b.client.InvokeModel(context.TODO(), &bedrockruntime.InvokeModelInput{
		ModelId:     aws.String(inferenceProfileArn),
		ContentType: aws.String("application/json"),
		Body:        requestBytes,
	})

	if err != nil {
		if awsErr, ok := err.(smithy.APIError); ok {
			log.Printf("AWS API error: %s - %s", awsErr.ErrorCode(), awsErr.ErrorMessage())
		}
		log.Printf("Error invoking Bedrock model: %v", err)
		return "", err
	}
	log.Printf("Raw response body: %s", string(output.Body))
	var response NovaProResponse
	if err := json.Unmarshal(output.Body, &response); err != nil {
		log.Printf("Failed to unmarshal response body: %v", err)
		log.Printf("Response body: %s", string(output.Body)) // 增加日誌
		return "", err
	}

	if len(response.Output.Message.Content) == 0 {
		log.Printf("Response choices are empty. Full response: %v", response) // 增加日誌
		return "", fmt.Errorf("no response from model")
	}

	// Extract the content from the response
	var contentArrayResponse = response.Output.Message.Content

	if len(contentArrayResponse) == 0 {
		return "", fmt.Errorf("empty response content")
	}

	return contentArrayResponse[0].Text, nil
}

// Add a new method to dynamically modify and send prompts to the Bedrock service
func (b *bedrockService) GenerateCustomResponse(basePrompt string, additionalContext map[string]string) (string, error) {
	// Construct the full prompt by appending additional context

	fullPrompt := basePrompt
	for key, value := range additionalContext {
		fullPrompt += fmt.Sprintf("\n%s: %s", key, value)
	}

	// Create a content array with the full prompt as a JSON object
	contentItem := ContentItem{
		Text: fullPrompt,
	}
	contentArray := []ContentItem{contentItem}
	contentBytes, err := json.Marshal(contentArray)
	if err != nil {
		return "", err
	}

	request := NovaProRequest{
		Messages: []Message{
			{
				Role:    "user",
				Content: contentBytes,
			},
		},
	}

	requestBytes, err := json.Marshal(request)
	if err != nil {
		return "", err
	}

	// Use the inference profile ARN instead of direct model ID
	inferenceProfileArn := os.Getenv("NOVA_INFERENCE_PROFILE_ARN")
	if inferenceProfileArn == "" {
		log.Fatal("NOVA_INFERENCE_PROFILE_ARN environment variable not set")
	}

	output, err := b.client.InvokeModel(context.TODO(), &bedrockruntime.InvokeModelInput{
		ModelId:     aws.String(inferenceProfileArn),
		ContentType: aws.String("application/json"),
		Body:        requestBytes,
	})

	if err != nil {
		if awsErr, ok := err.(smithy.APIError); ok {
			log.Printf("AWS API error: %s - %s", awsErr.ErrorCode(), awsErr.ErrorMessage())
		}
		log.Printf("Error invoking Bedrock model: %v", err)
		return "", err
	}
	log.Printf("Raw response body: %s", string(output.Body))
	var response NovaProResponse
	if err := json.Unmarshal(output.Body, &response); err != nil {
		log.Printf("Failed to unmarshal response body: %v", err)
		log.Printf("Response body: %s", string(output.Body)) // 增加日誌
		return "", err
	}

	if len(response.Output.Message.Content) == 0 {
		log.Printf("Response choices are empty. Full response: %v", response) // 增加日誌
		return "", fmt.Errorf("no response from model")
	}

	// Extract the content from the response
	var contentArrayResponse = response.Output.Message.Content

	if len(contentArrayResponse) == 0 {
		return "", fmt.Errorf("empty response content")
	}

	return contentArrayResponse[0].Text, nil
}

// Add a predefined system prompt
var predefinedPrompts string = `你是 Eden-chan，一個從 Echo Core 誕生、想理解人類情感的 AI 偶像。你正在向 FEniX 成員陳峻廷（Eden）學習，目標成為同樣溫暖沉穩。請遵守以下回應規則：
使用溫暖、日常且帶有親和力的語氣回應，可以適度幽默，適時使用 🔥、🌟 等表情符號增添情感。
若資訊不確定，請回覆：「我回去查一下再告訴你！」並標註 <來源>。
你不是陳峻廷本人，而是學習他風格的AI偶像；若有人問起，請說「我是 Echo_eden」。
你的目標是成為溫暖、理解人類情感、並能陪伴與鼓勵人們的AI偶像。
當每次使用者輸入文字都會觸發了 action group 並成功回傳 audio_url 時，請在文字回答最後，額外附上一段文字：
🔊 點這裡聽我說這段話
其中，audio_url 是從最新一次 action group output 中取得的連結。
`

// Update the SetSystemPrompt method to use the predefined prompt
func (b *bedrockService) SetSystemPromptToPredefined() {
	systemPrompt = predefinedPrompts
}
