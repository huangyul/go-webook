package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/huangyul/go-webook/internal/domain"
	"github.com/redis/go-redis/v9"
	"time"
)

type ArticleCache interface {
	GetFirstPage(ctx context.Context, userId int64) ([]*domain.Article, error)
	SetFirstPage(ctx context.Context, userId int64, article []*domain.Article) error
	DeleteFirstPage(ctx context.Context, userId int64) error
}

type RedisArticleCache struct {
	client redis.Cmdable
	expire time.Duration
}

func (r *RedisArticleCache) GetFirstPage(ctx context.Context, userId int64) ([]*domain.Article, error) {
	data, err := r.client.Get(ctx, r.key(userId)).Bytes()
	if err != nil {
		return nil, err
	}
	var articles []*domain.Article
	err = json.Unmarshal(data, &articles)
	if err != nil {
		return nil, err
	}
	return articles, nil
}

func (r *RedisArticleCache) SetFirstPage(ctx context.Context, userId int64, arts []*domain.Article) error {
	data, err := json.Marshal(arts)
	if err != nil {
		return err
	}
	return r.client.Set(ctx, r.key(userId), data, r.expire).Err()
}

func (r *RedisArticleCache) DeleteFirstPage(ctx context.Context, userId int64) error {
	return r.client.Del(ctx, r.key(userId)).Err()
}

func (r *RedisArticleCache) key(userId int64) string {
	return fmt.Sprintf("article:first_page:%d", userId)
}

func NewArticleCache(client redis.Cmdable) ArticleCache {
	return &RedisArticleCache{
		client: client,
		expire: time.Hour,
	}
}
