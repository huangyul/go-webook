package repository

import (
	"context"
	"time"

	"github.com/huangyul/go-webook/internal/domain"
	"github.com/huangyul/go-webook/internal/repository/dao"
)

type JobRepository interface {
	Preempt(ctx context.Context) (domain.Job, error)
	Release(ctx context.Context, jId int64) error
	UpdateTime(ctx context.Context, jId int64) error
	UpdateNextTime(ctx context.Context, jId int64, nextTime time.Time) error
}

var _ JobRepository = (*jobRespository)(nil)

type jobRespository struct {
	dao dao.JobDAO
}

// Preempt implements JobRepository.
func (repo *jobRespository) Preempt(ctx context.Context) (domain.Job, error) {
	job, err := repo.dao.Preempt(ctx)
	if err != nil {
		return domain.Job{}, err
	}
	return domain.Job{
		Id:         job.Id,
		Expression: job.Expression,
		Executor:   job.Executor,
		Name:       job.Name,
	}, nil
}

func (repo *jobRespository) Release(ctx context.Context, jId int64) error {
	return repo.dao.Release(ctx, jId)
}

func (repo *jobRespository) UpdateNextTime(ctx context.Context, jId int64, nextTime time.Time) error {
	return repo.dao.UpdateNextTime(ctx, jId, nextTime)
}

func (repo *jobRespository) UpdateTime(ctx context.Context, jId int64) error {
	return repo.dao.UpdateTime(ctx, jId)
}

func NewJobRespository(dao dao.JobDAO) JobRepository {
	return &jobRespository{dao: dao}
}
