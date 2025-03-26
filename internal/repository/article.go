package repository

import (
	"context"
	"github.com/huangyul/go-webook/internal/domain"
	"github.com/huangyul/go-webook/internal/repository/dao"
)

type ArticleRepository interface {
	Insert(ctx context.Context, art *domain.Article) (int64, error)
	UpdateById(ctx context.Context, art *domain.Article) error
	Sync(ctx context.Context, art *domain.Article) error
}

type articleRepository struct {
	dao dao.ArticleDAO
}

func (a *articleRepository) Insert(ctx context.Context, art *domain.Article) (int64, error) {
	return a.dao.Insert(ctx, a.toEntity(art))
}

func (a *articleRepository) UpdateById(ctx context.Context, art *domain.Article) error {
	return a.dao.UpdateById(ctx, a.toEntity(art))
}

func (a *articleRepository) Sync(ctx context.Context, art *domain.Article) error {
	return a.dao.Sync(ctx, a.toEntity(art))
}

func NewArticleRepository(dao dao.ArticleDAO) ArticleRepository {
	return &articleRepository{
		dao: dao,
	}
}

func (a *articleRepository) toEntity(art *domain.Article) *dao.Article {
	return &dao.Article{
		Id:       art.Id,
		Title:    art.Title,
		Content:  art.Content,
		AuthorId: art.Author.Id,
	}
}

func (a *articleRepository) toDomain(art *dao.Article) *domain.Article {
	return &domain.Article{
		Id:      art.Id,
		Title:   art.Title,
		Content: art.Content,
		Author: domain.Author{
			Id: art.AuthorId,
		},
		CreatedAt: art.CreatedAt,
		UpdatedAt: art.UpdatedAt,
	}
}
