package job

import (
	"context"
	"errors"

	"github.com/huangyul/go-webook/internal/domain"
)

type Executor interface {
	Name() string
	Exec(ctx context.Context, j domain.Job) error
}

type LocalFuncExecutor struct {
	funcs map[string]func(ctx context.Context, j domain.Job) error
}

func (l *LocalFuncExecutor) RegisterFunc(name string, fn func(ctx context.Context, j domain.Job) error) {
	l.funcs[name] = fn
}

func (l *LocalFuncExecutor) Exec(ctx context.Context, j domain.Job) error {
	fn, ok := l.funcs[j.Executor]
	if !ok {
		return errors.New("not exist this execurot")
	}
	return fn(ctx, j)
}

func (l *LocalFuncExecutor) Name() string {
	return "local"
}

func NewLocalFuncExecutor() Executor {
	return &LocalFuncExecutor{
		funcs: map[string]func(ctx context.Context, j domain.Job) error{},
	}
}
