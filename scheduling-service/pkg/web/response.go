package web

import (
	"github.com/gin-gonic/gin"
	"time"
)

type errorResponse struct {
	StatusCode int    `json:"status_code"`
	Status     string `json:"status"`
	Message    string `json:"message"`
	TimeStamp  string `json:"time_stamp"`
}

type response struct {
	Data interface{}
}

func BadResponse(ctx *gin.Context, statusCode int, status, message string) {
	ctx.JSON(statusCode, errorResponse{
		StatusCode: statusCode,
		Status:     status,
		Message:    message,
		TimeStamp:  time.Now().String(),
	})
}

func DeleteResponse(ctx *gin.Context, statusCode int, message string) {
	ctx.JSON(statusCode, errorResponse{
		StatusCode: statusCode,
		Status:     "success",
		Message:    message,
	})
}

func ResponseOK(ctx *gin.Context, statusCode int, data interface{}) {
	ctx.JSON(statusCode, data)
}
