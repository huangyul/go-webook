package repository

import (
	"context"

	"github.com/huangyul/go-webook/internal/repository/cache"
	"github.com/huangyul/go-webook/internal/repository/dao"
)

type InteractiveRepository interface {
	IncrReadCnt(ctx context.Context, biz string, bizId int64) error
}

func NewInteractiveRepository(dao dao.InteractiveDAO, cache cache.InteractiveCache) InteractiveRepository {
	return &InteractiveRepositoryImpl{
		dao:   dao,
		cache: cache,
	}
}

type InteractiveRepositoryImpl struct {
	dao   dao.InteractiveDAO
	cache cache.InteractiveCache
}

func (repo *InteractiveRepositoryImpl) IncrReadCnt(ctx context.Context, biz string, bizId int64) error {
	err := repo.dao.IncrReadCnt(ctx, biz, bizId)
	if err != nil {
		return err
	}
	return repo.cache.IncrReadCntIfPresent(ctx, biz, bizId)
}
