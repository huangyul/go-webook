package repository

import (
	"context"
	"github.com/huangyul/go-blog/internal/repository/cache"
	"github.com/huangyul/go-blog/internal/repository/dao"
)

type InteractiveRepository interface {
	IncrReadCnt(ctx context.Context, bizID int64, biz string) error
	IncrLikeCnt(ctx context.Context, userID, bizID int64, biz string) error
	DecrLikeCnt(ctx context.Context, userID, bizID int64, biz string) error
	AddCollectItem(ctx context.Context, uid int64, id int64, cid int64, biz string) error
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

func (repo *interactiveRepository) IncrLikeCnt(ctx context.Context, userID, bizID int64, biz string) error {
	err := repo.dao.InsertLikeBiz(ctx, userID, bizID, biz)
	if err != nil {
		return err
	}
	return repo.cache.IncrLikeCntIfPresent(ctx, bizID, biz)
}

func (repo *interactiveRepository) AddCollectItem(ctx context.Context, uid int64, id int64, cid int64, biz string) error {
	err := repo.dao.InsertCollectBiz(ctx, uid, id, cid, biz)
	if err != nil {
		return err
	}
	return repo.cache.IncrCollectCntIfPresent(ctx, uid, id, biz)
}

func (repo *interactiveRepository) DecrLikeCnt(ctx context.Context, userID, bizID int64, biz string) error {
	err := repo.dao.DeleteLikeBiz(ctx, userID, bizID, biz)
	if err != nil {
		return err
	}
	return repo.cache.DecrLikeCntIfPresent(ctx, bizID, biz)
}

func (repo *interactiveRepository) IncrReadCnt(ctx context.Context, bizID int64, biz string) error {
	err := repo.dao.IncrReadCnt(ctx, bizID, biz)
	if err != nil {
		return err
	}
	return repo.cache.IncrReadCntIfPresent(ctx, bizID, biz)
}
