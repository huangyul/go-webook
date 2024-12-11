package cache

import (
	"context"
	_ "embed"
	"fmt"
	"github.com/redis/go-redis/v9"
)

const (
	readKey    = "read_cnt"
	likeKey    = "like_cnt"
	collectKey = "collect_cnt"
)

//go:embed lua/incr_cnt.lua
var incrScript string

type InteractiveCache interface {
	IncrReadCntIfPresent(ctx context.Context, bizID int64, biz string) error
	IncrLikeCntIfPresent(ctx context.Context, bizID int64, biz string) error
	DecrLikeCntIfPresent(ctx context.Context, bizID int64, biz string) error
}

var _ InteractiveCache = (*RedisInteractiveCache)(nil)

type RedisInteractiveCache struct {
	client redis.Cmdable
}

func (cache *RedisInteractiveCache) DecrLikeCntIfPresent(ctx context.Context, bizID int64, biz string) error {
	return cache.client.Eval(ctx, incrScript, []string{cache.key(biz, bizID)}, likeKey, -1).Err()
}

func (cache *RedisInteractiveCache) IncrLikeCntIfPresent(ctx context.Context, bizID int64, biz string) error {
	return cache.client.Eval(ctx, incrScript, []string{cache.key(biz, bizID)}, likeKey, 1).Err()
}

func NewInteractiveCache(client redis.Cmdable) InteractiveCache {
	return &RedisInteractiveCache{client: client}
}

func (cache *RedisInteractiveCache) IncrReadCntIfPresent(ctx context.Context, bizID int64, biz string) error {
	return cache.client.Eval(ctx, incrScript, []string{cache.key(biz, bizID)}, readKey, 1).Err()
}

func (cache *RedisInteractiveCache) key(biz string, bizID int64) string {
	return fmt.Sprintf("interactive:%s:%d", biz, bizID)
}
