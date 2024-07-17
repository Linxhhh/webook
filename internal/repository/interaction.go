package repository

import (
	"context"

	"github.com/Linxhhh/webook/internal/domain"
	"github.com/Linxhhh/webook/internal/repository/cache"
	"github.com/Linxhhh/webook/internal/repository/dao"
)

type InteractionRepository interface {
	IncrReadCnt(ctx context.Context, biz string, bizId int64) error
	Like(ctx context.Context, biz string, bizId int64, uid int64) error
	CancelLike(ctx context.Context, biz string, bizId int64, uid int64) error
	Collect(ctx context.Context, biz string, bizId int64, uid int64) error
	CancelCollect(ctx context.Context, biz string, bizId int64, uid int64) error
	Get(ctx context.Context, biz string, bizId int64) (domain.Interaction, error)
	GetLike(ctx context.Context, biz string, bizId int64, uid int64) (bool, error)
	GetCollection(ctx context.Context, biz string, bizId int64, uid int64) (bool, error)
	GetCollectionList(ctx context.Context, biz string, uid int64) ([]int64, error)
}

type CacheInteractionRepository struct {
	dao   dao.InteractionDAO
	cache cache.InteractionCache
}

func NewInteractionRepository(dao dao.InteractionDAO, cache cache.InteractionCache) InteractionRepository {
	return &CacheInteractionRepository{
		dao:   dao,
		cache: cache,
	}
}

// -------------------------------------------------------------------------------------------------------------------------

func (repo *CacheInteractionRepository) IncrReadCnt(ctx context.Context, biz string, bizId int64) error {
	err := repo.dao.IncrReadCnt(ctx, biz, bizId)
	if err != nil {
		return err
	}
	return repo.cache.IncrReadCnt(ctx, biz, bizId)
}

// -------------------------------------------------------------------------------------------------------------------------

func (repo *CacheInteractionRepository) Like(ctx context.Context, biz string, bizId int64, uid int64) error {
	err := repo.dao.InsertLike(ctx, biz, bizId, uid)
	if err != nil {
		return err
	}
	return repo.cache.IncrLikeCnt(ctx, biz, bizId)
}

func (repo *CacheInteractionRepository) CancelLike(ctx context.Context, biz string, bizId int64, uid int64) error {
	err := repo.dao.DeleteLike(ctx, biz, bizId, uid)
	if err != nil {
		return err
	}
	return repo.cache.DecrLikeCnt(ctx, biz, bizId)
}

// -------------------------------------------------------------------------------------------------------------------------

func (repo *CacheInteractionRepository) Collect(ctx context.Context, biz string, bizId int64, uid int64) error {
	return repo.dao.InsertCollection(ctx, biz, bizId, uid)
}

func (repo *CacheInteractionRepository) CancelCollect(ctx context.Context, biz string, bizId int64, uid int64) error {
	return repo.dao.DeleteCollection(ctx, biz, bizId, uid)
}

func (repo *CacheInteractionRepository) GetCollectionList(ctx context.Context, biz string, uid int64) ([]int64, error) {
	
	collectionList, err := repo.dao.GetCollectionList(ctx, biz, uid)
	if err != nil {
		return nil, err
	}

	var aids []int64
	for _, c := range collectionList {
		aids = append(aids, c.BizId)
	}
	return aids, err
}

// -------------------------------------------------------------------------------------------------------------------------

func (repo *CacheInteractionRepository) Get(ctx context.Context, biz string, bizId int64) (domain.Interaction, error) {
	// 查询缓存
	i, err := repo.cache.Get(ctx, biz, bizId)
	if err == nil {
		return i, err
	}

	// 查询数据库
	interaction, err := repo.dao.Get(ctx, biz, bizId)
	if err != nil {
		return domain.Interaction{}, err
	}

	// 类型转换
	i = domain.Interaction{
		ReadCnt:    interaction.ReadCnt,
		LikeCnt:    interaction.LikeCnt,
		CollectCnt: interaction.CollectCnt,
	}

	// 回写缓存
	go func() {
		repo.cache.Set(ctx, biz, bizId, i)
	}()

	return i, err
}

func (repo *CacheInteractionRepository) GetLike(ctx context.Context, biz string, bizId int64, uid int64) (bool, error) {
	_, err := repo.dao.GetLike(ctx, biz, bizId, uid)
	switch err {
	case nil:
		return true, nil
	case dao.ErrRecordNotFound:
		return false, nil
	default:
		return false, err
	}
}

func (repo *CacheInteractionRepository) GetCollection(ctx context.Context, biz string, bizId int64, uid int64) (bool, error) {
	_, err := repo.dao.GetCollection(ctx, biz, bizId, uid)
	switch err {
	case nil:
		return true, nil
	case dao.ErrRecordNotFound:
		return false, nil
	default:
		return false, err
	}
}