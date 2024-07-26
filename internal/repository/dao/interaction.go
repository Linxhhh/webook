package dao

import (
	"context"
	"math/rand"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type InteractionDAO interface {
	Get(ctx context.Context, biz string, id int64) (Interaction, error)

	// 阅读模块
	IncrReadCnt(ctx context.Context, biz string, bizId int64) error
	BatchIncrReadCnt(ctx context.Context, bizs []string, bizIds []int64) error

	// 点赞模块
	GetLike(ctx context.Context, biz string, id int64, uid int64) (UserLike, error)
	InsertLike(ctx context.Context, biz string, id int64, uid int64) error
	DeleteLike(ctx context.Context, biz string, id int64, uid int64) error

	// 收藏模块
	GetCollection(ctx context.Context, biz string, id int64, uid int64) (UserCollection, error)
	GetCollectionList(ctx context.Context, biz string, uid int64) ([]UserCollection, error)
	InsertCollection(ctx context.Context, biz string, id int64, uid int64) error
	DeleteCollection(ctx context.Context, biz string, id int64, uid int64) error
}

type GORMInteractionDAO struct {
	master *gorm.DB
	slaves []*gorm.DB
}

func NewInteractionDAO(m *gorm.DB, s []*gorm.DB) InteractionDAO {
	return &GORMInteractionDAO{
		master: m,
		slaves: s,
	}
}

func (dao *GORMInteractionDAO) RandSalve() *gorm.DB {
	rand.Seed(time.Now().UnixNano())
    randomSlave := dao.slaves[rand.Intn(len(dao.slaves))]
    return randomSlave
}

// Get 获取（阅读、点赞、收藏）的数据
func (dao *GORMInteractionDAO) Get(ctx context.Context, biz string, id int64) (Interaction, error) {
	var res Interaction
	err := dao.RandSalve().WithContext(ctx).Where("biz = ? AND biz_id = ?", biz, id).First(&res).Error
	return res, err
}

// IncrReadCnt 增加阅读量
func (dao *GORMInteractionDAO) IncrReadCnt(ctx context.Context, biz string, bizId int64) error {
	now := time.Now().UnixMilli()

	// upsert 语义
	return dao.master.WithContext(ctx).Clauses(clause.OnConflict{
		DoUpdates: clause.Assignments(map[string]interface{}{
			"read_cnt": gorm.Expr("`read_cnt` + 1"),
			"utime":    now,
		}),
	}).Create(&Interaction{
		Biz:     biz,
		BizId:   bizId,
		ReadCnt: 1,
		Ctime:   now,
		Utime:   now,
	}).Error
}

// BatchIncrReadCnt 批量增加阅读量
func (dao *GORMInteractionDAO) BatchIncrReadCnt(ctx context.Context, bizs []string, bizIds []int64) error {
	return dao.master.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		txDAO := NewInteractionDAO(tx, nil)
		for i := 0; i < len(bizs); i++ {
			err := txDAO.IncrReadCnt(ctx, bizs[i], bizIds[i])
			if err != nil {
				return err
			}
		}
		return nil
	})
}

// GetLike 获取点赞信息（是否点赞）
func (dao *GORMInteractionDAO) GetLike(ctx context.Context, biz string, id int64, uid int64) (UserLike, error) {
	var res UserLike
	err := dao.RandSalve().WithContext(ctx).
		Where("biz = ? AND biz_id = ? AND uid = ? AND status = ?", biz, id, uid, 1).
		First(&res).Error
	return res, err
}

// InsertLike 插入点赞记录
func (dao *GORMInteractionDAO) InsertLike(ctx context.Context, biz string, id int64, uid int64) error {
	now := time.Now().UnixMilli()

	// 开启事务
	return dao.master.WithContext(ctx).Transaction(func(tx *gorm.DB) error {

		// 创建点赞记录（upsert 语义）
		err := tx.Clauses(clause.OnConflict{
			DoUpdates: clause.Assignments(map[string]interface{}{
				"utime":  now,
				"status": 1,
			}),
		}).Create(&UserLike{
			Uid:    uid,
			Biz:    biz,
			BizId:  id,
			Status: 1,
			Utime:  now,
			Ctime:  now,
		}).Error
		if err != nil {
			return err
		}

		// 点赞量 +1
		return tx.WithContext(ctx).Clauses(clause.OnConflict{
			DoUpdates: clause.Assignments(map[string]interface{}{
				"like_cnt": gorm.Expr("`like_cnt` + 1"),
				"utime":    now,
			}),
		}).Create(&Interaction{
			Biz:     biz,
			BizId:   id,
			LikeCnt: 1,
			Ctime:   now,
			Utime:   now,
		}).Error
	})
}

// DeleteLike 删除点赞记录（软删除）
func (dao *GORMInteractionDAO) DeleteLike(ctx context.Context, biz string, id int64, uid int64) error {
	now := time.Now().UnixMilli()

	// 开启事务
	return dao.master.WithContext(ctx).Transaction(func(tx *gorm.DB) error {

		// 软删除用户点赞记录
		err := tx.Model(&UserLike{}).
			Where("uid=? AND biz_id = ? AND biz=?", uid, id, biz).
			Updates(map[string]interface{}{
				"utime":  now,
				"status": 0,
			}).Error
		if err != nil {
			return err
		}

		// 点赞量 -1
		return tx.Model(&Interaction{}).
			Where("biz =? AND biz_id=?", biz, id).
			Updates(map[string]interface{}{
				"like_cnt": gorm.Expr("`like_cnt` - 1"),
				"utime":    now,
			}).Error
	})
}

// GetCollection 获取收藏信息（是否收藏）
func (dao *GORMInteractionDAO) GetCollection(ctx context.Context, biz string, bizId int64, uid int64) (UserCollection, error) {
	var res UserCollection
	err := dao.RandSalve().WithContext(ctx).Where("biz = ? AND biz_id = ? AND uid = ?", biz, bizId, uid).First(&res).Error
	return res, err
}

// GetCollectionList 获取收藏列表
func (dao *GORMInteractionDAO) GetCollectionList(ctx context.Context, biz string, uid int64) ([]UserCollection, error) {
	var res []UserCollection
	err := dao.RandSalve().WithContext(ctx).Where("biz = ? AND uid = ? AND status = 1", biz, uid).Find(&res).Error
	return res, err
}

// InsertCollection 插入收藏记录
func (dao *GORMInteractionDAO) InsertCollection(ctx context.Context, biz string, bizId int64, uid int64) error {
	now := time.Now().UnixMilli()
	c := UserCollection{
		Biz:    biz,
		BizId:  bizId,
		Uid:    uid,
		Status: 1,
		Ctime:  now,
		Utime:  now,
	}

	// 开启事务
	return dao.master.WithContext(ctx).Transaction(func(tx *gorm.DB) error {

		// 创建记录
		err := tx.Create(&c).Error
		if err != nil {
			return err
		}

		// 收藏量 +1（upsert 语义）
		return tx.WithContext(ctx).Clauses(clause.OnConflict{
			DoUpdates: clause.Assignments(map[string]interface{}{
				"collect_cnt": gorm.Expr("`collect_cnt` + 1"),
				"utime":       now,
			}),
		}).Create(&Interaction{
			Biz:        c.Biz,
			BizId:      c.BizId,
			CollectCnt: 1,
			Ctime:      now,
			Utime:      now,
		}).Error
	})
}

// DeleteCollection 删除收藏记录（软删除）
func (dao *GORMInteractionDAO) DeleteCollection(ctx context.Context, biz string, id int64, uid int64) error {
	now := time.Now().UnixMilli()

	// 开启事务
	return dao.master.WithContext(ctx).Transaction(func(tx *gorm.DB) error {

		// 软删除用户收藏记录
		err := tx.Model(&UserCollection{}).
			Where("uid=? AND biz_id = ? AND biz=?", uid, id, biz).
			Updates(map[string]interface{}{
				"utime":  now,
				"status": 0,
			}).Error
		if err != nil {
			return err
		}

		// 收藏量 -1
		return tx.Model(&Interaction{}).
			Where("biz =? AND biz_id=?", biz, id).
			Updates(map[string]interface{}{
				"collect_cnt": gorm.Expr("`collect_cnt` - 1"),
				"utime":       now,
			}).Error
	})
}

type UserLike struct {
	Id     int64  `gorm:"primaryKey,autoIncrement"`
	Uid    int64  `gorm:"uniqueIndex:uid_biz_type_id"`
	BizId  int64  `gorm:"uniqueIndex:uid_biz_type_id"`
	Biz    string `gorm:"type:varchar(128);uniqueIndex:uid_biz_type_id"`
	Status int
	Utime  int64
	Ctime  int64
}

type UserCollection struct {
	Id     int64  `gorm:"primaryKey,autoIncrement"`
	Uid    int64  `gorm:"uniqueIndex:uid_biz_type_id"`
	BizId  int64  `gorm:"uniqueIndex:uid_biz_type_id"`
	Biz    string `gorm:"type:varchar(128);uniqueIndex:uid_biz_type_id"`
	Status int
	Utime  int64
	Ctime  int64
}

type Interaction struct {
	Id int64 `gorm:"primaryKey,autoIncrement"`

	// <bizid, biz>
	BizId int64  `gorm:"uniqueIndex:biz_type_id"`
	Biz   string `gorm:"type:varchar(128);uniqueIndex:biz_type_id"`

	ReadCnt    int64
	LikeCnt    int64
	CollectCnt int64
	Utime      int64
	Ctime      int64
}