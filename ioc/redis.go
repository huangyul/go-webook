package ioc

import (
	"context"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
)

func InitRedis() redis.Cmdable {
	type RedisConfg struct {
		Addr     string `mapstructure:"addr"`
		Password string `mapstructure:"password"`
		Database int    `mapstructure:"db"`
	}

	var cfg RedisConfg
	if err := viper.UnmarshalKey("redis", &cfg); err != nil {
		panic(err)
	}

	cmd := redis.NewClient(&redis.Options{
		Addr:     cfg.Addr,
		Password: cfg.Password,
		DB:       cfg.Database,
	})
	if err := cmd.Ping(context.Background()).Err(); err != nil {
		panic(err)
	}
	return cmd
}
