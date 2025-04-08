package repository

import (
	"context"

	"github.com/huangyul/go-webook/internal/repository/cache"
	"github.com/huangyul/go-webook/internal/repository/dao"
)

type InteractiveRepository interface {
	IncrReadCnt(ctx context.Context, biz string, bizId int64) error
	IncrLikeCnt(ctx context.Context, biz string, bizId int64, userId int64) error
	DecrLikeCnt(ctx context.Context, biz string, bizId int64, userId int64) error
	IncrCollect(ctx context.Context, biz string, bizId int64, userId int64) error
	DecrCollect(ctx context.Context, biz string, bizId int64, userId int64) error
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

func (repo *InteractiveRepositoryImpl) IncrCollect(ctx context.Context, biz string, bizId int64, userId int64) error {
	err := repo.dao.AddCollectBiz(ctx, biz, bizId, userId)
	if err != nil {
		return err
	}
	return repo.cache.IncrCollectCntIfPresent(ctx, biz, bizId)
}

func (repo *InteractiveRepositoryImpl) DecrCollect(ctx context.Context, biz string, bizId int64, userId int64) error {
	err := repo.dao.DelCollectBiz(ctx, biz, bizId, userId)
	if err != nil {
		return err
	}
	return repo.cache.DecrCollectCntIfPresent(ctx, biz, bizId)
}

func (repo *InteractiveRepositoryImpl) IncrLikeCnt(ctx context.Context, biz string, bizId int64, userId int64) error {
	err := repo.dao.InsertLikeInfo(ctx, biz, bizId, userId)
	if err != nil {
		return err
	}
	return repo.cache.IncrLikeCntIfPresent(ctx, biz, bizId)
}

func (repo *InteractiveRepositoryImpl) DecrLikeCnt(ctx context.Context, biz string, bizId int64, userId int64) error {
	err := repo.dao.DeleteLikeInfo(ctx, biz, bizId, userId)
	if err != nil {
		return err
	}
	return repo.cache.DecrLikeCntIfPresent(ctx, biz, bizId)
}

func (repo *InteractiveRepositoryImpl) IncrReadCnt(ctx context.Context, biz string, bizId int64) error {
	err := repo.dao.IncrReadCnt(ctx, biz, bizId)
	if err != nil {
		return err
	}
	return repo.cache.IncrReadCntIfPresent(ctx, biz, bizId)
}
