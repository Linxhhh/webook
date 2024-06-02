package service

import (
	"context"
	"fmt"
	"math/rand"

	"github.com/Linxhhh/webook/internal/repository"
	"github.com/Linxhhh/webook/internal/service/sms"
)

var (
	ErrSendCodeTooMany   = repository.ErrSendCodeTooMany
	ErrVerifyCodeFailed  = repository.ErrVerifyCodeFailed
	ErrVerifyCodeTooMany = repository.ErrVerifyCodeTooMany
)

type CodeService struct {
	repo  *repository.CodeRepository
	sms   sms.Service
	tplId string
}

func NewCodeService(repo *repository.CodeRepository, sms sms.Service) *CodeService {
	return &CodeService{
		repo:  repo,
		sms:   sms,
		tplId: "1234567",
	}
}

/*
Send 验证码发送服务：
传入业务类型 biz，用户手机号码 phone
*/
func (svc *CodeService) Send(ctx context.Context, biz string, phone string) error {

	// 生成验证码
	code := svc.generateCode()

	// 存储验证码
	err := svc.repo.Store(ctx, biz, phone, code)
	if err != nil {
		return err
	}

	// 发送验证码
	err = svc.sms.Send(ctx, svc.tplId, []string{code}, phone)
	if err != nil {
		return err
	}
	return nil
}

func (svc *CodeService) generateCode() string {
	number := rand.Intn(1000000)
	return fmt.Sprintf("%6d", number)
}

/*
Verify 验证码校验服务：
传入业务类型 biz，用户手机号码 phone，用户输入的验证码 code
*/
func (svc *CodeService) Verify(ctx context.Context, biz string, phone string, code string) error {
	return svc.repo.Verify(ctx, biz, phone, code)
}
