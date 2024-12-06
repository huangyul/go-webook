package ioc

import (
	"github.com/gin-gonic/gin"
	"github.com/huangyul/go-blog/internal/web"
	"github.com/huangyul/go-blog/internal/web/middleware"
	ginxjwt "github.com/huangyul/go-blog/pkg/ginx/jwt"
	"github.com/huangyul/go-blog/pkg/ginx/middleware/ratelimit"
	"github.com/huangyul/go-blog/pkg/limiter"
	"github.com/redis/go-redis/v9"
	"time"
)

func InitServer(mdls []gin.HandlerFunc, userHdl *web.UserHandler) *gin.Engine {
	server := gin.Default()
	server.Use(mdls...)
	userHdl.RegisterRoutes(server)
	return server
}

func InitGinMiddlewares(cmd redis.Cmdable, jwt ginxjwt.JWT) []gin.HandlerFunc {
	return []gin.HandlerFunc{
		middleware.NewJWTLoginMiddlewareBuild(jwt).AddWhiteList("/user/login", "/user/signup", "/user/login-sms").Build(),
		ratelimit.NewBuilder(limiter.NewRedisSlideWindowRedis(cmd, time.Minute*10, 10)).Build(),
	}
}
