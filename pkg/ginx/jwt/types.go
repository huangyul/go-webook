package ginxjwt

import "github.com/gin-gonic/gin"

type JWT interface {
	GenToken(ctx *gin.Context, uid int64) (string, string, error)
	ExtractToken(ctx *gin.Context) (JwtClaims, error)
	ClearToken(ctx *gin.Context) error
	CheckToken(ctx *gin.Context) (JwtClaims, error)
	RefreshToken(ctx *gin.Context) (string, error)
}
