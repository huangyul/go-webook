package localstorage

import (
	"context"
	"fmt"
	"github.com/huangyul/go-blog/internal/service/sms"
)

type SmsLocalStorageService struct {
}

func NewSmsLocalStorageService() sms.SMSService {
	return &SmsLocalStorageService{}
}

func (s *SmsLocalStorageService) Send(ctx context.Context, tplId string, args []string, number ...string) error {
	fmt.Printf("code: %v \r\n", args)
	return nil
}
