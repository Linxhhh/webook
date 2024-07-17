package cache

import (
	"context"
	_ "embed"
	"errors"
	"fmt"

	"github.com/go-redis/redis"
)

//go:embed lua/setCode.lua
var LuaSetCode string

//go:embed lua/verifyCode.lua
var LuaVerifyCode string

var (
	ErrSendCodeTooMany = errors.New("短信发送频繁")
	ErrVerifyCodeFailed = errors.New("短信验证失败")
	ErrVerifyCodeTooMany = errors.New("短信验证频繁")
)

type CodeCache interface {
	Set(ctx context.Context, biz, phone, code string) error
	Verify(ctx context.Context, biz, phone, code string) error
}

type RedisCodeCache struct {
	cmd redis.Cmdable
}

func NewCodeCache(cmd redis.Cmdable) CodeCache {
	return &RedisCodeCache{
		cmd: cmd,
	}
}

func (cc *RedisCodeCache) Key(biz, phone string) string {
	// 格式化 Key
	return fmt.Sprintf("code:%s:%s", biz, phone)
}

func (cc *RedisCodeCache) Set(ctx context.Context, biz, phone, code string) error {
	
	res, err := cc.cmd.Eval(LuaSetCode, []string{cc.Key(biz, phone)}, code).Int()
	if err != nil {
		return err
	}

	switch res {
	case 0:
		return nil
	case 1:
		return ErrSendCodeTooMany
	default:
		return errors.New("系统错误")
	}
}

func (cc *RedisCodeCache) Verify(ctx context.Context, biz, phone, code string) error {
	
	res, err := cc.cmd.Eval(LuaVerifyCode, []string{cc.Key(biz, phone)}, code).Int()
	if err != nil {
		return err
	}

	switch res {
	case 0:
		return nil
	case 1:
		return ErrVerifyCodeFailed
	case -1:
		return ErrVerifyCodeTooMany
	default:
		return errors.New("系统错误")
	}
}