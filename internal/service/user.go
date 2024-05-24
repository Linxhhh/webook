package service

import (
	"context"
	"errors"

	"github.com/Linxhhh/webook/internal/domain"
	"github.com/Linxhhh/webook/internal/repository"
	"github.com/Linxhhh/webook/internal/repository/dao"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrDuplicateEmail = dao.ErrDuplicateEmail
	ErrInvalidEmailOrPassword = errors.New("邮箱或密码错误")
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
SignUp 用户注册服务：
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

/*
Login 用户登录服务：
对邮箱和密码进行校验
*/
func (svc *UserService) Login(ctx context.Context, email string, password string) (domain.User, error) {

	// 根据邮箱查找用户
	user, err := svc.repo.SearchByEmail(ctx, email)
	if err == repository.ErrUserNotFound {
		return user, ErrInvalidEmailOrPassword
	}
	if err != nil {
		return user, err
	}

	// 检查密码是否正确
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return user, ErrInvalidEmailOrPassword
	}
	
	return user, err
}

/*
edit 信息编辑服务：
直接调用存储层，进行数据更新
*/
func (us *UserService) Edit(ctx context.Context, u domain.User) error {
	return us.repo.Update(ctx, u)
}