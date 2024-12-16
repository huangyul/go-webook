package ioc

import "github.com/huangyul/go-blog/internal/pkg/log"

func InitLogger() log.Logger {
	log := log.NewLogger()
	return log
}
