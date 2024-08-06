package dao

import (
	"context"
	"database/sql"
	"errors"
	"math/rand"
	"time"

	"github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
)

var (
	ErrDuplicateEmailorPhone = errors.New("邮箱或手机号码冲突")
	ErrRecordNotFound        = gorm.ErrRecordNotFound
)

type UserDAO interface {
	Insert(ctx context.Context, u User) (int64, error)
	SearchById(ctx context.Context, id int64) (User, error)
	SearchByEmail(ctx context.Context, email string) (User, error)
	SearchByPhone(ctx context.Context, phone string) (User, error)
	Update(ctx context.Context, u User) error
}

// UserDAO 数据库存储实例
type GormUserDAO struct {
	master *gorm.DB
	slaves []*gorm.DB
}

// NewUserDAO 新建一个数据库存储实例
func NewUserDAO(m *gorm.DB, s []*gorm.DB) UserDAO {
	return &GormUserDAO{
		master: m,
		slaves: s,
	}
}

// RandSalve 随机获取从数据库
func (dao *GormUserDAO) RandSalve() *gorm.DB {
	rand.Seed(time.Now().UnixNano())
    randomSlave := dao.slaves[rand.Intn(len(dao.slaves))]
    return randomSlave
}

// Insert 往数据库 User 表中，插入一条新记录
func (dao GormUserDAO) Insert(ctx context.Context, u User) (int64, error) {

	// 存储毫秒数
	now := time.Now().UnixMilli()
	u.CTime = now
	u.UTime = now

	// 插入新记录
	err := dao.master.WithContext(ctx).Create(&u).Error
	if mysqlErr, ok := err.(*mysql.MySQLError); ok {
		const duplicateErr uint16 = 1062
		if mysqlErr.Number == duplicateErr {
			// 用户冲突，邮箱/手机号码冲突
			return -1, ErrDuplicateEmailorPhone
		}
	}
	return u.Id, err
}

// SearchById 通过 id 查找用户
func (dao *GormUserDAO) SearchById(ctx context.Context, id int64) (User, error) {
	var user User
	err := dao.RandSalve().WithContext(ctx).Where(id).First(&user).Error
	if err == gorm.ErrRecordNotFound {
		return user, ErrRecordNotFound
	}
	return user, err
}

// SearchByEmail 通过邮箱查找用户
func (dao *GormUserDAO) SearchByEmail(ctx context.Context, email string) (User, error) {
	var user User
	err := dao.RandSalve().WithContext(ctx).Where("email = ?", email).First(&user).Error
	if err == gorm.ErrRecordNotFound {
		return user, ErrRecordNotFound
	}
	return user, err
}

// SearchByPhone 通过手机号码查找用户
func (dao *GormUserDAO) SearchByPhone(ctx context.Context, phone string) (User, error) {
	var user User
	err := dao.RandSalve().WithContext(ctx).Where("phone = ?", phone).First(&user).Error
	if err == gorm.ErrRecordNotFound {
		return user, ErrRecordNotFound
	}
	return user, err
}

// Update 更新用户个人信息
func (dao *GormUserDAO) Update(ctx context.Context, u User) error {

	var user User
	if result := dao.master.WithContext(ctx).First(&user, User{Id: u.Id}); result.Error != nil {
        return result.Error
    }
	if u.NickName != "" {
		user.NickName = u.NickName
	}
	if u.Birthday != 0 {
		user.Birthday = u.Birthday
	}
	if u.Introduction != "" {
		user.Introduction = u.Introduction
	}
	user.UTime = time.Now().UnixMilli()
	return dao.master.WithContext(ctx).Save(&user).Error
}

// User 数据库表结构
type User struct {
	Id           int64          `gorm:"primaryKey"`
	Email        sql.NullString `gorm:"unique"`
	Password     string
	Phone        sql.NullString `gorm:"unique"`
	NickName     string
	Birthday     int64
	Introduction string
	CTime        int64 // 创建时间
	UTime        int64 // 更新时间
}