package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type BedrockRequest struct {
	Prompt string `json:"prompt"`
}

func (ops *BaseController) GenerateResponse(c *gin.Context) {
	var request BedrockRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON data"})
		return
	}

	response, err := ops.Service.GenerateResponse(request.Prompt)
	if err != nil {
		HandleFailedResponse(c, http.StatusInternalServerError, err)
		return
	}

	HandleSucccessResponse(c, "", response)
} 