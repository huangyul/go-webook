package ioc

import (
	"github.com/huangyul/go-blog/internal/service/sms"
	"github.com/huangyul/go-blog/internal/service/sms/localstorage"
	"github.com/huangyul/go-blog/internal/service/sms/ratelimit"
	"github.com/huangyul/go-blog/pkg/limiter"
	"github.com/redis/go-redis/v9"
	"time"
)

func InitSMSService(cmd redis.Cmdable) sms.SMSService {
	return ratelimit.NewSMSService(localstorage.NewSmsLocalStorageService(), limiter.NewRedisSlideWindowRedis(cmd, time.Hour*2, 2))
}
