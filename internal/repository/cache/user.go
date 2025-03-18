package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/huangyul/go-webook/internal/domain"
	"github.com/redis/go-redis/v9"
)

type UserCache interface {
	Get(ctx context.Context, id int64) (*domain.User, error)
	Set(ctx context.Context, user *domain.User) error
	Delete(ctx context.Context, id int64) error
}

var _ UserCache = (*RedisUserCache)(nil)

type RedisUserCache struct {
	cmd redis.Cmdable
}

func (r *RedisUserCache) Delete(ctx context.Context, id int64) error {
	return r.cmd.Del(ctx, r.key(id)).Err()
}

func (r *RedisUserCache) Get(ctx context.Context, id int64) (*domain.User, error) {
	var user domain.User
	key := r.key(id)
	res, err := r.cmd.Get(ctx, key).Result()
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal([]byte(res), &user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *RedisUserCache) Set(ctx context.Context, user *domain.User) error {
	key := r.key(user.ID)
	data, err := json.Marshal(user)
	if err != nil {
		return err
	}
	return r.cmd.Set(ctx, key, data, 0).Err()
}

func (r *RedisUserCache) key(id int64) string {
	return fmt.Sprintf("user:%d", id)
}

func NewRedisUserCache(cmd redis.Cmdable) UserCache {
	return &RedisUserCache{
		cmd: cmd,
	}
}
