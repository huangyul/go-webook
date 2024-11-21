package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

type LoginMiddleBuilder struct{}

func (m LoginMiddleBuilder) Build() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if strings.Contains(ctx.Request.URL.Path, "/login") || strings.Contains(ctx.Request.URL.Path, "/signup") {
			ctx.Next()
			return
		}
		sess := sessions.Default(ctx)
		userID := sess.Get("user_id")
		if userID == nil {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		ctx.Next()
	}
}
