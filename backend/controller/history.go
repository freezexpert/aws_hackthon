package controller

import (
	"backend/models"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (ops *BaseController) GetHistory(c *gin.Context) {
	var request models.History
	log.Printf("%v", c.Request)
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON data"})
		return
	}
	log.Println("Valid JSON data")
	ok, history := ops.Service.Search_chat(request.UserID)
	if ok {
		HandleSucccessResponse(c, "", history)
		return
	} else {
		HandleFailedResponse(c, http.StatusNotFound, fmt.Errorf("user %s not found", request.UserID))
	}
}

func (ops *BaseController) PostHistory(c *gin.Context) {
	var request models.History
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON data"})
		return
	}
	log.Println("Valid JSON data")
	is_existed, _ := ops.Service.Search_chat(request.UserID)
	if is_existed {
		log.Print("user already existed")
		err := ops.Service.Insert_chat(request.UserID, request.Chats)
		if err != nil {
			HandleFailedResponse(c, http.StatusBadRequest, err)
			return
		}
		HandleSucccessResponse(c, "")
		return
	} else {
		log.Print("user does not exist")
		err := ops.Service.Create_chat(request)
		if err != nil {
			HandleFailedResponse(c, http.StatusBadRequest, err)
			return
		}
		HandleSucccessResponse(c, "")
		return
	}
}
