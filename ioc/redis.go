package ioc

import (
	"context"
	"github.com/redis/go-redis/v9"
)

func InitRedis() redis.Cmdable {
	cmd := redis.NewClient(&redis.Options{
		Addr:     "127.0.0.1:16379",
		Password: "",
		DB:       0,
	})
	if err := cmd.Ping(context.Background()).Err(); err != nil {
		panic(err)
	}
	return cmd
}
