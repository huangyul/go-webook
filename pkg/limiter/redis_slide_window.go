package limiter

import (
	"context"
	_ "embed"
	"github.com/redis/go-redis/v9"
	"time"
)

//go:embed slide_window.lua
var luaScript string

type RedisSlideWindowRedis struct {
	cmd      redis.Cmdable
	interval time.Duration
	rate     int
}

func NewRedisSlideWindowRedis(cmd redis.Cmdable, interval time.Duration, rate int) *RedisSlideWindowRedis {
	return &RedisSlideWindowRedis{cmd: cmd, interval: interval, rate: rate}
}

func (r *RedisSlideWindowRedis) Limit(ctx context.Context, key string) (bool, error) {
	return r.cmd.Eval(ctx, luaScript, []string{key}, r.interval.Milliseconds(), r.rate, time.Now().UnixMilli()).Bool()
}
