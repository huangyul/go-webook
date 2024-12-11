package service

import (
	"context"
	"github.com/huangyul/go-blog/internal/repository"
)

type InteractiveService interface {
	IncrReadCnt(ctx context.Context, bizID int64, biz string) error
	Like(ctx context.Context, userID, bizID int64, biz string) error
	CancelLike(ctx context.Context, userID, bizID int64, biz string) error
}

var _ InteractiveService = (*interactiveService)(nil)

type interactiveService struct {
	repo repository.InteractiveRepository
}

func (svc *interactiveService) Like(ctx context.Context, userID, bizID int64, biz string) error {
	return svc.repo.IncrLikeCnt(ctx, userID, bizID, biz)
}

func (svc *interactiveService) CancelLike(ctx context.Context, userID, bizID int64, biz string) error {
	return svc.repo.DecrLikeCnt(ctx, userID, bizID, biz)
}

func NewInteractiveService(repo repository.InteractiveRepository) InteractiveService {
	return &interactiveService{
		repo: repo,
	}
}

func (svc *interactiveService) IncrReadCnt(ctx context.Context, bizID int64, biz string) error {
	return svc.repo.IncrReadCnt(ctx, bizID, biz)
}
