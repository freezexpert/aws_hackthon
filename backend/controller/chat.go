package controller

import (
	"backend/models"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type ChatRequest struct {
	UserID  string `json:"user_id"`
	Message string `json:"message"`
	Type    string `json:"type"`
}

type ChatResponse struct {
	Text     string `json:"text"`
	AudioURL string `json:"audio_url"`
}

func (ops *BaseController) ProcessChat(c *gin.Context) {
	var request ChatRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON data"})
		return
	}

	// Get user history
	exists, chats := ops.Service.Search_chat(request.UserID)
	if !exists {
		// Create new history if user doesn't exist
		history := models.History{
			UserID:      request.UserID,
			Type:        request.Type,
			Chats:       []models.Chat{},
			LastUpdated: time.Now(),
		}
		if err := ops.Service.Create_chat(history); err != nil {
			HandleFailedResponse(c, http.StatusInternalServerError, err)
			return
		}
	}

	// Add user message to history
	userChat := models.Chat{
		Role:      "user",
		Content:   request.Message,
		Time:      time.Now().Format(time.RFC3339),
		Timestamp: time.Now(),
	}

	// Get response from Bedrock
	response, err := ops.Service.GenerateResponse(request.Message)
	if err != nil {
		HandleFailedResponse(c, http.StatusInternalServerError, err)
		return
	}

	// Generate speech from Vyin AI
	audioURL, err := ops.Service.GenerateSpeech(response, 1, "max") // You can customize the voice ID
	if err != nil {
		HandleFailedResponse(c, http.StatusInternalServerError, err)
		return
	}

	// Add assistant response to history
	assistantChat := models.Chat{
		Role:      "assistant",
		Content:   response,
		Time:      time.Now().Format(time.RFC3339),
		AudioURL:  audioURL,
		Timestamp: time.Now(),
	}

	// Update history with new chats
	chats = append(chats, userChat, assistantChat)
	if err := ops.Service.Insert_chat(request.UserID, chats); err != nil {
		HandleFailedResponse(c, http.StatusInternalServerError, err)
		return
	}

	// Return response to frontend
	HandleSucccessResponse(c, "", ChatResponse{
		Text:     response,
		AudioURL: audioURL,
	})
}
