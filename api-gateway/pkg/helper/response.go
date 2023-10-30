package helper

import "github.com/gin-gonic/gin"

type Responses struct {
	Status  string      `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func CustomResponse(ctx *gin.Context, status string, message string, data interface{}) *Responses {
	return &Responses{
		Status:  status,
		Message: message,
		Data:    data,
	}
}
