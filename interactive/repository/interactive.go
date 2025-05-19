package repository

import (
	"context"
	"errors"

	"github.com/huangyul/go-webook/interactive/domain"
	"github.com/huangyul/go-webook/interactive/repository/cache"
	"github.com/huangyul/go-webook/interactive/repository/dao"
	"gorm.io/gorm"
)

type InteractiveRepository interface {
	IncrReadCnt(ctx context.Context, biz string, bizId int64) error
	IncrLikeCnt(ctx context.Context, biz string, bizId int64, userId int64) error
	DecrLikeCnt(ctx context.Context, biz string, bizId int64, userId int64) error
	IncrCollect(ctx context.Context, biz string, bizId int64, userId int64) error
	DecrCollect(ctx context.Context, biz string, bizId int64, userId int64) error
	Get(ctx context.Context, biz string, bizId int64) (*domain.Interactive, error)
	Liked(ctx context.Context, biz string, bizId int64, userId int64) (bool, error)
	Collected(ctx context.Context, biz string, bizId int64, userId int64) (bool, error)
	GetByIds(ctx context.Context, biz string, ids []int64) ([]*domain.Interactive, error)
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

// GetByIds
func (repo *InteractiveRepositoryImpl) GetByIds(ctx context.Context, biz string, ids []int64) ([]*domain.Interactive, error) {
	res, err := repo.dao.GetByIds(ctx, biz, ids)
	if err != nil {
		return nil, err
	}
	interactives := make([]*domain.Interactive, 0, len(res))
	for _, r := range res {
		interactives = append(interactives, repo.toDomain(r))
	}
	return interactives, nil
}

func (repo *InteractiveRepositoryImpl) Get(ctx context.Context, biz string, bizId int64) (*domain.Interactive, error) {
	r, err := repo.cache.Get(ctx, biz, bizId)
	if err == nil {
		return r, nil
	}
	res, err := repo.dao.Get(ctx, biz, bizId)
	if err != nil {
		return nil, err
	}
	r = repo.toDomain(res)
	go func() {
		repo.cache.Set(ctx, biz, bizId, r)
	}()
	return r, nil
}

func (repo *InteractiveRepositoryImpl) Liked(ctx context.Context, biz string, bizId int64, userId int64) (bool, error) {
	_, err := repo.dao.GetLikeInfo(ctx, biz, bizId, userId)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}

func (repo *InteractiveRepositoryImpl) Collected(ctx context.Context, biz string, bizId int64, userId int64) (bool, error) {
	_, err := repo.dao.GetCollectInfo(ctx, biz, bizId, userId)
	switch {
	case errors.Is(err, gorm.ErrRecordNotFound):
		return false, nil
	case err == nil:
		return true, nil
	default:
		return false, err
	}
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

func (repo *InteractiveRepositoryImpl) toDomain(inter *dao.Interactive) *domain.Interactive {
	return &domain.Interactive{
		Id:         inter.Id,
		CollectCnt: inter.CollectCnt,
		LikeCnt:    inter.LikeCnt,
		ReadCnt:    inter.ReadCnt,
	}
}
