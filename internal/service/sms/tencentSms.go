package sms

import (
	"context"
	"fmt"

	TCSms "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/sms/v20210111"
)

type TentcentService struct {
	client   *TCSms.Client
	appId    *string
	signName *string
}

func NewTentcentService(client *TCSms.Client, appId string, signName string) *TentcentService {
	return &TentcentService{
		client:   client,
		appId:    &appId,
		signName: &signName,
	}
}

/*
Send 短信发送服务：
传入模板 ID，模板参数，手机号码
*/
func (svc *TentcentService) Send(ctx context.Context, tplId string, args []string, numbers ...string) error {
	
	// 定义请求
	req := TCSms.NewSendSmsRequest()
	req.SetContext(ctx)
	req.SmsSdkAppId = svc.appId
	req.SignName = svc.signName
	req.TemplateId = &tplId
	req.TemplateParamSet = svc.toPtrSlice(args)
	req.PhoneNumberSet = svc.toPtrSlice(numbers)

	// 发送短信
	response, err := svc.client.SendSms(req)
	if err != nil {
		return err
	}

	// 检查短信状态
	for _, statusPtr := range response.Response.SendStatusSet {
		if statusPtr == nil {
			continue
		}
		status := *statusPtr
		if status.Code == nil || *(status.Code) != "Ok" {
			// 发送失败
			return fmt.Errorf("发送短信失败 code: %s, msg: %s", *status.Code, *status.Message)
		}
	}
	return nil
}

func (svc *TentcentService) toPtrSlice(strs []string) []*string {
    ptrSlice := make([]*string, len(strs))
    for i, s := range strs {
        ptrSlice[i] = &s
    }
    return ptrSlice
}
