package service

import (
	"context"
	"github.com/huangyul/go-webook/internal/domain"
	"github.com/huangyul/go-webook/internal/repository"
)

type HistoryService interface {
	GetListByUserId(ctx context.Context, userId int64) ([]*domain.History, error)
}

type HistoryRepositoryImpl struct {
	repo       repository.HistoryRepository
	userSvc    UserService
	articleSvc ArticleService
}

func NewHistoryService(repo repository.HistoryRepository, userSvc UserService, articleSvc ArticleService) HistoryService {
	return &HistoryRepositoryImpl{
		repo:       repo,
		userSvc:    userSvc,
		articleSvc: articleSvc,
	}
}

func (svc *HistoryRepositoryImpl) GetListByUserId(ctx context.Context, userId int64) ([]*domain.History, error) {
	res, err := svc.repo.ListByUserId(ctx, userId)
	if err != nil {
		return nil, err
	}
	for _, item := range res {
		user, err := svc.userSvc.FindById(ctx, item.AuthorId)
		if err != nil {
			item.AuthorName = "用户不存在"
		} else {
			item.AuthorName = user.Nickname
		}
		arti, err := svc.articleSvc.GetPudDetailById(ctx, item.ArticleId)
		if err != nil {
			item.ArticleTitle = "文章已删除"
		} else {
			item.ArticleTitle = arti.Title
		}
	}
	return res, nil
}
