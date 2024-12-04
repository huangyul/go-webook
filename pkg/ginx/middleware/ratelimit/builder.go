package ratelimit

import (
	_ "embed"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/huangyul/go-blog/pkg/limiter"
	"net/http"
)

//go:embed slide_window.lua
var luaString string

type Builder struct {
	prefix string
	l      limiter.Limiter
}

func NewBuilder(l limiter.Limiter) *Builder {
	return &Builder{
		prefix: "ip-limiter",
		l:      l,
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
	key := fmt.Sprintf("%s_%s", b.prefix, ctx.ClientIP())
	return b.l.Limit(ctx, key)
}
