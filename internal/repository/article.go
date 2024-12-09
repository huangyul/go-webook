package repository

import (
	"context"
	"github.com/huangyul/go-blog/internal/domain"
	"github.com/huangyul/go-blog/internal/repository/dao"
	"time"
)

type ArticleRepository interface {
	Create(ctx context.Context, art domain.Article) (int64, error)
	Update(ctx context.Context, art domain.Article) error
	Sync(ctx context.Context, art domain.Article) (int64, error)
}

var _ ArticleRepository = (*articleRepository)(nil)

type articleRepository struct {
	dao dao.ArticleDao
}

func NewArticleRepository(dao dao.ArticleDao) ArticleRepository {
	return &articleRepository{dao: dao}
}

func (repo *articleRepository) Create(ctx context.Context, art domain.Article) (int64, error) {
	return repo.dao.Insert(ctx, repo.toDao(art))
}

func (repo *articleRepository) Update(ctx context.Context, art domain.Article) error {
	return repo.dao.UpdateByID(ctx, repo.toDao(art))
}

func (repo *articleRepository) Sync(ctx context.Context, art domain.Article) (int64, error) {
	return repo.dao.Sync(ctx, repo.toDao(art))
}

func (repo *articleRepository) toDomain(art dao.Article) domain.Article {
	return domain.Article{
		ID:        art.ID,
		Title:     art.Title,
		Content:   art.Content,
		CreatedAt: time.UnixMilli(art.CreatedAt),
		UpdatedAt: time.UnixMilli(art.UpdatedAt),
		Author: domain.Author{
			ID: art.AuthorID,
		},
	}
}

func (repo *articleRepository) toDao(art domain.Article) dao.Article {
	return dao.Article{
		ID:        art.ID,
		Title:     art.Title,
		Content:   art.Content,
		AuthorID:  art.Author.ID,
		CreatedAt: art.CreatedAt.UnixMilli(),
		UpdatedAt: art.UpdatedAt.UnixMilli(),
	}
}
