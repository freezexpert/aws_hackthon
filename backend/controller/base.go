package controller

import (
	"backend/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

type BaseController struct {
	Service models.Service
}

type ResponseMessage struct {
	// in: body
	Body ResponseContent
}
type ResponseContent struct {
	Status  string      `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func HandleFailedResponse(c *gin.Context, errCode int, err error) {
	if err == nil {
		panic("err is nil")
	}
	c.Error(err)

	if e, ok := err.(*apiError); ok {
		c.AbortWithStatusJSON(errCode, e)
		return
	}

	c.AbortWithStatusJSON(errCode, gin.H{
		"code":    errCode,
		"message": err.Error(),
		"data":    false,
	})
}

// HandleSucccessResponse handles success http response which returns http 200 statusOK code.
func HandleSucccessResponse(c *gin.Context, msg string, arg ...interface{}) {

	if len(c.Errors) > 0 {
		panic("exist errors in context")
	}

	var data interface{}
	if len(arg) > 0 {
		data = arg[0]
	}
	c.JSON(http.StatusOK, newStatusOKResponse(msg, data))
}

func newStatusOKResponse(msg string, data interface{}) *ResponseContent {
	if msg == "" {
		msg = "Success"
	}

	return &ResponseContent{
		Status:  "0",
		Message: msg,
		Data:    data,
	}
}
