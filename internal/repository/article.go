package repository

import (
	"context"
	"github.com/huangyul/go-blog/internal/domain"
	"github.com/huangyul/go-blog/internal/repository/cache"
	"github.com/huangyul/go-blog/internal/repository/dao"
	"time"
)

//go:generate mockgen -source=./article.cgo -package=repomocks -destination=./mocks/article_repo_mock.go ArticleRepository
type ArticleRepository interface {
	Create(ctx context.Context, art domain.Article) (int64, error)
	Update(ctx context.Context, art domain.Article) error
	Sync(ctx context.Context, art domain.Article) (int64, error)
	SyncStatus(ctx context.Context, uid int64, id int64, status domain.ArticleStatus) error
	ListByAuthor(ctx context.Context, uid int64, page int64, pageSize int64) ([]domain.Article, error)
	GetById(ctx context.Context, uid int64, id int64) (domain.Article, error)
	GetPubById(ctx context.Context, uid int64, id int64) (domain.Article, error)
}

var _ ArticleRepository = (*articleRepository)(nil)

type articleRepository struct {
	dao   dao.ArticleDao
	cache cache.ArticleCache
}

func NewArticleRepository(dao dao.ArticleDao, cache cache.ArticleCache) ArticleRepository {
	return &articleRepository{dao: dao, cache: cache}
}

func (repo *articleRepository) GetPubById(ctx context.Context, uid int64, id int64) (domain.Article, error) {
	dart, err := repo.cache.GetPubDetail(ctx, uid, id)
	if err == nil {
		return dart, err
	}
	art, err := repo.dao.GetPubById(ctx, uid, id)
	if err != nil {
		return domain.Article{}, err
	}
	dart = repo.toDomain(dao.Article(art))

	go func() {
		ct, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()
		_ = repo.cache.SetPubDetail(ct, uid, id, dart)
	}()

	return dart, nil
}

func (repo *articleRepository) GetById(ctx context.Context, uid int64, id int64) (domain.Article, error) {
	dart, err := repo.cache.GetDetail(ctx, uid, id)
	if err == nil {
		return dart, nil
	}
	art, err := repo.dao.GetById(ctx, uid, id)
	if err != nil {
		return domain.Article{}, err
	}
	dart = repo.toDomain(art)
	go func() {
		ct, cancel := context.WithTimeout(context.Background(), time.Second*2)
		defer cancel()
		repo.cache.SetDetail(ct, uid, id, dart)

	}()
	return dart, nil
}

func (repo *articleRepository) SyncStatus(ctx context.Context, uid int64, id int64, status domain.ArticleStatus) error {
	return repo.dao.SyncStatus(ctx, uid, id, uint8(status))
}

func (repo *articleRepository) Create(ctx context.Context, art domain.Article) (int64, error) {
	return repo.dao.Insert(ctx, repo.toDao(art))
}

func (repo *articleRepository) Update(ctx context.Context, art domain.Article) error {
	err := repo.dao.UpdateByID(ctx, repo.toDao(art))
	// TODO 使用kafka优化
	go repo.cache.SetDetail(ctx, art.Author.ID, art.ID, art)
	return err
}

func (repo *articleRepository) Sync(ctx context.Context, art domain.Article) (int64, error) {
	return repo.dao.Sync(ctx, repo.toDao(art))
}

func (repo *articleRepository) ListByAuthor(ctx context.Context, uid int64, page int64, pageSize int64) ([]domain.Article, error) {
	arts, err := repo.dao.ListByAuthor(ctx, uid, page, pageSize)
	if err != nil {
		return nil, err
	}
	res := make([]domain.Article, 0, len(arts))
	for _, art := range arts {
		res = append(res, repo.toDomain(art))
	}
	return res, nil
}

func (repo *articleRepository) toDomain(art dao.Article) domain.Article {
	return domain.Article{
		ID:        art.ID,
		Title:     art.Title,
		Content:   art.Content,
		Status:    art.Status,
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
		Status:    art.Status,
		CreatedAt: art.CreatedAt.UnixMilli(),
		UpdatedAt: art.UpdatedAt.UnixMilli(),
	}
}
