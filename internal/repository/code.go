package repository

import (
	"context"

	"github.com/huangyul/go-blog/internal/repository/cache"
)

type CodeRepository interface {
	Set(ctx context.Context, biz, phone, code string) error
	Verify(ctx context.Context, biz, phone, code string) (bool, error)
}

var _ CodeRepository = (*codeRepository)(nil)

type codeRepository struct {
	cache cache.CodeCache
}

func NewCodeRepository(cache cache.CodeCache) CodeRepository {
	return &codeRepository{cache: cache}
}

// Set implements CodeRespotitory.
func (c *codeRepository) Set(ctx context.Context, biz, phone, code string) error {
	return c.cache.Set(ctx, biz, phone, code)
}

// Verify implements CodeRespotitory.
func (c *codeRepository) Verify(ctx context.Context, biz, phone, code string) (bool, error) {
	return c.cache.Verify(ctx, biz, phone, code)
}
