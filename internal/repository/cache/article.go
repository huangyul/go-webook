package cache

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/huangyul/go-webook/internal/domain"
	"github.com/redis/go-redis/v9"
)

type ArticleCache interface {
	GetFirstPage(ctx context.Context, userId int64) ([]*domain.Article, error)
	SetFirstPage(ctx context.Context, userId int64, article []*domain.Article) error
	DeleteFirstPage(ctx context.Context, userId int64) error
	GetPubById(ctx context.Context, id int64, userId int64) (*domain.Article, error)
	SetPubById(ctx context.Context, id int64, article *domain.Article) error
	GetById(ctx context.Context, id int64, userId int64) (*domain.Article, error)
	SetById(ctx context.Context, id int64, article *domain.Article) error
}

type RedisArticleCache struct {
	client redis.Cmdable
	expire time.Duration
}

// GetById
func (r *RedisArticleCache) GetById(ctx context.Context, id int64, userId int64) (*domain.Article, error) {
	var art domain.Article
	data, err := r.client.Get(ctx, r.keyArt(id)).Bytes()
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(data, &art)
	if err != nil {
		return nil, err
	}
	if art.Author.Id != userId {
		return nil, errors.New("article not found")
	}
	return &art, nil
}

// GetPubById implements ArticleCache.
func (r *RedisArticleCache) GetPubById(ctx context.Context, id int64, userId int64) (*domain.Article, error) {
	var art domain.Article
	data, err := r.client.Get(ctx, r.keyPubArt(id)).Bytes()
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(data, &art)
	if err != nil {
		return nil, err
	}
	if art.Author.Id != userId {
		return nil, errors.New("pub article not found")
	}
	return &art, nil
}

// SetById implements ArticleCache.
func (r *RedisArticleCache) SetById(ctx context.Context, id int64, article *domain.Article) error {
	data, err := json.Marshal(article)
	if err != nil {
		return err
	}
	return r.client.Set(ctx, r.keyArt(id), data, r.expire).Err()
}

// SetPubById implements ArticleCache.
func (r *RedisArticleCache) SetPubById(ctx context.Context, id int64, article *domain.Article) error {
	data, err := json.Marshal(article)
	if err != nil {
		return err
	}
	return r.client.Set(ctx, r.keyPubArt(id), data, r.expire).Err()
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

func (r *RedisArticleCache) keyPubArt(id int64) string {
	return fmt.Sprintf("article:pub_article:%d", id)
}

func (r *RedisArticleCache) keyArt(id int64) string {
	return fmt.Sprintf("article:article:%d", id)
}

func NewArticleCache(client redis.Cmdable) ArticleCache {
	return &RedisArticleCache{
		client: client,
		expire: time.Hour,
	}
}
