package cache

import (
	"context"
	_ "embed"
	"fmt"
	"github.com/redis/go-redis/v9"
)

//go:embed lua/incr_cnt.lua
var IncrCntScript string

const (
	ReadBiz    = "read_cnt"
	LikeBiz    = "like_cnt"
	CollectBiz = "collect_cnt"
)

type InteractiveCache interface {
	IncrReadCntIfPresent(ctx context.Context, biz string, bizId int64) error
}

func NewInteractiveCache(rds redis.Cmdable) InteractiveCache {
	return &RedisInteractiveCache{
		rds: rds,
	}
}

type RedisInteractiveCache struct {
	rds redis.Cmdable
}

func (cache *RedisInteractiveCache) IncrReadCntIfPresent(ctx context.Context, biz string, bizId int64) error {
	return cache.rds.Eval(ctx, IncrCntScript, []string{cache.key(biz, bizId)}, ReadBiz, 1).Err()
}

func (cache *RedisInteractiveCache) key(biz string, bizId int64) string {
	return fmt.Sprintf("interactive:%s:%d", biz, bizId)
}
