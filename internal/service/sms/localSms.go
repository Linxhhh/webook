package sms

import (
	"context"
	"log"
)

type LocalService struct {
}

func NewLocalService() *LocalService {
	return &LocalService{}
}

func (s *LocalService) Send(ctx context.Context, tplId string, args []string, numbers ...string) error {
	log.Println("验证码是", args)
	return nil
}