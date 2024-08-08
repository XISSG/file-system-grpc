package router

import (
	"github.com/gin-gonic/gin"
	"github.com/xissg/file-system-grpc/gateway/handler"
	"net/http"
)

func Router() *gin.Engine {
	r := gin.Default()
	r.LoadHTMLGlob("gateway/template/*")
	user := handler.NewUserClient()
	{
		r.GET("/", func(ctx *gin.Context) {
			ctx.HTML(http.StatusFound, "index.html", nil)
		})
		r.GET("/upload", func(ctx *gin.Context) {
			ctx.HTML(http.StatusFound, "upload.html", nil)
		})
		r.GET("/register", func(ctx *gin.Context) {
			ctx.HTML(http.StatusFound, "register.html", nil)
		})
		r.GET("/login", func(ctx *gin.Context) {
			ctx.HTML(http.StatusFound, "login.html", nil)
		})

		r.POST("/register", user.Register)
		r.POST("/login", user.Login)
	}
	return r
}
