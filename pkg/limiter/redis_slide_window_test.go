package limiter

import (
	"context"
	"testing"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
)

// initRedis need to start a real redis server
func initRedis() redis.Cmdable {
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:16379",
		Password: "",
		DB:       0,
	})
	if err := client.Ping(context.Background()).Err(); err != nil {
		panic(err)
	}
	return client
}

func TestRedisSlideWindowLimiter_Limit(t *testing.T) {
	client := initRedis()
	testCase := []struct {
		name      string
		before    func(t *testing.T, client redis.Cmdable)
		key       string
		wantLimit bool
		wantErr   error
	}{
		{
			name: "limit",
			before: func(t *testing.T, cmd redis.Cmdable) {
				cmd.Del(context.Background(), "test_limiter_limit")
				for _ = range 100 {
					now := time.Now().UnixMilli()
					cmd.ZAdd(context.Background(), "test_limiter_limit", redis.Z{
						Score:  float64(now),
						Member: float64(now),
					})
				}
			},
			key:       "test_limiter_limit",
			wantLimit: true,
			wantErr:   nil,
		},
		{
			name: "no limit",
			before: func(t *testing.T, cmd redis.Cmdable) {
				cmd.Del(context.Background(), "test_limiter_limit")
				for _ = range 8 {
					now := time.Now().UnixMilli()
					cmd.ZAdd(context.Background(), "test_limiter_limit", redis.Z{
						Score:  float64(now),
						Member: float64(now),
					})
				}
			},
			key:       "test_limiter_limit",
			wantLimit: false,
			wantErr:   nil,
		},
	}

	for _, tt := range testCase {
		t.Run(tt.name, func(t *testing.T) {
			l := NewRedisSlideWindowLimiter(client, time.Second, 10)
			tt.before(t, client)
			limit, err := l.Limit(context.Background(), tt.key)
			assert.Equal(t, tt.wantLimit, limit)
			assert.Equal(t, tt.wantErr, err)
		})
	}
}
