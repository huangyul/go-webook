package repository

import (
	"context"
	"github.com/huangyul/go-webook/internal/repository/cache"
)

var (
	ErrCodeSendTooMany   = cache.ErrCodeSendTooMany
	ErrCodeVerifyTooMany = cache.ErrCodeVerifyTooMany
)

type CodeRepository interface {
	Set(ctx context.Context, biz string, phone string, code string) error
	Verify(ctx context.Context, biz string, phone string, inputCode string) (bool, error)
}

type codeRepository struct {
	cache cache.CodeCache
}

func (repo *codeRepository) Set(ctx context.Context, biz string, phone string, code string) error {
	return repo.cache.Set(ctx, biz, phone, code)
}

func (repo *codeRepository) Verify(ctx context.Context, biz string, phone string, inputCode string) (bool, error) {
	return repo.cache.Verify(ctx, biz, phone, inputCode)
}

func NewCodeRepository(cache cache.CodeCache) CodeRepository {
	return &codeRepository{
		cache: cache,
	}
}
