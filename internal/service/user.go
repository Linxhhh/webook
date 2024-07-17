package service

import (
	"context"
	"errors"

	"github.com/Linxhhh/webook/internal/domain"
	"github.com/Linxhhh/webook/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrDuplicateEmailorPhone = repository.ErrDuplicateEmailorPhone
	ErrInvalidEmailOrPassword = errors.New("邮箱或密码错误")
)

/* 
这里不使用接口，是因为 <用户服务> 的可替换性不高
一般来说，可能替换的是下层的数据存储方式，即 repository.UserRepository
*/

type UserService struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) *UserService {
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
	return svc.repo.CreateByEmail(ctx, u)
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
Edit 信息编辑服务：
直接调用存储层，进行数据更新
*/
func (us *UserService) Edit(ctx context.Context, u domain.User) error {
	return us.repo.Update(ctx, u)
}

/*
Profile 信息获取服务：
直接调用存储层，然后返回信息
*/
func (us *UserService) Profile(ctx context.Context, id int64) (domain.User, error) {
	return us.repo.SearchById(ctx, id)
}

/*
FindOrCreate 查找或创建用户：
调用存储层，先查找，再创建
*/
func (us *UserService) FindOrCreate(ctx context.Context, phone string) (int64, error) {

	// 查询用户
	uid, err := us.repo.SearchByPhone(ctx, phone)
	if err == repository.ErrUserNotFound {
		// 创建用户
		err = us.repo.CreateByPhone(ctx, phone)
	}
	if err != nil {
		return -1, err
	}
	return uid, nil
}