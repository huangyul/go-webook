//go:build wireinject

package main

import (
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	"github.com/huangyul/go-webook/internal/pkg/authz"
	"github.com/huangyul/go-webook/internal/repository"
	"github.com/huangyul/go-webook/internal/repository/cache"
	"github.com/huangyul/go-webook/internal/repository/dao"
	"github.com/huangyul/go-webook/internal/service"
	"github.com/huangyul/go-webook/internal/service/sms"
	"github.com/huangyul/go-webook/internal/web"
	"github.com/huangyul/go-webook/ioc"
)

var thirdPartySet = wire.NewSet(
	ioc.InitDB,
	ioc.InitRedis)

var userSet = wire.NewSet(
	dao.NewUserDAO,
	cache.NewRedisUserCache,
	repository.NewUserRepository,
	service.NewUserService,
	web.NewUserHandler,
)

var codeSet = wire.NewSet(
	cache.NewCodeCache,
	repository.NewCodeRepository,
	service.NewCodeService)

var smsSet = wire.NewSet(
	sms.NewLocalService)

var articleSet = wire.NewSet(
	dao.NewArticleDAO,
	cache.NewArticleCache,
	repository.NewArticleRepository,
	service.NewArticleService,
	web.NewArticleHandler,
)

func InitService() *gin.Engine {
	wire.Build(
		thirdPartySet,

		userSet,
		codeSet,
		smsSet,
		articleSet,

		authz.NewAuthz,

		ioc.InitMiddlewares,
		ioc.InitWebServer,
	)

	return gin.Default()
}
