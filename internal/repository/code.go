package repository

import (
	"context"

	"github.com/Linxhhh/webook/internal/repository/cache"
)

var (
	ErrSendCodeTooMany = cache.ErrSendCodeTooMany
	ErrVerifyCodeFailed = cache.ErrVerifyCodeFailed
	ErrVerifyCodeTooMany = cache.ErrVerifyCodeTooMany
)

type CodeRepository struct {
	cache *cache.CodeCache
}

func NewCodeRepository(cache *cache.CodeCache) *CodeRepository {
	return &CodeRepository{
		cache: cache,
	}
}

func (repo *CodeRepository) Store(ctx context.Context, biz, phone, code string) error {
	return repo.cache.Set(ctx, biz, phone, code)
}

func (repo *CodeRepository) Verify(ctx context.Context, biz, phone, code string) error {
	return repo.cache.Verify(ctx, biz, phone, code)
}