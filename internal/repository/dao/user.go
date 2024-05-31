package dao

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
)

var (
	ErrDuplicateEmailorPhone = errors.New("邮箱或手机号码冲突")
	ErrRecordNotFound        = gorm.ErrRecordNotFound
)

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

// UserDAO 数据库存储实例
type UserDAO struct {
	db *gorm.DB
}

// NewUserDAO 新建一个数据库存储实例
func NewUserDAO(db *gorm.DB) *UserDAO {
	return &UserDAO{
		db: db,
	}
}

// Insert 往数据库 User 表中，插入一条新记录
func (dao UserDAO) Insert(ctx context.Context, u User) error {

	// 存储毫秒数
	now := time.Now().UnixMilli()
	u.CTime = now
	u.UTime = now

	// 插入新记录
	err := dao.db.WithContext(ctx).Create(&u).Error
	if mysqlErr, ok := err.(*mysql.MySQLError); ok {
		const duplicateErr uint16 = 1062
		if mysqlErr.Number == duplicateErr {
			// 用户冲突，邮箱/手机号码冲突
			return ErrDuplicateEmailorPhone
		}
	}
	return err
}

// SearchById 通过 id 查找用户
func (dao *UserDAO) SearchById(ctx context.Context, id int64) (User, error) {
	var user User
	err := dao.db.Where(id).First(&user).Error
	if err == gorm.ErrRecordNotFound {
		return user, ErrRecordNotFound
	}
	return user, err
}

// SearchByEmail 通过邮箱查找用户
func (dao *UserDAO) SearchByEmail(ctx context.Context, email string) (User, error) {
	var user User
	err := dao.db.Where("email = ?", email).First(&user).Error
	if err == gorm.ErrRecordNotFound {
		return user, ErrRecordNotFound
	}
	return user, err
}

// SearchByPhone 通过手机号码查找用户
func (dao *UserDAO) SearchByPhone(ctx context.Context, phone string) (User, error) {
	var user User
	err := dao.db.Where("phone = ?", phone).First(&user).Error
	if err == gorm.ErrRecordNotFound {
		return user, ErrRecordNotFound
	}
	return user, err
}

// Update 更新用户个人信息
func (dao *UserDAO) Update(ctx context.Context, u User) error {

	// 查找用户
	var user User
	err := dao.db.Where(u.Id).First(&user).Error
	if err == gorm.ErrRecordNotFound {
		return ErrRecordNotFound
	}
	if err != nil {
		return err
	}

	// 更新信息
	if u.NickName != "" {
		user.NickName = u.NickName
	}
	if u.Birthday != 0 {
		user.Birthday = u.Birthday
	}
	if u.Introduction != "" {
		user.Introduction = u.Introduction
	}
	now := time.Now().UnixMilli()
	user.UTime = now
	err = dao.db.Save(&user).Error
	return err
}
