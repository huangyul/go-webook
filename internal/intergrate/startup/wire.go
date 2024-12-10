//go:build wireinject

package startup

import (
	"github.com/google/wire"
	"github.com/huangyul/go-blog/internal/repository"
	"github.com/huangyul/go-blog/internal/repository/dao"
	"github.com/huangyul/go-blog/internal/service"
	"github.com/huangyul/go-blog/internal/web"
)

var thirdPartySet = wire.NewSet(
	InitDB, InitRedis,
)

func InitArticleHandler() *web.ArticleHandler {
	wire.Build(
		thirdPartySet,
		dao.NewArticleDao,
		repository.NewArticleRepository,
		service.NewArticleService,
		web.NewArticleHandler,
	)
	return &web.ArticleHandler{}
}
