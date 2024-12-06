package middleware

import (
	"encoding/gob"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	ginxjwt "github.com/huangyul/go-blog/pkg/ginx/jwt"
)

type JWTLoginMiddlewareBuild struct {
	whiteList []string
	jwt       ginxjwt.JWT
}

func NewJWTLoginMiddlewareBuild(jwt ginxjwt.JWT) *JWTLoginMiddlewareBuild {
	return &JWTLoginMiddlewareBuild{
		jwt: jwt,
	}
}

func (b *JWTLoginMiddlewareBuild) AddWhiteList(whiteList ...string) *JWTLoginMiddlewareBuild {
	b.whiteList = append(b.whiteList, whiteList...)
	return b
}

func (b *JWTLoginMiddlewareBuild) Build() gin.HandlerFunc {
	gob.Register(time.Time{})
	return func(ctx *gin.Context) {
		for _, r := range b.whiteList {
			if r == ctx.Request.URL.Path {
				ctx.Next()
				return
			}
		}

		c, err := b.jwt.CheckToken(ctx)

		if err != nil {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		ctx.Set("user_id", c.Uid)

		ctx.Next()

	}
}
