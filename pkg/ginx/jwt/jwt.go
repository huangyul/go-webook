package ginxjwt

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

var (
	ErrSessionInvalid  = errors.New("session invalid")
	ErrNoAuthorization = errors.New("no authorization")
	ErrTokenInvalid    = errors.New("token invalid")
)

var (
	JWTKey   = []byte("k6CswdUm77WKcbM68UQUuxVsHSpTCwgK")
	RfJWTKey = []byte("k6CswdUm77WKcbM68UQUuxVsHSpTCwgA")
)

var _ JWT = (*Handler)(nil)

type Handler struct {
	client     redis.Cmdable
	jwtMethod  jwt.SigningMethod
	rfDuration time.Duration
}

func NewJWT(client redis.Cmdable) JWT {
	return &Handler{
		client:     client,
		jwtMethod:  jwt.SigningMethodHS512,
		rfDuration: time.Hour * 24 * 7,
	}
}

type JwtClaims struct {
	Uid  int64
	ssid string // the id stored in redis
	jwt.RegisteredClaims
}

func (h *Handler) RefreshToken(ctx *gin.Context) (string, error) {
	var c JwtClaims
	authHeader := ctx.Request.Header.Get("Authorization")
	if authHeader == "" {
		return "", ErrNoAuthorization
	}
	authStrs := strings.Split(authHeader, " ")
	if len(authStrs) != 2 || authStrs[0] != "Bearer" {
		return "", ErrNoAuthorization
	}
	token, err := jwt.ParseWithClaims(authStrs[1], &c, func(token *jwt.Token) (interface{}, error) {
		return RfJWTKey, nil
	})
	if err != nil || !token.Valid {
		return "", ErrTokenInvalid
	}
	cnt, err := h.client.Exists(ctx, h.genKey(c.ssid)).Result()
	if err != nil {
		return "", err
	}
	if cnt > 0 {
		return "", ErrSessionInvalid
	}
	jc := JwtClaims{
		Uid:  c.Uid,
		ssid: c.ssid,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 2)),
		},
	}
	tokenStr, err := jwt.NewWithClaims(h.jwtMethod, jc).SignedString(JWTKey)
	if err != nil {
		return "", err
	}
	return tokenStr, nil
}

// GenToken will generate jwtToken and refreshToken
func (h *Handler) GenToken(ctx *gin.Context, uid int64) (string, string, error) {
	ssid := h.ssid()
	tokenClaims := &JwtClaims{
		Uid:  uid,
		ssid: ssid,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 2)),
		},
	}
	tokenStr, err := jwt.NewWithClaims(h.jwtMethod, tokenClaims).SignedString(JWTKey)
	if err != nil {
		return "", "", err
	}
	refreshTokenClaims := &JwtClaims{
		ssid: ssid,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24 * 7)),
		},
	}
	refreshTokenStr, err := jwt.NewWithClaims(h.jwtMethod, refreshTokenClaims).SignedString(RfJWTKey)
	if err != nil {
		return "", "", err
	}
	return tokenStr, refreshTokenStr, nil
}

// ExtractToken extracts the token from auth header
func (h *Handler) ExtractToken(ctx *gin.Context) (JwtClaims, error) {
	var c JwtClaims
	authHeader := ctx.Request.Header.Get("Authorization")
	if authHeader == "" {
		return c, ErrNoAuthorization
	}
	authStrs := strings.Split(authHeader, " ")
	if len(authStrs) != 2 || authStrs[0] != "Bearer" {
		return c, ErrNoAuthorization
	}
	token, err := jwt.ParseWithClaims(authStrs[1], &c, func(token *jwt.Token) (interface{}, error) {
		return JWTKey, nil
	})
	if err != nil || !token.Valid {
		return c, ErrTokenInvalid
	}
	return c, nil
}

// ClearToken add ssid to redis to mark token as logged out
func (h *Handler) ClearToken(ctx *gin.Context) error {
	claims, err := h.ExtractToken(ctx)
	if err != nil {
		return err
	}
	return h.client.Set(ctx, h.genKey(claims.ssid), claims.ssid, h.rfDuration).Err()
}

// CheckToken check if the token is valid and token is not logged out
func (h *Handler) CheckToken(ctx *gin.Context) (JwtClaims, error) {
	claims, err := h.ExtractToken(ctx)
	if err != nil {
		return JwtClaims{}, err
	}
	cnt, err := h.client.Exists(ctx, h.genKey(claims.ssid)).Result()
	if err != nil {
		return JwtClaims{}, err
	}
	if cnt > 0 {
		return JwtClaims{}, ErrSessionInvalid
	}
	return claims, nil
}

func (h *Handler) ssid() string {
	return uuid.New().String()
}

func (h *Handler) genKey(ssid string) string {
	return fmt.Sprintf("jwt:ssid:%s", ssid)
}
