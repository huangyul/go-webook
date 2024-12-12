package service

import (
	"context"
	"github.com/huangyul/go-blog/internal/domain"
	"github.com/huangyul/go-blog/internal/repository"
	"golang.org/x/sync/errgroup"
)

type InteractiveService interface {
	IncrReadCnt(ctx context.Context, bizID int64, biz string) error
	Like(ctx context.Context, userID, bizID int64, biz string) error
	CancelLike(ctx context.Context, userID, bizID int64, biz string) error
	Collect(ctx context.Context, uid int64, id int64, cid int64, biz string) error
	Get(ctx context.Context, uid int64, id int64, biz string) (domain.Interactive, error)
}

var _ InteractiveService = (*interactiveService)(nil)

type interactiveService struct {
	repo repository.InteractiveRepository
}

func NewInteractiveService(repo repository.InteractiveRepository) InteractiveService {
	return &interactiveService{
		repo: repo,
	}
}

func (svc *interactiveService) Get(ctx context.Context, uid int64, id int64, biz string) (domain.Interactive, error) {
	inte, err := svc.repo.Get(ctx, id, biz)
	if err != nil {
		return domain.Interactive{}, err
	}
	var eg errgroup.Group
	eg.Go(func() error {
		var er error
		inte.Liked, er = svc.repo.Liked(ctx, uid, id, biz)
		return er
	})
	eg.Go(func() error {
		var er error
		inte.Collected, er = svc.repo.Collect(ctx, uid, id, biz)
		return er
	})
	err = eg.Wait()
	if err != nil {
		return domain.Interactive{}, err
	}
	return inte, nil
}

func (svc *interactiveService) Collect(ctx context.Context, uid int64, id int64, cid int64, biz string) error {
	return svc.repo.AddCollectItem(ctx, uid, id, id, biz)
}

func (svc *interactiveService) Like(ctx context.Context, userID, bizID int64, biz string) error {
	return svc.repo.IncrLikeCnt(ctx, userID, bizID, biz)
}

func (svc *interactiveService) CancelLike(ctx context.Context, userID, bizID int64, biz string) error {
	return svc.repo.DecrLikeCnt(ctx, userID, bizID, biz)
}

func (svc *interactiveService) IncrReadCnt(ctx context.Context, bizID int64, biz string) error {
	return svc.repo.IncrReadCnt(ctx, bizID, biz)
}
