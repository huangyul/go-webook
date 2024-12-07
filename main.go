package main

import (
	"fmt"
	"github.com/huangyul/go-blog/internal/pkg/log"
	"github.com/spf13/viper"
)

func main() {
	log.Init()
	defer log.Sync()
	initViper()

	s := InitWebServer()
	log.Infow("server run", "port", viper.GetInt("server.port"))

	err := s.Run(fmt.Sprintf("127.0.0.1:%d", viper.GetInt("server.port")))
	if err != nil {
		panic(err)
	}
}

func initViper() {
	viper.SetConfigName("config")
	viper.AddConfigPath("config")
	viper.SetConfigType("yaml")
	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}
}
