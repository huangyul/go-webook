//go:build wireinject

package main

import (
	"github.com/google/wire"
	"github.com/huangyul/go-blog/interactive/events"
	"github.com/huangyul/go-blog/interactive/grpc"
	"github.com/huangyul/go-blog/interactive/ioc"
	"github.com/huangyul/go-blog/interactive/repository"
	"github.com/huangyul/go-blog/interactive/repository/cache"
	"github.com/huangyul/go-blog/interactive/repository/dao"
	"github.com/huangyul/go-blog/interactive/service"
)

var interactiveWireSet = wire.NewSet(
	dao.NewInteractiveDao,
	cache.NewInteractiveCache,
	repository.NewInteractiveRepository,
	service.NewInteractiveService,
)

func InitApp() *App {
	wire.Build(
		ioc.InitDB,
		ioc.InitConsumers,
		ioc.InitGrpc,
		ioc.InitRedis,
		ioc.InitSaramaClient,
		grpc.NewInteractiveServiceServer,
		events.NewInteractiveReadEventConsumer,
		interactiveWireSet,
		wire.Struct(new(App), "*"),
	)
	return &App{}
}
