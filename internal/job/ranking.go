package job

import (
	"context"
	rlock "github.com/gotomicro/redis-lock"
	"sync"
	"time"

	"github.com/huangyul/go-webook/internal/service"
)

type RankingJob struct {
	svc     service.RankingService
	timeout time.Duration
	client  *rlock.Client
	key     string

	localLock *sync.Mutex
	lock      *rlock.Lock
}

// Name return Job name
func (r *RankingJob) Name() string {
	return "ranking"
}

// Run
func (r *RankingJob) Run() error {

	r.localLock.Lock()
	lock := r.lock
	if lock == nil {
		// 尝试抢锁
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*4)
		defer cancel()
		lock, err := r.client.Lock(ctx, r.key, r.timeout, &rlock.FixIntervalRetry{
			Interval: time.Millisecond * 100,
			Max:      3,
		}, time.Second)
		if err != nil {
			// 抢分布式锁失败
			r.localLock.Unlock()
			return err
		}
		r.lock = lock
		r.localLock.Unlock()
		go func() {
			// 续期
			er := lock.AutoRefresh(r.timeout/2, r.timeout)
			if er != nil {
				// 续期失败
				r.localLock.Lock()
				r.lock = nil
				r.localLock.Unlock()
			}
		}()
	}

	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()

	return r.svc.TopN(ctx)
}

func (r *RankingJob) Close() error {
	r.localLock.Lock()
	lock := r.lock
	r.localLock.Unlock()
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	return lock.Unlock(ctx)
}

func NewRankingJob(svc service.RankingService, timeout time.Duration, client *rlock.Client) *RankingJob {
	return &RankingJob{
		svc:       svc,
		timeout:   timeout,
		client:    client,
		localLock: new(sync.Mutex),
	}
}
