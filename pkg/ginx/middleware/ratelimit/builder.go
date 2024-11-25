package ratelimit

import (
	_ "embed"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"net/http"
	"time"
)

//go:embed slide_window.lua
var luaString string

type Builder struct {
	prefix   string
	cmd      redis.Cmdable
	interval time.Duration
	rate     int
}

func NewBuilder(cmd redis.Cmdable, rate int, interval time.Duration) *Builder {
	return &Builder{
		prefix:   "ip-limiter",
		cmd:      cmd,
		rate:     rate,
		interval: interval,
	}
}

func (b *Builder) Build() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		limited, err := b.limit(ctx)
		if err != nil {
			ctx.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		if limited {
			ctx.AbortWithStatus(http.StatusTooManyRequests)
			return
		}
		ctx.Next()
	}
}

func (b *Builder) limit(ctx *gin.Context) (bool, error) {
	key := fmt.Sprintf("%s:%s", b.prefix, ctx.ClientIP())
	return b.cmd.Eval(ctx, luaString, []string{key}, b.interval.Milliseconds(), b.rate, time.Now().UnixMilli()).Bool()
}
