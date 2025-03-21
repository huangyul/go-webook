package sms

import (
	"context"
	"errors"
	"github.com/huangyul/go-webook/pkg/limiter"
)

type RateLimitSMSService struct {
	svc   Service
	limit limiter.Limiter
	key   string
}

func NewRateLimitSMSService(svc Service, key string) *RateLimitSMSService {
	return &RateLimitSMSService{svc: svc, key: key}
}

func (r *RateLimitSMSService) Send(ctx context.Context, tplId string, args []string, numbers ...string) error {
	limited, err := r.limit.Limit(ctx, r.key)
	if err != nil {
		return err
	}
	if limited {
		return errors.New("limit reached")
	}
	return r.svc.Send(ctx, tplId, args, numbers...)
}
