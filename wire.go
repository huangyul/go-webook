//go:build wireinject

package main

import (
	"github.com/google/wire"
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
		dao.NewInteractiveDao,
		cache.NewInteractiveCache,
		repository.NewInteractiveRepository,
		service.NewInteractiveService,
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

		CodeSet,
		UserSet,
		ArticleSet,
		InteractiveSet,
		HistorySet,

		article.NewSaramaSyncProducer,
		article.NewInteractiveReadConsumer,
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
