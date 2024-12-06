//go:build wireinject

package main

import (
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	"github.com/huangyul/go-blog/internal/ioc"
	"github.com/huangyul/go-blog/internal/repository"
	"github.com/huangyul/go-blog/internal/repository/cache"
	"github.com/huangyul/go-blog/internal/repository/dao"
	"github.com/huangyul/go-blog/internal/service"
	"github.com/huangyul/go-blog/internal/web"
	ginxjwt "github.com/huangyul/go-blog/pkg/ginx/jwt"
)

var (
	UserSet = wire.NewSet(
		dao.NewUserDAOGORM,
		cache.NewRedisUserCache,
		repository.NewUserRepository,
		service.NewUserService,
		web.NewUserHandler,
	)
	CodeSet = wire.NewSet(
		repository.NewCodeRepository,
		cache.NewRedisCodeCache,
		service.NewCodeService,
	)
)

func InitWebServer() *gin.Engine {

	wire.Build(
		ioc.InitDB,
		ioc.InitRedis,

		CodeSet,
		UserSet,

		ginxjwt.NewJWT,

		ioc.InitSMSService,
		ioc.InitGinMiddlewares,
		ioc.InitServer,
	)

	return gin.Default()
}
