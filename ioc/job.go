package ioc

import (
	"fmt"
	rlock "github.com/gotomicro/redis-lock"
	"time"

	"github.com/huangyul/go-webook/internal/job"
	"github.com/huangyul/go-webook/internal/service"
	"github.com/robfig/cron/v3"
)

func InitRankingJob(svc service.RankingService, rockClient *rlock.Client) *job.RankingJob {
	return job.NewRankingJob(svc, time.Second*30, rockClient)
}

func InitJobs(rJob *job.RankingJob) *cron.Cron {
	jobBuild := job.NewCronJobAdapterBuilder()

	expr := cron.New(cron.WithSeconds())
	_, err := expr.AddJob("0 1-59/2 * * * *", jobBuild.Build(rJob))
	if err != nil {
		fmt.Printf("run cron error: %s", err.Error())
	}
	return expr
}
