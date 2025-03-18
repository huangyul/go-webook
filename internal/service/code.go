package service

import (
	"context"
	"fmt"
	"github.com/huangyul/go-webook/internal/repository"
	"github.com/huangyul/go-webook/internal/service/sms"
	"math/rand"
	"time"
)

var (
	ErrCodeSendTooMany   = repository.ErrCodeSendTooMany
	ErrCodeVerifyTooMany = repository.ErrCodeVerifyTooMany
)

//go:generate mockgen -source=./code.go -package=svcmock -destination=./mocks/code.mock.go
type CodeService interface {
	Send(ctx context.Context, biz, phone string) error
	Verify(ctx context.Context, biz, phone, inputCode string) (bool, error)
}

var _ CodeService = (*codeService)(nil)

type codeService struct {
	sms  sms.Service
	repo repository.CodeRepository
}

func NewCodeService(sms sms.Service, repo repository.CodeRepository) CodeService {
	return &codeService{
		sms:  sms,
		repo: repo,
	}
}

func (svc *codeService) Send(ctx context.Context, biz, phone string) error {
	code := svc.generateCode()
	err := svc.repo.Set(ctx, biz, phone, code)
	if err != nil {
		return err
	}
	tplID := "1877556"
	return svc.sms.Send(ctx, tplID, []string{code}, phone)
}

func (svc *codeService) Verify(ctx context.Context, biz, phone, inputCode string) (bool, error) {
	return svc.repo.Verify(ctx, biz, phone, inputCode)
}

func init() {
	rand.New(rand.NewSource(time.Now().UnixMilli()))
}

func (svc *codeService) generateCode() string {
	code := rand.Intn(10000000)
	return fmt.Sprintf("%06d", code)

	// a more secure way is to ues crypto/rand
	//n, _ := rand.Int(rand.Reader, big.NewInt(1000000)) // 0 - 999999
	//return fmt.Sprintf("%06d", n.Int64())             // always 6 digits
}
