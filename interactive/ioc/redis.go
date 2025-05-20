package ioc

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
)

func InitRedis() redis.Cmdable {
	type Config struct {
		Addr     string
		Password string
		DB       int
	}
	var cfg Config
	if err := viper.UnmarshalKey("redis", &cfg); err != nil {
		panic(err)
	}
	cmd := redis.NewClient(&redis.Options{
		Addr:     cfg.Addr,
		Password: cfg.Password,
		DB:       cfg.DB,
	})
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	_, err := cmd.Ping(ctx).Result()
	defer cancel()
	if err != nil {
		panic(err)
	}
	return cmd
}
