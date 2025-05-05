package repository

import (
	"context"
	"time"

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
	GetById(ctx context.Context, id int64, userId int64) (*domain.Article, error)
	GetPubById(ctx context.Context, id int64, userId int64) (*domain.Article, error)
	GetPudDetailById(ctx context.Context, id int64) (*domain.Article, error)
	ListPub(ctx context.Context, start time.Time, offset, size int) ([]*domain.Article, error)
}

type articleRepository struct {
	dao   dao.ArticleDAO
	cache cache.ArticleCache
}

// ListPub implements ArticleRepository.
func (a *articleRepository) ListPub(ctx context.Context, start time.Time, offset int, size int) ([]*domain.Article, error) {
	res, err := a.dao.ListPub(ctx, start, offset, size)
	if err != nil {
		return nil, err
	}
	arts := make([]*domain.Article, 0, len(res))
	for _, r := range res {
		arts = append(arts, a.toDomain(r))
	}
	return arts, nil
}

func (a *articleRepository) GetPudDetailById(ctx context.Context, id int64) (*domain.Article, error) {
	res, err := a.dao.GetPudDetailById(ctx, id)
	if err != nil {
		return nil, err
	}
	return a.toDomain(res), nil
}

// GetById
func (a *articleRepository) GetById(ctx context.Context, id int64, userId int64) (*domain.Article, error) {
	dArt, err := a.cache.GetById(ctx, id, userId)
	if err == nil {
		return dArt, nil
	}
	art, err := a.dao.GetById(ctx, id, userId)
	if err != nil {
		return nil, err
	}
	go func() {
		_ = a.cache.SetById(ctx, id, a.toDomain(art))
	}()
	return a.toDomain(art), nil
}

// GetPubById
func (a *articleRepository) GetPubById(ctx context.Context, id int64, userId int64) (*domain.Article, error) {
	dArt, err := a.cache.GetPubById(ctx, id, userId)
	if err == nil {
		return dArt, nil
	}
	art, err := a.dao.GetPubById(ctx, id, userId)
	if err != nil {
		return nil, err
	}
	go func() {
		_ = a.cache.SetPubById(ctx, id, a.toDomain(art))
	}()
	return a.toDomain(art), nil
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
	art, err := a.dao.SyncStatus(ctx, userId, id, status.ToUint8())
	if err != nil {
		return err
	}
	go func() {
		if status == domain.ArticleStatusPublished {
			_ = a.cache.SetPubById(ctx, id, a.toDomain(art))
		}
	}()
	return nil
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
		Status:   art.Status.ToUint8(),
		Content:  art.Content,
		AuthorId: art.Author.Id,
	}
}

func (a *articleRepository) toDomain(art *dao.Article) *domain.Article {
	return &domain.Article{
		Id:      art.Id,
		Title:   art.Title,
		Content: art.Content,
		Status:  domain.ArticleStatus(art.Status),
		Author: domain.Author{
			Id: art.AuthorId,
		},
		CreatedAt: art.CreatedAt,
		UpdatedAt: art.UpdatedAt,
	}
}
