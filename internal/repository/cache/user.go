package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/huangyul/go-blog/internal/domain"
	"github.com/redis/go-redis/v9"
)

type UserCache interface {
	Get(ctx context.Context, uID int64) (domain.User, error)
	Set(ctx context.Context, user domain.User) error
}

var _ UserCache = (*RedisUserCache)(nil)

type RedisUserCache struct {
	client     redis.Cmdable
	expiration time.Duration
}

func NewRedisUserCache(client redis.Cmdable) UserCache {
	return &RedisUserCache{client: client, expiration: time.Minute * 15}
}

func (cache *RedisUserCache) Get(ctx context.Context, uID int64) (domain.User, error) {
	data, err := cache.client.Get(ctx, cache.Key(uID)).Result()
	if err != nil {
		return domain.User{}, err
	}
	var user domain.User
	err = json.Unmarshal([]byte(data), &user)
	if err != nil {
		return domain.User{}, err
	}
	return user, nil
}

func (cache *RedisUserCache) Set(ctx context.Context, user domain.User) error {
	data, err := json.Marshal(user)
	if err != nil {
		return err
	}
	return cache.client.Set(ctx, cache.Key(user.ID), data, cache.expiration).Err()
}

func (cache *RedisUserCache) Key(id int64) string {
	return fmt.Sprintf("go-blog:user:%d", id)
}
