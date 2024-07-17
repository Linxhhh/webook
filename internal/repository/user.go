package repository

import (
	"context"
	"database/sql"
	"time"

	"github.com/Linxhhh/webook/internal/domain"
	"github.com/Linxhhh/webook/internal/repository/cache"
	"github.com/Linxhhh/webook/internal/repository/dao"
)

var (
	ErrDuplicateEmailorPhone = dao.ErrDuplicateEmailorPhone
	ErrUserNotFound          = dao.ErrRecordNotFound
)

type UserRepository interface {
	CreateByEmail(ctx context.Context, u domain.User) error
	CreateByPhone(ctx context.Context, phone string) error
	SearchById(ctx context.Context, id int64) (domain.User, error)
	SearchByEmail(ctx context.Context, email string) (domain.User, error)
	SearchByPhone(ctx context.Context, phone string) (int64, error)
	Update(ctx context.Context, u domain.User) error
}

type CacheUserRepository struct {
	dao   dao.UserDAO
	cache cache.UserCache
}

func NewUserRepository(dao dao.UserDAO, cache cache.UserCache) UserRepository {
	return &CacheUserRepository{
		dao:   dao,
		cache: cache,
	}
}

func (repo *CacheUserRepository) CreateByEmail(ctx context.Context, u domain.User) error {
	return repo.dao.Insert(ctx, dao.User{
		Email: sql.NullString{
			String: u.Email,
			Valid:  u.Email != "",
		},
		Password: u.Password,
	})
}

func (repo *CacheUserRepository) CreateByPhone(ctx context.Context, phone string) error {
	return repo.dao.Insert(ctx, dao.User{
		Phone: sql.NullString{
			String: phone,
			Valid:  phone != "",
		},
	})
}

func (repo *CacheUserRepository) SearchById(ctx context.Context, id int64) (domain.User, error) {

	// 查询缓存
	if user, err := repo.cache.Get(ctx, id); err == nil {
		return user, err
	}

	// 查询数据库
	user, err := repo.dao.SearchById(ctx, id)
	if err == dao.ErrRecordNotFound {
		return domain.User{}, ErrUserNotFound
	}

	u := domain.User{
		Email:        user.Email.String,
		Phone:        user.Phone.String,
		NickName:     user.NickName,
		Birthday:     time.UnixMilli(user.Birthday),
		Introduction: user.Introduction,
	}

	// 回写缓存
	go func() {
		repo.cache.Set(ctx, u)
	}()

	return u, err
}

func (repo *CacheUserRepository) SearchByEmail(ctx context.Context, email string) (domain.User, error) {
	user, err := repo.dao.SearchByEmail(ctx, email)
	if err == dao.ErrRecordNotFound {
		return domain.User{}, ErrUserNotFound
	}
	return domain.User{
		Id:       user.Id,
		Email:    user.Email.String,
		Password: user.Password,
	}, err
}

func (repo *CacheUserRepository) SearchByPhone(ctx context.Context, phone string) (int64, error) {
	user, err := repo.dao.SearchByPhone(ctx, phone)
	if err == dao.ErrRecordNotFound {
		return -1, ErrUserNotFound
	}
	return user.Id, err
}

func (repo *CacheUserRepository) Update(ctx context.Context, u domain.User) error {
	err := repo.dao.Update(ctx, dao.User{
		Id:           u.Id,
		NickName:     u.NickName,
		Birthday:     u.Birthday.UnixMilli(),
		Introduction: u.Introduction,
	})
	if err == nil {
		go func() {
			// 清除缓存
			repo.cache.Del(ctx, u.Id)
		}()
	}
	return err
}
