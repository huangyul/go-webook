package ioc

import (
	"github.com/huangyul/go-blog/internal/job"
	"github.com/robfig/cron/v3"
)

func InitJobs(j1 *job.LogJob) *cron.Cron {
	builder := job.NewCronJobBuilder()

	expr := cron.New(cron.WithSeconds())
	expr.AddJob("@every 10s", builder.Build(j1))
	return expr
}

func InitLogJob() *job.LogJob {
	return job.NewLogJob()
}
