package main

import (
	"github.com/spf13/viper"
)

func main() {
	initViper()
	app := InitApp()

	for _, cc := range app.consumers {
		err := cc.Start()
		if err != nil {
			panic(err)
		}
	}

	err := app.server.Run("127.0.0.1:8088")
	if err != nil {
		panic(err)
	}
}

func initViper() {
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	viper.AddConfigPath("config")
	viper.SetConfigType("yaml")
	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}
}
