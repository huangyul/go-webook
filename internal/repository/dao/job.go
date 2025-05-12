package dao

import (
	"context"
	"time"

	"gorm.io/gorm"
)

type JobDAO interface {
	Preempt(ctx context.Context) (Job, error)
	Release(ctx context.Context, jobId int64) error
	UpdateTime(ctx context.Context, jobId int64) error
	UpdateNextTime(ctx context.Context, jobId int64, nextTime time.Time) error
}

var _ JobDAO = (*GORMJobDAO)(nil)

type GORMJobDAO struct {
	db *gorm.DB
}

func (dao *GORMJobDAO) UpdateNextTime(ctx context.Context, jobId int64, nextTime time.Time) error {
	now := time.Now().UnixMilli()
	return dao.db.WithContext(ctx).Model(&Job{}).Where("id = ?", jobId).Updates(map[string]any{
		"next_time": nextTime.UnixMilli(),
		"utime":     now,
	}).Error
}

func (dao *GORMJobDAO) UpdateTime(ctx context.Context, jobId int64) error {
	now := time.Now().UnixMilli()
	return dao.db.WithContext(ctx).Model(&Job{}).Where("id = ?", jobId).Updates(map[string]any{
		"utime": now,
	}).Error
}

func (dao *GORMJobDAO) Release(ctx context.Context, jobId int64) error {
	now := time.Now().UnixMilli()
	return dao.db.WithContext(ctx).Model(&Job{}).Where("id = ?", jobId).Updates(map[string]any{
		"status": jobStatusWaiting,
		"utime":  now,
	}).Error
}

func (dao *GORMJobDAO) Preempt(ctx context.Context) (Job, error) {
	now := time.Now().UnixMilli()
	var j Job
	for {
		// find job
		err := dao.db.WithContext(ctx).Model(&Job{}).Where("status = ? AND next_time < ?", jobStatusWaiting, now).First(&j).Error
		if err != nil {
			return Job{}, err
		}

		// preempt job
		res := dao.db.WithContext(ctx).Where("id = ? AND version = ?", j.Id, j.Version).Updates(map[string]any{
			"version": j.Version + 1,
			"utime":   now,
			"status":  jobStatusRunning,
		})
		if res.RowsAffected == 0 {
			continue
		}
		if res.Error != nil {
			return Job{}, res.Error
		}
		return j, nil
	}
}

type Job struct {
	Id         int64  `gorm:"primaryKey,autoIncrement"`
	Name       string `gorm:"type:varchar(128);unique"`
	Executor   string
	Expression string
	Cfg        string
	// 状态来表达，是不是可以抢占，有没有被人抢占
	Status int

	Version int

	NextTime int64 `gorm:"index"`

	Utime int64
	Ctime int64
}

const (
	jobStatusWaiting = iota
	jobStatusRunning
	jobSgtatusPaused
)
