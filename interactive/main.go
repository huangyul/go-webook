package main

import (
	"fmt"
	"github.com/spf13/viper"
)

func main() {
	initViper()

	app := InitApp()
	for _, c := range app.consumers {
		err := c.Start()
		if err != nil {
			fmt.Println(err)
		}
	}
	err := app.server.Serve()
	if err != nil {
		panic(err)
	}
}

func initViper() {
	viper.SetConfigName("config")
	viper.AddConfigPath("./interactive/config/")
	viper.SetConfigType("yaml")
	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}
}
