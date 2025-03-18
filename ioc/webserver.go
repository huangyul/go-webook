package ioc

import (
	"github.com/gin-gonic/gin"
	"github.com/huangyul/go-webook/internal/web"
	"github.com/huangyul/go-webook/internal/web/middleware"
	"github.com/huangyul/go-webook/pkg/ginx/middleware/ratelimit"
	"github.com/redis/go-redis/v9"
	"time"
)

func InitWebServer(mdls []gin.HandlerFunc, userHdl *web.UserHandler) *gin.Engine {
	server := gin.Default()

	server.Use(mdls...)

	userHdl.RegisterRoutes(server)

	return server
}

func InitMiddlewares(cmd redis.Cmdable) []gin.HandlerFunc {
	return []gin.HandlerFunc{
		ratelimit.NewBuilder(cmd,
			ratelimit.SetRate(10000),
			ratelimit.SetInterval(time.Minute*10)).Build(),
		middleware.NewJWTLoginMiddlewareBuild(
			middleware.AddWhiteList("/user/login", "/user/register", "/user/sms/login", "/user/sms/login")).Build(),
	}
}
