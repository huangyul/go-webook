package cache

import (
	"context"
	_ "embed"
	"errors"
	"fmt"
	"github.com/huangyul/go-blog/internal/domain"
	"github.com/redis/go-redis/v9"
	"strconv"
	"time"
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
	IncrCollectCntIfPresent(ctx context.Context, bizID int64, id int64, biz string) error
	Get(ctx context.Context, id int64, biz string) (domain.Interactive, error)
	Set(ctx context.Context, id int64, biz string, int domain.Interactive) error
}

var _ InteractiveCache = (*RedisInteractiveCache)(nil)

type RedisInteractiveCache struct {
	client redis.Cmdable
}

func (cache *RedisInteractiveCache) Get(ctx context.Context, id int64, biz string) (domain.Interactive, error) {
	res, err := cache.client.HGetAll(ctx, cache.key(biz, id)).Result()
	if err != nil {
		return domain.Interactive{}, err
	}
	if len(res) == 0 {
		return domain.Interactive{}, errors.New("key not found")
	}
	var int domain.Interactive
	int.CollectCnt, _ = strconv.Atoi(res[collectKey])
	int.ReadCnt, _ = strconv.Atoi(res[readKey])
	int.LikeCnt, _ = strconv.Atoi(res[likeKey])
	return int, nil
}

func (cache *RedisInteractiveCache) Set(ctx context.Context, id int64, biz string, int domain.Interactive) error {
	err := cache.client.HSet(ctx, cache.key(biz, id), readKey, int.ReadCnt, likeKey, int.LikeCnt, collectKey, int.CollectCnt).Err()
	if err != nil {
		return err
	}
	return cache.client.Expire(ctx, cache.key(biz, id), time.Minute*15).Err()
}

func NewInteractiveCache(client redis.Cmdable) InteractiveCache {
	return &RedisInteractiveCache{client: client}
}

func (cache *RedisInteractiveCache) DecrLikeCntIfPresent(ctx context.Context, bizID int64, biz string) error {
	return cache.client.Eval(ctx, incrScript, []string{cache.key(biz, bizID)}, likeKey, -1).Err()
}

func (cache *RedisInteractiveCache) IncrCollectCntIfPresent(ctx context.Context, bizID int64, id int64, biz string) error {
	return cache.client.Eval(ctx, incrScript, []string{cache.key(biz, bizID)}, collectKey, 1).Err()
}

func (cache *RedisInteractiveCache) IncrLikeCntIfPresent(ctx context.Context, bizID int64, biz string) error {
	return cache.client.Eval(ctx, incrScript, []string{cache.key(biz, bizID)}, likeKey, 1).Err()
}

func (cache *RedisInteractiveCache) IncrReadCntIfPresent(ctx context.Context, bizID int64, biz string) error {
	return cache.client.Eval(ctx, incrScript, []string{cache.key(biz, bizID)}, readKey, 1).Err()
}

func (cache *RedisInteractiveCache) key(biz string, bizID int64) string {
	return fmt.Sprintf("interactive:%s:%d", biz, bizID)
}
