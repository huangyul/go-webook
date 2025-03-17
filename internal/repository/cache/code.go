package cache

import (
	"context"
	_ "embed"
	"errors"
	"fmt"
	"github.com/redis/go-redis/v9"
)

var (
	ErrCodeSendTooMany   = errors.New("code send too many")
	ErrCodeVerifyTooMany = errors.New("code verify too many")
)

type CodeCache interface {
	Set(ctx context.Context, biz, phone, code string) error
	Verify(ctx context.Context, biz, phone, inputCode string) (bool, error)
}

//go:embed lua/set_code.lua
var luaSetCode string

//go:embed lua/verify_code.lua
var luaVerifyCode string

type codeCache struct {
	cmd redis.Cmdable
}

func (cache *codeCache) Set(ctx context.Context, biz, phone, code string) error {
	key := cache.key(biz, phone)
	res, err := cache.cmd.Eval(ctx, luaSetCode, []string{key}, code).Int()
	if err != nil {
		return err
	}
	switch res {
	case -2:
		// not expiration time
		return errors.New("not expiration time")
	case -1:
		// set too frequently
		return ErrCodeSendTooMany
	default:
		return nil
	}
}

func (cache *codeCache) Verify(ctx context.Context, biz, phone, inputCode string) (bool, error) {
	key := cache.key(biz, phone)
	res, err := cache.cmd.Eval(ctx, luaVerifyCode, []string{key}, inputCode).Int()
	if err != nil {
		return false, err
	}
	switch res {
	case -2:
		// wrong verification code
		return false, nil
	case -1:
		// verification attempts exceeded
		return false, ErrCodeVerifyTooMany
	default:
		// verification success
		return true, nil
	}
}

func NewCodeCache(cmd redis.Cmdable) CodeCache {
	return &codeCache{
		cmd: cmd,
	}
}

func (cache *codeCache) key(biz, phone string) string {
	return fmt.Sprintf("phone_code:%s:%s", biz, phone)
}
