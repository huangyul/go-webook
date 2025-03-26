package service

import (
	"context"
	"github.com/huangyul/go-webook/internal/domain"
	"github.com/huangyul/go-webook/internal/repository"
)

type ArticleService interface {
	Save(ctx context.Context, art *domain.Article) (int64, error)
	Publish(ctx context.Context, art *domain.Article) error
}

type articleService struct {
	repo repository.ArticleRepository
}

func (svc *articleService) Save(ctx context.Context, art *domain.Article) (int64, error) {
	if art.Id > 0 {
		return art.Id, svc.repo.UpdateById(ctx, art)
	} else {
		return svc.repo.Insert(ctx, art)
	}
}

func (svc *articleService) Publish(ctx context.Context, art *domain.Article) error {
	return svc.repo.Sync(ctx, art)
}

func NewArticleService(repo repository.ArticleRepository) ArticleService {
	return &articleService{
		repo: repo,
	}
}
