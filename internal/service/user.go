package service

import (
	"context"

	"github.com/Linxhhh/webook/internal/domain"
	"github.com/Linxhhh/webook/internal/repository"
	"github.com/Linxhhh/webook/internal/repository/dao"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrDuplicateEmail = dao.ErrDuplicateEmail
)

type UserService struct {
	repo *repository.UserRepository
}

func NewUserService(repo *repository.UserRepository) *UserService {
	return &UserService{
		repo: repo,
	}
}

/* 
用户注册服务：
先进行密码加密，再调用存储层，进行数据存储
*/
func (svc *UserService) SignUp(ctx context.Context, u domain.User) error {
	
	// 密码加密
	hashPwd, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hashPwd)
	
	// 数据存储
	return svc.repo.Create(ctx, u)
}

