package service

import (
	"context"
	"github.com/huangyul/go-webook/internal/domain"

	"github.com/huangyul/go-webook/internal/repository"
	"golang.org/x/sync/errgroup"
)

type InteractiveService interface {
	IncrReadCnt(ctx context.Context, biz string, bizId int64) error
	Like(ctx context.Context, biz string, bizId int64, userId int64) error
	CancelLike(ctx context.Context, biz string, bizId int64, userId int64) error
	Collect(ctx context.Context, biz string, bizId int64, userId int64) error
	CancelCollect(ctx context.Context, biz string, bizId int64, userId int64) error
	Get(ctx context.Context, biz string, bizId int64, userId int64) (*domain.Interactive, error)
	GetByIds(ctx context.Context, bix string, ids []int64) (map[int64]domain.Interactive, error)
}

func NewInteractiveService(repo repository.InteractiveRepository) InteractiveService {
	return &InteractiveServiceImpl{
		repo: repo,
	}
}

type InteractiveServiceImpl struct {
	repo repository.InteractiveRepository
}

func (svc *InteractiveServiceImpl) Get(ctx context.Context, biz string, bizId int64, userId int64) (*domain.Interactive, error) {
	res, err := svc.repo.Get(ctx, biz, bizId)
	if err != nil {
		return nil, err
	}
	var eg errgroup.Group
	var er error
	eg.Go(func() error {
		res.Liked, er = svc.repo.Liked(ctx, biz, bizId, userId)
		return er
	})
	eg.Go(func() error {
		res.Collectd, er = svc.repo.Collected(ctx, biz, bizId, userId)
		return er
	})
	return res, eg.Wait()
}

func (svc *InteractiveServiceImpl) Collect(ctx context.Context, biz string, bizId int64, userId int64) error {
	return svc.repo.IncrCollect(ctx, biz, bizId, userId)
}

func (svc *InteractiveServiceImpl) CancelCollect(ctx context.Context, biz string, bizId int64, userId int64) error {
	return svc.repo.DecrCollect(ctx, biz, bizId, userId)
}

func (svc *InteractiveServiceImpl) Like(ctx context.Context, biz string, bizId int64, userId int64) error {
	return svc.repo.IncrLikeCnt(ctx, biz, bizId, userId)
}

func (svc *InteractiveServiceImpl) CancelLike(ctx context.Context, biz string, bizId int64, userId int64) error {
	return svc.repo.DecrLikeCnt(ctx, biz, bizId, userId)
}

func (svc *InteractiveServiceImpl) IncrReadCnt(ctx context.Context, biz string, bizId int64) error {
	return svc.repo.IncrReadCnt(ctx, biz, bizId)
}
