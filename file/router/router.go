package router

import (
	"github.com/gin-gonic/gin"
	"github.com/xissg/file-system-grpc/file/handler"
	"github.com/xissg/file-system-grpc/file/middlewares"
)

func Router() *gin.Engine {
	r := gin.Default()
	file := handler.NewFileClient()
	auth := r.Group("/auth")
	auth.Use(middlewares.CheckLogin())
	{
		auth.POST("/upload", file.Upload)
		auth.POST("/multipart-upload", file.MultipartUpload)
		auth.GET("/files", file.GetFile)
		auth.GET("/download", file.Download)
		auth.GET("/delete", file.Delete)
	}
	return r
}
