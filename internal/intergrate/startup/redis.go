package startup

import (
	"context"
	"github.com/redis/go-redis/v9"
)

func InitRedis() redis.Cmdable {
	cmd := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	err := cmd.Ping(context.Background()).Err()
	if err != nil {
		panic(err)
	}

	return cmd

}
