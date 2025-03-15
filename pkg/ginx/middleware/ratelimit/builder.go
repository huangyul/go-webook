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
var limitScript string

type Option func(*Builder)

type Builder struct {
	prefix   string
	cmd      redis.Cmdable
	rate     int
	interval time.Duration
}

func NewBuilder(cmd redis.Cmdable, opts ...Option) *Builder {
	b := &Builder{
		cmd:      cmd,
		prefix:   "ratelimit",
		rate:     10,
		interval: time.Minute,
	}

	for _, opt := range opts {
		opt(b)
	}

	return b
}

func SetRate(rate int) Option {
	return func(b *Builder) {
		b.rate = rate
	}
}

func SetInterval(interval time.Duration) Option {
	return func(b *Builder) {
		b.interval = interval
	}
}

func (b *Builder) Build() gin.HandlerFunc {
	return func(ctx *gin.Context) {

		isLimit, err := b.limit(ctx)
		if err != nil {
			// redis crash
			ctx.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		if isLimit {
			ctx.AbortWithStatus(http.StatusTooManyRequests)
			return
		}
		ctx.Next()
	}
}

// limit rate limit based on ip
func (b *Builder) limit(ctx *gin.Context) (bool, error) {
	key := fmt.Sprintf("%s:limit:%s", b.prefix, ctx.ClientIP())
	return b.cmd.Eval(ctx, limitScript, []string{key}, b.rate, b.interval.Milliseconds(), time.Now().UnixMilli()).Bool()
}
