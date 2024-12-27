package validator

import (
	"context"
	"errors"
	"github.com/huangyul/go-blog/pkg/migrator"
	"github.com/huangyul/go-blog/pkg/migrator/events"
	"golang.org/x/sync/errgroup"
	"gorm.io/gorm"
	"time"
)

type Validator[T migrator.Entity] struct {
	base      *gorm.DB
	target    *gorm.DB
	producer  events.Producer
	batchSize int
}

func (v *Validator[T]) Validate(ctx context.Context) error {
	var egg errgroup.Group
	egg.Go(func() error {
		return v.validateTargetToBase(ctx)
	})
	egg.Go(func() error {
		return v.validateBaseToTarget(ctx)
	})
	err := egg.Wait()
	return err
}

func (v *Validator[T]) validateBaseToTarget(ctx context.Context) error {
	offset := -1
	for {
		offset++
		var src T
		err := v.base.WithContext(ctx).Order("id").Offset(offset).First(&src).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil
		}
		if err != nil {
			continue
		}
		var dst T
		err = v.base.WithContext(ctx).Where("id = ?", src.ID()).First(&dst).Error
		switch {
		case errors.Is(err, gorm.ErrRecordNotFound):
			v.notify(src.ID(), events.InconsistentEventTypeTargetMissing)
		case err == nil:
			eq := src.CompareTo(dst)
			if !eq {
				v.notify(src.ID(), events.InconsistentEventTypeNEQ)
			}
		default:
			// logger
		}

	}
}
func (v *Validator[T]) validateTargetToBase(ctx context.Context) error {
	offset := -v.batchSize
	for {
		offset += v.batchSize
		var dsts []T
		err := v.target.WithContext(ctx).Order("id").Offset(offset).Find(&dsts).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil
		}
		if err != nil {
			continue
		}

		var ids []int64
		for _, dst := range dsts {
			ids = append(ids, dst.ID())
		}
		var srcs []T
		err = v.target.WithContext(ctx).Where("id in (?)", ids).Find(&srcs).Error
		if errors.Is(err, gorm.ErrRecordNotFound) || len(srcs) == 0 {
			return nil
		}
		if err != nil {
			continue
		}

		var diff []int64
		for _, src := range srcs {
			for _, dst := range dsts {
				if src.ID() != dst.ID() {
					diff = append(diff, src.ID())
				}
			}
		}
		if len(diff) > 0 {
			v.notifyBaseMissing(diff)
		}
		if len(srcs) < v.batchSize {
			return nil
		}

	}
}

func (v *Validator[T]) notify(id int64, direction string) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	err := v.producer.ProduceEvent(ctx, events.InconsistentEvent{
		ID:        id,
		Direction: direction,
	})
	if err != nil {
		// log
		return
	}
}

func (v *Validator[T]) notifyBaseMissing(ids []int64) {
	for _, id := range ids {
		v.notify(id, events.InconsistentEventTypeBaseMissing)
	}
}
