package service

import (
	"context"
	"fmt"
	"math/rand"

	"github.com/huangyul/go-blog/internal/repository"
	"github.com/huangyul/go-blog/internal/service/sms"
)

type CodeService interface {
	Send(ctx context.Context, biz, phone string) error
	Verify(ctx context.Context, biz, phone, code string) (bool, error)
}

var _ CodeService = (*codeService)(nil)

type codeService struct {
	smsService sms.SMSService
	repo       repository.CodeRepository
}

func NewCodeService(repo repository.CodeRepository, smsService sms.SMSService) CodeService {
	return &codeService{repo: repo, smsService: smsService}
}

func (c *codeService) Verify(ctx context.Context, biz string, phone string, code string) (bool, error) {
	return c.repo.Verify(ctx, biz, phone, code)
}

func (c *codeService) Send(ctx context.Context, biz, number string) error {
	code := c.genCode()
	err := c.repo.Set(ctx, biz, number, code)
	if err != nil {
		return err
	}
	return c.smsService.Send(ctx, "SMS_253950001", []string{code}, number)
}

func (c *codeService) genCode() string {
	code := rand.Intn(1000000)
	return fmt.Sprintf("%06d", code)
}
