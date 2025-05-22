//go:build wireinject

package main

import (
	"github.com/google/wire"
	"github.com/huangyul/go-webook/interactive/events"
	"github.com/huangyul/go-webook/interactive/grpc"
	"github.com/huangyul/go-webook/interactive/ioc"
	"github.com/huangyul/go-webook/interactive/repository"
	"github.com/huangyul/go-webook/interactive/repository/cache"
	"github.com/huangyul/go-webook/interactive/repository/dao"
	"github.com/huangyul/go-webook/interactive/service"
)

func InitApp() *App {

	wire.Build(
		ioc.InitDB,
		ioc.InitConsumers,
		ioc.InitSaramaClient,
		ioc.InitRedis,
		ioc.InitEtcd,
		ioc.InitGrpcServer,

		grpc.NewInteractiveService,
		events.NewArticleReadConsumer,

		dao.NewInteractiveDAO,
		cache.NewInteractiveCache,
		repository.NewInteractiveRepository,
		service.NewInteractiveService,
		wire.Struct(new(App), "*"),
	)

	return new(App)
}
