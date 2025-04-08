package service

import (
	"context"

	"github.com/huangyul/go-webook/internal/repository"
)

type InteractiveService interface {
	IncrReadCnt(ctx context.Context, biz string, bizId int64) error
	Like(ctx context.Context, biz string, bizId int64, userId int64) error
	CancelLike(ctx context.Context, biz string, bizId int64, userId int64) error
}

func NewInteractiveService(repo repository.InteractiveRepository) InteractiveService {
	return &InteractiveServiceImpl{
		repo: repo,
	}
}

type InteractiveServiceImpl struct {
	repo repository.InteractiveRepository
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
