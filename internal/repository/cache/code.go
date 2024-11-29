package cache

import (
	"context"
	_ "embed"
	"fmt"

	"github.com/huangyul/go-blog/internal/pkg/errno"
	"github.com/redis/go-redis/v9"
)

//go:embed lua/set_code.lua
var SetCodeScript string

//go:embed lua/verify_code.lua
var VerifyCodeScript string

type CodeCache interface {
	Set(ctx context.Context, biz, phone, code string) error
	Verify(ctx context.Context, biz, phone, code string) (bool, error)
}

var _ CodeCache = (*RedisCodeCache)(nil)

type RedisCodeCache struct {
	cmd redis.Cmdable
}

func NewRedisCodeCache(cmd redis.Cmdable) CodeCache {
	return &RedisCodeCache{cmd: cmd}
}

func (cache *RedisCodeCache) Set(ctx context.Context, biz string, phone string, code string) error {
	resNum, err := cache.cmd.Eval(ctx, SetCodeScript, []string{cache.Key(biz, phone)}, code).Int()
	if err != nil {
		return err
	}
	switch resNum {
	case -2:
		return errno.ErrCodeNotExist
	case -1:
		return errno.ErrCodeSendTooFrequent
	default:
		return nil
	}
}

func (cache *RedisCodeCache) Verify(ctx context.Context, biz string, phone string, code string) (bool, error) {
	resNum, err := cache.cmd.Eval(ctx, VerifyCodeScript, []string{cache.Key(biz, phone)}, code).Int()
	if err != nil {
		return false, err
	}
	switch resNum {
	case -2:
		return false, errno.ErrCodeVerifyFailed
	case -1:
		return false, errno.ErrCodeNotExist
	default:
		return true, nil
	}
}

func (cache *RedisCodeCache) Key(biz string, phone string) string {
	return fmt.Sprintf("phone_code:%s:%s", biz, phone)
}
