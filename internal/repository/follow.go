package repository

import (
	"context"

	"github.com/Linxhhh/webook/internal/domain"
	"github.com/Linxhhh/webook/internal/repository/cache"
	"github.com/Linxhhh/webook/internal/repository/dao"
	"gorm.io/gorm"
)

type FollowRepository interface {
	Follow(ctx context.Context, follower_id, followee_id int64) error
	CancelFollow(ctx context.Context, follower_id, followee_id int64) error
	GetFollowed(ctx context.Context, follower_id, followee_id int64) (bool, error)
	GetFollowData(ctx context.Context, uid int64) (domain.FollowData, error)
	GetFolloweeList(ctx context.Context, follower_id int64, limit, offset int) ([]domain.FollowRelation, error)
	GetFollowerList(ctx context.Context, followee_id int64, limit, offset int) ([]domain.FollowRelation, error)
}

type CacheFollowRepository struct {
	dao   dao.FollowDAO
	cache cache.FollowCache
}

func NewFollowRepository(dao dao.FollowDAO, cache cache.FollowCache) FollowRepository {
	return &CacheFollowRepository{
		dao:   dao,
		cache: cache,
	}
}

func (repo *CacheFollowRepository) Follow(ctx context.Context, follower_id, followee_id int64) error {
	return repo.dao.InsertFollow(ctx, follower_id, followee_id)
}

func (repo *CacheFollowRepository) CancelFollow(ctx context.Context, follower_id, followee_id int64) error {
	return repo.dao.DeleteFollow(ctx, follower_id, followee_id)
}

func (repo *CacheFollowRepository) GetFollowed(ctx context.Context, follower_id, followee_id int64) (bool, error) {
	_, err := repo.dao.GetFollowed(ctx, follower_id, followee_id)
	switch err {
	case nil:
		return true, nil
	case dao.ErrRecordNotFound:
		return false, nil
	default:
		return false, err
	}
}

func (repo *CacheFollowRepository) GetFollowData(ctx context.Context, uid int64) (domain.FollowData, error) {

	// 查询缓存
	data, err := repo.cache.Get(ctx, uid)
	if err == nil {
		return data, err
	}

	// 查询数据库
	_data, err := repo.dao.GetFollowData(ctx, uid)
	if err != nil && err != gorm.ErrRecordNotFound {
		return domain.FollowData{}, err
	}

	// 类型转换
	data.Uid = uid
	data.Followees = _data.Followees
	data.Followers = _data.Followers

	// 回写缓存
	go func() {
		repo.cache.Set(ctx, data)
	}()

	return data, err
}

func (repo *CacheFollowRepository) GetFolloweeList(ctx context.Context, follower_id int64, limit, offset int) ([]domain.FollowRelation, error) {
	
	// 查询缓存（待补充）

	// 查询数据库
	list, err := repo.dao.GetFolloweeList(ctx, follower_id, limit, offset)
	if err != nil {
		return []domain.FollowRelation{}, err
	}

	// 类型转换
	followeeList := make([]domain.FollowRelation, 0, len(list))
	for _, elem := range list {
		followeeList = append(followeeList, domain.FollowRelation{
			Follower: elem.Follower,
			Followee: elem.Followee,
		})
	}
	return followeeList, err
}

func (repo *CacheFollowRepository) GetFollowerList(ctx context.Context, followee_id int64, limit, offset int) ([]domain.FollowRelation, error) {
	
	// 查询缓存（待补充）

	// 查询数据库
	list, err := repo.dao.GetFollowerList(ctx, followee_id, limit, offset)
	if err != nil {
		return []domain.FollowRelation{}, err
	}

	// 类型转换
	followeeList := make([]domain.FollowRelation, 0, len(list))
	for _, elem := range list {
		followeeList = append(followeeList, domain.FollowRelation{
			Follower: elem.Follower,
			Followee: elem.Followee,
		})
	}
	return followeeList, err
}
