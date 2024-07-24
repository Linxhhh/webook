package dao

import (
	"context"
	"math/rand"
	"time"

	"gorm.io/gorm"
)

// ----------------------------------------------- FeedPullEventDAO 拉模型 ----------------------------------------------------------

type FeedPullEventDAO interface {
	CreatePullEvent(ctx context.Context, event FeedPullEvent) error
	FindPullEvents(ctx context.Context, uids []int64, timestamp, limit int64) ([]FeedPullEvent, error)
	FindPullEventListWithTyp(ctx context.Context, typ string, uids []int64, timestamp, limit int64) ([]FeedPullEvent, error)
}

type feedPullEventDAO struct {
	master *gorm.DB
	slaves []*gorm.DB
}

func NewFeedPullEventDAO(master *gorm.DB, slaves []*gorm.DB) FeedPullEventDAO {
	return &feedPullEventDAO{
		master: master,
		slaves: slaves,
	}
}

func (dao *feedPullEventDAO) RandSalve() *gorm.DB {
	rand.Seed(time.Now().UnixNano())
    randomSlave := dao.slaves[rand.Intn(len(dao.slaves))]
    return randomSlave
}

func (f *feedPullEventDAO) CreatePullEvent(ctx context.Context, event FeedPullEvent) error {
	return f.master.WithContext(ctx).Create(&event).Error
}

func (f *feedPullEventDAO) FindPullEventListWithTyp(ctx context.Context, typ string, uids []int64, timestamp, limit int64) ([]FeedPullEvent, error) {
	var events []FeedPullEvent
	err := f.RandSalve().WithContext(ctx).
		Where("uid in ?", uids).
		Where("ctime < ?", timestamp).
		Where("type = ?", typ).
		Order("ctime desc").
		Limit(int(limit)).
		Find(&events).Error
	return events, err
}

func (f *feedPullEventDAO) FindPullEvents(ctx context.Context, uids []int64, timestamp, limit int64) ([]FeedPullEvent, error) {
	var events []FeedPullEvent
	err := f.RandSalve().WithContext(ctx).
		Where("uid in ?", uids).
		Where("ctime < ?", timestamp).
		Order("ctime desc").
		Limit(int(limit)).
		Find(&events).Error
	return events, err
}

type FeedPullEvent struct {
	Id      int64 `gorm:"primaryKey"`
	Uid     int64 `gorm:"index"`
	Type    string
	Ctime   int64
	Content string  // 存放一个大的 Json
}

// ----------------------------------------------- FeedPushEventDAO 推模型 ----------------------------------------------------------

type FeedPushEventDAO interface {
	CreatePushEvents(ctx context.Context, events []FeedPushEvent) error
	GetPushEvents(ctx context.Context, uid int64, timestamp, limit int64) ([]FeedPushEvent, error)
	GetPushEventsWithTyp(ctx context.Context, typ string, uid int64, timestamp, limit int64) ([]FeedPushEvent, error)
}

type feedPushEventDAO struct {
	master *gorm.DB
	slaves []*gorm.DB
}

func NewFeedPushEventDAO(master *gorm.DB, slaves []*gorm.DB) FeedPushEventDAO {
	return &feedPushEventDAO{
		master: master,
		slaves: slaves,
	}
}

func (dao *feedPushEventDAO) RandSalve() *gorm.DB {
	rand.Seed(time.Now().UnixNano())
    randomSlave := dao.slaves[rand.Intn(len(dao.slaves))]
    return randomSlave
}

func (f *feedPushEventDAO) CreatePushEvents(ctx context.Context, events []FeedPushEvent) error {
	return f.master.WithContext(ctx).Create(events).Error
}

func (f *feedPushEventDAO) GetPushEventsWithTyp(ctx context.Context, typ string, uid int64, timestamp, limit int64) ([]FeedPushEvent, error) {
	var events []FeedPushEvent
	err := f.RandSalve().WithContext(ctx).
		Where("uid = ?", uid).
		Where("ctime < ?", timestamp).
		Where("type = ?", typ).
		Order("ctime desc").
		Limit(int(limit)).
		Find(&events).Error
	return events, err
}

func (f *feedPushEventDAO) GetPushEvents(ctx context.Context, uid int64, timestamp, limit int64) ([]FeedPushEvent, error) {
	var events []FeedPushEvent
	err := f.RandSalve().WithContext(ctx).
		Where("uid = ?", uid).
		Where("ctime < ?", timestamp).
		Order("ctime desc").
		Limit(int(limit)).
		Find(&events).Error
	return events, err
}

type FeedPushEvent struct {
	Id      int64 `gorm:"primaryKey"`
	Uid     int64 `gorm:"index"`
	Type    string
	Ctime   int64
	Content string  // 存放一个大的 Json
}
