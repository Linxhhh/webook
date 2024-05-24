package repository

import (
	"context"
	"time"

	"github.com/Linxhhh/webook/internal/domain"
	"github.com/Linxhhh/webook/internal/repository/dao"
)

var (
	ErrDuplicateEmail = dao.ErrDuplicateEmail
	ErrUserNotFound   = dao.ErrRecordNotFound
)

type UserRepository struct {
	dao *dao.UserDAO
}

func NewUserRepository(dao *dao.UserDAO) *UserRepository {
	return &UserRepository{
		dao: dao,
	}
}

func (repo *UserRepository) Create(ctx context.Context, u domain.User) error {
	return repo.dao.Insert(ctx, dao.User{
		Email:    u.Email,
		Password: u.Password,
	})
}

func (repo *UserRepository) SearchById(ctx context.Context, id int64) (domain.User, error) {
	user, err := repo.dao.SearchById(ctx, id)
	if err == dao.ErrRecordNotFound {
		return domain.User{}, ErrUserNotFound
	}
	return domain.User{
		Email:        user.Email,
		NickName:     user.NickName,
		Birthday:     time.UnixMilli(user.Birthday),
		Introduction: user.Introduction,
	}, err
}

func (repo *UserRepository) SearchByEmail(ctx context.Context, email string) (domain.User, error) {
	user, err := repo.dao.SearchByEmail(ctx, email)
	if err == dao.ErrRecordNotFound {
		return domain.User{}, ErrUserNotFound
	}
	return domain.User{
		Id:       user.Id,
		Email:    user.Email,
		Password: user.Password,
	}, err
}

func (repo *UserRepository) Update(ctx context.Context, u domain.User) error {
	return repo.dao.Update(ctx, dao.User{
		Id:           u.Id,
		NickName:     u.NickName,
		Birthday:     u.Birthday.UnixMilli(),
		Introduction: u.Introduction,
	})
}
