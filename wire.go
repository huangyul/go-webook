//go:build wireinject

package main

import (
	"github.com/google/wire"
	articleEvents "github.com/huangyul/go-webook/internal/events/article"
	"github.com/huangyul/go-webook/internal/events/history"
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
	ioc.InitRedis,
	ioc.InitSaramaClient,
	ioc.InitSaramaProducer,
	ioc.InitConsumers)

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

var interactiveSet = wire.NewSet(
	dao.NewInteractiveDAO,
	cache.NewInteractiveCache,
	repository.NewInteractiveRepository,
	service.NewInteractiveService,
)

var historySet = wire.NewSet(
	dao.NewHistoryDAO,
	repository.NewHistoryRepository,
	service.NewHistoryService,
)

func InitApp() *App {
	wire.Build(
		thirdPartySet,

		userSet,
		codeSet,
		smsSet,
		articleSet,
		interactiveSet,
		historySet,

		authz.NewAuthz,

		articleEvents.NewArticleReadProducer,
		articleEvents.NewArticleReadConsumer,
		history.NewHistoryProducer,
		history.NewConsumer,
		ioc.InitMiddlewares,
		ioc.InitWebServer,

		wire.Struct(new(App), "*"),
	)

	return new(App)
}
