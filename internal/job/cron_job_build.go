package job

import (
	"fmt"
	"time"

	"github.com/robfig/cron"
)

type CronJobAdapterBuilder struct{}

func NewCronJobAdapterBuilder() *CronJobAdapterBuilder {
	return &CronJobAdapterBuilder{}
}

func (c *CronJobAdapterBuilder) Build(job Job) cron.Job {
	return CronJobAdapterFunc(func() {
		fmt.Printf("任务开始执行，当前时间: %s \r\n", time.Now().Format(time.DateTime))
		job.Run()
	})
}

type CronJobAdapterFunc func()

func (c CronJobAdapterFunc) Run() {
	c()
}
