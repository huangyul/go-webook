//go:build wireinject

package main

import (
	"github.com/google/wire"
	InteractiveEvents "github.com/huangyul/go-webook/interactive/events"
	repository2 "github.com/huangyul/go-webook/interactive/repository"
	cache2 "github.com/huangyul/go-webook/interactive/repository/cache"
	dao2 "github.com/huangyul/go-webook/interactive/repository/dao"
	service2 "github.com/huangyul/go-webook/interactive/service"
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
	ioc.InitConsumers,
	ioc.InitRedisLock)

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

var rankingSet = wire.NewSet(
	cache.NewRankingCache,
	repository.NewRankingRepository,
	service.NewRankingService,
)

var interactiveSet = wire.NewSet(
	dao2.NewInteractiveDAO,
	cache2.NewInteractiveCache,
	repository2.NewInteractiveRepository,
	service2.NewInteractiveService,
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
		rankingSet,

		// grpc client
		ioc.InitInteractiveClient,

		authz.NewAuthz,

		articleEvents.NewArticleReadProducer,
		InteractiveEvents.NewArticleReadConsumer,
		history.NewHistoryProducer,
		history.NewConsumer,
		ioc.InitMiddlewares,
		ioc.InitWebServer,
		ioc.InitRankingJob,
		ioc.InitJobs,

		wire.Struct(new(App), "*"),
	)

	return new(App)
}
