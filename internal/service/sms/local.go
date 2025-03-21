package sms

import (
	"context"
	"errors"
	"fmt"
	"github.com/huangyul/go-webook/pkg/limiter"
)

var _ Service = (*LocalService)(nil)

type LocalService struct {
	limit limiter.Limiter
}

func NewLocalService(limit limiter.Limiter) Service {
	return &LocalService{}
}

func (svc *LocalService) Send(ctx context.Context, tplId string, args []string, numbers ...string) error {
	limited, err := svc.limit.Limit(ctx, "xxx")
	if err != nil {
		return err
	}
	if limited {
		return errors.New("rate limited")
	}
	fmt.Printf("发送短信：内容：%s，手机号：%v \n", args, numbers)
	return nil
}
