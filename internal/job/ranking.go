package job

import (
	"context"
	"time"

	"github.com/huangyul/go-webook/internal/service"
)

type RankingJob struct {
	svc     service.RankingService
	timeout time.Duration
}

// Name return Job name
func (r *RankingJob) Name() string {
	return "ranking"
}

// Run
func (r *RankingJob) Run() error {
	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()

	return r.svc.TopN(ctx)
}

func NewRankingJob(svc service.RankingService, timeout time.Duration) *RankingJob {
	return &RankingJob{
		svc:     svc,
		timeout: timeout,
	}
}
