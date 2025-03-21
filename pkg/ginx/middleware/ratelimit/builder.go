package ratelimit

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/huangyul/go-webook/pkg/limiter"
	"net/http"
)

type Builder struct {
	prefix string
	limit  limiter.Limiter
}

func NewBuilder(prefix string, limit limiter.Limiter) *Builder {
	return &Builder{
		prefix: prefix,
		limit:  limit,
	}
}

func (b *Builder) Build() gin.HandlerFunc {
	return func(ctx *gin.Context) {

		isLimit, err := b.limit.Limit(ctx, fmt.Sprintf("%s:%s", b.prefix, ctx.ClientIP()))
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
