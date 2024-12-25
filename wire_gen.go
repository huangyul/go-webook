// Code generated by Wire. DO NOT EDIT.

//go:generate go run -mod=mod github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package main

import (
	"github.com/google/wire"
	"github.com/huangyul/go-blog/interactive/events"
	repository2 "github.com/huangyul/go-blog/interactive/repository"
	cache2 "github.com/huangyul/go-blog/interactive/repository/cache"
	dao2 "github.com/huangyul/go-blog/interactive/repository/dao"
	service2 "github.com/huangyul/go-blog/interactive/service"
	"github.com/huangyul/go-blog/internal/event/article"
	"github.com/huangyul/go-blog/internal/event/history"
	"github.com/huangyul/go-blog/internal/ioc"
	"github.com/huangyul/go-blog/internal/repository"
	"github.com/huangyul/go-blog/internal/repository/cache"
	"github.com/huangyul/go-blog/internal/repository/dao"
	"github.com/huangyul/go-blog/internal/service"
	"github.com/huangyul/go-blog/internal/web"
	"github.com/huangyul/go-blog/pkg/ginx/jwt"
)

// Injectors from wire.go:

func InitApp() *App {
	cmdable := ioc.InitRedis()
	jwt := ginxjwt.NewJWT(cmdable)
	v := ioc.InitGinMiddlewares(cmdable, jwt)
	db := ioc.InitDB()
	userDAO := dao.NewUserDAOGORM(db)
	userCache := cache.NewRedisUserCache(cmdable)
	userRepository := repository.NewUserRepository(userDAO, userCache)
	userService := service.NewUserService(userRepository)
	codeCache := cache.NewRedisCodeCache(cmdable)
	codeRepository := repository.NewCodeRepository(codeCache)
	smsService := ioc.InitSMSService(cmdable)
	codeService := service.NewCodeService(codeRepository, smsService)
	userHandler := web.NewUserHandler(userService, codeService, jwt)
	articleDao := dao.NewArticleDao(db)
	articleCache := cache.NewRedisArticleCache(cmdable)
	articleRepository := repository.NewArticleRepository(articleDao, articleCache)
	client := ioc.InitSaramaClient()
	syncProducer := ioc.InitProducer(client)
	producer := article.NewSaramaSyncProducer(syncProducer)
	historyDao := dao.NewHistoryDao(db)
	historyRepository := repository.NewHistoryRepository(historyDao)
	logger := ioc.InitLogger()
	historyProducer := history.NewSaramaProducer(syncProducer, logger)
	articleService := service.NewArticleService(articleRepository, userRepository, producer, historyRepository, historyProducer, logger)
	interactiveCache := cache2.NewInteractiveCache(cmdable)
	interactiveDao := dao2.NewInteractiveDao(db)
	interactiveRepository := repository2.NewInteractiveRepository(interactiveCache, interactiveDao)
	interactiveService := service2.NewInteractiveService(interactiveRepository)
	interactiveServiceClient := ioc.InitInteractiveGrpcClient(interactiveService)
	articleHandler := web.NewArticleHandler(articleService, interactiveServiceClient)
	engine := ioc.InitServer(v, userHandler, articleHandler)
	interactiveReadEventConsumer := events.NewInteractiveReadEventConsumer(interactiveRepository, client)
	consumer := history.NewConsumer(client, historyRepository, logger)
	v2 := ioc.InitConsumers(interactiveReadEventConsumer, consumer)
	logJob := ioc.InitLogJob()
	cron := ioc.InitJobs(logJob)
	app := &App{
		server:    engine,
		consumers: v2,
		jobs:      cron,
	}
	return app
}

// wire.go:

var (
	UserSet        = wire.NewSet(dao.NewUserDAOGORM, cache.NewRedisUserCache, repository.NewUserRepository, service.NewUserService, web.NewUserHandler)
	CodeSet        = wire.NewSet(repository.NewCodeRepository, cache.NewRedisCodeCache, service.NewCodeService)
	ArticleSet     = wire.NewSet(dao.NewArticleDao, cache.NewRedisArticleCache, repository.NewArticleRepository, service.NewArticleService, web.NewArticleHandler)
	InteractiveSet = wire.NewSet(dao2.NewInteractiveDao, cache2.NewInteractiveCache, repository2.NewInteractiveRepository, service2.NewInteractiveService, events.NewInteractiveReadEventConsumer)
	HistorySet     = wire.NewSet(dao.NewHistoryDao, repository.NewHistoryRepository)
)
