package repository

import (
	"context"
	"github.com/huangyul/go-blog/internal/repository/cache"
	"github.com/huangyul/go-blog/internal/repository/dao"
)

type InteractiveRepository interface {
	IncrReadCnt(ctx context.Context, bizID int64, biz string) error
}

var _ InteractiveRepository = (*interactiveRepository)(nil)

type interactiveRepository struct {
	cache cache.InteractiveCache
	dao   dao.InteractiveDao
}

func NewInteractiveRepository(cache cache.InteractiveCache, dao dao.InteractiveDao) InteractiveRepository {
	return &interactiveRepository{
		cache: cache,
		dao:   dao,
	}
}

func (repo *interactiveRepository) IncrReadCnt(ctx context.Context, bizID int64, biz string) error {
	err := repo.dao.IncrReadCnt(ctx, bizID, biz)
	if err != nil {
		return err
	}
	return repo.cache.IncrReadCntIfPresent(ctx, bizID, biz)
}
