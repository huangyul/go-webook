package cache

import (
	"context"
	_ "embed"
	"encoding/json"
	"fmt"
	"time"

	"github.com/huangyul/go-webook/interactive/domain"
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
	IncrLikeCntIfPresent(ctx context.Context, biz string, bizId int64) error
	DecrLikeCntIfPresent(ctx context.Context, biz string, bizId int64) error
	IncrCollectCntIfPresent(ctx context.Context, biz string, bizId int64) error
	DecrCollectCntIfPresent(ctx context.Context, biz string, bizId int64) error
	Set(ctx context.Context, biz string, bizId int64, inter *domain.Interactive) error
	Get(ctx context.Context, biz string, bizId int64) (*domain.Interactive, error)
}

func NewInteractiveCache(rds redis.Cmdable) InteractiveCache {
	return &RedisInteractiveCache{
		rds: rds,
	}
}

type RedisInteractiveCache struct {
	rds redis.Cmdable
}

func (cache *RedisInteractiveCache) Set(ctx context.Context, biz string, bizId int64, inter *domain.Interactive) error {
	key := cache.key(biz, bizId)
	err := cache.rds.HSet(ctx, key, ReadBiz, inter.ReadCnt, LikeBiz, inter.LikeCnt, CollectBiz, inter.CollectCnt).Err()
	if err != nil {
		return err
	}
	return cache.rds.Expire(ctx, key, time.Minute*15).Err()
}

func (cache *RedisInteractiveCache) Get(ctx context.Context, biz string, bizId int64) (*domain.Interactive, error) {
	key := cache.key(biz, bizId)
	val, err := cache.rds.HGet(ctx, key, ReadBiz).Result()
	if err != nil {
		return nil, err
	}
	var inter *domain.Interactive
	err = json.Unmarshal([]byte(val), &inter)
	if err != nil {
		return nil, err
	}
	return inter, nil
}

func (cache *RedisInteractiveCache) IncrCollectCntIfPresent(ctx context.Context, biz string, bizId int64) error {
	return cache.rds.Eval(ctx, IncrCntScript, []string{cache.key(biz, bizId)}, CollectBiz, 1).Err()
}

func (cache *RedisInteractiveCache) DecrCollectCntIfPresent(ctx context.Context, biz string, bizId int64) error {
	return cache.rds.Eval(ctx, IncrCntScript, []string{cache.key(biz, bizId)}, CollectBiz, -1).Err()
}

func (cache *RedisInteractiveCache) IncrLikeCntIfPresent(ctx context.Context, biz string, bizId int64) error {
	return cache.rds.Eval(ctx, IncrCntScript, []string{cache.key(biz, bizId)}, ReadBiz, 1).Err()
}

func (cache *RedisInteractiveCache) DecrLikeCntIfPresent(ctx context.Context, biz string, bizId int64) error {
	return cache.rds.Eval(ctx, IncrCntScript, []string{cache.key(biz, bizId)}, LikeBiz, -1).Err()
}

func (cache *RedisInteractiveCache) IncrReadCntIfPresent(ctx context.Context, biz string, bizId int64) error {
	return cache.rds.Eval(ctx, IncrCntScript, []string{cache.key(biz, bizId)}, LikeBiz, 1).Err()
}

func (cache *RedisInteractiveCache) key(biz string, bizId int64) string {
	return fmt.Sprintf("interactive:%s:%d", biz, bizId)
}
