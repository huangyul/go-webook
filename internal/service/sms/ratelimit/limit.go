package ratelimit

import (
	"context"
	"errors"
	"github.com/huangyul/go-blog/internal/service/sms"
	"github.com/huangyul/go-blog/pkg/limiter"
)

var ErrLimitExceeded = errors.New("rate limit exceeded")

var _ sms.SMSService = (*SMSService)(nil)

type SMSService struct {
	svc     sms.SMSService
	limiter limiter.Limiter
	key     string
}

func NewSMSService(svc sms.SMSService, limiter limiter.Limiter) *SMSService {
	return &SMSService{svc: svc, limiter: limiter, key: "sms-limit"}
}

func (r *SMSService) Send(ctx context.Context, tplId string, args []string, number ...string) error {
	ok, err := r.limiter.Limit(ctx, r.key)
	if err != nil {
		return err
	}
	if ok {
		return ErrLimitExceeded
	}
	return r.svc.Send(ctx, tplId, args, number...)
}
