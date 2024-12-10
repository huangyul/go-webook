package service

import (
	"context"
	"github.com/huangyul/go-blog/internal/domain"
	"github.com/huangyul/go-blog/internal/repository"
)

type ArticleService interface {
	Save(ctx context.Context, art domain.Article) (int64, error)
	Publish(ctx context.Context, art domain.Article) (int64, error)
	Withdraw(ctx context.Context, uid, id int64) error
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

func (svc *articleService) Withdraw(ctx context.Context, uid int64, id int64) error {
	return svc.repo.SyncStatus(ctx, uid, id, domain.ArticleStatusWithdraw)
}

func (svc *articleService) Save(ctx context.Context, art domain.Article) (int64, error) {
	art.Status = domain.ArticleStatusUnPublished
	if art.ID == 0 {
		return svc.repo.Create(ctx, art)
	} else {
		err := svc.repo.Update(ctx, art)
		return art.ID, err
	}
}

func (svc *articleService) Publish(ctx context.Context, art domain.Article) (int64, error) {
	art.Status = domain.ArticleStatusPublished
	return svc.repo.Sync(ctx, art)
}
