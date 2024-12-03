package ioc

import (
	"github.com/huangyul/go-blog/internal/service/sms"
	"github.com/huangyul/go-blog/internal/service/sms/localstorage"
)

func InitSMSService() sms.SMSService {
	return localstorage.NewSmsLocalStorageService()
}
