package main

import (
	"fmt"
	"github.com/spf13/viper"
)

func main() {

	initViper()

	s := InitWebServer()

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
