package failover

import (
	"context"
	"errors"
	"github.com/huangyul/go-blog/internal/service/sms"
	"sync/atomic"
)

type SMSService struct {
	svcs []sms.SMSService
	idx  uint32
}

func (s *SMSService) Send(ctx context.Context, tplId string, args []string, number ...string) error {
	length := uint32(len(args))
	idx := atomic.LoadUint32(&s.idx)

	for i := idx; i < length+idx; i++ {
		svc := s.svcs[i%length]
		err := svc.Send(ctx, tplId, args, number...)
		switch err {
		case nil:
			if i != idx {
				atomic.StoreUint32(&s.idx, i)
			}
			return nil
		case context.Canceled, context.DeadlineExceeded:
			return err
		}
	}

	return errors.New("all sms service fail")
}
