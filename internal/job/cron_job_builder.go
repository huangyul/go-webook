package job

import "github.com/robfig/cron/v3"

type CronJobBuilder struct{}

func NewCronJobBuilder() *CronJobBuilder {
	return &CronJobBuilder{}
}

func (c *CronJobBuilder) Build(job Job) cron.Job {
	return cronJobAdapterFunc(func() {
		job.Run()
	})
}

type cronJobAdapterFunc func()

func (c cronJobAdapterFunc) Run() {
	c()
}
