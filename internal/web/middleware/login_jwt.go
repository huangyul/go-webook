package middleware

import (
	"encoding/gob"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/huangyul/go-blog/internal/web"
	"log"
	"net/http"
	"strings"
	"time"
)

var (
	JWT_TOKEN_KEY = []byte("JWT_TOKEN_KEY")
)

type JWTLoginMiddlewareBuild struct {
	whiteList []string
}

func NewJWTLoginMiddlewareBuild() *JWTLoginMiddlewareBuild {
	return &JWTLoginMiddlewareBuild{}
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

		tokenStr := ctx.GetHeader("Authorization")
		if tokenStr == "" {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		tokenStr = strings.Split(tokenStr, " ")[1]

		if tokenStr == "" {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		var c web.JWTClaims

		token, err := jwt.ParseWithClaims(tokenStr, &c, func(token *jwt.Token) (interface{}, error) {
			return JWT_TOKEN_KEY, nil
		})
		if token == nil || !token.Valid || err != nil {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		if c.UserAgent != ctx.Request.UserAgent() {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		ctx.Set("user_id", c.UserID)
		// renewal of token validity time
		expiredAt := c.ExpiresAt
		fmt.Println(time.Now().Sub(expiredAt.Time))
		if expiredAt.Sub(time.Now()) < time.Second {
			newC := web.JWTClaims{
				UserID:    c.UserID,
				UserAgent: c.UserAgent,
				RegisteredClaims: jwt.RegisteredClaims{
					ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 2)),
				},
			}
			token = jwt.NewWithClaims(jwt.SigningMethodHS256, newC)
			tokenStr, err = token.SignedString([]byte(JWT_TOKEN_KEY))
			if err != nil {
				log.Println(err)
			}
			ctx.Header("x-jwt-token", tokenStr)
		}

		ctx.Next()

	}
}
