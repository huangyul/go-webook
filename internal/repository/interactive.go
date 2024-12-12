package repository

import (
	"context"
	"github.com/huangyul/go-blog/internal/domain"
	"github.com/huangyul/go-blog/internal/repository/cache"
	"github.com/huangyul/go-blog/internal/repository/dao"
	"time"
)

type InteractiveRepository interface {
	IncrReadCnt(ctx context.Context, bizID int64, biz string) error
	IncrLikeCnt(ctx context.Context, userID, bizID int64, biz string) error
	DecrLikeCnt(ctx context.Context, userID, bizID int64, biz string) error
	AddCollectItem(ctx context.Context, uid int64, id int64, cid int64, biz string) error
	Get(ctx context.Context, id int64, biz string) (domain.Interactive, error)
	Liked(ctx context.Context, uid int64, id int64, biz string) (bool, error)
	Collect(ctx context.Context, uid int64, id int64, biz string) (bool, error)
}

var _ InteractiveRepository = (*interactiveRepository)(nil)

type interactiveRepository struct {
	cache cache.InteractiveCache
	dao   dao.InteractiveDao
}

func (repo *interactiveRepository) Get(ctx context.Context, id int64, biz string) (domain.Interactive, error) {
	dInt, err := repo.cache.Get(ctx, id, biz)
	if err == nil {
		return dInt, nil
	}
	int, err := repo.dao.Get(ctx, id, biz)
	if err != nil {
		return domain.Interactive{}, err
	}
	dInt = repo.toDomain(int)
	go func() {
		ct, cancel := context.WithTimeout(context.Background(), time.Second*2)
		defer cancel()
		repo.cache.Set(ct, id, biz, dInt)
	}()
	return dInt, err
}

func (repo *interactiveRepository) Liked(ctx context.Context, uid int64, id int64, biz string) (bool, error) {
	like, err := repo.dao.GetLikedInfo(ctx, uid, id, biz)
	if err != nil {
		return false, err
	}
	return like.Status == 1, nil
}

func (repo *interactiveRepository) Collect(ctx context.Context, uid int64, id int64, biz string) (bool, error) {
	_, err := repo.dao.GetCollectInfo(ctx, uid, id, biz)
	if err != nil {
		return false, err
	}
	return true, nil
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

func (repo *interactiveRepository) toDomain(int dao.Interactive) domain.Interactive {
	return domain.Interactive{
		ReadCnt:    int.ReadCnt,
		LikeCnt:    int.LikeCnt,
		CollectCnt: int.CollectCnt,
	}
}
