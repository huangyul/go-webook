package repository

import (
	"context"
	"github.com/huangyul/go-webook/internal/domain"
	"github.com/huangyul/go-webook/internal/repository/cache"
	"github.com/huangyul/go-webook/internal/repository/dao"
)

type ArticleRepository interface {
	Insert(ctx context.Context, art *domain.Article) (int64, error)
	UpdateById(ctx context.Context, art *domain.Article) error
	Sync(ctx context.Context, art *domain.Article) error
	SyncStatus(ctx context.Context, userId, id int64, status domain.ArticleStatus) error
	GetByAuthorId(ctx context.Context, userId, page, pageSize int64) ([]*domain.Article, error)
}

type articleRepository struct {
	dao   dao.ArticleDAO
	cache cache.ArticleCache
}

func (a *articleRepository) GetByAuthorId(ctx context.Context, userId, page, pageSize int64) ([]*domain.Article, error) {
	limit := (page - 1) * pageSize
	if page == 1 && limit <= 100 {
		arts, err := a.cache.GetFirstPage(ctx, userId)
		if err == nil {
			return arts, nil
		}
	}
	arts, err := a.dao.GetByAuthorId(ctx, userId, page, pageSize)
	if err != nil {
		return nil, err
	}
	var res []*domain.Article
	for _, art := range arts {
		res = append(res, &domain.Article{
			Id:        art.Id,
			Title:     art.Title,
			Content:   art.Content,
			Status:    domain.ArticleStatus(art.Status),
			CreatedAt: art.CreatedAt,
			UpdatedAt: art.UpdatedAt,
			Author: domain.Author{
				Id: userId,
			},
		})
	}
	go func() {
		if page == 1 && limit <= 100 {
			_ = a.cache.SetFirstPage(ctx, userId, res)
		}
	}()
	return res, nil
}

func (a *articleRepository) SyncStatus(ctx context.Context, userId, id int64, status domain.ArticleStatus) error {
	return a.dao.SyncStatus(ctx, userId, id, status.ToUint8())
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

func NewArticleRepository(dao dao.ArticleDAO, cache cache.ArticleCache) ArticleRepository {
	return &articleRepository{
		dao:   dao,
		cache: cache,
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
