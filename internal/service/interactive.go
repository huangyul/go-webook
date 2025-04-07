package service

import (
	"context"

	"github.com/huangyul/go-webook/internal/repository"
)

type InteractiveService interface {
	IncrReadCnt(ctx context.Context, biz string, bizId int64) error
}

func NewInteractiveService(repo repository.InteractiveRepository) InteractiveService {
	return &InteractiveServiceImpl{
		repo: repo,
	}
}

type InteractiveServiceImpl struct {
	repo repository.InteractiveRepository
}

func (svc *InteractiveServiceImpl) IncrReadCnt(ctx context.Context, biz string, bizId int64) error {
	return svc.repo.IncrReadCnt(ctx, biz, bizId)
}
