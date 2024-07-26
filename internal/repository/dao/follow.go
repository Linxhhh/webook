package dao

import (
	"context"
	"math/rand"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type FollowDAO interface {
	InsertFollow(ctx context.Context, follower_id, followee_id int64) error
	DeleteFollow(ctx context.Context, follower_id, followee_id int64) error
	GetFollowed(ctx context.Context, follower_id, followee_id int64) (FollowRelation, error)
	GetFollowData(ctx context.Context, uid int64) (FollowData, error)
	GetFolloweeList(ctx context.Context, follower_id int64, limit, offset int) ([]FollowRelation, error)
	GetFollowerList(ctx context.Context, followee_id int64, limit, offset int) ([]FollowRelation, error)
}

type GormFollowDAO struct {
	master *gorm.DB
	slaves []*gorm.DB
}

func NewFollowDAO(m *gorm.DB, s []*gorm.DB) FollowDAO {
	return &GormFollowDAO{
		master: m,
		slaves: s,
	}
}

func (dao *GormFollowDAO) RandSalve() *gorm.DB {
	rand.Seed(time.Now().UnixNano())
    randomSlave := dao.slaves[rand.Intn(len(dao.slaves))]
    return randomSlave
}

// InsertFollow 往数据库中插入一条记录
func (dao *GormFollowDAO) InsertFollow(ctx context.Context, follower_id, followee_id int64) error {
	now := time.Now().UnixMilli()

	// 开启事务
	return dao.master.WithContext(ctx).Transaction(func(tx *gorm.DB) error {

		// upsert 语义
		err := tx.WithContext(ctx).Clauses(clause.OnConflict{
			DoUpdates: clause.Assignments(map[string]interface{}{
				"status": true,
				"utime":  now,
			}),
		}).Create(&FollowRelation{
			Follower: follower_id,
			Followee: followee_id,
			Status:   true,
			Ctime:    now,
			Utime:    now,
		}).Error
		if err != nil {
			return err
		}

		// 关注数量 +1
		err = tx.WithContext(ctx).Clauses(clause.OnConflict{
			DoUpdates: clause.Assignments(map[string]interface{}{
				"followees": gorm.Expr("`followees` + 1"),
				"utime":     now,
			}),
		}).Create(&FollowData{
			Uid:       follower_id,
			Followers: 0,
			Followees: 1,
			Ctime:     now,
			Utime:     now,
		}).Error
		if err != nil {
			return err
		}

		// 粉丝数量 +1
		return tx.WithContext(ctx).Clauses(clause.OnConflict{
			DoUpdates: clause.Assignments(map[string]interface{}{
				"followers": gorm.Expr("`followers` + 1"),
				"utime":     now,
			}),
		}).Create(&FollowData{
			Uid:       followee_id,
			Followers: 1,
			Followees: 0,
			Ctime:     now,
			Utime:     now,
		}).Error
	})
}

// DeleteFollow 软删除一条关注记录
func (dao *GormFollowDAO) DeleteFollow(ctx context.Context, follower_id, followee_id int64) error {
	now := time.Now().UnixMilli()

	// 开启事务
	return dao.master.WithContext(ctx).Transaction(func(tx *gorm.DB) error {

		// 软删除用户关注记录
		err := dao.master.WithContext(ctx).Model(&FollowRelation{}).
			Where("follower = ? AND followee = ?", follower_id, followee_id).
			Updates(map[string]any{
				"utime":  now,
				"status": false,
			}).Error
		if err != nil {
			return err
		}

		// 关注数量 -1
		err = tx.Model(&FollowData{}).
			Where("uid =?", follower_id).
			Updates(map[string]interface{}{
				"followees": gorm.Expr("`followees` - 1"),
				"utime":     now,
			}).Error
		if err != nil {
			return err
		}

		// 粉丝数量 -1
		return tx.Model(&FollowData{}).
			Where("uid =?", followee_id).
			Updates(map[string]interface{}{
				"followers": gorm.Expr("`followers` - 1"),
				"utime":     now,
			}).Error
	})
}

// GetFollowed 查询是否关注某人
func (dao *GormFollowDAO) GetFollowed(ctx context.Context, follower_id, followee_id int64) (FollowRelation, error) {
	var res FollowRelation
	err := dao.RandSalve().WithContext(ctx).Where("follower = ? AND followee = ?", follower_id, followee_id).First(&res).Error
	return res, err
}

// GetFollowData 获取关注数据（粉丝数，关注数）
func (dao *GormFollowDAO) GetFollowData(ctx context.Context, uid int64) (FollowData, error) {
	var res FollowData
	err := dao.RandSalve().WithContext(ctx).Where("uid = ?", uid).First(&res).Error
	return res, err
}

// GetFolloweeList 获取关注列表
func (dao *GormFollowDAO) GetFolloweeList(ctx context.Context, follower_id int64, limit, offset int) ([]FollowRelation, error) {
	var res []FollowRelation
	// 使用联合索引 "follower_followee"
	err := dao.RandSalve().WithContext(ctx).Select("follower, followee").
		Where("follower = ? AND status = 1", follower_id).Limit(limit).Offset(offset).Find(&res).Error
	return res, err
}

// GetFollowerList 获取粉丝列表
func (dao *GormFollowDAO) GetFollowerList(ctx context.Context, followee_id int64, limit, offset int) ([]FollowRelation, error) {
	var res []FollowRelation
	// 使用联合索引 "follower_followee"
	err := dao.RandSalve().WithContext(ctx).Select("follower, followee").
		Where("followee = ? AND status = 1", followee_id).Limit(limit).Offset(offset).Find(&res).Error
	return res, err
}

type FollowRelation struct {
	Id       int64 `gorm:"primaryKey"`
	Follower int64 `gorm:"not null;uniqueIndex:follower_followee"`// 粉丝
	Followee int64 `gorm:"not null;uniqueIndex:follower_followee"`// 博主
	Status   bool
	Ctime    int64
	Utime    int64
}

type FollowData struct {
	Id        int64 `gorm:"primaryKey"`
	Uid       int64 `gorm:"unique"`
	Followers int64 // 粉丝数量
	Followees int64 // 关注数量
	Ctime     int64
	Utime     int64
}
