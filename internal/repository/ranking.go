package repository

import (
	"context"
	"github.com/huangyul/go-webook/internal/domain"
	"github.com/huangyul/go-webook/internal/repository/cache"
)

type RankingRepository interface {
	GetTopN(ctx context.Context) ([]domain.Article, error)
	ReplaceTopN(ctx context.Context, arts []domain.Article) error
}

func NewRankingRepository(cache cache.RankingCache) RankingRepository {
	return &RankingRepositoryImpl{
		cache: cache,
	}
}

type RankingRepositoryImpl struct {
	cache cache.RankingCache
}

func (r *RankingRepositoryImpl) GetTopN(ctx context.Context) ([]domain.Article, error) {
	return r.cache.Get(ctx)
}
func (r *RankingRepositoryImpl) ReplaceTopN(ctx context.Context, arts []domain.Article) error {
	return r.cache.Set(ctx, arts)
}
