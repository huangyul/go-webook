package cache

import (
	"context"
	_ "embed"
	"fmt"
	"github.com/redis/go-redis/v9"
)

//go:embed lua/incr_cnt.lua
var incrScript string

type InteractiveCache interface {
	IncrReadCntIfPresent(ctx context.Context, bizID int64, biz string) error
}

var _ InteractiveCache = (*RedisInteractiveCache)(nil)

type RedisInteractiveCache struct {
	client redis.Cmdable
}

func NewInteractiveCache(client redis.Cmdable) InteractiveCache {
	return &RedisInteractiveCache{client: client}
}

func (cache *RedisInteractiveCache) IncrReadCntIfPresent(ctx context.Context, bizID int64, biz string) error {
	return cache.client.Eval(ctx, incrScript, []string{cache.key(biz, bizID)}, 1).Err()
}

func (cache *RedisInteractiveCache) key(biz string, bizID int64) string {
	return fmt.Sprintf("interactive:%s:%d", biz, bizID)
}
