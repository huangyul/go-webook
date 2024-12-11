package cache

import (
	"context"
	"errors"
	"fmt"
	lru "github.com/hashicorp/golang-lru"
	"github.com/huangyul/go-blog/internal/pkg/errno"
	"sync"
	"time"
)

type LocalCodeCache struct {
	cache      *lru.Cache
	l          sync.Mutex
	expiration time.Duration
}

func (l *LocalCodeCache) Set(ctx context.Context, biz, phone, code string) error {
	l.l.Lock()
	defer l.l.Unlock()
	val, ok := l.cache.Get(l.key(biz, phone))
	if !ok {
		// if code not exist, set new cache
		l.cache.Add(l.key(biz, phone), codeItem{
			code:   code,
			cnt:    3,
			expire: time.Now().Add(l.expiration),
		})
		return nil
	}
	itm, ok := val.(codeItem)
	if !ok {
		return errors.New("system error")
	}
	if itm.expire.Sub(time.Now()) > time.Minute*9 {
		return errno.ErrCodeSendTooFrequent
	}
	// resend code
	l.cache.Add(l.key(biz, phone), codeItem{
		code:   code,
		cnt:    3,
		expire: time.Now().Add(l.expiration),
	})
	return nil
}

func (l *LocalCodeCache) Verify(ctx context.Context, biz, phone, code string) (bool, error) {
	l.l.Lock()
	defer l.l.Unlock()
	val, ok := l.cache.Get(l.key(biz, phone))
	if !ok {
		return false, errno.ErrCodeNotExist
	}
	itm, ok := val.(codeItem)
	if !ok {
		return false, errors.New("system error")
	}
	if time.Now().After(itm.expire) {
		l.cache.Remove(l.key(biz, phone))
		return false, errno.ErrCodeVerifyFailed
	}
	if itm.cnt == 0 {
		l.cache.Remove(l.key(biz, phone))
		return false, errno.ErrCodeVerifyFailed
	}
	if itm.code != code {
		l.cache.Add(l.key(biz, phone), codeItem{
			code:   code,
			cnt:    itm.cnt - 1,
			expire: itm.expire,
		})
		return false, errno.ErrCodeVerifyFailed
	}
	l.cache.Remove(l.key(biz, phone))
	return true, nil
}

func (l *LocalCodeCache) key(biz, phone string) string {
	return fmt.Sprintf("phone_code:%s:%s", biz, phone)
}

type codeItem struct {
	code   string
	cnt    int
	expire time.Time
}
