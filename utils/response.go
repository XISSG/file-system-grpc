package utils

import (
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc/codes"
)

func Response(c *gin.Context, statusCode int, code codes.Code, msg string, data any) {
	c.JSON(statusCode, gin.H{
		"code": code,
		"msg":  msg,
		"data": data,
	})
}

func ErrResponse(c *gin.Context, statusCode int, code codes.Code, msg string) {
	c.JSON(statusCode, gin.H{
		"code": code,
		"msg":  msg,
	})
}

const (
	WaitingBackup = 0
	SuccessBackup = 1
)
