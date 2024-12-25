//go:build wireinject

package main

import (
	"github.com/google/wire"
	interEvents "github.com/huangyul/go-blog/interactive/events"
	interRepo "github.com/huangyul/go-blog/interactive/repository"
	interCache "github.com/huangyul/go-blog/interactive/repository/cache"
	interDao "github.com/huangyul/go-blog/interactive/repository/dao"
	interService "github.com/huangyul/go-blog/interactive/service"
	"github.com/huangyul/go-blog/internal/event/article"
	"github.com/huangyul/go-blog/internal/event/history"
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
	ArticleSet = wire.NewSet(
		dao.NewArticleDao,
		cache.NewRedisArticleCache,
		repository.NewArticleRepository,
		service.NewArticleService,
		web.NewArticleHandler,
	)
	InteractiveSet = wire.NewSet(
		interDao.NewInteractiveDao,
		interCache.NewInteractiveCache,
		interRepo.NewInteractiveRepository,
		interService.NewInteractiveService,
		interEvents.NewInteractiveReadEventConsumer,
	)
	HistorySet = wire.NewSet(
		dao.NewHistoryDao,
		repository.NewHistoryRepository,
	)
)

func InitApp() *App {

	wire.Build(
		ioc.InitDB,
		ioc.InitRedis,
		ioc.InitLogger,
		ioc.InitSaramaClient,
		ioc.InitProducer,
		ioc.InitConsumers,
		ioc.InitLogJob,
		ioc.InitJobs,

		ioc.InitInteractiveGrpcClient,

		CodeSet,
		UserSet,
		ArticleSet,
		InteractiveSet,
		HistorySet,

		article.NewSaramaSyncProducer,
		history.NewConsumer,
		history.NewSaramaProducer,

		ginxjwt.NewJWT,

		ioc.InitSMSService,
		ioc.InitGinMiddlewares,
		ioc.InitServer,

		wire.Struct(new(App), "*"),
	)

	return new(App)
}
