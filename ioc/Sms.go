package ioc

import "github.com/Linxhhh/webook/internal/service/sms"

func InitSmsService() sms.Service {
	return sms.NewLocalService()
}