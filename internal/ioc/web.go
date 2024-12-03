package ioc

import (
	"github.com/gin-gonic/gin"
	"github.com/huangyul/go-blog/internal/web"
	"github.com/huangyul/go-blog/internal/web/middleware"
)

func InitServer(mdls []gin.HandlerFunc, userHdl *web.UserHandler) *gin.Engine {
	server := gin.Default()
	server.Use(mdls...)
	userHdl.RegisterRoutes(server)
	return server
}

func InitGinMiddlewares() []gin.HandlerFunc {
	return []gin.HandlerFunc{
		middleware.NewJWTLoginMiddlewareBuild().AddWhiteList("/user/login", "/user/signup", "/user/login-sms").Build(),
	}
}
