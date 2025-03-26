package authz

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

var (
	accessTokenSecret  = "secret1"
	refreshTokenSecret = "secret2"
)

var (
	ErrInvalidToken = errors.New("invalid token")
)

type Authz interface {
	GenerateToken(userId int64, userAgent string) (string, string, error)
	VerifyToken(token string) (*AccessTokenClaims, error)
	RefreshToken(token string) (string, string, error)
	SetLogout(token string) error
	CheckToken(tokenStr string) (bool, error)
	Logout(tokenStr string) error
}

var _ Authz = (*authz)(nil)

type authz struct {
	rds redis.Cmdable
}

func NewAuthz(rds redis.Cmdable) Authz {
	return &authz{rds: rds}
}

// Logout store the ssid in redis and mark is as logged out
func (a *authz) Logout(tokenStr string) error {
	var c AccessTokenClaims
	jwt.ParseWithClaims(tokenStr, &c, func(token *jwt.Token) (interface{}, error) {
		return []byte(accessTokenSecret), nil
	})
	return a.rds.Set(context.Background(), fmt.Sprintf("jwt:%s", c.Ssid), "", 24*7*time.Hour).Err()
}

// GenerateToken long and short token
func (a *authz) GenerateToken(userId int64, userAgent string) (string, string, error) {
	ssid := uuid.New().String()
	aClaims := &AccessTokenClaims{
		UserId:    userId,
		UserAgent: userAgent,
		Ssid:      ssid,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
		},
	}
	tokenStr, err := jwt.NewWithClaims(jwt.SigningMethodHS512, aClaims).SignedString([]byte(accessTokenSecret))
	if err != nil {
		return "", "", err
	}

	rClaims := &RefreshTokenClaims{
		UserId: userId,
		Ssid:   ssid,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * 7 * time.Hour)),
		},
	}
	rTokenStr, err := jwt.NewWithClaims(jwt.SigningMethodHS512, rClaims).SignedString([]byte(refreshTokenSecret))
	if err != nil {
		return "", "", err
	}
	return tokenStr, rTokenStr, nil
}

func (a *authz) VerifyToken(tokenStr string) (*AccessTokenClaims, error) {
	var claims AccessTokenClaims
	token, err := jwt.ParseWithClaims(tokenStr, &claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(accessTokenSecret), nil
	})
	if err != nil || !token.Valid {
		return nil, ErrInvalidToken
	}
	return &claims, nil
}

func (a *authz) RefreshToken(tokenStr string) (string, string, error) {
	var claims *RefreshTokenClaims
	token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(accessTokenSecret), nil
	})
	if err != nil || token.Valid {
		return "", "", ErrInvalidToken
	}
	return a.GenerateToken(claims.UserId, "")
}

func (a *authz) SetLogout(tokenStr string) error {
	c, err := a.VerifyToken(tokenStr)
	if err != nil {
		return err
	}
	return a.rds.Set(context.Background(), fmt.Sprintf("jwt:%s", c.Ssid), c.Ssid, time.Hour*30*24).Err()
}

func (a *authz) CheckToken(tokenStr string) (bool, error) {
	c, err := a.VerifyToken(tokenStr)
	if err != nil {
		return false, err
	}
	cnt, err := a.rds.Exists(context.Background(), fmt.Sprintf("jwt:%s", c.Ssid)).Result()
	if err != nil {
		return false, err
	}
	if cnt > 0 {
		return false, nil
	}
	return true, nil
}

type AccessTokenClaims struct {
	jwt.RegisteredClaims
	UserId    int64
	Ssid      string
	UserAgent string
}

type RefreshTokenClaims struct {
	jwt.RegisteredClaims
	UserId int64
	Ssid   string
}
