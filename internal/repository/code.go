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

type CodeRepository interface {
	Store(ctx context.Context, biz, phone, code string) error
	Verify(ctx context.Context, biz, phone, code string) error
}

type CacheCodeRepository struct {
	cache cache.CodeCache
}

func NewCodeRepository(cache cache.CodeCache) CodeRepository {
	return &CacheCodeRepository{
		cache: cache,
	}
}

func (repo *CacheCodeRepository) Store(ctx context.Context, biz, phone, code string) error {
	return repo.cache.Set(ctx, biz, phone, code)
}

func (repo *CacheCodeRepository) Verify(ctx context.Context, biz, phone, code string) error {
	return repo.cache.Verify(ctx, biz, phone, code)
}