package domain

import (
	"time"

	"github.com/robfig/cron/v3"
)

type Job struct {
	Id         int64
	Name       string
	Expression string
	Executor   string
	Cfg        string
	CancelFunc func()
}

func (j Job) NextTime() time.Time {
	p := cron.NewParser(cron.Second | cron.Minute | cron.Hour |
		cron.Dom | cron.Month | cron.Dow | cron.Descriptor)
	s, _ := p.Parse(j.Expression)
	return s.Next(time.Now())
}
