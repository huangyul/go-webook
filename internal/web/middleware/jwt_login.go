package middleware

import (
	"github.com/gin-gonic/gin"
	jwt "github.com/golang-jwt/jwt/v5"
	"net/http"
	"strings"
)

type Option func(*JWTLoginMiddlewareBuild)

type JWTLoginMiddlewareBuild struct {
	WhiteList []string
}

func AddWhiteList(whiteList ...string) Option {
	return func(b *JWTLoginMiddlewareBuild) {
		b.WhiteList = whiteList
	}
}

func NewJWTLoginMiddlewareBuild(opts ...Option) *JWTLoginMiddlewareBuild {
	m := &JWTLoginMiddlewareBuild{}
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
		var c LoginClaims
		token, err := jwt.ParseWithClaims(tokenStr, &c, func(token *jwt.Token) (interface{}, error) {
			return []byte("secret"), nil
		})
		if err != nil || !token.Valid {
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
