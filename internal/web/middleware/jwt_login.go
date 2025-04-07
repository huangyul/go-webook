package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/huangyul/go-webook/internal/pkg/authz"
)

type Option func(*JWTLoginMiddlewareBuild)

type JWTLoginMiddlewareBuild struct {
	WhiteList []string
	jwt       authz.Authz
}

func AddWhiteList(whiteList ...string) Option {
	return func(b *JWTLoginMiddlewareBuild) {
		b.WhiteList = whiteList
	}
}

func NewJWTLoginMiddlewareBuild(jwt authz.Authz, opts ...Option) *JWTLoginMiddlewareBuild {
	m := &JWTLoginMiddlewareBuild{
		jwt: jwt,
	}
	for _, opt := range opts {
		opt(m)
	}
	return m
}

func (m *JWTLoginMiddlewareBuild) Build() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		path := ctx.Request.URL.Path
		for _, whitePath := range m.WhiteList {
			if whitePath == path {
				ctx.Next()
				return
			}
		}

		tokenStr := ctx.GetHeader("Authorization")

		if tokenStr == "" {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		tokenStr = strings.TrimPrefix(tokenStr, "Bearer ")
		if tokenStr == "" {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		c, err := m.jwt.VerifyToken(tokenStr)
		if err != nil {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		ok, err := m.jwt.CheckToken(tokenStr)
		if err != nil || !ok {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		// can add a check here: if the token is about to expire, renew it
		// and set it to the header
		//if c.ExpiresAt.Sub(time.Now()) < time.Minute*5 {
		//	c.ExpiresAt = jwt.NewNumericDate(time.Now().Add(time.Minute * 30))
		//	tokenStr, err = token.SignedString([]byte("secret"))
		//	if err!= nil {
		//		ctx.AbortWithStatus(http.StatusUnauthorized)
		//		return
		//	}
		//	ctx.Header("x-jwt-token", tokenStr)

		ctx.Set("user_id", c.UserId)

		ctx.Next()
	}
}

type LoginClaims struct {
	UserId int64
	jwt.RegisteredClaims
}
