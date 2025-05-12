package service

import (
	"context"
	"time"

	"github.com/huangyul/go-webook/internal/domain"
	"github.com/huangyul/go-webook/internal/repository"
)

type CronJobService interface {
	Preempt(ctx context.Context) (domain.Job, error)
	ResetNextTime(ctx context.Context, job domain.Job) error
}

type cronJobService struct {
	repo            repository.JobRepository
	refreshDuration time.Duration
}

func (svc *cronJobService) Preempt(ctx context.Context) (domain.Job, error) {
	j, err := svc.repo.Preempt(ctx)
	if err != nil {
		return domain.Job{}, err
	}
	ticker := time.NewTicker(svc.refreshDuration)
	go func() {
		for range ticker.C {
			svc.refresh()
		}
	}()
	j.CancelFunc = func() {
		ticker.Stop()
		ct, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		svc.repo.Release(ct, j.Id)
	}
	return j, nil
}

func (svc *cronJobService) ResetNextTime(ctx context.Context, job domain.Job) error {
	nextTime := job.NextTime()
	return svc.repo.UpdateNextTime(ctx, job.Id, nextTime)
}

func (svc *cronJobService) refresh() {}

func NewCronJobService(repo repository.JobRepository) CronJobService {
	return &cronJobService{
		repo:            repo,
		refreshDuration: time.Minute,
	}
}
