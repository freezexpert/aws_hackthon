package models

import (
	"os"
	"testing"
	"time"
)

func TestBedrockAndTTSServices(t *testing.T) {
	// Check for required environment variables
	if os.Getenv("AWS_ACCESS_KEY_ID") == "" || os.Getenv("AWS_SECRET_ACCESS_KEY") == "" {
		t.Skip("AWS credentials not set, skipping Bedrock test")
	}
	if os.Getenv("VYIN_API_KEY") == "" {
		t.Skip("VYIN_API_KEY not set, skipping TTS test")
	}

	// Initialize services
	service, err := New()
	if err != nil {
		t.Fatalf("Failed to initialize services: %v", err)
	}

	// Test Bedrock service
	testPrompt := "Hello, how are you today?"
	response, err := service.GenerateResponse(testPrompt)
	if err != nil {

		t.Fatalf("Bedrock service failed: %v", err)
	}
	t.Logf("Bedrock response: %s", response)

	// Test TTS service
	audioURL, err := service.GenerateSpeech("Hello, how are you?", 2, "max")
	if err != nil {
		t.Fatalf("TTS service failed: %v", err)
	}
	t.Logf("Audio URL: %s", audioURL)

	// Test chat history
	userID := "test_user_123"
	exists, chats := service.Search_chat(userID)
	if !exists {
		// Create new history
		history := History{
			UserID:      userID,
			Type:        "test",
			Chats:       []Chat{},
			LastUpdated: time.Now(),
		}
		if err := service.Create_chat(history); err != nil {
			t.Fatalf("Failed to create chat history: %v", err)
		}
	}

	// Add test chat
	userChat := Chat{
		Role:      "user",
		Content:   testPrompt,
		Time:      time.Now().Format(time.RFC3339),
		Timestamp: time.Now(),
	}

	assistantChat := Chat{
		Role:      "assistant",
		Content:   response,
		Time:      time.Now().Format(time.RFC3339),
		AudioURL:  audioURL,
		Timestamp: time.Now(),
	}

	chats = append(chats, userChat, assistantChat)
	if err := service.Insert_chat(userID, chats); err != nil {
		t.Fatalf("Failed to insert chat: %v", err)
	}

	// Verify chat was saved
	exists, savedChats := service.Search_chat(userID)
	if !exists {
		t.Fatal("Chat history not found after saving")
	}
	if len(savedChats) != 2 {
		t.Fatalf("Expected 2 chats, got %d", len(savedChats))
	}
	t.Logf("Successfully saved and retrieved chat history")
}
