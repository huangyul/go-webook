package sms

import (
	"context"
	"fmt"
)

var _ Service = (*LocalService)(nil)

type LocalService struct {
}

func NewLocalService() Service {
	return &LocalService{}
}

func (svc *LocalService) Send(ctx context.Context, tplId string, args []string, numbers ...string) error {
	fmt.Printf("发送短信：内容：%s，手机号：%v \n", args, numbers)
	return nil
}
