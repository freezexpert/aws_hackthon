package models

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
)

type TTSService interface {
	GenerateSpeech(text string, model_id int, speaker_name string) (string, error)
}

type ttsService struct {
	apiKey string
}

func NewTTSService() (TTSService, error) {
	apiKey := os.Getenv("VYIN_API_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("VYIN_API_KEY environment variable not set")
	}
	// Remove "Bearer " prefix if it exists to avoid duplication
	apiKey = strings.TrimPrefix(apiKey, "Bearer ")
	return &ttsService{apiKey: apiKey}, nil
}

type VyinRequest struct {
	Text        string `json:"text"`
	ModelID     int    `json:"model_id"`
	SpeakerName string `json:"speaker_name"`
}

type VyinResponse struct {
	AudioURL string `json:"audio_url"`
}

func (t *ttsService) GenerateSpeech(text string, model_id int, speaker_name string) (string, error) {
	if text == "" {
		return "", fmt.Errorf("text cannot be empty")
	}
	if speaker_name == "" {
		return "", fmt.Errorf("speaker_name cannot be empty")
	}
	if model_id <= 0 {
		return "", fmt.Errorf("model_id must be positive")
	}
	// 構建查詢參數 URL，補齊 speed_factor 與 mode
	url := fmt.Sprintf("https://uat-persona-sound.data.gamania.com/api/v1/public/voice?text=%s&model_id=%d&speaker_name=%s&speed_factor=1&mode=stream",
		url.QueryEscape(text), model_id, speaker_name)

	// 創建 GET 請求
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}

	// 設置標頭
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", t.apiKey))
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")

	// 發送請求
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// 記錄狀態碼和回應主體
	log.Printf("TTS API status: %d", resp.StatusCode)
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	log.Printf("TTS API response body: %s", string(body))

	// 確保回應是 JSON 格式
	if resp.Header.Get("Content-Type") != "application/json" {
		return "", fmt.Errorf("unexpected content type: %s", resp.Header.Get("Content-Type"))
	}

	// 解析回應
	var response VyinResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return "", err
	}

	return response.AudioURL, nil
}
