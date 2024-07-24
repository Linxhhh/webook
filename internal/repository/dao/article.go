package dao

import (
	"context"
	"errors"
	"math/rand"
	"time"

	"gorm.io/gorm"
)

var ErrIncorrectArticleorAuthor = errors.New("帖子或作者ID错误")

type ArticleDAO interface {
	Insert(ctx context.Context, article Article) (int64, error)
	Update(ctx context.Context, article Article) error
	Sync(ctx context.Context, article Article) (int64, error)
	SyncStatus(ctx context.Context, uid int64, aid int64, status uint8) error
	GetListByAuthor(ctx context.Context, uid int64, offset, limit int) ([]Article, error)
	GetById(ctx context.Context, aid int64) (Article, error)
	GetPubById(ctx context.Context, aid int64) (PublishedArticle, error) 
	GetPubList(ctx context.Context, startTime time.Time, offset, limit int) ([]PublishedArticle, error)
}

type GormArticleDAO struct {
	master *gorm.DB
	slaves []*gorm.DB
}

// NewArticleDAO 新建一个数据库存储实例
func NewArticleDAO(m *gorm.DB, s []*gorm.DB) ArticleDAO {
	return &GormArticleDAO{
		master: m,
		slaves: s,
	}
}

// RandSalve 随机获取从数据库
func (dao *GormArticleDAO) RandSalve() *gorm.DB {
	rand.Seed(time.Now().UnixNano())
    randomSlave := dao.slaves[rand.Intn(len(dao.slaves))]
    return randomSlave
}

// Insert 往数据库 Article 表中，插入一条新记录
func (dao *GormArticleDAO) Insert(ctx context.Context, article Article) (int64, error) {

	// 存储毫秒数
	now := time.Now().UnixMilli()
	article.Ctime = now
	article.Utime = now

	// 插入新记录
	err := dao.master.WithContext(ctx).Create(&article).Error
	return article.Id, err
}

// Update 更新帖子
func (dao *GormArticleDAO) Update(ctx context.Context, article Article) error {

	// Gorm 在更新时，会忽略为空值的字段
	now := time.Now().UnixMilli()
	res := dao.master.WithContext(ctx).Model(&article).
		Where("id = ? AND author_id = ?", article.Id, article.AuthorId).Updates(map[string]any{
		"title":   article.Title,
		"content": article.Content,
		"status": article.Status,
		"utime":   now,
	})
	if res.Error != nil {
		return res.Error
	}

	// 如果没有更新数据，则是 ID 或 AuthorId 错误
	if res.RowsAffected == 0 {
		return errors.New("ArticleId 或者 AuthorId 错误")
	}
	return nil
}

// Sync 使用事务，先存储制作库，再同步线上库
func (dao *GormArticleDAO) Sync(ctx context.Context, article Article) (int64, error) {

	// 使用事务
	err := dao.master.WithContext(ctx).Transaction(func(tx *gorm.DB) error {

		var err error

		// 存储制作库
		if article.Id > 0 {
			err = dao.Update(ctx, article)
		} else {
			article.Id, err = dao.Insert(ctx, article)
		}
		if err != nil {
			return err
		}

		// 同步到线上库（Insert or Update）
		err = dao.upsert(ctx, tx, article)
		return err
	})
	return article.Id, err
}

// upsert 新建帖子，或者更新帖子到线上库中
func (dao *GormArticleDAO) upsert(ctx context.Context, tx *gorm.DB, article Article) error {

	// 类型转换
	pa := PublishedArticle(article)

	// 查找帖子，如果不存在，则插入一条新记录
	err := tx.WithContext(ctx).Where("id = ?", pa.Id).First(&pa).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		pa.Ctime = time.Now().UnixMilli()
		pa.Utime = pa.Ctime
		return tx.WithContext(ctx).Create(&pa).Error
	}
	if err != nil {
		return err
	}

	// 如果帖子存在，则更新
	pa.Utime = time.Now().UnixMilli()
	return tx.WithContext(ctx).Save(&pa).Error
}

// SyncStatus 使用事务，先撤销制作库，再撤销线上库
func (dao *GormArticleDAO) SyncStatus(ctx context.Context, uid int64, aid int64, status uint8) error {

	now := time.Now().UnixMilli()

	// 使用事务
	err := dao.master.Transaction(func(tx *gorm.DB) error {

		// 撤销制作库的帖子
		result := tx.Model(&Article{}).Where("id = ? AND author_id = ?", aid, uid).Updates(map[string]any{
			"utime":  now,
			"status": status,
		})
		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected == 0 {
			return ErrIncorrectArticleorAuthor
		}

		// 撤销线上库的帖子
		return tx.Model(&PublishedArticle{}).Where("id = ?", uid).Updates(map[string]any{
			"utime":  now,
			"status": status,
		}).Error
	})
	return err
}

// GetListByAuthor 获取作者的制作库帖子列表
func (dao *GormArticleDAO) GetListByAuthor(ctx context.Context, uid int64, offset, limit int) ([]Article, error) {
	var arts []Article
	err := dao.RandSalve().WithContext(ctx).Where("author_id = ?", uid).Offset(offset).Limit(limit).Order("utime DESC").Find(&arts).Error
	return arts, err
}

// GetById 获取制作库中指定的帖子信息
func (dao *GormArticleDAO) GetById(ctx context.Context, aid int64) (Article, error) {
	var art Article
	err := dao.RandSalve().WithContext(ctx).Where("id = ?", aid).First(&art).Error
	return art, err
}

// GetPubById 获取线上库中指定的帖子信息
func (dao *GormArticleDAO) GetPubById(ctx context.Context, aid int64) (PublishedArticle, error) {
	var art PublishedArticle
	err := dao.RandSalve().WithContext(ctx).Where("id = ?", aid).First(&art).Error
	return art, err
}

// GetPubList 获取首页内容
func (dao *GormArticleDAO) GetPubList(ctx context.Context, startTime time.Time, offset, limit int) ([]PublishedArticle, error) {
	var res []PublishedArticle
	err := dao.RandSalve().WithContext(ctx).Order("utime DESC").
		Where("utime > ?", startTime.UnixMilli()).Limit(limit).Offset(offset).Find(&res).Error
	return res, err
}

// Article 制作库
type Article struct {
	Id       int64 `gorm:"primaryKey"`
	Title    string
	Content  string
	AuthorId int64
	Status   uint8
	Ctime    int64 // 创建时间
	Utime    int64 // 更新时间
}

// PublishedArticle 线上库
type PublishedArticle Article
