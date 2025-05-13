package job

import (
	"context"
	"fmt"
	"time"

	"github.com/huangyul/go-webook/internal/domain"
	"github.com/huangyul/go-webook/internal/service"
	"golang.org/x/sync/semaphore"
)

type Executor interface {
	Name() string
	Exec(ctx context.Context, j domain.Job) error
}

type LocalExecutor struct {
	funcs map[string]func(ctx context.Context, j domain.Job) error
}

func (l *LocalExecutor) RegisterFunc(name string, fn func(ctx context.Context, j domain.Job) error) {
	l.funcs[name] = fn
}

func (l *LocalExecutor) Exec(ctx context.Context, j domain.Job) error {
	fn, ok := l.funcs[j.Executor]
	if !ok {
		return fmt.Errorf("%s not exist", j.Executor)
	}
	return fn(ctx, j)
}

func (l *LocalExecutor) Name() string {
	return "local"
}

func NewLocalExecutor() Executor {
	return &LocalExecutor{
		funcs: map[string]func(ctx context.Context, j domain.Job) error{},
	}
}

type Scheduler struct {
	dbTimeout time.Duration
	svc       service.CronJobService
	exectors  map[string]Executor
	limit     *semaphore.Weighted
}

func NewScheduler(svc service.CronJobService) *Scheduler {
	return &Scheduler{
		svc:       svc,
		dbTimeout: time.Second,
		exectors:  map[string]Executor{},
		limit:     semaphore.NewWeighted(150),
	}
}

func (s *Scheduler) RegisterExec(e Executor) {
	s.exectors[e.Name()] = e
}

func (s *Scheduler) Schedule(ctx context.Context) error {
	for {
		if ctx.Err() != nil {
			return ctx.Err()
		}

		err := s.limit.Acquire(ctx, 1)
		if err != nil {
			return err
		}

		dctx, cancel := context.WithTimeout(ctx, s.dbTimeout)
		j, err := s.svc.Preempt(dctx)
		cancel()
		if err != nil {
			continue
		}

		e, ok := s.exectors[j.Executor]
		if !ok {
			continue
		}

		go func() {
			defer s.limit.Release(1)
			err1 := e.Exec(ctx, j)
			if err1 != nil {
				return
			}
			_ = s.svc.ResetNextTime(ctx, j)
		}()
	}
}
