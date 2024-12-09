package service

import (
	"context"
	"github.com/huangyul/go-blog/internal/domain"
	"github.com/huangyul/go-blog/internal/repository"
)

type ArticleService interface {
	Save(ctx context.Context, art domain.Article) (int64, error)
	Publish(ctx context.Context, art domain.Article) (int64, error)
}

var _ ArticleService = (*articleService)(nil)

type articleService struct {
	repo repository.ArticleRepository
}

func NewArticleService(repo repository.ArticleRepository) ArticleService {
	return &articleService{
		repo: repo,
	}
}

func (svc *articleService) Save(ctx context.Context, art domain.Article) (int64, error) {
	if art.ID == 0 {
		return svc.repo.Create(ctx, art)
	} else {
		err := svc.repo.Update(ctx, art)
		return art.ID, err
	}
}

func (svc *articleService) Publish(ctx context.Context, art domain.Article) (int64, error) {
	return svc.repo.Sync(ctx, art)
}
