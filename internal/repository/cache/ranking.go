package cache

import (
	"context"
	"encoding/json"
	"github.com/huangyul/go-webook/internal/domain"
	"github.com/redis/go-redis/v9"
	"time"
)

type RankingCache interface {
	Get(ctx context.Context) ([]domain.Article, error)
	Set(ctx context.Context, items []domain.Article) error
}

func NewRankingCache(client redis.Cmdable) RankingCache {
	return &RedisRankingCache{
		client:     client,
		key:        "ranking:top_n",
		expiration: time.Minute * 3,
	}
}

type RedisRankingCache struct {
	client     redis.Cmdable
	key        string
	expiration time.Duration
}

func (r *RedisRankingCache) Get(ctx context.Context) ([]domain.Article, error) {
	data, err := r.client.Get(ctx, "ranking").Bytes()
	if err != nil {
		return nil, err
	}
	var result []domain.Article
	err = json.Unmarshal(data, &result)
	return result, err
}

func (r *RedisRankingCache) Set(ctx context.Context, items []domain.Article) error {
	for _, item := range items {
		item.Content = item.Abstract()
	}
	data, err := json.Marshal(items)
	if err != nil {
		return err
	}
	return r.client.Set(ctx, "ranking", data, 0).Err()
}
