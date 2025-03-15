package middleware

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"net/http"
)

type LoginBuilder struct {
}

func NewLoginBuilder() *LoginBuilder {
	return &LoginBuilder{}
}

func (b *LoginBuilder) Build() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		path := ctx.Request.URL.Path
		if path == "/user/register" || path == "/user/login" {
			return
		}

		session := sessions.Default(ctx)
		if session.Get("userId") == nil {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		ctx.Next()
	}
}
